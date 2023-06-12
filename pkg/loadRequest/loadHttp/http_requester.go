// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

// this is copied from Google project, and make some code modification
// Copyright 2014 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package loadHttp

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/kdoctor-io/kdoctor/pkg/k8s/apis/system/v1beta1"
	"golang.org/x/net/http2"
	"io"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"math"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"strconv"
	"sync"
	"time"

	config "github.com/kdoctor-io/kdoctor/pkg/types"
	"github.com/kdoctor-io/kdoctor/pkg/utils/stats"
)

// Max size of the buffer of result channel.
const MaxResultChannelSize = 1000000

type result struct {
	err           error
	duration      time.Duration
	connDuration  time.Duration // connection setup(DNS lookup + Dial up) duration
	dnsDuration   time.Duration // dns lookup duration
	reqDuration   time.Duration // request "write" duration
	resDuration   time.Duration // response "read" duration
	statusCode    int
	contentLength int64
}

type Metrics struct {
	StartTime time.Time `json:"start"`
	EndTime   time.Time `json:"end"`
	Duration  string    `json:"duration"`
	// total requests times
	Requests int64 `json:"requestCount"`
	// success times
	Success int64 `json:"successCount"`
	// transactions Per Second
	TPS float64 `json:"tps"`
	// request data size
	TotalDataSize string              `json:"total_request_data"`
	Latencies     latencyDistribution `json:"latencies"`
	StatusCodes   map[int]int         `json:"status_codes"`
	Errors        map[string]int      `json:"errors"`
}

type latencyDistribution struct {
	// P50 is the 50th percentile request latency.
	P50 float32 `json:"P50_inMs"`
	// P90 is the 90th percentile request latency.
	P90 float32 `json:"P90_inMs"`
	// P95 is the 95th percentile request latency.
	P95 float32 `json:"P95_inMs"`
	// P99 is the 99th percentile request latency.
	P99 float32 `json:"P99_inMs"`
	// Max is the maximum observed request latency.
	Max float32 `json:"Max_inMx"`
	// Min is the minimum observed request latency.
	Min float32 `json:"Min_inMs"`
	// Mean is the mean request latency.
	Mean float32 `json:"Mean_inMs"`
}

type Work struct {
	// Request is the request to be made.
	Request *http.Request

	RequestBody []byte

	Cert tls.Certificate

	CertPool *x509.CertPool
	// RequestFunc is a function to generate requests. If it is nil, then
	// Request and RequestData are cloned for each request.
	RequestFunc func() *http.Request

	// Concurrency is the concurrency level, the number of concurrent workers to run.
	Concurrency int

	// Http2 is an option to make HTTP/2 requests
	Http2 bool

	// Timeout in seconds.
	Timeout int

	// Qps is the rate limit in queries per second.
	QPS int

	// DisableCompression is an option to disable compression in response
	DisableCompression bool

	// DisableKeepAlives is an option to prevents re-use of TCP connections between different HTTP requests
	DisableKeepAlives bool

	// DisableRedirects is an option to prevent the following of HTTP redirects
	DisableRedirects bool

	// Output represents the output type. If "csv" is provided, the
	// output will be dumped as a csv stream.
	Output string

	// ProxyAddr is the address of HTTP proxy server in the format on "host:port".
	// Optional.
	ProxyAddr *url.URL

	initOnce  sync.Once
	results   chan *result
	stopCh    chan struct{}
	start     time.Duration
	startTime metav1.Time
	report    *report
}

// Init initializes internal data-structures
func (b *Work) Init() {
	b.initOnce.Do(func() {
		b.results = make(chan *result, MaxResultChannelSize)
		b.stopCh = make(chan struct{}, b.Concurrency)
	})
}

// Run makes all the requests, prints the summary. It blocks until
// all work is done.
func (b *Work) Run() {
	b.Init()
	b.startTime = metav1.Now()
	b.start = time.Since(b.startTime.Time)
	b.report = newReport(b.results, math.MaxInt32)
	// Run the reporter first, it polls the result channel until it is closed.
	go func() {
		runReporter(b.report)
	}()

	b.runWorkers()
	b.Finish()
}

func (b *Work) Stop() {
	// Send stop signal so that workers can stop gracefully.
	for i := 0; i < b.Concurrency; i++ {
		b.stopCh <- struct{}{}
	}
}

func (b *Work) Finish() {
	close(b.results)
	total := b.now() - b.start
	// Wait until the reporter is done.
	<-b.report.done
	b.report.finalize(total)
}

