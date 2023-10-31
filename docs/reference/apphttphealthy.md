# AppHttpHealthy

[**简体中文**](./apphttphealthy-zh_CN.md) | **English**

## Basic description 

For this kind of task, kdoctor-controller generates the corresponding [agent](../concepts/runtime.md) and other resources. Each agent pod sends an http request to the specified target and gets the success rate and average latency. It can specify the success condition to determine whether the result is successful or not. And, detailed reports can be obtained through the aggregation API.

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
| name | Name of the AppHttpHealthy resource | string | mandatory |

### Spec

| fields | description | structure | validation | values | default |
|-----------|-------------|--------------------------------------------|---------|-------|------|
| agentSpec | task execution agent configuration | [agentSpec](./apphttphealthy.md#agentspec) | optional | | |
| schedule | Scheduling Task Execution | [schedule](./apphttphealthy.md#schedule) | optional | | |
| request | Request configuration for a destination address | [request](./apphttphealthy.md#request) | optional | | | |
| target | Request target settings | [target](./apphttphealthy.md#target) | optional | | | |
| expect | Task success condition judgment | [expect](./apphttphealthy.md#expect) | optional | | |

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
|--------------------|---------------------------------------|--------|-----|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-------|
| roundNumber       | task execution rounds | int | optional | greater than or equal to -1, -1 means permanent, greater than 0 means rounds will be executed | 1 | schedule | task execution time, execution time should be less than roundTimeoutMinute | string | optional | Support linux crontab and interval method.<br/>[linux crontab](https://linuxhandbook.com/crontab/) : */1 * * * * * means execute every minute <br/>Interval method: written in "M N" format, M takes the value of a number, which means how many minutes after opening the task, N takes the value of a number, indicating how many minutes between each round of task execution, for example, "0 1" means start the task immediately, each round of task interval 1min | "0 1" | roundTimeoutMininute
| roundTimeoutMinute | Task timeout, need to be greater than durationInSecond and task execution time | int | optional | greater than or equal to 1 | 60 |

#### Request

| Fields | Description | Structure | Validation | Values | Defaults |
|------------------------|---------------------------------------|--------|-----|---------------|---------------|
| durationInSecond | Duration of request send pressure for each round of tasks less than roundTimeoutMinute | int | optional | greater than or equal to 1 | 2 |
| perRequestTimeoutInMS | timeout per request, not greater than durationInSecond | int | optional | greater than or equal to 1 | 500 |
| qps | Number of requests per second per agent | int | optional | greater than or equal to 1 | 5 |

> Note: When using agent requests, all agents will make requests to the destination address, so the actual qps received by the server is equal to the number of agents * the set qps.
#### Target

| Fields | Description | Structures | Validation | Values | Defaults |
|------------------------|---------------------------------------------------------------------------------------------------------------------|--------|---------|---------------------------|-------|
| host | http request address | string | | Required | | | | |
| method | http request method | string | mandatory | GET, POST, PUT, DELETE, CONNECT, OPTIONS, PATCH, HEAD | | bodyConfigmap
| bodyConfigmapName | The name of the configmap stored in the body of the http request, [configmap content reference](./apphttphealthy.md#body)
| bodyConfigmapNamespace | http request body's configmap namespace, if bodyConfigmapName is not empty, need to set this field | string | optional | | | | | | | | | tlsSecretName
| tlsSecretName | The name of the secret where the https request certificate is stored, with a secret of type kubernetes.io/tls, [secret content reference](./apphttphealthy.md#tls), ignore this field if using a protocol other than https | string | | optional | | | | | tlsSecretName | tlsSecretName
| tlsSecretNamespace | The secret namespace where the https request certificate is stored, if the tlsSecretName field is not null, you need to set this field | string | optional | | | |
| header | http request header, as an array, for example "Content-Type: application/json" | elements are an array of strings | optional | | | |
| http2 | Use the request using http2 protocol switch | bool | optional | true,false | false |
| enableLatencyMetric | Statistical demo distribution, turn it on to increase memory usage | bool | optional | true,false | false |

#### Expect

Task success condition, if the task result does not meet the expected condition, the task fails.

| Fields | Description | Structures | Validation | Values | Defaults |
|--------------------|---------------------------------|-------|-----|--------|------|
| meanAccessDelayInMs | meanAccessDelayInMs | The average delay, if the final result exceeds this value, the task will be judged as failed | int | optional | greater than or equal to 1 | 5000 |
| successRate | The success rate of the http request, if the final result is less than this value, the task will fail | float | optional | 0-1 | 1 | statusCode | The status code to expect.
| statusCode | The status code to expect from the http request, if the final result is not equal to this value, the task will be judged as failed | int | optional | 0-600 | 200 |

#### Body

Carrying a body request, example of how to write a body

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

#### Tls

If you don't fill in the ca certificate when using certificate for https request, the certificate security will not be verified, https request certificate example.

``` yaml
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


### status

| fields | description | structures | values |
|--------------------|----------|------------------------------------------------|--------------|
| doneRound | number of completed task rounds | int | | |
| expectedRound | number of rounds expected to be performed | int | | |
| finish | Whether the task is complete or not | bool | true, false |
| lastRoundStatus | lastRoundStatus | string | notstarted, on-going, succeed, fail |
| history | Task history | Element is [history](./apphttphealthy.md#history) array | |

#### History

| Fields | Description | Structure | Values |
|----------------------------------|-----------------|--------------|---------------------------|
| roundNumber | Task round number | int | |
| status | Task status | string | notstarted, ongoing, succeed, fail |
| startTimeStamp | startTimeStamp | startTimeStamp | string | |
| startTimeStamp | startTimeStamp | startTimeStamp | endTimeStamp | endTimeStamp | endTimeStamp | string |
| deadLineTimeStamp | deadLineTimeStamp | deadLineTimeStamp | deadLineTimeStamp | deadline | string | |
| failedAgentNodeList | failedAgentNodeList | Array of failed agents | string |
| notReportAgentNodeList | agent who did not upload a task report | array of elements as string | | notReportAgentNodeList | agent who failed to upload a task report | array of elements as string|
