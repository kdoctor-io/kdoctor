# kdoctor
[![Auto Release Version](https://github.com/kdoctor-io/kdoctor/actions/workflows/auto-release.yaml/badge.svg)](https://github.com/kdoctor-io/kdoctor/actions/workflows/auto-release.yaml)
[![Auto Nightly CI](https://github.com/kdoctor-io/kdoctor/actions/workflows/auto-nightly-ci.yaml/badge.svg)](https://github.com/kdoctor-io/kdoctor/actions/workflows/auto-nightly-ci.yaml)
[![codecov](https://codecov.io/gh/kdoctor-io/kdoctor/branch/main/graph/badge.svg?token=rLmsuiBLM2)](https://codecov.io/gh/kdoctor-io/kdoctor)
[![Go Report Card](https://goreportcard.com/badge/github.com/kdoctor-io/kdoctor)](https://goreportcard.com/report/github.com/kdoctor-io/kdoctor)
[![CodeFactor](https://www.codefactor.io/repository/github/kdoctor-io/kdoctor/badge)](https://www.codefactor.io/repository/github/kdoctor-io/kdoctor)
![badge](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/ii2day/0300d0a99d701fec02909d843792e67d/raw/e2ereport.json)

***

**简体中文** | [**English**](./README.md)

## 介绍

kdoctor 是一个基于主动式压力注入的 Kubernetes 数据面测试组件，对集群进行功能、性能的测试。通过调研和抽象运维人员的常见需求，kdoctor 将网络、存储、应用等运维任务以云原生的方式实现。此外，还采用了基于 CRD 的设计，能够对接观测性组件。

**kdoctor 主要包含以下 3 个类型的任务：**
* [AppHttpHealthy](./docs/reference/apphttphealthy-zh_CN.md): 根据任务配置，使用 HTTP、HTTPS 协议对集群内外指定访问地址进行连通性检查，支持 PUT、GET、POST 等多种请求方式。
* [NetReach](./docs/reference/netreach-zh_CN.md): 根据任务配置对集群内 Pod IP、ClusterIP、NodePort、Loadbalancer IP、Ingress IP, 甚至是 Pod 多网卡、双栈 IP进行连通性巡检。
* [NetDns](./docs/reference/netdns-zh_CN.md): 根据任务配置，对集群内外的指定 DNS Server 进行连通性检测，支持 UDP、TCP、TCP-TLS 协议。

**kdoctor 较传统的测试组件有哪些优势:**
* 通过下发 CRD 配置巡检任务需求，使用者只需要关注巡检目标、巡检频率、发压参数以及期望巡检结果。
* 通过读取任务配置，以 Deployment 或 DaemonSet 的方式运行发压 agent，以达到多台发压机器的效果。
* 根据任务的 spec 配置，使用 default agent 或创建新的 agent 执行任务，以达到资源重复利用和任务资源隔离。
* 绑定相对应的资源目标，如 ingress 、service，每一个 agent Pod 根据任务配置相互访问绑定的资源，根据请求结果得出结论。
* 发压 client 通过性能调优，大大降低了发压请求时的资源消耗。
* 巡检报告通过日志、聚合 API、文件落盘等方式输出。

## 架构

<div style="text-align:center">
  <img src="./docs/images/arch.png" alt="Your Image Description">
</div>

组件构成：
* kdoctor controller: 以 Deployment 形式常驻，实施 CR 监控，任务创建，任务报告汇聚等。
* kdoctor agent: 以 Deployment 或 DaemonSet 形式按需动态创建，任务的执行者。

## 快速开始

**安装**
* 参考[安装 kdoctor](./docs/usage/install-zh_CN.md) 或 [kind 快速开始](./docs/usage/get-started-kind-zh_CN.md)

**开始任务**
* [开始任务 AppHttpHealthy](./docs/usage/apphttphealthy-zh_CN.md)
* [开始任务 NetReach](./docs/usage/netreach-zh_CN.md)
* [开始任务 NetDNS](./docs/usage/netdns-zh_CN.md)

## 参与开发

可参考 [开发搭建文档](./docs/develop/contributing.md).

## License

kdoctor is licensed under the Apache License, Version 2.0. See [LICENSE](./LICENSE) for the full license text.
