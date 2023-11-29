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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

var _ = Describe("schedule tracing unit test", Label("schedule"), func() {

	It("schedule tracing appHttpHealthy", Label("schedule"), func() {
		config := TrackerConfig{
			ItemChannelBuffer:     500,
			MaxDatabaseCap:        5000,
			ExecutorWorkers:       3,
			SignalTimeOutDuration: time.Second * 3,
			TraceGapDuration:      time.Second * 10,
		}
		types.ControllerConfig.Configmap.EnableIPv4 = true
		types.ControllerConfig.Configmap.EnableIPv6 = true
		types.ControllerConfig.PodNamespace = "default"
		track := NewTracker(c, c, config, logger.NewStdoutLogger("debug", "tracker"))
		tgm := int64(1)
		ctx := context.Background()
		taskName := "test-apphttp"
		task := &v1beta1.AppHttpHealthy{}
		task.SetName(taskName)
		replicas := int32(1)
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
		task.Spec.AgentSpec = &agentSpec
		err := c.Create(ctx, task)
		Expect(err).To(BeNil(), "create task")

		schedule := NewScheduler(c, c, types.KindNameAppHttpHealthy, taskName, "", logger.NewStdoutLogger("debug", "schedule"))
		resource, err := schedule.CreateTaskRuntimeIfNotExist(ctx, task, agentSpec)
		Expect(err).To(BeNil(), "create not exist runtime")
		deleteTime := metav1.NewTime(metav1.Now().Add(time.Second * 1))
		item := BuildItem(resource, types.KindNameAppHttpHealthy, taskName, deleteTime.DeepCopy())
		err = track.DB.Apply(item)
		Expect(err).To(BeNil(), "track db apply")

		task.Status.Resource = &resource
		err = c.Update(ctx, task)
		Expect(err).To(BeNil(), "update resource")
		ctx, cancel := context.WithTimeout(ctx, time.Second*10)
		defer cancel()
		time.Sleep(time.Second * 2)
		track.Start(ctx)
		time.Sleep(time.Second * 15)
	})

	It("schedule tracing netReach", Label("schedule"), func() {
		config := TrackerConfig{
			ItemChannelBuffer:     500,
			MaxDatabaseCap:        5000,
			ExecutorWorkers:       3,
			SignalTimeOutDuration: time.Second * 3,
			TraceGapDuration:      time.Second * 10,
		}
		types.ControllerConfig.Configmap.EnableIPv4 = true
		types.ControllerConfig.Configmap.EnableIPv6 = true
		types.ControllerConfig.PodNamespace = "default"
		track := NewTracker(c, c, config, logger.NewStdoutLogger("debug", "tracker"))
		tgm := int64(1)
		ctx := context.Background()
		taskName := "test-netreach"
		task := &v1beta1.NetReach{}
		enable := true
		task.Spec.Target = &v1beta1.NetReachTarget{
			Ingress:         &enable,
			IPv4:            &enable,
			IPv6:            &enable,
			ClusterIP:       &enable,
			NodePort:        &enable,
			MultusInterface: &enable,
			LoadBalancer:    &enable,
		}
		task.SetName(taskName)
		replicas := int32(1)
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
		task.Spec.AgentSpec = &agentSpec
		err := c.Create(ctx, task)
		Expect(err).To(BeNil(), "create task")

		schedule := NewScheduler(c, c, types.KindNameNetReach, taskName, "", logger.NewStdoutLogger("debug", "schedule"))
		resource, err := schedule.CreateTaskRuntimeIfNotExist(ctx, task, agentSpec)
		Expect(err).To(BeNil(), "create not exist runtime")
		deleteTime := metav1.NewTime(metav1.Now().Add(time.Second * 1))
		item := BuildItem(resource, types.KindNameNetReach, taskName, deleteTime.DeepCopy())
		err = track.DB.Apply(item)
		Expect(err).To(BeNil(), "track db apply")

		task.Status.Resource = &resource
		err = c.Update(ctx, task)
		Expect(err).To(BeNil(), "update resource")
		ctx, cancel := context.WithTimeout(ctx, time.Second*10)
		defer cancel()
		time.Sleep(time.Second * 2)
		track.Start(ctx)
		time.Sleep(time.Second * 15)
	})

	It("schedule tracing netDNS", Label("schedule"), func() {
		config := TrackerConfig{
			ItemChannelBuffer:     500,
			MaxDatabaseCap:        5000,
			ExecutorWorkers:       3,
			SignalTimeOutDuration: time.Second * 3,
			TraceGapDuration:      time.Second * 10,
		}
		types.ControllerConfig.Configmap.EnableIPv4 = true
		types.ControllerConfig.Configmap.EnableIPv6 = true
		types.ControllerConfig.PodNamespace = "default"
		track := NewTracker(c, c, config, logger.NewStdoutLogger("debug", "tracker"))
		tgm := int64(1)
		ctx := context.Background()
		taskName := "test-dns"
		task := &v1beta1.Netdns{}
		task.SetName(taskName)
		replicas := int32(1)
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
		task.Spec.AgentSpec = &agentSpec
		err := c.Create(ctx, task)
		Expect(err).To(BeNil(), "create task")

		schedule := NewScheduler(c, c, types.KindNameNetdns, taskName, "", logger.NewStdoutLogger("debug", "schedule"))
		resource, err := schedule.CreateTaskRuntimeIfNotExist(ctx, task, agentSpec)
		Expect(err).To(BeNil(), "create not exist runtime")
		deleteTime := metav1.NewTime(metav1.Now().Add(time.Second * 1))
		item := BuildItem(resource, types.KindNameNetdns, taskName, deleteTime.DeepCopy())
		err = track.DB.Apply(item)
		Expect(err).To(BeNil(), "track db apply")

		task.Status.Resource = &resource
		err = c.Update(ctx, task)
		Expect(err).To(BeNil(), "update resource")
		ctx, cancel := context.WithTimeout(ctx, time.Second*10)
		defer cancel()
		time.Sleep(time.Second * 2)
		track.Start(ctx)
		time.Sleep(time.Second * 15)
	})

	It("schedule tracing appHttpHealthy update resource status", Label("schedule"), func() {
		config := TrackerConfig{
			ItemChannelBuffer:     500,
			MaxDatabaseCap:        5000,
			ExecutorWorkers:       3,
			SignalTimeOutDuration: time.Second * 3,
			TraceGapDuration:      time.Second * 10,
		}
		types.ControllerConfig.Configmap.EnableIPv4 = true
		types.ControllerConfig.Configmap.EnableIPv6 = true
		types.ControllerConfig.PodNamespace = "default"
		track := NewTracker(c, c, config, logger.NewStdoutLogger("debug", "tracker"))
		tgm := int64(1)
		ctx := context.Background()
		taskName := "runtime-apphttp"
		task := &v1beta1.AppHttpHealthy{}
		task.SetName(taskName)
		replicas := int32(1)
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
		task.Spec.AgentSpec = &agentSpec
		err := c.Create(ctx, task)
		Expect(err).To(BeNil(), "create task")

		schedule := NewScheduler(c, c, types.KindNameAppHttpHealthy, taskName, "", logger.NewStdoutLogger("debug", "schedule"))
		resource, err := schedule.CreateTaskRuntimeIfNotExist(ctx, task, agentSpec)
		Expect(err).To(BeNil(), "create not exist runtime")
		deleteTime := metav1.NewTime(metav1.Now().Add(time.Second * 1))
		item := BuildItem(resource, types.KindNameAppHttpHealthy, taskName, deleteTime.DeepCopy())
		err = track.DB.Apply(item)
		Expect(err).To(BeNil(), "track db apply")

		task.Status.Resource = &resource
		task.Status.Resource.RuntimeStatus = v1beta1.RuntimeCreating
		err = c.Update(ctx, task)
		Expect(err).To(BeNil(), "update resource")
		ctx, cancel := context.WithTimeout(ctx, time.Second*10)
		defer cancel()
		time.Sleep(time.Second * 2)
		track.Start(ctx)
		time.Sleep(time.Second * 15)
	})
})
