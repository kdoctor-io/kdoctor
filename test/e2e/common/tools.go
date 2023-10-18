// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package common

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"reflect"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"

	"github.com/docker/docker/api/types"
	docker_client "github.com/docker/docker/client"
	"github.com/onsi/ginkgo/v2"
	frame "github.com/spidernet-io/e2eframework/framework"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kdoctor-io/kdoctor/api/v1/agentServer/models"
	"github.com/kdoctor-io/kdoctor/pkg/k8s/apis/kdoctor.io/v1beta1"
	kdoctor_report "github.com/kdoctor-io/kdoctor/pkg/k8s/apis/system/v1beta1"
	"github.com/kdoctor-io/kdoctor/pkg/pluginManager"
	kdoctor_types "github.com/kdoctor-io/kdoctor/pkg/types"
)

func WaitKdoctorTaskDone(f *frame.Framework, task client.Object, taskKind string, timeout int) error {
	interval := time.Duration(10)
	switch taskKind {
	case pluginManager.KindNameNetReach:
		fake := &v1beta1.NetReach{
			ObjectMeta: metav1.ObjectMeta{
				Name: task.GetName(),
			},
		}
		key := client.ObjectKeyFromObject(fake)
		after := time.After(time.Duration(timeout) * time.Second)

		for {
			select {
			case <-after:
				return fmt.Errorf("timeout wait task NetReach %s finish", task.GetName())
			default:
				rs := &v1beta1.NetReach{}
				if err := f.GetResource(key, rs); err != nil {
					return fmt.Errorf("failed get resource NetReach %s ,err : %v", task.GetName(), err)
				}
				if rs.Status.Finish {
					return nil
				}
				time.Sleep(time.Second * interval)
			}
		}
	case pluginManager.KindNameAppHttpHealthy:
		fake := &v1beta1.AppHttpHealthy{
			ObjectMeta: metav1.ObjectMeta{
				Name: task.GetName(),
			},
		}
		key := client.ObjectKeyFromObject(fake)
		after := time.After(time.Duration(timeout) * time.Second)

		for {
			select {
			case <-after:
				return fmt.Errorf("timeout wait task AppHttpHealthy %s finish", task.GetName())
			default:
				rs := &v1beta1.AppHttpHealthy{}
				if err := f.GetResource(key, rs); err != nil {
					return fmt.Errorf("failed get resource AppHttpHealthy %s,err: %v ", task.GetName(), err)
				}
				if rs.Status.Finish {
					return nil
				}
				time.Sleep(time.Second * interval)
			}
		}
	case pluginManager.KindNameNetdns:
		fake := &v1beta1.Netdns{
			ObjectMeta: metav1.ObjectMeta{
				Name: task.GetName(),
			},
		}
		key := client.ObjectKeyFromObject(fake)
		after := time.After(time.Duration(timeout) * time.Second)

		for {
			select {
			case <-after:
				return fmt.Errorf("timeout wait task Netdns %s finish", task.GetName())
			default:
				rs := &v1beta1.Netdns{}
				if err := f.GetResource(key, rs); err != nil {
					return fmt.Errorf("failed get resource Netdns %s, err: %v", task.GetName(), err)
				}
				if rs.Status.Finish {
					return nil
				}
				time.Sleep(time.Second * interval)
			}
		}
	default:
		return fmt.Errorf("unknown task type: %s", task.GetObjectKind().GroupVersionKind().Kind)
	}

}

func GetKdoctorToken(f *frame.Framework) (string, error) {

	if KdoctroTestToken != "" {
		return KdoctroTestToken, nil
	}

	fake := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      KdoctorTestTokenSecretName,
			Namespace: TestNameSpace,
		},
	}
	key := client.ObjectKeyFromObject(fake)

	s := &corev1.Secret{}
	if err := f.KClient.Get(context.TODO(), key, s); err != nil {
		return "", fmt.Errorf("get kdoctor token secret failed, err: %v", err)
	}
	KdoctroTestToken = string(s.Data["token"])

	return KdoctroTestToken, nil
}

