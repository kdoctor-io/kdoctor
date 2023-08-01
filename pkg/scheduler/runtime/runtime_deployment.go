// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package runtime

import (
	"context"

	"go.uber.org/zap"
	appsv1 "k8s.io/api/apps/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type runtimeDeployment struct {
	client    client.Client
	apiReader client.Reader
	Namespace string
	Name      string

	log *zap.Logger
}

func NewDeploymentRuntime(c client.Client, apiReader client.Reader, namespace, name string, log *zap.Logger) TaskRuntime {
	rd := &runtimeDeployment{
		client:    c,
		apiReader: apiReader,
		Namespace: namespace,
		Name:      name,
		log:       log,
	}
	return rd
}

func (rd *runtimeDeployment) IsReady(ctx context.Context) bool {
	var deploy appsv1.Deployment

	err := rd.apiReader.Get(ctx, types.NamespacedName{
		Namespace: rd.Namespace,
		Name:      rd.Name,
	}, &deploy)
	if nil != err {
		rd.log.Sugar().Errorf("failed to get deployment %s/%s, error: %v", rd.Namespace, rd.Name, err)
		return false
	}

	if deploy.Spec.Replicas != nil {
		if *deploy.Spec.Replicas == deploy.Status.ReadyReplicas {
			return true
		}
	}

	return false
}

func (rd *runtimeDeployment) Delete(ctx context.Context) error {
	var deploy appsv1.Deployment

	err := rd.apiReader.Get(ctx, types.NamespacedName{
		Namespace: rd.Namespace,
		Name:      rd.Name,
	}, &deploy)
	if nil != err {
		if apierrors.IsNotFound(err) {
			return nil
		}
		return err
	}

	if deploy.DeletionTimestamp != nil {
		return nil
	}

	err = rd.client.Delete(ctx, &deploy)
	if nil != err {
		return err
	}

	return nil
}
