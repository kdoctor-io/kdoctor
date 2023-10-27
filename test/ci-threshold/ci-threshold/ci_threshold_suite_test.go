// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package ci_threshold_test

import (
	kdoctor_v1beta1 "github.com/kdoctor-io/kdoctor/pkg/k8s/apis/kdoctor.io/v1beta1"
	"github.com/kdoctor-io/kdoctor/test/e2e/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	e2e "github.com/spidernet-io/e2eframework/framework"
	"k8s.io/apimachinery/pkg/runtime"
	"testing"
	"time"
)

func TestCIThreshold(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CI Threshold Suite")
}

var (
	frame     *e2e.Framework
	httpPort  = 80
	testSvcIP string
)

var _ = BeforeSuite(func() {
	defer GinkgoRecover()
	var e error
	frame, e = e2e.NewFramework(GinkgoT(), []func(*runtime.Scheme) error{kdoctor_v1beta1.AddToScheme})
	Expect(e).NotTo(HaveOccurred())

	//  get test app service ip and pod ip
	svc, e := frame.GetService(common.TestAppName, common.TestNameSpace)
	Expect(e).NotTo(HaveOccurred(), "get test app service")
	testSvcIP = svc.Spec.ClusterIP
	GinkgoWriter.Printf("get test service ip %v \n", testSvcIP)

	_, e = frame.WaitDeploymentReadyAndCheckIP(common.TestAppName, common.TestNameSpace, time.Second*60)
	Expect(e).NotTo(HaveOccurred(), "wait test app deploy ready")

})
