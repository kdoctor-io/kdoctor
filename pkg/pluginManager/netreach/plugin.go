// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package netreach

import (
	crd "github.com/kdoctor-io/kdoctor/pkg/k8s/apis/kdoctor.io/v1beta1"
	"github.com/kdoctor-io/kdoctor/pkg/pluginManager/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type PluginNetReach struct {
}

var _ types.ChainingPlugin = &PluginNetReach{}

func (s *PluginNetReach) GetApiType() client.Object {
	return &crd.NetReach{}
}
