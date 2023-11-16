# AppHttpHealthy

[**简体中文**](./apphttphealthy-zh_CN.md) | **English**

## Basic description 

For this kind of task, kdoctor-controller generates the corresponding [agent](../concepts/runtime.md) and other resources. Each agent Pod sends an HTTP request to the specified target and gets the success rate and average latency. It can specify the success condition to determine whether the result is successful or not. And, detailed reports can be obtained through the aggregation API.

## AppHttpHealthy Example

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

## AppHttpHealthy Definitions

### Metadata

| Fields | Descriptions | Structures | Validations |
|-----|---------------|--------|-----|
| Name | Name of the AppHttpHealthy Resource | String | Required |

### Spec

| Fields | Description | Structure | Validation | Values | Default |
|-----------|-------------|--------------------------------------------|---------|-------|------|
| agentSpec | Task Execution Agent Configuration | [agentSpec](./apphttphealthy.md#agentspec) | Optional | | |
| Schedule | Scheduling Task Execution | [schedule](./apphttphealthy.md#schedule) | Optional | | |
| Request | Request Configuration for a Destination Address | [request](./apphttphealthy.md#request) | Optional | | | |
| Target | Request Target Settings | [target](./apphttphealthy.md#target) | Optional | | | |
| Expect | Task Success Condition Judgment | [expect](./apphttphealthy.md#expect) | Optional | | |

#### AgentSpec

| Fields | Description | Structure | Validation | Values | Default |
|-------------------------------|------------------------|----------------------------------------------------------------------------------------------------------------------------------|-----|----------------------|-------------------------------|
|Annotation |Annotation of Agent Workload | Map[string]String | Optional | | | |
| Kind | Type of Agent Workload | String | Optional | Deployment, DaemonSet | DaemonSet |
| deploymentReplicas | The expected number of replicas when the agent workload type is deployment | int | Optional |Greater than or Equal to 0 | 0 |
|Affinity | Agent Workload Affinity | labelSelector | Optional | | |
| env | Agent Workload Environment Variable | env | Optional | | | | hostNetwork | Agent
| hostNetwork | Whether or not the agent workload uses the host network | Bool | Optional |True, false |False |
| Resources | Agent Workload Resource Usage Configuration | Resources | Optional | | Limit CPU:1000m,Memory:1024Mi |
| terminationGracePeriodMinutes | the minutes after a agent workload completes a task before it terminates | int | Optional | Greater than or equal to 0 | 60 |

#### Schedule

| Fields | Description | Structure | Validation | Values | Defaults |
|--------------------|---------------------------------------|--------|-----|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-------|
| roundNumber       | Task Execution Rounds | int |Optional |A value greater than or equal to -1 indicates indefinite execution, with -1 representing permanent execution. A value greater than 0 represents the number of rounds to be executed | 1 | Schedule | Task execution time which should be less than roundTimeoutMinute | String | Optional | Support linux crontab and interval method<br/>[linux crontab](https://linuxhandbook.com/crontab/) : */1 * * * * * means execute every minute <br/>Interval method: writing format "M N". M is a number that indicates how many minutes after the task is started; N is a number that indicates how many minutes between each round of tasks. For example, "0 1" means start the task immediately, 1min between each round of tasks.| "0 1" |
| roundTimeoutMinute | Task timeout which needs to be greater than durationInSecond and task execution time | int | optional | greater than or equal to 1 | 60 |

#### Request

| Fields | Description | Structure | Validation | Values | Defaults |
|------------------------|---------------------------------------|--------|-----|---------------|---------------|
| durationInSecond | Duration of request send pressure for each round of tasks which is less than roundTimeoutMinute | int |Optional |Greater than or equal to 1 | 2 |
| perRequestTimeoutInMS |Timeout per request, not greater than durationInSecond | int |Optional | Greater than or equal to 1 | 500 |
| QPS | Number of requests per second per agent | int | Optional | Greater than or equal to 1 | 5 |

> When using agent requests, all agents will make requests to the destination address, so the actual QPS received by the server is equal to the number of agents multiplied by the set QPS.

#### Target

| Fields | Description | Structures | Validation | Values | Defaults |
|------------------------|---------------------------------------------------------------------------------------------------------------------|--------|---------|---------------------------|-------|
| Host | HTTP Request Address |String | | Required | | | | |
|Method | HTTP Request Method | String | Required | GET, POST, PUT, DELETE, CONNECT, OPTIONS, PATCH, HEAD | |
| bodyConfigmapName | The name of the configmap stored in the body of the HTTP request. Refer to [configmap](./apphttphealthy.md#body). If you don't need a body request, ignore this field.| String| Optional|         |       |
| bodyConfigmapNamespace | HTTP request body's configmap namespace. If bodyConfigmapName is not empty, you need to set this field | string | optional | | |
| tlsSecretName | The name of the secret where the HTTPs request certificate is stored, with a secret of type kubernetes.io/tls. Refer to [secret](./apphttphealthy.md#tls). Ignore this field if using a protocol other than HTTPs |String | optional | | |
| tlsSecretNamespace | The secret namespace where the HTTPs request certificate is stored. If the tlsSecretName field is not null, you need to set this field |String | Optional | | |
| header | HTTP request header, in the form of an array, for example "Content-Type: application/json" | Elements are an array of strings |Optional | | |
| HTTP2 | Use the request HTTP2 protocol switch | bool | Optional | True,false | False |
| enableLatencyMetric | Statistics demo distribution, which increases memory usage when turned on | Bool | Optional | True,false | False |

#### Expect

Task success condition. If the task result does not meet the expected condition, the task fails.

| Fields | Description | Structures | Validation | Values | Defaults |
|--------------------|---------------------------------|-------|-----|--------|------|
| meanAccessDelayInMs | The average delay. If the final result exceeds this value, the task will be judged as failed | int | Optional | Greater than or equal to 1 | 5000 |
| successRate | Success rate of the HTTP request. If the final result is less than this value, the task will fail | Float | Optional | 0-1 | 1 |
| statusCode | The expected HTTP return status code. If the final result is not equal to this value, the task will be determined to have failed | int | Optional | 0-600 | 200 |

#### Body

Carry a body request. Example of how to write a body

```yaml
apiVersion: v1
kind: ConfigMap
metadata: name: kdoctor-test-body-5656-89207258
  name: kdoctor-test-body-5656-89207258
  namespace: kdoctor
data: kdoctor-test-body-5656-89207258 namespace: kdoctor
  test1: test1
  test2: test2
```

#### TLS

If you don't fill in the CA certificate when using certificate for HTTPs request, the certificate security will not be verified, HTTPs request certificate example:

```yaml
apiVersion: v1
data: ca.crt: xxxxxxx
  ca.crt: xxxxxxxxxxxbase64
  tls.crt: xxxxxxxxxxxbase64
  tls.key: xxxxxxxxxxxbase64
kind: Secret
metadata.
  name: kdoctor-client-cert
  namespace: kdoctor
type: kubernetes.io/tls
```


### Status

| Fields | Description | Structures | Values |
|--------------------|----------|------------------------------------------------|--------------|
| doneRound | Number of completed task rounds | int | |
| expectedRound | Number of rounds expected to be performed | int | |
| Finish | Whether the task is complete or not |Bool | True, false |
| lastRoundStatus | lastRoundStatus | String |Notstarted, on-going, succeed, fail |
| History | Task History | Element is [history](./apphttphealthy.md#history) array | |

#### History

| Fields | Description | Structure | Values |
|----------------------------------|-----------------|--------------|---------------------------|
| roundNumber | Task Round Number | int | |
| Status | Task Status | String | Notstarted, ongoing, succeed, fail |
| startTimeStamp | Start of the current round of tasks | String | |
| endTimeStamp | End of the current round of tasks | string |  |
| duration |Execution time of the current round of tasks |string | |
| deadLineTimeStamp | Deadline of the current round of tasks | string | |
| failedAgentNodeList | Agent whose tasks failed |Array of elements as string | |
| succeedAgentNodeList |Agent whose task succeeded | Array of elements as string | |
| notReportAgentNodeList |Agent who did not upload a task report | Array of elements as string | |
