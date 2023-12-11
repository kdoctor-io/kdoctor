// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package v1beta1

import (
	"github.com/kdoctor-io/kdoctor/pkg/k8s/apis/kdoctor.io/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KdoctorReport
// +k8s:openapi-gen=true
type KdoctorReport struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Report Reports `json:"report,omitempty"`

	Status Status `json:"status,omitempty"`

	Task TaskInfo `json:"task,omitempty"`
}

type Reports struct {
	LatestRoundReport *[]Report `json:"latestRoundReport,omitempty"`
}

// KdoctorReportList
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type KdoctorReportList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []KdoctorReport `json:"items"`
}

type Report struct {
	RoundNumber    int64       `json:"roundNumber"`
	RoundResult    string      `json:"roundResult"`
	NodeName       string      `json:"nodeName"`
	PodName        string      `json:"podName"`
	FailedReason   *string     `json:"reasonsForFailure,omitempty"`
	StartTimeStamp metav1.Time `json:"roundStartTimeStamp"`
	EndTimeStamp   metav1.Time `json:"roundEndTimeStamp"`
	RoundDuration  string      `json:"roundDuration"`

	TaskNetReach *NetReachTask `json:"taskNetReach,omitempty"`

	TaskAppHttpHealthy *AppHttpHealthyTask `json:"taskAppHealthy,omitempty"`

	TaskNetDNS *NetDNSTask `json:"taskNetDns,omitempty"`
}

type Status struct {
	ToTalRoundNumber    int64  `json:"totalRoundNumber"`
	FinishedRoundNumber int64  `json:"roundFinishedNumber"`
	Status              string `json:"status"`
	RoundNumber         int64  `json:"roundNumber"`
}

type TaskInfo struct {
	TaskName string   `json:"name"`
	TaskType string   `json:"kind"`
	Spec     TaskSpec `json:"spec"`
}

type TaskSpec struct {
	NetReachTaskSpec *v1beta1.NetReachSpec `json:"netReach,omitempty"`

	AppHttpHealthyTaskSpec *v1beta1.AppHttpHealthySpec `json:"appHttpHealthy,omitempty"`

	NetDNSTaskSpec *v1beta1.NetdnsSpec `json:"netDns,omitempty"`
}
