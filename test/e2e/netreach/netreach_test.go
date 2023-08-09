// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package netreach_test

import (
	"github.com/kdoctor-io/kdoctor/pkg/k8s/apis/kdoctor.io/v1beta1"
	"github.com/kdoctor-io/kdoctor/pkg/pluginManager"
	"github.com/kdoctor-io/kdoctor/test/e2e/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spidernet-io/e2eframework/tools"
)

var _ = Describe("testing netReach ", Label("netReach"), func() {
	var termMin = int64(3)
	It("success testing netReach", Label("B00001", "C00004"), func() {
		var e error
		successRate := float64(1)
		successMean := int64(1500)
		crontab := "0 1"
		netReachName := "netreach-" + tools.RandomName()

		netReach := new(v1beta1.NetReach)
		netReach.Name = netReachName

		// agentSpec
		agentSpec := new(v1beta1.AgentSpec)
		agentSpec.TerminationGracePeriodMinutes = &termMin
		netReach.Spec.AgentSpec = *agentSpec

		// successCondition
		successCondition := new(v1beta1.NetSuccessCondition)
		successCondition.SuccessRate = &successRate
		successCondition.MeanAccessDelayInMs = &successMean
		netReach.Spec.SuccessCondition = successCondition

		// target
		target := new(v1beta1.NetReachTarget)
		if !common.TestIPv4 && common.TestIPv6 {
			target.Ingress = false
		} else {
			target.Ingress = true
		}
		target.LoadBalancer = true
		target.ClusterIP = true
		target.Endpoint = true
		target.NodePort = true
		target.MultusInterface = false
		target.IPv4 = &common.TestIPv4
		target.IPv6 = &common.TestIPv6
		netReach.Spec.Target = target

		// request
		request := new(v1beta1.NetHttpRequest)
		request.PerRequestTimeoutInMS = 1000
		request.QPS = 10
		request.DurationInSecond = 10
		netReach.Spec.Request = request

		// Schedule
		Schedule := new(v1beta1.SchedulePlan)
		Schedule.Schedule = &crontab
		Schedule.RoundNumber = 1
		Schedule.RoundTimeoutMinute = 1
		netReach.Spec.Schedule = Schedule

		e = frame.CreateResource(netReach)
		Expect(e).NotTo(HaveOccurred(), "create netReach resource")

		e = common.WaitKdoctorTaskDone(frame, netReach, pluginManager.KindNameNetReach, 120)
		Expect(e).NotTo(HaveOccurred(), "wait netReach task finish")

		success, e := common.CompareResult(frame, netReachName, pluginManager.KindNameNetReach, []string{}, reportNum)
		Expect(e).NotTo(HaveOccurred(), "compare report and task")
		Expect(success).To(BeTrue(), "compare report and task result")

	})
})
