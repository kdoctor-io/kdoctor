// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0

package types

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

var ControllerEnvMapping = []EnvMapping{
	{"ENV_ENABLED_METRIC", "false", &ControllerConfig.EnableMetric},
	{"ENV_METRIC_HTTP_PORT", "", &ControllerConfig.MetricPort},
	{"ENV_HTTP_PORT", "80", &ControllerConfig.HttpPort},
	{"ENV_GOPS_LISTEN_PORT", "", &ControllerConfig.GopsPort},
	{"ENV_WEBHOOK_PORT", "", &ControllerConfig.WebhookPort},
	{"ENV_PYROSCOPE_PUSH_SERVER_ADDRESS", "", &ControllerConfig.PyroscopeServerAddress},
	{"ENV_POD_NAME", "", &ControllerConfig.PodName},
	{"ENV_POD_NAMESPACE", "", &ControllerConfig.PodNamespace},
	{"ENV_GOLANG_MAXPROCS", "8", &ControllerConfig.GolangMaxProcs},
	{"ENV_AGENT_GRPC_LISTEN_PORT", "3000", &ControllerConfig.AgentGrpcListenPort},
	{"ENV_ENABLE_AGGREGATE_AGENT_REPORT", "false", &ControllerConfig.EnableAggregateAgentReport},
	{"ENV_CONTROLLER_REPORT_STORAGE_PATH", "/report", &ControllerConfig.DirPathControllerReport},
	{"ENV_AGENT_REPORT_STORAGE_PATH", "", &ControllerConfig.DirPathAgentReport},
	{"ENV_CLEAN_AGED_REPORT_INTERVAL_IN_MINUTE", "10", &ControllerConfig.CleanAgedReportInMinute},
	{"ENV_CONTROLLER_REPORT_AGE_IN_DAY", "30", &ControllerConfig.ReportAgeInDay},
	{"ENV_COLLECT_AGENT_REPORT_INTERVAL_IN_SECOND", "600", &ControllerConfig.CollectAgentReportIntervalInSecond},
	{"ENV_RESOURCE_TRACKER_CHANNEL_BUFFER", "500", &ControllerConfig.ResourceTrackerChannelBuffer},
	{"ENV_RESOURCE_TRACKER_MAX_DATABASE_CAP", "5000", &ControllerConfig.ResourceTrackerMaxDatabaseCap},
	{"ENV_RESOURCE_TRACKER_EXECUTOR_WORKERS", "3", &ControllerConfig.ResourceTrackerExecutorWorkers},
	{"ENV_RESOURCE_TRACKER_SIGNAL_TIMEOUT_SECONDS", "3", &ControllerConfig.ResourceTrackerSignalTimeoutSeconds},
	{"ENV_RESOURCE_TRACKER_TRACE_GAP_SECONDS", "5", &ControllerConfig.ResourceTrackerTraceGapSeconds},
}

type ControllerConfigStruct struct {
	// ------- from env
	EnableMetric           bool
	MetricPort             int32
	HttpPort               int32
	GopsPort               int32
	WebhookPort            int32
	AgentGrpcListenPort    int32
	PyroscopeServerAddress string
	GolangMaxProcs         int32

	PodName      string
	PodNamespace string

	EnableAggregateAgentReport         bool
	CleanAgedReportInMinute            int32
	DirPathControllerReport            string
	DirPathAgentReport                 string
	ReportAgeInDay                     int32
	CollectAgentReportIntervalInSecond int32

	ResourceTrackerChannelBuffer        int32
	ResourceTrackerMaxDatabaseCap       int32
	ResourceTrackerExecutorWorkers      int32
	ResourceTrackerSignalTimeoutSeconds int32
	ResourceTrackerTraceGapSeconds      int32

	// -------- from flags
	ConfigMapPath     string
	TlsCaCertPath     string
	TlsServerCertPath string
	TlsServerKeyPath  string

	ConfigMapDeploymentPath string
	ConfigMapDaemonsetPath  string
	ConfigMapPodPath        string
	ConfigMapServicePath    string

	// -------- from configmap
	Configmap ConfigmapConfig
}

var ControllerConfig ControllerConfigStruct

// singleton for Application templates from configmap
var (
	DeploymentTempl = new(appsv1.Deployment)
	DaemonsetTempl  = new(appsv1.DaemonSet)
	PodTempl        = new(corev1.Pod)
	ServiceTempl    = new(corev1.Service)
)
