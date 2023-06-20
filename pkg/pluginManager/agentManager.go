// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package pluginManager

import (
	"fmt"
	"github.com/kdoctor-io/kdoctor/pkg/fileManager"
	k8sObjManager "github.com/kdoctor-io/kdoctor/pkg/k8ObjManager"
	crd "github.com/kdoctor-io/kdoctor/pkg/k8s/apis/kdoctor.io/v1beta1"
	"github.com/kdoctor-io/kdoctor/pkg/taskStatusManager"
	"github.com/kdoctor-io/kdoctor/pkg/types"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/utils/strings/slices"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

func (s *pluginManager) RunAgentController() {
	logger := s.logger
	logger.Sugar().Infof("setup agent reconcile")

	scheme := runtime.NewScheme()
	if e := clientgoscheme.AddToScheme(scheme); e != nil {
		logger.Sugar().Fatalf("failed to add k8s scheme, reason=%v", e)
	}
	if e := crd.AddToScheme(scheme); e != nil {
		logger.Sugar().Fatalf("failed to add scheme for plugins, reason=%v", e)
	}

	var fm fileManager.FileManager
	var e error
	if types.AgentConfig.EnableAggregateAgentReport {
		gcInterval := time.Duration(types.AgentConfig.CleanAgedReportInMinute) * time.Minute
		logger.Sugar().Infof("save report to %v, clean interval %v", types.AgentConfig.DirPathAgentReport, gcInterval.String())
		fm, e = fileManager.NewManager(logger.Named("fileManager"), types.AgentConfig.DirPathAgentReport, gcInterval)
		if e != nil {
			logger.Sugar().Fatalf("failed to new fileManager , reason=%v", e)
		}
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
		logger.Sugar().Fatalf("failed to NewManager, reason=%v", err)
	}

	if len(types.AgentConfig.LocalNodeName) == 0 {
		logger.Sugar().Fatalf("local node name is empty")
	}

	if e := k8sObjManager.Initk8sObjManager(mgr.GetClient()); e != nil {
		logger.Sugar().Fatalf("failed to Initk8sObjManager, error=%v", e)
	}

	if !slices.Contains(types.TaskKinds, types.AgentConfig.TaskKind) {
		logger.Sugar().Fatalf("unsupport TaskKind %s in %v", types.AgentConfig.TaskKind, types.TaskKinds)
	}

	for name, plugin := range s.chainingPlugins {
		if name != types.AgentConfig.TaskKind {
			continue
		}

		logger.Sugar().Infof("run controller for plugin %v", name)
		k := &pluginAgentReconciler{
			logger:        logger.Named(name + "Reconciler"),
			plugin:        plugin,
			client:        mgr.GetClient(),
			crdKind:       name,
			taskRoundData: taskStatusManager.NewTaskStatus(),
			localNodeName: types.AgentConfig.LocalNodeName,
			fm:            fm,
		}
		if e := k.SetupWithManager(mgr); e != nil {
			s.logger.Sugar().Fatalf("failed to builder reconcile for plugin %v, error=%v", name, e)
		}
	}

	// before mgr.Start, it should not use mgr.GetClient() to get api obj, because "the controller cache is not started, can not read objects"
	go func() {
		msg := "reconcile of plugin down"
		if e := mgr.Start(ctrl.SetupSignalHandler()); e != nil {
			msg += fmt.Sprintf(", error=%v", e)
		}
		s.logger.Error(msg)
		time.Sleep(5 * time.Second)
	}()

}
