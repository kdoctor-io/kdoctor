# NetReach

[**English**](./netreach.md) | **简体中文**

## 介绍 

对于这种任务，kdoctor-controller 会根据 agentSpec 生成对应的 [agent](../concepts/runtime-zh_CN.md) 等资源，每一个 agent pod 都会以一定的压力相互发送http请求，请求地址为每一个 agent 的 pod ip 、service ip、ingress ip 等等，并获得成功率和平均延迟。根据成功条件来判断结果是否成功。并且，可以通过聚合API获取详细的报告。

1.应用场景：
* 每1min心跳，监控集群内每个角落的连通性
* 大规模集群部署后，巡检集群每个角落的连通性
* 给集群所有角落、每种通信方式注入流量，配合 bug 复现等场景
* 生产或 E2E 环境下，巡检 CNI、Loadbalancer、kube-proxy、Ingress、Multus 等组件的异常。

2.关于 NetReach CRD 的更多描述，可参考[NetReach](../reference/netreach-zh_CN.md)

3.功能列表:
    
* 支持 ClusterIP、Endpoint、Ingress、NodePort、LoadBalancer、Multus 多网卡、IPv4 IPv6

## 开始

接下来将展示 `NetReach` 的使用示例

### 安装 kdoctor 

参照[安装教程](./install-zh_CN.md)安装 kdoctor

### 创建 NetReach 

创建 `NetReach` ，该任务将执行一轮持续 10s 的任务，每个节点的 agent 会相互使用 http 协议访问 ClusterIP、Endpoint、Ingress、NodePort、LoadBalancer 的 IPv4 地址，并立即执行。

```shell
cat <<EOF | kubectl apply -f -
apiVersion: kdoctor.io/v1beta1
kind: NetReach
metadata:
  name: task
spec:
  expect:
    meanAccessDelayInMs: 1500
    successRate: 1
  request:
    durationInSecond: 10
    perRequestTimeoutInMS: 1000
    qps: 10
  schedule:
    roundNumber: 1
    roundTimeoutMinute: 1
    schedule: 0 1
  target:
    clusterIP: true
    endpoint: true
    ingress: true
    ipv4: true
    loadBalancer: true
    multusInterface: false
    nodePort: true
EOF
```

### 查看任务状态

当执行完成一轮后就可以使用 kdoctor 聚合 api 查看当前轮的报告，当 FINISH 为 true 时任务全部完成，可查看整体报告

```shell
kubectl get netreach
NAME   FINISH   EXPECTEDROUND   DONEROUND   LASTROUNDSTATUS   SCHEDULE
task   true     1               1           succeed           0 1
```

* FINISH：任务是否完成
* EXPECTEDROUND：希望任务执行轮数
* DONEROUND：当前执行完成轮数
* LASTROUNDSTATUS：最后一轮任务执行情况
* SCHEDULE：任务的调度规则


### 查看任务报告

1.查看已有报告

```shell
kubectl get kdoctorreport
NAME   CREATED AT
task   0001-01-01T00:00:00Z
```

2.查看具体任务报告

节点 kdoctor-control-plane 和节点 kdoctor-worker 上 agent 分别都执行一轮互相发压后，将 agent 报告聚合而成。

