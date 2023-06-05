// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package fileManager

import (
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	// defaults to 100 megabytes
	defaultMaxSizeMB = 100
	// 0 for no limit
	defaultMaxAgeDays = 0
	// o for no limit
	defaultMaxBackups = 0
)

func DefaultFileWriter(maxSizeMB int, maxAgeDays int, maxBackups int) {
	defaultMaxSizeMB = maxSizeMB
	defaultMaxAgeDays = maxAgeDays
	defaultMaxBackups = maxBackups
}

// the logger will auto create the directory for the file
// lumberjack will handle the max size / age / rotate when write or close .
// So, if you will not use lumberjack to operate a file, it needs manually remove the file with MaxAge
func NewFileWriter(filePath string) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    defaultMaxSizeMB,
		MaxAge:     defaultMaxAgeDays,
		MaxBackups: defaultMaxBackups,
		Compress:   false,
	}
}
