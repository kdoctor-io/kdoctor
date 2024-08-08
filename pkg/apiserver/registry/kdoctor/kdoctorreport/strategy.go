// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package kdoctorreport

import (
	"context"
	"errors"

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

type kdoctorReportStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

var _ rest.RESTCreateStrategy = &kdoctorReportStrategy{}
var _ rest.RESTUpdateStrategy = &kdoctorReportStrategy{}

func NewStrategy(typer runtime.ObjectTyper) kdoctorReportStrategy {
	return kdoctorReportStrategy{
		ObjectTyper:   typer,
		NameGenerator: names.SimpleNameGenerator,
	}
}

func (p kdoctorReportStrategy) NamespaceScoped() bool {
	return true
}

func (p kdoctorReportStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
}

func (p kdoctorReportStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
}

func (p kdoctorReportStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	return field.ErrorList{}
}

func (p kdoctorReportStrategy) WarningsOnCreate(ctx context.Context, obj runtime.Object) []string {
	return []string{}
}

func (p kdoctorReportStrategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	return []string{}
}

func (p kdoctorReportStrategy) AllowCreateOnUpdate() bool {
	return false
}

func (p kdoctorReportStrategy) AllowUnconditionalUpdate() bool {
	return false
}

func (p kdoctorReportStrategy) Canonicalize(obj runtime.Object) {
}

func (p kdoctorReportStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return field.ErrorList{}
}

func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	kdoctorReport, ok := obj.(*v1beta1.KdoctorReport)
	if !ok {
		return nil, nil, errors.New("given object is not a KdoctorReport")
	}
	return labels.Set(kdoctorReport.ObjectMeta.Labels), SelectableFields(kdoctorReport), nil
}

func MatchKdoctorReport(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
	}
}

func SelectableFields(obj *v1beta1.KdoctorReport) fields.Set {
	return generic.ObjectMetaFieldsSet(&obj.ObjectMeta, true)
}
