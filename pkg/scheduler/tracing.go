// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package scheduler

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8types "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	crd "github.com/kdoctor-io/kdoctor/pkg/k8s/apis/kdoctor.io/v1beta1"
	scheduleruntime "github.com/kdoctor-io/kdoctor/pkg/scheduler/runtime"
	"github.com/kdoctor-io/kdoctor/pkg/types"
)

type TrackerConfig struct {
	ItemChannelBuffer int
	MaxDatabaseCap    int
	ExecutorWorkers   int

	SignalTimeOutDuration time.Duration // default: 5s
	TraceGapDuration      time.Duration // default: 10s
}

type Tracker struct {
	client     client.Client
	apiReader  client.Reader
	DB         DB
	itemSignal chan Item
	log        *zap.Logger

	TrackerConfig
}

func NewTracker(c client.Client, apiReader client.Reader, config TrackerConfig, log *zap.Logger) *Tracker {
	tracker := &Tracker{
		client:        c,
		apiReader:     apiReader,
		DB:            NewDB(config.MaxDatabaseCap, log.Named("Database")),
		itemSignal:    make(chan Item, config.ItemChannelBuffer),
		log:           log,
		TrackerConfig: config,
	}

	return tracker
}

func (t *Tracker) Start(ctx context.Context) {
	// trace db
	go t.trace(ctx)

	for i := 1; i <= t.ExecutorWorkers; i++ {
		go t.executor(ctx, i)
	}

	t.log.Info("Tracker is running")
}

func (t *Tracker) trace(ctx context.Context) {
	t.log.Info("starting tracing db")

	for {
		select {
		case <-ctx.Done():
			t.log.Warn("received ctx done, stop tracing")
			return
		default:
			items := t.DB.List()
			for _, item := range items {
				t.signaling(item)
			}

			time.Sleep(t.TraceGapDuration)
		}
	}
}

func (t *Tracker) signaling(item Item) {
	select {
	case t.itemSignal <- item:
		t.log.Sugar().Debugf("sending signal to handle Task %s", item.Task)

	case <-time.After(t.SignalTimeOutDuration):
		t.log.Sugar().Warnf("failed to send signal, itemSignal length %d, item %v will be dropped", len(t.itemSignal), item)
	}
}

func (t *Tracker) executor(ctx context.Context, workerIndex int) {
	innerLog := t.log.With(zap.Any("Executor_Index", workerIndex))
	innerLog.Info("start running executor")

	for {
		select {
		case item := <-t.itemSignal:
			runtime, err := scheduleruntime.FindRuntime(t.client, t.apiReader, item.RuntimeKey.RuntimeKind, item.RuntimeKey.RuntimeName, t.log.Named("Runtime"))
			if nil != err {
				innerLog.Error(err.Error())
				continue
			}

			// deletion
			if item.RuntimeDeletionTime != nil {
				if metav1.Now().After(item.RuntimeDeletionTime.Time) {
					// 1. delete runtime
					innerLog.Sugar().Debugf("resource %v is already out of deletiontime, try to delete it", item.RuntimeKey)
					err := runtime.Delete(ctx)
					if client.IgnoreNotFound(err) != nil {
						innerLog.Sugar().Errorf("failed to delete resource %v, error: %v", item.RuntimeKey, err)
						continue
					}

					// 2. update status
					innerLog.Sugar().Infof("try to update task %s resource status from %s to %s", item.Task, item.RuntimeStatus, crd.RuntimeDeleted)
					err = t.updateRuntimeStatus(ctx, item, crd.RuntimeDeleted)
					if client.IgnoreNotFound(err) != nil {
						innerLog.Error(err.Error())
						continue
					}

					// 3. clean up db
					t.DB.Delete(item)
					continue
				} else {
					innerLog.Sugar().Debugf("resource %v isn't out of date, skip it", item.RuntimeKey)
				}
			}

			// update created
			if item.RuntimeStatus == crd.RuntimeCreating && runtime.IsReady(ctx) {
				innerLog.Sugar().Infof("try to update task %s resource status from %s to %s", item.Task, item.RuntimeStatus, crd.RuntimeCreated)
				err := t.updateRuntimeStatus(ctx, item, crd.RuntimeCreated)
				if nil != err {
					innerLog.Error(err.Error())
					continue
				}

				// clean up db
				t.DB.Delete(item)
			}

		case <-ctx.Done():
			innerLog.Warn("received ctx done, stop running executor")
			return
		}
	}
}

// TODO: optimize here with interface?
func (t *Tracker) updateRuntimeStatus(ctx context.Context, item Item, status string) error {
	resource := &crd.TaskResource{
		RuntimeName:   item.RuntimeName,
		RuntimeType:   item.RuntimeKind,
		ServiceNameV4: item.ServiceNameV4,
		ServiceNameV6: item.ServiceNameV6,
		RuntimeStatus: status,
	}

	for taskName, taskKind := range item.Task {
		switch taskKind {
		case types.KindNameNetReach:
			instance := crd.NetReach{}
			err := t.apiReader.Get(ctx, k8types.NamespacedName{Name: taskName}, &instance)
			if nil != err {
				return err
			}

			// check the resource whether is already equal
			if reflect.DeepEqual(instance.Status.Resource, resource) {
				t.log.Sugar().Debugf("task %v resource already updatede, skip it", item.RuntimeKey)
				return nil
			}

			t.log.Sugar().Debugf("task %v old resource is %v, the new resource is %v", item.RuntimeKey, *instance.Status.Resource, *resource)
			instance.Status.Resource = resource
			err = t.client.Status().Update(ctx, &instance)
			if nil != err {
				return err
			}

		case types.KindNameAppHttpHealthy:
			instance := crd.AppHttpHealthy{}
			err := t.apiReader.Get(ctx, k8types.NamespacedName{Name: taskName}, &instance)
			if nil != err {
				return err
			}

			// check the resource whether is already equal
			if reflect.DeepEqual(instance.Status.Resource, resource) {
				t.log.Sugar().Debugf("task %v resource already updatede, skip it", item.RuntimeKey)
				return nil
			}

			t.log.Sugar().Debugf("task %v old resource is %v, the new resource is %v", item.RuntimeKey, *instance.Status.Resource, *resource)
			instance.Status.Resource = resource
			err = t.client.Status().Update(ctx, &instance)
			if nil != err {
				return err
			}

		case types.KindNameNetdns:
			instance := crd.Netdns{}
			err := t.apiReader.Get(ctx, k8types.NamespacedName{Name: taskName}, &instance)
			if nil != err {
				return err
			}

			// check the resource whether is already equal
			if reflect.DeepEqual(instance.Status.Resource, resource) {
				t.log.Sugar().Debugf("task %v resource already updatede, skip it", item.RuntimeKey)
				return nil
			}

			t.log.Sugar().Debugf("task %v old resource is %v, the new resource is %v", item.RuntimeKey, *instance.Status.Resource, *resource)
			instance.Status.Resource = resource
			err = t.client.Status().Update(ctx, &instance)
			if nil != err {
				return err
			}

		default:
			return fmt.Errorf("unsupported task '%s/%s'", taskKind, taskName)
		}
	}

	return nil
}
