// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package netreach

import (
	"context"
	"fmt"
	"reflect"

	"go.uber.org/zap"
	networkingv1 "k8s.io/api/networking/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/strings/slices"

	k8sObjManager "github.com/kdoctor-io/kdoctor/pkg/k8ObjManager"
	crd "github.com/kdoctor-io/kdoctor/pkg/k8s/apis/kdoctor.io/v1beta1"
	"github.com/kdoctor-io/kdoctor/pkg/pluginManager/tools"
	"github.com/kdoctor-io/kdoctor/pkg/types"
)

func (s *PluginNetReach) WebhookMutating(logger *zap.Logger, ctx context.Context, obj runtime.Object) error {
	req, ok := obj.(*crd.NetReach)
	if !ok {
		s := "failed to get NetReach obj"
		logger.Error(s)
		return apierrors.NewBadRequest(s)
	}

	if req.DeletionTimestamp != nil {
		return nil
	}

	if req.Spec.Target == nil {
		var agentV4Url, agentV6Url *k8sObjManager.ServiceAccessUrl
		var e error

		testIngress := false
		var agentIngress *networkingv1.Ingress
		agentIngress, e = k8sObjManager.GetK8sObjManager().GetIngress(ctx, types.ControllerConfig.Configmap.AgentIngressName, types.ControllerConfig.PodNamespace)
		if e != nil {
			logger.Sugar().Errorf("failed to get ingress , error=%v", e)
		}
		if agentIngress != nil && len(agentIngress.Status.LoadBalancer.Ingress) > 0 {
			testIngress = true
		}

		serviceAccessPortName := "http"
		testLoadBalancer := false
		if types.ControllerConfig.Configmap.EnableIPv4 {
			// TODO: AgentSerivceIpv4Name???
			agentV4Url, e = k8sObjManager.GetK8sObjManager().GetServiceAccessUrl(ctx, types.ControllerConfig.Configmap.AgentSerivceIpv4Name, types.ControllerConfig.PodNamespace, serviceAccessPortName)
			if e != nil {
				logger.Sugar().Errorf("failed to get agent ipv4 service url , error=%v", e)
			}
			if len(agentV4Url.LoadBalancerUrl) > 0 {
				testLoadBalancer = true
			}
		}
		if types.ControllerConfig.Configmap.EnableIPv6 {
			// TODO: AgentSerivceIpv6Name???
			agentV6Url, e = k8sObjManager.GetK8sObjManager().GetServiceAccessUrl(ctx, types.ControllerConfig.Configmap.AgentSerivceIpv6Name, types.ControllerConfig.PodNamespace, serviceAccessPortName)
			if e != nil {
				logger.Sugar().Errorf("failed to get agent ipv6 service url , error=%v", e)
			}
			if len(agentV6Url.LoadBalancerUrl) > 0 {
				testLoadBalancer = true
			}
		}

		enableIpv4 := types.ControllerConfig.Configmap.EnableIPv4
		enableIpv6 := types.ControllerConfig.Configmap.EnableIPv6
		enable := true
		disable := false
		m := &crd.NetReachTarget{
			Endpoint:        &enable,
			MultusInterface: &disable,
			ClusterIP:       &enable,
			NodePort:        &enable,
			LoadBalancer:    &testLoadBalancer,
			Ingress:         &testIngress,
			IPv6:            &enableIpv6,
			IPv4:            &enableIpv4,
		}
		req.Spec.Target = m
		logger.Sugar().Debugf("set default target for NetReach %v", req.Name)
	}

	if req.Spec.Schedule == nil {
		req.Spec.Schedule = tools.GetDefaultSchedule()
		logger.Sugar().Debugf("set default SchedulePlan for NetReach %v", req.Name)
	}

	if req.Spec.Request == nil {
		m := &crd.NetHttpRequest{
			DurationInSecond:      types.ControllerConfig.Configmap.NethttpDefaultRequestDurationInSecond,
			QPS:                   types.ControllerConfig.Configmap.NethttpDefaultRequestQps,
			PerRequestTimeoutInMS: types.ControllerConfig.Configmap.NethttpDefaultRequestPerRequestTimeoutInMS,
		}
		req.Spec.Request = m
		logger.Sugar().Debugf("set default Request for NetReach %v", req.Name)
	}

	if req.Spec.SuccessCondition == nil {
		req.Spec.SuccessCondition = tools.GetDefaultNetSuccessCondition()
		logger.Sugar().Debugf("set default SuccessCondition for NetReach %v", req.Name)
	}

	// agentSpec
	if true {
		if req.Spec.AgentSpec != nil {
			if req.Spec.AgentSpec.TerminationGracePeriodMinutes == nil {
				req.Spec.AgentSpec.TerminationGracePeriodMinutes = &types.ControllerConfig.Configmap.AgentDefaultTerminationGracePeriodMinutes
			}
		}
	}
	return nil
}

