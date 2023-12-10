// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0
//
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

// Based on https://github.com/rakyll/hey/blob/master/hey.go
//
// Changes:
// - remove param that we don't use

package loadDns

import (
	"fmt"
	"github.com/kdoctor-io/kdoctor/pkg/k8s/apis/system/v1beta1"
	"github.com/miekg/dns"
	"go.uber.org/zap"
	"time"
)

type RequestProtocol string

const (
	RequestMethodUdp    = RequestProtocol("udp")
	RequestMethodTcp    = RequestProtocol("tcp")
	RequestMethodTcpTls = RequestProtocol("tcp-tls")

	DefaultDnsConfPath = "/etc/resolv.conf"
)

type DnsRequestData struct {
	Protocol RequestProtocol
	// dns.TypeA or dns.TypeAAAA
	DnsType uint16
	// must be full domain
	TargetDomain string
	// empty, or specified to be format "2.2.2.2:53"
	DnsServerAddr         string
	PerRequestTimeoutInMs int
	Qps                   int
	DurationInSecond      int
	EnableLatencyMetric   bool
}

func DnsRequest(logger *zap.Logger, reqData *DnsRequestData) (result *v1beta1.DNSMetrics, err error) {

	logger.Sugar().Infof("dns ServerAddress=%v, request=%v, ", reqData.DnsServerAddr, reqData)

	if _, ok := dns.IsDomainName(reqData.TargetDomain); !ok {
		return nil, fmt.Errorf("invalid domain name: %v", reqData.TargetDomain)
	}
	// if not fqdn, the dns library will report error, so convert the format
	if !dns.IsFqdn(reqData.TargetDomain) {
		reqData.TargetDomain = dns.Fqdn(reqData.TargetDomain)
		logger.Sugar().Debugf("convert target domain to fqdn %v", reqData.TargetDomain)
	}
	duration := time.Duration(reqData.DurationInSecond) * time.Second

	w := &Work{
		RequestTimeSecond:   reqData.DurationInSecond,
		QPS:                 reqData.Qps,
		Timeout:             reqData.PerRequestTimeoutInMs,
		Msg:                 new(dns.Msg).SetQuestion(reqData.TargetDomain, reqData.DnsType),
		Protocol:            string(reqData.Protocol),
		ServerAddr:          reqData.DnsServerAddr,
		EnableLatencyMetric: reqData.EnableLatencyMetric,
		Logger:              logger.Named("dns-client"),
	}
	w.Init()
	logger.Sugar().Infof("begin to request %v for duration %v ", w.ServerAddr, duration.String())
	w.Run()
	logger.Sugar().Infof("finish all request %v for %s ", w.report.totalCount, w.ServerAddr)
	// Collect metric reports
	metrics := w.AggregateMetric()

	logger.Sugar().Infof("result : %v ", metrics)
	return metrics, nil

}
