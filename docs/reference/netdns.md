# netdns

[**简体中文**](./netdns-zh_CN.md) | **English**

## Basic description

For this kind of task, kdoctor-controller generates the corresponding [agent](../concepts/runtime.md) and other resources. Each agent pod sends a Dns request to the specified target and gets the success rate and average latency. It can specify a success condition to inform the result of success or failure.

## netdns Example

### Cluster Dns Server Check

Sends a corresponding request to a dns server (coredns) in a cluster to get the performance status of the dns server in the cluster.

```yaml
apiVersion: kdoctor.io/v1beta1
kind: Netdns
metadata:
  name: netdns-cluster
spec:
  agentSpec:
    hostNetwork: false
    kind: DaemonSet
    terminationGracePeriodMinutes: 60
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

### Specify dns server to check

Send a corresponding request to a dns server outside the cluster to get the performance status of the dns server outside the cluster.

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
## netdns definition

### Metadata

| fields | description | structure | validation |
|-----|---------------|--------|-----|
| name | Name of the netdns resource | string | required |

### Spec

| fields | description | structure | validation | take values | default |
|-----------|-------------|--------------------------------------------|---------|-------|------|
| agentSpec | task execution agent configuration | [agentSpec](./apphttphealthy.md#agentspec) | optional | | |
| schedule | Schedule Task Execution | [schedule](./apphttphealthy.md#schedule) | optional | | |
| request | Request configuration for destination address | [request](./netdns-zh_CN.md#Request) | Optional | | |
| target | Request target settings | [target](./apphttphealthy.md#target) | Optional | | |
| expect | Task success condition judgment | [expect](./apphttphealthy.md#expect) | Optional | | |

#### AgentSpec

| Fields | Description | Structure | Validation | Values | Default |
|-------------------------------|------------------------|----------------------------------------------------------------------------------------------------------------------------------|-----|----------------------|-------------------------------|
| annotation | annotation of agent workload | map[string]string | optional | | | |
| kind | type of agent workload | string | optional | Deployment, DaemonSet | DaemonSet |
| deploymentReplicas | The expected number of replicas when the agent workload type is deployment | int | optional | greater than or equal to 0 | 0 |
| affinity | agent workload affinity | labelSelector | optional | | |
| env | agent workload environment variable | env | optional | | | | hostNetwork | agent
| hostNetwork | agent Whether or not the workload uses the host network | bool | optional | true, false | false |
| resources | agent workload resource usage configuration | resources | optional | | limit cpu:1000m,memory:1024Mi |
| terminationGracePeriodMinutes | agent How many minutes after a workload completes a task before it terminates | int | optional | greater than or equal to 0 | 60 |

#### Schedule

| Fields | Description | Structure | Validation | Values | Defaults |
|--------------------|---------------------------------------|--------|-----|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|------|
| roundNumber        | task execution rounds | int | optional | greater than or equal to -1, -1 means permanent, greater than 0 means rounds will be executed | 1 | schedule | task execution time, execution time should be less than roundTimeoutMinute | string | optional | Support linux crontab and interval method.<br/>[linux crontab](https://linuxhandbook.com/crontab/) : */1 * * * * * means execute every minute <br/>Interval method: written in "M N" format, M takes the value of a number, indicating how many minutes after opening the task, N takes the value of a number, indicating how many minutes between each round of task execution, for example, "0 1" means start the task immediately, each round of task interval 1min | "0 60" | roundTimeoutMininute
| roundTimeoutMinute | Task timeout, needs to be greater than durationInSecond and task execution time | int | optional | greater than or equal to 1 | 60 |

#### Request

| Fields | Description | Structure | Validation | Values | Defaults |
|------------------------|---------------------------------------|--------|-----|---------------|---------------|
| durationInSecond | Duration of request send pressure for each round of tasks less than roundTimeoutMinute | int | optional | greater than or equal to 1 | 2 |
| perRequestTimeoutInMS | timeout per request, not greater than durationInSecond | int | optional | greater than or equal to 1 | 500 |
| qps | requests per second per agent | int | optional | greater than or equal to 1 | 5 | protocol | request protocol
| protocol | request protocol | string | optional | udp, tcp, tcp-tls | udp | domain | dns Requests for
| domain | dns The domain name for which the request is to be resolved | string | optional | | kubernetes.default.svc.cluster.local | | kubernetes.default.svc.

> Note: When using agent requests, all agents make requests to the destination address, so the actual qps received by the server is equal to the number of agents * the set qps.

#### Target

| Fields | Description | Structures | Validation | Values | Defaults |
|--------------------|-------------------------|----------------------------------------------|-----|------------|-------|
| targetUser | dns request to user-defined dns server | [targetUser](./netdns.md#targetuser) | Optional | | true |
| targetDns | Make a dns request to the cluster's dns server (coredns) | [targetDns](./netdns.md#targetuser) | Optional | | true |
| enableLatencyMetric | statisticsDemoDistribution,enable to increase memory usage | bool | optional | true,false | false |

#### Expect

Task success condition, if the task result does not meet the expected condition, the task will fail.

| fields | description | structure | validation | values | default |
| --------------------| ---------------------------------| -------| -----| --------| ------|
| meanAccessDelayInMs | meanAccessDelayInMs | The average delay, if the final result exceeds this value, the task will be judged as failed | int | optional | greater than or equal to 1 | 5000 |
| successRate | The success rate of the http request, if the final result is less than this value, the task will fail | float | optional | 0-1 | 1 |

#### TargetUser

Test user customized dns server

| Fields | Description | Structure | Validation | Values | Defaults |
|-------------|---------------|--------|-----|---------|---|
| server | dns server address | string | mandatory | | |
| port | dns server port | int | required | 1-65535 | |

#### TargetDns

Test dns server in cluster

| fields | description | structure | validation | values | defaults |
|----------|--------------------------|--------|-----|-----------|-------|
| testIPv4 | testIPv4 | Test ipv4 address request A record | bool | optional | true,false | true |
| testIPv6 | Test ipv6 address request AAAA record | bool | optional | true,false | false |
| serviceName | cluster dns server service address | string | optional | | |
| serviceNamespace | cluster dns server service port | string | optional | | |

### status

| fields | description | structure | values |
|--------------------|----------|------------------------------------------|-----------------|
| doneRound | number of completed task rounds | int | |
| expectedRound | number of rounds expected to be performed | int | | |
| finish | Whether the task is complete or not | bool | true, false |
| lastRoundStatus | lastRoundStatus | string | notstarted, on-going, succeed, fail |
| history | Task history | Element is [history](./apphttphealthy.md#history) array | |

#### History

| Fields | Description | Structure | Values |
| ----------------------------------|-----------------|--------------|--------------------------------|
| roundNumber | Task round number | int | |
| status | Task Status | string | notstarted, on-going, succeed, fail |
| startTimeStamp | startTimeStamp | startTimeStamp | string | |
| startTimeStamp | startTimeStamp | startTimeStamp | endTimeStamp | endTimeStamp | endTimeStamp | string | |
| deadLineTimeStamp | deadLineTimeStamp | deadLineTimeStamp | deadLineTimeStamp | deadline | string | |
| failedAgentNodeList | failedAgentNodeList | Array of failed agents | string | |
| notReportAgentNodeList | agent who did not upload a task report | array of elements as string | | notReportAgentNodeList | agent who failed to upload a task report | array of elements as string｜
