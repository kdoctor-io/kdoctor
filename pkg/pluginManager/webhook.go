// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package pluginManager

import (
	"context"
	"fmt"
	plugintypes "github.com/kdoctor-io/kdoctor/pkg/pluginManager/types"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// --------------------

type pluginWebhookhander struct {
	logger  *zap.Logger
	plugin  plugintypes.ChainingPlugin
	client  client.Client
	crdKind string
}

var _ webhook.CustomValidator = (*pluginWebhookhander)(nil)

const (
	MsgErrCtx = "miss admission Request in ctx "
)

// mutating webhook
func (s *pluginWebhookhander) Default(ctx context.Context, obj runtime.Object) error {
	req, err := admission.RequestFromContext(ctx)
	if err != nil {
		return fmt.Errorf("%s : %w", MsgErrCtx, err)
	}
	s.logger.Sugar().Debugf("mutating kind %v", req.Kind.Kind)

	// ------ add crd ------
	// switch s.crdKind {
	// case KindNameNethttp:
	// 	instance, ok := obj.(*crd.Nethttp)
	// 	if !ok {
	// 		s.logger.Error(ApiMsgGetFailure)
	// 		return apierrors.NewBadRequest(ApiMsgGetFailure)
	// 	}
	// 	s.logger.Sugar().Debugf("nethppt instance: %+v", instance)
	//
	// case KindNameNetdns:
	// 	instance, ok := obj.(*crd.Netdns)
	// 	if !ok {
	// 		s.logger.Error(ApiMsgGetFailure)
	// 		return apierrors.NewBadRequest(ApiMsgGetFailure)
	// 	}
	// 	s.logger.Sugar().Debugf("netdns instance: %+v", instance)
	// 	*(instance.Status.ExpectedRound) = instance.Spec.Schedule.RoundNumber
	//
	// default:
	// 	s.logger.Sugar().Errorf("%s, support kind=%v, objkind=%v, obj=%+v", ApiMsgUnknowCRD, s.crdKind, obj.GetObjectKind(), obj)
	// 	return apierrors.NewBadRequest(ApiMsgUnknowCRD)
	// }

	return s.plugin.WebhookMutating(s.logger.Named("mutatingWebhook"), ctx, obj)
}

func (s *pluginWebhookhander) ValidateCreate(ctx context.Context, obj runtime.Object) error {
	req, err := admission.RequestFromContext(ctx)
	if err != nil {
		return fmt.Errorf("%s : %w", MsgErrCtx, err)
	}
	s.logger.Sugar().Debugf("create kind %v", req.Kind.Kind)

	return s.plugin.WebhookValidateCreate(s.logger.Named("validatingCreateWebhook"), ctx, obj)
}

func (s *pluginWebhookhander) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) error {
	req, err := admission.RequestFromContext(ctx)
	if err != nil {
		return fmt.Errorf("%s : %w", MsgErrCtx, err)
	}
	s.logger.Sugar().Debugf("update kind %v", req.Kind.Kind)

	return s.plugin.WebhookValidateUpdate(s.logger.Named("validatingCreateWebhook"), ctx, oldObj, newObj)
}

// not registered in ValidatingWebhookConfiguration
func (s *pluginWebhookhander) ValidateDelete(ctx context.Context, obj runtime.Object) error {
	return nil
}

// --------------------

func (s *pluginWebhookhander) SetupWebhook(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).For(s.plugin.GetApiType()).WithDefaulter(s).WithValidator(s).RecoverPanic().Complete()
}
