// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package netdns_test

import (
	"fmt"
	"github.com/kdoctor-io/kdoctor/pkg/k8s/apis/kdoctor.io/v1beta1"
	"github.com/kdoctor-io/kdoctor/pkg/pluginManager"
	"github.com/kdoctor-io/kdoctor/test/e2e/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spidernet-io/e2eframework/tools"
)

var _ = Describe("testing netDns ", Label("netDns"), func() {

	var targetDomain = "%s.kubernetes.default.svc.cluster.local"
	var termMin = int64(1)
	It("Successfully testing Cluster Dns Server case", Label("D00001", "C00005", "E00003"), func() {
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
		agentSpec.TerminationGracePeriodMinutes = &termMin
		netDns.Spec.AgentSpec = *agentSpec

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
		netDns.Spec.Target = target

		// request
		request := new(v1beta1.NetdnsRequest)
		var perRequestTimeoutInMS = uint64(1000)
		var qps = uint64(10)
		var durationInSecond = uint64(10)
		request.PerRequestTimeoutInMS = &perRequestTimeoutInMS
		request.QPS = &qps
		request.DurationInSecond = &durationInSecond
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

		success, e := common.CompareResult(frame, netDnsName, pluginManager.KindNameNetdns, []string{}, reportNum, netDns)
		Expect(e).NotTo(HaveOccurred(), "compare report and task")
		Expect(success).NotTo(BeFalse(), "compare report and task result")

		e = common.CheckRuntimeDeadLine(frame, netDnsName, pluginManager.KindNameNetdns, 120)
		Expect(e).NotTo(HaveOccurred(), "check task runtime resource delete")

	})

	It("Successfully testing User Define Dns server case", Label("D00002"), func() {
		var e error
		successRate := float64(1)
		successMean := int64(1500)
		crontab := "0 1"
		netDnsName := "netdns-e2e-" + tools.RandomName()

		netDns := new(v1beta1.Netdns)
		netDns.Name = netDnsName

		// agentSpec
		agentSpec := new(v1beta1.AgentSpec)
		agentSpec.TerminationGracePeriodMinutes = &termMin
		netDns.Spec.AgentSpec = *agentSpec

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
		netDns.Spec.Target = target

		// request
		request := new(v1beta1.NetdnsRequest)
		var perRequestTimeoutInMS = uint64(1000)
		var qps = uint64(10)
		var durationInSecond = uint64(10)
		request.PerRequestTimeoutInMS = &perRequestTimeoutInMS
		request.QPS = &qps
		request.DurationInSecond = &durationInSecond
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

		e = common.CheckRuntimeDeadLine(frame, netDnsName, pluginManager.KindNameNetdns, 120)
		Expect(e).NotTo(HaveOccurred(), "check task runtime resource delete")
	})
})