func GetPluginReportResult(f *frame.Framework, name string, n int) (*kdoctor_report.KdoctorReport, error) {
	var err error
	url := fmt.Sprintf("%s%s%s", APISERVICEADDR, PluginReportPath, name)
	ginkgo.GinkgoWriter.Println(url)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	tr.TLSNextProto = make(map[string]func(string, *tls.Conn) http.RoundTripper)
	c := &http.Client{Transport: tr}
	req, _ := http.NewRequest("GET", url, nil)
	token, err := GetKdoctorToken(f)
	if err != nil {
		return nil, fmt.Errorf("get kdoctor test token failed,err : %v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get plugin report failed,err : %v", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("plugin report not found ")
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("read plugin report body failed,err : %v", err)
	}
	ginkgo.GinkgoWriter.Println(string(body))
	report := new(kdoctor_report.KdoctorReport)

	err = json.Unmarshal(body, report)
	if err != nil {
		return nil, fmt.Errorf("unmarshal plugin report failed,err : %v", err)
	}

	if len(*report.Spec.Report) != n {
		return nil, fmt.Errorf("have agent not upload report")
	}

	return report, nil
}

func CompareResult(f *frame.Framework, name, taskKind string, podIPs []string, n int, object client.Object) (bool, error) {

	// get Aggregate API report
	var r *kdoctor_report.KdoctorReport
	var err error
	c := time.After(time.Second * 60)
	r, err = GetPluginReportResult(f, name, n)
	for err != nil {
		select {
		case <-c:
			return false, fmt.Errorf("get %s %s report time out,err: %v ", taskKind, name, err)
		default:
			time.Sleep(time.Second * 5)
			r, err = GetPluginReportResult(f, name, n)
			break
		}
	}
	switch taskKind {
	case pluginManager.KindNameNetReach:
		obj := object.(*v1beta1.NetReach)
		fake := &v1beta1.NetReach{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			},
		}
		key := client.ObjectKeyFromObject(fake)
		rs := &v1beta1.NetReach{}
		if err = f.GetResource(key, rs); err != nil {
			return GetResultFromReport(r), fmt.Errorf("failed get resource AppHttpHealthy %s", name)
		}

		if *obj.Spec.Target.Ingress != *rs.Spec.Target.Ingress {
			return GetResultFromReport(r), fmt.Errorf("spec target ingress not equal input %v,output %v", *obj.Spec.Target.Ingress, *rs.Spec.Target.Ingress)
		}
		if obj.Spec.Target.EnableLatencyMetric != rs.Spec.Target.EnableLatencyMetric {
			return GetResultFromReport(r), fmt.Errorf("spec target EnableLatencyMetric not equal input %v,output %v", obj.Spec.Target.EnableLatencyMetric, rs.Spec.Target.EnableLatencyMetric)
		}
		if *obj.Spec.Target.Endpoint != *rs.Spec.Target.Endpoint {
			return GetResultFromReport(r), fmt.Errorf("spec target Endpoint not equal input %v,output %v", *obj.Spec.Target.Endpoint, *rs.Spec.Target.Endpoint)
		}
		if *obj.Spec.Target.ClusterIP != *rs.Spec.Target.ClusterIP {
			return GetResultFromReport(r), fmt.Errorf("spec target ClusterIP not equal input %v,output %v", *obj.Spec.Target.ClusterIP, *rs.Spec.Target.ClusterIP)
		}
		if *obj.Spec.Target.LoadBalancer != *rs.Spec.Target.LoadBalancer {
			return GetResultFromReport(r), fmt.Errorf("spec target LoadBalancer not equal input %v,output %v", *obj.Spec.Target.LoadBalancer, *rs.Spec.Target.LoadBalancer)
		}
		if *obj.Spec.Target.MultusInterface != *rs.Spec.Target.MultusInterface {
			return GetResultFromReport(r), fmt.Errorf("spec target MultusInterface not equal input %v,output %v", *obj.Spec.Target.MultusInterface, *rs.Spec.Target.MultusInterface)
		}
		if *obj.Spec.Target.IPv4 != *rs.Spec.Target.IPv4 {
			return GetResultFromReport(r), fmt.Errorf("spec target IPv4 not equal input %v,output %v", *obj.Spec.Target.IPv4, *rs.Spec.Target.IPv4)
		}
		if *obj.Spec.Target.IPv6 != *rs.Spec.Target.IPv6 {
			return GetResultFromReport(r), fmt.Errorf("spec target IPv6 not equal input %v,output %v", *obj.Spec.Target.IPv6, *rs.Spec.Target.IPv6)
		}

		for _, v := range *r.Spec.Report {
			for _, m := range v.NetReachTask.Detail {
				// qps
				expectRequestCount := float64(rs.Spec.Request.QPS * rs.Spec.Request.DurationInSecond)
				realRequestCount := float64(m.Metrics.RequestCounts)
				if math.Abs(realRequestCount-expectRequestCount)/expectRequestCount > 0.1 {
					return GetResultFromReport(r), fmt.Errorf("The error in the number of requests is greater than 0.1,real request count: %d,expect request count:%d", int(realRequestCount), int(expectRequestCount))
				}
				if float64(m.Metrics.SuccessCounts)/float64(m.Metrics.RequestCounts) != m.SucceedRate {
					return GetResultFromReport(r), fmt.Errorf("succeedRate not equal")
				}
			}
			// startTime
			shcedule := pluginManager.NewSchedule(*rs.Spec.Schedule.Schedule)
			startTime := shcedule.StartTime(rs.CreationTimestamp.Time)
			if v.StartTimeStamp.Time.Compare(startTime) > 5 {
				return GetResultFromReport(r), fmt.Errorf("The task start time error is greater than 5 seconds ")
			}

		}

		// roundNumber
		rounds := []int64{r.Spec.FinishedRoundNumber, r.Spec.ReportRoundNumber, r.Spec.ToTalRoundNumber, rs.Spec.Schedule.RoundNumber}
		for i := 0; i < len(rounds); i++ {
			for j := i + 1; j < len(rounds); j++ {
				if rounds[i] != rounds[j] {
					return GetResultFromReport(r), fmt.Errorf("roundNumber not equal ")
				}
			}
		}

		return GetResultFromReport(r), nil
	case pluginManager.KindNameAppHttpHealthy:
		fake := &v1beta1.NetReach{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			},
		}
		key := client.ObjectKeyFromObject(fake)
		rs := &v1beta1.AppHttpHealthy{}
		if err = f.GetResource(key, rs); err != nil {
			return GetResultFromReport(r), fmt.Errorf("failed get resource AppHttpHealthy %s", name)
		}

		if r.Spec.Report == nil {
			return false, fmt.Errorf("failed get resource AppHttpHealthy %s report", name)
		}

		var reportRequestCount int64
		var realRequestCount int64
		for _, v := range *r.Spec.Report {
			for _, m := range v.HttpAppHealthyTask.Detail {
				// qps
				expectCount := float64(rs.Spec.Request.QPS * rs.Spec.Request.DurationInSecond)
				realCount := float64(m.Metrics.RequestCounts)
				// report request count
				reportRequestCount += m.Metrics.RequestCounts
				if math.Abs(realCount-expectCount)/expectCount > 0.1 {
					return GetResultFromReport(r), fmt.Errorf("The error in the number of requests is greater than 0.1,real request count: %d,expect request count:%d", int(realCount), int(expectCount))
				}
				if float64(m.Metrics.SuccessCounts)/float64(m.Metrics.RequestCounts) != m.SucceedRate {
					return GetResultFromReport(r), fmt.Errorf("succeedRate not equal")
				}
			}
			// startTime
			shcedule := pluginManager.NewSchedule(*rs.Spec.Schedule.Schedule)
			startTime := shcedule.StartTime(rs.CreationTimestamp.Time)
			if v.StartTimeStamp.Time.Compare(startTime) > 5 {
				return GetResultFromReport(r), fmt.Errorf("The task start time error is greater than 5 seconds ")
			}
		}
		// real request count
		if len(podIPs) != 0 {
			for _, ip := range podIPs {
				count, e := GetRealRequestCount(name, ip)
				if e != nil {
					return GetResultFromReport(r), e
				}
				realRequestCount += count
			}

			if realRequestCount != reportRequestCount {
				return GetResultFromReport(r), fmt.Errorf("real request count %d not equal report request count %d ", realRequestCount, reportRequestCount)
			}
		}

		// roundNumber
		rounds := []int64{r.Spec.FinishedRoundNumber, r.Spec.ReportRoundNumber, r.Spec.ToTalRoundNumber, rs.Spec.Schedule.RoundNumber}
		for i := 0; i < len(rounds); i++ {
			for j := i + 1; j < len(rounds); j++ {
				if rounds[i] != rounds[j] {
					return GetResultFromReport(r), fmt.Errorf("roundNumber not equal ")
				}
			}
		}

		return GetResultFromReport(r), nil
	case pluginManager.KindNameNetdns:
		fake := &v1beta1.Netdns{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			},
		}
		key := client.ObjectKeyFromObject(fake)
		rs := &v1beta1.Netdns{}
		if err = f.GetResource(key, rs); err != nil {
			return GetResultFromReport(r), fmt.Errorf("failed get resource AppHttpHealthy %s", name)
		}

		if r.Spec.Report == nil {
			return false, fmt.Errorf("failed get resource AppHttpHealthy %s report", name)
		}

		var reportRequestCount int64
		var realRequestCount int64
		for _, v := range *r.Spec.Report {
			for _, m := range v.NetDNSTask.Detail {
				// qps
				expectCount := float64((*rs.Spec.Request.QPS) * (*rs.Spec.Request.DurationInSecond))
				realCount := float64(m.Metrics.RequestCounts)
				// report request count
				reportRequestCount += m.Metrics.RequestCounts
				if math.Abs(realCount-expectCount)/expectCount > 0.1 {
					return GetResultFromReport(r), fmt.Errorf("The error in the number of requests is greater than 0.1, real request count: %d,expect request count:%d ", int(realCount), int(expectCount))
				}
				if float64(m.Metrics.SuccessCounts)/float64(m.Metrics.RequestCounts) != m.SucceedRate {
					return GetResultFromReport(r), fmt.Errorf("succeedRate not equal")
				}
			}
			// startTime
			shcedule := pluginManager.NewSchedule(*rs.Spec.Schedule.Schedule)
			startTime := shcedule.StartTime(rs.CreationTimestamp.Time)
			if v.StartTimeStamp.Time.Compare(startTime) > 5 {
				return GetResultFromReport(r), fmt.Errorf("The task start time error is greater than 5 seconds ")
			}
		}
		// real request count
		if len(podIPs) != 0 {
			for _, ip := range podIPs {
				count, e := GetRealRequestCount(name, ip)
				if e != nil {
					return GetResultFromReport(r), e
				}
				realRequestCount += count
			}

			if realRequestCount != reportRequestCount {
				return GetResultFromReport(r), fmt.Errorf("real request count %d not equal report request count %d ", int(realRequestCount), int(reportRequestCount))
			}
		}

		// roundNumber
		rounds := []int64{r.Spec.FinishedRoundNumber, r.Spec.ReportRoundNumber, r.Spec.ToTalRoundNumber, rs.Spec.Schedule.RoundNumber}
		for i := 0; i < len(rounds); i++ {
			for j := i + 1; j < len(rounds); j++ {
				if rounds[i] != rounds[j] {
					return GetResultFromReport(r), fmt.Errorf("roundNumber not equal ")
				}
			}
		}
		return GetResultFromReport(r), nil
	default:
		return GetResultFromReport(r), fmt.Errorf("unknown task type: %s", taskKind)
	}

}

