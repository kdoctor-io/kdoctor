// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package ci_threshold_test

import (
	"fmt"
	"github.com/kdoctor-io/kdoctor/pkg/k8s/apis/kdoctor.io/v1beta1"
	"github.com/kdoctor-io/kdoctor/pkg/pluginManager"
	"github.com/kdoctor-io/kdoctor/pkg/types"
	"github.com/kdoctor-io/kdoctor/test/e2e/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spidernet-io/e2eframework/tools"
	"net"
	"time"
)

var _ = Describe("testing ci threshold test ", Serial, Label("ci threshold"), func() {
	var termMin = int64(1)
	var requestTimeout = 3000
	var replicas = int32(1)
	It("deployment apphttp threshold", Label("appHealth threshold"), func() {
		var e error
		successRate := float64(1)
		successMean := int64(1500)
		crontab := "0 1"

		tmp := new(v1beta1.AppHttpHealthy)

		// agent
		agentSpec := new(v1beta1.AgentSpec)
		agentSpec.TerminationGracePeriodMinutes = &termMin
		agentSpec.DeploymentReplicas = &replicas
		agentSpec.Kind = types.KindDeployment

		tmp.Spec.AgentSpec = agentSpec

		// successCondition
		successCondition := new(v1beta1.NetSuccessCondition)
		successCondition.SuccessRate = &successRate
		successCondition.MeanAccessDelayInMs = &successMean
		tmp.Spec.SuccessCondition = successCondition

		// target
		target := new(v1beta1.AppHttpHealthyTarget)
		target.Method = "GET"
		if net.ParseIP(testSvcIP).To4() == nil {
			target.Host = fmt.Sprintf("http://[%s]:%d", testSvcIP, httpPort)
		} else {
			target.Host = fmt.Sprintf("http://%s:%d", testSvcIP, httpPort)
		}
		tmp.Spec.Target = target

		// request
		request := new(v1beta1.NetHttpRequest)
		request.PerRequestTimeoutInMS = requestTimeout
		request.DurationInSecond = 10
		tmp.Spec.Request = request

		// Schedule
		Schedule := new(v1beta1.SchedulePlan)
		Schedule.Schedule = &crontab
		Schedule.RoundNumber = 1
		Schedule.RoundTimeoutMinute = 1
		tmp.Spec.Schedule = Schedule

		maxQPS := 100
		tmp.Spec.Request.QPS = maxQPS

		for {
			appHttpHealth := new(v1beta1.AppHttpHealthy)
			appHttpHealthName := "apphttphealth-get" + tools.RandomName()
			appHttpHealth.Name = appHttpHealthName
			appHttpHealth.Spec = tmp.Spec

			e = frame.CreateResource(appHttpHealth)
			Expect(e).NotTo(HaveOccurred(), "create appHttpHealth resource")

			e = common.WaitKdoctorTaskDone(frame, appHttpHealth, pluginManager.KindNameAppHttpHealthy, 120)
			Expect(e).NotTo(HaveOccurred(), "wait appHttpHealth task finish")

			time.Sleep(time.Second * 10)
			r, e := common.GetPluginReportResult(frame, appHttpHealthName, int(replicas))
			if e != nil {
				GinkgoWriter.Printf("failed get report \n")
				break
			}

			Reports := *r.Spec.Report

			if !Reports[0].HttpAppHealthyTask.Succeed {
				e = frame.DeleteResource(appHttpHealth)
				Expect(e).NotTo(HaveOccurred(), "delete resource failed")
				break
			} else {
				maxQPS = appHttpHealth.Spec.Request.QPS
				request.QPS += 100
			}
			e = frame.DeleteResource(appHttpHealth)
			Expect(e).NotTo(HaveOccurred(), "delete resource failed")
		}

		GinkgoWriter.Printf("max QPS in AppHttpHealthy is %d \n", maxQPS)
	})
	It("deployment netreach threshold", Label("netreach threshold"), func() {
		var e error
		successRate := float64(1)
		successMean := int64(1500)
		crontab := "0 1"

		tmp := new(v1beta1.NetReach)

		// agentSpec
		agentSpec := new(v1beta1.AgentSpec)
		agentSpec.TerminationGracePeriodMinutes = &termMin
		agentSpec.DeploymentReplicas = &replicas
		agentSpec.Kind = types.KindDeployment
		tmp.Spec.AgentSpec = agentSpec

		// successCondition
		successCondition := new(v1beta1.NetSuccessCondition)
		successCondition.SuccessRate = &successRate
		successCondition.MeanAccessDelayInMs = &successMean
		tmp.Spec.SuccessCondition = successCondition
		enable := true
		disable := false
		// target
		target := new(v1beta1.NetReachTarget)
		if !common.TestIPv4 && common.TestIPv6 {
			target.Ingress = &disable
		} else {
			target.Ingress = &enable
		}
		target.LoadBalancer = &enable
		target.ClusterIP = &enable
		target.Endpoint = &enable
		target.NodePort = &enable
		target.MultusInterface = &disable
		target.IPv4 = &common.TestIPv4
		target.IPv6 = &common.TestIPv6
		tmp.Spec.Target = target

		// request
		request := new(v1beta1.NetHttpRequest)
		request.PerRequestTimeoutInMS = requestTimeout
		request.DurationInSecond = 10
		tmp.Spec.Request = request

		// Schedule
		Schedule := new(v1beta1.SchedulePlan)
		Schedule.Schedule = &crontab
		Schedule.RoundNumber = 1
		Schedule.RoundTimeoutMinute = 1
		tmp.Spec.Schedule = Schedule

		maxQPS := 100
		tmp.Spec.Request.QPS = maxQPS

		for {
			netReach := new(v1beta1.NetReach)
			netReachName := "netreach-" + tools.RandomName()
			netReach.Name = netReachName
			netReach.Spec = tmp.Spec
			e = frame.CreateResource(netReach)
			Expect(e).NotTo(HaveOccurred(), "create appHttpHealth resource")

			e = common.WaitKdoctorTaskDone(frame, netReach, pluginManager.KindNameNetReach, 120)
			Expect(e).NotTo(HaveOccurred(), "wait appHttpHealth task finish")

			time.Sleep(time.Second * 10)
			r, e := common.GetPluginReportResult(frame, netReachName, int(replicas))
			if e != nil {
				GinkgoWriter.Printf("failed get report \n")
				break
			}

			Reports := *r.Spec.Report

			if !Reports[0].NetReachTask.Succeed {
				e = frame.DeleteResource(netReach)
				Expect(e).NotTo(HaveOccurred(), "delete resource failed")
				break
			} else {
				maxQPS = netReach.Spec.Request.QPS
				request.QPS += 100
			}
			e = frame.DeleteResource(netReach)
			Expect(e).NotTo(HaveOccurred(), "delete resource failed")
		}

		GinkgoWriter.Printf("max QPS in netReach is %d \n", maxQPS)
	})

	It("deployment netdns threshold", Label("netdns threshold"), func() {
		var e error
		successRate := float64(1)
		successMean := int64(1500)
		crontab := "0 1"

		tmp := new(v1beta1.Netdns)

		// agentSpec
		agentSpec := new(v1beta1.AgentSpec)
		agentSpec.TerminationGracePeriodMinutes = &termMin
		agentSpec.DeploymentReplicas = &replicas
		agentSpec.Kind = types.KindDeployment
		tmp.Spec.AgentSpec = agentSpec

		// successCondition
		successCondition := new(v1beta1.NetSuccessCondition)
		successCondition.SuccessRate = &successRate
		successCondition.MeanAccessDelayInMs = &successMean
		tmp.Spec.SuccessCondition = successCondition

		// target
		target := new(v1beta1.NetDnsTarget)
		targetDnsUser := new(v1beta1.NetDnsTargetUserSpec)
		targetDnsUser.Server = &testSvcIP
		port := 53
		targetDnsUser.Port = &port
		target.NetDnsTargetUser = targetDnsUser
		tmp.Spec.Target = target

		// request
		request := new(v1beta1.NetdnsRequest)
		request.PerRequestTimeoutInMS = 1000
		request.QPS = 10
		request.DurationInSecond = 10

		protocol := "udp"
		request.Protocol = &protocol
		tmp.Spec.Request = request

		// Schedule
		Schedule := new(v1beta1.SchedulePlan)
		Schedule.Schedule = &crontab
		Schedule.RoundNumber = 1
		Schedule.RoundTimeoutMinute = 1
		tmp.Spec.Schedule = Schedule

		maxQPS := 100
		tmp.Spec.Request.QPS = maxQPS
		for {
			netDns := new(v1beta1.Netdns)
			netDnsName := "netdns-e2e-" + tools.RandomName()
			netDns.Name = netDnsName
			netDns.Spec = tmp.Spec
			request.Domain = fmt.Sprintf("%s.kubernetes.default.svc.cluster.local", netDnsName)
			e = frame.CreateResource(netDns)
			Expect(e).NotTo(HaveOccurred(), "create appHttpHealth resource")

			e = common.WaitKdoctorTaskDone(frame, netDns, pluginManager.KindNameNetdns, 120)
			Expect(e).NotTo(HaveOccurred(), "wait appHttpHealth task finish")

			time.Sleep(time.Second * 10)
			r, e := common.GetPluginReportResult(frame, netDnsName, int(replicas))
			if e != nil {
				GinkgoWriter.Printf("failed get report \n")
				break
			}

			Reports := *r.Spec.Report

			if !Reports[0].NetDNSTask.Succeed {
				e = frame.DeleteResource(netDns)
				Expect(e).NotTo(HaveOccurred(), "delete resource failed")
				break
			} else {
				maxQPS = netDns.Spec.Request.QPS
				request.QPS += 100
			}
			e = frame.DeleteResource(netDns)
			Expect(e).NotTo(HaveOccurred(), "delete resource failed")
		}

		GinkgoWriter.Printf("max QPS in netDns is %d \n", maxQPS)
	})

})
