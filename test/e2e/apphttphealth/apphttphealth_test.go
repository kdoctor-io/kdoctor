// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package apphttphealth_test

import (
	"fmt"
	"github.com/kdoctor-io/kdoctor/pkg/k8s/apis/kdoctor.io/v1beta1"
	"github.com/kdoctor-io/kdoctor/pkg/pluginManager"
	"github.com/kdoctor-io/kdoctor/test/e2e/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spidernet-io/e2eframework/tools"
	"net"
)

var _ = Describe("testing appHttpHealth test ", Label("appHttpHealth"), func() {
	var termMin = int64(3)

	It("success http testing appHttpHealth method GET", Label("A00001", "A00011", "C00006"), func() {
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
		appHttpHealth.Spec.AgentSpec = *agentSpec

		// successCondition
		successCondition := new(v1beta1.NetSuccessCondition)
		successCondition.SuccessRate = &successRate
		successCondition.MeanAccessDelayInMs = &successMean
		appHttpHealth.Spec.SuccessCondition = successCondition

		// target
		target := new(v1beta1.AppHttpHealthyTarget)
		target.Method = "GET"
		if net.ParseIP(testSvcIP).To4() == nil {
			target.Host = fmt.Sprintf("http://[%s]:%d/?task=%s", testSvcIP, httpPort, appHttpHealthName)
		} else {
			target.Host = fmt.Sprintf("http://%s:%d?task=%s", testSvcIP, httpPort, appHttpHealthName)
		}
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

		e = common.WaitKdoctorTaskDone(frame, appHttpHealth, pluginManager.KindNameAppHttpHealthy, 120)
		Expect(e).NotTo(HaveOccurred(), "wait appHttpHealth task finish")

		success, e := common.CompareResult(frame, appHttpHealthName, pluginManager.KindNameAppHttpHealthy, testPodIPs, reportNum)
		Expect(e).NotTo(HaveOccurred(), "compare report and task")
		Expect(success).To(BeTrue(), "compare report and task result")
	})

	It("failed http testing appHttpHealth due to status code", Label("A00002"), func() {
		var e error
		successRate := float64(1)
		successMean := int64(1500)
		crontab := "0 1"
		expectStatusCode := 205
		appHttpHealthName := "apphttphealth-get" + tools.RandomName()

		appHttpHealth := new(v1beta1.AppHttpHealthy)
		appHttpHealth.Name = appHttpHealthName

		// agentSpec
		agentSpec := new(v1beta1.AgentSpec)
		agentSpec.TerminationGracePeriodMinutes = &termMin
		appHttpHealth.Spec.AgentSpec = *agentSpec

		// successCondition
		successCondition := new(v1beta1.NetSuccessCondition)
		successCondition.SuccessRate = &successRate
		successCondition.MeanAccessDelayInMs = &successMean
		successCondition.StatusCode = &expectStatusCode
		appHttpHealth.Spec.SuccessCondition = successCondition

		// target
		target := new(v1beta1.AppHttpHealthyTarget)
		target.Method = "GET"
		if net.ParseIP(testSvcIP).To4() == nil {
			target.Host = fmt.Sprintf("http://[%s]:%d/?task=%s", testSvcIP, httpPort, appHttpHealthName)
		} else {
			target.Host = fmt.Sprintf("http://%s:%d/?task=%s", testSvcIP, httpPort, appHttpHealthName)
		}
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

		e = common.WaitKdoctorTaskDone(frame, appHttpHealth, pluginManager.KindNameAppHttpHealthy, 120)
		Expect(e).NotTo(HaveOccurred(), "wait appHttpHealth task finish")

		success, e := common.CompareResult(frame, appHttpHealthName, pluginManager.KindNameAppHttpHealthy, testPodIPs, reportNum)
		Expect(e).NotTo(HaveOccurred(), "compare report and task")
		Expect(success).To(BeFalse(), "compare report and task result")

	})

	It("Failed http testing appHttpHealth due to delay ", Label("A00003"), func() {
		var e error
		successRate := float64(1)
		successMean := int64(200)
		crontab := "0 1"
		appHttpHealthName := "apphttphealth-get" + tools.RandomName()

		appHttpHealth := new(v1beta1.AppHttpHealthy)
		appHttpHealth.Name = appHttpHealthName

		// agentSpec
		agentSpec := new(v1beta1.AgentSpec)
		agentSpec.TerminationGracePeriodMinutes = &termMin
		appHttpHealth.Spec.AgentSpec = *agentSpec

		// successCondition
		successCondition := new(v1beta1.NetSuccessCondition)
		successCondition.SuccessRate = &successRate
		successCondition.MeanAccessDelayInMs = &successMean
		appHttpHealth.Spec.SuccessCondition = successCondition

		// target
		target := new(v1beta1.AppHttpHealthyTarget)
		target.Method = "GET"
		if net.ParseIP(testSvcIP).To4() == nil {
			target.Host = fmt.Sprintf("http://[%s]:%d/?delay=1&task=%s", testSvcIP, httpPort, appHttpHealthName)
		} else {
			target.Host = fmt.Sprintf("http://%s:%d/?delay=1&task=%s", testSvcIP, httpPort, appHttpHealthName)
		}
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

		e = common.WaitKdoctorTaskDone(frame, appHttpHealth, pluginManager.KindNameAppHttpHealthy, 120)
		Expect(e).NotTo(HaveOccurred(), "wait appHttpHealth task finish")

		success, e := common.CompareResult(frame, appHttpHealthName, pluginManager.KindNameAppHttpHealthy, testPodIPs, reportNum)
		Expect(e).NotTo(HaveOccurred(), "compare report and task")
		Expect(success).To(BeFalse(), "compare report and task result")

	})

	It("success https testing appHttpHealth method GET", Label("A00004"), func() {
		var e error
		successRate := float64(1)
		successMean := int64(1500)
		crontab := "0 1"
		appHttpHealthName := "apphttphealth-get" + tools.RandomName()
		appHttpHealth := new(v1beta1.AppHttpHealthy)
		appHttpHealth.Name = appHttpHealthName

		// agentSpec
		agentSpec := new(v1beta1.AgentSpec)
		agentSpec.TerminationGracePeriodMinutes = &termMin
		appHttpHealth.Spec.AgentSpec = *agentSpec

		// successCondition
		successCondition := new(v1beta1.NetSuccessCondition)
		successCondition.SuccessRate = &successRate
		successCondition.MeanAccessDelayInMs = &successMean
		appHttpHealth.Spec.SuccessCondition = successCondition

		// target
		target := new(v1beta1.AppHttpHealthyTarget)
		target.Method = "GET"
		if net.ParseIP(testPodIPs[0]).To4() == nil {
			target.Host = fmt.Sprintf("https://[%s]/?task=%s", testPodIPs[0], appHttpHealthName)
		} else {
			target.Host = fmt.Sprintf("https://%s/?task=%s", testPodIPs[0], appHttpHealthName)
		}
		target.TlsSecretName = &common.TlsClientName
		target.TlsSecretNamespace = &testAppNamespace
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

		e = common.WaitKdoctorTaskDone(frame, appHttpHealth, pluginManager.KindNameAppHttpHealthy, 120)
		Expect(e).NotTo(HaveOccurred(), "wait appHttpHealth task finish")

		success, e := common.CompareResult(frame, appHttpHealthName, pluginManager.KindNameAppHttpHealthy, testPodIPs, reportNum)
		Expect(e).NotTo(HaveOccurred(), "compare report and task")
		Expect(success).To(BeTrue(), "compare report and task result")

	})

	It("failed https testing appHttpHealth due to tls", Label("A00005"), func() {
		var e error
		successRate := float64(1)
		successMean := int64(1500)
		crontab := "0 1"
		appHttpHealthName := "apphttphealth-get" + tools.RandomName()

		appHttpHealth := new(v1beta1.AppHttpHealthy)
		appHttpHealth.Name = appHttpHealthName

		// agentSpec
		agentSpec := new(v1beta1.AgentSpec)
		agentSpec.TerminationGracePeriodMinutes = &termMin
		appHttpHealth.Spec.AgentSpec = *agentSpec

		// successCondition
		successCondition := new(v1beta1.NetSuccessCondition)
		successCondition.SuccessRate = &successRate
		successCondition.MeanAccessDelayInMs = &successMean
		appHttpHealth.Spec.SuccessCondition = successCondition

		// target
		target := new(v1beta1.AppHttpHealthyTarget)
		target.Method = "GET"
		// The service IP is not signed in the test pod, so using the service IP will fail certificate validation
		if net.ParseIP(testSvcIP).To4() == nil {
			target.Host = fmt.Sprintf("https://[%s]:%d/?task=%s", testSvcIP, httpsPort, appHttpHealthName)
		} else {
			target.Host = fmt.Sprintf("https://%s:%d/?task=%s", testSvcIP, httpsPort, appHttpHealthName)
		}
		target.TlsSecretName = &common.TlsClientName
		target.TlsSecretNamespace = &testAppNamespace
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

		e = common.WaitKdoctorTaskDone(frame, appHttpHealth, pluginManager.KindNameAppHttpHealthy, 120)
		Expect(e).NotTo(HaveOccurred(), "wait appHttpHealth task finish")

		success, e := common.CompareResult(frame, appHttpHealthName, pluginManager.KindNameAppHttpHealthy, []string{}, reportNum)
		Expect(e).NotTo(HaveOccurred(), "compare report and task")
		Expect(success).To(BeFalse(), "compare report and task result")

	})

	It("Successfully http testing appHttpHealth method PUT ", Label("A00006"), func() {
		var e error
		successRate := float64(1)
		successMean := int64(1500)
		crontab := "0 1"
		appHttpHealthName := "apphttphealth-put" + tools.RandomName()

		appHttpHealth := new(v1beta1.AppHttpHealthy)
		appHttpHealth.Name = appHttpHealthName

		// agentSpec
		agentSpec := new(v1beta1.AgentSpec)
		agentSpec.TerminationGracePeriodMinutes = &termMin
		appHttpHealth.Spec.AgentSpec = *agentSpec

		// successCondition
		successCondition := new(v1beta1.NetSuccessCondition)
		successCondition.SuccessRate = &successRate
		successCondition.MeanAccessDelayInMs = &successMean
		appHttpHealth.Spec.SuccessCondition = successCondition

		// target
		target := new(v1beta1.AppHttpHealthyTarget)
		target.Method = "PUT"
		if net.ParseIP(testSvcIP).To4() == nil {
			target.Host = fmt.Sprintf("http://[%s]:%d/?task=%s", testSvcIP, httpPort, appHttpHealthName)
		} else {
			target.Host = fmt.Sprintf("http://%s:%d/?task=%s", testSvcIP, httpPort, appHttpHealthName)
		}
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

		e = common.WaitKdoctorTaskDone(frame, appHttpHealth, pluginManager.KindNameAppHttpHealthy, 120)
		Expect(e).NotTo(HaveOccurred(), "wait appHttpHealth task finish")

		success, e := common.CompareResult(frame, appHttpHealthName, pluginManager.KindNameAppHttpHealthy, testPodIPs, reportNum)
		Expect(e).NotTo(HaveOccurred(), "compare report and task")
		Expect(success).To(BeTrue(), "compare report and task result")

	})

	It("Successfully http testing appHttpHealth method POST With Body", Label("A00007"), func() {
		var e error
		successRate := float64(1)
		successMean := int64(1500)
		crontab := "0 1"
		appHttpHealthName := "apphttphealth-post" + tools.RandomName()

		appHttpHealth := new(v1beta1.AppHttpHealthy)
		appHttpHealth.Name = appHttpHealthName

		// agentSpec
		agentSpec := new(v1beta1.AgentSpec)
		agentSpec.TerminationGracePeriodMinutes = &termMin
		appHttpHealth.Spec.AgentSpec = *agentSpec

		// successCondition
		successCondition := new(v1beta1.NetSuccessCondition)
		successCondition.SuccessRate = &successRate
		successCondition.MeanAccessDelayInMs = &successMean
		appHttpHealth.Spec.SuccessCondition = successCondition

		// target
		target := new(v1beta1.AppHttpHealthyTarget)
		target.Method = "POST"
		if net.ParseIP(testSvcIP).To4() == nil {
			target.Host = fmt.Sprintf("http://[%s]:%d/?task=%s", testSvcIP, httpPort, appHttpHealthName)
		} else {
			target.Host = fmt.Sprintf("http://%s:%d/?task=%s", testSvcIP, httpPort, appHttpHealthName)
		}
		target.BodyConfigName = &bodyConfigMapName
		target.BodyConfigNamespace = &common.TestNameSpace
		target.Header = []string{"Content-Type: application/json"}
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

		e = common.WaitKdoctorTaskDone(frame, appHttpHealth, pluginManager.KindNameAppHttpHealthy, 120)
		Expect(e).NotTo(HaveOccurred(), "wait appHttpHealth task finish")

		success, e := common.CompareResult(frame, appHttpHealthName, pluginManager.KindNameAppHttpHealthy, testPodIPs, reportNum)
		Expect(e).NotTo(HaveOccurred(), "compare report and task")
		Expect(success).To(BeTrue(), "compare report and task result")
	})

	It("Successfully http testing appHttpHealth method HEAD", Label("A00008"), func() {
		var e error
		successRate := float64(1)
		successMean := int64(1500)
		crontab := "0 1"
		appHttpHealthName := "apphttphealth-head" + tools.RandomName()

		appHttpHealth := new(v1beta1.AppHttpHealthy)
		appHttpHealth.Name = appHttpHealthName

		// agentSpec
		agentSpec := new(v1beta1.AgentSpec)
		agentSpec.TerminationGracePeriodMinutes = &termMin
		appHttpHealth.Spec.AgentSpec = *agentSpec

		// successCondition
		successCondition := new(v1beta1.NetSuccessCondition)
		successCondition.SuccessRate = &successRate
		successCondition.MeanAccessDelayInMs = &successMean
		appHttpHealth.Spec.SuccessCondition = successCondition

		// target
		target := new(v1beta1.AppHttpHealthyTarget)
		target.Method = "HEAD"
		if net.ParseIP(testSvcIP).To4() == nil {
			target.Host = fmt.Sprintf("http://[%s]:%d/?task=%s", testSvcIP, httpPort, appHttpHealthName)
		} else {
			target.Host = fmt.Sprintf("http://%s:%d/?task=%s", testSvcIP, httpPort, appHttpHealthName)
		}
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

		e = common.WaitKdoctorTaskDone(frame, appHttpHealth, pluginManager.KindNameAppHttpHealthy, 120)
		Expect(e).NotTo(HaveOccurred(), "wait appHttpHealth task finish")

		success, e := common.CompareResult(frame, appHttpHealthName, pluginManager.KindNameAppHttpHealthy, testPodIPs, reportNum)
		Expect(e).NotTo(HaveOccurred(), "compare report and task")
		Expect(success).To(BeTrue(), "compare report and task result")

	})

	It("Successfully http testing appHttpHealth method PATCH", Label("A00009"), func() {
		var e error
		successRate := float64(1)
		successMean := int64(1500)
		crontab := "0 1"
		appHttpHealthName := "apphttphealth-patch" + tools.RandomName()

		appHttpHealth := new(v1beta1.AppHttpHealthy)
		appHttpHealth.Name = appHttpHealthName

		// agentSpec
		agentSpec := new(v1beta1.AgentSpec)
		agentSpec.TerminationGracePeriodMinutes = &termMin
		appHttpHealth.Spec.AgentSpec = *agentSpec

		// successCondition
		successCondition := new(v1beta1.NetSuccessCondition)
		successCondition.SuccessRate = &successRate
		successCondition.MeanAccessDelayInMs = &successMean
		appHttpHealth.Spec.SuccessCondition = successCondition

		// target
		target := new(v1beta1.AppHttpHealthyTarget)
		target.Method = "PATCH"
		if net.ParseIP(testSvcIP).To4() == nil {
			target.Host = fmt.Sprintf("http://[%s]:%d/?task=%s", testSvcIP, httpPort, appHttpHealthName)
		} else {
			target.Host = fmt.Sprintf("http://%s:%d/?task=%s", testSvcIP, httpPort, appHttpHealthName)
		}
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

		e = common.WaitKdoctorTaskDone(frame, appHttpHealth, pluginManager.KindNameAppHttpHealthy, 120)
		Expect(e).NotTo(HaveOccurred(), "wait appHttpHealth task finish")

		success, e := common.CompareResult(frame, appHttpHealthName, pluginManager.KindNameAppHttpHealthy, testPodIPs, reportNum)
		Expect(e).NotTo(HaveOccurred(), "compare report and task")
		Expect(success).To(BeTrue(), "compare report and task result")

	})

	It("Successfully http testing appHttpHealth method OPTIONS", Label("A00010"), func() {
		var e error
		successRate := float64(1)
		successMean := int64(1500)
		crontab := "0 1"
		appHttpHealthName := "apphttphealth-options" + tools.RandomName()

		appHttpHealth := new(v1beta1.AppHttpHealthy)
		appHttpHealth.Name = appHttpHealthName

		// agentSpec
		agentSpec := new(v1beta1.AgentSpec)
		agentSpec.TerminationGracePeriodMinutes = &termMin
		appHttpHealth.Spec.AgentSpec = *agentSpec

		// successCondition
		successCondition := new(v1beta1.NetSuccessCondition)
		successCondition.SuccessRate = &successRate
		successCondition.MeanAccessDelayInMs = &successMean
		appHttpHealth.Spec.SuccessCondition = successCondition

		// target
		target := new(v1beta1.AppHttpHealthyTarget)
		target.Method = "OPTIONS"
		if net.ParseIP(testSvcIP).To4() == nil {
			target.Host = fmt.Sprintf("http://[%s]:%d/?task=%s", testSvcIP, httpPort, appHttpHealthName)
		} else {
			target.Host = fmt.Sprintf("http://%s:%d/?task=%s", testSvcIP, httpPort, appHttpHealthName)
		}
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

		e = common.WaitKdoctorTaskDone(frame, appHttpHealth, pluginManager.KindNameAppHttpHealthy, 120)
		Expect(e).NotTo(HaveOccurred(), "wait appHttpHealth task finish")

		success, e := common.CompareResult(frame, appHttpHealthName, pluginManager.KindNameAppHttpHealthy, testPodIPs, reportNum)
		Expect(e).NotTo(HaveOccurred(), "compare report and task")
		Expect(success).To(BeTrue(), "compare report and task result")

	})

	It("Successfully http testing appHttpHealth due to success rate", Label("A00012"), func() {
		var e error
		successRate := float64(0.2)
		successMean := int64(1200)
		crontab := "0 1"
		appHttpHealthName := "apphttphealth-get" + tools.RandomName()

		appHttpHealth := new(v1beta1.AppHttpHealthy)
		appHttpHealth.Name = appHttpHealthName

		// agentSpec
		agentSpec := new(v1beta1.AgentSpec)
		agentSpec.TerminationGracePeriodMinutes = &termMin
		appHttpHealth.Spec.AgentSpec = *agentSpec

		// successCondition
		successCondition := new(v1beta1.NetSuccessCondition)
		successCondition.SuccessRate = &successRate
		successCondition.MeanAccessDelayInMs = &successMean
		appHttpHealth.Spec.SuccessCondition = successCondition

		// target
		target := new(v1beta1.AppHttpHealthyTarget)
		target.Method = "GET"
		if net.ParseIP(testSvcIP).To4() == nil {
			target.Host = fmt.Sprintf("http://[%s]:%d/?delay=1&task=%s", testSvcIP, httpPort, appHttpHealthName)
		} else {
			target.Host = fmt.Sprintf("http://%s:%d/?delay=1&task=%s", testSvcIP, httpPort, appHttpHealthName)
		}
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

		e = common.WaitKdoctorTaskDone(frame, appHttpHealth, pluginManager.KindNameAppHttpHealthy, 120)
		Expect(e).NotTo(HaveOccurred(), "wait appHttpHealth task finish")

		success, e := common.CompareResult(frame, appHttpHealthName, pluginManager.KindNameAppHttpHealthy, testPodIPs, reportNum)
		Expect(e).NotTo(HaveOccurred(), "compare report and task")
		Expect(success).To(BeTrue(), "compare report and task result")

	})

	It("Successfully https testing appHttpHealth method GET Protocol Http2", Label("A00013"), func() {
		var e error
		successRate := float64(1)
		successMean := int64(1500)
		crontab := "0 1"
		appHttpHealthName := "apphttphealth-get" + tools.RandomName()

		appHttpHealth := new(v1beta1.AppHttpHealthy)
		appHttpHealth.Name = appHttpHealthName

		// agentSpec
		agentSpec := new(v1beta1.AgentSpec)
		agentSpec.TerminationGracePeriodMinutes = &termMin
		appHttpHealth.Spec.AgentSpec = *agentSpec

		// successCondition
		successCondition := new(v1beta1.NetSuccessCondition)
		successCondition.SuccessRate = &successRate
		successCondition.MeanAccessDelayInMs = &successMean
		appHttpHealth.Spec.SuccessCondition = successCondition

		// target
		target := new(v1beta1.AppHttpHealthyTarget)
		target.Method = "GET"
		if net.ParseIP(testPodIPs[0]).To4() == nil {
			target.Host = fmt.Sprintf("https://[%s]/?task=%s", testPodIPs[0], appHttpHealthName)
		} else {
			target.Host = fmt.Sprintf("https://%s/?task=%s", testPodIPs[0], appHttpHealthName)
		}
		target.TlsSecretName = &common.TlsClientName
		target.TlsSecretNamespace = &testAppNamespace
		target.Http2 = true
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

		e = common.WaitKdoctorTaskDone(frame, appHttpHealth, pluginManager.KindNameAppHttpHealthy, 120)
		Expect(e).NotTo(HaveOccurred(), "wait appHttpHealth task finish")

		success, e := common.CompareResult(frame, appHttpHealthName, pluginManager.KindNameAppHttpHealthy, testPodIPs, reportNum)
		Expect(e).NotTo(HaveOccurred(), "compare report and task")
		Expect(success).To(BeTrue(), "compare report and task result")

	})
})
