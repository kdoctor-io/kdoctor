// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package httpapphealthy

import (
	"context"
	"fmt"
	crd "github.com/kdoctor-io/kdoctor/pkg/k8s/apis/kdoctor.io/v1beta1"
	"github.com/kdoctor-io/kdoctor/pkg/pluginManager/tools"
	"github.com/kdoctor-io/kdoctor/pkg/types"
	"go.uber.org/zap"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

func (s *PluginHttpAppHealthy) WebhookMutating(logger *zap.Logger, ctx context.Context, obj runtime.Object) error {
	req, ok := obj.(*crd.HttpAppHealthy)
	if !ok {
		s := "failed to get HttpAppHealthy obj"
		logger.Error(s)
		return apierrors.NewBadRequest(s)
	}

	if req.DeletionTimestamp != nil {
		return nil
	}

	if req.Spec.Schedule == nil {
		req.Spec.Schedule = tools.GetDefaultSchedule()
		logger.Sugar().Debugf("set default SchedulePlan for HttpAppHealthy %v", req.Name)
	}

	if req.Spec.Request == nil {
		m := &crd.NetHttpRequest{
			DurationInSecond:      types.ControllerConfig.Configmap.NethttpDefaultRequestDurationInSecond,
			QPS:                   types.ControllerConfig.Configmap.NethttpDefaultRequestQps,
			PerRequestTimeoutInMS: types.ControllerConfig.Configmap.NethttpDefaultRequestPerRequestTimeoutInMS,
		}
		req.Spec.Request = m
		logger.Sugar().Debugf("set default Request for HttpAppHealthy %v", req.Name)
	}

	if req.Spec.SuccessCondition == nil {
		req.Spec.SuccessCondition = tools.GetDefaultNetSuccessCondition()
		logger.Sugar().Debugf("set default SuccessCondition for HttpAppHealthy %v", req.Name)
	}

	return nil
}

func (s *PluginHttpAppHealthy) WebhookValidateCreate(logger *zap.Logger, ctx context.Context, obj runtime.Object) error {
	r, ok := obj.(*crd.HttpAppHealthy)
	if !ok {
		s := "failed to get HttpAppHealthy obj"
		logger.Error(s)
		return apierrors.NewBadRequest(s)
	}
	logger.Sugar().Debugf("HttpAppHealthy: %+v", r)

	// validate Schedule
	if true {
		if err := tools.ValidataCrdSchedule(r.Spec.Schedule); err != nil {
			s := fmt.Sprintf("HttpAppHealthy %v : %v", r.Name, err)
			logger.Error(s)
			return apierrors.NewBadRequest(s)
		}
	}

	// validate request
	if true {
		if r.Spec.Request.QPS >= types.ControllerConfig.Configmap.NethttpDefaultRequestMaxQps {
			s := fmt.Sprintf("HttpAppHealthy %v requires qps %v bigger than maximum %v", r.Name, r.Spec.Request.QPS, types.ControllerConfig.Configmap.NethttpDefaultRequestMaxQps)
			logger.Error(s)
			return apierrors.NewBadRequest(s)
		}
		if r.Spec.Request.PerRequestTimeoutInMS > int(r.Spec.Schedule.RoundTimeoutMinute*60*1000) {
			s := fmt.Sprintf("HttpAppHealthy %v requires PerRequestTimeoutInMS %v ms smaller than Schedule.RoundTimeoutMinute %vm ", r.Name, r.Spec.Request.PerRequestTimeoutInMS, r.Spec.Schedule.RoundTimeoutMinute)
			logger.Error(s)
			return apierrors.NewBadRequest(s)
		}
		if r.Spec.Request.DurationInSecond > int(r.Spec.Schedule.RoundTimeoutMinute*60) {
			s := fmt.Sprintf("HttpAppHealthy %v requires request.DurationInSecond %vs smaller than Schedule.RoundTimeoutMinute %vm ", r.Name, r.Spec.Request.DurationInSecond, r.Spec.Schedule.RoundTimeoutMinute)
			logger.Error(s)
			return apierrors.NewBadRequest(s)
		}
	}

	// validate target
	if true {
		// TODO (ii2day) validate host body method tls header
		if r.Spec.Target == nil {
			s := fmt.Sprintf("HttpAppHealthy %v, no target specified in the spec", r.Name)
			logger.Error(s)
			return apierrors.NewBadRequest(s)
		}
	}

	// validate SuccessCondition
	if true {
		if r.Spec.SuccessCondition.SuccessRate == nil && r.Spec.SuccessCondition.MeanAccessDelayInMs == nil {
			s := fmt.Sprintf("HttpAppHealthy %v, no SuccessCondition specified in the spec", r.Name)
			logger.Error(s)
			return apierrors.NewBadRequest(s)
		}
		if r.Spec.SuccessCondition.SuccessRate != nil && (*(r.Spec.SuccessCondition.SuccessRate) > 1) {
			s := fmt.Sprintf("HttpAppHealthy %v, SuccessCondition.SuccessRate %v must not be bigger than 1", r.Name, *(r.Spec.SuccessCondition.SuccessRate))
			logger.Error(s)
			return apierrors.NewBadRequest(s)
		}
		if r.Spec.SuccessCondition.SuccessRate != nil && (*(r.Spec.SuccessCondition.SuccessRate) < 0) {
			s := fmt.Sprintf("HttpAppHealthy %v, SuccessCondition.SuccessRate %v must not be smaller than 0 ", r.Name, *(r.Spec.SuccessCondition.SuccessRate))
			logger.Error(s)
			return apierrors.NewBadRequest(s)
		}
	}

	return nil
}

// this will not be called, it is not allowed to modify crd
func (s *PluginHttpAppHealthy) WebhookValidateUpdate(logger *zap.Logger, ctx context.Context, oldObj, newObj runtime.Object) error {

	return nil
}
