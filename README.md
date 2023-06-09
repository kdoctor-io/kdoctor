# kdoctor
[![Auto Release Version](https://github.com/kdoctor-io/kdoctor/actions/workflows/auto-release.yaml/badge.svg)](https://github.com/kdoctor-io/kdoctor/actions/workflows/auto-release.yaml)
[![Auto Nightly CI](https://github.com/kdoctor-io/kdoctor/actions/workflows/auto-nightly-ci.yaml/badge.svg)](https://github.com/kdoctor-io/kdoctor/actions/workflows/auto-nightly-ci.yaml)
[![codecov](https://codecov.io/gh/kdoctor-io/kdoctor/branch/main/graph/badge.svg?token=rLmsuiBLM2)](https://codecov.io/gh/kdoctor-io/kdoctor)
[![Go Report Card](https://goreportcard.com/badge/github.com/kdoctor-io/kdoctor)](https://goreportcard.com/report/github.com/kdoctor-io/kdoctor)
[![CodeFactor](https://www.codefactor.io/repository/github/kdoctor-io/kdoctor/badge)](https://www.codefactor.io/repository/github/kdoctor-io/kdoctor)

***

**English** | [**简体中文**](./README-zh_CN.md)

## Introduction

kdoctor is a cloud native project of data plane test. Through the pressure injection, it realizes the active inspection for the function and performance of the cluster.

For the traditional operation and maintenance , the status of clusters and applications is confirmed by collecting information such as metrics, logs, and application status, 
which could be called passive inspection. However, in some special scenarios, this method may not meet the expected purpose, timeliness, or cluster range, 
administrators need to manually inject some pressure into the cluster and checkout the cluster status, which could be called active inspection. 
When the cluster scale is large, or the inspection frequency is high, or the inspection process is complicated, it is hard to implement  manually. These scenarios include:

* After deploying a large-scale cluster, administrators want to confirm the network connectivity between all nodes, to find out network failures on a certain 
    node, or occasional packet loss. In addition, there are many communication ways including POD IP, clusterIP, nodePort, loadbalancerIP, ingress, or even POD multiple network interface, dual-stack IP.

* It is desired to make sure that PODs on all nodes can access the coredns service, or the resource configuration and the replica number of the coredns are enough to support expected access pressure.

* Disks are consumables and applications like etcd are sensitive to disk performance. In daily maintenance, administrators want to periodically confirm that local disks performance of all nodes are normal.

* Actively inject pressure on a service like registry, mysql or api-server, to cooperate with BUG reproduce, or to confirm service performance

kdoctor is a cloud native project of data plane test, which is derived from practices of the production operation and maintenance. Through the pressure injection, it realizes the active inspection for the function and performance of the cluster. kdoctor can be applied to scenarios:

* inspection after creating new cluster, daily operation and maintenance. 

* E2E testing, bug reproduction, chaos testing.

## Architecture

## Quick Start

## Feature

## License

Spiderpool is licensed under the Apache License, Version 2.0. See [LICENSE](./LICENSE) for the full license text.

