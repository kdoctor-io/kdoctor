// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package scheduler

import (
	"fmt"
	"reflect"

	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	crd "github.com/kdoctor-io/kdoctor/pkg/k8s/apis/kdoctor.io/v1beta1"
	"github.com/kdoctor-io/kdoctor/pkg/lock"
)

const DefaultMaxDBCap = 10000

type DB interface {
	Apply(item Item) error
	List() []Item
	Delete(item Item)
}

func NewDB(maxCap int, log *zap.Logger) DB {
	if maxCap <= 0 {
		maxCap = DefaultMaxDBCap
	}

	return &Database{
		cache:  make(map[RuntimeKey]Item),
		maxCap: maxCap,
		log:    log,
	}
}

type Database struct {
	lock.RWMutex
	cache  map[RuntimeKey]Item
	maxCap int
	log    *zap.Logger
}

type RuntimeKey struct {
	RuntimeKind string
	RuntimeName string
}

type Item struct {
	RuntimeKey

	RuntimeStatus       string
	RuntimeDeletionTime *metav1.Time
	ServiceNameV4       *string
	ServiceNameV6       *string

	TaskKind string
	TaskName string
}

func BuildItem(resource crd.TaskResource, taskKind, taskName string, deletionTime *metav1.Time) Item {
	item := Item{
		RuntimeKey: RuntimeKey{
			RuntimeKind: resource.RuntimeType,
			RuntimeName: resource.RuntimeName,
		},
		RuntimeStatus:       resource.RuntimeStatus,
		RuntimeDeletionTime: deletionTime,
		ServiceNameV4:       resource.ServiceNameV4,
		ServiceNameV6:       resource.ServiceNameV6,
		TaskKind:            taskKind,
		TaskName:            taskName,
	}

	return item
}

func (d *Database) Apply(item Item) error {
	d.Lock()

	old, ok := d.cache[item.RuntimeKey]
	if !ok {
		if len(d.cache) == d.maxCap {
			d.Unlock()
			return fmt.Errorf("database is out of capacity, discard item %v", item)
		}

		d.cache[item.RuntimeKey] = item
		d.Unlock()
		d.log.Sugar().Debugf("create item %v", item)
		return nil
	} else {
		if !reflect.DeepEqual(old, item) {
			d.cache[item.RuntimeKey] = item
			d.Unlock()
			d.log.Sugar().Debugf("item %v has changed, the old one is %v, and the new one is %v",
				item.RuntimeKey, old, item)
			return nil
		}
	}

	d.Unlock()
	return nil
}

func (d *Database) List() []Item {
	d.RLock()
	defer d.RUnlock()

	items := make([]Item, 0, len(d.cache))
	for k := range d.cache {
		items = append(items, d.cache[k])
	}

	return items
}

func (d *Database) Delete(item Item) {
	d.Lock()

	_, ok := d.cache[item.RuntimeKey]
	if !ok {
		d.Unlock()
		d.log.Sugar().Debugf("item %v already deleted", item.RuntimeKey)
	}

	delete(d.cache, item.RuntimeKey)
	d.Unlock()
	d.log.Sugar().Debugf("delete item %v successfully", item.RuntimeKey)
}
