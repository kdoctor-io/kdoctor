// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package common

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/kdoctor-io/kdoctor/api/v1/agentServer/models"
	"github.com/onsi/ginkgo/v2"
	"io"
	"math"
	"net"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	docker_client "github.com/docker/docker/client"
	"github.com/kdoctor-io/kdoctor/pkg/k8s/apis/kdoctor.io/v1beta1"
	kdoctor_report "github.com/kdoctor-io/kdoctor/pkg/k8s/apis/system/v1beta1"
	"github.com/kdoctor-io/kdoctor/pkg/pluginManager"
	frame "github.com/spidernet-io/e2eframework/framework"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
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

func CompareResult(f *frame.Framework, name, taskKind string, podIPs []string, n int) (bool, error) {

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

		for _, v := range *r.Spec.Report {
			for _, m := range v.NetReachTask.Detail {
				// qps
				expectRequestCount := float64(rs.Spec.Request.QPS * rs.Spec.Request.DurationInSecond)
				realRequestCount := float64(m.Metrics.RequestCounts)
				if math.Abs(realRequestCount-expectRequestCount)/expectRequestCount > 0.1 {
					return GetResultFromReport(r), fmt.Errorf("The error in the number of requests is greater than 0.1,real request count: %d,expect request count:%d", int(realRequestCount), int(expectRequestCount))
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
			ginkgo.GinkgoWriter.Println("reach")
			return v.NetReachTask.Succeed
		}
		if v.HttpAppHealthyTask != nil {
			ginkgo.GinkgoWriter.Println("app")
			return v.HttpAppHealthyTask.Succeed
		}
		if v.NetDNSTask != nil {
			ginkgo.GinkgoWriter.Println("dns")
			return v.NetDNSTask.Succeed
		}
	}
	ginkgo.GinkgoWriter.Println("none")
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

func CreateTestApp(name, namespace string, o []string) error {

	var stdout, stderr bytes.Buffer
	cmd := exec.Command(
		"helm",
		"install",
		name,
		AppChartDir,
		fmt.Sprintf("--namespace=%s", namespace),
		"--create-namespace=true",
		"--wait=true",
		fmt.Sprintf("--kubeconfig=%s", KubeConfigPath),
	)
	cmd.Args = append(cmd.Args, o...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		ginkgo.GinkgoWriter.Println(stderr.String())
		ginkgo.GinkgoWriter.Println(stdout.String())
		return fmt.Errorf("start cmd [%s] failed,err: %v ", cmd.String(), err)
	}

	if err := cmd.Wait(); err != nil {
		ginkgo.GinkgoWriter.Println(stderr.String())
		ginkgo.GinkgoWriter.Println(stdout.String())
		return fmt.Errorf("run cmd [%s] failed,err: %v ", cmd.String(), err)
	}

	return nil
}
