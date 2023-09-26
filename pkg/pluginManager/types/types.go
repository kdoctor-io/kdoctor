// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package types

import (
	"context"
	"github.com/kdoctor-io/kdoctor/pkg/k8s/apis/system/v1beta1"
	"github.com/kdoctor-io/kdoctor/pkg/resource"

	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

type Task interface {
	KindTask() string
}

type ChainingPlugin interface {
	GetApiType() client.Object

	AgentExecuteTask(logger *zap.Logger, ctx context.Context, obj runtime.Object, r *resource.UsedResource) (failureReason string, task Task, err error)
	SetReportWithTask(report *v1beta1.Report, crdSpec interface{}, task Task) error

	// ControllerReconcile(*zap.Logger, client.Client, context.Context, reconcile.Request) (reconcile.Result, error)
	// AgentReconcile(*zap.Logger, client.Client, context.Context, reconcile.Request) (reconcile.Result, error)

	WebhookMutating(logger *zap.Logger, ctx context.Context, obj runtime.Object) error
	WebhookValidateCreate(logger *zap.Logger, ctx context.Context, obj runtime.Object) error
	WebhookValidateUpdate(logger *zap.Logger, ctx context.Context, oldObj, newObj runtime.Object) error
}

type RoundResultStatus string

const (
	RoundResultSucceed = RoundResultStatus("succeed")
	RoundResultFail    = RoundResultStatus("fail")
)

type PluginReport struct {
	TaskName       string
	TaskSpec       interface{}
	RoundNumber    int
	RoundResult    RoundResultStatus
	NodeName       string
	PodName        string
	FailedReason   string
	StartTimeStamp time.Time
	EndTimeStamp   time.Time
	RoundDuraiton  string
	ReportType     ReportTypeType
	Detail         interface{}
}

type ReportTypeType string

const (
	ReportTypeSummary = "round summray report"
	ReportTypeAgent   = "agent test report"
)

type PluginRoundDetail map[string]interface{}

const (
	ApiMsgGetFailure      = "failed to get instance"
	ApiMsgUnknowCRD       = "unsupported crd type"
	ApiMsgUnsupportModify = "unsupported modify spec"
)
