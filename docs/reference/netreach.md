# NetReach

[**简体中文**](./netreach-zh_CN.md) | **English**

## Basic description 

For this kind of task, kdoctor-controller will generate corresponding [agent](../concepts/runtime.md) and other resources. Each agent Pod sends http requests to each other with the request address of each agent's Pod IP, service IP, ingress IP and so on, and obtains the success rate and average latency. It can specify the success condition to determine whether the result is successful or not. Detailed reports can be obtained through the aggregation API.

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

| Fields | Description | Structure | Validation |
|-----|---------------|--------|-----|
| Name | Name of the NetReach Resource | String | Required |

### Spec

| Fields | Description | Structure | Validation |  Values | Default |
|-----------|-------------|--------------------------------------------|---------|-------|------|
|  agentSpec | Task Execution Agent Configuration | [agentSpec](./apphttphealthy.md#agentspec) | Optional |       |      |
| Schedule  |Schedule Task Execution | [schedule](./apphttphealthy.md#schedule) | Optional |       |      |
|Request   |Request Configuration for Destination Address | [request](./netdns.md#request) | Optional |       |      |
|Target    | Request Target Settings | [target](./apphttphealthy.md#target) | Optional |       |      |
|Expect    |Task Success Condition Judgment | [expect](./apphttphealthy.md#expect) | Optional |       |      |


#### AgentSpec

| Fields | Description | Structure | Validation | Values | Default |
|-------------------------------|------------------------|----------------------------------------------------------------------------------------------------------------------------------|-----|----------------------|-------------------------------|
|Annotation | Annotation of Agent Workload |Map[string]String | Optional | | |
| kind |Type of Agent Workload | String | Optional | Deployment, DaemonSet | DaemonSet |
| deploymentReplicas | The expected number of replicas when the agent workload type is deployment | int | Optional | Greater than or equal to 0 | 0 |
|Affinity | Agent Workload Affinity | labelSelector | Optional | | |
| env | Agent Workload Environment Variable | env |Optional | | |
| hostNetwork | Whether or not the agent workload uses the host network | Bool | Optional | True, false | False |
| Resources | Agent Workload Resource Usage Configuration | Resources | Optional | | Limit cpu: 1000m,Memory:1024Mi |
| terminationGracePeriodMinutes | the minutes after a agent workload completes a task before it terminates | int | Optional |Greater than or equal to 0 | 60 |


#### Schedule

| Fields | Description | Structure | Validation | Values | Defaults |
|--------------------|---------------------------------------|--------|-----|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|------|
| roundNumber        |Task Execution Rounds | int | Optional | A value greater than or equal to -1 indicates indefinite execution, with -1 representing permanent execution. A value greater than 0 represents the number of rounds to be executed | 1 | Schedule | Task execution time which should be less than roundTimeoutMinute | String | Optional | Support linux crontab and interval method<br/>[linux crontab](https://linuxhandbook.com/crontab/) : */1 * * * * * means execute every minute <br/>Interval method: writing format "M N". M is a number that indicates how many minutes after the task is started; N is a number that indicates how many minutes between each round of tasks. For example, "0 1" means start the task immediately, 1min between each round of tasks. | "0 60" |
| roundTimeoutMinute | Task timeout which needs to be greater than durationInSecond and task execution time | int | Optional | Greater than or equal to 1 | 60 |

#### Request

| Fields | Description | Structure | Validation | Values | Defaults |
|------------------------|---------------------------------------|--------|-----|---------------|---------------|
| durationInSecond | Duration of request send pressure for each round of tasks which is less than roundTimeoutMinute | int |Optional | Greater than or equal to 1 | 2 |
| perRequestTimeoutInMS | Timeout per request, not greater than durationInSecond | int |Optional | Greater than or equal to 1 | 500 |
| QPS | Requests per second per agent | int | Optional | Greater than or equal to 1 | 5 |

> When using agent requests, all agents will make requests to the destination address, so the actual QPS received by the server is equal to the number of agents multiplied by the set QPS.

#### Target

| Fields | Descriptions | Structures | Validations | Values | Defaults |
|--------------------|-------------------------|--------|-----|------------|-------|
| ClusterIP        | Test cluster service's cluster IP | Bool   | Optional  | True,false |True  |
| Endpoint           | Test cluster Pod endpoint       | Bool | Optional   | True,false   | True  |
| multusInterface | Test cluster Pod Multus multi-NIC IP  | Bool | Optional   | True,false  | False |
| IPv4 | Test IPv4                 | Bool | Optional   | True,false  | True  |
| IPv6 | Test IPv6                 | Bool | Optional   | True,false  |False |
|Ingress | Test Ingress Address           | Bool | Optional   | True,false  | False |
| nodePort | Test Service Node Port    | Bool | Optional  | True,false  | True  |
| enableLatencyMetric | Statistics demo distribution, which increases memory usage when turned on      | Bool | Optional   | True,false  |False |

#### Expect

Task success condition. If the task result does not meet the expected condition, the task will fail.

| Fields | Description | Structures | Validation | Values | Default |
| --------------------| ---------------------------------| -------| -----| --------| ------|
|meanAccessDelayInMs | The average delay. If the final result exceeds this value, the task will be judged as failed | int | Optional | Greater than or equal to 1 | 5000 |
| successRate | Success rate of the HTTP request. If the final result is less than this value, the task will fail | Float | Optional | 0-1 | 1 |

### status

| Fields | Description | Structures | Values |
|--------------------|----------|------------------------------------------|-----------------|
| doneRound | Number of completed task rounds | int | |
| expectedRound | Number of rounds expected to be performed | int | |
| Finish | Whether the task is complete or not |Bool |True, false |
| lastRoundStatus | lastRoundStatus | String | Notstarted, on-going, succeed, fail |
| History | Task History | Element is [history](./apphttphealthy.md#history) array | |

#### History

| Fields | Description | Structure | Values |
| ----------------------------------|-----------------|--------------|--------------------------------|
| roundNumber | Task Round Number | int | |
| Status | Task Status | String | Notstarted, on-going, succeed, fail |
| startTimeStamp | Start of the current round of tasks | String | |
| endTimeStamp | End of the current round of tasks | string |  |
| duration |Execution time of the current round of tasks |string | |
| deadLineTimeStamp | Deadline of the current round of tasks | string | |
| failedAgentNodeList | Agent whose tasks failed |Array of elements as string | |
| succeedAgentNodeList |Agent whose task succeeded | Array of elements as string | |
| notReportAgentNodeList |Agent who did not upload a task report | Array of elements as string | |
