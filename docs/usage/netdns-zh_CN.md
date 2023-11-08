# NetDns

[**English**](./netdns.md) | **简体中文**

## 介绍 

对于这种任务， kdoctor-controller 会根据 agentSpec 生成对应的 [agent](../concepts/runtime-zh_CN.md) 等资源 ，每一个 agent Pod 都会向指定的 DNS server 发送 DNS 请求，默认并发量为 50 可覆盖多副本情况，并发量可在 kdoctor 的 configmap 中设置，并获得成功率和平均延迟。根据成功条件来判断结果是否成功。并且，可以通过聚合 API 获取详细的报告。

1. 应用场景：
    * 生产或 E2E 环境下，检测集群每个角落可访问 CoreDNS 服务
    * 在应用部署阶段，用以配合调整 CoreDNS 的资源和副本数量，以确认能够支撑期望的访问压力
    * 给 CoreDNS 注入压力，配合 CoreDNS 升级测试、混沌测试、bug 复现等目的
    * 测试集群外部的 DNS 服务

2. 关于 NetDns CRD 的更多描述，可参考[NetDns](../reference/netdns-zh_CN.md)

3. 功能列表:

    * 支持集群内外 DNS server 测试
    * 支持 typeA 、typeAAAA 记录
    * 支持 UDP、TCP、TCP-TLS 协议

## 使用步骤

接下来将展示 `NetDNS` 的使用示例

### 安装 kdoctor 

参照[安装教程](./install-zh_CN.md)安装 kdoctor

### 安装测试 server (选做)

kdoctor 官方仓库中包含了一个名为 server 的应用，内包含 http server，https server, DNS server，可用来测试 kdoctor 功能，若存在其他测试的 server 可跳过安装。

```shell
helm repo add kdoctor https://kdoctor-io.github.io/kdoctor
helm repo update kdoctor
helm install server kdoctor/server -n kdoctor-test-server --wait --debug --create-namespace 
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

### 创建 NetDns

创建 `NetDns`，该任务将执行一轮持续 10s 的任务，任务会向集群内 DNS server 以 QPS 为 10 的速度使用 UDP 协议，请求解析 `kubernetes.default.svc.cluster.local` 域名的 typeA 记录，并且立即执行。

```shell
cat <<EOF | kubectl apply -f -
apiVersion: kdoctor.io/v1beta1
kind: Netdns
metadata:
  name: netdns-cluster
spec:
  expect:
    meanAccessDelayInMs: 1500
    successRate: 1
  request:
    domain: kubernetes.default.svc.cluster.local
    durationInSecond: 10
    perRequestTimeoutInMS: 1000
    protocol: udp
    qps: 10
  schedule:
    roundNumber: 1
    roundTimeoutMinute: 1
    schedule: 0 1
  target:
    targetDns:
      serviceName: kube-dns
      serviceNamespace: kube-system
      testIPv4: true