```shell
root@kdoctor-control-plane:/# kubectl get kdoctorreport task -oyaml
apiVersion: system.kdoctor.io/v1beta1
kind: KdoctorReport
metadata:
  creationTimestamp: null
  name: task
spec:
  FailedRoundNumber: null
  FinishedRoundNumber: 1
  Report:
  - NodeName: kdoctor-control-plane
    NetReachTask:
      Detail:
      - TargetName: AgentLoadbalancerV4IP_172.18.0.51:80
        Metrics:
          Duration: 10.032286878s
          EndTime: "2023-08-01T08:37:06Z"
          Errors: {}
          Latencies:
            Max_inMx: 0
            Mean_inMs: 23.08
            Min_inMs: 0
            P50_inMs: 0
            P90_inMs: 0
            P95_inMs: 0
            P99_inMs: 0
          RequestCounts: 100
          StartTime: "2023-08-01T08:36:56Z"
          StatusCodes:
            "200": 100
          SuccessCounts: 100
          TPS: 9.967817030760152
          TotalDataSize: 36968 byte
        Succeed: true
        SucceedRate: 1
        TargetMethod: GET
        TargetUrl: http://172.18.0.51:80
        MeanDelay: 23.08
      - TargetName: AgentNodePortV4IP_172.18.0.3_32713
        Metrics:
          ...
        Succeed: true
        SucceedRate: 1
        TargetMethod: GET
        TargetUrl: http://172.18.0.3:32713
        MeanDelay: 68.42
      - TargetName: AgentPodV4IP_kdoctor-agent-ntp9l_172.40.0.6
        Metrics:
          ...
        Succeed: true
        SucceedRate: 1
        TargetMethod: GET
        TargetUrl: http://172.40.0.6:80
        MeanDelay: 44.049503
      - TargetName: AgentClusterV4IP_172.41.249.6:80
        Metrics:
          ...
        Succeed: true
        SucceedRate: 1
        TargetMethod: GET
        TargetUrl: http://172.41.249.6:80
        MeanDelay: 26.307692
      - TargetName: AgentPodV4IP_kdoctor-agent-krrnp_172.40.1.5
        Metrics:
          ...
        Succeed: true
        SucceedRate: 1
        TargetMethod: GET
        TargetUrl: http://172.40.1.5:80
        MeanDelay: 61.564358
      - TargetName: AgentIngress_http://172.18.0.50/kdoctoragent
        Metrics:
          ...
        Succeed: true
        SucceedRate: 1
        TargetMethod: GET
        TargetUrl: http://172.18.0.50/kdoctoragent
        MeanDelay: 65.47059
      Succeed: true
      TargetNumber: 6
      TargetType: NetReach
    NetReachTaskSpec:
    ...
    PodName: kdoctor-agent-ntp9l
    ReportType: agent test report
    RoundDuration: 11.178657432s
    RoundNumber: 1
    RoundResult: succeed
    StartTimeStamp: "2023-08-01T08:36:56Z"
    EndTimeStamp: "2023-08-01T08:37:07Z"
    TaskName: netreach.task
    TaskType: NetReach
  - NodeName: kdoctor-worker
    NetReachTask:
      Detail:
      - TargetName: AgentClusterV4IP_172.41.249.6:80
        Metrics:
          ...
        Succeed: true
        SucceedRate: 1
        TargetMethod: GET
        TargetUrl: http://172.41.249.6:80
        MeanDelay: 47.25
      - TargetName: AgentNodePortV4IP_172.18.0.2_32713
        Metrics:
          ...
        Succeed: true
        SucceedRate: 1
        TargetMethod: GET
        TargetUrl: http://172.18.0.2:32713
        MeanDelay: 13.480392
      - TargetName: AgentPodV4IP_kdoctor-agent-krrnp_172.40.1.5
        Metrics:
          ...
        Succeed: true
        SucceedRate: 1
        TargetMethod: GET
        TargetUrl: http://172.40.1.5:80
        MeanDelay: 39.637257
      - TargetName: AgentPodV4IP_kdoctor-agent-ntp9l_172.40.0.6
        Metrics:
          ...
        Succeed: true
        SucceedRate: 1
        TargetMethod: GET
        TargetUrl: http://172.40.0.6:80
        MeanDelay: 51.38614
      - TargetName: AgentLoadbalancerV4IP_172.18.0.51:80
        Metrics:
          ...
        Succeed: true
        SucceedRate: 1
        TargetMethod: GET
        TargetUrl: http://172.18.0.51:80
        MeanDelay: 41.735847
      - TargetName: AgentIngress_http://172.18.0.50/kdoctoragent
        Metrics:
          ...
        Succeed: true
        SucceedRate: 1
        TargetMethod: GET
        TargetUrl: http://172.18.0.50/kdoctoragent
        MeanDelay: 60.463634
      Succeed: true
      TargetNumber: 6
      TargetType: NetReach
    NetReachTaskSpec:
    ...
    PodName: kdoctor-agent-krrnp
    ReportType: agent test report
    RoundDuration: 11.180813761s
    RoundNumber: 1
    RoundResult: succeed
    StartTimeStamp: "2023-08-01T08:36:56Z"
    EndTimeStamp: "2023-08-01T08:37:07Z"
    TaskName: netreach.task
    TaskType: NetReach
  ReportRoundNumber: 1
  RoundNumber: 1
  Status: Finished
  TaskName: task
  TaskType: NetReach
```

## 环境清理

```shell
kubectl delete netreach task
```