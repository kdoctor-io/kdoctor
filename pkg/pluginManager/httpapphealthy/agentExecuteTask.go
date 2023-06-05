// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package httpapphealthy

import (
	"context"
	"fmt"
	crd "github.com/kdoctor-io/kdoctor/pkg/k8s/apis/kdoctor.io/v1beta1"
	"github.com/kdoctor-io/kdoctor/pkg/k8s/apis/system/v1beta1"
	"github.com/kdoctor-io/kdoctor/pkg/loadRequest/loadHttp"
	"github.com/kdoctor-io/kdoctor/pkg/pluginManager/types"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"
)

func ParseSuccessCondition(successCondition *crd.NetSuccessCondition, metricResult *v1beta1.HttpMetrics) (failureReason string, err error) {
	switch {
	case successCondition.SuccessRate != nil && float64(metricResult.SuccessCounts)/float64(metricResult.RequestCounts) < *(successCondition.SuccessRate):
		failureReason = fmt.Sprintf("Success Rate %v is lower than request %v", float64(metricResult.SuccessCounts)/float64(metricResult.RequestCounts), *(successCondition.SuccessRate))
	case successCondition.MeanAccessDelayInMs != nil && int64(metricResult.Latencies.Mean) > *(successCondition.MeanAccessDelayInMs):
		failureReason = fmt.Sprintf("mean delay %v ms is bigger than request %v ms", metricResult.Latencies.Mean, *(successCondition.MeanAccessDelayInMs))
	default:
		failureReason = ""
		err = nil
	}
	return
}

func SendRequestAndReport(logger *zap.Logger, targetName string, req *loadHttp.HttpRequestData, successCondition *crd.NetSuccessCondition) (failureReason string, report v1beta1.HttpAppHealthyTaskDetail) {
	report.TargetName = targetName
	report.TargetUrl = req.Url
	report.TargetMethod = string(req.Method)

	result := loadHttp.HttpRequest(logger, req)
	report.MeanDelay = result.Latencies.Mean
	report.SucceedRate = float64(result.SuccessCounts) / float64(result.RequestCounts)

	var err error
	failureReason, err = ParseSuccessCondition(successCondition, result)
	if err != nil {
		failureReason = fmt.Sprintf("%v", err)
		logger.Sugar().Errorf("internal error for target %v, error=%v", req.Url, err)
		report.FailureReason = pointer.String(failureReason)
		return
	}

	// generate report
	// notice , upper case for first character of key, or else fail to parse json
	report.Metrics = *result
	report.FailureReason = pointer.String(failureReason)
	if report.FailureReason == nil {
		report.Succeed = true
		logger.Sugar().Infof("succeed to test %v", req.Url)
	} else {
		report.Succeed = false
		logger.Sugar().Warnf("failed to test %v", req.Url)
	}

	return
}

type TestTarget struct {
	Name   string
	Url    string
	Method loadHttp.HttpMethod
}

func (s *PluginHttpAppHealthy) AgentExecuteTask(logger *zap.Logger, ctx context.Context, obj runtime.Object) (finalfailureReason string, finalReport types.Task, err error) {
	finalfailureReason = ""
	task := &v1beta1.HttpAppHealthyTask{}
	err = nil

	instance, ok := obj.(*crd.HttpAppHealthy)
	if !ok {
		msg := "failed to get instance"
		logger.Error(msg)
		err = fmt.Errorf(msg)
		return
	}

	logger.Sugar().Infof("plugin implement task round, instance=%+v", instance)

	target := instance.Spec.Target
	request := instance.Spec.Request
	successCondition := instance.Spec.SuccessCondition

	logger.Sugar().Infof("load test custom target: Method=%v, Url=%v , qps=%v, PerRequestTimeout=%vs, Duration=%vs", target.Method, target.Host, request.QPS, request.PerRequestTimeoutInMS, request.DurationInSecond)
	task.TargetType = "HttpAppHealthy"
	task.TargetNumber = 1
	d := &loadHttp.HttpRequestData{
		Method:              loadHttp.HttpMethod(target.Method),
		Url:                 target.Host,
		Qps:                 request.QPS,
		PerRequestTimeoutMS: request.PerRequestTimeoutInMS,
		RequestTimeSecond:   request.DurationInSecond,
		Http2:               target.Http2,
	}
	failureReason, itemReport := SendRequestAndReport(logger, "HttpAppHealthy target", d, successCondition)
	if len(failureReason) > 0 {
		finalfailureReason = fmt.Sprintf("test HttpAppHealthy target: %v", failureReason)
	}

	task.Detail = []v1beta1.HttpAppHealthyTaskDetail{itemReport}
	if len(finalfailureReason) > 0 {
		logger.Sugar().Errorf("plugin finally failed, %v", finalfailureReason)
		task.FailureReason = pointer.String(finalfailureReason)
		task.Succeed = false
	} else {
		task.Succeed = true
	}

	return finalfailureReason, task, err

}

func (s *PluginHttpAppHealthy) SetReportWithTask(report *v1beta1.Report, crdSpec interface{}, task types.Task) error {
	httpAppHealthySpec, ok := crdSpec.(*crd.HttpAppHealthySpec)
	if !ok {
		return fmt.Errorf("the given crd spec %#v doesn't match HttpAppHealthySpec", crdSpec)
	}

	httpAppHealthyTask, ok := task.(*v1beta1.HttpAppHealthyTask)
	if !ok {
		return fmt.Errorf("task type %v doesn't match HttpAppHealthyTask", task.KindTask())
	}

	report.HttpAppHealthyTaskSpec = httpAppHealthySpec
	report.HttpAppHealthyTask = httpAppHealthyTask
	return nil
}
