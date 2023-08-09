// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package apphttphealth_test

import (
	"context"
	"fmt"
	kdoctor_v1beta1 "github.com/kdoctor-io/kdoctor/pkg/k8s/apis/kdoctor.io/v1beta1"
	"github.com/kdoctor-io/kdoctor/test/e2e/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	e2e "github.com/spidernet-io/e2eframework/framework"
	"github.com/spidernet-io/e2eframework/tools"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"testing"
	"time"
)

func TestNetReach(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AppHttpHealth Suite")
}

var (
	frame             *e2e.Framework
	httpPort          = 80
	httpsPort         = 443
	bodyConfigMapName string
	caSecret          *v1.Secret
	reportNum         int
	testAppName       string
	testAppNamespace  string
	testSvcIP         string
	testPodIPs        []string
)

var _ = BeforeSuite(func() {
	defer GinkgoRecover()
	var e error
	frame, e = e2e.NewFramework(GinkgoT(), []func(*runtime.Scheme) error{kdoctor_v1beta1.AddToScheme})
	Expect(e).NotTo(HaveOccurred())

	// test request body
	bodyConfigMapName = "kdoctor-test-body-" + tools.RandomName()
	cm := new(v1.ConfigMap)
	cm.SetName(bodyConfigMapName)
	cm.SetNamespace(common.TestNameSpace)
	body := make(map[string]string, 0)
	body["test1"] = "test1"
	body["test2"] = "test2"
	cm.Data = body
	e = frame.CreateConfigmap(cm)
	Expect(e).NotTo(HaveOccurred(), "create test body configmap")

	caSecret = new(v1.Secret)
	key := types.NamespacedName{
		Name:      common.KDoctorCaName,
		Namespace: common.TestNameSpace,
	}
	e = frame.KClient.Get(context.Background(), key, caSecret)
	Expect(e).NotTo(HaveOccurred(), "get kdoctor ca secret")

	nodeLIst, e := frame.GetNodeList()
	Expect(e).NotTo(HaveOccurred(), "get node list")
	reportNum = len(nodeLIst.Items)

	testAppName = "app-" + tools.RandomName()
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
})

var _ = AfterSuite(func() {
	defer GinkgoRecover()
	Expect(frame.DeleteConfigmap(bodyConfigMapName, common.TestNameSpace)).NotTo(HaveOccurred())
	Expect(frame.DeleteNamespace(testAppNamespace)).NotTo(HaveOccurred())
})
