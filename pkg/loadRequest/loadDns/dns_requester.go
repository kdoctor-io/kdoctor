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
	"fmt"
	"github.com/kdoctor-io/kdoctor/pkg/k8s/apis/system/v1beta1"
	"github.com/kdoctor-io/kdoctor/pkg/utils/stats"
	"github.com/miekg/dns"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net"
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

func (b *Work) makeRequest(client *dns.Client, msg *dns.Msg, conn *dns.Conn, wg *sync.WaitGroup) {
	defer wg.Done()

	// Due to the time limitation on long-lived TCP connections by CoreDNS, connection reuse is not adopted for the TCP protocol.
	if RequestProtocol(b.Protocol) == RequestMethodUdp {
		err := client.ExchangeWithReuseConn(msg, conn)
		if err != nil {
			b.results <- &result{
				duration: 0,
				err:      err,
				msg:      nil,
			}
		}
	} else {
		r, rtt, err := client.Exchange(msg, b.ServerAddr)
		b.results <- &result{
			duration: rtt,
			err:      err,
			msg:      r,
		}
	}

}

func (b *Work) runWorker() {
	var conn *dns.Conn
	var err error
	client := new(dns.Client)
	client.Net = b.Protocol
	client.Timeout = time.Duration(b.Timeout) * time.Millisecond
	if RequestProtocol(b.Protocol) == RequestMethodTcpTls {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
		}
		client.TLSConfig = tlsConfig
	}

	if RequestProtocol(b.Protocol) == RequestMethodUdp {
		conn, err = b.makeConn(client)
		if err != nil {
			b.Logger.Sugar().Errorf("failed create dns conn,err=%v", err)
			return
		}
		go conn.Receiver()
	} else {
		conn = new(dns.Conn)
	}

	wg := &sync.WaitGroup{}
	for {
		// Check if application is stopped. Do not send into a closed channel.
		select {
		case <-b.stopCh:
			wg.Wait()
			// Wait for the last request to return
			time.Sleep(time.Duration(b.Timeout) * time.Millisecond)
			if RequestProtocol(b.Protocol) == RequestMethodUdp {
				if len(conn.ResponseReceiver) > 0 {
					for i := 0; i < len(conn.ResponseReceiver); i++ {
						resp := <-conn.ResponseReceiver
						e := resp.Err
						if resp.Rtt > time.Duration(b.Timeout)*time.Millisecond {
							e = fmt.Errorf("timeout for request, %d more than %d", resp.Rtt.Milliseconds(), b.Timeout)
						}
						b.results <- &result{
							duration: resp.Rtt,
							err:      e,
							msg:      resp.Msg,
						}
					}
				}
				conn.ShutDownReceiver()
				conn.Close()
			}
			return
		case <-b.qosTokenBucket:
			wg.Add(1)
			msg := new(dns.Msg)
			*msg = *b.Msg
			msg.Id = dns.Id()
			go b.makeRequest(client, msg, conn, wg)
		case resp := <-conn.ResponseReceiver:
			e := resp.Err
			if resp.Rtt > time.Duration(b.Timeout)*time.Millisecond {
				e = fmt.Errorf("timeout for request, %d more than %d", resp.Rtt.Milliseconds(), b.Timeout)
			}
			b.results <- &result{
				duration: resp.Rtt,
				err:      e,
				msg:      resp.Msg,
			}
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

func (b *Work) makeConn(c *dns.Client) (*dns.Conn, error) {
	var err error
	d := new(net.Dialer)
	conn := new(dns.Conn)
	conn.ResponseReceiver = make(chan dns.Response, b.QPS)
	conn.ShutDown = make(chan struct{})
	conn.Conn, err = d.Dial(b.Protocol, b.ServerAddr)

	// A zero value for t means Read and Write will not time out.
	_ = conn.SetWriteDeadline(time.Time{})
	_ = conn.SetReadDeadline(time.Time{})

	conn.TsigSecret, conn.TsigProvider = c.TsigSecret, c.TsigProvider
	return conn, err
}
