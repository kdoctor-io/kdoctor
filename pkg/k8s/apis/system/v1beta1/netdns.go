// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const NetDNSTaskName = "Netdns"

type NetDNSTask struct {
	TargetType       string             `json:"targetType"`
	TargetNumber     int64              `json:"targetNumber"`
	FailureReason    *string            `json:"reasonsForFailure,omitempty"`
	Succeed          bool               `json:"roundSucceed"`
	SystemResource   SystemResource     `json:"systemResource"`
	TotalRunningLoad TotalRunningLoad   `json:"runningLoadTotal"`
	Detail           []NetDNSTaskDetail `json:"roundTaskDetail"`
}

type NetDNSTaskDetail struct {
	TargetName     string     `json:"name"`
	TargetServer   string     `json:"requestServer"`
	TargetProtocol string     `json:"protocol"`
	Succeed        bool       `json:"requestSucceed"`
	FailureReason  *string    `json:"reasonsForFailure"`
	MeanDelay      float32    `json:"requestMeanDelay"`
	SucceedRate    float64    `json:"requestSucceedRate"`
	Metrics        DNSMetrics `json:"requestTargetMetrics"`
}

type DNSMetrics struct {
	StartTime             metav1.Time         `json:"requestStartTime"`
	EndTime               metav1.Time         `json:"requestEndTime"`
	Duration              string              `json:"requestDuration"`
	RequestCounts         int64               `json:"requestCounts"`
	SuccessCounts         int64               `json:"successCounts"`
	TPS                   float64             `json:"tps"`
	Errors                map[string]int      `json:"errors"`
	Latencies             LatencyDistribution `json:"latencies"`
	ExistsNotSendRequests bool                `json:"existsNotSendRequests"`
	TargetDomain          string              `json:"domain"`
	DNSServer             string              `json:"server"`
	DNSMethod             string              `json:"method"`
	FailedCounts          int64               `json:"failedCounts"`
	ReplyCode             map[string]int      `json:"replyCode"`
}

func (n *NetDNSTask) KindTask() string {
	return NetDNSTaskName
}
