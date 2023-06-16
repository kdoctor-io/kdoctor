# Nethttp

## concept 

Fo this kind task, each kdoctor agent will send http request to specified target, and get success rate and mean delay. 
It could specify success condition to tell the result succeed or fail. 
And, more detailed report will print to kdoctor agent stdout, or save to disc by kdoctor controller.

the following is the spec of nethttp
```shell
apiVersion: kdoctor.io/v1beta1
kind: NetReach
metadata:
  creationTimestamp: "2023-05-24T08:11:13Z"
  generation: 1
  name: netreach
  resourceVersion: "1427617"
  uid: 2f5da0d6-0252-4229-adca-a43a5d2ac4ff
spec:
  request:
    durationInSecond: 10
    perRequestTimeoutInMS: 1000
    qps: 10
  schedule:
    roundNumber: 2
    roundTimeoutMinute: 1
    schedule: 1 1
  expect:
    meanAccessDelayInMs: 10000
    successRate: 1
  target:
    clusterIP: true
    endpoint: true
    ingress: false
    ipv4: true
    ipv6: true
    loadBalancer: false
    multusInterface: false
    nodePort: true
status:
  doneRound: 2
  expectedRound: 2
  finish: true
  history:
  - deadLineTimeStamp: "2023-05-24T08:14:13Z"
    duration: 20.089468009s
    endTimeStamp: "2023-05-24T08:13:33Z"
    expectedActorNumber: 2
    failedAgentNodeList: []
    notReportAgentNodeList: []
    roundNumber: 2
    startTimeStamp: "2023-05-24T08:13:13Z"
    status: succeed
    succeedAgentNodeList:
    - kdoctor-worker
    - kdoctor-control-plane
```

* spec.schedule: set how to schedule the task.

      roundNumber: how many rounds it should be to run this task

      schedule: Support Linux crontab syntax for scheduling tasks, while also supporting simple writing. 
                The first digit represents how long the task will start, and the second digit represents the interval time between each round of tasks,
                separated by spaces. Example: "1 2" indicates that the task will start in 1 minute, and the interval time between each round of tasks.

      roundTimeoutMinute: the timeout in minute for each round, when the rask does not finish in time, it results to be failuire

      sourceAgentNodeSelector [optional]: set the node label selector, then, the kdoctor agent who locates on these nodes will implement the task. If not set this field, all kdoctor agent will execute the task

* spec.request: how each kdoctor agent should send the http request

    durationInSecond: for each round, the duration in second how long the http request lasts

    perRequestTimeoutInMS: timeout in ms for each http request 

    qps: qps

* spec.target: set the target of http request. it could not set targetUser and targetAgent at the same time

      target: [optional]: set the http tareget to kdoctor agents

        clusterIP: send http request to the cluster ipv4 or ipv6 address of kdoctor agnent, according to ipv4 and ipv6.

        endpoint: send http request to other kdoctor agnent ipv4 or ipv6 address according to testIPv4 and testIPv6.

        multusInterface: whether send http request to all interfaces ip in testEndpoint case.

        ipv4: test any IPv4 address. Notice, the 'enableIPv4' in configmap  spiderdocter must be enabled

        ipv6: test any IPv6 address. Notice, the 'enableIPv6' in configmap  spiderdocter must be enabled

        ingress: send http request to the ingress ipv4 or ipv6 address of kdoctor agnent

        nodePort: send http request to the nodePort ipv4 or ipv6 address with each local node of kdoctor agnent , according to testIPv4 and testIPv6.

        >notice: when test targetAgent case, it will send http request to all targets at the same time with spec.request.qps for each one. That meaning, the actually QPS may be bigger than spec.request.qps

* spec.expect: define the success condition of the task result 

    meanAccessDelayInMs: mean access delay in MS, if the actual delay is bigger than this, it results to be failure

    successRate: the success rate of all http requests. Notice, when a http response code is >=200 and < 400, it's treated as success. if the actual whole success rate is smaller than successRate, the task results to be failure

* status: the status of the task
    doneRound: how many rounds have finished

    expectedRound: how many rounds the task expect

    finish: whether all rounds of this task have finished

    lastRoundStatus: the result of last round

    history:
        roundNumber: the round number

        status: the status of this round

        startTimeStamp: when this round begins

        endTimeStamp: when this round finally finished

        duration: how long the round spent

        deadLineTimeStamp: the time deadline of a round 

        failedAgentNodeList: the node list where failed kdoctor agent locate

        notReportAgentNodeList: the node list where uknown kdoctor agent locate. This means these agents have problems.

        succeedAgentNodeList: the node list where successful kdoctor agent locate


## example 

a quick task to test kdoctor agent, to verify the whole network is ok, each agent could reach each other

```shell

cat <<EOF > netreachhealthy-test-agent.yaml
apiVersion: kdoctor.io/v1beta1
kind: NetReach
metadata:
  generation: 1
  name: netreach
spec:
  request:
    durationInSecond: 10
    perRequestTimeoutInMS: 1000
    qps: 10
  schedule:
    roundNumber: 2
    roundTimeoutMinute: 1
    schedule: 1 1
  expect:
    meanAccessDelayInMs: 10000
    successRate: 1
  target:
    clusterIP: true
    endpoint: true
    ingress: false
    ipv4: true
    ipv6: true
    loadBalancer: false
    multusInterface: false
    nodePort: true
EOF
kubectl apply -f netreachhealthy-test-agent.yaml

```

## debug

when something wrong happen, see the log for your task with following command
```shell
#get log 
CRD_KIND="netreachhealthy"
CRD_NAME="netreach"
kubectl logs -n kube-system  kdoctor-agent-v4vzx | grep -i "${CRD_KIND}.${CRD_NAME}"

```


## report

when the kdoctor is not enabled to aggerate reports, all reports will be printed in the stdout of kdoctor agent.
Use the following command to get its report
```shell
kubectl logs -n kube-system  kdoctor-agent-v4vzx | jq 'select( .TaskName=="netreachhealthy.netreach" )'
```

when the kdoctor is enabled to aggregate reports, all reports will be collected in the PVC or hostPath of kdoctor controller.


metric introduction
```shell
      {
        "FailureReason": "",
        "MeanDelay": 106.84,
        "Metrics": {
          "start": "2023-05-24T08:13:13.530015031Z",
          "end": "2023-05-24T08:13:28.560982373Z",
          "duration": "15.030967342s",
          "requestCount": 150,
          "successCount": 150,
          "tps": 9.979397638691244,
          "total_request_data": "34866 byte",
          "latencies": {
            "P50_inMs": 103,
            "P90_inMs": 199,
            "P95_inMs": 202,
            "P99_inMs": 204,
            "Max_inMx": 205,
            "Min_inMs": 3,
            "Mean_inMs": 106.84
          },
          "status_codes": {
            "200": 150
          },
          "errors": {}
        },
        "Succeed": "true",
        "SucceedRate": "1",
        "TargetMethod": "GET",
        "TargetName": "AgentClusterV4IP_172.41.156.187:80",
        "TargetUrl": "http://172.41.156.187:80"
      }
```
