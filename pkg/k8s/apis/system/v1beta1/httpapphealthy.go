// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const HttpAppHealthyTaskName = "HttpAppHealthy"

type HttpAppHealthyTask struct {
	TargetType    string                     `json:"TargetType"`
	TargetNumber  int64                      `json:"TargetNumber"`
	FailureReason *string                    `json:"FailureReason,omitempty"`
	Succeed       bool                       `json:"Succeed"`
	Detail        []HttpAppHealthyTaskDetail `json:"Detail"`
}

type HttpAppHealthyTaskDetail struct {
	TargetName    string      `json:"TargetName"`
	TargetUrl     string      `json:"TargetUrl"`
	TargetMethod  string      `json:"TargetMethod"`
	Succeed       bool        `json:"Succeed"`
	MeanDelay     float32     `json:"MeanDelay"`
	SucceedRate   float64     `json:"SucceedRate"`
	FailureReason *string     `json:"FailureReason,omitempty"`
	Metrics       HttpMetrics `json:"Metrics"`
}

type HttpMetrics struct {
	StartTime     metav1.Time         `json:"StartTime"`
	EndTime       metav1.Time         `json:"EndTime"`
	Duration      string              `json:"Duration"`
	RequestCounts int64               `json:"RequestCounts"`
	SuccessCounts int64               `json:"SuccessCounts"`
	TPS           float64             `json:"TPS"`
	Errors        map[string]int      `json:"Errors"`
	Latencies     LatencyDistribution `json:"Latencies"`

	// request data size
	TotalDataSize string      `json:"TotalDataSize"`
	StatusCodes   map[int]int `json:"StatusCodes"`
}

func (h *HttpAppHealthyTask) KindTask() string {
	return HttpAppHealthyTaskName
}
