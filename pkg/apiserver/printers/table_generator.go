// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package printers

import (
	"context"

	metatable "k8s.io/apimachinery/pkg/api/meta/table"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/registry/rest"

	"github.com/kdoctor-io/kdoctor/pkg/k8s/apis/system/v1beta1"
)

var (
	tableColumnDefinitions = []metav1.TableColumnDefinition{
		{Name: "Name", Type: "string", Format: "name", Description: swaggerMetadataDescriptions["name"]},
		{Name: "RoundNumber", Type: "int", Description: ""},
		{Name: "RoundResult", Type: "string", Description: ""},
		{Name: "NodeName", Type: "string", Description: ""},
		{Name: "Started At", Type: "date", Description: swaggerMetadataDescriptions["startTimestamp"]},
	}
)

var _ rest.TableConvertor = &TableGenerator{}

type TableGenerator struct {
	defaultQualifiedResource schema.GroupResource
}

func NewTableGenerator(defaultQualifiedResource schema.GroupResource) *TableGenerator {
	return &TableGenerator{defaultQualifiedResource: defaultQualifiedResource}
}

var swaggerMetadataDescriptions = metav1.ObjectMeta{}.SwaggerDoc()

func (t TableGenerator) ConvertToTable(ctx context.Context, obj runtime.Object, tableOptions runtime.Object) (*metav1.Table, error) {
	table := &metav1.Table{
		ColumnDefinitions: tableColumnDefinitions,
	}

	var err error
	table.Rows, err = metatable.MetaToTableRow(obj, func(obj runtime.Object, m metav1.Object, name, age string) ([]interface{}, error) {
		pluginReport := obj.(*v1beta1.KdoctorReport)
		return []interface{}{
			name,
			pluginReport.Spec.ToTalRoundNumber,
			// pluginReport.Spec.StartTimeStamp.Time.UTC().Format(time.RFC3339),
		}, nil
	})
	return table, err
}
