// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package netdns_test

import (
	"fmt"
	"github.com/kdoctor-io/kdoctor/pkg/k8s/apis/kdoctor.io/v1beta1"
	"github.com/kdoctor-io/kdoctor/pkg/pluginManager"
	"github.com/kdoctor-io/kdoctor/pkg/types"
	"github.com/kdoctor-io/kdoctor/test/e2e/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spidernet-io/e2eframework/tools"
	"k8s.io/utils/pointer"
)

var _ = Describe("testing netDns ", Label("netDns"), func() {

	var targetDomain = "%s.kubernetes.default.svc.cluster.local"
	successMean := int64(1500)
	successRate := float64(1)
	It("Successfully testing Cluster Dns Server case", Label("D00001", "C00005", "E00015", "E00018"), Serial, func() {
		var e error

		crontab := "0 1"
		netDnsName := "netdns-e2e-" + tools.RandomName()

		// netdns
		netDns := new(v1beta1.Netdns)
		netDns.Name = netDnsName

		// successCondition
		successCondition := new(v1beta1.NetSuccessCondition)
		successCondition.SuccessRate = &successRate
		successCondition.MeanAccessDelayInMs = &successMean
		netDns.Spec.SuccessCondition = successCondition

		// target
		target := new(v1beta1.NetDnsTarget)
		targetDns := new(v1beta1.NetDnsTargetDnsSpec)
		targetDns.TestIPv4 = &common.TestIPv4
		targetDns.TestIPv6 = &common.TestIPv6
		targetDns.ServiceName = &KubeDnsName
		targetDns.ServiceNamespace = &KubeDnsNamespace
		target.NetDnsTargetDns = targetDns
		target.EnableLatencyMetric = true
		netDns.Spec.Target = target

		// request
		request := new(v1beta1.NetdnsRequest)
		request.PerRequestTimeoutInMS = 2000
		request.QPS = 10
		request.DurationInSecond = 10
		request.Domain = "kubernetes.default.svc.cluster.local"
		protocol := "udp"
		request.Protocol = &protocol
		netDns.Spec.Request = request

		// Schedule
		Schedule := new(v1beta1.SchedulePlan)
		Schedule.Schedule = &crontab
		Schedule.RoundNumber = 1
		Schedule.RoundTimeoutMinute = 1
		netDns.Spec.Schedule = Schedule

		e = frame.CreateResource(netDns)
		Expect(e).NotTo(HaveOccurred(), "create netDns resource")

		e = common.CheckRuntime(frame, netDns, pluginManager.KindNameNetdns, 60)
		Expect(e).NotTo(HaveOccurred(), "check task runtime spec")

		e = common.WaitKdoctorTaskDone(frame, netDns, pluginManager.KindNameNetdns, 120)
		Expect(e).NotTo(HaveOccurred(), "wait netDns task finish")

		success, e := common.CompareResult(frame, netDnsName, pluginManager.KindNameNetdns, []string{}, reportNum, netDns)
		Expect(e).NotTo(HaveOccurred(), "compare report and task")
		Expect(success).NotTo(BeFalse(), "compare report and task result")

	})

	It("Successfully testing User Define Dns server case", Label("D00002"), Serial, func() {
		var e error
		crontab := "0 1"
		netDnsName := "netdns-e2e-" + tools.RandomName()

		netDns := new(v1beta1.Netdns)
		netDns.Name = netDnsName

		// successCondition
		successCondition := new(v1beta1.NetSuccessCondition)
		successCondition.SuccessRate = &successRate
		successCondition.MeanAccessDelayInMs = &successMean
		netDns.Spec.SuccessCondition = successCondition

		// target
		target := new(v1beta1.NetDnsTarget)
		targetDnsUser := new(v1beta1.NetDnsTargetUserSpec)
		targetDnsUser.Server = &testSvcIP
		port := 53
		targetDnsUser.Port = &port
		target.NetDnsTargetUser = targetDnsUser
		target.EnableLatencyMetric = true
		netDns.Spec.Target = target

		// request
		request := new(v1beta1.NetdnsRequest)
		request.PerRequestTimeoutInMS = 2000
		request.QPS = 10
		request.DurationInSecond = 10
		request.Domain = fmt.Sprintf(targetDomain, netDnsName)
		protocol := "udp"
		request.Protocol = &protocol
		netDns.Spec.Request = request

		// Schedule
		Schedule := new(v1beta1.SchedulePlan)
		Schedule.Schedule = &crontab
		Schedule.RoundNumber = 1
		Schedule.RoundTimeoutMinute = 1
		netDns.Spec.Schedule = Schedule

		e = frame.CreateResource(netDns)
		Expect(e).NotTo(HaveOccurred(), "create netDns resource")

		e = common.CheckRuntime(frame, netDns, pluginManager.KindNameNetdns, 60)
		Expect(e).NotTo(HaveOccurred(), "check task runtime spec")

		e = common.WaitKdoctorTaskDone(frame, netDns, pluginManager.KindNameNetdns, 120)
		Expect(e).NotTo(HaveOccurred(), "wait netDns task finish")

		success, e := common.CompareResult(frame, netDnsName, pluginManager.KindNameNetdns, testPodIPs, reportNum, netDns)

		Expect(e).NotTo(HaveOccurred(), "compare report and task")
		Expect(success).To(BeTrue(), "compare report and task result")

	})

	It("Successfully testing User Define Dns server case use tcp protocol", Serial, Label("D00003"), func() {
		var e error
		crontab := "0 1"
		netDnsName := "netdns-e2e-" + tools.RandomName()

		netDns := new(v1beta1.Netdns)
		netDns.Name = netDnsName

		// successCondition
		successCondition := new(v1beta1.NetSuccessCondition)
		successCondition.SuccessRate = &successRate
		successCondition.MeanAccessDelayInMs = &successMean
		netDns.Spec.SuccessCondition = successCondition

		// target
		target := new(v1beta1.NetDnsTarget)
		targetDnsUser := new(v1beta1.NetDnsTargetUserSpec)
		targetDnsUser.Server = &testSvcIP
		port := 53
		targetDnsUser.Port = &port
		target.NetDnsTargetUser = targetDnsUser
		target.EnableLatencyMetric = true
		netDns.Spec.Target = target

		// request
		request := new(v1beta1.NetdnsRequest)
		request.PerRequestTimeoutInMS = 2000
		request.QPS = 10
		request.DurationInSecond = 10
		request.Domain = fmt.Sprintf(targetDomain, netDnsName)
		protocol := "tcp"
		request.Protocol = &protocol
		netDns.Spec.Request = request

		// Schedule
		Schedule := new(v1beta1.SchedulePlan)
		Schedule.Schedule = &crontab
		Schedule.RoundNumber = 1
		Schedule.RoundTimeoutMinute = 1
		netDns.Spec.Schedule = Schedule

		e = frame.CreateResource(netDns)
		Expect(e).NotTo(HaveOccurred(), "create netDns resource")

		e = common.CheckRuntime(frame, netDns, pluginManager.KindNameNetdns, 60)
		Expect(e).NotTo(HaveOccurred(), "check task runtime spec")

		e = common.WaitKdoctorTaskDone(frame, netDns, pluginManager.KindNameNetdns, 120)
		Expect(e).NotTo(HaveOccurred(), "wait netDns task finish")

		success, e := common.CompareResult(frame, netDnsName, pluginManager.KindNameNetdns, testPodIPs, reportNum, netDns)

		Expect(e).NotTo(HaveOccurred(), "compare report and task")
		Expect(success).To(BeTrue(), "compare report and task result")
	})

	It("Successfully testing User Define Dns server case use tcp-tls protocol", Serial, Label("D00004"), func() {
		var e error
		crontab := "0 1"
		netDnsName := "netdns-e2e-" + tools.RandomName()

		netDns := new(v1beta1.Netdns)
		netDns.Name = netDnsName

		// successCondition
		successCondition := new(v1beta1.NetSuccessCondition)
		successCondition.SuccessRate = &successRate
		successCondition.MeanAccessDelayInMs = &successMean
		netDns.Spec.SuccessCondition = successCondition

		// target
		target := new(v1beta1.NetDnsTarget)
		targetDnsUser := new(v1beta1.NetDnsTargetUserSpec)
		targetDnsUser.Server = &testSvcIP
		port := 853
		targetDnsUser.Port = &port
		target.NetDnsTargetUser = targetDnsUser
		target.EnableLatencyMetric = true
		netDns.Spec.Target = target

		// request
		request := new(v1beta1.NetdnsRequest)
		request.PerRequestTimeoutInMS = 2000
		request.QPS = 10
		request.DurationInSecond = 10
		request.Domain = fmt.Sprintf(targetDomain, netDnsName)
		protocol := "tcp-tls"
		request.Protocol = &protocol
		netDns.Spec.Request = request

		// Schedule
		Schedule := new(v1beta1.SchedulePlan)
		Schedule.Schedule = &crontab
		Schedule.RoundNumber = 1
		Schedule.RoundTimeoutMinute = 1
		netDns.Spec.Schedule = Schedule

		e = frame.CreateResource(netDns)
		Expect(e).NotTo(HaveOccurred(), "create netDns resource")

		e = common.CheckRuntime(frame, netDns, pluginManager.KindNameNetdns, 60)
		Expect(e).NotTo(HaveOccurred(), "check task runtime spec")

		e = common.WaitKdoctorTaskDone(frame, netDns, pluginManager.KindNameNetdns, 120)
		Expect(e).NotTo(HaveOccurred(), "wait netDns task finish")

		success, e := common.CompareResult(frame, netDnsName, pluginManager.KindNameNetdns, testPodIPs, reportNum, netDns)
		Expect(e).NotTo(HaveOccurred(), "compare report and task")
		Expect(success).To(BeTrue(), "compare report and task result")
	})

	It("Successfully testing NetDns crontab case", Label("C00002"), Serial, func() {
		var e error
		crontab := "* * * * *"
		netDnsName := "netdns-e2e-" + tools.RandomName()

		netDns := new(v1beta1.Netdns)
		netDns.Name = netDnsName

		// successCondition
		successCondition := new(v1beta1.NetSuccessCondition)
		successCondition.SuccessRate = &successRate
		successCondition.MeanAccessDelayInMs = &successMean
		netDns.Spec.SuccessCondition = successCondition

		// target
		target := new(v1beta1.NetDnsTarget)
		targetDnsUser := new(v1beta1.NetDnsTargetUserSpec)
		targetDnsUser.Server = &testSvcIP
		port := 53
		targetDnsUser.Port = &port
		target.NetDnsTargetUser = targetDnsUser
		target.EnableLatencyMetric = true
		netDns.Spec.Target = target

		// request
		request := new(v1beta1.NetdnsRequest)
		request.PerRequestTimeoutInMS = 2000
		request.QPS = 10
		request.DurationInSecond = 10
		request.Domain = fmt.Sprintf(targetDomain, netDnsName)
		protocol := "udp"
		request.Protocol = &protocol
		netDns.Spec.Request = request

		// Schedule
		Schedule := new(v1beta1.SchedulePlan)
		Schedule.Schedule = &crontab
		Schedule.RoundNumber = 1
		Schedule.RoundTimeoutMinute = 1
		netDns.Spec.Schedule = Schedule

		e = frame.CreateResource(netDns)
		Expect(e).NotTo(HaveOccurred(), "create netDns resource")

		e = common.CheckRuntime(frame, netDns, pluginManager.KindNameNetdns, 60)
		Expect(e).NotTo(HaveOccurred(), "check task runtime spec")

		e = common.WaitKdoctorTaskDone(frame, netDns, pluginManager.KindNameNetdns, 120)
		Expect(e).NotTo(HaveOccurred(), "wait netDns task finish")

		success, e := common.CompareResult(frame, netDnsName, pluginManager.KindNameNetdns, testPodIPs, reportNum, netDns)
		Expect(e).NotTo(HaveOccurred(), "compare report and task")
		Expect(success).To(BeTrue(), "compare report and task result")

	})

	It("Successfully testing Task NetDns Runtime Deployment Service creation", Serial, Label("E00006"), func() {
		var e error
		successRate := float64(1)
		successMean := int64(1500)
		crontab := "0 1"
		netDnsName := "netdns-e2e-" + tools.RandomName()

		// netdns
		netDns := new(v1beta1.Netdns)
		netDns.Name = netDnsName

		// agentSpec
		agentSpec := new(v1beta1.AgentSpec)
		agentSpec.TerminationGracePeriodMinutes = pointer.Int64(1)
		agentSpec.Kind = types.KindDeployment
		agentSpec.DeploymentReplicas = pointer.Int32(2)
		netDns.Spec.AgentSpec = agentSpec

		// successCondition
		successCondition := new(v1beta1.NetSuccessCondition)
		successCondition.SuccessRate = &successRate
		successCondition.MeanAccessDelayInMs = &successMean
		netDns.Spec.SuccessCondition = successCondition

		// target
		target := new(v1beta1.NetDnsTarget)
		targetDns := new(v1beta1.NetDnsTargetDnsSpec)
		targetDns.TestIPv4 = &common.TestIPv4
		targetDns.TestIPv6 = &common.TestIPv6
		targetDns.ServiceName = &KubeDnsName
		targetDns.ServiceNamespace = &KubeDnsNamespace
		target.NetDnsTargetDns = targetDns
		target.EnableLatencyMetric = true
		netDns.Spec.Target = target

		// request
		request := new(v1beta1.NetdnsRequest)
		request.PerRequestTimeoutInMS = 1000
		request.QPS = 10
		request.DurationInSecond = 10
		request.Domain = "kubernetes.default.svc.cluster.local"
		protocol := "udp"
		request.Protocol = &protocol
		netDns.Spec.Request = request

		// Schedule
		Schedule := new(v1beta1.SchedulePlan)
		Schedule.Schedule = &crontab
		Schedule.RoundNumber = 1
		Schedule.RoundTimeoutMinute = 1
		netDns.Spec.Schedule = Schedule

		e = frame.CreateResource(netDns)
		Expect(e).NotTo(HaveOccurred(), "create netDns resource")

		e = common.CheckRuntime(frame, netDns, pluginManager.KindNameNetdns, 60)
		Expect(e).NotTo(HaveOccurred(), "check task runtime spec")

		e = common.WaitKdoctorTaskDone(frame, netDns, pluginManager.KindNameNetdns, 120)
		Expect(e).NotTo(HaveOccurred(), "wait netDns task finish")

		success, e := common.CompareResult(frame, netDnsName, pluginManager.KindNameNetdns, []string{}, reportNum, netDns)
		Expect(e).NotTo(HaveOccurred(), "compare report and task")
		Expect(success).NotTo(BeFalse(), "compare report and task result")

		e = common.CheckRuntimeDeadLine(frame, netDnsName, pluginManager.KindNameNetdns, 120)
		Expect(e).NotTo(HaveOccurred(), "check task runtime resource delete")

	})

	It("Successfully testing Task NetDns Runtime DaemonSet Service creation", Serial, Label("E00003"), func() {
		var e error
		successRate := float64(1)
		successMean := int64(1500)
		crontab := "0 1"
		netDnsName := "netdns-e2e-" + tools.RandomName()

		// netdns
		netDns := new(v1beta1.Netdns)
		netDns.Name = netDnsName

		// agentSpec
		agentSpec := new(v1beta1.AgentSpec)
		agentSpec.TerminationGracePeriodMinutes = pointer.Int64(1)
		netDns.Spec.AgentSpec = agentSpec

		// successCondition
		successCondition := new(v1beta1.NetSuccessCondition)
		successCondition.SuccessRate = &successRate
		successCondition.MeanAccessDelayInMs = &successMean
		netDns.Spec.SuccessCondition = successCondition

		// target
		target := new(v1beta1.NetDnsTarget)
		targetDns := new(v1beta1.NetDnsTargetDnsSpec)
		targetDns.TestIPv4 = &common.TestIPv4
		targetDns.TestIPv6 = &common.TestIPv6
		targetDns.ServiceName = &KubeDnsName
		targetDns.ServiceNamespace = &KubeDnsNamespace
		target.NetDnsTargetDns = targetDns
		target.EnableLatencyMetric = true
		netDns.Spec.Target = target

		// request
		request := new(v1beta1.NetdnsRequest)
		request.PerRequestTimeoutInMS = 1000
		request.QPS = 10
		request.DurationInSecond = 10
		request.Domain = "kubernetes.default.svc.cluster.local"
		protocol := "udp"
		request.Protocol = &protocol
		netDns.Spec.Request = request

		// Schedule
		Schedule := new(v1beta1.SchedulePlan)
		Schedule.Schedule = &crontab
		Schedule.RoundNumber = 1
		Schedule.RoundTimeoutMinute = 1
		netDns.Spec.Schedule = Schedule

		e = frame.CreateResource(netDns)
		Expect(e).NotTo(HaveOccurred(), "create netDns resource")

		e = common.CheckRuntime(frame, netDns, pluginManager.KindNameNetdns, 60)
		Expect(e).NotTo(HaveOccurred(), "check task runtime spec")

		e = common.WaitKdoctorTaskDone(frame, netDns, pluginManager.KindNameNetdns, 120)
		Expect(e).NotTo(HaveOccurred(), "wait netDns task finish")

		success, e := common.CompareResult(frame, netDnsName, pluginManager.KindNameNetdns, []string{}, reportNum, netDns)
		Expect(e).NotTo(HaveOccurred(), "compare report and task")
		Expect(success).NotTo(BeFalse(), "compare report and task result")

		e = common.CheckRuntimeDeadLine(frame, netDnsName, pluginManager.KindNameNetdns, 120)
		Expect(e).NotTo(HaveOccurred(), "check task runtime resource delete")

	})
})
