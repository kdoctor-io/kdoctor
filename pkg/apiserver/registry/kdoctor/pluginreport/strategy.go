// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package pluginreport

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/names"

	"github.com/kdoctor-io/kdoctor/pkg/k8s/apis/system/v1beta1"
)

type pluginReportStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

var _ rest.RESTCreateStrategy = &pluginReportStrategy{}
var _ rest.RESTUpdateStrategy = &pluginReportStrategy{}

func NewStrategy(typer runtime.ObjectTyper) pluginReportStrategy {
	return pluginReportStrategy{
		ObjectTyper:   typer,
		NameGenerator: names.SimpleNameGenerator,
	}
}

func (p pluginReportStrategy) NamespaceScoped() bool {
	return true
}

func (p pluginReportStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
}

func (p pluginReportStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
}

func (p pluginReportStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	return field.ErrorList{}
}

func (p pluginReportStrategy) WarningsOnCreate(ctx context.Context, obj runtime.Object) []string {
	return []string{}
}

func (p pluginReportStrategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	return []string{}
}

func (p pluginReportStrategy) AllowCreateOnUpdate() bool {
	return false
}

func (p pluginReportStrategy) AllowUnconditionalUpdate() bool {
	return false
}

func (p pluginReportStrategy) Canonicalize(obj runtime.Object) {
}

func (p pluginReportStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return field.ErrorList{}
}

func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	pluginReport, ok := obj.(*v1beta1.PluginReport)
	if !ok {
		return nil, nil, fmt.Errorf("given object is not a PluginReport")
	}
	return labels.Set(pluginReport.ObjectMeta.Labels), SelectableFields(pluginReport), nil
}

func MatchPluginReport(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
	}
}

func SelectableFields(obj *v1beta1.PluginReport) fields.Set {
	return generic.ObjectMetaFieldsSet(&obj.ObjectMeta, true)
}
