// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package pluginManager

import (
	"github.com/robfig/cron"
	"strconv"
	"strings"
	"time"
)

type Schedule interface {
	Next(time.Time) time.Time
	StartTime(time.Time) time.Time
}

type schedule struct {
	IsCron           bool
	StartAfterMinute int
	IntervalMinute   int
	CronSchedule     cron.Schedule
}

func NewSchedule(s string) Schedule {
	scheduler := &schedule{}
	args := strings.Split(s, " ")

	// crontab
	if len(args) == 5 {
		scheduler.IsCron = true
		cronSchedule, _ := cron.ParseStandard(s)
		scheduler.CronSchedule = cronSchedule
		// simple
	} else if len(args) == 2 {
		scheduler.IsCron = false
		StartAfterMinute, _ := strconv.Atoi(args[0])
		scheduler.StartAfterMinute = StartAfterMinute
		intervalMinute, _ := strconv.Atoi(args[1])
		scheduler.IntervalMinute = intervalMinute
	}
	return scheduler

}

func (s *schedule) Next(t time.Time) time.Time {

	if s.IsCron {
		return s.CronSchedule.Next(t)
	}

	return t.Add(time.Duration(s.IntervalMinute) * time.Minute)
}

func (s *schedule) StartTime(t time.Time) time.Time {

	if s.IsCron {
		return s.CronSchedule.Next(t)
	}

	return t.Add(time.Duration(s.StartAfterMinute) * time.Minute)
}
