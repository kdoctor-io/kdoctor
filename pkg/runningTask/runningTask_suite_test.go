// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package runningTask

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestRunningTask(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "RunningTask Suite")
}

var _ = BeforeSuite(func() {
	// nothing to do
})
