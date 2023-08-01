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

type runtimeDaemonSet struct {
	client    client.Client
	apiReader client.Reader
	Namespace string
	Name      string

	log *zap.Logger
}

func NewDaemonSetRuntime(c client.Client, apiReader client.Reader, namespace, name string, log *zap.Logger) TaskRuntime {
	rd := &runtimeDaemonSet{
		client:    c,
		apiReader: apiReader,
		Namespace: namespace,
		Name:      name,
		log:       log,
	}

	return rd
}

func (rd *runtimeDaemonSet) IsReady(ctx context.Context) bool {
	var daemonSet appsv1.DaemonSet

	err := rd.apiReader.Get(ctx, types.NamespacedName{
		Namespace: rd.Namespace,
		Name:      rd.Name,
	}, &daemonSet)
	if nil != err {
		rd.log.Sugar().Errorf("failed to get deployment %s/%s, error: %v", rd.Namespace, rd.Name, err)
		return false
	}

	if daemonSet.Status.DesiredNumberScheduled == daemonSet.Status.NumberReady {
		return true
	}

	return false
}

func (rd *runtimeDaemonSet) Delete(ctx context.Context) error {
	var daemonSet appsv1.DaemonSet

	err := rd.apiReader.Get(ctx, types.NamespacedName{
		Namespace: rd.Namespace,
		Name:      rd.Name,
	}, &daemonSet)
	if nil != err {
		if apierrors.IsNotFound(err) {
			return nil
		}
		return err
	}

	if daemonSet.DeletionTimestamp != nil {
		return nil
	}

	err = rd.client.Delete(ctx, &daemonSet)
	if nil != err {
		return err
	}

	return nil
}