func (s *PluginNetReach) WebhookValidateCreate(logger *zap.Logger, ctx context.Context, obj runtime.Object) error {
	r, ok := obj.(*crd.NetReach)
	if !ok {
		s := "failed to get NetReach obj"
		logger.Error(s)
		return apierrors.NewBadRequest(s)
	}
	logger.Sugar().Debugf("NetReach: %+v", r)

	// validate Schedule
	if true {
		if err := tools.ValidataCrdSchedule(r.Spec.Schedule); err != nil {
			s := fmt.Sprintf("NetReach %v : %v", r.Name, err)
			logger.Error(s)
			return apierrors.NewBadRequest(s)
		}
	}

	// validate request
	if true {
		if r.Spec.Request.QPS >= types.ControllerConfig.Configmap.NethttpDefaultRequestMaxQps {
			s := fmt.Sprintf("NetReach %v requires qps %v bigger than maximum %v", r.Name, r.Spec.Request.QPS, types.ControllerConfig.Configmap.NethttpDefaultRequestMaxQps)
			logger.Error(s)
			return apierrors.NewBadRequest(s)
		}
		if r.Spec.Request.PerRequestTimeoutInMS > int(r.Spec.Schedule.RoundTimeoutMinute*60*1000) {
			s := fmt.Sprintf("NetReach %v requires PerRequestTimeoutInMS %v ms smaller than Schedule.RoundTimeoutMinute %vm ", r.Name, r.Spec.Request.PerRequestTimeoutInMS, r.Spec.Schedule.RoundTimeoutMinute)
			logger.Error(s)
			return apierrors.NewBadRequest(s)
		}
		if r.Spec.Request.DurationInSecond > int(r.Spec.Schedule.RoundTimeoutMinute*60) {
			s := fmt.Sprintf("NetReach %v requires request.DurationInSecond %vs smaller than Schedule.RoundTimeoutMinute %vm ", r.Name, r.Spec.Request.DurationInSecond, r.Spec.Schedule.RoundTimeoutMinute)
			logger.Error(s)
			return apierrors.NewBadRequest(s)
		}
	}

	// validate target
	if true {
		if r.Spec.Target != nil {
			// validate target
			if r.Spec.Target.IPv4 != nil && *(r.Spec.Target.IPv4) && !types.ControllerConfig.Configmap.EnableIPv4 {
				s := fmt.Sprintf("NetReach %v TestIPv4, but kdoctor ipv4 feature is disabled", r.Name)
				logger.Error(s)
				return apierrors.NewBadRequest(s)
			}
			if r.Spec.Target.IPv6 != nil && *(r.Spec.Target.IPv6) && !types.ControllerConfig.Configmap.EnableIPv6 {
				s := fmt.Sprintf("NetReach %v TestIPv6, but kdoctor ipv6 feature is disabled", r.Name)
				logger.Error(s)
				return apierrors.NewBadRequest(s)
			}
		}
	}

	// validate SuccessCondition
	if true {
		if r.Spec.SuccessCondition.SuccessRate == nil && r.Spec.SuccessCondition.MeanAccessDelayInMs == nil {
			s := fmt.Sprintf("NetReach %v, no SuccessCondition specified in the spec", r.Name)
			logger.Error(s)
			return apierrors.NewBadRequest(s)
		}
		if r.Spec.SuccessCondition.SuccessRate != nil && (*(r.Spec.SuccessCondition.SuccessRate) > 1) {
			s := fmt.Sprintf("NetReach %v, SuccessCondition.SuccessRate %v must not be bigger than 1", r.Name, *(r.Spec.SuccessCondition.SuccessRate))
			logger.Error(s)
			return apierrors.NewBadRequest(s)
		}
		if r.Spec.SuccessCondition.SuccessRate != nil && (*(r.Spec.SuccessCondition.SuccessRate) < 0) {
			s := fmt.Sprintf("NetReach %v, SuccessCondition.SuccessRate %v must not be smaller than 0 ", r.Name, *(r.Spec.SuccessCondition.SuccessRate))
			logger.Error(s)
			return apierrors.NewBadRequest(s)
		}
	}

	// validate AgentSpec
	if true {
		if r.Spec.AgentSpec != nil {
			if !slices.Contains(types.TaskRuntimes, r.Spec.AgentSpec.Kind) {
				return apierrors.NewBadRequest(fmt.Sprintf("Invalid agent runtime kind %s", r.Spec.AgentSpec.Kind))
			}
		}
	}

	return nil
}

// this will not be called, it is not allowed to modify crd
func (s *PluginNetReach) WebhookValidateUpdate(logger *zap.Logger, ctx context.Context, oldObj, newObj runtime.Object) error {
	oldNetReach := oldObj.(*crd.NetReach)
	newNetReach := newObj.(*crd.NetReach)

	if !reflect.DeepEqual(oldNetReach.Spec, newNetReach.Spec) {
		return apierrors.NewBadRequest(fmt.Sprintf("it's not allowed to modify NetReach %s Spec", oldNetReach.Name))
	}

	return nil
}