func GetResultFromReport(r *kdoctor_report.KdoctorReport) bool {
	for _, v := range *r.Spec.Report {
		if v.NetReachTask != nil {
			return v.NetReachTask.Succeed
		}
		if v.HttpAppHealthyTask != nil {
			return v.HttpAppHealthyTask.Succeed
		}
		if v.NetDNSTask != nil {
			return v.NetDNSTask.Succeed
		}
	}
	return true
}

func GetRealRequestCount(name string, ip string) (int64, error) {
	cmd := []string{"curl"}
	if net.ParseIP(ip).To4() == nil {
		cmd = append(cmd, fmt.Sprintf("http://[%s]/?task=%s", ip, name))
		cmd = append(cmd, "-6g")
	} else {
		cmd = append(cmd, fmt.Sprintf("http://%s/?task=%s", ip, name))
	}

	ctx := context.Background()
	cli, err := docker_client.NewClientWithOpts(
		docker_client.FromEnv,
		docker_client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return 0, fmt.Errorf("new docker client failed : %v", err)
	}
	listOpt := types.ContainerListOptions{}
	containers, err := cli.ContainerList(ctx, listOpt)
	if err != nil {
		return 0, fmt.Errorf("list docker containers failed : %v", err)
	}
	var containerID string
	for _, container := range containers {
		for _, name := range container.Names {
			if strings.Contains(name, KindClusterName) {
				containerID = container.ID
				break
			}
		}
		if containerID != "" {
			break
		}
	}
	execCfg := types.ExecConfig{
		AttachStderr: true,
		AttachStdout: true,
		AttachStdin:  true,
		Cmd:          cmd,
		Tty:          true,
	}
	exec, err := cli.ContainerExecCreate(ctx, containerID, execCfg)
	if err != nil {
		return 0, fmt.Errorf("create docker container cmd failed : %v", err)
	}
	resp, err := cli.ContainerExecAttach(ctx, exec.ID, types.ExecStartCheck{})
	if err != nil {
		return 0, fmt.Errorf("exec docker container cmd failed : %v", err)
	}
	defer resp.Close()
	r, err := io.ReadAll(resp.Reader)
	if err != nil {
		return 0, fmt.Errorf("read docker server response failed : %v", err)
	}
	var start int
	if string(r)[strings.Index(string(r), "{")+1] == '{' {
		start = strings.Index(string(r), "{") + 1
	} else {
		start = strings.Index(string(r), "{")
	}
	s := string(r)[start : strings.LastIndex(string(r), "}")+1]
	ginkgo.GinkgoWriter.Println(s)
	echoRes := new(models.EchoRes)
	err = json.Unmarshal([]byte(s), echoRes)
	if err != nil {
		return 0, fmt.Errorf("unmarshal docker server response failed : %v", err)
	}
	return echoRes.RequestCount - 1, nil
}

