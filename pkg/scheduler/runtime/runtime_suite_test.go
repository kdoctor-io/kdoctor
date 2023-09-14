// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package runtime

import (
	networkingv1 "k8s.io/api/networking/v1"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestRuntimeTask(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Runtime Suite")
}

var c client.WithWatch

var _ = BeforeSuite(func() {
	scheme := runtime.NewScheme()
	err := networkingv1.AddToScheme(scheme)
	Expect(err).To(BeNil(), "add ingress scheme failed")
	err = clientgoscheme.AddToScheme(scheme)
	Expect(err).To(BeNil(), "add client go scheme failed")
	builder := fake.NewClientBuilder().WithScheme(scheme)
	c = builder.Build()

})
