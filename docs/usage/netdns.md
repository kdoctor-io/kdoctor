# netdns

[**简体中文**](./netdns-zh_CN.md) | **English**

## Introduction

kdoctor-controller creates the necessary resources, including [agent](../concepts/runtime.md), based on the agentSpec. Each agent Pod sends DNS requests to a specified DNS server. By default, the concurrency level is set to 50, which can handle scenarios with multiple replicas. The concurrency level can be configured in the kdoctor configmap. The success rate and average latency are measured, and the results are evaluated based on predefined success criteria. Detailed reports can be obtained using the aggregate API.

1. Use Cases:
    * Test the accessibility of the CoreDNS service from all corners of the cluster in production or E2E environments.
    * Adjust CoreDNS resources and replica count during application deployment to ensure it can handle the expected load.
    * Apply stress to CoreDNS for purposes like testing upgrades, chaos tests, bug reproduction, etc
    * Test external DNS services from within the cluster.

2. For more detailed information about the NetDns CRD, refer to [NetDns](../reference/netdns.md)

3. Features:

    * Support testing DNS servers both inside and outside the cluster
    * Supports typeA and typeAAAA records
    * Supports UDP, TCP, and TCP-TLS protocols

## Steps

The following example demonstrates how to use `NetDNS`.

### Install kdoctor

Follow the [installation guide](./install.md) to install kdoctor.

### Install Test Server (Optional)

The official kdoctor repository includes an application called "server" that contains an HTTP server, HTTPS server, and DNS server. This server can be employed to test the functionality of kdoctor. If you have other test servers available, you can skip this installation step.

```shell
helm repo add kdoctor https://kdoctor-io.github.io/kdoctor
helm repo update kdoctor
helm install server kdoctor/server -n kdoctor-test-server --wait --debug --create-namespace 
```

Check the status of test server

```shell
kubectl get pod -n kdoctor -owide
NAME                                READY   STATUS    RESTARTS   AGE   IP            NODE                    NOMINATED NODE   READINESS GATES
server-7649566ff9-dv4jc   1/1     Running   0          76s   172.40.1.45   kdoctor-worker          <none>           <none>
server-7649566ff9-qc5dh   1/1     Running   0          76s   172.40.0.35   kdoctor-control-plane   <none>           <none>
```

Obtain the service address of the test server

```shell
kubectl get service -n kdoctor 
NAME               TYPE        CLUSTER-IP    EXTERNAL-IP   PORT(S)                                AGE
server   ClusterIP   172.41.71.0   <none>        80/TCP,443/TCP,53/UDP,53/TCP,853/TCP   2m31s
```

### Create NetDns

Create `NetDns` object that will execute a 10-second continuous task. The task will send UDP requests to the cluster's internal DNS server at a rate of 10 QPS. It will request the typeA records whose domain name is `kubernetes.default.svc.cluster.local`.

```shell
cat <<EOF | kubectl apply -f -
apiVersion: kdoctor.io/v1beta1
kind: Netdns
metadata:
  name: netdns-cluster
spec:
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
    targetDns:
      serviceName: kube-dns
      serviceNamespace: kube-system
      testIPv4: true
EOF
```

### Check Task Status

After completing a round of tasks, you can use the kdoctor aggregate API to view the report for the current round. When the FINISH field is set to true, it indicates that all tasks have been completed, and you can access the overall report.

```shell
kubectl get netdns
NAME             FINISH   EXPECTEDROUND   DONEROUND   LASTROUNDSTATUS   SCHEDULE
netdns-cluster   true     1               1           succeed           0 1
```

* FINISH: indicate whether the task has been completed
* EXPECTEDROUND: number of expected task rounds
* DONEROUND: number of completed task rounds
* LASTROUNDSTATUS: execution status of the last round of tasks
* SCHEDULE: schedule rules for the task

### View Task Reports

1. View existed reports

    ```shell
    kubectl get kdoctorreport
    NAME             CREATED AT
    netdns-cluster   0001-01-01T00:00:00Z
    ```

