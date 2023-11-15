# AppHttpHealthy

[**English**](./apphttphealthy.md) | **简体中文**

## 基本描述

对于这种任务，kdoctor-controller 会根据 agentSpec 生成对应的 [agent](../concepts/runtime-zh_CN.md) 等资源，每一个 agent Pod 都会向指定的目标发送 HTTP请求，并获得成功率和平均延迟。它可以指定成功条件来判断结果是否成功。并且，可以通过聚合 API 获取详细的报告。

## AppHttpHealthy 示例

```yaml
apiVersion: kdoctor.io/v1beta1
kind: AppHttpHealthy
metadata:
  name: apphttphealth
spec:
  agentSpec:
    hostNetwork: false
    kind: DaemonSet
    terminationGracePeriodMinutes: 60
  expect:
    meanAccessDelayInMs: 1500
    successRate: 1
  request:
    durationInSecond: 10
    perRequestTimeoutInMS: 2000
    qps: 10
  schedule:
    roundNumber: 1
    roundTimeoutMinute: 1
    schedule: 0 1
  target:
    enableLatencyMetric: false
    host: http://www.baidu.com
    http2: false
    method: GET
status:
  doneRound: 1
  expectedRound: 1
  finish: true
  history:
    - deadLineTimeStamp: "2023-07-28T09:58:41Z"
      duration: 17.272005445s
      endTimeStamp: "2023-07-28T09:57:58Z"
      expectedActorNumber: 2
      failedAgentNodeList: []
      notReportAgentNodeList: []
      roundNumber: 1
      startTimeStamp: "2023-07-28T09:57:41Z"
      status: succeed
      succeedAgentNodeList:
        - kdoctor-worker
        - kdoctor-control-plane
  lastRoundStatus: succeed
```

## AppHttpHealthy 定义

### Metadata

| 字段  | 描述            | 结构     | 验证  |
|-----|---------------|--------|-----|
| name | AppHttpHealthy 资源的名称 | string | 必填  |

### Spec

