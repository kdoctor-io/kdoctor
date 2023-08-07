# 安装文档

[**English**](./install.md) | **简体中文**

## 介绍

安装 kdoctor 对集群内外的网络及性能进行检查

## 实施要求

1.一套完整的 k8s 集群

2.已安装 [Helm](https://helm.sh/docs/intro/install/)

3.storageClass(可选) 如果需要 kdoctor-controller 高可用且需要报告持久化

## 安装

### 添加 helm 仓库

```shell
helm repo add kdoctor https://kdoctor-io.github.io/kdoctor
helm repo update kdoctor
```

### 安装 kdoctor
kdoctor 可以根据不同的需求进行安装，以下为几个场景的推荐安装方式

#### 1.非高可用安装

以下方法 kdoctor agent 只将报告打印到标准输出
```shell 
helm install kdoctor kdoctor/kdoctor \
    -n kdoctor --debug --create-namespace 
```
#### 2.高可用安装

以下方法将 kdoctor-controller 的收集报告引导到存储，因此,需要安装storageClass

```shell 

helm  install kdoctor kdoctor/kdoctor \
    -n kdoctor --debug --create-namespace \
    --set kdoctorController.replicas=2 \
    --set feature.aggregateReport.controller.pvc.enabled=true \
    --set feature.aggregateReport.controller.pvc.storageClass=local-path  \
    --set feature.aggregateReport.controller.pvc.storageRequests="100Mi" \
    --set feature.aggregateReport.controller.pvc.storageLimits="500Mi"
```

### 确认 kdoctor 所有组件正常运行

```shell
kubectl get pod -n kdoctor
NAME                                  READY   STATUS    RESTARTS   AGE
kdoctor-agent-gp5mh                   1/1     Running   0          137m
kdoctor-agent-xkjn4                   1/1     Running   0          137m
kdoctor-controller-686b75d6d7-k4dcq   1/1     Running   0          137m
```

### 卸载 kdoctor

```shell
helm uninstall kdoctor -n kdoctor
```
