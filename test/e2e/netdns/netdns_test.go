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
	"time"
)

var _ = Describe("testing netDns ", Label("netDns"), func() {
	var testSvcIP string
	var testPodIPs []string
	var testAppNamespace string
	var targetDomain = "kubernetes.default.svc.cluster.local"
	BeforeEach(func() {
		var e error

		testAppName := "app-" + tools.RandomName()
		testAppNamespace = "ns-" + tools.RandomName()

		// create test app
		args := []string{
			fmt.Sprintf("--set=image.tag=%s", common.AppImageTag),
			fmt.Sprintf("--set=appName=%s", testAppName),
		}
		e = common.CreateTestApp(testAppName, testAppNamespace, args)
		Expect(e).NotTo(HaveOccurred(), "create test app")

		//  get test app service ip and pod ip
		svc, e := frame.GetService(testAppName, testAppNamespace)
		Expect(e).NotTo(HaveOccurred(), "get test app service")
		testSvcIP = svc.Spec.ClusterIP
		GinkgoWriter.Printf("get test service ip %v \n", testSvcIP)

		podLIst, e := frame.WaitDeploymentReadyAndCheckIP(testAppName, testAppNamespace, time.Second*60)
		Expect(e).NotTo(HaveOccurred(), "wait test app deploy ready")

		testPodIPs = make([]string, 0)
		for _, v := range podLIst.Items {
			testPodIPs = append(testPodIPs, v.Status.PodIP)
		}

		GinkgoWriter.Printf("get test pod ips %v \n", testPodIPs)

		// Clean test env
		DeferCleanup(func() {
			GinkgoWriter.Printf("delete namespace %v \n", testAppNamespace)
			Expect(frame.DeleteNamespace(testAppNamespace)).NotTo(HaveOccurred())
		})

	})
	It("success testing netDns", Label("D00001"), func() {
		var e error
		successRate := float64(1)
		successMean := int64(1500)
		crontab := "0 1"
		netDnsName := "netdns-" + tools.RandomName()

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
		netDns.Spec.Target = target

		// request
		request := new(v1beta1.NetdnsRequest)
		var perRequestTimeoutInMS = uint64(1000)
		var qps = uint64(10)
		var durationInSecond = uint64(10)
		request.PerRequestTimeoutInMS = &perRequestTimeoutInMS
		request.QPS = &qps
		request.DurationInSecond = &durationInSecond
		request.Domain = targetDomain
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

		e = common.WaitKdoctorTaskDone(frame, netDns, pluginManager.KindNameNetdns, 120)
		Expect(e).NotTo(HaveOccurred(), "wait netDns task finish")

		success, e := common.CompareResult(frame, netDnsName, pluginManager.KindNameNetdns, []string{}, reportNum)
		Expect(success).NotTo(BeFalse(), "compare report and task result")
		Expect(e).NotTo(HaveOccurred(), "compare report and task")

	})

	It("success testing netDns User", Label("D00002"), func() {
		var e error
		successRate := float64(1)
		successMean := int64(1500)
		crontab := "0 1"
		netDnsName := "netdns-" + tools.RandomName()

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
		netDns.Spec.Target = target

		// request
		request := new(v1beta1.NetdnsRequest)
		var perRequestTimeoutInMS = uint64(1000)
		var qps = uint64(10)
		var durationInSecond = uint64(10)
		request.PerRequestTimeoutInMS = &perRequestTimeoutInMS
		request.QPS = &qps
		request.DurationInSecond = &durationInSecond
		request.Domain = targetDomain
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

		e = common.WaitKdoctorTaskDone(frame, netDns, pluginManager.KindNameNetdns, 120)
		Expect(e).NotTo(HaveOccurred(), "wait netDns task finish")

		success, e := common.CompareResult(frame, netDnsName, pluginManager.KindNameNetdns, testPodIPs, reportNum)
		Expect(success).NotTo(BeFalse(), "compare report and task result")
		Expect(e).NotTo(HaveOccurred(), "compare report and task")

	})
})
