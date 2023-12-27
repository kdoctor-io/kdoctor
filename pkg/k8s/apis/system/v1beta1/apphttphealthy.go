// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const AppHttpHealthyTaskName = "AppHttpHealthy"

type AppHttpHealthyTask struct {
	TargetType       string                     `json:"targetType"`
	TargetNumber     int64                      `json:"targetNumber"`
	FailureReason    *string                    `json:"reasonsForFailure,omitempty"`
	Succeed          bool                       `json:"roundSucceed"`
	SystemResource   SystemResource             `json:"systemResource"`
	TotalRunningLoad TotalRunningLoad           `json:"runningLoadTotal"`
	Detail           []AppHttpHealthyTaskDetail `json:"roundTaskDetail"`
}

type AppHttpHealthyTaskDetail struct {
	TargetName    string      `json:"name"`
	TargetUrl     string      `json:"url"`
	TargetMethod  string      `json:"method"`
	Succeed       bool        `json:"requestSucceed"`
	MeanDelay     float32     `json:"requestMeanDelay"`
	SucceedRate   float64     `json:"requestSucceedRate"`
	FailureReason *string     `json:"reasonsForFailure,omitempty"`
	Metrics       HttpMetrics `json:"requestTargetMetrics"`
}

type HttpMetrics struct {
	StartTime             metav1.Time         `json:"requestStartTime"`
	EndTime               metav1.Time         `json:"requestEndTime"`
	Duration              string              `json:"requestDuration"`
	RequestCounts         int64               `json:"requestCounts"`
	SuccessCounts         int64               `json:"successCounts"`
	TPS                   float64             `json:"tps"`
	Errors                map[string]int      `json:"errors"`
	Latencies             LatencyDistribution `json:"latencies"`
	ExistsNotSendRequests bool                `json:"existsNotSendRequests"`

	// request data size
	TotalDataSize string      `json:"totalDataSize"`
	StatusCodes   map[int]int `json:"statusCodes"`
}

func (h *AppHttpHealthyTask) KindTask() string {
	return AppHttpHealthyTaskName
}
