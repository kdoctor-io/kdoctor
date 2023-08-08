// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package netreach_test

import (
	kdoctor_v1beta1 "github.com/kdoctor-io/kdoctor/pkg/k8s/apis/kdoctor.io/v1beta1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	e2e "github.com/spidernet-io/e2eframework/framework"
	"k8s.io/apimachinery/pkg/runtime"
	"testing"
)

func TestNetReach(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "NetReach Suite")
}

var frame *e2e.Framework
var reportNum int
var _ = BeforeSuite(func() {
	defer GinkgoRecover()
	var e error
	frame, e = e2e.NewFramework(GinkgoT(), []func(*runtime.Scheme) error{kdoctor_v1beta1.AddToScheme})
	Expect(e).NotTo(HaveOccurred())
	nodeLIst, e := frame.GetNodeList()
	Expect(e).NotTo(HaveOccurred(), "get node list")
	reportNum = len(nodeLIst.Items)
	// TODO (ii2day): add agent multus network

})
