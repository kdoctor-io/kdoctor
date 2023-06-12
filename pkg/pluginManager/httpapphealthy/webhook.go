// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package httpapphealthy

import (
	"context"
	"fmt"
	k8sObjManager "github.com/kdoctor-io/kdoctor/pkg/k8ObjManager"
	crd "github.com/kdoctor-io/kdoctor/pkg/k8s/apis/kdoctor.io/v1beta1"
	"github.com/kdoctor-io/kdoctor/pkg/pluginManager/tools"
	"github.com/kdoctor-io/kdoctor/pkg/types"
	"github.com/kdoctor-io/kdoctor/pkg/utils"
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
		// TODO (ii2day): validate host body method header
		if r.Spec.Target == nil {
			s := fmt.Sprintf("HttpAppHealthy %v, no target specified in the spec", r.Name)
			logger.Error(s)
			return apierrors.NewBadRequest(s)
		}

		// tls
		if r.Spec.Target.TlsSecret != nil {
			name, namespace, err := utils.GetObjNameNamespace(*r.Spec.Target.TlsSecret)
			if err != nil {
				s := fmt.Sprintf("HttpAppHealthy %v requires Target.TlsCert enter correctly err: %v", r.Name, err)
				logger.Error(s)
				return apierrors.NewBadRequest(s)
			}
			tlsData, err := k8sObjManager.GetK8sObjManager().GetSecret(ctx, name, namespace)
			if err != nil {
				s := fmt.Sprintf("HttpAppHealthy %v failed get secret %s err: %v", r.Name, *r.Spec.Target.TlsSecret, err)
				logger.Error(s)
				return apierrors.NewBadRequest(s)
			}

			for k, v := range tlsData.Data {
				switch k {
				case "ca.crt":
					if len(v) == 0 {
						s := fmt.Sprintf("HttpAppHealthy %v get tls secret %s ca.crt value is nil", r.Name, *r.Spec.Target.TlsSecret)
						logger.Error(s)
						return apierrors.NewBadRequest(s)
					}
				case "tls.crt":
					if len(v) == 0 {
						s := fmt.Sprintf("HttpAppHealthy %v get tls secret %s tls.crt value is nil", r.Name, *r.Spec.Target.TlsSecret)
						logger.Error(s)
						return apierrors.NewBadRequest(s)
					}
				case "tls.key":
					if len(v) == 0 {
						s := fmt.Sprintf("HttpAppHealthy %v get tls secret %s tls.key value is nil", r.Name, *r.Spec.Target.TlsSecret)
						logger.Error(s)
						return apierrors.NewBadRequest(s)
					}
				default:
					s := fmt.Sprintf("HttpAppHealthy %v get tls secret %s key %s,keys other than ca.crt, tls.crt, and tls.key cannot be used", r.Name, v, *r.Spec.Target.TlsSecret)
					logger.Error(s)
					return apierrors.NewBadRequest(s)
				}
			}

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
