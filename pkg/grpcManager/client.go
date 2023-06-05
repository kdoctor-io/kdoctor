// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package grpcManager

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"github.com/kdoctor-io/kdoctor/api/v1/agentGrpc"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	"time"
)

const (
	DefaultDialTimeOut                 = 30 * time.Second
	DefaultClientKeepAliveTimeInterval = 30 * time.Second
	DefaultClientKeepAliveTimeOut      = 10 * time.Second

	LBPolicyFistPick = "pick_first"
	LBPolicyRR       = "round_robin"
)

type GrpcClientManager interface {
	SendRequestForExecRequest(ctx context.Context, serverAddress []string, requestList []*agentGrpc.ExecRequestMsg) ([]*agentGrpc.ExecResponseMsg, error)
	GetFileList(ctx context.Context, serverAddress, directory string) ([]string, error)
	SaveRemoteFileToLocal(ctx context.Context, serverAddress, remoteFilePath, localFilePath string) error
}

type grpcClientManager struct {
	logger *zap.Logger
	opts   []grpc.DialOption
	client *grpc.ClientConn
}

func NewGrpcClient(logger *zap.Logger, enableTls bool) GrpcClientManager {
	s := &grpcClientManager{
		logger: logger.Named("GrpcClientManager"),
	}

	s.opts = append(s.opts, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(DefaultMaxRecvMsgSize), grpc.MaxCallSendMsgSize(DefaultMaxRecvMsgSize)))
	s.opts = append(s.opts, grpc.WithBlock())

	if enableTls {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
		}
		s.opts = append(s.opts, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	} else {
		// no tls
		s.opts = append(s.opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	return s
}

// serverAddress:=[]string{"1.1.1.1:456"}
func (s *grpcClientManager) clientDial(ctx context.Context, serverAddress []string) error {
	opts := []grpc.DialOption{}
	s.logger.Sugar().Infof("grpc dial for %+v", serverAddress)

	addr := serverAddress[0]

	// --------
	serviceConfig := map[string]interface{}{}
	serviceConfig["LoadBalancingConfig"] = []map[string]map[string]string{
		{LBPolicyRR: {}},
	}
	serviceConfig["healthCheckConfig"] = map[string]string{
		// An empty string (`""`) typically indicates the overall health of a server should be reported
		"serviceName": "",
	}
	if jsongByte, e := json.Marshal(serviceConfig); e != nil {
		s.logger.Sugar().Fatalf("failed to parse serviceConfig, error=%v", e)
	} else {
		s.logger.Sugar().Debugf("grpc client serviceConfig = %+v \n ", string(jsongByte))
		opts = append(opts, grpc.WithDefaultServiceConfig(string(jsongByte)))
	}

	// --------
	kacp := keepalive.ClientParameters{
		Time:                DefaultClientKeepAliveTimeInterval, // send pings every 10 seconds if there is no activity
		Timeout:             DefaultClientKeepAliveTimeOut,
		PermitWithoutStream: false, // send pings even without active streams
	}
	opts = append(opts, grpc.WithKeepaliveParams(kacp))

	opts = append(opts, s.opts...)

	if c, err := grpc.DialContext(ctx, addr, opts...); err != nil {
		return errors.Errorf("grpc failed to dial")
	} else {
		s.client = c
	}
	return nil
}
