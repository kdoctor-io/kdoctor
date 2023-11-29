// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0
package taskStatusManager_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTaskStatusManager(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "task status Suite")
}

var _ = BeforeSuite(func() {
	// nothing to do
})
