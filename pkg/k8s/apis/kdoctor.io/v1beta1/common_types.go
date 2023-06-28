// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package v1beta1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type SchedulePlan struct {

	// +kubebuilder:validation:Optional
	Schedule *string `json:"schedule,omitempty"`

	// +kubebuilder:default=60
	// +kubebuilder:validation:Minimum=1
	RoundTimeoutMinute int64 `json:"roundTimeoutMinute"`

	// +kubebuilder:default=1
	// +kubebuilder:validation:Minimum=-1
	RoundNumber int64 `json:"roundNumber"`
}

type TaskStatus struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Minimum=-1
	ExpectedRound *int64 `json:"expectedRound,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Minimum=0
	DoneRound *int64 `json:"doneRound,omitempty"`

	Finish bool `json:"finish"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=succeed;fail;unknown
	LastRoundStatus *string `json:"lastRoundStatus,omitempty"`

	History []StatusHistoryRecord `json:"history"`
}

const (
	StatusHistoryRecordStatusSucceed    = "succeed"
	StatusHistoryRecordStatusFail       = "fail"
	StatusHistoryRecordStatusOngoing    = "ongoing"
	StatusHistoryRecordStatusNotstarted = "notstarted"
)

type StatusHistoryRecord struct {

	// +kubebuilder:validation:Enum=succeed;fail;ongoing;notstarted
	Status string `json:"status"`

	// +kubebuilder:validation:Optional
	FailureReason string `json:"failureReason,omitempty"`

	RoundNumber int `json:"roundNumber"`

	// +kubebuilder:validation:Type:=string
	// +kubebuilder:validation:Format:=date-time
	StartTimeStamp metav1.Time `json:"startTimeStamp"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type:=string
	// +kubebuilder:validation:Format:=date-time
	EndTimeStamp *metav1.Time `json:"endTimeStamp,omitempty"`

	// +kubebuilder:validation:Optional
	Duration *string `json:"duration,omitempty"`

	// +kubebuilder:validation:Type:=string
	// +kubebuilder:validation:Format:=date-time
	DeadLineTimeStamp metav1.Time `json:"deadLineTimeStamp"`

	// +kubebuilder:validation:Optional
	// expected how many agents should involve
	ExpectedActorNumber *int `json:"expectedActorNumber,omitempty"`

	FailedAgentNodeList []string `json:"failedAgentNodeList"`

	SucceedAgentNodeList []string `json:"succeedAgentNodeList"`

	NotReportAgentNodeList []string `json:"notReportAgentNodeList"`
}

type NetSuccessCondition struct {

	// +kubebuilder:default=1
	// +kubebuilder:validation:Maximum=1
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Optional
	SuccessRate *float64 `json:"successRate,omitempty"`

	// +kubebuilder:default=5000
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Optional
	MeanAccessDelayInMs *int64 `json:"meanAccessDelayInMs,omitempty"`
}

type NetHttpRequest struct {

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=2
	// +kubebuilder:validation:Minimum=1
	DurationInSecond int `json:"durationInSecond,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=5
	// +kubebuilder:validation:Maximum=100
	// +kubebuilder:validation:Minimum=1
	QPS int `json:"qps,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=5
	// +kubebuilder:validation:Minimum=1
	PerRequestTimeoutInMS int `json:"perRequestTimeoutInMS,omitempty"`
}

type AgentSpec struct {
	// +kubebuilder:validation:Optional
	Annotation map[string]string `json:"annotation,omitempty"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=Deployment;DaemonSet
	Kind string `json:"kind"`

	// +kubebuilder:validation:Optional
	DeploymentReplicas *int32 `json:"deploymentReplicas,omitempty"`

	// +kubebuilder:validation:Optional
	Affinity *v1.Affinity `json:"affinity,omitempty"`

	// +kubebuilder:validation:Optional
	Env []v1.EnvVar `json:"env,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=false
	HostNetwork bool `json:"hostNetwork,omitempty"`

	// +kubebuilder:validation:Optional
	Resources *v1.ResourceRequirements `json:"resources,omitempty"`

	// +kubebuilder:validation:Optional
	PersistentVolumeClaim *v1.PersistentVolumeClaimVolumeSource `json:"persistentVolumeClaim,omitempty"`

	// +kubebuilder:validation:Optional
	TerminationGracePeriodSeconds *int64 `json:"terminationGracePeriodSeconds,omitempty"`
}
