// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package pluginManager

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/strings/slices"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	crd "github.com/kdoctor-io/kdoctor/pkg/k8s/apis/kdoctor.io/v1beta1"
	plugintypes "github.com/kdoctor-io/kdoctor/pkg/pluginManager/types"
	"github.com/kdoctor-io/kdoctor/pkg/reportManager"
	"github.com/kdoctor-io/kdoctor/pkg/scheduler"
	"github.com/kdoctor-io/kdoctor/pkg/types"
)

func (s *pluginControllerReconciler) GetSpiderAgentNodeNotInRecord(ctx context.Context, succeedNodeList []string, podMatchLabel client.MatchingLabels) ([]string, error) {
	allNodeList, failNodeList := []string{}, []string{}

	podList := corev1.PodList{}

	err := s.client.List(ctx, &podList,
		podMatchLabel,
		client.InNamespace(types.ControllerConfig.PodNamespace),
	)
	if nil != err {
		return nil, err
	}

	if len(podList.Items) == 0 {
		return nil, fmt.Errorf("failed to find agent node with matchLabels '%s' namespace '%s'", podMatchLabel, types.ControllerConfig.PodNamespace)
	}

	for index := range podList.Items {
		allNodeList = append(allNodeList, podList.Items[index].Spec.NodeName)
	}

	// if the runtime is deployment, we may get duplicated node
	allNodeList = RemoveDuplicates[string](allNodeList)
	s.logger.Sugar().Debugf("all agent nodes: %v", allNodeList)

	// gather the failure Node list
	failNodeList = slices.Filter(failNodeList, allNodeList, func(s string) bool {
		return !slices.Contains(succeedNodeList, s)
	})

	return failNodeList, nil
}

func (s *pluginControllerReconciler) UpdateRoundFinalStatus(logger *zap.Logger, ctx context.Context, newStatus *crd.TaskStatus, runtimePodSelector client.MatchingLabels, deadline bool) (roundDone bool, err error) {
	latestRecord := &(newStatus.History[0])
	roundNumber := latestRecord.RoundNumber

	if latestRecord.Status == crd.StatusHistoryRecordStatusFail ||
		latestRecord.Status == crd.StatusHistoryRecordStatusSucceed ||
		latestRecord.Status == crd.StatusHistoryRecordStatusNotstarted {
		return true, nil
	}

	// when not reach deadline, ignore when nothing report
	if !deadline && len(latestRecord.SucceedAgentNodeList) == 0 && len(latestRecord.FailedAgentNodeList) == 0 {
		logger.Sugar().Debugf("round %v not report anything", roundNumber)
		return false, nil
	}

	// update result in latestRecord
	reportNode := []string{}
	reportNode = append(reportNode, latestRecord.SucceedAgentNodeList...)
	reportNode = append(reportNode, latestRecord.FailedAgentNodeList...)
	unknownReportNodeList, e := s.GetSpiderAgentNodeNotInRecord(ctx, reportNode, runtimePodSelector)
	if e != nil {
		logger.Sugar().Errorf("round %v failed to GetSpiderAgentNodeNotInSucceedRecord, error=%v", roundNumber, e)
		return false, e
	}

	if len(unknownReportNodeList) > 0 && !deadline {
		// when not reach the deadline, ignore
		logger.Sugar().Debugf("round %v , partial agents did not reported, wait for deadline", roundNumber)
		return false, nil
	}

	// it's ok to collect round status (we meet the deadline)
	if len(unknownReportNodeList) > 0 || len(latestRecord.FailedAgentNodeList) > 0 {
		latestRecord.NotReportAgentNodeList = unknownReportNodeList
		n := crd.StatusHistoryRecordStatusFail
		latestRecord.Status = n
		newStatus.LastRoundStatus = &n
		logger.Sugar().Errorf("round %v failed , failedNode=%v, unknowReportNode=%v", roundNumber, latestRecord.FailedAgentNodeList, unknownReportNodeList)

		if len(latestRecord.FailedAgentNodeList) > 0 {
			latestRecord.FailureReason = "some agents failed"
		} else if len(unknownReportNodeList) > 0 {
			latestRecord.FailureReason = "some agents did not report"
		}
	} else {
		n := crd.StatusHistoryRecordStatusSucceed
		latestRecord.Status = n
		newStatus.LastRoundStatus = &n
		logger.Sugar().Infof("round %v succeeded ", latestRecord.RoundNumber)
	}
	cnt := len(reportNode) + len(unknownReportNodeList)
	latestRecord.ExpectedActorNumber = &cnt
	latestRecord.EndTimeStamp = &metav1.Time{
		Time: time.Now(),
	}
	i := time.Since(latestRecord.StartTimeStamp.Time).String()
	latestRecord.Duration = &i

	return true, nil
}

