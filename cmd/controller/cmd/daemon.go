// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"context"
	"os"

	"go.opentelemetry.io/otel/attribute"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/klog/v2"
	"path/filepath"

	"github.com/kdoctor-io/kdoctor/pkg/apiserver"
	"github.com/kdoctor-io/kdoctor/pkg/debug"
	"github.com/kdoctor-io/kdoctor/pkg/pluginManager"
	"github.com/kdoctor-io/kdoctor/pkg/types"
)

func SetupUtility() {

	// run gops
	d := debug.New(rootLogger)
	if types.ControllerConfig.GopsPort != 0 {
		d.RunGops(int(types.ControllerConfig.GopsPort))
	}

	if types.ControllerConfig.PyroscopeServerAddress != "" {
		d.RunPyroscope(types.ControllerConfig.PyroscopeServerAddress, types.ControllerConfig.PodName)
	}
}

func DaemonMain() {

	rootLogger.Sugar().Infof("config: %+v", types.ControllerConfig)

	SetupUtility()

	// SetupHttpServer()

	// ------

	RunMetricsServer(types.ControllerConfig.PodName)
	MetricGaugeEndpoint.Add(context.Background(), 100)
	MetricGaugeEndpoint.Add(context.Background(), -10)
	MetricGaugeEndpoint.Add(context.Background(), 5)

	attrs := []attribute.KeyValue{
		attribute.Key("pod1").String("value1"),
	}
	MetricCounterRequest.Add(context.Background(), 10, attrs...)
	attrs = []attribute.KeyValue{
		attribute.Key("pod2").String("value1"),
	}
	MetricCounterRequest.Add(context.Background(), 5, attrs...)

	MetricHistogramDuration.Record(context.Background(), 10)
	MetricHistogramDuration.Record(context.Background(), 20)

	// ----------
	s := pluginManager.InitPluginManager(rootLogger.Named("pluginsManager"))
	s.RunControllerController(int(types.ControllerConfig.HttpPort), int(types.ControllerConfig.WebhookPort), filepath.Dir(types.ControllerConfig.TlsServerCertPath))

	// ------------
	rootLogger.Info("finish kdoctor-controller initialization")

	// start apiserver
	stopCh := genericapiserver.SetupSignalHandler()
	apiserverConfig, err := apiserver.NewkdoctorServerOptions().Config()
	if nil != err {
		rootLogger.Sugar().Fatal("Error creating server configuration: %v", err)
	}
	server, err := apiserverConfig.Complete().New()
	if nil != err {
		rootLogger.Sugar().Fatal("Error creating server: %v", err)
	}

	rootLogger.Info("running kdoctor-apiserver")
	err = server.Run(stopCh)
	if nil != err {
		klog.Errorf("Error creating server: %v", err)
		os.Exit(1)
	}

	// sleep forever
	// select {}
}
