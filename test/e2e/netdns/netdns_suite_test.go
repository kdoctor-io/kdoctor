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
)

func TestNetReach(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "NetReach Suite")
}

var frame *e2e.Framework
var reportNum int
var KubeDnsName string
var KubeDnsNamespace string

var _ = BeforeSuite(func() {
	defer GinkgoRecover()
	var e error
	frame, e = e2e.NewFramework(GinkgoT(), []func(*runtime.Scheme) error{kdoctor_v1beta1.AddToScheme})
	Expect(e).NotTo(HaveOccurred())
	ds, e := frame.GetDaemonSet(common.KDoctorAgentDSName, common.TestNameSpace)
	Expect(e).NotTo(HaveOccurred(), "get kdoctor-agent daemonset")
	reportNum = int(ds.Status.NumberReady)

	KubeServiceList := &v1.ServiceList{}
	ops := []client.ListOption{
		client.MatchingLabels(map[string]string{"k8s-app": "kube-dns"}),
	}
	e = frame.KClient.List(context.Background(), KubeServiceList, ops...)
	Expect(e).NotTo(HaveOccurred(), "get kube dns service")
	KubeDnsName = KubeServiceList.Items[0].Name
	KubeDnsNamespace = KubeServiceList.Items[0].Namespace
})