2. View specific task reports

    The reports are aggregated from the agents running on both the kdoctor-control-plane node and the kdoctor-worker nodes after performing two rounds of stress testing respectively.

    ```shell
    root@kdoctor-control-plane:/# kubectl get kdoctorreport netdns-cluster -oyaml
    apiVersion: system.kdoctor.io/v1beta1
    kind: KdoctorReport
    metadata:
      creationTimestamp: null
      name: netdns-cluster
    spec:
      FailedRoundNumber: null
      FinishedRoundNumber: 1
      Report:
      - NodeName: kdoctor-control-plane
        PodName: kdoctor-agent-ntp9l
        ReportType: agent test report
        RoundDuration: 11.025723086s
        RoundNumber: 1
        RoundResult: succeed
        StartTimeStamp: "2023-08-01T09:09:39Z"
        EndTimeStamp: "2023-08-01T09:09:50Z"
        TaskName: netdns.netdns-cluster
        TaskType: Netdns
        netDNSTask:
          detail:
          - FailureReason: null
            MeanDelay: 0.2970297
            Metrics:
              DNSMethod: udp
              DNSServer: 172.41.0.10:53
              Duration: 11.002666395s
              EndTime: "2023-08-01T09:09:50Z"
              Errors: {}
              FailedCounts: 0
              Latencies:
                MaxInMs: 0
                MeanInMs: 0.2970297
                MinInMs: 0
                P50InMs: 0
                P90InMs: 0
                P95InMs: 0
                P99InMs: 0
              ReplyCode:
                NOERROR: 101
              RequestCounts: 101
              StartTime: "2023-08-01T09:09:39Z"
              SuccessCounts: 101
              TPS: 9.179593052634765
              TargetDomain: kubernetes.default.svc.cluster.local.
            Succeed: true
            SucceedRate: 1
            TargetName: typeA_172.41.0.10:53_kubernetes.default.svc.cluster.local
            TargetProtocol: udp
            TargetServer: 172.41.0.10:53
          succeed: true
          targetNumber: 1
          targetType: kdoctor agent
          MaxCPU: 30.651%
          MaxMemory: 97.00MB
        netDNSTaskSpec:
          ...
      - NodeName: kdoctor-worker
        PodName: kdoctor-agent-krrnp
        ReportType: agent test report
        RoundDuration: 10.024533428s
        RoundNumber: 1
        RoundResult: succeed
        StartTimeStamp: "2023-08-01T09:09:39Z"
        EndTimeStamp: "2023-08-01T09:09:49Z"
        TaskName: netdns.netdns-cluster
        TaskType: Netdns
        netDNSTask:
          detail:
          - FailureReason: null
            MeanDelay: 0.58
            Metrics:
              ...
            Succeed: true
            SucceedRate: 1
            TargetName: typeA_172.41.0.10:53_kubernetes.default.svc.cluster.local
            TargetProtocol: udp
            TargetServer: 172.41.0.10:53
          succeed: true
          targetNumber: 1
          targetType: kdoctor agent
          MaxCPU: 30.651%
          MaxMemory: 97.00MB
        netDNSTaskSpec:
          ...
      ReportRoundNumber: 1
      RoundNumber: 1
      Status: Finished
      TaskName: netdns-cluster
      TaskType: Netdns
    ```

> If the reports do not align with the expected results, check the MaxCPU and MaxMemory fields in the report to verify if there are available resources of the agents and adjust the resource limits for the agents accordingly.

## Test an External DNS Server

Below are examples of HTTP and HTTPS requests with body:

1. Create `NetDns` object that will execute a 10-second continuous task. The task will send UDP requests to the cluster's internal DNS server at a rate of 10 QPS. It will request the typeA records whose domain name is `kubernetes.default.svc.cluster.local`.

    We are using the service address of the test server. If you have a different server address available, feel free to use it instead.

    Create `NetDns`

    ```shell
    SERVER="172.41.71.0"
    apiVersion: kdoctor.io/v1beta1
    kind: Netdns
    metadata:
      name: netdns- user
    spec:
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
        targetUser:
          port: 53
          server: ${SERVER}
    EOF
    ```

## Environment Cleanup

```shell
kubectl delete netdns netdns-cluster netdns- user
```
