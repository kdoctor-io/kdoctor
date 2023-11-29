// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package taskStatusManager_test

import (
	"github.com/kdoctor-io/kdoctor/pkg/taskStatusManager"
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("test task status", Label("task status"), func() {
	It("test status  ", func() {
		taskManager := taskStatusManager.NewTaskStatus()
		taskManager.SetTask("test", taskStatusManager.RoundStatusOngoing)
		taskManager.CheckTask("test")
		taskManager.DeleteTask("test")
	})
})
