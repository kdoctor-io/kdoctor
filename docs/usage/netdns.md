# netdns

## concept

Fo this kind task, each kdoctor agent will send dns request to specified target, and get success rate and mean delay.
It could specify success condition to tell the result succeed or fail.
And, more detailed report will print to kdoctor agent stdout, or save to disc by kdoctor controller.

the following is the spec of netdns

```shell

cat <<EOF > netdns.yaml
apiVersion: kdoctor.io/v1beta1
kind: Netdns
metadata:
  name: testdns
spec:
  schedule:
    schedule: "1 1"
    roundNumber: 2
    roundTimeoutMinute: 1
  target:
    targetDns:
      testIPv4: true
      testIPv6: false
      serviceName: coredns
      serviceNamespace: kube-system
    targetUser:
      server: 172.18.0.1
      port: 53
  request:
    durationInSecond: 10
    qps: 20
    perRequestTimeoutInMS: 500
    domain: "kube-dns.kube-system.svc.cluster.local"
    protocol: udp
  expect:
    successRate: 1
    meanAccessDelayInMs: 10000
EOF

kubectl apply -f netdns.yaml

```

* spec.schedule: set how to schedule the task.

      roundNumber: how many rounds it should be to run this task

      schedule: Support Linux crontab syntax for scheduling tasks, while also supporting simple writing. 
                The first digit represents how long the task will start, and the second digit represents the interval time between each round of tasks,
                separated by spaces. Example: "1 2" indicates that the task will start in 1 minute, and the interval time between each round of tasks.

      roundTimeoutMinute: the timeout in minute for each round, when the rask does not finish in time, it results to be failuire

      sourceAgentNodeSelector [optional]: set the node label selector, then, the kdoctor agent who locates on these nodes will implement the task. If not set this field, all kdoctor agent will execute the task

* spec.request: how each kdoctor agent should send the dns request

      durationInSecond: for each round, the duration in second how long the dns request lasts

      perRequestTimeoutInMS: timeout in ms for each dns request

      qps: qps
     
      domain： resolved domain

* spec.target: set the target of dns request. it could not set targetUser and targetDns at the same time

      targetUser [optional]: set an user-defined DNS server for the dns request

        server: the address for dns server

        port: the port for dns server

      targetDns: [optional]: set cluster dns server for the dns request

        testIPv4: test DNS server IPv4 address and request is type A. 

        testIPv6: test DNS server IPv6 address and request is type AAAA.

        serviceName: Specify the name of the DNS to be tested
* 
        serviceNamespace: Specify the namespace of the DNS to be tested
 
      protocol: Specify request protocol,Optional value udp，tcp，tcp-tls,default udp.

* spec.expect: define the success condition of the task result

  meanAccessDelayInMs: mean access delay in MS, if the actual delay is bigger than this, it results to be failure

  successRate: the success rate of all dns requests. Notice, when a dns response code is >=200 and < 400, it's treated as success. if the actual whole success rate is smaller than successRate, the task results to be failure

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

test custom dns server by crontab

```shell

cat <<EOF > netdns1.yaml
apiVersion: kdoctor.io/v1beta1
kind: Netdns
metadata:
  name: testdns
spec:
  schedule:
    schedule: "*/1 * * * *"
    roundNumber: 2
    roundTimeoutMinute: 1
  target:
    targetUser:
      server: 172.18.0.1
      port: 53
  request:
    durationInSecond: 10
    qps: 10
    perRequestTimeoutInMS: 500
    domain: "baidu.com"
    protocol: udp
  expect:
    successRate: 1
    meanAccessDelayInMs: 1000
EOF

kubectl apply -f netdns1.yaml

```

test custom dns server by simple

```shell

cat <<EOF > netdns1.yaml
apiVersion: kdoctor.io/v1beta1
kind: Netdns
metadata:
  name: testdns
spec:
  schedule:
    schedule: "1 1"
    roundNumber: 2
    roundTimeoutMinute: 1
  target:
    protocol: udp
    targetUser:
      server: 172.18.0.1
      port: 53
  request:
    durationInSecond: 10
    qps: 10
    perRequestTimeoutInMS: 500
    domain: "baidu.com"
  expect:
    successRate: 1
    meanAccessDelayInMs: 1000
EOF

kubectl apply -f netdns1.yaml

```

test cluster dns server by crontab

```shell

cat <<EOF > netdns.yaml
apiVersion: kdoctor.io/v1beta1
kind: Netdns
metadata:
  name: testdns
spec:
  schedule:
    schedule: "*/1 * * * *"
    roundNumber: 2
    roundTimeoutMinute: 1
  target:
    targetDns:
      testIPv4: true
      testIPv6: false
      serviceNamespaceName: kube-system/kube-dns
    protocol: udp
  request:
    durationInSecond: 10
    qps: 20
    perRequestTimeoutInMS: 500
    domain: "kube-dns.kube-system.svc.cluster.local"
  expect:
    successRate: 1
    meanAccessDelayInMs: 10000
EOF

kubectl apply -f netdns.yaml

```

test cluster dns server by simple

```shell

cat <<EOF > netdns.yaml
apiVersion: kdoctor.io/v1beta1
kind: Netdns
metadata:
  name: testdns
spec:
  schedule:
    schedule: "1 1"
    roundNumber: 2
    roundTimeoutMinute: 1
  target:
    targetDns:
      testIPv4: true
      testIPv6: false
      serviceNamespaceName: kube-system/test-app
    protocol: udp
  request:
    durationInSecond: 10
    qps: 20
    perRequestTimeoutInMS: 500
    domain: "kube-dns.kube-system.svc.cluster.local"
  expect:
    successRate: 1
    meanAccessDelayInMs: 10000
EOF

kubectl apply -f netdns.yaml

```


## report

when the kdoctor is not enabled to aggerate reports, all reports will be printed in the stdout of kdoctor agent.
Use the following command to get its report
```shell
kubectl logs -n kube-system  kdoctor-agent-lwhtm | jq 'select( .TaskName=="netdns.testdns" )'
```

when the kdoctor is enabled to aggregate reports, all reports will be collected in the PVC or hostPath of kdoctor controller.


metric introduction
```json
{
  "TaskName": "netdns.testdns",
  "TaskSpec": {
    "schedule": {
      "schedule": "1 1",
      "roundTimeoutMinute": 1,
      "roundNumber": 2
    },
    "target": {
      "protocol": "tcp"
    },
    "request": {
      "durationInSecond": 10,
      "qps": 20,
      "perRequestTimeoutInMS": 500,
      "domain": "kube-dns.kube-system.svc.cluster.local"
    },
    "success": {
      "successRate": 1,
      "meanAccessDelayInMs": 10000
    }
  },
  "RoundNumber": 1,
  "RoundResult": "succeed",
  "NodeName": "kdoctor-control-plane",
  "PodName": "kdoctor-agent-lwhtm",
  "FailedReason": "",
  "StartTimeStamp": "2023-04-27T07:07:32.032814878Z",
  "EndTimeStamp": "2023-04-27T07:07:32.070513569Z",
  "RoundDuraiton": "37.69869ms",
  "ReportType": "agent test report",
  "Detail": {}
}
```