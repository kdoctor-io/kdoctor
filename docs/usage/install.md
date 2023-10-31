# install

[**简体中文**](./install-zh_CN.md) | **English**

## Introduction

Install kdoctor to check network and performance inside and outside the cluster.

## Implementation requirements

1. A complete k8s cluster.

2. installed [Helm](https://helm.sh/docs/intro/install/)

3. storageClass (optional) if kdoctor-controller is required for high availability and report persistence is required.

## Installation

### Add helm repository

```shell
helm repo add kdoctor https://kdoctor-io.github.io/kdoctor
helm repo update kdoctor
```

### Installing kdoctor

kdoctor can be installed according to different needs, the following are the recommended installation methods for several scenarios

#### 1. Non-highly available installation

The kdoctor agent only prints reports to standard output in the following way

```shell
helm install kdoctor kdoctor/kdoctor \
    -n kdoctor --debug --create-namespace 
```

#### 2. Highly available installation

The following method directs the collection reports from kdoctor-controller to storage, so you need to install storageClass.

```shell

helm install kdoctor kdoctor/kdoctor \
    -n kdoctor --debug --create-namespace \
    --set kdoctorController.replicas=2 \
    --set feature.aggregateReport.controller.pvc.enabled=true \
    --set feature.aggregateReport.controller.pvc.storageClass=local-path \
    --set feature.aggregateReport.controller.pvc.storageRequests="100Mi" \
    --set feature.aggregateReport.controller.pvc.storageLimits="500Mi"
ðŸ "ðŸ "ðŸ "ðŸ "ðŸ "ðŸ "ñ
```

### Verify that all components of kdoctor are running properly

```shell
kubectl get pod -n kdoctor
NAME READY STATUS RESTARTS AGE
kdoctor-controller-686b75d6d7-k4dcq 1/1 Running 0 137m
```

### Uninstall kdoctor

```shell
helm uninstall kdoctor -n kdoctor
```
