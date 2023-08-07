# netdns

[**English**](./netdns.md) | **简体中文**

## 基本描述

对于这种任务，每个 kdoctor 代理都会向指定的目标发送 Dns 请求，并获得成功率和平均延迟。 它可以指定成功条件来告知结果成功或失败。

## netdns 示例

### 集群 Dns Server 检查

对集群内的 dns server（coredns）发送对应请求，获取集群内 dns server 性能状态。

```yaml
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
    enableLatencyMetric: false
    targetDns:
      serviceName: kube-dns
      serviceNamespace: kube-system
      testIPv4: true
      testIPv6: true
status:
  doneRound: 1
  expectedRound: 1
  finish: true
  history:
    - deadLineTimeStamp: "2023-07-28T09:45:03Z"
      duration: 15.809063339s
      endTimeStamp: "2023-07-28T09:44:18Z"
      expectedActorNumber: 2
      failedAgentNodeList: []
      notReportAgentNodeList: []
      roundNumber: 1
      startTimeStamp: "2023-07-28T09:44:03Z"
      status: succeed
      succeedAgentNodeList:
        - kdoctor-control-plane
        - kdoctor-worker
  lastRoundStatus: succeed
```

### 指定 dns server 检查

对集群外部的 dns server 发送对应请求，获取集群外部 dns server 性能状态。

```Yaml
apiVersion: kdoctor.io/v1beta1
kind: Netdns
metadata:
  name: netdns- user
spec:
  expect:
    meanAccessDelayInMs: 1500
    successRate: 1
  request:
    domain: www.baidu.com
    durationInSecond: 10
    perRequestTimeoutInMS: 1000
    protocol: udp
    qps: 10
  schedule:
    roundNumber: 1
    roundTimeoutMinute: 1
    schedule: 0 1
  target:
    enableLatencyMetric: false
    targetUser:
      port: 53
      server: 172.41.54.83
```

##  netdns 定义

### Metadata

| 字段  | 描述            | 结构     | 验证  |
|-----|---------------|--------|-----|
| name | netdns 资源的名称 | string | 必填  |

### Spec