func CheckRuntime(f *frame.Framework, task client.Object, taskKind string, timeout int) error {
	interval := time.Duration(10)
	switch taskKind {
	case pluginManager.KindNameNetReach:
		fake := &v1beta1.NetReach{
			ObjectMeta: metav1.ObjectMeta{
				Name: task.GetName(),
			},
		}
		key := client.ObjectKeyFromObject(fake)
		rs := &v1beta1.NetReach{}
		after := time.After(time.Duration(timeout) * time.Second)
		create := false
		for !create {
			select {
			case <-after:
				return fmt.Errorf("timeout wait task %s %s runtime create", pluginManager.KindNameNetReach, task.GetName())
			default:
				if err := f.GetResource(key, rs); err != nil {
					return fmt.Errorf("failed get resource %s %s, err: %v", pluginManager.KindNameNetReach, task.GetName(), err)
				}
				if rs.Status.Resource == nil {
					continue
				}
				if rs.Status.Resource.RuntimeStatus == v1beta1.RuntimeCreated {
					create = true
					break
				}
				time.Sleep(time.Second * interval)
			}
		}
		if err := checkAgentSpec(f, rs, rs.Spec.AgentSpec, rs.Status, pluginManager.KindNameNetReach); err != nil {
			return fmt.Errorf("check AgentSpce err,reason= %v ", err)
		}
	case pluginManager.KindNameAppHttpHealthy:
		fake := &v1beta1.AppHttpHealthy{
			ObjectMeta: metav1.ObjectMeta{
				Name: task.GetName(),
			},
		}
		key := client.ObjectKeyFromObject(fake)
		rs := &v1beta1.AppHttpHealthy{}
		after := time.After(time.Duration(timeout) * time.Second)
		create := false
		for !create {
			select {
			case <-after:
				return fmt.Errorf("timeout wait task %s %s runtime create", pluginManager.KindNameAppHttpHealthy, task.GetName())
			default:
				if err := f.GetResource(key, rs); err != nil {
					return fmt.Errorf("failed get resource %s %s, err: %v", pluginManager.KindNameAppHttpHealthy, task.GetName(), err)
				}
				if rs.Status.Resource == nil {
					continue
				}
				if rs.Status.Resource.RuntimeStatus == v1beta1.RuntimeCreated {
					create = true
					break
				}
				time.Sleep(time.Second * interval)
			}
		}
		if err := checkAgentSpec(f, rs, rs.Spec.AgentSpec, rs.Status, pluginManager.KindNameAppHttpHealthy); err != nil {
			return fmt.Errorf("check AgentSpce err,reason= %v ", err)
		}
	case pluginManager.KindNameNetdns:
		fake := &v1beta1.Netdns{
			ObjectMeta: metav1.ObjectMeta{
				Name: task.GetName(),
			},
		}
		key := client.ObjectKeyFromObject(fake)
		rs := &v1beta1.Netdns{}
		after := time.After(time.Duration(timeout) * time.Second)
		create := false
		for !create {
			select {
			case <-after:
				return fmt.Errorf("timeout wait task %s %s runtime create", pluginManager.KindNameNetdns, task.GetName())
			default:
				if err := f.GetResource(key, rs); err != nil {
					return fmt.Errorf("failed get resource %s %s, err: %v", pluginManager.KindNameNetdns, task.GetName(), err)
				}
				if rs.Status.Resource == nil {
					continue
				}
				if rs.Status.Resource.RuntimeStatus == v1beta1.RuntimeCreated {
					create = true
					break
				}
				time.Sleep(time.Second * interval)
			}
		}
		if err := checkAgentSpec(f, rs, rs.Spec.AgentSpec, rs.Status, pluginManager.KindNameNetdns); err != nil {
			return fmt.Errorf("check AgentSpce err,reason= %v ", err)
		}
	default:
		return fmt.Errorf("unknown task type: %s", task.GetObjectKind().GroupVersionKind().Kind)
	}
	return nil
}

