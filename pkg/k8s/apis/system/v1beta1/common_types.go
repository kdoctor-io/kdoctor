// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package v1beta1

type LatencyDistribution struct {
	// P50 is the 50th percentile request latency.
	P50 float32 `json:"P50InMs"`
	// P90 is the 90th percentile request latency.
	P90 float32 `json:"P90InMs"`
	// P95 is the 95th percentile request latency.
	P95 float32 `json:"P95InMs"`
	// P99 is the 99th percentile request latency.
	P99 float32 `json:"P99InMs"`
	// Max is the maximum observed request latency.
	Max float32 `json:"MaxInMs"`
	// Min is the minimum observed request latency.
	Min float32 `json:"MinInMs"`
	// Mean is the mean request latency.
	Mean float32 `json:"MeanInMs"`
}

type TotalRunningLoad struct {
	AppHttpHealthyQPS int64 `json:"AppHttpHealthyQPS"`
	NetReachQPS       int64 `json:"NetReachQPS"`
	NetDnsQPS         int64 `json:"NetDnsQPS"`
}

type SystemResource struct {
	MaxCPU    string `json:"MaxCPU"`
	MeanCPU   string `json:"MeanCPU"`
	MaxMemory string `json:"MaxMemory"`
}
