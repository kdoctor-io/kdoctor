// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package loadDns_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kdoctor-io/kdoctor/pkg/loadRequest/loadDns"
	"github.com/kdoctor-io/kdoctor/pkg/logger"
	"github.com/miekg/dns"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("test dns ", Label("dns"), func() {

	It("test udp ", func() {

		dnsServer := "8.8.8.8:53"
		req := &loadDns.DnsRequestData{
			Protocol:              loadDns.RequestMethodUdp,
			DnsType:               dns.TypeA,
			TargetDomain:          "www.baidu.com",
			DnsServerAddr:         dnsServer,
			PerRequestTimeoutInMs: 5000,
			DurationInSecond:      1,
			Qps:                   10,
		}

		log := logger.NewStdoutLogger("debug", "test")
		result, e := loadDns.DnsRequest(log, req)
		Expect(e).NotTo(HaveOccurred(), "failed to execute , error=%v", e)
		Expect(int(result.FailedCounts)).To(Equal(0))
		Expect(len(result.ReplyCode)).To(Equal(1))
		Expect(result.ReplyCode).Should(HaveKey(dns.RcodeToString[dns.RcodeSuccess]))

		jsongByte, e := json.Marshal(result)
		Expect(e).NotTo(HaveOccurred(), "failed to Marshal , error=%v", e)

		var out bytes.Buffer
		e = json.Indent(&out, jsongByte, "", "\t")
		Expect(e).NotTo(HaveOccurred(), "failed to Indent , error=%v", e)
		fmt.Printf("%s\n", out.String())

	})

	It("test latency ", func() {

		dnsServer := "8.8.8.8:53"
		req := &loadDns.DnsRequestData{
			Protocol:              loadDns.RequestMethodUdp,
			DnsType:               dns.TypeA,
			TargetDomain:          "www.baidu.com",
			DnsServerAddr:         dnsServer,
			PerRequestTimeoutInMs: 5000,
			DurationInSecond:      1,
			Qps:                   10,
			EnableLatencyMetric:   true,
		}

		log := logger.NewStdoutLogger("debug", "test")
		result, e := loadDns.DnsRequest(log, req)
		Expect(e).NotTo(HaveOccurred(), "failed to execute , error=%v", e)
		Expect(int(result.FailedCounts)).To(Equal(0))
		Expect(len(result.ReplyCode)).To(Equal(1))
		Expect(result.ReplyCode).Should(HaveKey(dns.RcodeToString[dns.RcodeSuccess]))

		jsongByte, e := json.Marshal(result)
		Expect(e).NotTo(HaveOccurred(), "failed to Marshal , error=%v", e)

		var out bytes.Buffer
		e = json.Indent(&out, jsongByte, "", "\t")
		Expect(e).NotTo(HaveOccurred(), "failed to Indent , error=%v", e)
		fmt.Printf("%s\n", out.String())

	})

	It("test tcp ", func() {

		dnsServer := "223.5.5.5:53"
		req := &loadDns.DnsRequestData{
			Protocol:              loadDns.RequestMethodTcp,
			DnsType:               dns.TypeA,
			TargetDomain:          "www.baidu.com",
			DnsServerAddr:         dnsServer,
			PerRequestTimeoutInMs: 5000,
			DurationInSecond:      1,
			Qps:                   10,
		}

		log := logger.NewStdoutLogger("debug", "test")
		result, e := loadDns.DnsRequest(log, req)
		Expect(e).NotTo(HaveOccurred(), "failed to execute , error=%v", e)
		Expect(int(result.FailedCounts)).To(Equal(0))
		Expect(len(result.ReplyCode)).To(Equal(1))
		Expect(result.ReplyCode).Should(HaveKey(dns.RcodeToString[dns.RcodeSuccess]))

		jsongByte, e := json.Marshal(result)
		Expect(e).NotTo(HaveOccurred(), "failed to Marshal , error=%v", e)

		var out bytes.Buffer
		e = json.Indent(&out, jsongByte, "", "\t")
		Expect(e).NotTo(HaveOccurred(), "failed to Indent , error=%v", e)
		fmt.Printf("%s\n", out.String())
	})

	It("test bad domain ", func() {

		dnsServer := "8.8.8.8:53"
		req := &loadDns.DnsRequestData{
			Protocol:              loadDns.RequestMethodUdp,
			DnsType:               dns.TypeA,
			TargetDomain:          "www.no-existed.com",
			DnsServerAddr:         dnsServer,
			PerRequestTimeoutInMs: 5000,
			DurationInSecond:      1,
			Qps:                   10,
		}

		log := logger.NewStdoutLogger("debug", "test")
		result, e := loadDns.DnsRequest(log, req)
		Expect(e).NotTo(HaveOccurred(), "failed to execute , error=%v", e)
		Expect(result.ReplyCode["NXDOMAIN"]).To(Equal(int(result.RequestCounts)))
		Expect(len(result.ReplyCode)).To(Equal(1))
		Expect(result.ReplyCode).Should(HaveKey(dns.RcodeToString[dns.RcodeNameError]))

		jsongByte, e := json.Marshal(result)
		Expect(e).NotTo(HaveOccurred(), "failed to Marshal , error=%v", e)

		var out bytes.Buffer
		e = json.Indent(&out, jsongByte, "", "\t")
		Expect(e).NotTo(HaveOccurred(), "failed to Indent , error=%v", e)
		fmt.Printf("%s\n", out.String())

	})

	It("test aaaa ", Label("aaaa"), func() {
		dnsServer := "8.8.8.8:53"
		req := &loadDns.DnsRequestData{
			Protocol:              loadDns.RequestMethodUdp,
			DnsType:               dns.TypeAAAA,
			TargetDomain:          "wikipedia.org",
			DnsServerAddr:         dnsServer,
			PerRequestTimeoutInMs: 5000,
			DurationInSecond:      1,
			Qps:                   10,
		}

		log := logger.NewStdoutLogger("debug", "test")
		result, e := loadDns.DnsRequest(log, req)
		Expect(e).NotTo(HaveOccurred(), "failed to execute , error=%v", e)
		Expect(int(result.FailedCounts)).To(Equal(0))
		Expect(len(result.ReplyCode)).To(Equal(1))
		Expect(result.ReplyCode).Should(HaveKey(dns.RcodeToString[dns.RcodeSuccess]))

		jsongByte, e := json.Marshal(result)
		Expect(e).NotTo(HaveOccurred(), "failed to Marshal , error=%v", e)

		var out bytes.Buffer
		e = json.Indent(&out, jsongByte, "", "\t")
		Expect(e).NotTo(HaveOccurred(), "failed to Indent , error=%v", e)
		fmt.Printf("%s\n", out.String())

	})
})
