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

package loadDns

import (
	"crypto/tls"
	"github.com/kdoctor-io/kdoctor/pkg/k8s/apis/system/v1beta1"
	"github.com/kdoctor-io/kdoctor/pkg/utils/stats"
	"github.com/miekg/dns"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sync"
	"time"
)

// Max size of the buffer of result channel.
const maxResult = 1000000

type result struct {
	err      error
	duration time.Duration
	msg      *dns.Msg
}

type Work struct {
	ServerAddr string

	Msg *dns.Msg

	Protocol string

	// C is the concurrency level, the number of concurrent workers to run.
	Concurrency int

	// Timeout in seconds.
	Timeout int

	/// RequestTimeSecond request in second
	RequestTimeSecond int

	// Qps is the rate limit in queries per second.
	QPS int

	EnableLatencyMetric bool

	Logger *zap.Logger

	initOnce       sync.Once
	results        chan *result
	stopCh         chan struct{}
	qosTokenBucket chan struct{}
	startTime      metav1.Time
	report         *report
}

// Init initializes internal data-structures
func (b *Work) Init() {
	b.initOnce.Do(func() {
		b.results = make(chan *result, maxResult)
		b.stopCh = make(chan struct{}, b.Concurrency)
		b.qosTokenBucket = make(chan struct{}, b.QPS)
	})
}

// Run makes all the requests, prints the summary. It blocks until
// all work is done.
func (b *Work) Run() {
	b.Init()
	b.startTime = metav1.Now()
	b.report = newReport(b.results, b.EnableLatencyMetric)
	// Run the reporter first, it polls the result channel until it is closed.
	go func() {
		runReporter(b.report)
	}()

	// Send qps number of tokens to the channel qosTokenBucket every second to the coroutine for execution
	go func() {
		c := time.After(time.Duration(b.RequestTimeSecond) * time.Second)
		ticker := time.NewTicker(time.Second)
		// The request should be sent immediately at 0 seconds
		for i := 0; i < b.QPS; i++ {
			b.qosTokenBucket <- struct{}{}
		}
		b.Logger.Sugar().Debugf("request token channel len: %d", len(b.qosTokenBucket))
		for {
			select {
			case <-c:
				// Reach request duration stop request
				if len(b.qosTokenBucket) > 0 {
					b.Logger.Sugar().Errorf("request finish remaining number of tokens len: %d", len(b.qosTokenBucket))
					b.report.existsNotSendRequests = true
				}
				b.Stop()
				return
			case <-ticker.C:
				b.Logger.Sugar().Debugf("request token channel len: %d", len(b.qosTokenBucket))
				for i := 0; i < b.QPS; i++ {
					b.qosTokenBucket <- struct{}{}
				}
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
	total := metav1.Now().Sub(b.startTime.Time)
	// Wait until the reporter is done.
	<-b.report.done
	b.report.finalize(total)
}

func (b *Work) makeRequest(client *dns.Client, conn *dns.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	var msg *dns.Msg
	var rtt time.Duration
	var err error

	if b.Protocol == "tcp" || b.Protocol == "tcp-tls" {
		msg, rtt, err = client.Exchange(b.Msg, b.ServerAddr)

	} else {
		if conn == nil {
			conn, _ = client.Dial(b.ServerAddr)
		}
		msg, rtt, err = client.ExchangeWithConn(b.Msg, conn)
	}

	b.results <- &result{
		duration: rtt,
		err:      err,
		msg:      msg,
	}
}

func (b *Work) runWorker() {
	client := new(dns.Client)
	client.Net = b.Protocol
	client.Timeout = time.Duration(b.Timeout) * time.Millisecond
	conn, _ := client.Dial(b.ServerAddr)
	if b.Protocol == "tcp-tls" {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
		}
		client.TLSConfig = tlsConfig
	}
	wg := &sync.WaitGroup{}
	for {
		// Check if application is stopped. Do not send into a closed channel.
		select {
		case <-b.stopCh:
			wg.Wait()
			return
		case <-b.qosTokenBucket:
			wg.Add(1)
			go b.makeRequest(client, conn, wg)
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

func (b *Work) AggregateMetric() *v1beta1.DNSMetrics {
	latency := v1beta1.LatencyDistribution{}

	if b.EnableLatencyMetric {
		t, _ := stats.Mean(b.report.lats)
		latency.Mean = t

		t, _ = stats.Max(b.report.lats)
		latency.Max = t

		t, _ = stats.Min(b.report.lats)
		latency.Min = t

		t, _ = stats.Percentile(b.report.lats, 50)
		latency.P50 = t

		t, _ = stats.Percentile(b.report.lats, 90)
		latency.P90 = t

		t, _ = stats.Percentile(b.report.lats, 95)
		latency.P95 = t

		t, _ = stats.Percentile(b.report.lats, 99)
		latency.P99 = t
	} else {
		latency.Mean = b.report.totalLatencies / float32(b.report.totalCount)
	}

	metric := &v1beta1.DNSMetrics{
		StartTime:             b.startTime,
		EndTime:               metav1.NewTime(b.startTime.Add(b.report.total)),
		Duration:              b.report.total.String(),
		RequestCounts:         b.report.totalCount,
		SuccessCounts:         b.report.successCount,
		TPS:                   b.report.tps,
		Errors:                b.report.errorDist,
		Latencies:             latency,
		TargetDomain:          b.Msg.Question[0].Name,
		DNSServer:             b.ServerAddr,
		DNSMethod:             b.Protocol,
		FailedCounts:          b.report.failedCount,
		ReplyCode:             b.report.ReplyCode,
		ExistsNotSendRequests: b.report.existsNotSendRequests,
	}

	return metric
}
