// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package types

var AgentEnvMapping = []EnvMapping{
	{"ENV_ENABLED_METRIC", "false", &AgentConfig.EnableMetric},
	{"ENV_METRIC_HTTP_PORT", "", &AgentConfig.MetricPort},
	{"ENV_HTTP_PORT", "80", &AgentConfig.HttpPort},
	{"ENV_GOPS_LISTEN_PORT", "", &AgentConfig.GopsPort},
	{"ENV_WEBHOOK_PORT", "", &AgentConfig.WebhookPort},
	{"ENV_PYROSCOPE_PUSH_SERVER_ADDRESS", "", &AgentConfig.PyroscopeServerAddress},
	{"ENV_POD_NAME", "", &AgentConfig.PodName},
	{"ENV_POD_NAMESPACE", "", &AgentConfig.PodNamespace},
	{"ENV_GOLANG_MAXPROCS", "8", &AgentConfig.GolangMaxProcs},
	{"ENV_AGENT_GRPC_LISTEN_PORT", "3000", &AgentConfig.AgentGrpcListenPort},
	{"ENV_ENABLE_AGGREGATE_AGENT_REPORT", "false", &AgentConfig.EnableAggregateAgentReport},
	{"ENV_AGENT_REPORT_STORAGE_PATH", "", &AgentConfig.DirPathAgentReport},
	{"ENV_CLEAN_AGED_REPORT_INTERVAL_IN_MINUTE", "10", &AgentConfig.CleanAgedReportInMinute},
	{"ENV_CLUSTER_DNS_DOMAIN", "cluster.local", &AgentConfig.ClusterDnsDomain},
	{"ENV_LOCAL_NODE_IP", "", &AgentConfig.LocalNodeIP},
	{"ENV_LOCAL_NODE_NAME", "", &AgentConfig.LocalNodeName},
}

type AgentConfigStruct struct {
	// ------- from env
	EnableMetric           bool
	MetricPort             int32
	HttpPort               int32
	GopsPort               int32
	WebhookPort            int32
	AgentGrpcListenPort    int32
	PyroscopeServerAddress string
	GolangMaxProcs         int32

	PodName          string
	PodNamespace     string
	ClusterDnsDomain string
	LocalNodeIP      string
	LocalNodeName    string

	EnableAggregateAgentReport bool
	DirPathAgentReport         string
	CleanAgedReportInMinute    int32

	// ------- from flags
	ConfigMapPath     string
	TlsCaCertPath     string
	TlsServerCertPath string
	TlsServerKeyPath  string
	TaskKind          string
	TaskName          string

	AppMode bool

	// from configmap
	Configmap ConfigmapConfig
}

var AgentConfig AgentConfigStruct