func (s *pluginControllerReconciler) WriteSummaryReport(taskName string, roundNumber int, newStatus *crd.TaskStatus) {
	if s.fm == nil {
		return
	}

	kindName := strings.Split(taskName, ".")[0]
	instanceName := strings.TrimPrefix(taskName, kindName+".")
	t := time.Duration(types.ControllerConfig.ReportAgeInDay*24) * time.Hour
	endTime := newStatus.History[0].StartTimeStamp.Add(t)

	if !s.fm.CheckTaskFileExisted(kindName, instanceName, roundNumber) {
		// add to workqueue to collect all report of last round, for node latestRecord.FailedAgentNodeList and latestRecord.SucceedAgentNodeList
		reportManager.TriggerSyncReport(fmt.Sprintf("%s.%d", taskName, roundNumber))

		// TODO (Icarus9913): change to use v1beta1.Report ?
		// write controller summary report
		msg := plugintypes.PluginReport{
			TaskName:       strings.ToLower(taskName),
			TaskSpec:       "",
			RoundNumber:    roundNumber,
			RoundResult:    plugintypes.RoundResultStatus(newStatus.History[0].Status),
			FailedReason:   newStatus.History[0].FailureReason,
			NodeName:       "",
			PodName:        types.ControllerConfig.PodName,
			StartTimeStamp: newStatus.History[0].StartTimeStamp.Time,
			EndTimeStamp:   time.Now(),
			RoundDuraiton:  time.Since(newStatus.History[0].StartTimeStamp.Time).String(),
			Detail:         newStatus.History[0],
			ReportType:     plugintypes.ReportTypeSummary,
		}

		if jsongByte, e := json.Marshal(msg); e != nil {
			s.logger.Sugar().Errorf("failed to generate round summary report for kind %v task %v round %v, json marsha error=%v", kindName, instanceName, roundNumber, e)
		} else {
			// print to stdout for human reading
			fmt.Printf("%+v\n ", string(jsongByte))

			var out bytes.Buffer
			if e := json.Indent(&out, jsongByte, "", "\t"); e != nil {
				s.logger.Sugar().Errorf("failed to generate round summary report for kind %v task %v round %v, json Indent error=%v", kindName, instanceName, roundNumber, e)
			} else {
				// file name format: fmt.Sprintf("%s_%s_round%d_%s_%s", kindName, taskName, roundNumber, nodeName, suffix)
				if e := s.fm.WriteTaskFile(kindName, instanceName, roundNumber, "summary", endTime, out.Bytes()); e != nil {
					s.logger.Sugar().Errorf("failed to generate round summary report for kind %v task %v round %v, write file error=%v", kindName, instanceName, roundNumber, e)
				} else {
					s.logger.Sugar().Debugf("succeeded to generate round summary report for kind %v task %v round %v", kindName, instanceName, roundNumber)
				}
			}
		}
	}
}

