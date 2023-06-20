// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package pluginManager

import (
	"github.com/kdoctor-io/kdoctor/pkg/lock"
	"github.com/kdoctor-io/kdoctor/pkg/pluginManager/apphttphealthy"
	"github.com/kdoctor-io/kdoctor/pkg/pluginManager/netdns"
	"github.com/kdoctor-io/kdoctor/pkg/pluginManager/netreach"
	plugintypes "github.com/kdoctor-io/kdoctor/pkg/pluginManager/types"
	"go.uber.org/zap"
)

var pluginLock = &lock.Mutex{}

type pluginManager struct {
	chainingPlugins map[string]plugintypes.ChainingPlugin
	logger          *zap.Logger
}
type PluginManager interface {
	RunAgentController()
	RunControllerController(healthPort int, webhookPort int, webhookTlsDir string)
}

var globalPluginManager *pluginManager

// -------------------------

func InitPluginManager(logger *zap.Logger) PluginManager {
	pluginLock.Lock()
	defer pluginLock.Unlock()

	globalPluginManager.logger = logger

	return globalPluginManager
}

const (
	// ------ add crd ------
	KindNameAppHttpHealthy = "AppHttpHealthy"
	KindNameNetReach       = "NetReach"
	KindNameNetdns         = "Netdns"
)

var TaskKinds = []string{KindNameAppHttpHealthy, KindNameNetReach, KindNameNetdns}

func init() {
	globalPluginManager = &pluginManager{
		chainingPlugins: map[string]plugintypes.ChainingPlugin{},
	}

	// ------ add crd ------
	globalPluginManager.chainingPlugins[KindNameAppHttpHealthy] = &apphttphealthy.PluginAppHttpHealthy{}
	globalPluginManager.chainingPlugins[KindNameNetReach] = &netreach.PluginNetReach{}
	globalPluginManager.chainingPlugins[KindNameNetdns] = &netdns.PluginNetDns{}

}
