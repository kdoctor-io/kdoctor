// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"github.com/kdoctor-io/kdoctor/pkg/grpcManager"
	"github.com/kdoctor-io/kdoctor/pkg/types"
)

func initGrpcServer() {
	// ---- grpc server
	rootLogger.Info("start grpc server")

	t := grpcManager.NewGrpcServer(rootLogger, TlsCertPath, TlsKeyPath)
	listenAddr := fmt.Sprintf(":%d", types.AgentConfig.AgentGrpcListenPort)
	t.Run(listenAddr)
}
