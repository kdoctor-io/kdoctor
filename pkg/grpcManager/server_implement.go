// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package grpcManager

import (
	"context"
	"fmt"
	"github.com/kdoctor-io/kdoctor/api/v1/agentGrpc"
	"github.com/kdoctor-io/kdoctor/pkg/utils"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"time"
)

// ------ implement
type myGrpcServer struct {
	agentGrpc.UnimplementedCmdServiceServer
	logger *zap.Logger
}

func (s *myGrpcServer) ExecRemoteCmd(stream agentGrpc.CmdService_ExecRemoteCmdServer) error {
	logger := s.logger
	var finalError error

	handler := func(ctx context.Context, r *agentGrpc.ExecRequestMsg) (*agentGrpc.ExecResponseMsg, error) {
		if len(r.Command) == 0 {
			logger.Error("grpc server ExecRemoteCmd: got empty command \n")
			return nil, status.Error(codes.InvalidArgument, "request command is empty")
		}

		if r.Timeoutsecond == 0 {
			logger.Error("grpc server ExecRemoteCmd: got empty timeout \n")
			return nil, status.Error(codes.InvalidArgument, "request command is empty")
		}

		clientctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.Timeoutsecond)*time.Second)
		defer cancel()
		go func() {
			select {
			case <-clientctx.Done():
			case <-ctx.Done():
				cancel()
			}
		}()

		StdoutMsg, StderrMsg, exitedCode, e := utils.RunFrondendCmd(clientctx, r.Command, nil, "")

		logger.Sugar().Debugf("stderrMsg = %v", StderrMsg)
		logger.Sugar().Debugf("StdoutMsg = %v", StdoutMsg)
		logger.Sugar().Debugf("exitedCode = %v", exitedCode)
		logger.Sugar().Debugf("error = %v", e)

		w := &agentGrpc.ExecResponseMsg{
			Stdmsg: StdoutMsg,
			Stderr: StderrMsg,
			Code:   int32(exitedCode),
		}

		return w, nil
	}

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// client has finish sending all message
			break
		}
		if err != nil {
			c := fmt.Sprintf("recv error, %v", err)
			logger.Error(c)
			finalError = status.Error(codes.Unknown, c)
			break
		}

		re, e := handler(stream.Context(), req)
		if e != nil {
			finalError = e
			break
		}

		if e := stream.Send(re); e != nil {
			c := fmt.Sprintf("grpc server failed to send msg: %v", err)
			logger.Error(c)
			finalError = fmt.Errorf(c)
			break
		}
	}

	return finalError
}

// ------------
func (t *grpcServer) registerService() {
	agentGrpc.RegisterCmdServiceServer(t.server, &myGrpcServer{
		logger: t.logger,
	})
}