func (s *pluginControllerReconciler) UpdateStatus(logger *zap.Logger, ctx context.Context, oldStatus *crd.TaskStatus, schedulePlan *crd.SchedulePlan, runtimePodMatchLabels client.MatchingLabels, taskName string) (result *reconcile.Result, taskStatus *crd.TaskStatus, e error) {
	newStatus := oldStatus.DeepCopy()
	nextInterval := time.Duration(types.ControllerConfig.Configmap.TaskPollIntervalInSecond) * time.Second
	nowTime := time.Now()
	var startTime time.Time

	// init new instance first
	scheduler := NewSchedule(*schedulePlan.Schedule)
	if newStatus.ExpectedRound == nil || len(newStatus.History) == 0 {
		startTime = scheduler.StartTime(nowTime)
		m := int64(0)
		newStatus.DoneRound = &m
		newStatus.ExpectedRound = &schedulePlan.RoundNumber

		newRecord := NewStatusHistoryRecord(startTime, 1, schedulePlan)
		newStatus.History = append(newStatus.History, *newRecord)
		logger.Sugar().Debugf("initialize the first round of task : %v ", taskName, *newRecord)
		// trigger
		result = &reconcile.Result{
			Requeue: true,
		}
		// updating status firstly , it will trigger to handle it next round
		return result, newStatus, nil
	}

	// done task
	if *newStatus.DoneRound == *newStatus.ExpectedRound {
		return nil, nil, nil
	}

	latestRecord := &(newStatus.History[0])
	roundNumber := latestRecord.RoundNumber
	logger.Sugar().Debugf("current time:%v , latest history record: %+v", nowTime, latestRecord)
	logger.Sugar().Debugf("all history record: %+v", newStatus.History)

	switch {
	case nowTime.After(latestRecord.StartTimeStamp.Time) && nowTime.Before(latestRecord.DeadLineTimeStamp.Time):
		if latestRecord.Status == crd.StatusHistoryRecordStatusNotstarted {
			latestRecord.Status = crd.StatusHistoryRecordStatusOngoing
			// requeue immediately to make sure the update succeed , not conflicted
			result = &reconcile.Result{
				Requeue: true,
			}

		} else if latestRecord.Status == crd.StatusHistoryRecordStatusOngoing {
			logger.Debug("try to poll the status of task " + taskName)
			roundDone, e := s.UpdateRoundFinalStatus(logger, ctx, newStatus, runtimePodMatchLabels, false)
			if e != nil {
				return nil, nil, e
			}

			if roundDone {
				logger.Sugar().Infof("round %v get reports from all agents ", roundNumber)

				// before insert new record, write summary of last round
				s.WriteSummaryReport(taskName, roundNumber, newStatus)

				// add new round record
				if *(newStatus.DoneRound) < *(newStatus.ExpectedRound) || *newStatus.ExpectedRound == -1 {
					n := *(newStatus.DoneRound) + 1
					newStatus.DoneRound = &n
					startTime = scheduler.Next(latestRecord.StartTimeStamp.Time)
					if n < *(newStatus.ExpectedRound) || *newStatus.ExpectedRound == -1 {

						newRecord := NewStatusHistoryRecord(startTime, int(n+1), schedulePlan)

						tmp := append([]crd.StatusHistoryRecord{*newRecord}, newStatus.History...)
						if len(tmp) > types.ControllerConfig.Configmap.CrdMaxHistory {
							tmp = tmp[:(types.ControllerConfig.Configmap.CrdMaxHistory)]
						}
						newStatus.History = tmp

						logger.Sugar().Infof("insert new record for next round : %+v", *newRecord)
					} else {
						newStatus.Finish = true
						now := metav1.Now()
						newStatus.FinishTime = &now
					}
				}

				// requeue immediately to make sure the update succeed , not conflicted
				result = &reconcile.Result{
					Requeue: true,
				}
			} else {
				// trigger after interval
				result = &reconcile.Result{
					RequeueAfter: nextInterval,
				}
			}
		} else {
			logger.Debug("ignore poll the finished round of task " + taskName)

			// trigger when deadline
			result = &reconcile.Result{
				RequeueAfter: time.Until(latestRecord.DeadLineTimeStamp.Time),
			}
		}

	case nowTime.Before(latestRecord.StartTimeStamp.Time):
		fallthrough
	case nowTime.After(latestRecord.DeadLineTimeStamp.Time):
		if *newStatus.DoneRound == *newStatus.ExpectedRound {
			logger.Sugar().Debugf("task %s finish, ignore ", taskName)
			newStatus.Finish = true
			now := metav1.Now()
			newStatus.FinishTime = &now
			result = nil

		} else {

			// when task not finish , once we update the status succeed , we will not get here , it should go to case nowTime.Before(latestRecord.StartTimeStamp.Time)
			if latestRecord.Status == crd.StatusHistoryRecordStatusOngoing {
				// here, we should update last round status

				if _, e := s.UpdateRoundFinalStatus(logger, ctx, newStatus, runtimePodMatchLabels, true); e != nil {
					return nil, nil, e
				} else {
					// all agent finished, so try to update the summary
					logger.Sugar().Infof("round %v got reports from all agents, try to summarize", roundNumber)

					// before insert new record, write summary of last round
					s.WriteSummaryReport(taskName, roundNumber, newStatus)

					// add new round record
					if *(newStatus.DoneRound) < *(newStatus.ExpectedRound) || *newStatus.ExpectedRound == -1 {
						n := *(newStatus.DoneRound) + 1
						newStatus.DoneRound = &n

						if n < *(newStatus.ExpectedRound) || *newStatus.ExpectedRound == -1 {
							startTime = scheduler.Next(latestRecord.StartTimeStamp.Time)
							newRecord := NewStatusHistoryRecord(startTime, int(n+1), schedulePlan)
							tmp := append([]crd.StatusHistoryRecord{*newRecord}, newStatus.History...)
							if len(tmp) > types.ControllerConfig.Configmap.CrdMaxHistory {
								tmp = tmp[:(types.ControllerConfig.Configmap.CrdMaxHistory)]
							}
							newStatus.History = tmp

							logger.Sugar().Infof("insert new record for next round : %+v", *newRecord)
						} else {
							newStatus.Finish = true
							now := metav1.Now()
							newStatus.FinishTime = &now
						}
					}

					// requeue immediately to make sure the update succeed , not conflicted
					result = &reconcile.Result{
						Requeue: true,
					}
				}
			} else {
				// round finish
				// trigger when next round start
				currentLatestRecord := &(newStatus.History[0])
				logger.Sugar().Infof("task %v wait for next round %v at %v", taskName, currentLatestRecord.RoundNumber, currentLatestRecord.StartTimeStamp)
				result = &reconcile.Result{
					RequeueAfter: time.Until(currentLatestRecord.StartTimeStamp.Time),
				}
			}
		}
	}

	return result, newStatus, nil

}

