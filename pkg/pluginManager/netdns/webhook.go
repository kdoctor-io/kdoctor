// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package netdns

import (
	"context"
	"fmt"
	"reflect"

	"go.uber.org/zap"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/strings/slices"

	crd "github.com/kdoctor-io/kdoctor/pkg/k8s/apis/kdoctor.io/v1beta1"
	"github.com/kdoctor-io/kdoctor/pkg/pluginManager/tools"
	"github.com/kdoctor-io/kdoctor/pkg/types"
)

func (s *PluginNetDns) WebhookMutating(logger *zap.Logger, ctx context.Context, obj runtime.Object) error {
	r, ok := obj.(*crd.Netdns)
	if !ok {
		s := "failed to get netnds obj"
		logger.Error(s)
		return apierrors.NewBadRequest(s)
	}
	logger.Sugar().Infof("obj: %+v", r)

	// agentSpec
	// agentSpec
	if true {
		if r.Spec.AgentSpec != nil {
			if r.Spec.AgentSpec.TerminationGracePeriodMinutes == nil {
				r.Spec.AgentSpec.TerminationGracePeriodMinutes = &types.ControllerConfig.Configmap.AgentDefaultTerminationGracePeriodMinutes
			}
		}
	}
	// TODO: mutating default value

	return nil
}

func (s *PluginNetDns) WebhookValidateCreate(logger *zap.Logger, ctx context.Context, obj runtime.Object) error {
	r, ok := obj.(*crd.Netdns)
	if !ok {
		s := "failed to get netdns obj"
		logger.Error(s)
		return apierrors.NewBadRequest(s)
	}
	logger.Sugar().Infof("obj: %+v", r)

	// validate Schedule
	if true {
		if err := tools.ValidataCrdSchedule(r.Spec.Schedule); err != nil {
			s := fmt.Sprintf("netdns %v : %v", r.Name, err)
			logger.Error(s)
			return apierrors.NewBadRequest(s)
		}
	}

	// validate request
	if true {

		if int(*r.Spec.Request.QPS) >= types.ControllerConfig.Configmap.NethttpDefaultRequestMaxQps {
			s := fmt.Sprintf("netdns %v requires qps %v bigger than maximum %v", r.Name, r.Spec.Request.QPS, types.ControllerConfig.Configmap.NethttpDefaultRequestMaxQps)
			logger.Error(s)
			return apierrors.NewBadRequest(s)
		}
		if int(*r.Spec.Request.PerRequestTimeoutInMS) > int(r.Spec.Schedule.RoundTimeoutMinute*60*1000) {
			s := fmt.Sprintf("netdns %v requires PerRequestTimeoutInMS %v ms smaller than Schedule.RoundTimeoutMinute %vm ", r.Name, r.Spec.Request.PerRequestTimeoutInMS, r.Spec.Schedule.RoundTimeoutMinute)
			logger.Error(s)
			return apierrors.NewBadRequest(s)
		}
		if int(*r.Spec.Request.DurationInSecond) > int(r.Spec.Schedule.RoundTimeoutMinute*60) {
			s := fmt.Sprintf("netdns %v requires request.DurationInSecond %vs smaller than Schedule.RoundTimeoutMinute %vm ", r.Name, r.Spec.Request.DurationInSecond, r.Spec.Schedule.RoundTimeoutMinute)
			logger.Error(s)
			return apierrors.NewBadRequest(s)
		}
	}

	// validate target
	if true {
		if r.Spec.Target.NetDnsTargetDns != nil && r.Spec.Target.NetDnsTargetUser != nil {
			s := fmt.Sprintf("netdns %v, forbid to set TargetUser and TargetDns at same time", r.Name)
			logger.Error(s)
			return apierrors.NewBadRequest(s)
		}
	}

	// validate SuccessCondition
	if true {
		if r.Spec.SuccessCondition.SuccessRate == nil && r.Spec.SuccessCondition.MeanAccessDelayInMs == nil {
			s := fmt.Sprintf("netdns %v, no SuccessCondition specified in the spec", r.Name)
			logger.Error(s)
			return apierrors.NewBadRequest(s)
		}
		if r.Spec.SuccessCondition.SuccessRate != nil && (*(r.Spec.SuccessCondition.SuccessRate) > 1) {
			s := fmt.Sprintf("netdns %v, SuccessCondition.SuccessRate %v must not be bigger than 1", r.Name, *(r.Spec.SuccessCondition.SuccessRate))
			logger.Error(s)
			return apierrors.NewBadRequest(s)
		}
		if r.Spec.SuccessCondition.SuccessRate != nil && (*(r.Spec.SuccessCondition.SuccessRate) < 0) {
			s := fmt.Sprintf("netdns %v, SuccessCondition.SuccessRate %v must not be smaller than 0 ", r.Name, *(r.Spec.SuccessCondition.SuccessRate))
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

func (s *PluginNetDns) WebhookValidateUpdate(logger *zap.Logger, ctx context.Context, oldObj, newObj runtime.Object) error {
	oldNetdns := oldObj.(*crd.Netdns)
	newNetdns := newObj.(*crd.Netdns)

	if !reflect.DeepEqual(oldNetdns.Spec, newNetdns.Spec) {
		return apierrors.NewBadRequest(fmt.Sprintf("it's not allowed to modify Netdns %s Spec", oldNetdns.Name))
	}

	return nil
}
