// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package pluginManager

import (
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	crd "github.com/kdoctor-io/kdoctor/pkg/k8s/apis/kdoctor.io/v1beta1"
)

func NewStatusHistoryRecord(startTime time.Time, RoundNumber int, schedulePlan *crd.SchedulePlan) *crd.StatusHistoryRecord {
	newRecod := crd.StatusHistoryRecord{
		Status:                 crd.StatusHistoryRecordStatusNotstarted,
		FailureReason:          "",
		RoundNumber:            RoundNumber,
		SucceedAgentNodeList:   []string{},
		FailedAgentNodeList:    []string{},
		NotReportAgentNodeList: []string{},
	}
	newRecod.StartTimeStamp = metav1.NewTime(startTime)

	adder := time.Duration(schedulePlan.RoundTimeoutMinute) * time.Minute
	endTime := startTime.Add(adder)
	newRecod.DeadLineTimeStamp = metav1.NewTime(endTime)

	return &newRecod
}

func CheckItemInList(item string, checklist []string) (bool, error) {
	if len(item) == 0 {
		return false, fmt.Errorf("empty item")
	}
	if len(checklist) == 0 {
		return false, nil
	}
	for _, v := range checklist {
		if v == item {
			return true, nil
		}
	}
	return false, nil
}

type element interface {
	~int | ~int32 | ~int64 |
		~float64 | ~float32 |
		~string
}

// RemoveDuplicates would clean up the duplicate element in a given slice
func RemoveDuplicates[T element](arr []T) []T {
	var newArr []T
	dic := make(map[T]struct{})

	for _, tmpElement := range arr {
		dic[tmpElement] = struct{}{}
	}

	for tmpItem := range dic {
		newArr = append(newArr, tmpItem)
	}

	return newArr
}
