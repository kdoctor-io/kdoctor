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

kdoctor is a cloud native project of data plane test. Through pressure injection, it realizes active cluster inspections for functionality and performance.

Traditional cluster inspections rely on collecting metrics, logs, and application status to passively assess the state of the cluster and applications.
However, in certain situations, this approach may not achieve the desired inspection objectives in terms of timeliness, coverage, or effectiveness.
To overcome these limitations, operators often resort to manually injecting stress into the cluster for active inspections.
However, manual methods become cumbersome to sustain when dealing with large-scale clusters, frequent inspections, or complex processes. Some common scenarios include:

* After deploying a large-scale cluster, administrators want to confirm the network connectivity between all nodes, to find out network failures on a certain 
    node, or occasional packet loss. This involves considering various communication methods such as Pod IP, ClusterIP, NodePort, LoadBalancer IP, ingress IP, and even multiple NICs or dual-stack IPs.

* Active inspections are necessary to make sure that Pods on all nodes can access the CoreDNS service, and the resource configuration and the replica number of the CoreDNS are correct to support expected access pressure.

* Disks are consumables and applications like etcd are sensitive to disk performance. In daily maintenance, administrators want to periodically confirm that local disks performance of all nodes are normal.

* Actively inject pressure on a service like registry, mysql or api-server, to cooperate with bug reproduction, or to confirm service performance

kdoctor is a Kubernetes data plane testing project, which is derived from practices of the production operation and maintenance. Through the pressure injection, it realizes active cluster inspections for functionality and performance. kdoctor can be applied to the following scenarios:

* Support deployment checks and routine maintenance in production environments, significantly reducing manual inspection efforts.  

* Support E2E tests, bug reproduction, chaos tests with  little need for extensive programming efforts.

## Architecture

## Quick Start

## Feature

## License

kdoctor is licensed under the Apache License, Version 2.0. See [LICENSE](../LICENSE) for the full license text.
