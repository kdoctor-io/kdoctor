// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package netdns_test

import (
	"context"
	kdoctor_v1beta1 "github.com/kdoctor-io/kdoctor/pkg/k8s/apis/kdoctor.io/v1beta1"
	"github.com/kdoctor-io/kdoctor/test/e2e/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	e2e "github.com/spidernet-io/e2eframework/framework"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"testing"
	"time"
)

func TestNetReach(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "NetDns Suite")
}

var (
	frame            *e2e.Framework
	KubeDnsName      string
	KubeDnsNamespace string
	reportNum        int
	testSvcIP        string
	testPodIPs       []string
)

var _ = BeforeSuite(func() {
	defer GinkgoRecover()
	var e error
	frame, e = e2e.NewFramework(GinkgoT(), []func(*runtime.Scheme) error{kdoctor_v1beta1.AddToScheme})
	Expect(e).NotTo(HaveOccurred())
	nodeLIst, e := frame.GetNodeList()
	Expect(e).NotTo(HaveOccurred(), "get node list")
	reportNum = len(nodeLIst.Items)

	KubeServiceList := &v1.ServiceList{}
	ops := []client.ListOption{
		client.MatchingLabels(map[string]string{"k8s-app": "kube-dns"}),
	}
	e = frame.KClient.List(context.Background(), KubeServiceList, ops...)
	Expect(e).NotTo(HaveOccurred(), "get kube dns service")
	KubeDnsName = KubeServiceList.Items[0].Name
	KubeDnsNamespace = KubeServiceList.Items[0].Namespace

	//  get test app service ip and pod ip
	svc, e := frame.GetService(common.TestAppName, common.TestNameSpace)
	Expect(e).NotTo(HaveOccurred(), "get test app service")
	testSvcIP = svc.Spec.ClusterIP
	GinkgoWriter.Printf("get test service ip %v \n", testSvcIP)

	podLIst, e := frame.WaitDeploymentReadyAndCheckIP(common.TestAppName, common.TestNameSpace, time.Second*60)
	Expect(e).NotTo(HaveOccurred(), "wait test app deploy ready")

	testPodIPs = make([]string, 0)
	for _, v := range podLIst.Items {
		testPodIPs = append(testPodIPs, v.Status.PodIP)
	}

	GinkgoWriter.Printf("get test pod ips %v \n", testPodIPs)
})

var _ = AfterSuite(func() {
	defer GinkgoRecover()
})
