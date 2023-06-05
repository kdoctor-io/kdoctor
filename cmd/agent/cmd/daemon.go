// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"github.com/kdoctor-io/kdoctor/pkg/debug"
	"github.com/kdoctor-io/kdoctor/pkg/pluginManager"
	"github.com/kdoctor-io/kdoctor/pkg/types"
)

func SetupUtility() {
	// run gops
	d := debug.New(rootLogger)
	if types.AgentConfig.GopsPort != 0 {
		d.RunGops(int(types.AgentConfig.GopsPort))
	}

	if types.AgentConfig.PyroscopeServerAddress != "" {
		d.RunPyroscope(types.AgentConfig.PyroscopeServerAddress, types.AgentConfig.PodName)
	}
}

func DaemonMain() {
	rootLogger.Sugar().Infof("config: %+v", types.AgentConfig)

	SetupHttpServer()
	initGrpcServer()
	// TODO: udp server, tcp server, websocket server

	if types.AgentConfig.AppMode {
		// app mode, just used to debug
		rootLogger.Info("run in app mode")
		// sleep forever
		select {}
	} else {
		rootLogger.Info("run in agent mode")

		SetupUtility()

		RunMetricsServer(types.AgentConfig.PodName)

		s := pluginManager.InitPluginManager(rootLogger.Named("agentContorller"))
		s.RunAgentController()

		rootLogger.Info("finish initialization")
		// sleep forever
		select {}
	}

}
