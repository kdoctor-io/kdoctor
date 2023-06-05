// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package tools

import (
	"fmt"
	crd "github.com/kdoctor-io/kdoctor/pkg/k8s/apis/kdoctor.io/v1beta1"
	"github.com/robfig/cron"
	"strconv"
	"strings"
)

func ValidataCrdSchedule(plan *crd.SchedulePlan) error {

	if plan == nil {
		return fmt.Errorf("Schedule is empty ")
	}

	args := strings.Split(*plan.Schedule, " ")

	if len(args) == 2 {
		startAfterMinute, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("The format of the schedule is incorrect, it should be number ")
		}
		intervalMinute, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("The format of the schedule is incorrect, it should be number ")
		}
		if startAfterMinute < 0 {
			return fmt.Errorf("Schedule.StartAfterMinute %v must not be smaller than 0 ", startAfterMinute)
		}

		if intervalMinute < 1 {
			return fmt.Errorf("Schedule.IntervalMinute %v must not be smaller than 1 ", intervalMinute)
		}

		if int(plan.RoundTimeoutMinute) > intervalMinute {
			return fmt.Errorf("Schedule.RoundTimeoutMinute %v must not be bigger than Schedule.IntervalMinute %v ", plan.RoundTimeoutMinute, intervalMinute)
		}

	} else if len(args) == 5 {
		_, err := cron.ParseStandard(*plan.Schedule)
		if err != nil {
			return fmt.Errorf("Crontab configuration error,err: %v ", err)
		}

	} else {
		return fmt.Errorf("The format of the schedule is incorrect, it should be two or five ")
	}

	if plan.RoundTimeoutMinute < 1 {
		return fmt.Errorf("Schedule.RoundTimeoutMinute %v must not be smaller than 1 ", plan.RoundTimeoutMinute)
	}

	return nil
}

func GetDefaultSchedule() (plan *crd.SchedulePlan) {
	s := "0 60"
	return &crd.SchedulePlan{
		RoundTimeoutMinute: 60,
		Schedule:           &s,
		RoundNumber:        1,
	}
}

func GetDefaultNetSuccessCondition() (plan *crd.NetSuccessCondition) {
	n := float64(1)
	return &crd.NetSuccessCondition{
		SuccessRate: &n,
	}
}
