// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package k8sObjManager

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (nm *k8sObjManager) GetConfigMap(ctx context.Context, name, namespace string) (*corev1.ConfigMap, error) {
	d := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	key := client.ObjectKeyFromObject(d)
	if e := nm.client.Get(ctx, key, d); e != nil {
		return nil, fmt.Errorf("failed to get configmap %v/%v, reason=%v", namespace, name, e)
	}
	return d, nil
}
