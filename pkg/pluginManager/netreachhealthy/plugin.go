// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package netreachhealthy

import (
	crd "github.com/kdoctor-io/kdoctor/pkg/k8s/apis/kdoctor.io/v1beta1"
	"github.com/kdoctor-io/kdoctor/pkg/pluginManager/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type PluginNetReachHealthy struct {
}

var _ types.ChainingPlugin = &PluginNetReachHealthy{}

func (s *PluginNetReachHealthy) GetApiType() client.Object {
	return &crd.NetReachHealthy{}
}
