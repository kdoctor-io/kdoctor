// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"github.com/kdoctor-io/kdoctor/pkg/logger"
	"github.com/kdoctor-io/kdoctor/pkg/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"os"
	"reflect"
	"strconv"
)

func init() {

	viper.AutomaticEnv()
	if t := viper.GetString("ENV_LOG_LEVEL"); len(t) > 0 {
		rootLogger = logger.NewStdoutLogger(t, BinName).Named(BinName)
	} else {
		rootLogger = logger.NewStdoutLogger("", BinName).Named(BinName)
	}

	logger := rootLogger.Named("config")
	// env built in the image
	if t := viper.GetString("ENV_VERSION"); len(t) > 0 {
		logger.Info("app version " + t)
	}
	if t := viper.GetString("ENV_GIT_COMMIT_VERSION"); len(t) > 0 {
		logger.Info("git commit version " + t)
	}
	if t := viper.GetString("ENV_GIT_COMMIT_TIMESTAMP"); len(t) > 0 {
		logger.Info("git commit timestamp " + t)
	}

	for n, v := range types.AgentEnvMapping {
		m := v.DefaultValue
		if t := viper.GetString(v.EnvName); len(t) > 0 {
			m = t
		}
		if len(m) > 0 {
			switch v.P.(type) {
			case *int32:
				if s, err := strconv.ParseInt(m, 10, 64); err == nil {
					r := types.AgentEnvMapping[n].P.(*int32)
					*r = int32(s)
				} else {
					logger.Fatal("failed to parse env value of " + v.EnvName + " to int32, value=" + m)
				}
			case *string:
				r := types.AgentEnvMapping[n].P.(*string)
				*r = m
			case *bool:
				if s, err := strconv.ParseBool(m); err == nil {
					r := types.AgentEnvMapping[n].P.(*bool)
					*r = s
				} else {
					logger.Fatal("failed to parse env value of " + v.EnvName + " to bool, value=" + m)
				}
			default:
				logger.Sugar().Fatal("unsupported type to parse %v, config type=%v ", v.EnvName, reflect.TypeOf(v.P))
			}
		}

		logger.Info(v.EnvName + " = " + m)
	}

	// command flags
	globalFlag := rootCmd.PersistentFlags()
	globalFlag.StringVarP(&types.AgentConfig.ConfigMapPath, "config-path", "C", "", "configmap file path")
	globalFlag.StringVar(&types.AgentConfig.TaskKind, "task-kind", "", "")
	globalFlag.StringVar(&types.AgentConfig.TaskName, "task-name", "", "")
	if err := rootCmd.MarkPersistentFlagRequired("task-kind"); nil != err {
		logger.Sugar().Fatalf("failed to mark persistentFlag 'task-kind' as required, error: %v", err)
	}
	if err := rootCmd.MarkPersistentFlagRequired("task-name"); nil != err {
		logger.Sugar().Fatalf("failed to mark persistentFlag 'task-name' as required, error: %v", err)
	}

	globalFlag.BoolVarP(&types.AgentConfig.AppMode, "app-mode", "A", false, "app mode")
	globalFlag.BoolVarP(&types.AgentConfig.TlsInsecure, "tls-insecure", "K", true, "skip verify tls")
	globalFlag.StringVarP(&types.AgentConfig.TlsCaCertPath, "tls-ca-cert", "R", "/etc/tls/ca.crt", "ca file path")
	globalFlag.StringVarP(&types.AgentConfig.TlsCaKeyPath, "tls-ca-key", "Y", "/etc/tls/ca.key", "ca key file path")
	if e := viper.BindPFlags(globalFlag); e != nil {
		logger.Sugar().Fatalf("failed to BindPFlags, reason=%v", e)
	}
	printFlag := func() {
		logger.Info("config-path = " + types.AgentConfig.ConfigMapPath)

		// load configmap
		if len(types.AgentConfig.ConfigMapPath) > 0 {
			configmapBytes, err := os.ReadFile(types.AgentConfig.ConfigMapPath)
			if nil != err {
				logger.Sugar().Fatalf("failed to read configmap file %v, error: %v", types.AgentConfig.ConfigMapPath, err)
			}
			if err := yaml.Unmarshal(configmapBytes, &types.AgentConfig.Configmap); nil != err {
				logger.Sugar().Fatalf("failed to parse configmap data, error: %v", err)
			}
		}

		logger.Info("task-kind = " + types.AgentConfig.TaskKind)
		logger.Info("task-name = " + types.AgentConfig.TaskName)
	}
	cobra.OnInitialize(printFlag)

}