// checkAgentSpec check agentSpec generate deployment or daemonSet is right
func checkAgentSpec(f *frame.Framework, task client.Object, agentSpec *v1beta1.AgentSpec, taskStatus v1beta1.TaskStatus, taskKind string) error {

	if agentSpec == nil {
		return nil
	}

	switch agentSpec.Kind {
	case kdoctor_types.KindDaemonSet:
		if taskStatus.Resource.RuntimeType != kdoctor_types.KindDaemonSet {
			return fmt.Errorf("agent spec is %s,but status reource runtimeType is %s ", kdoctor_types.KindDaemonSet, taskStatus.Resource.RuntimeType)
		}
		runtime, err := f.GetDaemonSet(taskStatus.Resource.RuntimeName, TestNameSpace)
		if err != nil {
			return fmt.Errorf("failed get runtime daemonset %s ,reason=%v ", taskStatus.Resource.RuntimeName, err)
		}
		// compare annotation
		if !reflect.DeepEqual(runtime.Spec.Template.Annotations, agentSpec.Annotation) {
			return fmt.Errorf("create runtime daemonset annotation not equal,spec: %v real: %v", agentSpec.Annotation, runtime.Spec.Template.Annotations)
		}

		// TODO(ii2day): compare env resource affinity

		// compare HostNetwork
		if !reflect.DeepEqual(runtime.Spec.Template.Spec.HostNetwork, agentSpec.HostNetwork) {
			return fmt.Errorf("create runtime daemonset HostNetwork not equal,spec: %v real: %v", agentSpec.HostNetwork, runtime.Spec.Template.Spec.HostNetwork)
		}

	case kdoctor_types.KindDeployment:
		runtime, err := f.GetDeployment(taskStatus.Resource.RuntimeName, TestNameSpace)
		if err != nil {
			return fmt.Errorf("failed get runtime deployment %s ,reason=%v ", taskStatus.Resource.RuntimeName, err)
		}
		// compare annotation
		if !reflect.DeepEqual(runtime.Spec.Template.Annotations, agentSpec.Annotation) {
			return fmt.Errorf("create runtime deployment annotation not equal,spec: %v real: %v", agentSpec.Annotation, runtime.Spec.Template.Annotations)
		}
		// TODO(ii2day): compare env resource affinity

		// compare HostNetwork
		if !reflect.DeepEqual(runtime.Spec.Template.Spec.HostNetwork, agentSpec.HostNetwork) {
			return fmt.Errorf("create runtime deployment HostNetwork not equal,spec: %v real: %v", agentSpec.HostNetwork, runtime.Spec.Template.Spec.HostNetwork)
		}

		// compare replicas
		if *agentSpec.DeploymentReplicas != runtime.Status.Replicas {
			return fmt.Errorf("create runtime deployment Replicas not equal,spec: %d real: %d", *agentSpec.DeploymentReplicas, runtime.Status.Replicas)
		}
	default:
		return fmt.Errorf("unknown agent kind %s ", agentSpec.Kind)
	}

	if TestIPv4 {
		_, err := f.GetService(*taskStatus.Resource.ServiceNameV4, TestNameSpace)
		if err != nil {
			return fmt.Errorf("failed get service %s ,reason=%v ", *taskStatus.Resource.ServiceNameV4, err)
		}
	}

	if TestIPv6 {
		_, err := f.GetService(*taskStatus.Resource.ServiceNameV6, TestNameSpace)
		if err != nil {
			return fmt.Errorf("failed get service %s ,reason=%v ", *taskStatus.Resource.ServiceNameV6, err)
		}
	}
	if taskKind == kdoctor_types.KindNameNetReach {
		taskNr := task.(*v1beta1.NetReach)
		if *taskNr.Spec.Target.Ingress {
			ig := &networkingv1.Ingress{}
			fake := &networkingv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name:      *taskStatus.Resource.ServiceNameV4,
					Namespace: TestNameSpace,
				},
			}
			key := client.ObjectKeyFromObject(fake)
			if err := f.GetResource(key, ig); err != nil {
				return fmt.Errorf("task %s enable ingress test,but not found %s ingress err=%v", task.GetName(), *taskStatus.Resource.ServiceNameV4, err)
			}
		}
	}
	return nil
}

