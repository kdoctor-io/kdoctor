// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package logger_test

import (
	"github.com/kdoctor-io/kdoctor/pkg/logger"
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("test logger ", Label("logger"), func() {
	It("test debug logger  ", func() {
		logger.NewStdoutLogger("debug", "test")
	})
	It("test info logger  ", func() {
		logger.NewStdoutLogger("info", "test")
	})
	It("test warn logger  ", func() {
		logger.NewStdoutLogger("warn", "test")
	})
	It("test error logger  ", func() {
		logger.NewStdoutLogger("error", "test")
	})
	It("test fatal logger  ", func() {
		logger.NewStdoutLogger("fatal", "test")
	})
	It("test panic logger  ", func() {
		logger.NewStdoutLogger("panic", "test")
	})
	It("test empty logger  ", func() {
		logger.NewStdoutLogger("", "test")
	})
})
