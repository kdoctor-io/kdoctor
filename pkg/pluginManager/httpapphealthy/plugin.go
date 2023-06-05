// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package httpapphealthy

import (
	crd "github.com/kdoctor-io/kdoctor/pkg/k8s/apis/kdoctor.io/v1beta1"
	"github.com/kdoctor-io/kdoctor/pkg/pluginManager/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type PluginHttpAppHealthy struct {
}

var _ types.ChainingPlugin = &PluginHttpAppHealthy{}

func (s *PluginHttpAppHealthy) GetApiType() client.Object {
	return &crd.HttpAppHealthy{}
}
