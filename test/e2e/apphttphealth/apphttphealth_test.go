// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package apphttphealth_test

import (
	"context"
	"fmt"
	"github.com/kdoctor-io/kdoctor/pkg/k8s/apis/kdoctor.io/v1beta1"
	"github.com/kdoctor-io/kdoctor/pkg/pluginManager"
	"github.com/kdoctor-io/kdoctor/test/e2e/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spidernet-io/e2eframework/tools"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net"
	"time"
)

var _ = Describe("testing appHttpHealth test ", Label("appHttpHealth"), func() {
	var testSvcIP string
	var testPodIPs []string

	BeforeEach(func() {
		var e error

		TestDsName := "kdoctor-app-" + tools.RandomName()
		TestDsNamespace := "kdoctor-ns" + tools.RandomName()

		e = frame.CreateNamespaceUntilDefaultServiceAccountReady(TestDsNamespace, common.ServiceAccountReadyTimeout)
		Expect(e).NotTo(HaveOccurred(), "create test namespace")
		// copy kdoctor ca cert
		copyCaSecret := new(v1.Secret)
		copyCaSecret.Name = caSecret.Name
		copyCaSecret.Namespace = TestDsNamespace
		copyCaSecret.Type = caSecret.Type
		copyCaSecret.Data = caSecret.Data

		e = frame.KClient.Create(context.Background(), copyCaSecret)
		Expect(e).NotTo(HaveOccurred(), "copy kdoctor ca in test namespace")

		// create app daemonSet
		commonKey := metav1.ObjectMeta{
			Name:      TestDsName,
			Namespace: TestDsNamespace,
			Labels:    map[string]string{"app": TestDsName},
		}
		e = common.CreateApplicationDaemonSetAndService(frame, commonKey, true, false)
		Expect(e).NotTo(HaveOccurred(), "create test daemonSet and service")

		// wait test ds ready
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
		defer cancel()
		ds, e := frame.WaitDaemonSetReady(TestDsName, TestDsNamespace, ctx)
		Expect(e).NotTo(HaveOccurred(), "wait test daemonSet and service ready")
		// wait pod running
		ctx2, cancel2 := context.WithTimeout(context.Background(), time.Second*60)
		defer cancel2()

		e = frame.WaitPodListRunning(ds.Labels, int(ds.Status.NumberReady), ctx2)
		Expect(e).NotTo(HaveOccurred(), "wait test pod running")
		time.Sleep(time.Second * 10)
		//  get test app service ip and pod ip
		commonSvc, e := frame.GetService(TestDsName, TestDsNamespace)
		Expect(e).NotTo(HaveOccurred(), "get kdoctor test service")
		testSvcIP = commonSvc.Spec.ClusterIP
		GinkgoWriter.Printf("get common service ip %v \n", testSvcIP)

		commonDs, e := frame.GetDaemonSet(TestDsName, TestDsNamespace)
		Expect(e).NotTo(HaveOccurred(), "get kdoctor test daemonset")

		commonPods, e := frame.GetDaemonSetPodList(commonDs)
		Expect(e).NotTo(HaveOccurred(), "get kdoctor test pod list")

		testPodIPs = make([]string, 0)
		for _, v := range commonPods.Items {
			testPodIPs = append(testPodIPs, v.Status.PodIP)
		}

		GinkgoWriter.Printf("get test pod ips %v \n", testPodIPs)

		// Clean test env
		DeferCleanup(func() {
			GinkgoWriter.Printf("delete namespace %v \n", TestDsNamespace)
			Expect(frame.DeleteNamespace(TestDsNamespace)).NotTo(HaveOccurred())
		})

	})

	It("success http testing appHttpHealth method GET", Label("A00001"), func() {
		var e error
		successRate := float64(1)
		successMean := int64(1500)
		crontab := "0 1"
		appHttpHealthName := "apphttphealth-get" + tools.RandomName()

		appHttpHealth := new(v1beta1.AppHttpHealthy)
		appHttpHealth.Name = appHttpHealthName

		// successCondition
		successCondition := new(v1beta1.NetSuccessCondition)
		successCondition.SuccessRate = &successRate
		successCondition.MeanAccessDelayInMs = &successMean
		appHttpHealth.Spec.SuccessCondition = successCondition

		// target
		target := new(v1beta1.AppHttpHealthyTarget)
		target.Method = "GET"
		if net.ParseIP(testSvcIP).To4() == nil {
			target.Host = fmt.Sprintf("http://[%s]:%d", testSvcIP, httpPort)
		} else {
			target.Host = fmt.Sprintf("http://%s:%d", testSvcIP, httpPort)
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
		Expect(success).To(BeTrue(), "compare report and task result")
		Expect(e).NotTo(HaveOccurred(), "compare report and task")
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
			target.Host = fmt.Sprintf("http://[%s]:%d", testSvcIP, httpPort)
		} else {
			target.Host = fmt.Sprintf("http://%s:%d", testSvcIP, httpPort)
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
		Expect(success).To(BeFalse(), "compare report and task result")
		Expect(e).NotTo(HaveOccurred(), "compare report and task")

	})

	It("Failed http testing appHttpHealth due to delay ", Label("A00003"), func() {
		var e error
		successRate := float64(1)
		successMean := int64(200)
		crontab := "0 1"
		appHttpHealthName := "apphttphealth-get" + tools.RandomName()

		appHttpHealth := new(v1beta1.AppHttpHealthy)
		appHttpHealth.Name = appHttpHealthName

		// successCondition
		successCondition := new(v1beta1.NetSuccessCondition)
		successCondition.SuccessRate = &successRate
		successCondition.MeanAccessDelayInMs = &successMean
		appHttpHealth.Spec.SuccessCondition = successCondition

		// target
		target := new(v1beta1.AppHttpHealthyTarget)
		target.Method = "GET"
		if net.ParseIP(testSvcIP).To4() == nil {
			target.Host = fmt.Sprintf("http://[%s]:%d/?delay=1", testSvcIP, httpPort)
		} else {
			target.Host = fmt.Sprintf("http://%s:%d/?delay=1", testSvcIP, httpPort)
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
		Expect(success).To(BeFalse(), "compare report and task result")
		Expect(e).NotTo(HaveOccurred(), "compare report and task")

	})

	It("success https testing appHttpHealth method GET", Label("A00004"), func() {
		var e error
		successRate := float64(1)
		successMean := int64(1500)
		crontab := "0 1"
		appHttpHealthName := "apphttphealth-get" + tools.RandomName()
		appHttpHealth := new(v1beta1.AppHttpHealthy)
		appHttpHealth.Name = appHttpHealthName

		// successCondition
		successCondition := new(v1beta1.NetSuccessCondition)
		successCondition.SuccessRate = &successRate
		successCondition.MeanAccessDelayInMs = &successMean
		appHttpHealth.Spec.SuccessCondition = successCondition

		// target
		target := new(v1beta1.AppHttpHealthyTarget)
		target.Method = "GET"
		if net.ParseIP(testPodIPs[0]).To4() == nil {
			target.Host = fmt.Sprintf("https://[%s]", testPodIPs[0])
		} else {
			target.Host = fmt.Sprintf("https://%s", testPodIPs[0])
		}
		target.TlsSecretName = &common.TlsClientName
		target.TlsSecretNamespace = &common.TestNameSpace
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
		Expect(success).To(BeTrue(), "compare report and task result")
		Expect(e).NotTo(HaveOccurred(), "compare report and task")

	})

	It("failed https testing appHttpHealth due to tls", Label("A00005"), func() {
		var e error
		successRate := float64(1)
		successMean := int64(1500)
		crontab := "0 1"
		appHttpHealthName := "apphttphealth-get" + tools.RandomName()

		appHttpHealth := new(v1beta1.AppHttpHealthy)
		appHttpHealth.Name = appHttpHealthName

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
			target.Host = fmt.Sprintf("https://[%s]:%d", testSvcIP, httpsPort)
		} else {
			target.Host = fmt.Sprintf("https://%s:%d", testSvcIP, httpsPort)
		}

		target.TlsSecretName = &common.TlsClientName
		target.TlsSecretNamespace = &common.TestNameSpace
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
		Expect(success).To(BeFalse(), "compare report and task result")
		Expect(e).NotTo(HaveOccurred(), "compare report and task")

	})

	It("Successfully http testing appHttpHealth method PUT ", Label("A00006"), func() {
		var e error
		successRate := float64(1)
		successMean := int64(1500)
		crontab := "0 1"
		appHttpHealthName := "apphttphealth-put" + tools.RandomName()

		appHttpHealth := new(v1beta1.AppHttpHealthy)
		appHttpHealth.Name = appHttpHealthName

		// successCondition
		successCondition := new(v1beta1.NetSuccessCondition)
		successCondition.SuccessRate = &successRate
		successCondition.MeanAccessDelayInMs = &successMean
		appHttpHealth.Spec.SuccessCondition = successCondition

		// target
		target := new(v1beta1.AppHttpHealthyTarget)
		target.Method = "PUT"
		if net.ParseIP(testSvcIP).To4() == nil {
			target.Host = fmt.Sprintf("http://[%s]:%d", testSvcIP, httpPort)
		} else {
			target.Host = fmt.Sprintf("http://%s:%d", testSvcIP, httpPort)
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
		Expect(success).To(BeTrue(), "compare report and task result")
		Expect(e).NotTo(HaveOccurred(), "compare report and task")

	})

	It("Successfully http testing appHttpHealth method POST With Body", Label("A00007"), func() {
		var e error
		successRate := float64(1)
		successMean := int64(1500)
		crontab := "0 1"
		appHttpHealthName := "apphttphealth-post" + tools.RandomName()

		appHttpHealth := new(v1beta1.AppHttpHealthy)
		appHttpHealth.Name = appHttpHealthName

		// successCondition
		successCondition := new(v1beta1.NetSuccessCondition)
		successCondition.SuccessRate = &successRate
		successCondition.MeanAccessDelayInMs = &successMean
		appHttpHealth.Spec.SuccessCondition = successCondition

		// target
		target := new(v1beta1.AppHttpHealthyTarget)
		target.Method = "POST"
		if net.ParseIP(testSvcIP).To4() == nil {
			target.Host = fmt.Sprintf("http://[%s]:%d", testSvcIP, httpPort)
		} else {
			target.Host = fmt.Sprintf("http://%s:%d", testSvcIP, httpPort)
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
		Expect(success).To(BeTrue(), "compare report and task result")
		Expect(e).NotTo(HaveOccurred(), "compare report and task")
	})

	It("Successfully http testing appHttpHealth method HEAD", Label("A00008"), func() {
		var e error
		successRate := float64(1)
		successMean := int64(1500)
		crontab := "0 1"
		appHttpHealthName := "apphttphealth-head" + tools.RandomName()

		appHttpHealth := new(v1beta1.AppHttpHealthy)
		appHttpHealth.Name = appHttpHealthName

		// successCondition
		successCondition := new(v1beta1.NetSuccessCondition)
		successCondition.SuccessRate = &successRate
		successCondition.MeanAccessDelayInMs = &successMean
		appHttpHealth.Spec.SuccessCondition = successCondition

		// target
		target := new(v1beta1.AppHttpHealthyTarget)
		target.Method = "HEAD"
		if net.ParseIP(testSvcIP).To4() == nil {
			target.Host = fmt.Sprintf("http://[%s]:%d", testSvcIP, httpPort)
		} else {
			target.Host = fmt.Sprintf("http://%s:%d", testSvcIP, httpPort)
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
		Expect(success).To(BeTrue(), "compare report and task result")
		Expect(e).NotTo(HaveOccurred(), "compare report and task")

	})

	It("Successfully http testing appHttpHealth method PATCH", Label("A00009"), func() {
		var e error
		successRate := float64(1)
		successMean := int64(1500)
		crontab := "0 1"
		appHttpHealthName := "apphttphealth-patch" + tools.RandomName()

		appHttpHealth := new(v1beta1.AppHttpHealthy)
		appHttpHealth.Name = appHttpHealthName

		// successCondition
		successCondition := new(v1beta1.NetSuccessCondition)
		successCondition.SuccessRate = &successRate
		successCondition.MeanAccessDelayInMs = &successMean
		appHttpHealth.Spec.SuccessCondition = successCondition

		// target
		target := new(v1beta1.AppHttpHealthyTarget)
		target.Method = "PATCH"

		if net.ParseIP(testSvcIP).To4() == nil {
			target.Host = fmt.Sprintf("http://[%s]:%d", testSvcIP, httpPort)
		} else {
			target.Host = fmt.Sprintf("http://%s:%d", testSvcIP, httpPort)
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
		Expect(success).To(BeTrue(), "compare report and task result")
		Expect(e).NotTo(HaveOccurred(), "compare report and task")

	})

	It("Successfully http testing appHttpHealth method OPTIONS", Label("A00010"), func() {
		var e error
		successRate := float64(1)
		successMean := int64(1500)
		crontab := "0 1"
		appHttpHealthName := "apphttphealth-options" + tools.RandomName()

		appHttpHealth := new(v1beta1.AppHttpHealthy)
		appHttpHealth.Name = appHttpHealthName

		// successCondition
		successCondition := new(v1beta1.NetSuccessCondition)
		successCondition.SuccessRate = &successRate
		successCondition.MeanAccessDelayInMs = &successMean
		appHttpHealth.Spec.SuccessCondition = successCondition

		// target
		target := new(v1beta1.AppHttpHealthyTarget)
		target.Method = "OPTIONS"

		if net.ParseIP(testSvcIP).To4() == nil {
			target.Host = fmt.Sprintf("http://[%s]:%d", testSvcIP, httpPort)
		} else {
			target.Host = fmt.Sprintf("http://%s:%d", testSvcIP, httpPort)
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
		Expect(success).To(BeTrue(), "compare report and task result")
		Expect(e).NotTo(HaveOccurred(), "compare report and task")

	})

	It("Successfully http testing appHttpHealth due to success rate", Label("A000012"), func() {
		var e error
		successRate := float64(0.2)
		successMean := int64(1200)
		crontab := "0 1"
		appHttpHealthName := "apphttphealth-get" + tools.RandomName()

		appHttpHealth := new(v1beta1.AppHttpHealthy)
		appHttpHealth.Name = appHttpHealthName

		// successCondition
		successCondition := new(v1beta1.NetSuccessCondition)
		successCondition.SuccessRate = &successRate
		successCondition.MeanAccessDelayInMs = &successMean
		appHttpHealth.Spec.SuccessCondition = successCondition

		// target
		target := new(v1beta1.AppHttpHealthyTarget)
		target.Method = "GET"
		if net.ParseIP(testSvcIP).To4() == nil {
			target.Host = fmt.Sprintf("http://[%s]:%d/?delay=1", testSvcIP, httpPort)
		} else {
			target.Host = fmt.Sprintf("http://%s:%d/?delay=1", testSvcIP, httpPort)
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
		Expect(success).To(BeTrue(), "compare report and task result")
		Expect(e).NotTo(HaveOccurred(), "compare report and task")

	})

	It("Successfully https testing appHttpHealth method GET Protocol Http2", Label("A000013"), func() {
		var e error
		successRate := float64(1)
		successMean := int64(1500)
		crontab := "0 1"
		appHttpHealthName := "apphttphealth-get" + tools.RandomName()

		appHttpHealth := new(v1beta1.AppHttpHealthy)
		appHttpHealth.Name = appHttpHealthName

		// successCondition
		successCondition := new(v1beta1.NetSuccessCondition)
		successCondition.SuccessRate = &successRate
		successCondition.MeanAccessDelayInMs = &successMean
		appHttpHealth.Spec.SuccessCondition = successCondition

		// target
		target := new(v1beta1.AppHttpHealthyTarget)
		target.Method = "GET"

		if net.ParseIP(testPodIPs[0]).To4() == nil {
			target.Host = fmt.Sprintf("https://[%s]", testPodIPs[0])
		} else {
			target.Host = fmt.Sprintf("https://%s", testPodIPs[0])
		}
		target.TlsSecretName = &common.TlsClientName
		target.TlsSecretNamespace = &common.TestNameSpace
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
		Expect(success).To(BeTrue(), "compare report and task result")
		Expect(e).NotTo(HaveOccurred(), "compare report and task")

	})
})
