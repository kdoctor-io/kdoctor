// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package loadHttp_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kdoctor-io/kdoctor/pkg/loadRequest/loadHttp"
	"github.com/kdoctor-io/kdoctor/pkg/logger"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("test http ", Label("http"), func() {
	header := make(map[string]string, 0)
	header["test"] = "test"

	It("test no latency  ", func() {

		req := &loadHttp.HttpRequestData{
			Method:              "GET",
			Url:                 "https://github.com",
			PerRequestTimeoutMS: 10000,
			RequestTimeSecond:   10,
			Qps:                 10,
			Header:              header,
		}
		log := logger.NewStdoutLogger("debug", "test")
		result := loadHttp.HttpRequest(log, req)

		jsongByte, e := json.Marshal(result)
		Expect(e).NotTo(HaveOccurred(), "failed to Marshal , error=%v", e)

		var out bytes.Buffer
		e = json.Indent(&out, jsongByte, "", "\t")
		Expect(e).NotTo(HaveOccurred(), "failed to Indent , error=%v", e)
		fmt.Printf("%s\n", out.String())

		Expect(len(result.Errors)).To(Equal(0))

	})

	It("test latency ", func() {

		req := &loadHttp.HttpRequestData{
			Method:              "GET",
			Url:                 "https://github.com",
			PerRequestTimeoutMS: 10000,
			RequestTimeSecond:   10,
			EnableLatencyMetric: true,
			Qps:                 10,
			Header:              header,
		}
		log := logger.NewStdoutLogger("debug", "test")
		result := loadHttp.HttpRequest(log, req)
		jsongByte, e := json.Marshal(result)
		Expect(e).NotTo(HaveOccurred(), "failed to Marshal , error=%v", e)

		var out bytes.Buffer
		e = json.Indent(&out, jsongByte, "", "\t")
		Expect(e).NotTo(HaveOccurred(), "failed to Indent , error=%v", e)
		fmt.Printf("%s\n", out.String())

		Expect(len(result.Errors)).To(Equal(0))

	})
})
