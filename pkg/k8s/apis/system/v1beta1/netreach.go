// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package v1beta1

const NetReachTaskName = "NetReach"

type NetReachTask struct {
	TargetType    string               `json:"TargetType"`
	TargetNumber  int64                `json:"TargetNumber"`
	FailureReason *string              `json:"FailureReason,omitempty"`
	Succeed       bool                 `json:"Succeed"`
	Detail        []NetReachTaskDetail `json:"Detail"`
}

type NetReachTaskDetail struct {
	TargetName    string      `json:"TargetName"`
	TargetUrl     string      `json:"TargetUrl"`
	TargetMethod  string      `json:"TargetMethod"`
	Succeed       bool        `json:"Succeed"`
	MeanDelay     float32     `json:"MeanDelay"`
	SucceedRate   float64     `json:"SucceedRate"`
	FailureReason *string     `json:"FailureReason,omitempty"`
	Metrics       HttpMetrics `json:"Metrics"`
}

func (n *NetReachTask) KindTask() string {
	return NetReachTaskName
}
