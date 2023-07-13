// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package loadDns

import (
	"fmt"
	"github.com/kdoctor-io/kdoctor/pkg/k8s/apis/system/v1beta1"
	"github.com/miekg/dns"
	"go.uber.org/zap"
	"time"

	config "github.com/kdoctor-io/kdoctor/pkg/types"
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
		Concurrency:         config.AgentConfig.Configmap.NetdnsDefaultConcurrency,
		QPS:                 reqData.Qps,
		Timeout:             reqData.PerRequestTimeoutInMs,
		Msg:                 new(dns.Msg).SetQuestion(reqData.TargetDomain, reqData.DnsType),
		Protocol:            string(reqData.Protocol),
		ServerAddr:          reqData.DnsServerAddr,
		EnableLatencyMetric: reqData.EnableLatencyMetric,
	}
	w.Init()

	// The monitoring task timed out
	if duration > 0 {
		go func() {
			time.Sleep(duration)
			w.Stop()
		}()
	}
	logger.Sugar().Infof("begin to request %v for duration %v ", w.ServerAddr, duration.String())
	w.Run()
	logger.Sugar().Infof("finish all request %v for %s ", w.report.totalCount, w.ServerAddr)
	// Collect metric reports
	metrics := w.AggregateMetric()

	logger.Sugar().Infof("result : %v ", metrics)
	return metrics, nil

}
