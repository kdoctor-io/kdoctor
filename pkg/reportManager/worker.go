// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package reportManager

import (
	"context"
	"fmt"
	"github.com/kdoctor-io/kdoctor/pkg/grpcManager"
	k8sObjManager "github.com/kdoctor-io/kdoctor/pkg/k8ObjManager"
	crd "github.com/kdoctor-io/kdoctor/pkg/k8s/apis/kdoctor.io/v1beta1"
	"github.com/kdoctor-io/kdoctor/pkg/types"
	"github.com/kdoctor-io/kdoctor/pkg/utils"
	"go.uber.org/zap"
	"net"
	"path"
	"strings"
	"time"
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
		taskName := v[1]
		if !strings.Contains(podName, taskName) {
			logger.Sugar().Debugf("task %s not task of pod %s ,skip sync report %s", taskName, podName, remoteFileName)
			continue
		}
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

func (s *reportManager) runControllerAggregateReportOnce(ctx context.Context, logger *zap.Logger) error {

	// grpc client
	grpcClient := grpcManager.NewGrpcClient(s.logger.Named("grpc"), true)

	localFileList, e := utils.GetFileList(s.reportDir)
	if e != nil {
		logger.Sugar().Errorf("failed to get local report files underlay %v, error=%v ", s.reportDir, e)
		// ignore , no retry
		return nil
	}
	logger.Sugar().Debugf("before sync, local report files: %v", localFileList)
	// get all runtime obj
	for _, v := range s.runtimeDB {

		for _, m := range v.List() {
			// only aggregate created runtime report
			if m.RuntimeStatus != crd.RuntimeCreated {
				continue
			}
			var podIP k8sObjManager.PodIps
			var err error
			if m.RuntimeKind == types.KindDaemonSet {
				podIP, err = k8sObjManager.GetK8sObjManager().ListDaemonsetPodIPs(context.Background(), m.RuntimeName, types.ControllerConfig.PodNamespace)
			}
			if m.RuntimeKind == types.KindDeployment {
				podIP, err = k8sObjManager.GetK8sObjManager().ListDeploymentPodIPs(context.Background(), m.RuntimeName, types.ControllerConfig.PodNamespace)
			}
			logger.Sugar().Debugf("podIP : %v", podIP)
			if err != nil {
				m := fmt.Sprintf("failed to get kind %s name %s agent ip, error=%v", m.RuntimeKind, m.RuntimeName, err)
				logger.Error(m)
				// retry
				return fmt.Errorf(m)
			}

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
		}
	}

	return nil
}

// just one worker to sync all report and save to local disc of controller pod
func (s *reportManager) syncHandler(ctx context.Context, triggerName string) error {
	logger := s.logger.With(
		zap.String("triggerSource", triggerName),
	)
	return s.runControllerAggregateReportOnce(ctx, logger)
}
