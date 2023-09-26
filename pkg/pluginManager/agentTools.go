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

	"github.com/kdoctor-io/kdoctor/pkg/resource"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	crd "github.com/kdoctor-io/kdoctor/pkg/k8s/apis/kdoctor.io/v1beta1"
	systemv1beta1 "github.com/kdoctor-io/kdoctor/pkg/k8s/apis/system/v1beta1"
	plugintypes "github.com/kdoctor-io/kdoctor/pkg/pluginManager/types"
	"github.com/kdoctor-io/kdoctor/pkg/taskStatusManager"
	"github.com/kdoctor-io/kdoctor/pkg/types"
)

// call plugin to implement the round task and collect the report
func (s *pluginAgentReconciler) CallPluginImplementRoundTask(logger *zap.Logger, obj runtime.Object, schedulePlan *crd.SchedulePlan, taskName string, roundNumber int, crdObjSpec interface{}) {
	taskRoundName := fmt.Sprintf("%s.round%d", taskName, roundNumber)

	roundDuration := time.Duration(schedulePlan.RoundTimeoutMinute) * time.Minute
	ctx, cancel := context.WithTimeout(context.Background(), roundDuration)
	defer cancel()
	taskSucceed := make(chan bool)
	logger.Sugar().Infof("plugin begins to implement, expect deadline %v, ", roundDuration.String())

	go func() {
		// process mem cpu stats
		resourceStats := resource.InitResource(ctx)
		resourceStats.RunResourceCollector()
		startTime := metav1.Now()
		msg := &systemv1beta1.Report{
			TaskName:       strings.ToLower(taskName),
			RoundNumber:    int64(roundNumber),
			NodeName:       s.localNodeName,
			PodName:        types.AgentConfig.PodName,
			FailedReason:   nil,
			StartTimeStamp: startTime,
			ReportType:     plugintypes.ReportTypeAgent,
		}
		failureReason, report, e := s.plugin.AgentExecuteTask(logger, ctx, obj, resourceStats)

		if e != nil {
			logger.Sugar().Errorf("plugin failed to implement the round task, error=%v", e)
			taskSucceed <- false
			msg.RoundResult = string(plugintypes.RoundResultFail)
			msg.FailedReason = pointer.String(e.Error())
		} else {
			select {
			case <-ctx.Done():
				logger.Sugar().Errorf("plugin finished the round task, timeout, it takes %v , logger than expected %s", time.Since(startTime.Time).String(), roundDuration.String())
				taskSucceed <- false
				msg.RoundResult = string(plugintypes.RoundResultFail)
				msg.FailedReason = pointer.String("implementing timeout")
			default:
				logger.Sugar().Infof("plugin finished the round task, it takes %v , shorter than expected %s", time.Since(startTime.Time).String(), roundDuration.String())
				if len(failureReason) == 0 {
					taskSucceed <- true
					msg.RoundResult = string(plugintypes.RoundResultSucceed)
					msg.FailedReason = nil
				} else {
					taskSucceed <- false
					msg.RoundResult = string(plugintypes.RoundResultFail)
					msg.FailedReason = pointer.String(failureReason)
				}
			}
			msg.TaskType = report.KindTask()
		}
		endTime := metav1.Now()
		msg.EndTimeStamp = endTime
		msg.RoundDuration = endTime.Sub(startTime.Time).String()
		if report != nil {
			err := s.plugin.SetReportWithTask(msg, crdObjSpec, report)
			if nil != err {
				// TODO (Icarus9913): improve the error solution
				logger.Sugar().Errorf("failed to set task details to report, error: %v", err)
			}
		}

		if jsongByte, err := json.Marshal(msg); err != nil {
			logger.Sugar().Errorf("failed to generate round report , marsha json error=%v", err)
		} else {
			// print to stdout for human reading
			fmt.Printf("%+v\n ", string(jsongByte))

			// write report to disk for controller to collect
			if s.fm != nil {
				var out bytes.Buffer
				if e := json.Indent(&out, jsongByte, "", "\t"); e != nil {
					logger.Sugar().Errorf("failed to json Indent for report of %v, error=%v", taskRoundName, e)
				} else {
					kindName := strings.Split(taskName, ".")[0]
					instanceName := strings.TrimPrefix(taskName, kindName+".")
					// save with maximum age roundDuration , in this interval, the controller also will collect it
					t := time.Duration(schedulePlan.RoundTimeoutMinute+5) * time.Minute

					// file name format: fmt.Sprintf("%s_%s_round%d_%s_%s", kindName, taskName, roundNumber, nodeName, suffix)
					if e := s.fm.WriteTaskFile(kindName, instanceName, roundNumber, s.localNodeName, time.Now().Add(t), out.Bytes()); e != nil {
						logger.Sugar().Errorf("failed to write report of %v, error=%v", taskRoundName, e)
					} else {
						logger.Sugar().Debugf("succeed to write report for %v", taskRoundName)
					}
				}
			}
		}

	}()

	select {
	case <-ctx.Done():
		logger.Sugar().Errorf("timeout for getting result from plugin, the round task failed")
		s.taskRoundData.SetTask(taskRoundName, taskStatusManager.RoundStatusFail)
	case r := <-taskSucceed:
		logger.Sugar().Infof("succeed to call plugin to implement round task, succeed=%v", r)
		if r {
			s.taskRoundData.SetTask(taskRoundName, taskStatusManager.RoundStatusSucceeded)
		} else {
			s.taskRoundData.SetTask(taskRoundName, taskStatusManager.RoundStatusFail)
		}
	}

	// delete data
	go func() {
		time.Sleep(roundDuration)
		s.taskRoundData.DeleteTask(taskRoundName)
	}()
}

