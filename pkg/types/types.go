// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0
package types

type ConfigmapConfig struct {
	EnableIPv4 bool `yaml:"enableIPv4"`
	EnableIPv6 bool `yaml:"enableIPv6"`

	TaskPollIntervalInSecond int `yaml:"taskPollIntervalInSecond"`
	// nethttp
	NetHttpDefaultRequestQPS                   int `yaml:"netHttpDefaultRequestQPS"`
	NetHttpDefaultRequestDurationInSecond      int `yaml:"netHttpDefaultRequestDurationInSecond"`
	NetHttpDefaultRequestPerRequestTimeoutInMS int `yaml:"netHttpDefaultRequestPerRequestTimeoutInMS"`

	// netreach
	NetReachRequestMaxQPS int `yaml:"netReachRequestMaxQPS"`
	// apphttphealthy
	AppHttpHealthyRequestMaxQPS int `yaml:"appHttpHealthyRequestMaxQPS"`
	// netdns
	NetDnsRequestMaxQPS int `yaml:"netDnsRequestMaxQPS"`

	MultusPodAnnotationKey string `yaml:"multusPodAnnotationKey"`
	CrdMaxHistory          int    `yaml:"crdMaxHistory"`

	AgentSerivceIpv4Name string `yaml:"agentSerivceIpv4Name"`
	AgentSerivceIpv6Name string `yaml:"agentSerivceIpv6Name"`
	AgentIngressName     string `yaml:"agentIngressName"`
	AgentDaemonsetName   string `yaml:"agentDaemonsetName"`

	AgentDefaultTerminationGracePeriodMinutes int64              `yaml:"agentDefaultTerminationGracePeriodMinutes"`
	KdoctorAgent                              KdoctorAgentConfig `yaml:"kdoctorAgent"`
}

type KdoctorAgentConfig struct {
	UniqueMatchLabelKey string `json:"uniqueMatchLabelKey"`
}

type EnvMapping struct {
	EnvName      string
	DefaultValue string
	P            interface{}
}
