// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package k8s

import (
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	typedv1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
)

func NewEventRecord(scheme *runtime.Scheme, EventSourceName string, logger *zap.Logger) record.EventRecorder {
	// ------------- for generate event for the crd
	config, err := rest.InClusterConfig()
	if err != nil {
		logger.Sugar().Fatalf("failed to InClusterConfig, reason=%v", err)
	}
	clientset, err := kubernetes.NewForConfig(config) // 初始化 client
	if err != nil {
		logger.Sugar().Fatalf("failed to NewForConfig, reason=%v", err)
	}

	eventBroadcaster := record.NewBroadcaster()

	eventBroadcaster.StartLogging(logger.Named("event").Sugar().Infof)
	eventBroadcaster.StartRecordingToSink(&typedv1.EventSinkImpl{
		Interface: typedv1.New(clientset.CoreV1().RESTClient()).Events(""),
	})

	return eventBroadcaster.NewRecorder(scheme,
		corev1.EventSource{
			Component: EventSourceName,
		})
}
