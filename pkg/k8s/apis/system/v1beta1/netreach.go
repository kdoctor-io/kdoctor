// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package v1beta1

const NetReachTaskName = "NetReach"

type NetReachTask struct {
	TargetType       string               `json:"targetType"`
	TargetNumber     int64                `json:"targetNumber"`
	FailureReason    *string              `json:"reasonsForFailure,omitempty"`
	Succeed          bool                 `json:"roundSucceed"`
	SystemResource   SystemResource       `json:"systemResource"`
	TotalRunningLoad TotalRunningLoad     `json:"runningLoadTotal"`
	Detail           []NetReachTaskDetail `json:"roundTaskDetail"`
}

type NetReachTaskDetail struct {
	TargetName    string      `json:"name"`
	TargetUrl     string      `json:"url"`
	TargetMethod  string      `json:"method"`
	Succeed       bool        `json:"requestSucceed"`
	MeanDelay     float32     `json:"requestMeanDelay"`
	SucceedRate   float64     `json:"requestSucceedRate"`
	FailureReason *string     `json:"failureReason,omitempty"`
	Metrics       HttpMetrics `json:"requestTargetMetrics"`
}

func (n *NetReachTask) KindTask() string {
	return NetReachTaskName
}
