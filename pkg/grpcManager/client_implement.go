// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package grpcManager

import (
	"bufio"
	"context"
	"fmt"
	"github.com/kdoctor-io/kdoctor/api/v1/agentGrpc"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"os"
	"strings"
)

func (s *grpcClientManager) SendRequestForExecRequest(ctx context.Context, serverAddress []string, requestList []*agentGrpc.ExecRequestMsg) ([]*agentGrpc.ExecResponseMsg, error) {

	logger := s.logger.With(
		zap.String("server", fmt.Sprintf("%v", serverAddress)),
	)

	if e := s.clientDial(ctx, serverAddress); e != nil {
		return nil, errors.Errorf("failed to dial, error=%v", e)
	}
	defer s.client.Close()

	response := []*agentGrpc.ExecResponseMsg{}

	c := agentGrpc.NewCmdServiceClient(s.client)
	stream, err := c.ExecRemoteCmd(ctx)
	if err != nil {
		return nil, err
	}

	for n, request := range requestList {
		logger.Sugar().Debugf("send %v request ", n)
		if err := stream.Send(request); err != nil {
			return nil, err
		}

		if r, err := stream.Recv(); err != nil {
			return nil, err
		} else {
			logger.Sugar().Debugf("recv %v response ", n)
			response = append(response, r)
		}
	}

	logger.Debug("finish")
	if e := stream.CloseSend(); e != nil {
		logger.Sugar().Errorf("grpc failed to CloseSend error=%v ", e)
	}
	return response, nil

}

func (s *grpcClientManager) GetFileList(ctx context.Context, serverAddress, directory string) ([]string, error) {
	// get agent files list
	request := &agentGrpc.ExecRequestMsg{
		Timeoutsecond: 10,
		Command:       fmt.Sprintf("ls %v", directory),
	}
	responseList, e := s.SendRequestForExecRequest(ctx, []string{serverAddress}, []*agentGrpc.ExecRequestMsg{request})
	if e != nil || len(responseList) == 0 {
		return nil, fmt.Errorf("failed to get file list under directory %v of %v, error=%v", directory, serverAddress, e)
	}
	if responseList[0].Code != 0 {
		return nil, fmt.Errorf("failed to get file list under directory %v of %v, exit code=%v, stderr=%v", directory, serverAddress, responseList[0].Code, responseList[0].Stderr)
	}

	return strings.Fields(responseList[0].Stdmsg), nil
}

func (s *grpcClientManager) SaveRemoteFileToLocal(ctx context.Context, serverAddress, remoteFilePath, localFilePath string) error {

	// get agent files list
	request := &agentGrpc.ExecRequestMsg{
		Timeoutsecond: 10,
		Command:       fmt.Sprintf("cat %v", remoteFilePath),
	}
	responseList, e := s.SendRequestForExecRequest(ctx, []string{serverAddress}, []*agentGrpc.ExecRequestMsg{request})
	if e != nil || len(responseList) == 0 {
		return fmt.Errorf("failed to get remote file %v of %v, error=%v", remoteFilePath, serverAddress, e)
	}
	if responseList[0].Code != 0 {
		return fmt.Errorf("failed to get remote file %v of %v, exit code=%v, stderr=%v", remoteFilePath, serverAddress, responseList[0].Code, responseList[0].Stderr)
	}

	if len(responseList[0].Stdmsg) == 0 {
		return fmt.Errorf("got empty remote file %v of %v ", remoteFilePath, serverAddress)
	}

	f, e := os.OpenFile(localFilePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if e != nil {
		return fmt.Errorf("open file %v failed, error=%v", localFilePath, e)
	}
	defer f.Close()
	// _, e = f.Write([]byte(response.Stdmsg))
	// if e != nil {
	// 	return fmt.Errorf("failed to write file %v, error=%v", localFilePath, e)
	// }

	writer := bufio.NewWriter(f)
	_, e = writer.WriteString(responseList[0].Stdmsg)
	if e != nil {
		return fmt.Errorf("failed to write file %v, error=%v", localFilePath, e)
	}
	if e := writer.Flush(); e != nil {
		return fmt.Errorf("failed to flush file %v, error=%v", localFilePath, e)
	}

	return nil
}
