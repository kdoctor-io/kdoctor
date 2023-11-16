// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package resource

import (
	"context"
	"fmt"
	"github.com/kdoctor-io/kdoctor/pkg/k8s/apis/system/v1beta1"
	"runtime"
	"time"

	"github.com/kdoctor-io/kdoctor/pkg/lock"
	"github.com/kdoctor-io/kdoctor/pkg/types"

	"github.com/shirou/gopsutil/cpu"
)

type UsedResource struct {
	mem        uint64  // byte
	cpu        float64 // percent
	roundCount int     // Count the number of rounds of cpu mem
	totalCPU   float64 //Total cpu usage statistics
	l          lock.RWMutex
	ctx        context.Context
	stop       chan struct{}
}

func InitResource(ctx context.Context) *UsedResource {
	return &UsedResource{
		ctx:  ctx,
		l:    lock.RWMutex{},
		stop: make(chan struct{}, 1),
	}
}

func (r *UsedResource) RunResourceCollector() {
	interval := time.Duration(types.AgentConfig.CollectResourceInSecond) * time.Second
	go func() {
		for {
			select {
			case <-r.stop:
				return
			default:
				if r.ctx.Err() != nil {
					return
				}
				cpuStats, err := cpu.Percent(interval, false)
				if err == nil {
					r.l.Lock()
					r.roundCount++
					if r.cpu < cpuStats[0] {
						r.cpu = cpuStats[0]
					}
					r.totalCPU += cpuStats[0]
					r.l.Unlock()
				}
				m := &runtime.MemStats{}
				runtime.ReadMemStats(m)
				r.l.Lock()
				if r.mem < m.Sys {
					r.mem = m.Sys
				}
				r.l.Unlock()
			}
		}

	}()
}

func (r *UsedResource) Stats() v1beta1.SystemResource {
	r.l.Lock()
	useCPU := r.cpu
	totalCPU := r.totalCPU
	roundCount := r.roundCount
	mem := r.mem
	r.l.Unlock()
	resource := v1beta1.SystemResource{
		MaxCPU:    fmt.Sprintf("%.3f%%", useCPU),
		MeanCPU:   fmt.Sprintf("%.3f%%", totalCPU/float64(roundCount)),
		MaxMemory: fmt.Sprintf("%.2fMB", float64(mem/(1024*1024))),
	}

	return resource
}

func (r *UsedResource) Stop() {
	r.stop <- struct{}{}
}
