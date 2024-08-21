// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package apphttphealthy

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"strings"

	k8sObjManager "github.com/kdoctor-io/kdoctor/pkg/k8ObjManager"
	crd "github.com/kdoctor-io/kdoctor/pkg/k8s/apis/kdoctor.io/v1beta1"
	"github.com/kdoctor-io/kdoctor/pkg/k8s/apis/system/v1beta1"
	"github.com/kdoctor-io/kdoctor/pkg/loadRequest/loadHttp"
	"github.com/kdoctor-io/kdoctor/pkg/pluginManager/types"
	"github.com/kdoctor-io/kdoctor/pkg/resource"
	"github.com/kdoctor-io/kdoctor/pkg/runningTask"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"
)

func ParseSuccessCondition(successCondition *crd.NetSuccessCondition, metricResult *v1beta1.HttpMetrics) (failureReason string) {
	switch {
	case successCondition.SuccessRate != nil && float64(metricResult.SuccessCounts)/float64(metricResult.RequestCounts) < *(successCondition.SuccessRate):
		failureReason = fmt.Sprintf("Success Rate %v is lower than request %v", float64(metricResult.SuccessCounts)/float64(metricResult.RequestCounts), *(successCondition.SuccessRate))
	case successCondition.MeanAccessDelayInMs != nil && int64(metricResult.Latencies.Mean) > *(successCondition.MeanAccessDelayInMs):
		failureReason = fmt.Sprintf("mean delay %v ms is bigger than request %v ms", metricResult.Latencies.Mean, *(successCondition.MeanAccessDelayInMs))
	case metricResult.ExistsNotSendRequests:
		failureReason = "There are unsent requests after the execution time has been reached"
	default:
		failureReason = ""
	}
	return
}

func SendRequestAndReport(logger *zap.Logger, targetName string, req *loadHttp.HttpRequestData, successCondition *crd.NetSuccessCondition) (failureReason string, report v1beta1.AppHttpHealthyTaskDetail) {
	report.TargetName = targetName
	report.TargetUrl = req.Url
	report.TargetMethod = string(req.Method)

	result := loadHttp.HttpRequest(logger, req)
	report.MeanDelay = result.Latencies.Mean
	report.SucceedRate = float64(result.SuccessCounts) / float64(result.RequestCounts)

	failureReason = ParseSuccessCondition(successCondition, result)

	// generate report
	// notice , upper case for first character of key, or else fail to parse json
	report.Metrics = *result
	if len(failureReason) == 0 {
		report.FailureReason = nil
		report.Succeed = true
		logger.Sugar().Infof("succeed to test %v", req.Url)
	} else {
		report.FailureReason = pointer.String(failureReason)
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

func (s *PluginAppHttpHealthy) AgentExecuteTask(logger *zap.Logger, ctx context.Context, obj runtime.Object, rt *runningTask.RunningTask) (finalfailureReason string, finalReport types.Task, err error) {
	// process mem cpu stats
	resourceStats := resource.InitResource(ctx)
	resourceStats.RunResourceCollector()

	finalfailureReason = ""
	task := &v1beta1.AppHttpHealthyTask{}
	err = nil

	instance, ok := obj.(*crd.AppHttpHealthy)
	if !ok {
		msg := "failed to get instance"
		logger.Error(msg)
		err = fmt.Errorf("error: %v", msg)
		return finalfailureReason, task, err
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
		ExpectStatusCode:    instance.Spec.SuccessCondition.StatusCode,
		EnableLatencyMetric: instance.Spec.Target.EnableLatencyMetric,
	}

	// https cert
	if target.TlsSecretName != nil {

		tlsData, err := k8sObjManager.GetK8sObjManager().GetSecret(context.Background(), *target.TlsSecretName, *target.TlsSecretNamespace)
		if err != nil {
			err = fmt.Errorf("failed get [%s/%s] secret err : %v", *target.TlsSecretNamespace, *target.TlsSecretName, err)
			logger.Sugar().Errorf(err.Error())
			return finalfailureReason, task, err
		}
		ca, caOk := tlsData.Data["ca.crt"]
		if caOk {
			caCertPool := x509.NewCertPool()
			caCertPool.AppendCertsFromPEM(ca)
			d.CaCertPool = caCertPool
		}
		crt, crtOk := tlsData.Data["tls.crt"]
		key, keyOk := tlsData.Data["tls.key"]
		if crtOk && keyOk {
			cert, err := tls.X509KeyPair(crt, key)
			if err != nil {
				err := fmt.Errorf("failed to load certificate and key: %v", err)
				logger.Sugar().Errorf(err.Error())
				return finalfailureReason, task, err
			}
			d.ClientCert = cert
		}

	}

	// body
	if target.BodyConfigName != nil {
		bodyCM, err := k8sObjManager.GetK8sObjManager().GetConfigMap(context.Background(), *target.BodyConfigName, *target.BodyConfigNamespace)
		if err != nil {
			errMsg := fmt.Errorf("failed get [%s/%s] configmap err : %v", *target.BodyConfigNamespace, *target.BodyConfigName, err)
			logger.Sugar().Errorf(errMsg.Error())
			return finalfailureReason, task, errMsg
		}
		body, err := json.Marshal(bodyCM.Data)
		if err != nil {
			errMsg := fmt.Errorf("failed get body from [%s/%s] configmap err : %v", *target.BodyConfigNamespace, *target.BodyConfigName, err)
			logger.Sugar().Errorf(errMsg.Error())
			return finalfailureReason, task, errMsg
		}
		d.Body = body
	}

	if len(target.Header) != 0 {
		header := make(map[string]string, len(target.Header))
		for _, v := range target.Header {
			h := strings.Split(v, ":")
			header[h[0]] = h[1]
		}
		d.Header = header
	}

	failureReason, itemReport := SendRequestAndReport(logger, "HttpAppHealthy target", d, successCondition)
	if len(failureReason) > 0 {
		finalfailureReason = fmt.Sprintf("test HttpAppHealthy target: %v", failureReason)
	}

	task.Detail = []v1beta1.AppHttpHealthyTaskDetail{itemReport}
	if len(finalfailureReason) > 0 {
		logger.Sugar().Errorf("plugin finally failed, %v", finalfailureReason)
		task.FailureReason = pointer.String(finalfailureReason)
		task.Succeed = false
	} else {
		task.Succeed = true
	}
	task.SystemResource = resourceStats.Stats()
	resourceStats.Stop()
	task.TotalRunningLoad = rt.QpsStats()

	return finalfailureReason, task, err

}

func (s *PluginAppHttpHealthy) SetReportWithTask(report *v1beta1.Report, task types.Task) error {
	AppHttpHealthyTask, ok := task.(*v1beta1.AppHttpHealthyTask)
	if !ok {
		return fmt.Errorf("task type %v doesn't match HttpAppHealthyTask", task.KindTask())
	}

	report.TaskAppHttpHealthy = AppHttpHealthyTask
	return nil
}
