// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package resource

import (
	"context"
	. "github.com/onsi/ginkgo/v2"
	"time"
)

var _ = Describe("resource unit test", Label("resource"), func() {

	It("get system resource", func() {

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		r := InitResource(ctx)
		r.RunResourceCollector()
		time.Sleep(time.Second * 5)
		GinkgoWriter.Printf("resource stats %+v", r.Stats())
		time.Sleep(time.Second * 5)
		r.Stop()
	})

})
