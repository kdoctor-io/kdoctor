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

	for n, v := range types.ControllerEnvMapping {
		m := v.DefaultValue
		if t := viper.GetString(v.EnvName); len(t) > 0 {
			m = t
		}
		if len(m) > 0 {
			switch v.P.(type) {
			case *int32:
				if s, err := strconv.ParseInt(m, 10, 64); err == nil {
					r := types.ControllerEnvMapping[n].P.(*int32)
					*r = int32(s)
				} else {
					logger.Fatal("failed to parse env value of " + v.EnvName + " to int32, value=" + m)
				}
			case *string:
				r := types.ControllerEnvMapping[n].P.(*string)
				*r = m
			case *bool:
				if s, err := strconv.ParseBool(m); err == nil {
					r := types.ControllerEnvMapping[n].P.(*bool)
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
	globalFlag.StringVarP(&types.ControllerConfig.ConfigMapPath, "config-path", "C", "", "configmap file path")
	globalFlag.StringVarP(&types.ControllerConfig.TlsCaCertPath, "tls-ca-cert", "R", "", "ca file path")
	globalFlag.StringVarP(&types.ControllerConfig.TlsServerCertPath, "tls-server-cert", "T", "", "server cert file path")
	globalFlag.StringVarP(&types.ControllerConfig.TlsServerKeyPath, "tls-server-key", "Y", "", "server key file path")
	if e := viper.BindPFlags(globalFlag); e != nil {
		logger.Sugar().Fatalf("failed to BindPFlags, reason=%v", e)
	}
	printFlag := func() {
		logger.Info("config-path = " + types.ControllerConfig.ConfigMapPath)
		logger.Info("tls-ca-cert = " + types.ControllerConfig.TlsCaCertPath)
		logger.Info("tls-server-cert = " + types.ControllerConfig.TlsServerCertPath)
		logger.Info("tls-server-key = " + types.ControllerConfig.TlsServerKeyPath)

		// load configmap
		if len(types.ControllerConfig.ConfigMapPath) > 0 {
			configmapBytes, err := os.ReadFile(types.ControllerConfig.ConfigMapPath)
			if nil != err {
				logger.Sugar().Fatalf("failed to read configmap file %v, error: %v", types.ControllerConfig.ConfigMapPath, err)
			}
			if err := yaml.Unmarshal(configmapBytes, &types.ControllerConfig.Configmap); nil != err {
				logger.Sugar().Fatalf("failed to parse configmap data, error: %v", err)
			}
		}
	}
	cobra.OnInitialize(printFlag)

}
