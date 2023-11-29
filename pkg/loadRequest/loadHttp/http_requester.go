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
//
// Based on https://github.com/rakyll/hey/blob/master/requester/requester.go
//
// Changes:
// - add more metric
// - add more concurrency request

package loadHttp

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/kdoctor-io/kdoctor/pkg/k8s/apis/system/v1beta1"
	config "github.com/kdoctor-io/kdoctor/pkg/types"
	"github.com/kdoctor-io/kdoctor/pkg/utils/stats"
	"go.uber.org/zap"

	"golang.org/x/net/http2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Max size of the buffer of result channel.
const MaxResultChannelSize = 1000000

type result struct {
	err           error
	duration      time.Duration
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

	/// RequestTimeSecond request in second
	RequestTimeSecond int

	// Qps is the rate limit in queries per second.
	QPS int

	// DisableCompression is an option to disable compression in response
	DisableCompression bool

	// DisableKeepAlives is an option to prevents reuse of TCP connections between different HTTP requests
	DisableKeepAlives bool

	// DisableRedirects is an option to prevent the following of HTTP redirects
	DisableRedirects bool

	// Output represents the output type. If "csv" is provided, the
	// output will be dumped as a csv stream.
	Output string

	// ProxyAddr is the address of HTTP proxy server in the format on "host:port".
	// Optional.
	ProxyAddr *url.URL

	// ExpectStatusCode is the expect request return http code
	// Optional.
	ExpectStatusCode *int

	// EnableLatencyMetric is collect latency metric . default false
	// Optional.
	EnableLatencyMetric bool

	Logger *zap.Logger

	initOnce       sync.Once
	results        chan *result
	stopCh         chan struct{}
	qosTokenBucket chan struct{}
	start          time.Duration
	startTime      metav1.Time
	report         *report
}

// Init initializes internal data-structures
func (b *Work) Init() {
	b.initOnce.Do(func() {
		b.results = make(chan *result, MaxResultChannelSize)
		b.stopCh = make(chan struct{}, b.Concurrency)
		b.qosTokenBucket = make(chan struct{}, b.QPS)
	})
}

