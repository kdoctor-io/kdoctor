// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package runtime_test

import (
	"github.com/kdoctor-io/kdoctor/pkg/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spidernet-io/e2eframework/tools"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kdoctor-io/kdoctor/pkg/k8s/apis/kdoctor.io/v1beta1"
	"github.com/kdoctor-io/kdoctor/pkg/pluginManager"
	"github.com/kdoctor-io/kdoctor/test/e2e/common"
)

var _ = Describe("testing runtime ", Label("runtime"), func() {
	var termMin = int64(5)
	var replicas = int32(2)
	It("Successfully testing cascading deletion with Task NetReach DaemonSet Service and Ingress ", Label("E00007"), func() {
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
		netReach.Spec.AgentSpec = agentSpec

		// successCondition
		successCondition := new(v1beta1.NetSuccessCondition)
		successCondition.SuccessRate = &successRate
		successCondition.MeanAccessDelayInMs = &successMean
		netReach.Spec.SuccessCondition = successCondition

		// target
		target := new(v1beta1.NetReachTarget)
		enable := true
		disable := false
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

		e = common.CheckRuntime(frame, netReach, pluginManager.KindNameNetReach, 60)
		Expect(e).NotTo(HaveOccurred(), "check task runtime spec")

		// get runtime resource
		fake := &v1beta1.NetReach{
			ObjectMeta: metav1.ObjectMeta{
				Name: netReachName,
			},
		}
		key := client.ObjectKeyFromObject(fake)
		rs := &v1beta1.NetReach{}
		e = frame.GetResource(key, rs)
		Expect(e).NotTo(HaveOccurred(), "get task runtime status")

		e = frame.DeleteResource(rs)
		Expect(e).NotTo(HaveOccurred(), "delete task")
		// check runtime resource is deleted after deleting task
		e = common.GetRuntimeResource(frame, rs.Status.Resource, true)
		Expect(e).NotTo(HaveOccurred(), "get task runtime resource")
	})

	It("Successfully testing cascading deletion with Task NetAppHttpHealthy DaemonSet Service", Label("E00008"), func() {
		var e error
		successRate := float64(1)
		successMean := int64(1500)
		crontab := "0 1"
		appHttpHealthName := "apphttphealth-get" + tools.RandomName()

		appHttpHealth := new(v1beta1.AppHttpHealthy)
		appHttpHealth.Name = appHttpHealthName

		// agent
		agentSpec := new(v1beta1.AgentSpec)
		agentSpec.TerminationGracePeriodMinutes = &termMin
		appHttpHealth.Spec.AgentSpec = agentSpec

		// successCondition
		successCondition := new(v1beta1.NetSuccessCondition)
		successCondition.SuccessRate = &successRate
		successCondition.MeanAccessDelayInMs = &successMean
		appHttpHealth.Spec.SuccessCondition = successCondition

		// target
		target := new(v1beta1.AppHttpHealthyTarget)
		target.Method = "GET"
		target.Host = "www.baidu.com"
		appHttpHealth.Spec.Target = target

		// request
		request := new(v1beta1.NetHttpRequest)
		request.PerRequestTimeoutInMS = 2000
		request.QPS = 10
		request.DurationInSecond = 10
		appHttpHealth.Spec.Request = request

		// Schedule
		Schedule := new(v1beta1.SchedulePlan)
		Schedule.Schedule = &crontab
		Schedule.RoundNumber = 1
		Schedule.RoundTimeoutMinute = 1
		appHttpHealth.Spec.Schedule = Schedule

		e = frame.CreateResource(appHttpHealth)
		Expect(e).NotTo(HaveOccurred(), "create appHttpHealth resource")

		e = common.CheckRuntime(frame, appHttpHealth, pluginManager.KindNameAppHttpHealthy, 60)
		Expect(e).NotTo(HaveOccurred(), "check task runtime spec")

		// get runtime resource
		fake := &v1beta1.AppHttpHealthy{
			ObjectMeta: metav1.ObjectMeta{
				Name: appHttpHealthName,
			},
		}
		key := client.ObjectKeyFromObject(fake)
		rs := &v1beta1.AppHttpHealthy{}
		e = frame.GetResource(key, rs)
		Expect(e).NotTo(HaveOccurred(), "get task runtime status")

		e = frame.DeleteResource(rs)
		Expect(e).NotTo(HaveOccurred(), "delete task")
		// check runtime resource is deleted after deleting task
		e = common.GetRuntimeResource(frame, rs.Status.Resource, false)
		Expect(e).NotTo(HaveOccurred(), "get task runtime resource")
	})

	It("Successfully testing cascading deletion with Task NetDns DaemonSet Service  ", Label("E00009"), func() {
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
		ServiceName := "kube-dns"
		ServiceNamespace := "kube-system"
		targetDns.ServiceName = &ServiceName
		targetDns.ServiceNamespace = &ServiceNamespace
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
		request.Domain = "www.baidu.com"
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

		// get runtime resource
		fake := &v1beta1.Netdns{
			ObjectMeta: metav1.ObjectMeta{
				Name: netDnsName,
			},
		}
		key := client.ObjectKeyFromObject(fake)
		rs := &v1beta1.Netdns{}
		e = frame.GetResource(key, rs)
		Expect(e).NotTo(HaveOccurred(), "get task runtime status")

		e = frame.DeleteResource(rs)
		Expect(e).NotTo(HaveOccurred(), "delete task")
		// check runtime resource is deleted after deleting task
		e = common.GetRuntimeResource(frame, rs.Status.Resource, false)
		Expect(e).NotTo(HaveOccurred(), "get task runtime resource")

	})

	It("Successfully testing cascading deletion with Task NetReach Deployment Service and Ingress ", Label("E00010"), func() {
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
		agentSpec.Kind = types.KindDeployment
		agentSpec.DeploymentReplicas = &replicas
		netReach.Spec.AgentSpec = agentSpec

		// successCondition
		successCondition := new(v1beta1.NetSuccessCondition)
		successCondition.SuccessRate = &successRate
		successCondition.MeanAccessDelayInMs = &successMean
		netReach.Spec.SuccessCondition = successCondition

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

		e = common.CheckRuntime(frame, netReach, pluginManager.KindNameNetReach, 60)
		Expect(e).NotTo(HaveOccurred(), "check task runtime spec")

		// get runtime resource
		fake := &v1beta1.NetReach{
			ObjectMeta: metav1.ObjectMeta{
				Name: netReachName,
			},
		}
		key := client.ObjectKeyFromObject(fake)
		rs := &v1beta1.NetReach{}
		e = frame.GetResource(key, rs)
		Expect(e).NotTo(HaveOccurred(), "get task runtime status")

		e = frame.DeleteResource(rs)
		Expect(e).NotTo(HaveOccurred(), "delete task")
		// check runtime resource is deleted after deleting task
		e = common.GetRuntimeResource(frame, rs.Status.Resource, true)
		Expect(e).NotTo(HaveOccurred(), "get task runtime resource")
	})

	It("Successfully testing cascading deletion with Task NetAppHttpHealthy Deployment Service", Label("E00011"), func() {
		var e error
		successRate := float64(1)
		successMean := int64(1500)
		crontab := "0 1"
		appHttpHealthName := "apphttphealth-get" + tools.RandomName()

		appHttpHealth := new(v1beta1.AppHttpHealthy)
		appHttpHealth.Name = appHttpHealthName

		// agent
		agentSpec := new(v1beta1.AgentSpec)
		agentSpec.TerminationGracePeriodMinutes = &termMin
		agentSpec.Kind = types.KindDeployment
		agentSpec.DeploymentReplicas = &replicas
		appHttpHealth.Spec.AgentSpec = agentSpec

		// successCondition
		successCondition := new(v1beta1.NetSuccessCondition)
		successCondition.SuccessRate = &successRate
		successCondition.MeanAccessDelayInMs = &successMean
		appHttpHealth.Spec.SuccessCondition = successCondition

		// target
		target := new(v1beta1.AppHttpHealthyTarget)
		target.Method = "GET"
		target.Host = "www.baidu.com"
		appHttpHealth.Spec.Target = target

		// request
		request := new(v1beta1.NetHttpRequest)
		request.PerRequestTimeoutInMS = 2000
		request.QPS = 10
		request.DurationInSecond = 10
		appHttpHealth.Spec.Request = request

		// Schedule
		Schedule := new(v1beta1.SchedulePlan)
		Schedule.Schedule = &crontab
		Schedule.RoundNumber = 1
		Schedule.RoundTimeoutMinute = 1
		appHttpHealth.Spec.Schedule = Schedule

		e = frame.CreateResource(appHttpHealth)
		Expect(e).NotTo(HaveOccurred(), "create appHttpHealth resource")

		e = common.CheckRuntime(frame, appHttpHealth, pluginManager.KindNameAppHttpHealthy, 60)
		Expect(e).NotTo(HaveOccurred(), "check task runtime spec")

		// get runtime resource
		fake := &v1beta1.AppHttpHealthy{
			ObjectMeta: metav1.ObjectMeta{
				Name: appHttpHealthName,
			},
		}
		key := client.ObjectKeyFromObject(fake)
		rs := &v1beta1.AppHttpHealthy{}
		e = frame.GetResource(key, rs)
		Expect(e).NotTo(HaveOccurred(), "get task runtime status")

		e = frame.DeleteResource(rs)
		Expect(e).NotTo(HaveOccurred(), "delete task")
		// check runtime resource is deleted after deleting task
		e = common.GetRuntimeResource(frame, rs.Status.Resource, false)
		Expect(e).NotTo(HaveOccurred(), "get task runtime resource")
	})

	It("Successfully testing cascading deletion with Task NetDns Deployment Service  ", Label("E00012"), func() {
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
		agentSpec.Kind = types.KindDeployment
		agentSpec.DeploymentReplicas = &replicas
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
		ServiceName := "kube-dns"
		ServiceNamespace := "kube-system"
		targetDns.ServiceName = &ServiceName
		targetDns.ServiceNamespace = &ServiceNamespace
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
		request.Domain = "www.baidu.com"
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

		// get runtime resource
		fake := &v1beta1.Netdns{
			ObjectMeta: metav1.ObjectMeta{
				Name: netDnsName,
			},
		}
		key := client.ObjectKeyFromObject(fake)
		rs := &v1beta1.Netdns{}
		e = frame.GetResource(key, rs)
		Expect(e).NotTo(HaveOccurred(), "get task runtime status")

		e = frame.DeleteResource(rs)
		Expect(e).NotTo(HaveOccurred(), "delete task")
		// check runtime resource is deleted after deleting task
		e = common.GetRuntimeResource(frame, rs.Status.Resource, false)
		Expect(e).NotTo(HaveOccurred(), "get task runtime resource")

	})

})
