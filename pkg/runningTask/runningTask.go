// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package runningTask

import (
	"github.com/kdoctor-io/kdoctor/pkg/k8s/apis/system/v1beta1"
	"github.com/kdoctor-io/kdoctor/pkg/lock"
	"github.com/kdoctor-io/kdoctor/pkg/types"
)

type Task struct {
	Kind string
	Qps  int
	Name string
}

type RunningTask struct {
	task map[string]Task
	lock.RWMutex
}

func InitRunningTask() *RunningTask {
	return &RunningTask{
		task: make(map[string]Task, 0),
	}
}

func (rt *RunningTask) SetTask(task Task) {
	rt.Lock()
	defer rt.Unlock()
	rt.task[task.Name] = task
}

func (rt *RunningTask) DeleteTask(taskName string) {
	rt.Lock()
	defer rt.Unlock()
	delete(rt.task, taskName)
}

func (rt *RunningTask) QpsStats() v1beta1.TotalRunningLoad {
	rt.Lock()
	defer rt.Unlock()

	var appHealthQps int
	var netDNSQps int
	var netReachQps int

	for _, v := range rt.task {
		switch v.Kind {
		case types.KindNameAppHttpHealthy:
			appHealthQps += v.Qps
		case types.KindNameNetReach:
			netReachQps += v.Qps
		case types.KindNameNetdns:
			netDNSQps += v.Qps
		}
	}

	return v1beta1.TotalRunningLoad{
		AppHttpHealthyQPS: int64(appHealthQps),
		NetDnsQPS:         int64(netDNSQps),
		NetReachQPS:       int64(netReachQps),
	}
}
