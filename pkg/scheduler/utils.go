// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package scheduler

import (
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	k8svalidation "k8s.io/apimachinery/pkg/util/validation"
)

var UniqueMatchLabelValue = TaskRuntimeName

// TaskRuntimeName generates a unique name for task runtime.
// notice: different kind tasks could use the same CR object name, so we need add their kind to generate name.
func TaskRuntimeName(taskKind, taskName string) string {
	taskRuntimeName := fmt.Sprintf("%s-%s-%s", kdoctor, strings.ToLower(taskKind), taskName)
	if len(taskRuntimeName) > k8svalidation.DNS1123SubdomainMaxLength {
		taskRuntimeName = taskRuntimeName[:k8svalidation.DNS1123SubdomainMaxLength]
	}

	return taskRuntimeName
}

// TaskRuntimeServiceName generates a service name for the corresponding task runtime.
func TaskRuntimeServiceName(taskRuntimeName string, ipFamily corev1.IPFamily) string {
	// "*-ipv4" or "*-ipv6"
	prefixLen := k8svalidation.DNS1123SubdomainMaxLength - 5

	if len(taskRuntimeName) >= prefixLen {
		taskRuntimeName = taskRuntimeName[:prefixLen]
	}

	taskRuntimeName = fmt.Sprintf("%s-%s", taskRuntimeName, strings.ToLower(string(ipFamily)))
	return taskRuntimeName
}

// AppendAnnotationOrLabel will combine two annotations, and the later one will overwrite the same key-value.
func AppendAnnotationOrLabel(origin, addition map[string]string) map[string]string {
	if origin == nil {
		origin = make(map[string]string)
	}

	for key, val := range addition {
		origin[key] = val
	}

	return origin
}