// CheckRuntimeDeadLine check terminationGracePeriodMinutes delete deployment or daemonSet service ingress.
func CheckRuntimeDeadLine(f *frame.Framework, taskName, taskKind string, timeout int) error {
	interval := time.Duration(10)
	var runtimeResource *v1beta1.TaskResource
	var terminationGracePeriodMinutes int64
	var testIngress = false
	switch taskKind {
	case kdoctor_types.KindNameNetReach:
		fake := &v1beta1.NetReach{
			ObjectMeta: metav1.ObjectMeta{
				Name: taskName,
			},
		}
		key := client.ObjectKeyFromObject(fake)
		rs := &v1beta1.NetReach{}
		after := time.After(time.Duration(timeout) * time.Second)
		done := false
		for !done {
			select {
			case <-after:
				return fmt.Errorf("timeout wait task %s %s runtime deadline", pluginManager.KindNameNetReach, taskName)
			default:
				if err := f.GetResource(key, rs); err != nil {
					return fmt.Errorf("failed get resource %s %s, err: %v", pluginManager.KindNameNetReach, taskName, err)
				}
				if rs.Status.FinishTime != nil {
					done = true
					runtimeResource = rs.Status.Resource
					terminationGracePeriodMinutes = *rs.Spec.AgentSpec.TerminationGracePeriodMinutes
					testIngress = *rs.Spec.Target.Ingress
					break
				}
				time.Sleep(time.Second * interval)
			}
		}
	case kdoctor_types.KindNameNetdns:
		fake := &v1beta1.Netdns{
			ObjectMeta: metav1.ObjectMeta{
				Name: taskName,
			},
		}
		key := client.ObjectKeyFromObject(fake)
		rs := &v1beta1.Netdns{}
		after := time.After(time.Duration(timeout) * time.Minute)
		done := false
		for !done {
			select {
			case <-after:
				return fmt.Errorf("timeout wait task %s %s runtime deadline", pluginManager.KindNameNetdns, taskName)
			default:
				if err := f.GetResource(key, rs); err != nil {
					return fmt.Errorf("failed get resource %s %s, err: %v", pluginManager.KindNameNetdns, taskName, err)
				}
				if rs.Status.FinishTime != nil {
					done = true
					runtimeResource = rs.Status.Resource
					terminationGracePeriodMinutes = *rs.Spec.AgentSpec.TerminationGracePeriodMinutes
					break
				}
				time.Sleep(time.Second * interval)
			}
		}
	case kdoctor_types.KindNameAppHttpHealthy:
		fake := &v1beta1.AppHttpHealthy{
			ObjectMeta: metav1.ObjectMeta{
				Name: taskName,
			},
		}
		key := client.ObjectKeyFromObject(fake)
		rs := &v1beta1.AppHttpHealthy{}
		after := time.After(time.Duration(timeout) * time.Minute)
		done := false
		for !done {
			select {
			case <-after:
				return fmt.Errorf("timeout wait task %s %s runtime deadline", pluginManager.KindNameAppHttpHealthy, taskName)
			default:
				if err := f.GetResource(key, rs); err != nil {
					return fmt.Errorf("failed get resource %s %s, err: %v", pluginManager.KindNameAppHttpHealthy, taskName, err)
				}
				if rs.Status.FinishTime != nil {
					done = true
					runtimeResource = rs.Status.Resource
					terminationGracePeriodMinutes = *rs.Spec.AgentSpec.TerminationGracePeriodMinutes
					break
				}
				time.Sleep(time.Second * interval)
			}
		}
	default:
		return fmt.Errorf("unknown task %s type: %s", taskName, taskKind)
	}

	c := time.After(time.Duration(terminationGracePeriodMinutes) * time.Minute)
	// wait delete time
	<-c

	c2 := time.After(time.Minute)

	// check runtime delete
	runtimeDeleted := false
	for !runtimeDeleted {
		select {
		case <-c2:
			return fmt.Errorf("timeout delete %s task %s runtime", taskKind, taskName)
		default:
			if runtimeResource.RuntimeType == kdoctor_types.KindDaemonSet {
				_, err := f.GetDaemonSet(runtimeResource.RuntimeName, TestNameSpace)
				if err != nil {
					if errors.IsNotFound(err) {
						ginkgo.GinkgoWriter.Printf("task runtime daemonSet %s deleted \n", runtimeResource.RuntimeName)
						runtimeDeleted = true
						break
					}
					return fmt.Errorf("failed get task %s deamonSet reason=%v", taskName, err)
				}

			} else {
				_, err := f.GetDeployment(runtimeResource.RuntimeName, TestNameSpace)
				if err != nil {
					if errors.IsNotFound(err) {
						ginkgo.GinkgoWriter.Printf("task runtime deployment %s deleted \n", runtimeResource.RuntimeName)
						runtimeDeleted = true
						break
					}
					return fmt.Errorf("failed get task %s deamonSet reason=%v", taskName, err)
				}
			}
			time.Sleep(time.Second * interval)
		}
	}

	// check runtime service ingress delete
	c3 := time.After(time.Minute)
	serviceV4Deleted := false
	serviceV6Deleted := false
	ingressDeleted := false
	for !serviceV4Deleted || !serviceV6Deleted {
		select {
		case <-c3:
			return fmt.Errorf("timeout delete %s task %s runtime service", taskKind, taskName)
		default:
			if TestIPv4 {
				_, err := f.GetService(*runtimeResource.ServiceNameV4, TestNameSpace)
				if err != nil {
					if errors.IsNotFound(err) {
						ginkgo.GinkgoWriter.Printf("task runtime service v4 %s deleted \n", *runtimeResource.ServiceNameV4)
						serviceV4Deleted = true
					} else {
						return fmt.Errorf("failed get task %s service %s reason=%v", taskName, *runtimeResource.ServiceNameV4, err)
					}
				}
				if testIngress && !ingressDeleted {
					ig := &networkingv1.Ingress{}
					fake := &networkingv1.Ingress{
						ObjectMeta: metav1.ObjectMeta{
							Name:      *runtimeResource.ServiceNameV4,
							Namespace: TestNameSpace,
						},
					}
					key := client.ObjectKeyFromObject(fake)
					err = f.GetResource(key, ig)
					if err != nil {
						if errors.IsNotFound(err) {
							ginkgo.GinkgoWriter.Printf("task runtime ingress %s deleted \n", *runtimeResource.ServiceNameV4)
							ingressDeleted = true
						} else {
							return fmt.Errorf("task %s enable ingress test,but not found %s ingress err=%v", taskName, *runtimeResource.ServiceNameV4, err)
						}
					}
				} else {
					ingressDeleted = true
				}
			} else {
				serviceV4Deleted = true
			}
			if TestIPv6 {
				_, err := f.GetService(*runtimeResource.ServiceNameV6, TestNameSpace)
				if err != nil {
					if errors.IsNotFound(err) {
						ginkgo.GinkgoWriter.Printf("task runtime service v6 %s deleted \n", *runtimeResource.ServiceNameV6)
						serviceV6Deleted = true
						break
					}
					return fmt.Errorf("failed get task %s service %s reason=%v", taskName, *runtimeResource.ServiceNameV6, err)
				}
			} else {
				serviceV6Deleted = true
			}
			time.Sleep(time.Second * interval)
		}
	}

	return nil

}