| 字段       | 描述          | 结构                                         | 验证      | 取值    | 默认值  |
|-----------|-------------|--------------------------------------------|---------|-------|------|
| agentSpec  | 任务执行 agent 配置 | [agentSpec](./apphttphealthy-zh_CN.md#AgentSpec) | 可选      |       |      |
| schedule  | 调度任务执行      | [schedule](./apphttphealthy-zh_CN.md#Schedule)   | 可选      |       |      |
| request   | 对目标地址请求配置   | [request](./apphttphealthy-zh_CN.md#Request)     | 可选      |       |      |
| target    | 请求目标设置      | [target](./apphttphealthy-zh_CN.md#Target)       | 可选      |       |      |
| expect    | 任务成功条件判断    | [expect](./apphttphealthy-zh_CN.md#Expect)       | 可选      |       |      |


#### AgentSpec

| 字段                            | 描述                     | 结构                                                                                                                               | 验证  | 取值                   | 默认值                           |
|-------------------------------|------------------------|----------------------------------------------------------------------------------------------------------------------------------|-----|----------------------|-------------------------------|
| annotation                    | agent 工作负载的 annotation | map[string]string                                                                                                                | 可选  |                      |                               |
| kind                          | agent 工作负载的类型          | string                                                                                                                           | 可选  | Deployment、DaemonSet | DaemonSet                     |
| deploymentReplicas            | agent 工作负载类型为 deployment 时的期望副本数 | int                                                                                                                              | 可选  | 大于等于 0               | 0                             |
| affinity                      | agent 工作负载亲和性          | labelSelector | 可选  |                      |                               |
| env                           | agent 工作负载环境变量         | env                      | 可选  |                      |                               |
| hostNetwork                   | agent 工作负载是否使用宿主机网络    | bool                                                                                                                             | 可选  | true、false           | false                         |
| resources                     | agent 工作负载资源使用配置       | resources       | 可选  |                      | limit cpu:1000m,memory:1024Mi |
| terminationGracePeriodMinutes | agent 工作负载完成任务后多少分钟之后终止 | int                                                                                                                              | 可选  | 大于等于 0               | 60                            |

#### Schedule

| 字段                 | 描述                                    | 结构     | 验证  | 取值                                                                                                                                                                                                          | 默认值   |
|--------------------|---------------------------------------|--------|-----|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-------|
| roundNumber        | 任务执行轮数                                | int    | 可选  | 大于等于-1，为 -1 时表示永久执行,大于 0 表示将要执行的轮数                                                                                                                                                             | 1     |
| schedule           | 任务执行时间, 执行时间应小于roundTimeoutMinute     | string | 可选  | 支持 linux crontab 与间隔法两种写法<br/>[linux crontab](https://linuxhandbook.com/crontab/) ： */1 * * * * 表示每分钟执行一次 <br/>间隔法：书写格式为 “M N” ，M取值为一个数字，表示多少分钟之后开启任务，N取值为一个数字，表示每一轮任务的间隔多少分钟执行，例如 “0 1” 表示立即开始任务，每轮任务间隔1min | "0 1" |
| roundTimeoutMinute | 任务超时时间，需要大于 durationInSecond 和 任务执行时间 | int    | 可选  | 大于等于 1                                                                                                                                                                                                      | 60    |

#### Request

| 字段                     | 描述                                    | 结构     | 验证  | 取值            | 默认值           |
|------------------------|---------------------------------------|--------|-----|---------------|---------------|
| durationInSecond       | 每轮任务的请求发压的持续时间，小于roundTimeoutMinute   | int    | 可选  | 大于等于 1        | 2             |
| perRequestTimeoutInMS  | 每个请求的超时时间，不可大于 durationInSecond       | int    | 可选  | 大于等于 1        | 500           |
| qps                    | 每一个 agent 每秒请求数量                      | int    | 可选  | 大于等于 1        | 5             |

> 使用 agent 请求时，所有的 agent 都会向目标地址进行请求，因此实际 server 接收的 QPS 等于 agent 数量 * 设置的 QPS。

#### Target

| 字段                     | 描述                                                                                                                  | 结构     | 验证      | 取值                        | 默认值   |
|------------------------|---------------------------------------------------------------------------------------------------------------------|--------|---------|---------------------------|-------|
| host                   | HTTP 请求地址                                                                                                           | string | 必填      |                           |       |
| method                 | HTTP 请求方法                                                                                                           | string   | 必填      | GET、POST、PUT、DELETE、CONNECT、OPTIONS、PATCH、HEAD |       |
| bodyConfigmapName      | HTTP 请求 body 存放的 configmap 名称,[configmap 内容参考](./apphttphealthy-zh_CN.md#Body)，若不需要 body 请求，忽略此字段                   | string   | 可选      |          |       |
| bodyConfigmapNamespace | HTTP 请求 body 的 configmap 命名空间，如果 bodyConfigmapName 不为空，需要设置此字段                                                      | string   | 可选      |                           |       |
| tlsSecretName          | HTTPs 请求证书存放的 secret 名称，secret类型为 kubernetes.io/tls,[secret 内容参考](./apphttphealthy-zh_CN.md#Tls)，若使用协议非 https 忽略此字段 | string   | 可选      |                           |       |
| tlsSecretNamespace     | HTTPs 请求证书存放的 secret 命名空间，如果 tlsSecretName 字段不为空，需要设置此字段                                                            | string   | 可选      |                           |       |
| header                 | HTTP 请求头，数组形式,示例为 "Content-Type: application/json"                                                                  | 元素为字符串的数组  | 可选  |                           |       |
| HTTP2                  | 使用 HTTP2 协议进行请求开关                                                                                                 | bool   | 可选      | true,false                | false |
| enableLatencyMetric    | 统计演示分布,开启后会增加内存使用量                                                                                                  | bool   | 可选      | true,false                | false |

#### Expect

任务成功条件，若任务结果没有达到期望条件，任务失败

| 字段                 | 描述                              | 结构    | 验证  | 取值     | 默认值  |
|--------------------|---------------------------------|-------|-----|--------|------|
| meanAccessDelayInMs  | 平均延时,如果最终的结果 超过本值，任务会判定为失败      | int   | 可选  | 大于等于 1 | 5000 |
| successRate          | HTTP 请求成功率,如果最终的结果 小于本值，任务会判定为失败 | float | 可选  | 0-1    | 1    |
| statusCode           | 期待的 HTTP 返回状态码，如果最终的结果不等于本值，任务会判定为失败    | int   | 可选  | 0-600  | 200  |


#### Body

携带 body 请求，body 写法示例

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: kdoctor-test-body-5656-89207258
  namespace: kdoctor
data:
  test1: test1
  test2: test2
```


#### TLS

HTTPs 请求使用证书时，若不填写 CA 证书则不对证书安全进行校验，HTTPs 请求证书示例。

```yaml
apiVersion: v1
data:
  ca.crt:  xxxxxxxxxbase64
  tls.crt: xxxxxxxxxbase64
  tls.key: xxxxxxxxxbase64
kind: Secret
metadata:
  name: kdoctor-client-cert
  namespace: kdoctor
type: kubernetes.io/tls
```


### Status 

| 字段                 | 描述       | 结构                                             | 取值           |
|--------------------|----------|------------------------------------------------|--------------|
| doneRound        | 完成的任务轮数  | int                                            |              |
| expectedRound           | 期望执行的轮数  | int                                            |              |
| finish | 任务是否完成   | bool                                           | true、false   |
| lastRoundStatus | 最后一轮任务状态 | string        | notstarted、ongoing、succeed、fail |
| history | 任务历史     | 元素为[history](./apphttphealthy-zh_CN.md#History)的数组 |              |

#### History

| 字段                               | 描述              | 结构           | 取值                        |
|----------------------------------|-----------------|--------------|---------------------------|
| roundNumber                      | 任务轮数            | int          |                           |
| status                           | 任务状态            | string       | notstarted、ongoing、succeed、fail |
| startTimeStamp                   | 本轮任务开始时间        | string       |                           |
| endTimeStamp                     | 本轮任务结束时间        | string       |                           |
| duration                         | 本轮任务执行时间         | string       |                           |
| deadLineTimeStamp                | 本轮任务 deadline   | string       |                           |
| failedAgentNodeList              | 任务失败的 agent     | 元素为 string 的数组 |                           |
| succeedAgentNodeList             | 任务成功的 agent     | 元素为 string 的数组  |                           |
| notReportAgentNodeList           | 没有上传任务报告的 agent | 元素为 string 的数组  |                           |