func (s *pluginControllerReconciler) TaskResourceReconcile(ctx context.Context, taskKind string, ownerTask metav1.Object, agentSpec crd.AgentSpec, taskStatus *crd.TaskStatus, logger *zap.Logger) (*crd.TaskStatus, error) {
	var resource crd.TaskResource
	var err error
	var deletionTime *metav1.Time

	if taskStatus.Resource == nil {
		logger.Sugar().Debugf("task '%s/%s' just created, try to initial its corresponding runtime resource", taskKind, ownerTask.GetName())
		newScheduler := scheduler.NewScheduler(s.client, s.apiReader, taskKind, ownerTask.GetName(), s.runtimeUniqueMatchLabelKey, logger)
		// create the task corresponding resources(runtime,service) and record them to the task CR object subresource with 'Creating' status
		resource, err = newScheduler.CreateTaskRuntimeIfNotExist(ctx, ownerTask, agentSpec)
		if nil != err {
			return nil, err
		}
		taskStatus.Resource = &resource
	} else {
		// we need to track it again, in order to avoid controller restart
		resource = *taskStatus.Resource
		if taskStatus.FinishTime != nil {
			deletionTime = taskStatus.FinishTime.DeepCopy()
			if agentSpec.TerminationGracePeriodMinutes != nil {
				newTime := metav1.NewTime(deletionTime.Add(time.Duration(*agentSpec.TerminationGracePeriodMinutes) * time.Minute))
				deletionTime = newTime.DeepCopy()
			}
			logger.Sugar().Debugf("task '%s/%s' finish time '%s' and runtime deletion time '%s'",
				taskKind, ownerTask.GetName(), taskStatus.FinishTime, deletionTime)
		}
	}

	// record the task resource to the tracker DB, and the tracker will update the task subresource resource status or delete corresponding runtime asynchronously
	err = s.tracker.DB.Apply(scheduler.BuildItem(resource, taskKind, ownerTask.GetName(), deletionTime))
	if nil != err {
		return nil, err
	}

	return taskStatus, nil
}
