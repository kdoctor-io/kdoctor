// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package apphttphealthy

import (
	crd "github.com/kdoctor-io/kdoctor/pkg/k8s/apis/kdoctor.io/v1beta1"
	"github.com/kdoctor-io/kdoctor/pkg/pluginManager/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type PluginAppHttpHealthy struct {
}

var _ types.ChainingPlugin = &PluginAppHttpHealthy{}

func (s *PluginAppHttpHealthy) GetApiType() client.Object {
	return &crd.AppHttpHealthy{}
}
