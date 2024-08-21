// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package reportManager

import (
	"context"
	"fmt"
	"net"
	"path"
	"strings"
	"time"

	"github.com/kdoctor-io/kdoctor/pkg/grpcManager"
	k8sObjManager "github.com/kdoctor-io/kdoctor/pkg/k8ObjManager"
	"github.com/kdoctor-io/kdoctor/pkg/scheduler"
	"github.com/kdoctor-io/kdoctor/pkg/types"
	"github.com/kdoctor-io/kdoctor/pkg/utils"
	"go.uber.org/zap"
)

func GetMissRemoteReport(remoteFileList []string, localFileList []string) []string {
	remoteMissFileList := []string{}

	for _, remotefileName := range remoteFileList {
		// file name format: fmt.Sprintf("%s_%s_round%d_%s_%s", kindName, taskName, roundNumber, nodeName, suffix)
		v := strings.Split(remotefileName, "_")
		timeSuffix1 := v[len(v)-1]
		remoteFilePre := strings.TrimSuffix(remotefileName, "_"+timeSuffix1)

		bingo := false
		for _, localFileName := range localFileList {
			a := strings.Split(localFileName, "_")
			timeSuffix2 := a[len(a)-1]
			localFilePre := strings.TrimSuffix(localFileName, "_"+timeSuffix2)

			// fmt.Printf("compare: local %v  / %v , remote %v\n", localFileName, localFilePre, remoteFilePre)
			if localFilePre == remoteFilePre {
				bingo = true
				break
			}
		}
		if !bingo {
			remoteMissFileList = append(remoteMissFileList, remotefileName)
		}
	}
	return remoteMissFileList
}

func (s *reportManager) syncReportFromOneAgent(ctx context.Context, logger *zap.Logger, client grpcManager.GrpcClientManager, localFileList []string, podName, address string) {
	logger.Sugar().Debugf("sync report from agent %v with grpc address %v", podName, address)

	remoteFilesList, e := client.GetFileList(ctx, address, types.ControllerConfig.DirPathAgentReport)
	if e != nil {
		logger.Sugar().Errorf("%v", e)
		return
	}

	logger.Sugar().Debugf("agent pod %v has reports: %v", podName, remoteFilesList)
	logger.Sugar().Debugf("local has reports: %v", localFileList)
	missRemoteFileList := GetMissRemoteReport(remoteFilesList, localFileList)
	logger.Sugar().Debugf("try to sync pod %v reports: %v", podName, missRemoteFileList)

	for _, remoteFileName := range missRemoteFileList {
		// --
		remoteAbsPath := path.Join(types.ControllerConfig.DirPathAgentReport, remoteFileName)
		// --
		v := strings.Split(remoteFileName, "_")
		timeSuffix := v[len(v)-1]
		remoteFilePre := strings.TrimSuffix(remoteFileName, "_"+timeSuffix)
		// file name format: fmt.Sprintf("%s_%s_round%d_%s_%s", kindName, taskName, roundNumber, nodeName, suffix)
		t := time.Duration(types.ControllerConfig.ReportAgeInDay*24) * time.Hour
		suffix := time.Now().Add(t).Format(time.RFC3339)
		localFileName := fmt.Sprintf("%s_%v", remoteFilePre, suffix)
		localAbsPath := path.Join(types.ControllerConfig.DirPathControllerReport, localFileName)
		// --
		if e := client.SaveRemoteFileToLocal(ctx, address, remoteAbsPath, localAbsPath); e != nil {
			logger.Sugar().Errorf("failed to save remote file %v of pod %v to local file %v, error=%v", remoteAbsPath, podName, localAbsPath, e)
		} else {
			logger.Sugar().Infof("succeeded to save remote file %v of pod %v to local file %v", remoteAbsPath, podName, localAbsPath)
		}
	}

}

func (s *reportManager) runControllerAggregateReportOnce(ctx context.Context, logger *zap.Logger, taskKind string, taskName string) error {
	var task scheduler.Item
	var err error
	var podIP k8sObjManager.PodIps
	// grpc client
	grpcClient := grpcManager.NewGrpcClient(s.logger.Named("grpc"), true)

	localFileList, e := utils.GetFileList(s.reportDir)
	if e != nil {
		logger.Sugar().Errorf("failed to get local report files underlay %v, error=%v ", s.reportDir, e)
		// ignore , no retry
		return nil
	}
	logger.Sugar().Debugf("before sync, local report files: %v", localFileList)

	switch taskKind {
	case types.KindNameAppHttpHealthy:
		task, err = s.appHttpHealthyRuntimeDB.Get(taskName)
	case types.KindNameNetReach:
		task, err = s.netReachRuntimeDB.Get(taskName)
	case types.KindNameNetdns:
		task, err = s.netDNSRuntimeDB.Get(taskName)
	}
	if err != nil {
		return err
	}

	if task.RuntimeKind == types.KindDaemonSet {
		podIP, err = k8sObjManager.GetK8sObjManager().ListDaemonsetPodIPs(context.Background(), task.RuntimeName, types.ControllerConfig.PodNamespace)
	}
	if task.RuntimeKind == types.KindDeployment {
		podIP, err = k8sObjManager.GetK8sObjManager().ListDeploymentPodIPs(context.Background(), task.RuntimeName, types.ControllerConfig.PodNamespace)
	}
	if err != nil {
		errMsg := fmt.Errorf("failed to get kind %s name %s agent ip, error=%v", task.RuntimeKind, task.RuntimeName, err)
		logger.Error(errMsg.Error())
		// retry
		return errMsg
	}
	logger.Sugar().Debugf("podIP : %v", podIP)

	for podName, podIpInfo := range podIP {
		// get pod ip
		if len(podIpInfo) == 0 {
			logger.Sugar().Errorf("failed to get agent %s ip ", podName)
			continue
		}
		var podip string
		if types.ControllerConfig.Configmap.EnableIPv4 {
			podip = podIpInfo[0].IPv4
		} else {
			podip = podIpInfo[0].IPv6
		}
		if len(podip) == 0 {
			logger.Sugar().Errorf("failed to get agent %s ip ", podName)
			continue
		}

		ip := net.ParseIP(podip)
		var address string
		if ip.To4() == nil {
			address = fmt.Sprintf("[%s]:%d", podip, types.ControllerConfig.AgentGrpcListenPort)
		} else {
			address = fmt.Sprintf("%s:%d", podip, types.ControllerConfig.AgentGrpcListenPort)
		}
		s.syncReportFromOneAgent(ctx, logger, grpcClient, localFileList, podName, address)
	}

	return nil
}

// just one worker to sync all report and save to local disc of controller pod
func (s *reportManager) syncHandler(ctx context.Context, trigger string) error {
	logger := s.logger.With(
		zap.String("triggerSource", trigger),
	)
	return s.runControllerAggregateReportOnce(ctx, logger, strings.Split(trigger, ".")[0], strings.Split(trigger, ".")[1])
}
