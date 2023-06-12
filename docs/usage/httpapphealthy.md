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
    body: kube-system/http-body
    header:
    - Accept:text/html
    host: https://10.6.172.20:9443
    http2: false
    method: PUT
    tls-secret: kube-system/https-cert
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

        host: the host for http, example service ip, pod ip, service domain, an user-defined UR

        method: http method, must be one of GET POST PUT DELETE CONNECT OPTIONS PATCH HEAD
        
        body: The configmap format for logging HTTP requests is namespace/configmap-name
        
        tls-cert: The secret format for logging HTTPS request certificates is namespace/configmap-name

        header:  HTTP request header

        http2: Requests are made using the http2 protocol

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
    body: kube-system/http-body
    header:
    - Accept:text/html
    host: https://10.6.172.20:9443
    http2: false
    method: PUT
    tls-secret: kube-system/https-cert
EOF
kubectl apply -f test-httpapphealthy.yaml

```


### example body
```shell
cat <<EOF > http-body.yaml
apiVersion: v1
data:
  body: |
    {test:test}
kind: ConfigMap
metadata:
  name: http-body
  namespace: kube-system
EOF
kubectl apply -f http-body.yaml
```

### example https cert
```shell
cat <<EOF > https-cert.yaml
apiVersion: v1
data:
  ca.crt: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURWekNDQWorZ0F3SUJBZ0lKQU1vL3p5bGZZZzVSTUEwR0NTcUdTSWIzRFFFQkN3VUFNRUl4Q3pBSkJnTlYKQkFZVEFsaFlNUlV3RXdZRFZRUUhEQXhFWldaaGRXeDBJRU5wZEhreEhEQWFCZ05WQkFvTUUwUmxabUYxYkhRZwpRMjl0Y0dGdWVTQk1kR1F3SGhjTk1qTXdOakE0TURrd05URTVXaGNOTWpRd05qQTNNRGt3TlRFNVdqQkNNUXN3CkNRWURWUVFHRXdKWVdERVZNQk1HQTFVRUJ3d01SR1ZtWVhWc2RDQkRhWFI1TVJ3d0dnWURWUVFLREJORVpXWmgKZFd4MElFTnZiWEJoYm5rZ1RIUmtNSUlCSWpBTkJna3Foa2lHOXcwQkFRRUZBQU9DQVE4QU1JSUJDZ0tDQVFFQQpwWERBUGt1UzhvNW9lRTBBS0ZxL2Vjb1pjN2hFbXk1RlMvbWxlZCt2MFRFMlV5cGord1k0M0hvaEhpWjl3bVRECkpwTHZTKzgwTFpmMitrVkNBb05hTjdMdU1rdFJKaXlQUDc2TklSaUdPdzl6MmZpNHNIaUFnS0dvR1ZMb1c1YUMKa0RoK3dLKzh5NDVnVGZGR1VaWGpBa0pKSm1mWDd2TXllbkpyT2J5SUE0ajFuc294cDBNelFzNkYzREQ1TmdmZApQc1JtT3N6QlNLRTdNaFF4MEN5RlVQWjRTZ3U2N25MQytmRFBWVXBYY3pKQU1ZSTVrNWpyaElneDZnR2hKVFk0CmRLQ0VMWllwUmhCMWFFbVBIRjFlVUZ3MC9FcG5ldUdPd1ZqazZsSEp6QUxRUHBnR1dBZ0V1WFFVckxYb0dNclAKcWJrYU9WeitMelh0N1ZCaWJOZmdFUUlEQVFBQm8xQXdUakFkQmdOVkhRNEVGZ1FVTDlnL3FhZ2ptaGJ1K1pvQQpRVkFOdE1nd0cra3dId1lEVlIwakJCZ3dGb0FVTDlnL3FhZ2ptaGJ1K1pvQVFWQU50TWd3Rytrd0RBWURWUjBUCkJBVXdBd0VCL3pBTkJna3Foa2lHOXcwQkFRc0ZBQU9DQVFFQW5Jd0Jyc3paY0pRRGFrZnNVVHp1eEtORmdUNUoKNVJOUWQ1S1NZM01HTWVkQjNIU2dMTEVWM3FTM1pLNHlrT3NFaXJ6c1lmdGV2Q1JGL2VsTkVIZDREQ2tiRzlSeQpwTzAwYTk1RkdFNktUWk1iSTRDaW1kMmF6eVhkazMyYUtzeFpjaDBRMzlneHgwSGFnVW5CcDk1VWxaQkcyOTh2CnhkQVRtSXZJNmVpd2FCeWc2NjBKRklRZ1VkVmhUM0VNTUI2dUlxaWdKZTJlMEcvV0t1My9BNGhvL3hDZUVNWmYKbFJlRXJIeFo0TzZsVTVEM0pnNURWcUM1MlBPK0V2aTJpdm5ZTmNvbTRibU9HbE41RmhzdG5La3M2ZlVsV012RQpXeXZCb01OU2VFNllueVFUUURBZ09BN0NlVzhKSUk4b3JRN00vM0JnWlNSOUZ4OGdhY01jeEUydTZRPT0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=
  tls.crt: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURBRENDQWVnQ0NRRGljN284NjltaW16QU5CZ2txaGtpRzl3MEJBUXNGQURCQ01Rc3dDUVlEVlFRR0V3SlkKV0RFVk1CTUdBMVVFQnd3TVJHVm1ZWFZzZENCRGFYUjVNUnd3R2dZRFZRUUtEQk5FWldaaGRXeDBJRU52YlhCaApibmtnVEhSa01CNFhEVEl6TURZd09EQTVNak13TlZvWERUSTBNRFl3TnpBNU1qTXdOVm93UWpFTE1Ba0dBMVVFCkJoTUNXRmd4RlRBVEJnTlZCQWNNREVSbFptRjFiSFFnUTJsMGVURWNNQm9HQTFVRUNnd1RSR1ZtWVhWc2RDQkQKYjIxd1lXNTVJRXgwWkRDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUJBTm5RUnVVaApnbWhtdFV4aitJQXpIayswdU5Oa2ZnUWFVMHpYdG1ZMGI3T0g3VkZXd2tVWDNJYVU2Q2JJRzd2MVRPbllISmZhCloyZmRvNURudlpaSk5OWnZOWUdEbkU4d3JuNmVHSlFpMW5hbDlEbGdyaWJPSWJtOUU4N055VUFTRXpkSEFhN1oKOUp1bEVlRlE5VXJHV25KbXhEZjBTSGM4ckJGRnBsMVpORXpUbDFpUFFhaDRDVGcya1lWMWJtS0h2cmhUcDV6UQp4VE4wT0ZUa1JKWmprN284YTRxUXhBVWYvU1ZkQ3BkeXBHYzhPM3JWL0dPMTdneFlVK2lmRmY0OW84M2l5ZGluCm0zQ3NmOFRFVWlNTGZqcDgwQk5ZQTVScVJPSG15NTVNSG9TM2VnWEJWOG5va3hvVHkzSy8rdnd0L21BMmVLQWEKVHRXeEo1NVM2T0M5R0xrQ0F3RUFBVEFOQmdrcWhraUc5dzBCQVFzRkFBT0NBUUVBalYrWWpKSkU1c2o5ODZqQgpzenY3cXhkdXRGSHdIM1NpRXhzaEt0S3VpaE1BYjBJR1V4MEIxbmdTZ1gvNUVqUzEwYTZtRmhJZGxWS29PSDJ1CkQ3aVZtRkdweHBWaUtCTFMwVnhweGphVmZ0OExCd2k1cHN5eDZyWmEwaFMvek1NMEFlL1FuQXpoSzZDMDl5T08KN1g1R2orMjNQQjBVNkorZnRteThMYVpwK0ZBbWFobi9OYThJbmJNY1hEQjVEeEhNUWkzdjFrQUh6bnBGNU02KwpuSEkyR3B0RzR4UzlDTitFK2FBa3NBZGMzY2VjZ3JCL04vUFZNMWhFdkZtakw5SVpTdEJkYzBGQ1pGMHJPVy96CllhWVhLa3FRTm9Wa2FaMENsRVEybWdIMFh0ZktzQ0VZVGx0OGJncDgveDdKTmlqN2UvYkoyU0E0M015NTF1ZHYKZ0N5M1pRPT0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=
  tls.key: LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFb3dJQkFBS0NBUUVBMmRCRzVTR0NhR2ExVEdQNGdETWVUN1M0MDJSK0JCcFRUTmUyWmpSdnM0ZnRVVmJDClJSZmNocFRvSnNnYnUvVk02ZGdjbDlwblo5MmprT2U5bGtrMDFtODFnWU9jVHpDdWZwNFlsQ0xXZHFYME9XQ3UKSnM0aHViMFR6czNKUUJJVE4wY0JydG4wbTZVUjRWRDFTc1phY21iRU4vUklkenlzRVVXbVhWazBUTk9YV0k5QgpxSGdKT0RhUmhYVnVZb2UrdUZPbm5OREZNM1E0Vk9SRWxtT1R1anhyaXBERUJSLzlKVjBLbDNLa1p6dzdldFg4Clk3WHVERmhUNko4Vi9qMmp6ZUxKMktlYmNLeC94TVJTSXd0K09uelFFMWdEbEdwRTRlYkxua3dlaExkNkJjRlgKeWVpVEdoUExjci82L0MzK1lEWjRvQnBPMWJFbm5sTG80TDBZdVFJREFRQUJBb0lCQUM4QVVLd1ZCUXorVE5VRgpKWlNVYzFBRDBYWmNVdzBUbVRJVndsaGZyRkx6Vy9TWFlpaUNzNldlOEZHZUVNNElhdVp6S2doaXFybXhEQ0N5CndTaHk5NkhtTVllWEhOM0J4WVd4RytDcmU5ZnlpN2J0OCthUHlKdEovOEk2aWRqM2pZbjZHcFRlbDNnV3NMc00KTzBJOWR6c0VqZ2I5QWI0cEs0QTJwV1d6WUNQTGh3cVVBTXBOS0tSc215NkQzWlozL2lHMmJ6TVRzdlpUUTQ1NQpZS1RCRkxzZGV2N1NLd3UxbWtUQzQrWkFPUjFqMFNSemVWdmFTSnczY2pWVm00SnRod3c1M1o4TVRzSFBIT3dLCmlGM1dlazdaK3BmZkpSTE91T2h3VnBzcTRlVnpCYXhXeW1QNkg0UmQxc1YwQWRGM3QwQlRmRnVudXBCa2RkM1kKUzY2UG5RRUNnWUVBNytCQXpjNDJxaDluQVgwekNhWGZlQU5GOG5OOFp4TTRUTXZEL0I3VWl2Ty9VNWR0Rm15Rgo4emhMbWpXZ3R4ZTdraFRma3R6czJ0TVc4UmhpRlZ0bGRtWForT0pqWFJtR1pwMkhLcmJ2R1pxTXdyLzFjdzJsClhkaTIwZHU0Slp1emxMM3ArSzhFMnRCdEFNK0tyZ09yUy9iTkVPRjg4NDFPUDhCRmIxNmtMakVDZ1lFQTZIUmcKbzNTaGpJcHU1NDNDOCtzbXNNRTFPZVVheEJ6YzdCNmg4a010SEhlRi8yOGVKQ3JJV0NnaDVtTWZCR3AxTmlmSgo2N0ZlT21NM0QyS0JzK3N5cVVQMVlKcVFCY0JKeUVuUzR5SnJpR2diQ1ZpUkpHdm95VGdEV21YbmpHKzI1N2ZHCnZaV0tSTGNKUHhWR0NrckdZc01BK3R1YzJJTE5LUkRPMlpRanlRa0NnWUE5YjVFSlpPUkpSQXVzclBVeVptSksKcVlQenFiSlY3KzArZGYyM0IrcGx3REhqWmVnUmt5L25jQ2FrMDFGYk0xL2Q5U3lodjZXR0VnUlJNVzZGaThmNwp2L0JJdHlxOXdIalV0VW5XSGM0MUg0a25vK1JvV0RsZlJNN21Cc0V1R0tldzA4Y2w0eVY2S1dHUmtKWXpKVXR0CkJFUFhLL2xGbzQ1RDg2bVU4WWRaTVFLQmdRRFN1UDBKOENhcWtxdXErUldyckpYc1VabUFuREhCYWpEVFU0bVgKWmxJMHBoMHd5M2hWYlBzay8yeUx2M3RVczNVQjNOdnM3Mkx1SnhhNHVhRytpZzNvNTVRL09KNHF1SCtxTTFJYgpXUTZHSDJteTlUak4vWXlQTEZuTnp1Y3lwZXIyNytBWDZNSHBQTXdEQmJQeWpJcCs2U3V3UFBsWVJHcmJPVU5xCmRpSmlrUUtCZ0Y0dWtWazFjRkFJazlOeThpNXlrNHR2QzY1SXk1dkQvYWFMeE0yWFo5dnN2TTc5TkNzbUp4VkwKYnlTeWxJbi9rWnFFT0tkRHkxZnRYWnY1aGsrUWhvUzNsT0xWTjlUZ2k0Unhqcnl6QmJDYXdQNjlBZmxxN3dsdwpJYzZFNlZncXJhOU52Q0Zxbm1PMzBaQ1NteENBUEJ3d2hqQmN1K1JEMFVxT0ZGMXZwckNlCi0tLS0tRU5EIFJTQSBQUklWQVRFIEtFWS0tLS0tCg==
kind: Secret
metadata:
  name: https-cert
  namespace: kube-system
type: kubernetes.io/tls
EOF
kubectl apply -f https-cert.yaml
```

> notice: key body ca.crt tls.crt and tls.key are fixed field cannot be customized

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
