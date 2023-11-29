# kdoctor
[![Auto Release Version](https://github.com/kdoctor-io/kdoctor/actions/workflows/auto-release.yaml/badge.svg)](https://github.com/kdoctor-io/kdoctor/actions/workflows/auto-release.yaml)
[![Auto Nightly CI](https://github.com/kdoctor-io/kdoctor/actions/workflows/auto-nightly-ci.yaml/badge.svg)](https://github.com/kdoctor-io/kdoctor/actions/workflows/auto-nightly-ci.yaml)
[![codecov](https://codecov.io/gh/kdoctor-io/kdoctor/branch/main/graph/badge.svg?token=rLmsuiBLM2)](https://codecov.io/gh/kdoctor-io/kdoctor)
[![Go Report Card](https://goreportcard.com/badge/github.com/kdoctor-io/kdoctor)](https://goreportcard.com/report/github.com/kdoctor-io/kdoctor)
[![CodeFactor](https://www.codefactor.io/repository/github/kdoctor-io/kdoctor/badge)](https://www.codefactor.io/repository/github/kdoctor-io/kdoctor)
![badge](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/ii2day/0300d0a99d701fec02909d843792e67d/raw/e2ereport.json)

***

**English** | [**简体中文**](./README-zh_CN.md)

## Introduction

kdoctor is a Kubernetes data plane testing component that conducts functional and performance tests on clusters using proactive pressure injection. It addresses the operational needs of network, storage, and applications by adopting a cloud-native approach based on extensive research and abstraction. With its CRD design, kdoctor can seamlessly integrate with observability components.

**kdoctor mainly offers three types of tasks:**
* [AppHttpHealthy](./docs/reference/apphttphealthy.md): according to the task configuration, perform connectivity checks using HTTP and HTTPS protocols on specified addresses within or outside the cluster, supporting various request methods such as PUT, GET, and POST.
* [NetReach](./docs/reference/netreach.md): conduct connectivity inspections on Pod IP, ClusterIP, NodePort, LoadBalancer IP, Ingress IP, and even Pods with multiple network interfaces or dual-stack IPs.
* [NetDns](./docs/reference/netdns.md): perform connectivity checks on designated DNS servers within or outside the cluster, supporting UDP, TCP, and TCP-TLS protocols.

**Advantages of kdoctor over traditional testing components:**
* By configuring inspection tasks through CRDs, users only need to focus on the inspection targets, frequency, pressure parameters, and expected results.
* Pressure-injecting agents are dynamically run as Deployments or DaemonSets, achieving the effect of multiple pressure-injecting machines.
* The execution of tasks utilizes default agents or newly created agents based on the task's spec configurations, enabling resource reuse and task resource isolation.
* Agents are bound to corresponding resource targets such as ingress and service. Each agent Pod mutually accesses the bound resources according to the task configuration, deriving conclusions from the request results.
* Through performance optimization, the pressure-injecting client significantly reduces resource consumption during requests.
* Inspection reports can be generated through various means, including logging, aggregated APIs, and file storage.

## Architecture

<div style="text-align:center">
  <img src="./docs/images/arch.png" alt="Your Image Description">
</div>

Components:
* kdoctor agent: kdoctor controller: a persistent Deployment responsible for CR monitoring, task creation, and task report aggregation.
* kdoctor agent: dynamically created on-demand as Deployments or DaemonSets to execute tasks.

## Quick Start

**Install**
* Refer to [Install kdoctor](./docs/usage/install.md) 或 [kind Quick Start](./docs/usage//get-started-kind.md)

**Task Get Started**
* [AppHttpHealthy Get Started](./docs/usage/apphttphealthy.md)
* [NetReach Get Started](./docs/usage/netreach.md)
* [NetDNS Get Started](./docs/usage/netdns.md)

## Contribution

Refer to the [Contribution doc](./docs/develop/contributing.md).

## License

kdoctor is licensed under the Apache License, Version 2.0. See [LICENSE](./LICENSE) for the full license text.
