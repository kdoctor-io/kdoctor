// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package resource

import (
	"context"
	"runtime"
	"time"

	"github.com/kdoctor-io/kdoctor/pkg/lock"
	"github.com/kdoctor-io/kdoctor/pkg/types"

	"github.com/shirou/gopsutil/cpu"
)

type UsedResource struct {
	mem uint64  // byte
	cpu float64 // percent
	l   lock.RWMutex
	ctx context.Context
}

func InitResource(ctx context.Context) *UsedResource {
	return &UsedResource{
		ctx: ctx,
	}
}

func (r *UsedResource) RunResourceCollector() {
	interval := time.Duration(types.AgentConfig.CollectResourceInSecond) * time.Second
	go func() {
		for {
			if r.ctx.Err() != nil {
				return
			}
			cpuStats, err := cpu.Percent(interval, false)
			if err == nil {
				if r.cpu < cpuStats[0] {
					r.cpu = cpuStats[0]
				}
			}
			m := &runtime.MemStats{}
			runtime.ReadMemStats(m)
			if r.mem < m.Sys {
				r.mem = m.Sys
			}
		}
	}()
}

func (r *UsedResource) Stats() (uint64, float64) {
	r.l.RLock()
	defer r.l.RUnlock()
	m := r.mem
	c := r.cpu
	return m, c
}

func (r *UsedResource) CleanStats() {
	r.l.Lock()
	defer r.l.Unlock()
	r.mem = 0
	r.cpu = 0
}
