// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0
package types

const (
	TlsCaCommonName = "kdoctor.io"
)

const (
	KindNameAppHttpHealthy = "AppHttpHealthy"
	KindNameNetReach       = "NetReach"
	KindNameNetdns         = "Netdns"

	KindDeployment = "Deployment"
	KindDaemonSet  = "DaemonSet"
)

var TaskKinds = []string{KindNameAppHttpHealthy, KindNameNetReach, KindNameNetdns}
var TaskRuntimes = []string{KindDeployment, KindDaemonSet}
