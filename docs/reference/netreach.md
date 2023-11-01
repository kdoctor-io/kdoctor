# NetReach

[**简体中文**](./netreach-zh_CN.md) | **English**

## Basic description 

For this kind of task, kdoctor-controller will generate corresponding [agent](../concepts/runtime.md) and other resources, and each agent pod sends http requests to each other with the request address of each agent's pod ip, service ip, ingress ip and so on, and obtains the success rate and average latency. It can specify the success condition to determine whether the result is successful or not. And, detailed reports can be obtained through the aggregation API.

## NetReach example

```yaml
apiVersion: kdoctor.io/v1beta1
kind: NetReach
metadata:
  name: netreach
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
    perRequestTimeoutInMS: 1000
    qps: 10
  schedule:
    roundNumber: 1
    roundTimeoutMinute: 1
    schedule: 0 1
  target:
    clusterIP: true
    enableLatencyMetric: false
    endpoint: true
    ingress: true
    ipv4: true
    ipv6: true
    loadBalancer: true
    multusInterface: false
    nodePort: true
status:
  doneRound: 1
  expectedRound: 1
  finish: true
  history:
    - deadLineTimeStamp: "2023-07-28T09:59:12Z"
      duration: 15.462579243s
      endTimeStamp: "2023-07-28T09:58:27Z"
      expectedActorNumber: 2
      failedAgentNodeList:
        - kdoctor-worker
        - kdoctor-control-plane
      failureReason: some agents failed
      notReportAgentNodeList: []
      roundNumber: 1
      startTimeStamp: "2023-07-28T09:58:12Z"
      status: fail
      succeedAgentNodeList: []
  lastRoundStatus: fail
```

##  NetReach Definition

### Metadata

| fields | description | structure | validation |
|-----|---------------|--------|-----|
| name | Name of the NetReach resource | string | required |

### Spec

| fields | description | structure | validation | take values | default |
|-----------|-------------|--------------------------------------------|---------|-------|------|
|  agentSpec | task execution agent configuration | [agentSpec](./apphttphealthy.md#agentspec) | optional |       |      |
| schedule  |Schedule Task Execution | [schedule](./apphttphealthy.md#schedule) | optional |       |      |
| request   |Request configuration for destination address | [request](./netdns.md#request) | Optional |       |      |
| target    | Request target settings | [target](./apphttphealthy.md#target) | Optional |       |      |
| expect    |Task success condition judgment | [expect](./apphttphealthy.md#expect) | Optional |       |      |


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
| roundNumber        | task execution rounds | int | optional | greater than or equal to -1, -1 means permanent, greater than 0 means rounds will be executed | 1 | schedule | task execution time, execution time should be less than roundTimeoutMinute | string | optional | Support linux crontab and interval method.<br/>[linux crontab](https://linuxhandbook.com/crontab/) : */1 * * * * * means execute every minute <br/>Interval method: written in "M N" format, M takes the value of a number, indicating how many minutes after opening the task, N takes the value of a number, indicating how many minutes between each round of task execution, for example, "0 1" means start the task immediately, each round of task interval 1min | "0 1" | roundTimeoutMininute
| roundTimeoutMinute | Task timeout, needs to be greater than durationInSecond and task execution time | int | optional | greater than or equal to 1 | 60 |

#### Request

| Fields | Description | Structure | Validation | Values | Defaults |
|------------------------|---------------------------------------|--------|-----|---------------|---------------|
| durationInSecond | Duration of request send pressure for each round of tasks less than roundTimeoutMinute | int | optional | greater than or equal to 1 | 2 |
| perRequestTimeoutInMS | timeout per request, not greater than durationInSecond | int | optional | greater than or equal to 1 | 500 |
| qps | requests per second per agent | int | optional | greater than or equal to 1 | 5 |

> Note: When using agent requests, all agents make requests to the destination address, so the actual qps received by the server is equal to the number of agents * the set qps.

#### Target

| Fields | Descriptions | Structures | Validations | Values | Defaults |
|--------------------|-------------------------|--------|-----|------------|-------|
| clusterIP        | Test cluster service's cluster ip | bool   | Optional  | true,false | true  |
| endpoint           | Test cluster pod endpoint       | bool | Optional   | true,false   | true  |
| multusInterface | Test cluster pod multus multi-NIC ip  | bool | Optional   | true,false  | false |
| ipv4 | Test ipv4                 | bool | Optional   | true,false  | true  |
| ipv6 | Test ipv6                 | bool | Optional   | true,false  | false |
| ingress | Test ingress address           | bool | Optional   | true,false  | false |
| nodePort | test service node port    | bool | Optional  | true,false  | true  |
| enableLatencyMetric | Count demo distribution, when turned on it will increase memory usage      | bool | Optional   | true,false  | false |

#### Expect

Task success condition, if the task result does not meet the expected condition, the task will fail.

| fields | description | structure | validation | values | default |
| --------------------| ---------------------------------| -------| -----| --------| ------|
| meanAccessDelayInMs | meanAccessDelayInMs | The average delay, if the final result exceeds this value, the task will be judged as failed | int | optional | greater than or equal to 1 | 5000 |
| successRate | The success rate of the http request, if the final result is less than this value, the task will fail | float | optional | 0-1 | 1 |

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
