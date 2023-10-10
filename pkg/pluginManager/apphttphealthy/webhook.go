// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package apphttphealthy

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"go.uber.org/zap"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/strings/slices"

	k8sObjManager "github.com/kdoctor-io/kdoctor/pkg/k8ObjManager"
	crd "github.com/kdoctor-io/kdoctor/pkg/k8s/apis/kdoctor.io/v1beta1"
	"github.com/kdoctor-io/kdoctor/pkg/pluginManager/tools"
	"github.com/kdoctor-io/kdoctor/pkg/types"
)

func (s *PluginAppHttpHealthy) WebhookMutating(logger *zap.Logger, ctx context.Context, obj runtime.Object) error {
	req, ok := obj.(*crd.AppHttpHealthy)
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

	// target
	if true {
		protoclHttp := strings.Contains(req.Spec.Target.Host, "http")
		protoclHttps := strings.Contains(req.Spec.Target.Host, "https")
		// default http
		if !protoclHttp && !protoclHttps {
			req.Spec.Target.Host = fmt.Sprintf("http://%s", req.Spec.Target.Host)
			logger.Sugar().Debugf("set default target host protocol http for HttpAppHealthy %v", req.Name)
		}
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

func (s *PluginAppHttpHealthy) WebhookValidateCreate(logger *zap.Logger, ctx context.Context, obj runtime.Object) error {
	r, ok := obj.(*crd.AppHttpHealthy)
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
		if r.Spec.Target.TlsSecretName != nil {
			tlsData, err := k8sObjManager.GetK8sObjManager().GetSecret(ctx, *r.Spec.Target.TlsSecretName, *r.Spec.Target.TlsSecretNamespace)
			if err != nil {
				s := fmt.Sprintf("HttpAppHealthy %v failed get secret %s/%s err: %v", r.Name, *r.Spec.Target.TlsSecretNamespace, *r.Spec.Target.TlsSecretName, err)
				logger.Error(s)
				return apierrors.NewBadRequest(s)
			}

			for k, v := range tlsData.Data {
				switch k {
				case "ca.crt":
					if len(v) == 0 {
						s := fmt.Sprintf("HttpAppHealthy %v get tls secret %s/%s ca.crt value is nil", r.Name, *r.Spec.Target.TlsSecretNamespace, *r.Spec.Target.TlsSecretName)
						logger.Error(s)
						return apierrors.NewBadRequest(s)
					}
				case "tls.crt":
					if len(v) == 0 {
						s := fmt.Sprintf("HttpAppHealthy %v get tls secret %s/%s tls.crt value is nil", r.Name, *r.Spec.Target.TlsSecretNamespace, *r.Spec.Target.TlsSecretName)
						logger.Error(s)
						return apierrors.NewBadRequest(s)
					}
				case "tls.key":
					if len(v) == 0 {
						s := fmt.Sprintf("HttpAppHealthy %v get tls secret %s/%s tls.key value is nil", r.Name, *r.Spec.Target.TlsSecretNamespace, *r.Spec.Target.TlsSecretName)
						logger.Error(s)
						return apierrors.NewBadRequest(s)
					}
				default:
					s := fmt.Sprintf("HttpAppHealthy %v get tls secret %s/%s key %s,keys other than ca.crt, tls.crt, and tls.key cannot be used", r.Name, *r.Spec.Target.TlsSecretNamespace, *r.Spec.Target.TlsSecretName, v)
					logger.Error(s)
					return apierrors.NewBadRequest(s)
				}
			}

		}

		// host is ip or domain
		err := tools.ValidataAppHttpHealthyHost(r)
		if err != nil {
			return err
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
func (s *PluginAppHttpHealthy) WebhookValidateUpdate(logger *zap.Logger, ctx context.Context, oldObj, newObj runtime.Object) error {
	oldHealthy := oldObj.(*crd.AppHttpHealthy)
	newHealthy := newObj.(*crd.AppHttpHealthy)

	if !reflect.DeepEqual(oldHealthy.Spec, newHealthy.Spec) {
		return apierrors.NewBadRequest(fmt.Sprintf("it's not allowed to modify AppHttpHealthy %s Spec", oldHealthy.Name))
	}

	return nil
}