func (b *Work) makeRequest(c *http.Client, wg *sync.WaitGroup) {
	defer wg.Done()
	s := b.now()
	var size int64
	var dnsStart, connStart, resStart time.Duration
	var dnsDuration, connDuration, resDuration, reqDuration time.Duration
	var req *http.Request
	if b.RequestFunc != nil {
		req = b.RequestFunc()
	} else {
		req = cloneRequest(b.Request, b.RequestBody)
	}
	trace := &httptrace.ClientTrace{
		DNSStart: func(info httptrace.DNSStartInfo) {
			dnsStart = b.now()
		},
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			dnsDuration = b.now() - dnsStart
		},
		GetConn: func(h string) {
			connStart = b.now()
		},
		GotConn: func(connInfo httptrace.GotConnInfo) {
			if !connInfo.Reused {
				connDuration = b.now() - connStart
			}
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	ctx, cancel := context.WithTimeout(req.Context(), time.Duration(b.Timeout)*time.Millisecond)
	defer cancel()
	req = req.WithContext(ctx)
	resp, err := c.Do(req)
	t := b.now()
	resDuration = t - resStart
	finish := t - s
	var statusCode int
	if err == nil {
		size = resp.ContentLength
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
		statusCode = resp.StatusCode
	} else {
		statusCode = 0
	}

	b.results <- &result{
		duration:      finish,
		statusCode:    statusCode,
		err:           err,
		contentLength: size,
		connDuration:  connDuration,
		dnsDuration:   dnsDuration,
		reqDuration:   reqDuration,
		resDuration:   resDuration,
	}
}

func (b *Work) runWorker() {
	var ticker *time.Ticker
	if b.QPS > 0 {
		ticker = time.NewTicker(time.Duration(1e6*b.Concurrency/(b.QPS)) * time.Microsecond)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates:       []tls.Certificate{b.Cert},
			RootCAs:            b.CertPool,
			InsecureSkipVerify: true,
			ServerName:         b.Request.Host,
		},
		MaxIdleConnsPerHost: config.AgentConfig.Configmap.NethttpDefaultMaxIdleConnsPerHost,
		DisableCompression:  b.DisableCompression,
		DisableKeepAlives:   b.DisableKeepAlives,
		Proxy:               http.ProxyURL(b.ProxyAddr),
	}

	// verify ca
	if b.CertPool != nil {
		tr.TLSClientConfig.InsecureSkipVerify = false
	}

	if b.Http2 {
		_ = http2.ConfigureTransport(tr)
	} else {
		tr.TLSNextProto = make(map[string]func(string, *tls.Conn) http.RoundTripper)
	}
	// Each goroutine uses the same HTTP Client instance
	client := &http.Client{Transport: tr, Timeout: time.Duration(b.Timeout) * time.Second}
	if b.DisableRedirects {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	wg := &sync.WaitGroup{}
	for {
		// Check if application is stopped. Do not send into a closed channel.
		select {
		case <-b.stopCh:
			wg.Wait()
			return
		default:
			if b.QPS > 0 {
				<-ticker.C
			}
			wg.Add(1)
			go b.makeRequest(client, wg)
		}

	}
}

func (b *Work) runWorkers() {
	var wg sync.WaitGroup
	wg.Add(b.Concurrency)
	for i := 0; i < b.Concurrency; i++ {
		go func() {
			b.runWorker()
			wg.Done()
		}()
	}
	wg.Wait()

}

// Returns the time since the start of the task
func (b *Work) now() time.Duration { return time.Since(b.startTime.Time) }

// AggregateMetric Aggregate metric information and return
func (b *Work) AggregateMetric() *v1beta1.HttpMetrics {
	latency := v1beta1.LatencyDistribution{}

	t, _ := stats.Mean(b.report.latencies)
	latency.Mean = t

	t, _ = stats.Max(b.report.latencies)
	latency.Max = t

	t, _ = stats.Min(b.report.latencies)
	latency.Min = t

	t, _ = stats.Percentile(b.report.latencies, 50)
	latency.P50 = t

	t, _ = stats.Percentile(b.report.latencies, 90)
	latency.P90 = t

	t, _ = stats.Percentile(b.report.latencies, 95)
	latency.P95 = t

	t, _ = stats.Percentile(b.report.latencies, 99)
	latency.P99 = t

	var errNum int64
	for _, v := range b.report.errorDist {
		errNum += int64(v)
	}

	metric := &v1beta1.HttpMetrics{
		StartTime:     b.startTime,
		EndTime:       metav1.NewTime(b.startTime.Add(b.report.total)),
		Duration:      b.report.total.String(),
		RequestCounts: b.report.totalCount,
		SuccessCounts: b.report.totalCount - errNum,
		TPS:           b.report.tps,
		Errors:        b.report.errorDist,
		Latencies:     latency,
		TotalDataSize: strconv.Itoa(int(b.report.sizeTotal)) + " byte",
		StatusCodes:   b.report.statusCodes,
	}

	return metric
}

// cloneRequest returns a clone of the provided *http.Request.
// The clone is a shallow copy of the struct and its Header map.
func cloneRequest(r *http.Request, body []byte) *http.Request {
	// shallow copy of the struct
	r2 := new(http.Request)
	*r2 = *r
	// deep copy of the Header
	r2.Header = make(http.Header, len(r.Header))
	for k, s := range r.Header {
		r2.Header[k] = append([]string(nil), s...)
	}
	if len(body) > 0 {
		r2.Body = io.NopCloser(bytes.NewReader(body))

	}
	return r2
}
