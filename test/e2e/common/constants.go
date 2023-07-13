// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package common

import (
	"os"
	"time"
)

const (
	PluginReportPath = "/apis/system.kdoctor.io/v1beta1/namespaces/default/kdoctorreports/"

	KDoctorAgentDSName         = "kdoctor-agent"
	KDoctorCaName              = "kdoctor-ca"
	KdoctorTestTokenSecretName = "apiserver-token"
	ServiceAccountReadyTimeout = time.Second * 20
)

var (
	TlsClientName = "kdoctor-client-cert"

	// get from env
	KdoctroTestToken   = ""
	APISERVICEADDR     = ""
	TestNameSpace      = "kube-system"
	AgentImageRegistry = "ghcr.io/kdoctor-io/kdoctor-agent"
	AgentImageTag      = "v0.1.0"
	TestIPv4           = false
	TestIPv6           = false
	KindClusterName    = "kdoctor"
)

func init() {
	APISERVICEADDR = os.Getenv("APISERVER")
	TestNameSpace = os.Getenv("E2E_INSTALL_NAMESPACE")
	AgentImageTag = os.Getenv("GIT_COMMIT_VERSION")
	AgentImageRegistry = os.Getenv("AGENT_REGISTER")
	KindClusterName = os.Getenv("E2E_KIND_CLUSTER_NAME")
	TestIPv4 = os.Getenv("E2E_IPV4_ENABLED") == "true"
	TestIPv6 = os.Getenv("E2E_IPV6_ENABLED") == "true"
}
