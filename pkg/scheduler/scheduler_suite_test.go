// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package scheduler

import (
	"github.com/kdoctor-io/kdoctor/pkg/k8s/apis/kdoctor.io/v1beta1"
	"github.com/kdoctor-io/kdoctor/pkg/types"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestScheduler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "scheduler Suite")
}

var c client.WithWatch

var _ = BeforeSuite(func() {
	scheme := runtime.NewScheme()
	err := networkingv1.AddToScheme(scheme)
	Expect(err).To(BeNil(), "add ingress scheme failed")
	err = clientgoscheme.AddToScheme(scheme)
	Expect(err).To(BeNil(), "add client go scheme failed")
	err = v1beta1.AddToScheme(scheme)
	Expect(err).To(BeNil(), "add kdoctor scheme failed")
	builder := fake.NewClientBuilder().WithScheme(scheme)
	c = builder.Build()

	// workload template
	// ingress
	ingress := new(networkingv1.Ingress)
	rule := networkingv1.IngressRule{}
	path := networkingv1.HTTPIngressPath{}
	backend := networkingv1.IngressBackend{}
	backend.Service = new(networkingv1.IngressServiceBackend)
	path.Backend = backend
	rule.HTTP = new(networkingv1.HTTPIngressRuleValue)
	rule.HTTP.Paths = append(rule.HTTP.Paths, path)
	ingress.Spec.Rules = append([]networkingv1.IngressRule{}, rule)
	types.IngressTempl = ingress

	// service
	service := new(corev1.Service)
	service.SetLabels(map[string]string{"test": "test"})
	service.Spec.Selector = map[string]string{"test": "test"}
	types.ServiceTempl = service

	// pod
	pod := new(corev1.Pod)
	container := corev1.Container{Name: "test"}
	pod.Spec.Containers = []corev1.Container{container}
	types.PodTempl = pod

	// deployment
	deployment := new(appsv1.Deployment)
	deployment.SetLabels(map[string]string{"test": "test"})
	deployment.Spec.Selector = &v1.LabelSelector{}
	replicas := int32(1)
	deployment.Spec.Replicas = &replicas
	deployment.Spec.Template.Spec.HostNetwork = true
	types.DeploymentTempl = deployment

	// daemonSet
	daemonSet := new(appsv1.DaemonSet)
	daemonSet.SetLabels(map[string]string{"test": "test"})
	daemonSet.Spec.Selector = &v1.LabelSelector{}
	types.DaemonsetTempl = daemonSet
})
