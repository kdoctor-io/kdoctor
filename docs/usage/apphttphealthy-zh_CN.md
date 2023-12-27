# AppHttpHealthy

[**English**](./apphttphealthy.md) | **简体中文**

## 介绍

对于这种任务， kdoctor-controller 会根据 agentSpec 生成对应的 [agent](../concepts/runtime-zh_CN.md) 等资源 ，每一个 agent Pod 都会向指定的 DNS server 发送 DNS 请求，默认并发量为 50 可覆盖多副本情况，并发量可在 kdoctor 的 configmap 中设置，并获得成功率和平均延迟。根据成功条件来判断结果是否成功。并且，可以通过聚合 API 获取详细的报告。

1. 应用场景：

    * 测试连通性，确认指定应用能够被集群的每一个角落访问到
    * 大规模集群测试，模拟更多的 client 数量，以能够产生更大的压力，测试应用的抗压能力，模拟更多的源 IP 来产生更多的应用会话，测试应用的资源限制。
    * 给指定应用注入压力，配合灰度发布、混沌测试、bug 复现等目的
    * 测试集群外部服务，确认集群 egress 工作正常

2. 关于 AppHttpHealthy CRD 的更多描述，可参考[AppHttpHealthy](../reference/apphttphealthy-zh_CN.md)

3. 功能列表:

      * 支持 HTTP、HTTPS、HTTP2，能够定义 header、body

## 使用步骤

接下来将展示 `AppHttpHealthy` 的使用示例

### 安装 kdoctor

参照[安装教程](./install-zh_CN.md)安装 kdoctor

### 安装测试 server (选做)

kdoctor 官方仓库中包含了一个名为 server 的应用，内包含 http server，https server, DNS server，可用来测试 kdoctor 功能，若存在其他测试的 server 可跳过安装。

```shell
helm repo add kdoctor https://kdoctor-io.github.io/kdoctor
helm repo update kdoctor
helm install server kdoctor/server -n kdoctor --wait --debug --create-namespace 
```

查看测试 server 状态

```shell
kubectl get pod -n kdoctor -owide
NAME                                READY   STATUS    RESTARTS   AGE   IP            NODE                    NOMINATED NODE   READINESS GATES
server-7649566ff9-dv4jc   1/1     Running   0          76s   172.40.1.45   kdoctor-worker          <none>           <none>
server-7649566ff9-qc5dh   1/1     Running   0          76s   172.40.0.35   kdoctor-control-plane   <none>           <none>
```

获取测试 server 的 service 地址

```shell
kubectl get service -n kdoctor
NAME               TYPE        CLUSTER-IP    EXTERNAL-IP   PORT(S)                                AGE
server   ClusterIP   172.41.71.0   <none>        80/TCP,443/TCP,53/UDP,53/TCP,853/TCP   2m31s
```

### 创建 AppHttpHealthy

创建 http `AppHttpHealthy` ，该任务将执行一轮持续 10s 的任务，任务会向指定的 server 以 QPS 为 10 的速度发送 Get 请求，并且立即执行。

这里使用 server 的 service 地址，若有其他 server 地址 可使用其他 server 地址。

```shell
SERVER="172.41.71.0"
cat <<EOF | kubectl apply -f -
apiVersion: kdoctor.io/v1beta1
kind: AppHttpHealthy
metadata:
  name: http1
spec:
  request:
    durationInSecond: 10
    perRequestTimeoutInMS: 1000
    qps: 10
  schedule:
    roundNumber: 1
    roundTimeoutMinute: 1
    schedule: 0 1
  expect:
    meanAccessDelayInMs: 1000
    successRate: 1
  target:
    host: http://${SERVER}
    method: GET
EOF
```

### 查看任务状态

当执行完成一轮后就可以使用 kdoctor 聚合 api 查看当前轮的报告，当 FINISH 为 true 时任务全部完成，可查看整体报告

```shell
kubectl get apphttphealthy
NAME        FINISH   EXPECTEDROUND   DONEROUND   LASTROUNDSTATUS   SCHEDULE
http        true     1               1           succeed           0 1
```

* FINISH：任务是否完成
* EXPECTEDROUND：希望任务执行轮数
* DONEROUND：当前执行完成轮数
* LASTROUNDSTATUS：最后一轮任务执行情况
* SCHEDULE：任务的调度规则

### 查看任务报告

1. 查看已有报告

    ```shell
    kubectl get kdoctorreport
    NAME        CREATED AT
    http        0001-01-01T00:00:00Z
    ```

