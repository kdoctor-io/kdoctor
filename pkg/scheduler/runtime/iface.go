// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package runtime

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kdoctor-io/kdoctor/pkg/types"
)

type TaskRuntime interface {
	IsReady(ctx context.Context) bool
	Delete(ctx context.Context) error
}

func FindRuntime(client client.Client, apiReader client.Reader, runtimeKind, runtimeName string, log *zap.Logger) (TaskRuntime, error) {
	var tr TaskRuntime
	var err error

	switch runtimeKind {
	case types.KindDeployment:
		tr = NewDeploymentRuntime(client, apiReader, types.ControllerConfig.PodNamespace, runtimeName, log)
	case types.KindDaemonSet:
		tr = NewDaemonSetRuntime(client, apiReader, types.ControllerConfig.PodNamespace, runtimeName, log)
	default:
		err = fmt.Errorf("unrecognized runtime kind %s for %s", runtimeKind, runtimeName)
	}

	return tr, err
}