EOF
```

### 查看任务状态

当执行完成一轮后就可以使用 kdoctor 聚合 api 查看当前轮的报告，当 FINISH 为 true 时任务全部完成，可查看整体报告

```shell
kubectl get netdns
NAME             FINISH   EXPECTEDROUND   DONEROUND   LASTROUNDSTATUS   SCHEDULE
netdns-cluster   true     1               1           succeed           0 1
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
    NAME             CREATED AT
    netdns-cluster   0001-01-01T00:00:00Z
    ```

2. 查看具体任务报告

    节点 kdoctor-control-plane 和节点 kdoctor-worker 上 agent 分别都执行一轮发压后，将 agent 报告聚合而成。

    ```shell
    root@kdoctor-control-plane:/# kubectl get kdoctorreport netdns-cluster -oyaml
    apiVersion: system.kdoctor.io/v1beta1
    kind: KdoctorReport
    metadata:
      creationTimestamp: null
      name: netdns-cluster
    spec:
      FailedRoundNumber: null
      FinishedRoundNumber: 1
      Report:
      - NodeName: kdoctor-control-plane
        PodName: kdoctor-agent-ntp9l
        ReportType: agent test report
        RoundDuration: 11.025723086s
        RoundNumber: 1
        RoundResult: succeed
        StartTimeStamp: "2023-08-01T09:09:39Z"
        EndTimeStamp: "2023-08-01T09:09:50Z"
        TaskName: netdns.netdns-cluster
        TaskType: Netdns
        netDNSTask:
          detail:
          - FailureReason: null
            MeanDelay: 0.2970297
            Metrics:
              DNSMethod: udp
              DNSServer: 172.41.0.10:53
              Duration: 11.002666395s
              EndTime: "2023-08-01T09:09:50Z"
              Errors: {}
              FailedCounts: 0
              Latencies:
                Max_inMx: 0
                Mean_inMs: 0.2970297
                Min_inMs: 0
                P50_inMs: 0
                P90_inMs: 0
                P95_inMs: 0
                P99_inMs: 0
              ReplyCode:
                NOERROR: 101
              RequestCounts: 101
              StartTime: "2023-08-01T09:09:39Z"
              SuccessCounts: 101
              TPS: 9.179593052634765
              TargetDomain: kubernetes.default.svc.cluster.local.
            Succeed: true
            SucceedRate: 1
            TargetName: typeA_172.41.0.10:53_kubernetes.default.svc.cluster.local
            TargetProtocol: udp
            TargetServer: 172.41.0.10:53
          succeed: true
          targetNumber: 1
          targetType: kdoctor agent
          MaxCPU: 30.651%
          MaxMemory: 97.00MB
        netDNSTaskSpec:
          ...
      - NodeName: kdoctor-worker
        PodName: kdoctor-agent-krrnp
        ReportType: agent test report
        RoundDuration: 10.024533428s
        RoundNumber: 1
        RoundResult: succeed
        StartTimeStamp: "2023-08-01T09:09:39Z"
        EndTimeStamp: "2023-08-01T09:09:49Z"
        TaskName: netdns.netdns-cluster
        TaskType: Netdns
        netDNSTask:
          detail:
          - FailureReason: null
            MeanDelay: 0.58
            Metrics:
              ...
            Succeed: true
            SucceedRate: 1
            TargetName: typeA_172.41.0.10:53_kubernetes.default.svc.cluster.local
            TargetProtocol: udp
            TargetServer: 172.41.0.10:53
          succeed: true
          targetNumber: 1
          targetType: kdoctor agent
          MaxCPU: 30.651%
          MaxMemory: 97.00MB
        netDNSTaskSpec:
          ...
      ReportRoundNumber: 1
      RoundNumber: 1
      Status: Finished
      TaskName: netdns-cluster
      TaskType: Netdns
    ```

> 若报告与预期结果不符合，可关注报告中的 MaxCPU和 MaxMemory 字段，对比 agent 资源是否充足，调整 agent 的资源限制。

## 集群外 DNS server 测试

下面是携带 body 的 http 请求示例和 https 的请求示例：

1. 创建 `NetDns` 任务，该任务将执行一轮持续 10s 的任务，任务会向指定的 DNS server 以 QPS 为 10 的速度进行 UDP 请求 `kubernetes.default.svc.cluster.local` 域名的 typeAAAA，并且立即执行。

    这里使用 server 的 service 地址，若有其他 server 地址 可使用其他 server 地址。

    创建 `NetDns`

    ```shell
    SERVER="172.41.71.0"
    apiVersion: kdoctor.io/v1beta1
    kind: Netdns
    metadata:
      name: netdns- user
    spec:
      expect:
        meanAccessDelayInMs: 1500
        successRate: 1
      request:
        domain: kubernetes.default.svc.cluster.local
        durationInSecond: 10
        perRequestTimeoutInMS: 1000
        protocol: udp
        qps: 10
      schedule:
        roundNumber: 1
        roundTimeoutMinute: 1
        schedule: 0 1
      target:
        targetUser:
          port: 53
          server: ${SERVER}
    EOF
    ```

## 环境清理

```shell
kubectl delete netdns netdns-cluster netdns- user
```