2. 查看具体任务报告

    节点 kdoctor-control-plane 和节点 kdoctor-worker 上 agent 分别都执行一轮发压后，将 agent 报告聚合而成，报告名称由`${TaskKind}-${TaskName}`组成

    ```shell
    kubectl get kdoctorreport apphttphealthy-http -oyaml
    apiVersion: system.kdoctor.io/v1beta1
    kind: KdoctorReport
    metadata:
      creationTimestamp: null
      name: apphttphealthy-http
    spec:
      FailedRoundNumber: null
      FinishedRoundNumber: 1
      Report:
      - NodeName: kdoctor-control-plane
        HttpAppHealthyTask:
          Detail:
          - MeanDelay: 10.317307
            Metrics:
              Duration: 11.022081662s
              EndTime: "2023-07-31T07:25:23Z"
              Errors: {}
              Latencies:
                MaxInMs: 0
                MeanInMs: 10.317307
                MinInMs: 0
                P50InMs: 0
                P90InMs: 0
                P95InMs: 0
                P99InMs: 0
              RequestCounts: 104
              StartTime: "2023-07-31T07:25:12Z"
              StatusCodes:
                "200": 104
              SuccessCounts: 104
              TPS: 9.435604197939574
              TotalDataSize: 40040 byte
            Succeed: true
            SucceedRate: 1
            TargetMethod: GET
            TargetName: HttpAppHealthy target
            TargetUrl: http://172.41.71.0
          Succeed: true
          TargetNumber: 1
          TargetType: HttpAppHealthy
          MaxCPU: 30.651%
          MaxMemory: 97.00MB
        HttpAppHealthyTaskSpec:
        ...
        PodName: kdoctor-agent-fmr9m
        ReportType: agent test report
        RoundDuration: 11.038965547s
        RoundNumber: 1
        RoundResult: succeed
        StartTimeStamp: "2023-07-31T07:25:12Z"
        EndTimeStamp: "2023-07-31T07:25:23Z"
        TaskName: apphttphealthy.http
        TaskType: AppHttpHealthy
      - NodeName: kdoctor-worker
        HttpAppHealthyTask:
          Detail:
          - MeanDelay: 10.548077
            Metrics:
              ...
            Succeed: true
            SucceedRate: 1
            TargetMethod: GET
            TargetName: HttpAppHealthy target
            TargetUrl: http://172.41.71.0
          Succeed: true
          TargetNumber: 1
          TargetType: HttpAppHealthy
        HttpAppHealthyTaskSpec:
        ...
        PodName: kdoctor-agent-s468h
        ReportType: agent test report
        RoundDuration: 11.034140236s
        RoundNumber: 1
        RoundResult: succeed
        StartTimeStamp: "2023-07-31T07:25:12Z"
        EndTimeStamp: "2023-07-31T07:25:23Z"
        TaskName: apphttphealthy.http
        TaskType: AppHttpHealthy
      ReportRoundNumber: 1
      RoundNumber: 1
      Status: Finished
      TaskName: http
      TaskType: AppHttpHealthy
    ```

> 若报告与预期结果不符合，可关注报告中的 MaxCPU和 MaxMemory 字段，对比 agent 资源是否充足，调整 agent 的资源限制。

## 其他常用示例

下面是携带 body 的 http 请求示例和 https 的请求示例

1. 创建带有 body 的 http `AppHttpHealthy`，该任务将执行一轮持续 10s 的任务，任务会向指定的 server 以 qps 为 10 的速度携带body 进行 Post 请求，并且立即执行。

    这里使用 server 的 service 地址，若有其他 server 地址 可使用其他 server 地址。

    创建测试 body 数据

    ```shell
    cat <<EOF | kubectl apply -f -
    apiVersion: v1
    kind: ConfigMap
    metadata:
      name: kdoctor-test-body
      namespace: kdoctor-test-server
    data:
      test1: test1
      test2: test2
    EOF
    kubectl apply -f http-body.yaml
    ```

    创建 `AppHttpHealthy`

    ```shell
    SERVER="172.41.71.0"
    cat <<EOF | kubectl apply -f -
    apiVersion: kdoctor.io/v1beta1
    kind: AppHttpHealthy
    metadata:
      name: http-body
    spec:
      request:
        durationInSecond: 10
        perRequestTimeoutInMS: 1000
        qps: 10
      schedule:
        roundNumber: 1
        roundTimeoutMinute: 1
        schedule: 0 1
      expect:
        meanAccessDelayInMs: 1000
        successRate: 1
        statusCode: 200
      target:
        bodyConfigmapName: kdoctor-test-body
        bodyConfigmapNamespace: kdoctor-test-server
        header:
         - "Content-Type: application/json"
        host: http://${SERVER}
        method: POST
    EOF
    ```

2. 创建 https `AppHttpHealthy` ，该任务将执行一轮持续 10s 的任务，任务会向指定的 server 以 qps 为 10 的速度使用 https 协议携带证书发送 Get 请求，并且立即执行

    此 TLS 证书由 server 生成，证书只对 Pod 的 IP 进行了签名，因此我们 server 的 Pod IP 进行访问，若使用其他 server 请自行创建证书 secret。

    ```shell
    SERVER="172.40.0.35"
    cat <<EOF | kubectl apply -f -
    apiVersion: kdoctor.io/v1beta1
    kind: AppHttpHealthy
    metadata:
      name: https
    spec:
      request:
        durationInSecond: 10
        perRequestTimeoutInMS: 1000
        qps: 10
      schedule:
        roundNumber: 1
        roundTimeoutMinute: 1
        schedule: 0 1
      expect:
        meanAccessDelayInMs: 1000
        successRate: 1
        statusCode: 200
      target:
        host: https://${SERVER}
        method: GET
        tlsSecretName: https-client-cert
        tlsSecretNamespace: kdoctor-test-server
    EOF
    ```

## 环境清理

```shell
kubectl delete apphttphealthy http https http-body
```