// Run makes all the requests, prints the summary. It blocks until
// all work is done.
func (b *Work) Run() {
	b.Init()
	b.startTime = metav1.Now()
	b.start = time.Since(b.startTime.Time)
	b.report = newReport(b.results, b.EnableLatencyMetric)
	// Run the reporter first, it polls the result channel until it is closed.
	go func() {
		runReporter(b.report)
	}()

	// Send qps number of tokens to the channel qosTokenBucket every second to the coroutine for execution
	go func() {
		// Request token counter to avoid issuing multiple tokens due to errors
		requestRound := 0

		c := time.After(time.Duration(b.RequestTimeSecond) * time.Second)
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		// The request should be sent immediately at 0 seconds
		for i := 0; i < b.QPS; i++ {
			b.qosTokenBucket <- struct{}{}
		}
		requestRound++
		b.Logger.Sugar().Debugf("send token %d times", requestRound)

		b.Logger.Sugar().Debugf("request token channel len: %d", len(b.qosTokenBucket))
		for {
			select {
			case <-c:
				b.Logger.Sugar().Debugf("reach request duration time, stop request")
				// Reach request duration stop request
				if len(b.qosTokenBucket) > 0 {
					b.Logger.Sugar().Errorf("request finish remaining number of tokens len: %d", len(b.qosTokenBucket))
					b.report.existsNotSendRequests = true
				}
				b.Stop()
				return
			case <-ticker.C:
				if requestRound >= b.RequestTimeSecond {
					b.Logger.Sugar().Debugf("All request tokens have been sent and will not be sent again.")
					continue
				}
				b.Logger.Sugar().Debugf("request token channel len: %d", len(b.qosTokenBucket))
				for i := 0; i < b.QPS; i++ {
					b.qosTokenBucket <- struct{}{}
				}
				requestRound++
				b.Logger.Sugar().Debugf("send token %d times", requestRound)
			}
		}
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
	close(b.qosTokenBucket)
	total := b.now() - b.start
	// Wait until the reporter is done.
	<-b.report.done
	b.report.finalize(total)
}

func (b *Work) makeRequest(c *http.Client, wg *sync.WaitGroup) {
	defer wg.Done()
	s := b.now()
	var size int64
	req := genRequest(b.Request, b.RequestBody)
	var t0, t1, t2, t3, t4, t5, t6 time.Time
	trace := &httptrace.ClientTrace{
		DNSStart: func(_ httptrace.DNSStartInfo) { t0 = time.Now() },
		DNSDone:  func(_ httptrace.DNSDoneInfo) { t1 = time.Now() },
		ConnectStart: func(_, _ string) {
			if t1.IsZero() {
				t1 = time.Now()
			}
		},
		ConnectDone: func(net, addr string, err error) {
			t2 = time.Now()

		},
		GotConn:              func(_ httptrace.GotConnInfo) { t3 = time.Now() },
		GotFirstResponseByte: func() { t4 = time.Now() },
		TLSHandshakeStart:    func() { t5 = time.Now() },
		TLSHandshakeDone:     func(_ tls.ConnectionState, _ error) { t6 = time.Now() },
	}
	req = req.WithContext(httptrace.WithClientTrace(context.Background(), trace))
	ctx, cancel := context.WithTimeout(req.Context(), time.Duration(b.Timeout)*time.Millisecond)
	defer cancel()
	req = req.WithContext(ctx)
	resp, err := c.Do(req)
	t := b.now()
	t7 := time.Now()
	finish := t - s
	var statusCode int
	if err == nil {
		size = resp.ContentLength
		resp.Body.Close()
		statusCode = resp.StatusCode
	} else {
		statusCode = 0
		if t0.IsZero() {
			t0 = t1
		}
		b.Logger.Sugar().Debugf("request err: %v", err)
		b.Logger.Sugar().Debugf("dns lookup duration: %d", t1.Sub(t0).Milliseconds())
		b.Logger.Sugar().Debugf("tls handshake duration: %d", t6.Sub(t5).Milliseconds())
		b.Logger.Sugar().Debugf("tcp connection duration: %d", t2.Sub(t1).Milliseconds())
		b.Logger.Sugar().Debugf("server processing duration: %d", t4.Sub(t3).Milliseconds())
		b.Logger.Sugar().Debugf("content transfer duration: %d", t7.Sub(t4).Milliseconds())
	}
	if b.ExpectStatusCode != nil {
		if statusCode != *b.ExpectStatusCode {
			if err == nil {
				err = fmt.Errorf("The %d status code returned is not the expected %d ", statusCode, *b.ExpectStatusCode)
			}
		}
	}

	b.results <- &result{
		duration:      finish,
		statusCode:    statusCode,
		err:           err,
		contentLength: size,
	}
}

func (b *Work) runWorker() {
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
	client := &http.Client{Transport: tr}
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
			client.CloseIdleConnections()
			return
		case <-b.qosTokenBucket:
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

	if b.EnableLatencyMetric {
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
	} else {
		latency.Mean = b.report.totalLatencies / float32(b.report.totalCount)
	}

	var errNum int64
	for _, v := range b.report.errorDist {
		errNum += int64(v)
	}

	metric := &v1beta1.HttpMetrics{
		StartTime:             b.startTime,
		EndTime:               metav1.NewTime(b.startTime.Add(b.report.total)),
		Duration:              b.report.total.String(),
		RequestCounts:         b.report.totalCount,
		SuccessCounts:         b.report.totalCount - errNum,
		TPS:                   b.report.tps,
		Errors:                b.report.errorDist,
		Latencies:             latency,
		TotalDataSize:         strconv.Itoa(int(b.report.sizeTotal)) + " byte",
		StatusCodes:           b.report.statusCodes,
		ExistsNotSendRequests: b.report.existsNotSendRequests,
	}

	return metric
}

func genRequest(r *http.Request, body []byte) *http.Request {
	// shallow copy of the struct
	r2, _ := http.NewRequest(r.Method, r.URL.String(), nil)
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
