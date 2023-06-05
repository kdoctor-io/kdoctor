// Copyright 2014 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0
// Remove unwanted metrics, optimize memory consumption, and optimize QPS requests

package loadHttp

import (
	"time"
)

// We report for max 1M results.
// todo (ii2day) configmap limit
const maxRes = 1000000

type report struct {
	// transactions Per Second
	tps float64

	results chan *result
	done    chan bool

	total       time.Duration
	statusCodes map[int]int
	errorDist   map[string]int
	latencies   []float32
	sizeTotal   int64
	totalCount  int64
}

func newReport(results chan *result, n int) *report {
	return &report{
		results:     results,
		done:        make(chan bool, 1),
		errorDist:   make(map[string]int),
		latencies:   make([]float32, 0, maxRes),
		statusCodes: make(map[int]int),
	}
}

func runReporter(r *report) {
	// Loop will continue until channel is closed
	for res := range r.results {
		r.totalCount++
		r.statusCodes[res.statusCode]++
		if res.err != nil {
			r.errorDist[res.err.Error()]++
		} else {
			if len(r.latencies) < maxRes {
				r.latencies = append(r.latencies, float32(res.duration.Milliseconds()))
			}
			if res.contentLength > 0 {
				r.sizeTotal += res.contentLength
			}
		}
	}
	// Signal reporter is done.
	r.done <- true
}

func (r *report) finalize(total time.Duration) {
	r.total = total
	r.tps = float64(r.totalCount) / r.total.Seconds()
}
