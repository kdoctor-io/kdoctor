// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kdoctor-io/kdoctor/pkg/agentDnsServer"
	"github.com/kdoctor-io/kdoctor/pkg/agentHttpServer"
	"github.com/kdoctor-io/kdoctor/pkg/debug"
	k8sObjManager "github.com/kdoctor-io/kdoctor/pkg/k8ObjManager"
	"github.com/kdoctor-io/kdoctor/pkg/pluginManager"
	"github.com/kdoctor-io/kdoctor/pkg/types"
)

func SetupUtility() {
	// run gops
	d := debug.New(rootLogger)
	if types.AgentConfig.GopsPort != 0 {
		d.RunGops(int(types.AgentConfig.GopsPort))
	}

	if types.AgentConfig.PyroscopeServerAddress != "" {
		d.RunPyroscope(types.AgentConfig.PyroscopeServerAddress, types.AgentConfig.PodName)
	}
}

func DaemonMain() {
	rootLogger.Sugar().Infof("config: %+v", types.AgentConfig)

	// TODO: udp server, tcp server, websocket server
	if types.AgentConfig.AppMode {
		// app mode, just used to debug
		rootLogger.Info("run in app mode")
		scheme := runtime.NewScheme()
		if err := clientgoscheme.AddToScheme(scheme); err != nil {
			rootLogger.Sugar().Fatalf("failed to add to scheme, reason=%v", err)
		}

		n := ctrl.Options{
			Scheme:                 scheme,
			MetricsBindAddress:     "0",
			HealthProbeBindAddress: "0",
			LeaderElection:         false,
			// for this not watched obj, get directly from api-server
			ClientDisableCacheFor: []client.Object{
				&corev1.Node{},
				&corev1.Namespace{},
				&corev1.Pod{},
				&corev1.Service{},
				&appsv1.Deployment{},
				&appsv1.StatefulSet{},
				&appsv1.ReplicaSet{},
				&appsv1.DaemonSet{},
			},
		}
		mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), n)
		if err != nil {
			rootLogger.Sugar().Fatalf("failed to NewManager, reason=%v", err)
		}

		if e := k8sObjManager.Initk8sObjManager(mgr.GetClient()); e != nil {
			rootLogger.Sugar().Fatalf("failed to Initk8sObjManager, error=%v", e)
		}
	} else {
		rootLogger.Info("run in agent mode")

		SetupUtility()

		RunMetricsServer(types.AgentConfig.PodName)

		s := pluginManager.InitPluginManager(rootLogger.Named("agentContorller"))
		s.RunAgentController()

	}
	rootLogger.Sugar().Info("generate server cert and key")
	GenServerCert(rootLogger)
	agentHttpServer.SetupHealthHttpServer(rootLogger)
	agentHttpServer.SetupAppHttpServer(rootLogger, TlsCertPath, TlsKeyPath)
	initGrpcServer()
	agentDnsServer.SetupAppDnsServer(rootLogger, TlsCertPath, TlsKeyPath)

	rootLogger.Info("finish initialization")
	// sleep forever
	select {}
}
