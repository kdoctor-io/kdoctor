// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0
package example_test

import (
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("example ", Label("example"), func() {

	It("example", Label("example-1"), func() {
		GinkgoWriter.Printf("cluster information: %+v \n", frame.Info)

	})
})
