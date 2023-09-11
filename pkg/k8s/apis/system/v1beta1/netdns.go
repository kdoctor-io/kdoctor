// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const NetDNSTaskName = "Netdns"

type NetDNSTask struct {
	TargetType    string             `json:"targetType"`
	TargetNumber  int64              `json:"targetNumber"`
	FailureReason *string            `json:"failureReason,omitempty"`
	Succeed       bool               `json:"succeed"`
	MaxCPU        string             `json:"MaxCPU"`
	MaxMemory     string             `json:"MaxMemory"`
	Detail        []NetDNSTaskDetail `json:"detail"`
}

type NetDNSTaskDetail struct {
	TargetName     string     `json:"TargetName"`
	TargetServer   string     `json:"TargetServer"`
	TargetProtocol string     `json:"TargetProtocol"`
	Succeed        bool       `json:"Succeed"`
	FailureReason  *string    `json:"FailureReason"`
	MeanDelay      float32    `json:"MeanDelay"`
	SucceedRate    float64    `json:"SucceedRate"`
	Metrics        DNSMetrics `json:"Metrics"`
}

type DNSMetrics struct {
	StartTime     metav1.Time         `json:"StartTime"`
	EndTime       metav1.Time         `json:"EndTime"`
	Duration      string              `json:"Duration"`
	RequestCounts int64               `json:"RequestCounts"`
	SuccessCounts int64               `json:"SuccessCounts"`
	TPS           float64             `json:"TPS"`
	Errors        map[string]int      `json:"Errors"`
	Latencies     LatencyDistribution `json:"Latencies"`

	TargetDomain string         `json:"TargetDomain"`
	DNSServer    string         `json:"DNSServer"`
	DNSMethod    string         `json:"DNSMethod"`
	FailedCounts int64          `json:"FailedCounts"`
	ReplyCode    map[string]int `json:"ReplyCode"`
}

func (n *NetDNSTask) KindTask() string {
	return NetDNSTaskName
}
