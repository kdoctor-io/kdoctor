// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package v1beta1

type LatencyDistribution struct {
	// P50 is the 50th percentile request latency.
	P50 float32 `json:"p50InMs"`
	// P90 is the 90th percentile request latency.
	P90 float32 `json:"p90InMs"`
	// P95 is the 95th percentile request latency.
	P95 float32 `json:"p95InMs"`
	// P99 is the 99th percentile request latency.
	P99 float32 `json:"p99InMs"`
	// Max is the maximum observed request latency.
	Max float32 `json:"maxInMs"`
	// Min is the minimum observed request latency.
	Min float32 `json:"minInMs"`
	// Mean is the mean request latency.
	Mean float32 `json:"meanInMs"`
}

type TotalRunningLoad struct {
	AppHttpHealthyQPS int64 `json:"appHttpHealthyQPS"`
	NetReachQPS       int64 `json:"netReachQPS"`
	NetDnsQPS         int64 `json:"netDnsQPS"`
}

type SystemResource struct {
	MaxCPU    string `json:"maxCPU"`
	MeanCPU   string `json:"meanCPU"`
	MaxMemory string `json:"maxMemory"`
}