| 字段       | 描述               | 结构                                       | 验证      | 取值    | 默认值  |
|-----------|------------------|------------------------------------------|---------|-------|------|
| schedule  | 调度任务执行           | [schedule](./netdns-zh_CN.md#Schedule) | 可选      |       |      |
| request   | 对目标地址请求配置        | [request](./netdns-zh_CN.md#Request)   | 可选      |       |      |
| target    | 请求目标设置           | [target](./netdns-zh_CN.md#Target)     | 可选      |       |      |
| expect    | 任务成功条件判断         | [expect](./netdns-zh_CN.md#Expect)     | 可选      |       |      |


#### Schedule

| 字段                 | 描述                                    | 结构     | 验证  | 取值                                                                                                                                                                                                          | 默认值  |
|--------------------|---------------------------------------|--------|-----|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|------|
| roundNumber        | 任务执行轮数                                | int    | 可选  | 大于等于-1，为 -1 时表示永久执行,大于 0 表示将要执行的轮数                                                                                                                                                             | 1    |
| schedule           | 任务执行时间, 执行时间应小于roundTimeoutMinute     | string | 可选  | 支持 linux crontab 与间隔法两种写法<br/>[linux crontab](https://linuxhandbook.com/crontab/) ： */1 * * * * 表示每分钟执行一次 <br/>间隔法：书写格式为 “M N” ，M取值为一个数字，表示多少分钟之后开启任务，N取值为一个数字，表示每一轮任务的间隔多少分钟执行，例如 “0 1” 表示立即开始任务，每轮任务间隔1min | "0 60" |
| roundTimeoutMinute | 任务超时时间，需要大于 durationInSecond 和 任务执行时间 | int    | 可选  | 大于等于 1                                                                                                                                                                                                      | 60   |

#### Request

| 字段                     | 描述                                    | 结构     | 验证  | 取值            | 默认值           |
|------------------------|---------------------------------------|--------|-----|---------------|---------------|
| durationInSecond       | 每轮任务的请求发压的持续时间，小于roundTimeoutMinute   | int    | 可选  | 大于等于 1        | 2             |
| perRequestTimeoutInMS  | 每个请求的超时时间，不可大于 durationInSecond       | int    | 可选  | 大于等于 1        | 500           |
| qps                    | 每一个 agent 每秒请求数量                      | int    | 可选  | 大于等于 1        | 5             |
| protocol               | 请求协议                                  | string | 可选  | udp、tcp、tcp-tls | udp           |
| domain                 | dns 请求解析的域名                           | string | 可选  |               | kubernetes.default.svc.cluster.local |

> 注意：使用 agent 请求时，所有的 agent 都会向目标地址进行请求，因此实际 server 接收的 qps 等于 agent 数量 * 设置的qps。

#### Target

| 字段                 | 描述                      | 结构                                           | 验证  | 取值         | 默认值   |
|--------------------|-------------------------|----------------------------------------------|-----|------------|-------|
| targetUser        | 对用户自定义的 dns server 进行 dns 请求| [targetUser](./netdns-zh_CN.md#TargetUser) | 可选  |            | true  |
| targetDns           | 对集群的 dns server（coredns）进行 dns 请求 | [targetDns](./netdns-zh_CN.md#TargetDns)   | 可选  |            | true  |
| enableLatencyMetric | 统计演示分布,开启后会增加内存使用量      | bool                                         | 可选  | true,false | false |

#### Expect

任务成功条件，若任务结果没有达到期望条件，任务失败

| 字段                 | 描述                              | 结构    | 验证  | 取值     | 默认值  |
|--------------------|---------------------------------|-------|-----|--------|------|
| meanAccessDelayInMs        | 平均延时,如果最终的结果 超过本值，任务会判定为失败      | int   | 可选  | 大于等于 1 | 5000 |
| successRate           | http请求成功率,如果最终的结果 小于本值，任务会判定为失败 | float | 可选  | 0-1    | 1    |

#### TargetUser

测试用户自定义 dns server

| 字段          | 描述            | 结构     | 验证  | 取值      | 默认值 |
|-------------|---------------|--------|-----|---------|---|
| server  | dns server 地址 | string | 必填  |         |  |
| port    | dns server 端口 | int    | 必填  | 1-65535 |   |


#### TargetDns

测试集群内 dns server

| 字段       | 描述                       | 结构     | 验证  | 取值        | 默认值   |
|----------|--------------------------|--------|-----|-----------|-------|
| testIPv4 | 测试 ipv4 地址 请求 A 记录       | bool   | 可选  | true,false | true  |
| testIPv6 | 测试 ipv6 地址 请求 AAAA 记录    | bool   | 可选  | true,false    | false |
| serviceName     | 集群 dns server service 地址 | string | 可选  |    |       |
| serviceNamespace  | 集群 dns server service 端口 | string | 可选  |    |       |

### status

| 字段                 | 描述       | 结构                                       | 取值              |
|--------------------|----------|------------------------------------------|-----------------|
| doneRound        | 完成的任务轮数  | int                                      |                 |
| expectedRound           | 期望执行的轮数  | int                                      |                 |
| finish | 任务是否完成   | bool      | true、false      |
| lastRoundStatus | 最后一轮任务状态 | string    | notstarted、ongoing、succeed、fail |
| history | 任务历史     | 元素为[history](./apphttphealthy-zh_CN.md#History)的数组 |                 |

#### History

| 字段                               | 描述              | 结构           | 取值                             |
|----------------------------------|-----------------|--------------|--------------------------------|
| roundNumber                      | 任务轮数            | int          |                                |
| status                           | 任务状态            | string       | notstarted、ongoing、succeed、fail |
| startTimeStamp                   | 本轮任务开始时间        | string       |                                |
| endTimeStamp                     | 本轮任务结束时间        | string       |                                |
| duration                         | 本轮任执行时间         | string       |                                |
| deadLineTimeStamp                | 本轮任务 deadline   | string       |                                |
| failedAgentNodeList              | 任务失败的 agent     | 元素为string的数组 |                                |
| succeedAgentNodeList             | 任务成功的 agent     |   元素为string的数组            |                                |
| notReportAgentNodeList           | 没有上传任务报告的 agent |   元素为string的数组            |                                |