func (s *pluginAgentReconciler) HandleAgentTaskRound(logger *zap.Logger, ctx context.Context, oldStatus *crd.TaskStatus, schedulePlan *crd.SchedulePlan, obj runtime.Object, taskName string, crdObjSpec interface{}) (result *reconcile.Result, taskStatus *crd.TaskStatus, e error) {
	newStatus := oldStatus.DeepCopy()
	nowTime := time.Now()

	if newStatus.ExpectedRound == nil || len(newStatus.History) == 0 || *newStatus.DoneRound == *newStatus.ExpectedRound {
		// not start or all finish
		return nil, nil, nil
	}

	latestRecord := &(newStatus.History[0])
	logger.Sugar().Debugf("current time:%v , latest history record: %+v", nowTime, latestRecord)
	// logger.Sugar().Debugf("all history record: %+v", newStatus.History)

	if latestRecord.Status != crd.StatusHistoryRecordStatusOngoing {
		logger.Sugar().Debugf("ignore task %v , no opportunity to implement ", taskName)
		return nil, nil, nil
	}

	// check whether we have reported the round result
	if len(latestRecord.SucceedAgentNodeList) != 0 || len(latestRecord.FailedAgentNodeList) != 0 {
		v := []string{}
		v = append(v, latestRecord.SucceedAgentNodeList...)
		v = append(v, latestRecord.FailedAgentNodeList...)
		logger.Sugar().Debugf("check whether localNode %v has report ", s.localNodeName)

		if ok, e := CheckItemInList(s.localNodeName, v); e != nil {
			logger.Sugar().Errorf("failed to check local node in task record")
			// no need to requeue
			return nil, nil, nil
		} else {
			if ok {
				logger.Sugar().Debugf("ignore task %v , it has reported for the round result", taskName)
				return nil, nil, nil
			}
		}
	}

	taskRoundName := fmt.Sprintf("%s.round%d", taskName, latestRecord.RoundNumber)
	nextInterval := time.Duration(types.AgentConfig.Configmap.TaskPollIntervalInSecond) * time.Second

	if status, existed := s.taskRoundData.CheckTask(taskRoundName); !existed {
		// mark to started it
		s.taskRoundData.SetTask(taskRoundName, taskStatusManager.RoundStatusOngoing)

		// we still have not reported the result for an ongoing round. do it
		go s.CallPluginImplementRoundTask(logger.Named(taskRoundName), obj, schedulePlan, taskName, latestRecord.RoundNumber, crdObjSpec)
		logger.Sugar().Infof("task %v , trigger to implement task round, and try to poll report after %v second", taskRoundName, types.AgentConfig.Configmap.TaskPollIntervalInSecond)

		// trigger to poll result after interval
		result = &reconcile.Result{
			RequeueAfter: nextInterval,
		}

	} else {
		if status == taskStatusManager.RoundStatusOngoing {
			// the task is on going
			logger.Sugar().Infof("task %v is going, try to poll report after %v second", taskRoundName, types.AgentConfig.Configmap.TaskPollIntervalInSecond)

			// trigger to poll result after interval
			result = &reconcile.Result{
				RequeueAfter: nextInterval,
			}

		} else {
			// the task finish, report
			if status == taskStatusManager.RoundStatusSucceeded {
				logger.Sugar().Infof("task %v , report to succeed", taskRoundName)
				latestRecord.SucceedAgentNodeList = append(latestRecord.SucceedAgentNodeList, s.localNodeName)
			} else {
				logger.Sugar().Infof("task %v , report to fail", taskRoundName)
				latestRecord.FailedAgentNodeList = append(latestRecord.FailedAgentNodeList, s.localNodeName)
			}
			// requeue immediately to make sure the update succeed , not conflicted
			result = &reconcile.Result{
				Requeue: true,
			}

		}
	}

	return result, newStatus, nil
}
