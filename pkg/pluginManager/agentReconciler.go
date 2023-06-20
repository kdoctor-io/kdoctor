// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package pluginManager

import (
	"context"
	"github.com/kdoctor-io/kdoctor/pkg/fileManager"
	crd "github.com/kdoctor-io/kdoctor/pkg/k8s/apis/kdoctor.io/v1beta1"
	plugintypes "github.com/kdoctor-io/kdoctor/pkg/pluginManager/types"
	"github.com/kdoctor-io/kdoctor/pkg/taskStatusManager"
	"github.com/kdoctor-io/kdoctor/pkg/types"
	"go.uber.org/zap"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type pluginAgentReconciler struct {
	client        client.Client
	plugin        plugintypes.ChainingPlugin
	logger        *zap.Logger
	crdKind       string
	localNodeName string
	taskRoundData taskStatusManager.TaskStatus
	fm            fileManager.FileManager
}

var _ reconcile.Reconciler = &pluginAgentReconciler{}

func (s *pluginAgentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).For(s.plugin.GetApiType()).Complete(s)
}

// https://github.com/kubernetes-sigs/controller-runtime/blob/master/pkg/internal/controller/controller.go#L320
// when err!=nil , c.Queue.AddRateLimited(req) and log error
// when err==nil && result.Requeue, just c.Queue.AddRateLimited(req)
// when err==nil && result.RequeueAfter > 0 , c.Queue.Forget(obj) and c.Queue.AddAfter(req, result.RequeueAfter)
// or else, c.Queue.Forget(obj)
func (s *pluginAgentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// filter other tasks
	if req.NamespacedName.Name != types.AgentConfig.TaskName {
		s.logger.With(zap.String(types.AgentConfig.TaskKind, types.AgentConfig.TaskName)).
			Sugar().Debugf("ignore Task %s", req.NamespacedName.Name)
		return ctrl.Result{}, nil
	}

	// ------ add crd ------
	switch s.crdKind {
	case KindNameNetReach:
		instance := crd.NetReach{}
		if err := s.client.Get(ctx, req.NamespacedName, &instance); err != nil {
			s.logger.Sugar().Errorf("unable to fetch obj , error=%v", err)
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}

		logger := s.logger.With(zap.String(instance.Kind, instance.Name))
		logger.Sugar().Debugf("reconcile handle %v", instance)
		if instance.DeletionTimestamp != nil {
			s.logger.Sugar().Debugf("ignore deleting task %v", req)
			return ctrl.Result{}, nil
		}

		oldStatus := instance.Status.DeepCopy()
		taskName := instance.Kind + "." + instance.Name
		if result, newStatus, err := s.HandleAgentTaskRound(logger, ctx, oldStatus, instance.Spec.Schedule.DeepCopy(), &instance, taskName, instance.Spec.DeepCopy()); err != nil {
			// requeue
			logger.Sugar().Errorf("failed to HandleAgentTaskRound, will retry it, error=%v", err)
			return ctrl.Result{}, err

		} else {
			if newStatus != nil && !reflect.DeepEqual(newStatus, oldStatus) {
				instance.Status = *newStatus
				if err := s.client.Status().Update(ctx, &instance); err != nil {
					// requeue
					logger.Sugar().Errorf("failed to update status, will retry it, error=%v", err)
					return ctrl.Result{}, err
				}
				logger.Sugar().Debugf("succeeded update status, newStatus=%+v", newStatus)
			}

			if result != nil {
				return *result, nil
			}
		}

	case KindNameAppHttpHealthy:
		instance := crd.AppHttpHealthy{}
		if err := s.client.Get(ctx, req.NamespacedName, &instance); err != nil {
			s.logger.Sugar().Errorf("unable to fetch obj , error=%v", err)
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}

		logger := s.logger.With(zap.String(instance.Kind, instance.Name))
		logger.Sugar().Debugf("reconcile handle %v", instance)
		if instance.DeletionTimestamp != nil {
			s.logger.Sugar().Debugf("ignore deleting task %v", req)
			return ctrl.Result{}, nil
		}

		oldStatus := instance.Status.DeepCopy()
		taskName := instance.Kind + "." + instance.Name
		if result, newStatus, err := s.HandleAgentTaskRound(logger, ctx, oldStatus, instance.Spec.Schedule.DeepCopy(), &instance, taskName, instance.Spec.DeepCopy()); err != nil {
			// requeue
			logger.Sugar().Errorf("failed to HandleAgentTaskRound, will retry it, error=%v", err)
			return ctrl.Result{}, err

		} else {
			if newStatus != nil && !reflect.DeepEqual(newStatus, oldStatus) {
				instance.Status = *newStatus
				if err := s.client.Status().Update(ctx, &instance); err != nil {
					// requeue
					logger.Sugar().Errorf("failed to update status, will retry it, error=%v", err)
					return ctrl.Result{}, err
				}
				logger.Sugar().Debugf("succeeded update status, newStatus=%+v", newStatus)
			}

			if result != nil {
				return *result, nil
			}
		}
	case KindNameNetdns:
		instance := crd.Netdns{}
		if err := s.client.Get(ctx, req.NamespacedName, &instance); err != nil {
			s.logger.Sugar().Errorf("unable to fetch obj , error=%v", err)
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}
		logger := s.logger.With(zap.String(instance.Kind, instance.Name))
		logger.Sugar().Debugf("reconcile handle %v", instance)
		if instance.DeletionTimestamp != nil {
			s.logger.Sugar().Debugf("ignore deleting task %v", req)
			return ctrl.Result{}, nil
		}

		oldStatus := instance.Status.DeepCopy()
		taskName := instance.Kind + "." + instance.Name
		if result, newStatus, err := s.HandleAgentTaskRound(logger, ctx, oldStatus, instance.Spec.Schedule.DeepCopy(), &instance, taskName, instance.Spec.DeepCopy()); err != nil {
			// requeue
			logger.Sugar().Errorf("failed to HandleAgentTaskRound, will retry it, error=%v", err)
			return ctrl.Result{}, err

		} else {
			if newStatus != nil && !reflect.DeepEqual(newStatus, oldStatus) {
				instance.Status = *newStatus
				if err := s.client.Status().Update(ctx, &instance); err != nil {
					// requeue
					logger.Sugar().Errorf("failed to update status, will retry it, error=%v", err)
					return ctrl.Result{}, err
				}
				logger.Sugar().Debugf("succeeded update status, newStatus=%+v", newStatus)
			}

			if result != nil {
				return *result, nil
			}
		}

	default:
		s.logger.Sugar().Fatalf("unknown crd type , support kind=%v, detail=%+v", s.crdKind, req)
	}

	return ctrl.Result{}, nil
}
