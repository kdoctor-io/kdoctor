// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package scheduler

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kdoctor-io/kdoctor/pkg/k8s/apis/kdoctor.io/v1beta1"
	"github.com/kdoctor-io/kdoctor/pkg/types"
)

const (
	kdoctor = "kdoctor"

	containerArgTaskKind  = "--task-kind"
	containerArgTaskName  = "--task-name"
	containerArgServiceV4 = "--service-ipv4-name"
	containerArgServiceV6 = "--service-ipv6-name"
	UniqueMatchLabelKey   = "app.kubernetes.io/name"
)

type Scheduler struct {
	client    client.Client
	apiReader client.Reader

	taskKind string
	taskName string

	uniqueLabelKey string
	// this is a unique string variable that combined with taskKind and taskName, we can use it as a selector
	uniqueLabelValue string

	log *zap.Logger
}

func NewScheduler(client client.Client, apiReader client.Reader, taskKind, taskName, uniqueLabelKey string, log *zap.Logger) *Scheduler {
	s := &Scheduler{
		client:           client,
		apiReader:        apiReader,
		taskKind:         taskKind,
		taskName:         taskName,
		uniqueLabelKey:   uniqueLabelKey,
		uniqueLabelValue: UniqueMatchLabelValue(taskKind, taskName),
		log:              log,
	}

	return s
}

func (s *Scheduler) CreateTaskRuntimeIfNotExist(ctx context.Context, ownerTask metav1.Object, agentSpec v1beta1.AgentSpec) (v1beta1.TaskResource, error) {
	taskRuntimeName := TaskRuntimeName(s.taskKind, s.taskName)
	resource := v1beta1.TaskResource{
		RuntimeName:   taskRuntimeName,
		RuntimeType:   "",
		ServiceNameV4: nil,
		ServiceNameV6: nil,
		RuntimeStatus: v1beta1.RuntimeCreating,
	}

	var runtime client.Object

	switch agentSpec.Kind {
	case types.KindDeployment:
		runtime = &appsv1.Deployment{}
		resource.RuntimeType = types.KindDeployment

	case types.KindDaemonSet:
		runtime = &appsv1.DaemonSet{}
		resource.RuntimeType = types.KindDaemonSet

	default:
		return v1beta1.TaskResource{}, fmt.Errorf("unrecognized agent runtime kind '%s'", agentSpec.Kind)
	}

	var needCreate bool
	objectKey := client.ObjectKey{
		// reuse kdoctor-controller namespace
		Namespace: types.ControllerConfig.PodNamespace,
		Name:      taskRuntimeName,
	}

	s.log.Sugar().Debugf("try to get task '%s/%s' corresponding runtime '%s'", s.taskKind, s.taskName, taskRuntimeName)
	err := s.apiReader.Get(ctx, objectKey, runtime)
	if nil != err {
		if errors.IsNotFound(err) {
			s.log.Sugar().Infof("task '%s/%s' corresponding runtime '%s' not found, try to create one", s.taskKind, s.taskName, taskRuntimeName)
			needCreate = true
		} else {
			return v1beta1.TaskResource{}, err
		}
	}

	if needCreate {
		// create task Runtime
		if agentSpec.Kind == types.KindDeployment {
			runtime = s.generateDeployment(agentSpec)
		} else if agentSpec.Kind == types.KindDaemonSet {
			runtime = s.generateDaemonSet(agentSpec)
		}

		runtime.SetName(taskRuntimeName)
		runtime.SetNamespace(types.ControllerConfig.PodNamespace)

		err := controllerruntime.SetControllerReference(ownerTask, runtime, s.client.Scheme())
		if nil != err {
			return v1beta1.TaskResource{}, fmt.Errorf("failed to set task %s/%s corresponding runtime %s controllerReference, error: %v", s.taskKind, s.taskName, taskRuntimeName, err)
		}

		s.log.Sugar().Infof("try to create task %s/%s corresponding runtime %v", s.taskKind, s.taskName, runtime)
		err = s.client.Create(ctx, runtime)
		if nil != err {
			return v1beta1.TaskResource{}, err
		}
	}

	// Service
	if types.ControllerConfig.Configmap.EnableIPv4 {
		serviceNameV4, err := s.createService(ctx, taskRuntimeName, agentSpec, runtime, corev1.IPv4Protocol)
		if nil != err {
			return v1beta1.TaskResource{}, fmt.Errorf("failed to create runtime IPv4 service for task '%s/%s', error: %w", s.taskKind, s.taskName, err)
		}
		resource.ServiceNameV4 = pointer.String(serviceNameV4)
	}
	if types.ControllerConfig.Configmap.EnableIPv6 {
		serviceNameV6, err := s.createService(ctx, taskRuntimeName, agentSpec, runtime, corev1.IPv6Protocol)
		if nil != err {
			return v1beta1.TaskResource{}, fmt.Errorf("failed to create runtime IPv6 service for task '%s/%s', error: %w", s.taskKind, s.taskName, err)
		}
		resource.ServiceNameV6 = pointer.String(serviceNameV6)
	}

	return resource, nil
}

