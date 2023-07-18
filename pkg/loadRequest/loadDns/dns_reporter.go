// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

// this is copied from Google project, and make some code modification
// Copyright 2014 Google Inc. All Rights Reserved.

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

// Based on https://github.com/rakyll/hey/blob/master/requester/report.go
//
// Changes:
// - remove metrics that we don't use

package loadDns

import (
	"github.com/miekg/dns"
	"time"
)

// We report for max 1M results.
const maxRes = 1000000

type report struct {
	avgTotal float64
	average  float64
	tps      float64

	results      chan *result
	done         chan bool
	total        time.Duration
	errorDist    map[string]int
	lats         []float32
	totalCount   int64
	successCount int64
	failedCount  int64
	ReplyCode    map[string]int
}

func newReport(results chan *result) *report {
	return &report{
		results:   results,
		done:      make(chan bool, 1),
		errorDist: make(map[string]int),
		lats:      make([]float32, 0, maxRes),
		ReplyCode: make(map[string]int),
	}
}

func runReporter(r *report) {
	// Loop will continue until channel is closed
	for res := range r.results {
		r.totalCount++
		if res.err != nil {
			r.errorDist[res.err.Error()]++
			r.failedCount++
		} else {
			r.avgTotal += res.duration.Seconds()
			if len(r.lats) < maxRes {
				r.lats = append(r.lats, float32(res.duration.Milliseconds()))
			}
			rcodeStr := dns.RcodeToString[res.msg.Rcode]
			r.ReplyCode[rcodeStr]++
			r.successCount++
		}
	}
	// Signal reporter is done.
	r.done <- true
}

func (r *report) finalize(total time.Duration) {
	r.total = total
	r.tps = float64(r.totalCount) / r.total.Seconds()
	r.average = r.avgTotal / float64(len(r.lats))
}
