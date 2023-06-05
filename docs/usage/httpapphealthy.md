# Nethttp

## concept 

Fo this kind task, each kdoctor agent will send http request to specified target, and get success rate and mean delay. 
It could specify success condition to tell the result succeed or fail. 
And, more detailed report will print to kdoctor agent stdout, or save to disc by kdoctor controller.

the following is the spec of nethttp
```shell
apiVersion: kdoctor.io/v1beta1
kind: HttpAppHealthy
metadata:
  creationTimestamp: "2023-05-24T08:00:05Z"
  generation: 1
  name: httphealthy
  resourceVersion: "1426047"
  uid: baebfc2a-3dbc-4df6-b64e-f5066f82fdd6
spec:
  request:
    durationInSecond: 10
    perRequestTimeoutInMS: 1000
    qps: 10
  schedule:
    roundNumber: 2
    roundTimeoutMinute: 1
    schedule: 1 1
  success:
    meanAccessDelayInMs: 10000
    successRate: 1
  target:
    body: default/httpbody
    header:
    - Accept:text/html
    host: http://172.41.156.187
    http2: false
    method: GET
    tls-ca: default/httptls
status:
  doneRound: 2
  expectedRound: 2
  finish: true
  history:
  - deadLineTimeStamp: "2023-05-24T08:03:05Z"
    duration: 15.092667522s
    endTimeStamp: "2023-05-24T08:02:20Z"
    expectedActorNumber: 2
    failedAgentNodeList: []
    notReportAgentNodeList: []
    roundNumber: 2
    startTimeStamp: "2023-05-24T08:02:05Z"
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

* spec.target: set the target of http request.

      targetUser [optional]: set an user-defined URL for the http request

        host: the host for http, example service ip, pod ip, service domain, an user-defined UR

        method: http method, must be one of GET POST PUT DELETE CONNECT OPTIONS PATCH HEAD

        >notice: when test targetAgent case, it will send http request to all targets at the same time with spec.request.qps for each one. That meaning, the actually QPS may be bigger than spec.request.qps

* spec.success: define the success condition of the task result 

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

a quick task to test kdoctor agent, to verify the whole network is ok, each agent could reach specific host

```shell

cat <<EOF > test-httpapphealthy.yaml
apiVersion: kdoctor.io/v1beta1
kind: HttpAppHealthy
metadata:
  name: httphealthy
spec:
  request:
    durationInSecond: 10
    perRequestTimeoutInMS: 1000
    qps: 10
  schedule:
    roundNumber: 2
    roundTimeoutMinute: 1
    schedule: 1 1
  success:
    meanAccessDelayInMs: 10000
    successRate: 1
  target:
    body: default/httpbody
    header:
    - Accept:text/html
    host: http://kdoctor-agent-ipv4.kube-system.svc.cluster.local
    http2: false
    method: GET
    tls-ca: default/httptls
EOF
kubectl apply -f test-httpapphealthy.yaml

```

## debug

when something wrong happen, see the log for your task with following command
```shell
#get log 
CRD_KIND="httpapphealthy"
CRD_NAME="httphealthy"
kubectl logs -n kube-system  kdoctor-agent-v4vzx | grep -i "${CRD_KIND}.${CRD_NAME}"

```


## report

when the kdoctor is not enabled to aggerate reports, all reports will be printed in the stdout of kdoctor agent.
Use the following command to get its report
```shell
kubectl logs -n kube-system  kdoctor-agent-v4vzx | jq 'select( .TaskName=="httpapphealthy.httphealthy" )'
```

when the kdoctor is enabled to aggregate reports, all reports will be collected in the PVC or hostPath of kdoctor controller.


metric introduction
```shell
		"FailureReason": "",
		"MeanDelay": 34.36,
		"Metrics": {
			"start": "2023-05-24T09:00:35.03987095Z",
			"end": "2023-05-24T09:00:45.08774646Z",
			"duration": "10.04787551s",
			"requestCount": 100,
			"successCount": 100,
			"tps": 9.952352604336754,
			"total_request_data": "23247 byte",
			"latencies": {
				"P50_inMs": 35,
				"P90_inMs": 57,
				"P95_inMs": 58,
				"P99_inMs": 59,
				"Max_inMx": 62,
				"Min_inMs": 18,
				"Mean_inMs": 34.36
			},
			"status_codes": {
				"200": 100
			},
			"errors": {}
		},
		"Succeed": "true",
		"SucceedRate": "1",
		"TargetMethod": "GET",
		"TargetName": "HttpAppHealthy target",
		"TargetNumber": "1",
		"TargetType": "HttpAppHealthy",
		"TargetUrl": "http://kdoctor-agent-ipv4.kube-system.svc.cluster.local"
```
