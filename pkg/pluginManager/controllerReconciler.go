// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package pluginManager

import (
	"context"
	"reflect"
	"time"

	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/kdoctor-io/kdoctor/pkg/fileManager"
	crd "github.com/kdoctor-io/kdoctor/pkg/k8s/apis/kdoctor.io/v1beta1"
	plugintypes "github.com/kdoctor-io/kdoctor/pkg/pluginManager/types"
	"github.com/kdoctor-io/kdoctor/pkg/scheduler"
)

type pluginControllerReconciler struct {
	client      client.Client
	apiReader   client.Reader
	plugin      plugintypes.ChainingPlugin
	logger      *zap.Logger
	crdKind     string
	fm          fileManager.FileManager
	crdKindName string

	runtimeUniqueMatchLabelKey string
	tracker                    *scheduler.Tracker
}

// contorller reconcile
// (1) chedule all task time
// (2) update stauts result
// (3) collect report from agent
func (s *pluginControllerReconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {

	// ------ add crd ------
	switch s.crdKind {
	case KindNameNetReach:
		// ------ add crd ------
		instance := crd.NetReach{}

		if err := s.client.Get(ctx, req.NamespacedName, &instance); err != nil {
			s.logger.Sugar().Errorf("unable to fetch obj , error=%v", err)
			// since we have OwnerReference for task corresponding runtime and service, we could just delete the tracker DB record directly
			if errors.IsNotFound(err) && instance.Status.Resource != nil {
				s.tracker.DB.Delete(scheduler.BuildItem(*instance.Status.Resource, KindNameNetReach, instance.Name, nil))
			}
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}
		logger := s.logger.With(zap.String(instance.Kind, instance.Name))
		logger.Sugar().Debugf("reconcile handle %v", instance)

		if instance.DeletionTimestamp != nil {
			s.logger.Sugar().Debugf("ignore deleting task %v", req)
			return ctrl.Result{}, nil
		}

		newStatus, err := s.TaskResourceReconcile(ctx, KindNameNetReach, &instance, instance.Spec.AgentSpec, instance.Status.DeepCopy(), logger)
		if nil != err {
			logger.Sugar().Errorf(err.Error())
			return ctrl.Result{}, err
		}
		if !reflect.DeepEqual(newStatus, instance.Status.DeepCopy()) {
			instance.Status = *newStatus
			logger.Sugar().Infof("try to update %s/%s status with resource %v", KindNameNetReach, instance.Name, newStatus.Resource)
			err := s.client.Status().Update(ctx, &instance)
			if nil != err {
				logger.Sugar().Errorf("failed to update %s/%s status with resource %v, error: %v", KindNameNetReach, instance.Name, newStatus.Resource, err)
				return reconcile.Result{}, err
			}
		}

		// runtime creating status means the agent is not ready, so we don't need to initial the task right now.
		// the tracker DB will update the status asynchronously, and we would receive the task event after it updated.
		if instance.Status.Resource.RuntimeStatus == crd.RuntimeCreating {
			return ctrl.Result{}, nil
		}

		// the task corresponding agent pods have this unique label
		runtimePodMatchLabels := client.MatchingLabels{
			s.runtimeUniqueMatchLabelKey: scheduler.UniqueMatchLabelValue(KindNameNetReach, instance.Name),
		}

		oldStatus := instance.Status.DeepCopy()
		taskName := instance.Kind + "." + instance.Name
		if result, newStatus, err := s.UpdateStatus(logger, ctx, oldStatus, instance.Spec.Schedule.DeepCopy(), runtimePodMatchLabels, taskName); err != nil {
			// requeue
			logger.Sugar().Errorf("failed to UpdateStatus, will retry it, error=%v", err)
			return ctrl.Result{}, err
		} else {
			if newStatus != nil {
				if !reflect.DeepEqual(newStatus, oldStatus) {
					instance.Status = *newStatus
					if err := s.client.Status().Update(ctx, &instance); err != nil {
						// requeue
						logger.Sugar().Errorf("failed to update status, will retry it, error=%v", err)
						return ctrl.Result{}, err
					}
					logger.Sugar().Debugf("succeeded update status, newStatus=%+v", newStatus)
				}

				// update tracker database
				if newStatus.FinishTime != nil {
					deletionTime := newStatus.FinishTime.DeepCopy()
					if instance.Spec.AgentSpec.TerminationGracePeriodMinutes != nil {
						newTime := metav1.NewTime(deletionTime.Add(time.Duration(*instance.Spec.AgentSpec.TerminationGracePeriodMinutes) * time.Minute))
						deletionTime = newTime.DeepCopy()
					}
					logger.Sugar().Debugf("task finish time '%s' and runtime deletion time '%s'", newStatus.FinishTime, deletionTime)
					// record the task resource to the tracker DB, and the tracker will update the task subresource resource status asynchronously
					err := s.tracker.DB.Apply(scheduler.BuildItem(*instance.Status.Resource, KindNameNetReach, instance.Name, deletionTime))
					if nil != err {
						logger.Error(err.Error())
						return ctrl.Result{}, err
					}
				}
			}

			if result != nil {
				return *result, nil
			}
		}

	case KindNameAppHttpHealthy:
		// ------ add crd ------
		instance := crd.AppHttpHealthy{}

		if err := s.client.Get(ctx, req.NamespacedName, &instance); err != nil {
			s.logger.Sugar().Errorf("unable to fetch obj , error=%v", err)
			// since we have OwnerReference for task corresponding runtime and service, we could just delete the tracker DB record directly
			if errors.IsNotFound(err) && instance.DeletionTimestamp != nil {
				s.tracker.DB.Delete(scheduler.BuildItem(*instance.Status.Resource, KindNameAppHttpHealthy, instance.Name, nil))

			}
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}
		logger := s.logger.With(zap.String(instance.Kind, instance.Name))
		logger.Sugar().Debugf("reconcile handle %v", instance)

		if instance.DeletionTimestamp != nil {
			s.logger.Sugar().Debugf("ignore deleting task %v", req)
			return ctrl.Result{}, nil
		}

		newStatus, err := s.TaskResourceReconcile(ctx, KindNameAppHttpHealthy, &instance, instance.Spec.AgentSpec, instance.Status.DeepCopy(), logger)
		if nil != err {
			logger.Sugar().Errorf(err.Error())
			return ctrl.Result{}, err
		}
		if !reflect.DeepEqual(newStatus, instance.Status.DeepCopy()) {
			instance.Status = *newStatus
			logger.Sugar().Infof("try to update %s/%s status with resource %v", KindNameAppHttpHealthy, instance.Name, newStatus.Resource)
			err := s.client.Status().Update(ctx, &instance)
			if nil != err {
				logger.Sugar().Errorf("failed to update %s/%s status with resource %v, error: %v", KindNameAppHttpHealthy, instance.Name, newStatus.Resource, err)
				return reconcile.Result{}, err
			}
		}

		// runtime creating status means the agent is not ready, so we don't need to initial the task right now.
		// the tracker DB will update the status asynchronously, and we would receive the task event after it updated.
		if instance.Status.Resource.RuntimeStatus == crd.RuntimeCreating {
			return ctrl.Result{}, nil
		}

		// the task corresponding agent pods have this unique label
		runtimePodMatchLabels := client.MatchingLabels{
			s.runtimeUniqueMatchLabelKey: scheduler.UniqueMatchLabelValue(KindNameAppHttpHealthy, instance.Name),
		}

		oldStatus := instance.Status.DeepCopy()
		taskName := instance.Kind + "." + instance.Name
		if result, newStatus, err := s.UpdateStatus(logger, ctx, oldStatus, instance.Spec.Schedule.DeepCopy(), runtimePodMatchLabels, taskName); err != nil {
			// requeue
			logger.Sugar().Errorf("failed to UpdateStatus, will retry it, error=%v", err)
			return ctrl.Result{}, err
		} else {
			if newStatus != nil {
				if !reflect.DeepEqual(newStatus, oldStatus) {
					instance.Status = *newStatus
					if err := s.client.Status().Update(ctx, &instance); err != nil {
						// requeue
						logger.Sugar().Errorf("failed to update status, will retry it, error=%v", err)
						return ctrl.Result{}, err
					}
					logger.Sugar().Debugf("succeeded update status, newStatus=%+v", newStatus)
				}

				// update tracker database
				if newStatus.FinishTime != nil {
					deletionTime := newStatus.FinishTime.DeepCopy()
					if instance.Spec.AgentSpec.TerminationGracePeriodMinutes != nil {
						newTime := metav1.NewTime(deletionTime.Add(time.Duration(*instance.Spec.AgentSpec.TerminationGracePeriodMinutes) * time.Minute))
						deletionTime = newTime.DeepCopy()
					}
					logger.Sugar().Debugf("task finish time '%s' and runtime deletion time '%s'", newStatus.FinishTime, deletionTime)
					err := s.tracker.DB.Apply(scheduler.BuildItem(*instance.Status.Resource, KindNameAppHttpHealthy, instance.Name, deletionTime))
					if nil != err {
						logger.Error(err.Error())
						return ctrl.Result{}, err
					}
				}
			}

			if result != nil {
				return *result, nil
			}
		}

	case KindNameNetdns:
		// ------ add crd ------
		instance := crd.Netdns{}

		if err := s.client.Get(ctx, req.NamespacedName, &instance); err != nil {
			s.logger.Sugar().Errorf("unable to fetch obj , error=%v", err)
			// since we have OwnerReference for task corresponding runtime and service, we could just delete the tracker DB record directly
			if errors.IsNotFound(err) && instance.DeletionTimestamp != nil {
				s.tracker.DB.Delete(scheduler.BuildItem(*instance.Status.Resource, KindNameNetdns, instance.Name, nil))
			}
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}
		logger := s.logger.With(zap.String(instance.Kind, instance.Name))
		logger.Sugar().Debugf("reconcile handle %v", instance)

		if instance.DeletionTimestamp != nil {
			s.logger.Sugar().Debugf("ignore deleting task %v", req)
			return ctrl.Result{}, nil
		}

		newStatus, err := s.TaskResourceReconcile(ctx, KindNameNetdns, &instance, instance.Spec.AgentSpec, instance.Status.DeepCopy(), logger)
		if nil != err {
			logger.Sugar().Errorf(err.Error())
			return ctrl.Result{}, err
		}
		if !reflect.DeepEqual(newStatus, instance.Status.DeepCopy()) {
			instance.Status = *newStatus
			logger.Sugar().Infof("try to update %s/%s status with resource %v", KindNameNetdns, instance.Name, newStatus.Resource)
			err := s.client.Status().Update(ctx, &instance)
			if nil != err {
				logger.Sugar().Errorf("failed to update %s/%s status with resource %v, error: %v", KindNameNetdns, instance.Name, newStatus.Resource, err)
				return reconcile.Result{}, err
			}
		}

		// runtime creating status means the agent is not ready, so we don't need to initial the task right now.
		// the tracker DB will update the status asynchronously, and we would receive the task event after it updated.
		if instance.Status.Resource.RuntimeStatus == crd.RuntimeCreating {
			return ctrl.Result{}, nil
		}

		// the task corresponding agent pods have this unique label
		runtimePodMatchLabels := client.MatchingLabels{
			s.runtimeUniqueMatchLabelKey: scheduler.TaskRuntimeName(KindNameNetdns, instance.Name),
		}

		oldStatus := instance.Status.DeepCopy()
		taskName := instance.Kind + "." + instance.Name
		if result, newStatus, err := s.UpdateStatus(logger, ctx, oldStatus, instance.Spec.Schedule.DeepCopy(), runtimePodMatchLabels, taskName); err != nil {
			// requeue
			logger.Sugar().Errorf("failed to UpdateStatus, will retry it, error=%v", err)
			return ctrl.Result{}, err
		} else {
			if newStatus != nil {
				if !reflect.DeepEqual(newStatus, oldStatus) {
					instance.Status = *newStatus
					if err := s.client.Status().Update(ctx, &instance); err != nil {
						// requeue
						logger.Sugar().Errorf("failed to update status, will retry it, error=%v", err)
						return ctrl.Result{}, err
					}
					logger.Sugar().Debugf("succeeded update status, newStatus=%+v", newStatus)
				}

				// update tracker database
				if newStatus.FinishTime != nil {
					deletionTime := newStatus.FinishTime.DeepCopy()
					if instance.Spec.AgentSpec.TerminationGracePeriodMinutes != nil {
						newTime := metav1.NewTime(deletionTime.Add(time.Duration(*instance.Spec.AgentSpec.TerminationGracePeriodMinutes) * time.Minute))
						deletionTime = newTime.DeepCopy()
					}
					logger.Sugar().Debugf("task finish time '%s' and runtime deletion time '%s'", newStatus.FinishTime, deletionTime)
					err := s.tracker.DB.Apply(scheduler.BuildItem(*instance.Status.Resource, KindNameNetdns, instance.Name, deletionTime))
					if nil != err {
						logger.Error(err.Error())
						return ctrl.Result{}, err
					}
				}
			}
			if result != nil {
				return *result, nil
			}
		}

	default:
		s.logger.Sugar().Fatalf("unknown crd type , support kind=%v, detail=%+v", s.crdKind, req)
	}
	// forget this
	return ctrl.Result{}, nil

	// return s.plugin.ControllerReconcile(s.logger, s.client, ctx, req)
}

var _ reconcile.Reconciler = &pluginControllerReconciler{}

func (s *pluginControllerReconciler) SetupWithManager(mgr ctrl.Manager) error {

	return ctrl.NewControllerManagedBy(mgr).For(s.plugin.GetApiType()).Complete(s)
}
