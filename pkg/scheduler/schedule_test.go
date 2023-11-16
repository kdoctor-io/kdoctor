// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0
package scheduler

import (
	"context"
	"github.com/kdoctor-io/kdoctor/pkg/k8s/apis/kdoctor.io/v1beta1"
	"github.com/kdoctor-io/kdoctor/pkg/logger"
	"github.com/kdoctor-io/kdoctor/pkg/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
)

var _ = Describe("schedule unit test", Label("schedule"), func() {

	It("schedule appHttpHealthy", Label("schedule"), func() {
		types.ControllerConfig.Configmap.EnableIPv4 = true
		types.ControllerConfig.Configmap.EnableIPv6 = true
		tgm := int64(1)
		ctx := context.Background()
		taskName := "test"
		task := &v1beta1.AppHttpHealthy{}
		replicas := int32(2)
		agentSpec := v1beta1.AgentSpec{
			Kind:                          types.KindDeployment,
			DeploymentReplicas:            &replicas,
			Affinity:                      &v1.Affinity{},
			Env:                           []v1.EnvVar{{Name: "test", Value: "test"}},
			HostNetwork:                   true,
			Resources:                     &v1.ResourceRequirements{},
			Annotation:                    map[string]string{"test": "test"},
			TerminationGracePeriodMinutes: &tgm,
		}
		schedule := NewScheduler(c, c, types.KindNameAppHttpHealthy, taskName, "", logger.NewStdoutLogger("debug", "schedule"))
		_, err := schedule.CreateTaskRuntimeIfNotExist(ctx, task, agentSpec)
		Expect(err).To(BeNil(), "create not exist runtime")
	})

	It("schedule netReach", Label("schedule"), func() {
		types.ControllerConfig.Configmap.EnableIPv4 = true
		types.ControllerConfig.Configmap.EnableIPv6 = true
		tgm := int64(1)
		ctx := context.Background()
		taskName := "test"
		task := &v1beta1.NetReach{}
		enabelIngress := true

		target := new(v1beta1.NetReachTarget)
		target.Ingress = &enabelIngress
		task.Spec.Target = target

		schedule := NewScheduler(c, c, types.KindNameNetReach, taskName, "", logger.NewStdoutLogger("debug", "schedule"))
		agentSpec := v1beta1.AgentSpec{
			Kind:                          types.KindDaemonSet,
			Affinity:                      &v1.Affinity{},
			Env:                           []v1.EnvVar{{Name: "test", Value: "test"}},
			HostNetwork:                   false,
			Resources:                     &v1.ResourceRequirements{},
			Annotation:                    map[string]string{"test": "test"},
			TerminationGracePeriodMinutes: &tgm,
		}
		_, err := schedule.CreateTaskRuntimeIfNotExist(ctx, task, agentSpec)
		Expect(err).To(BeNil(), "create not exists runtime")
	})

	It("schedule unrecognized type", Label("schedule"), func() {
		types.ControllerConfig.Configmap.EnableIPv4 = true
		types.ControllerConfig.Configmap.EnableIPv6 = true
		tgm := int64(1)
		ctx := context.Background()
		taskName := "test"
		task := &v1beta1.NetReach{}
		enabelIngress := true

		target := new(v1beta1.NetReachTarget)
		target.Ingress = &enabelIngress
		task.Spec.Target = target

		schedule := NewScheduler(c, c, types.KindNameNetReach, taskName, "", logger.NewStdoutLogger("debug", "schedule"))
		agentSpec := v1beta1.AgentSpec{
			Kind:                          "test",
			Affinity:                      &v1.Affinity{},
			Env:                           []v1.EnvVar{{Name: "test", Value: "test"}},
			HostNetwork:                   false,
			Resources:                     &v1.ResourceRequirements{},
			Annotation:                    map[string]string{"test": "test"},
			TerminationGracePeriodMinutes: &tgm,
		}
		_, err := schedule.CreateTaskRuntimeIfNotExist(ctx, task, agentSpec)
		Expect(err).NotTo(BeNil(), "create not exists runtime,unrecognized type")
	})
})
