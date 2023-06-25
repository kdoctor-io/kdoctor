// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package apphttphealth_test

import (
	"context"
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
)

func TestNetReach(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AppHttpHealth Suite")
}

var frame *e2e.Framework

var httpPort = 80
var httpsPort = 443
var bodyConfigMapName string
var caSecret *v1.Secret
var reportNum int

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

	ds, e := frame.GetDaemonSet(common.KDoctorAgentDSName, common.TestNameSpace)
	Expect(e).NotTo(HaveOccurred(), "get kdoctor-agent daemonset")
	reportNum = int(ds.Status.NumberReady)
})

var _ = AfterSuite(func() {
	defer GinkgoRecover()
	_ = frame.DeleteConfigmap(bodyConfigMapName, common.TestNameSpace)
})
