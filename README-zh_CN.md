# kdoctor
[![Auto Release Version](https://github.com/kdoctor-io/kdoctor/actions/workflows/auto-release.yaml/badge.svg)](https://github.com/kdoctor-io/kdoctor/actions/workflows/auto-release.yaml)
[![Auto Nightly CI](https://github.com/kdoctor-io/kdoctor/actions/workflows/auto-nightly-ci.yaml/badge.svg)](https://github.com/kdoctor-io/kdoctor/actions/workflows/auto-nightly-ci.yaml)
[![codecov](https://codecov.io/gh/kdoctor-io/kdoctor/branch/main/graph/badge.svg?token=rLmsuiBLM2)](https://codecov.io/gh/kdoctor-io/kdoctor)
[![Go Report Card](https://goreportcard.com/badge/github.com/kdoctor-io/kdoctor)](https://goreportcard.com/report/github.com/kdoctor-io/kdoctor)
[![CodeFactor](https://www.codefactor.io/repository/github/kdoctor-io/kdoctor/badge)](https://www.codefactor.io/repository/github/kdoctor-io/kdoctor)

***

**简体中文** | [**English**](./README.md)

## Introduction

kdoctor 是一个 kubernetes 数据面测试项目，通过压力注入的方式，实现对集群进行功能、性能的主动式巡检。

传统的集群巡检，通过采集指标、日志、应用状态等信息来确认集群和应用的状态，实现被动式巡检。但是在一些特殊场景下，这种方式可能不能实现预期的巡检目的、时效性、集群范围，运维人员就需要采用手动方式给集群注入一些压力，进行主动式巡检，当集群规模很大、巡检频率高或巡检流程复杂时，手工方式难以持久实施。这些场景包括：

* 部署大规模集群后，希望确认所有节点间 POD 的网络连通性，避免某个节点存在网络故障，发现网络中是否存在偶发丢包问题，而通信渠道非常多，包括 pod IP、clusterIP、nodePort、loadbalancer ip、ingress ip, 甚至是 POD 多网卡、双栈IP

* 希望主动检测所有节点间上的 POD 能够正常访问 coredns 服务，希望确认 coredns 服务的资源配置和副本数量正确，其服务性能能欧支持预期的最大访问量

* 磁盘是易耗品，例如 etcd 等应用对磁盘性能是比较敏感的，在日常运维工作中，管理员希望周期地确认所有节点的本地磁盘是正常的，文件读写的吞吐量和延时是符合预期的

* 给某个服务主动注入压力，它可能是镜像仓库、mysql 或者 api-server，以配合 BUG 复现，或确认服务性能

kdoctor 是一个 kubernetes 数据面测试项目，来源于生产运维过程中的实践场景，通过压力注入的方式，实现对集群进行功能、性能的主动式巡检。 kdoctor 可以应用于:

* 生产环境的部署检查、日常运维等场景，能避免了人工巡检的工作负担。

* 能应用 E2E 测试、bug 复现、混沌测试等，减少编程工作。

## 架构

## 快速开始

## 核心功能

## License

Spiderpool is licensed under the Apache License, Version 2.0. See [LICENSE](./LICENSE) for the full license text.