func (s *Scheduler) createService(ctx context.Context, taskRuntimeName string, agentSpec v1beta1.AgentSpec, ownerRuntime client.Object, ipFamily corev1.IPFamily) (serviceName string, err error) {
	needCreate := false

	// service name
	svcName := TaskRuntimeServiceName(taskRuntimeName, ipFamily)

	var service corev1.Service
	objectKey := client.ObjectKey{
		// reuse kdoctor-controller namespace
		Namespace: types.ControllerConfig.PodNamespace,
		Name:      svcName,
	}
	s.log.Sugar().Debugf("try to get task '%s/%s' corresponding runtime service '%s'", s.taskKind, s.taskName, svcName)
	err = s.apiReader.Get(ctx, objectKey, &service)
	if nil != err {
		if errors.IsNotFound(err) {
			s.log.Sugar().Debugf("task '%s/%s' corresponding runtime service '%s' not found, try to create one", s.taskKind, s.taskName, svcName)
			needCreate = true
		} else {
			return "", err
		}
	}

	if needCreate {
		svc := s.generateService(agentSpec, ipFamily)
		svc.SetName(svcName)
		svc.SetNamespace(types.ControllerConfig.PodNamespace)

		err := controllerruntime.SetControllerReference(ownerRuntime, svc, s.client.Scheme())
		if nil != err {
			return "", fmt.Errorf("failed to set service %s/%s controllerReference with runtime %s, error: %v",
				types.ControllerConfig.PodNamespace, taskRuntimeName, ownerRuntime.GetName(), err)
		}

		s.log.Sugar().Infof("try to create task  %s/%s corresponding runtime service '%v'", s.taskKind, s.taskName, service)
		err = s.client.Create(ctx, svc)
		if nil != err {
			return "", err
		}
	}

	return svcName, nil
}

func (s *Scheduler) generateDaemonSet(agentSpec v1beta1.AgentSpec) *appsv1.DaemonSet {
	daemonSet := types.DaemonsetTempl.DeepCopy()

	// add unique selector
	var selector metav1.LabelSelector
	if daemonSet.Spec.Selector != nil {
		daemonSet.Spec.Selector.DeepCopyInto(&selector)
	}
	selector.MatchLabels = AppendAnnotationOrLabel(selector.MatchLabels, map[string]string{
		s.uniqueLabelKey: s.uniqueLabelValue,
	})
	daemonSet.Spec.Selector = &selector

	// replace AgentSpec properties
	if len(agentSpec.Annotation) != 0 {
		daemonSet.SetAnnotations(AppendAnnotationOrLabel(daemonSet.Annotations, agentSpec.Annotation))
	}

	// assemble
	podTemplateSpec := s.generatePodTemplateSpec(s.uniqueLabelValue, agentSpec)
	daemonSet.Spec.Template = podTemplateSpec

	return daemonSet
}

