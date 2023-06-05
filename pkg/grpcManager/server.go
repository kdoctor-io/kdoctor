// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package grpcManager

import (
	"context"
	"crypto/tls"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"math"
	"net"
	"time"
)

type GrpcServerManager interface {
	Run(listenAddress string)
	UpdateHealthStatus(status healthpb.HealthCheckResponse_ServingStatus)
}

type grpcServer struct {
	server             *grpc.Server
	healthcheckService *health.Server
	logger             *zap.Logger
}

const (
	DefaultServerKeepAliveTimeInterval      = 60 * time.Second
	DefaultServerKeepAliveTimeOut           = 10 * time.Second
	DefaultServerKeepAlivedAlowPeerInterval = 5 * time.Second
	// 100 MByte
	DefaultMaxRecvMsgSize = 1024 * 1024 * 100
	DefaultMaxSendMsgSize = math.MaxInt32
)

func NewGrpcServer(logger *zap.Logger, tlsCertPath, tlskeyPath string) GrpcServerManager {
	m := &grpcServer{}
	opts := []grpc.ServerOption{}

	m.logger = logger.Named("grpcManager")
	m.logger.Sugar().Infof("NewGrpcServer, tlsCertPath=%v, tlskeyPath=%v", tlsCertPath, tlskeyPath)

	// ----- tls
	cert, err := tls.LoadX509KeyPair(tlsCertPath, tlskeyPath)
	if err != nil {
		m.logger.Sugar().Fatalf("failed to load tls: %s", err)
	}
	opts = append(opts, grpc.Creds(credentials.NewServerTLSFromCert(&cert)))

	// https://godoc.org/google.golang.org/grpc/keepalive#EnforcementPolicy
	// Enforcement policy is a special setting on server side to protect server from malicious or misbehaving clients
	// for case : (1)Client sends too frequent pings (2)Client sends pings when there's no stream and this is disallowed by server config
	kaep := keepalive.EnforcementPolicy{
		MinTime:             DefaultServerKeepAlivedAlowPeerInterval,
		PermitWithoutStream: true,
	}
	opts = append(opts, grpc.KeepaliveEnforcementPolicy(kaep))

	// https://godoc.org/google.golang.org/grpc/keepalive#ServerParameters
	kasp := keepalive.ServerParameters{
		Time:    DefaultServerKeepAliveTimeInterval,
		Timeout: DefaultServerKeepAliveTimeOut,
	}
	// https://godoc.org/google.golang.org/grpc#KeepaliveParams
	opts = append(opts, grpc.KeepaliveParams(kasp))

	// https://godoc.org/google.golang.org/grpc#WithUnaryInterceptor
	// 添加 unary call 前的回调
	opts = append(opts, grpc.UnaryInterceptor(m.unaryInterceptor))
	// 对于 streaming RPCs , 注册 Interceptor  https://godoc.org/google.golang.org/grpc#StreamInterceptor
	opts = append(opts, grpc.StreamInterceptor(m.streamInterceptor))

	opts = append(opts, grpc.MaxRecvMsgSize(DefaultMaxRecvMsgSize))
	opts = append(opts, grpc.MaxSendMsgSize(DefaultMaxSendMsgSize))

	m.server = grpc.NewServer(opts...)
	if m.server == nil {
		m.logger.Fatal("failed to New Grpc Server ")
	}

	m.healthcheckService = health.NewServer()
	healthpb.RegisterHealthServer(m.server, m.healthcheckService)

	reflection.Register(m.server)

	m.registerService()

	return m
}

// address: "127.0.0.1:5000" or ":5000" (listen on ipv4 and ipv6)
func (t *grpcServer) Run(listenAddress string) {

	t.logger.Info("run grpc server at " + listenAddress)
	d, err := net.Listen("tcp", listenAddress)
	if err != nil {
		t.logger.Sugar().Fatalf("failed to listen: %v", err)
	}

	go func() {
		if err := t.server.Serve(d); err != nil {
			t.logger.Sugar().Fatalf("failed to run Grpc Server, reason=%v", err)
		}
	}()

}

// type HealthCheckResponse_ServingStatus int32
// const (
//
//	HealthCheckResponse_UNKNOWN         HealthCheckResponse_ServingStatus = 0
//	HealthCheckResponse_SERVING         HealthCheckResponse_ServingStatus = 1
//	HealthCheckResponse_NOT_SERVING     HealthCheckResponse_ServingStatus = 2
//	HealthCheckResponse_SERVICE_UNKNOWN HealthCheckResponse_ServingStatus = 3 // Used only by the Watch method.
//
// )
func (t *grpcServer) UpdateHealthStatus(status healthpb.HealthCheckResponse_ServingStatus) {
	if t.healthcheckService == nil {
		return
	}

	t.healthcheckService.SetServingStatus("", status)
	t.logger.Sugar().Infof("grpc server update health status to %v", status)

}

func (t *grpcServer) unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (i interface{}, e error) {
	// we can add more middleware heare

	// Continue execution of handler
	start := time.Now()
	i, e = handler(ctx, req)
	end := time.Now()
	t.logger.Sugar().Debugf("grpc server: rpc=%s , start_time=%s, end_time=%s, err=%v \n", info.FullMethod, start.Format(time.RFC3339Nano), end.Format(time.RFC3339Nano), e)

	return i, e
}

// 插入 每个 stream call 的回调
func (t *grpcServer) streamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (e error) {

	start := time.Now()
	e = handler(srv, ss)
	end := time.Now()
	t.logger.Sugar().Debugf("grpc server: rpc=%s , start_time=%s, end_time=%s, err=%v \n", info.FullMethod, start.Format(time.RFC3339Nano), end.Format(time.RFC3339Nano), e)

	return e
}