func GetRuntimeResource(f *frame.Framework, resource *v1beta1.TaskResource, ingress bool) error {
	if resource == nil {
		return fmt.Errorf("runtime resource is nil")
	}

	c := time.After(time.Minute)
	<-c

	switch resource.RuntimeType {
	case kdoctor_types.KindDaemonSet:
		_, err := f.GetDaemonSet(resource.RuntimeName, TestNameSpace)
		if !errors.IsNotFound(err) {
			return fmt.Errorf("after 1 min runtime daemonset %s not delete", resource.RuntimeName)
		}
	case kdoctor_types.KindDeployment:
		_, err := f.GetDeployment(resource.RuntimeName, TestNameSpace)
		if !errors.IsNotFound(err) {
			return fmt.Errorf("after 1 min runtime deployment %s not delete", resource.RuntimeName)
		}
	}

	if resource.ServiceNameV4 != nil {
		_, err := f.GetService(*resource.ServiceNameV4, TestNameSpace)
		if !errors.IsNotFound(err) {
			return fmt.Errorf("after 1 min runtime service %s not delete", *resource.ServiceNameV4)
		}

		if ingress {
			ig := &networkingv1.Ingress{}
			fake := &networkingv1.Ingress{
				ObjectMeta: metav1.ObjectMeta{
					Name:      *resource.ServiceNameV4,
					Namespace: TestNameSpace,
				},
			}
			key := client.ObjectKeyFromObject(fake)
			err := f.GetResource(key, ig)
			if !errors.IsNotFound(err) {
				return fmt.Errorf("after 1 min runtime ingress %s not delete", *resource.ServiceNameV4)
			}
		}
	}

	if resource.ServiceNameV6 != nil {
		_, err := f.GetService(*resource.ServiceNameV6, TestNameSpace)
		if !errors.IsNotFound(err) {
			return fmt.Errorf("after 1 min runtime service %s not delete", *resource.ServiceNameV6)
		}
	}

	return nil
}