func (s *Scheduler) generateDeployment(agentSpec v1beta1.AgentSpec) *appsv1.Deployment {
	deployment := types.DeploymentTempl.DeepCopy()

	// add unique selector
	var selector metav1.LabelSelector
	if deployment.Spec.Selector != nil {
		deployment.Spec.Selector.DeepCopyInto(&selector)
	}
	selector.MatchLabels = AppendAnnotationOrLabel(selector.MatchLabels, map[string]string{
		s.uniqueLabelKey: s.uniqueLabelValue,
	})
	deployment.Spec.Selector = &selector

	// replace AgentSpec properties
	if len(agentSpec.Annotation) != 0 {
		deployment.SetAnnotations(AppendAnnotationOrLabel(deployment.Annotations, agentSpec.Annotation))
	}
	if agentSpec.DeploymentReplicas != nil {
		deployment.Spec.Replicas = agentSpec.DeploymentReplicas
	}

	// assemble
	podTemplateSpec := s.generatePodTemplateSpec(s.uniqueLabelValue, agentSpec)
	deployment.Spec.Template = podTemplateSpec

	return deployment
}

func (s *Scheduler) generatePodTemplateSpec(uniqueLabelVal string, agentSpec v1beta1.AgentSpec) corev1.PodTemplateSpec {
	pod := types.PodTempl.DeepCopy()

	// 1. add container start parameters
	for index := range pod.Spec.Containers {
		tmpArgs := pod.Spec.Containers[index].Args
		tmpArgs = append(tmpArgs,
			fmt.Sprintf("%s=%s", containerArgTaskKind, s.taskKind),
			fmt.Sprintf("%s=%s", containerArgTaskName, s.taskName),
		)
		if types.ControllerConfig.Configmap.EnableIPv4 {
			tmpArgs = append(tmpArgs, fmt.Sprintf("%s=%s", containerArgServiceV4, TaskRuntimeServiceName(TaskRuntimeName(s.taskKind, s.taskName), corev1.IPv4Protocol)))
		}
		if types.ControllerConfig.Configmap.EnableIPv6 {
			tmpArgs = append(tmpArgs, fmt.Sprintf("%s=%s", containerArgServiceV6, TaskRuntimeServiceName(TaskRuntimeName(s.taskKind, s.taskName), corev1.IPv6Protocol)))
		}

		pod.Spec.Containers[index].Args = tmpArgs
	}

	// 2. add unique selector
	{
		podLabels := pod.GetLabels()
		if podLabels == nil {
			podLabels = make(map[string]string)
		}
		podLabels[UniqueMatchLabelKey] = uniqueLabelVal
		pod.SetLabels(podLabels)
	}

	// 3. replace AgentSpec properties
	{
		if len(agentSpec.Annotation) != 0 {
			pod.SetAnnotations(AppendAnnotationOrLabel(pod.Annotations, agentSpec.Annotation))
		}

		if agentSpec.Affinity != nil {
			pod.Spec.Affinity = agentSpec.Affinity
		}

		if pod.Spec.HostNetwork != agentSpec.HostNetwork {
			if agentSpec.HostNetwork {
				pod.Spec.HostNetwork = true
				pod.Spec.DNSPolicy = corev1.DNSClusterFirstWithHostNet
			} else {
				pod.Spec.HostNetwork = false
				pod.Spec.DNSPolicy = corev1.DNSClusterFirst
			}
		}

		for index := range pod.Spec.Containers {
			if len(agentSpec.Env) != 0 {
				pod.Spec.Containers[index].Env = append(pod.Spec.Containers[index].Env, agentSpec.Env...)
			}

			if agentSpec.Resources != nil {
				pod.Spec.Containers[index].Resources = *agentSpec.Resources
			}
		}
	}

	podTemplateSpec := corev1.PodTemplateSpec{
		ObjectMeta: pod.ObjectMeta,
		Spec:       pod.Spec,
	}

	return podTemplateSpec
}

func (s *Scheduler) generateService(agentSpec v1beta1.AgentSpec, ipFamily corev1.IPFamily) *corev1.Service {
	service := types.ServiceTempl.DeepCopy()

	// add unique selector
	selector := map[string]string{}
	if service.Spec.Selector != nil {
		selector = service.Spec.Selector
	}
	selector[s.uniqueLabelKey] = s.uniqueLabelValue
	service.Spec.Selector = selector

	// replace AgentSpec properties
	if len(agentSpec.Annotation) != 0 {
		service.SetAnnotations(AppendAnnotationOrLabel(service.Annotations, agentSpec.Annotation))
	}

	// set IP Family
	service.Spec.IPFamilies = []corev1.IPFamily{ipFamily}

	return service
}
