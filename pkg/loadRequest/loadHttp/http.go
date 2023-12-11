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

package loadHttp

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/kdoctor-io/kdoctor/pkg/k8s/apis/system/v1beta1"
)

type HttpMethod string

const (
	HttpMethodGet     = HttpMethod("GET")
	HttpMethodPost    = HttpMethod("POST")
	HttpMethodPut     = HttpMethod("PUT")
	HttpMethodDelete  = HttpMethod("DELETE")
	HttpMethodConnect = HttpMethod("CONNECT")
	HttpMethodOptions = HttpMethod("OPTIONS")
	HttpMethodPatch   = HttpMethod("PATCH")
	HttpMethodHead    = HttpMethod("HEAD")
)

type HttpRequestData struct {
	Method              HttpMethod
	Url                 string
	Qps                 int
	PerRequestTimeoutMS int
	RequestTimeSecond   int
	Header              map[string]string
	Body                []byte
	ClientCert          tls.Certificate
	CaCertPool          *x509.CertPool
	Http2               bool
	DisableKeepAlives   bool
	DisableCompression  bool
	ExpectStatusCode    *int
	EnableLatencyMetric bool
}

func HttpRequest(logger *zap.Logger, reqData *HttpRequestData) *v1beta1.HttpMetrics {
	logger.Sugar().Infof("http request=%v", reqData)
	req, _ := http.NewRequest(string(reqData.Method), reqData.Url, nil)

	duration := time.Duration(reqData.RequestTimeSecond) * time.Second
	for k, v := range reqData.Header {
		req.Header.Set(k, v)
	}

	w := &Work{
		Request:             req,
		RequestTimeSecond:   reqData.RequestTimeSecond,
		QPS:                 reqData.Qps,
		Timeout:             reqData.PerRequestTimeoutMS,
		DisableCompression:  reqData.DisableCompression,
		DisableKeepAlives:   reqData.DisableKeepAlives,
		Http2:               reqData.Http2,
		Cert:                reqData.ClientCert,
		CertPool:            reqData.CaCertPool,
		ExpectStatusCode:    reqData.ExpectStatusCode,
		RequestBody:         reqData.Body,
		EnableLatencyMetric: reqData.EnableLatencyMetric,
		Logger:              logger.Named("http-client"),
	}
	logger.Sugar().Infof("do http requests work=%v", w)
	w.Init()
	logger.Sugar().Infof("begin to request %v for duration %v ", w.Request.URL, duration.String())
	w.Run()
	logger.Sugar().Infof("finish all request %v for %s ", w.report.totalCount, w.Request.URL)
	// Collect metric reports
	metrics := w.AggregateMetric()
	return metrics
}
