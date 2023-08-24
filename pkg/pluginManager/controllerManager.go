// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package pluginManager

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/healthz"

	"github.com/kdoctor-io/kdoctor/pkg/fileManager"
	k8sObjManager "github.com/kdoctor-io/kdoctor/pkg/k8ObjManager"
	crd "github.com/kdoctor-io/kdoctor/pkg/k8s/apis/kdoctor.io/v1beta1"
	"github.com/kdoctor-io/kdoctor/pkg/reportManager"
	"github.com/kdoctor-io/kdoctor/pkg/scheduler"
	"github.com/kdoctor-io/kdoctor/pkg/types"
)

func (s *pluginManager) RunControllerController(healthPort int, webhookPort int, webhookTlsDir string) {

	logger := s.logger
	scheme := runtime.NewScheme()
	if e := clientgoscheme.AddToScheme(scheme); e != nil {
		logger.Sugar().Fatalf("failed to add k8s scheme, reason=%v", e)
	}
	if e := crd.AddToScheme(scheme); e != nil {
		logger.Sugar().Fatalf("failed to add scheme for plugins, reason=%v", e)
	}
	if e := networkingv1.AddToScheme(scheme); e != nil {
		logger.Sugar().Fatalf("failed to add scheme for plugins, reason=%v", e)
	}
	n := ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: "0",
		// health
		HealthProbeBindAddress: "0",
		// webhook
		Port:    webhookPort,
		CertDir: webhookTlsDir,
		// lease
		LeaderElection:          true,
		LeaderElectionNamespace: types.ControllerConfig.PodNamespace,
		LeaderElectionID:        types.ControllerConfig.PodName,
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
	if healthPort != 0 {
		n.HealthProbeBindAddress = fmt.Sprintf(":%d", healthPort)
		n.ReadinessEndpointName = "/healthy/readiness"
		n.LivenessEndpointName = "/healthy/liveness"
	}
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), n)
	if err != nil {
		logger.Sugar().Fatalf("failed to NewManager, reason=%v", err)
	}
	if healthPort != 0 {
		// could implement your checker , type Checker func(req *http.Request) error
		if err := mgr.AddHealthzCheck("/healthy/liveness", healthz.Ping); err != nil {
			logger.Sugar().Fatalf("failed to AddHealthzCheck, reason=%v", err)
		}
		if err := mgr.AddReadyzCheck("/healthy/readiness", healthz.Ping); err != nil {
			logger.Sugar().Fatalf("failed to AddReadyzCheck, reason=%v", err)
		}
		// add other route
		// mgr.GetWebhookServer().Register("/route", XXXX)
	}

	if e := k8sObjManager.Initk8sObjManager(mgr.GetClient()); e != nil {
		logger.Sugar().Fatalf("failed to Initk8sObjManager, error=%v", e)

	}

	var fm fileManager.FileManager
	var e error

	// TODO: runControllerAggregateReportOnce need agents' IPs
	if types.ControllerConfig.EnableAggregateAgentReport {
		// fileManager takes charge of writing and removing local report
		gcInterval := time.Duration(types.ControllerConfig.CleanAgedReportInMinute) * time.Minute
		logger.Sugar().Infof("save report to %v, clean interval %v", types.ControllerConfig.DirPathControllerReport, gcInterval.String())
		fm, e = fileManager.NewManager(logger.Named("fileManager"), types.ControllerConfig.DirPathControllerReport, gcInterval)
		if e != nil {
			logger.Sugar().Fatalf("failed to new fileManager , reason=%v", e)
		}
	}

	runtimeDB := make([]scheduler.DB, 0, len(s.chainingPlugins))
	ctx, cancelFunc := context.WithCancel(context.TODO())
	for name, plugin := range s.chainingPlugins {
		// setup reconcile
		logger.Sugar().Infof("run controller for plugin %v", name)

		uniqueMatchLabelKey := types.ControllerConfig.Configmap.KdoctorAgent.UniqueMatchLabelKey
		if len(uniqueMatchLabelKey) == 0 {
			logger.Sugar().Debugf("there's no uniqueMatchLabelKey in the configmap, try to use the default '%s'", scheduler.UniqueMatchLabelKey)
			uniqueMatchLabelKey = scheduler.UniqueMatchLabelKey
		}
		tracker := scheduler.NewTracker(mgr.GetClient(), mgr.GetAPIReader(), scheduler.TrackerConfig{
			ItemChannelBuffer:     int(types.ControllerConfig.ResourceTrackerChannelBuffer),
			MaxDatabaseCap:        int(types.ControllerConfig.ResourceTrackerMaxDatabaseCap),
			ExecutorWorkers:       int(types.ControllerConfig.ResourceTrackerExecutorWorkers),
			SignalTimeOutDuration: time.Duration(types.ControllerConfig.ResourceTrackerSignalTimeoutSeconds) * time.Second,
			TraceGapDuration:      time.Duration(types.ControllerConfig.ResourceTrackerTraceGapSeconds) * time.Second,
		}, logger.Named(name+"Tracker"))
		runtimeDB = append(runtimeDB, tracker.DB)
		k := &pluginControllerReconciler{
			logger:                     logger.Named(name + "Reconciler"),
			plugin:                     plugin,
			client:                     mgr.GetClient(),
			apiReader:                  mgr.GetAPIReader(),
			crdKind:                    name,
			fm:                         fm,
			crdKindName:                name,
			runtimeUniqueMatchLabelKey: uniqueMatchLabelKey,
			tracker:                    tracker,
		}
		k.tracker.Start(ctx)
		if e := k.SetupWithManager(mgr); e != nil {
			s.logger.Sugar().Fatalf("failed to builder reconcile for plugin %v, error=%v", name, e)
		}
		// setup webhook
		t := &pluginWebhookhander{
			logger:  logger.Named(name + "Webhook"),
			plugin:  plugin,
			client:  mgr.GetClient(),
			crdKind: name,
		}
		if e := t.SetupWebhook(mgr); e != nil {
			s.logger.Sugar().Fatalf("failed to builder webhook for plugin %v, error=%v", name, e)
		}
	}

	if types.ControllerConfig.EnableAggregateAgentReport {
		// reportManager takes charge of sync reports from remote agents
		interval := time.Duration(types.ControllerConfig.CollectAgentReportIntervalInSecond) * time.Second
		logger.Sugar().Infof("run report Sync manager, save to %v, collectInterval %v ", types.ControllerConfig.DirPathControllerReport, interval)
		reportManager.InitReportManager(logger.Named("reportSyncManager"), types.ControllerConfig.DirPathControllerReport, interval, runtimeDB)
	}

	go func() {
		msg := "reconcile of plugin down"
		if e := mgr.Start(ctrl.SetupSignalHandler()); e != nil {
			msg += fmt.Sprintf(", error=%v", e)
			cancelFunc()
		}
		s.logger.Error(msg)
		time.Sleep(5 * time.Second)
	}()

}
