// Copyright 2023 Authors of kdoctor-io
// SPDX-License-Identifier: Apache-2.0
package types

type ConfigmapConfig struct {
	EnableIPv4                                 bool   `yaml:"enableIPv4"`
	EnableIPv6                                 bool   `yaml:"enableIPv6"`
	TaskPollIntervalInSecond                   int    `yaml:"taskPollIntervalInSecond"`
	NethttpDefaultRequestQps                   int    `yaml:"nethttp_defaultRequest_Qps"`
	NethttpDefaultRequestMaxQps                int    `yaml:"nethttp_defaultRequest_MaxQps"`
	NethttpDefaultConcurrency                  int    `yaml:"nethttp_defaultConcurrency"`
	NethttpDefaultMaxIdleConnsPerHost          int    `yaml:"nethttp_defaultMaxIdleConnsPerHost"`
	NethttpDefaultRequestDurationInSecond      int    `yaml:"nethttp_defaultRequest_DurationInSecond"`
	NethttpDefaultRequestPerRequestTimeoutInMS int    `yaml:"nethttp_defaultRequest_PerRequestTimeoutInMS"`
	NetdnsDefaultConcurrency                   int    `yaml:"netdns_defaultConcurrency"`
	MultusPodAnnotationKey                     string `yaml:"multusPodAnnotationKey"`
	CrdMaxHistory                              int    `yaml:"crdMaxHistory"`

	AgentSerivceIpv4Name string `yaml:"agentSerivceIpv4Name"`
	AgentSerivceIpv6Name string `yaml:"agentSerivceIpv6Name"`
	AgentIngressName     string `yaml:"agentIngressName"`
	AgentDaemonsetName   string `yaml:"agentDaemonsetName"`

	KdoctorAgent KdoctorAgentConfig `yaml:"kdoctorAgent"`
}

type KdoctorAgentConfig struct {
	UniqueMatchLabelKey string `json:"uniqueMatchLabelKey"`
}

type EnvMapping struct {
	EnvName      string
	DefaultValue string
	P            interface{}
}
