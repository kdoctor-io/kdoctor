# Nethttp

[**简体中文**](./netreach-zh_CN.md) | **English**

## Introduction

kdoctor-controller creates the necessary resources, including [agent](../concepts/runtime.md), based on the agentSpec. Each agent Pod sends HTTP requests to each other with a certain level of stress. These requests target various addresses such as Pod IP, service IP, ingress IP, and more. The success rate and average latency are measured, and the results are evaluated based on predefined success criteria. Detailed reports can be obtained using the aggregate API.

1. Use cases:

    * Monitor connectivity in every corner of the cluster with a 1-minute heartbeat interval.
    * Perform connectivity inspections across all corners of the cluster after deploying at scale
    * Inject traffic to all corners of the cluster using different communication methods for scenarios like bug reproduction
    * Inspect anomalies in components such as CNI, LoadBalancer, kube-proxy, Ingress, Multus, etc., in production or E2E environments.

2. For more information on NetReach CRD, refer to[NetReach](../reference/netreach.md)

3. Features:

    Support ClusterIP, Endpoint, Ingress, NodePort, LoadBalancer, Multus with multiple network interfaces, IPv4, and IPv6.

## Steps

The following example demonstrates how to use `NetReach`.

### Install kdoctor

Follow the [installation guide](./install.md) to install kdoctor.

### Create NetReach

Create `NetReach` object that will execute a 10-second continuous task. Each agent on the nodes will use the HTTP protocol to access the IPv4 addresses of ClusterIP, Endpoint, Ingress, NodePort, and LoadBalancer immediately,Report name consists of `${TaskKind}-${TaskName}`

```shell
cat <<EOF | kubectl apply -f -
apiVersion: kdoctor.io/v1beta1
kind: NetReach
metadata:
  name: task
spec:
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
    endpoint: true
    ingress: true
    ipv4: true
    loadBalancer: true
    multusInterface: false
    nodePort: true
EOF
```

### Check Task Status

After completing a round of tasks, you can use the kdoctor aggregate API to view the report for the current round. When FINISH is true, it indicates that all tasks have been completed, and you can view the overall report.

```shell
kubectl get netreach
NAME   FINISH   EXPECTEDROUND   DONEROUND   LASTROUNDSTATUS   SCHEDULE
task   true     1               1           succeed           0 1
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
    NAME           CREATED AT
    neteach-task   0001-01-01T00:00:00Z
    ```

2. View specific task reports

   The reports are aggregated from the agents running on both the kdoctor-control-plane node and the kdoctor-worker nodes after performing two rounds of stress testing respectively.

    ```shell
    root@kdoctor-control-plane:/# kubectl get kdoctorreport neteach-task -oyaml
    apiVersion: system.kdoctor.io/v1beta1
    kind: KdoctorReport
    metadata:
      creationTimestamp: null
      name: neteach-task
    spec:
      FailedRoundNumber: null
      FinishedRoundNumber: 1
      Report:
      - NodeName: kdoctor-control-plane
        NetReachTask:
          Detail:
          - TargetName: AgentLoadbalancerV4IP_172.18.0.51:80
            Metrics:
              Duration: 10.032286878s
              EndTime: "2023-08-01T08:37:06Z"
              Errors: {}
              Latencies:
                MaxInMs: 0
                MeanInMs: 23.08
                MinInMs: 0
                P50InMs: 0
                P90InMs: 0
                P95InMs: 0
                P99InMs: 0
              RequestCounts: 100
              StartTime: "2023-08-01T08:36:56Z"
              StatusCodes:
                "200": 100
              SuccessCounts: 100
              TPS: 9.967817030760152
              TotalDataSize: 36968 byte
            Succeed: true
            SucceedRate: 1
            TargetMethod: GET
            TargetUrl: http://172.18.0.51:80
            MeanDelay: 23.08
          - TargetName: AgentNodePortV4IP_172.18.0.3_32713
            Metrics:
              ...
            Succeed: true
            SucceedRate: 1
            TargetMethod: GET
            TargetUrl: http://172.18.0.3:32713
            MeanDelay: 68.42
          - TargetName: AgentPodV4IP_kdoctor-agent-ntp9l_172.40.0.6
            Metrics:
              ...
            Succeed: true
            SucceedRate: 1
            TargetMethod: GET
            TargetUrl: http://172.40.0.6:80
            MeanDelay: 44.049503
          - TargetName: AgentClusterV4IP_172.41.249.6:80
            Metrics:
              ...
            Succeed: true
            SucceedRate: 1
            TargetMethod: GET
            TargetUrl: http://172.41.249.6:80
            MeanDelay: 26.307692
          - TargetName: AgentPodV4IP_kdoctor-agent-krrnp_172.40.1.5
            Metrics:
              ...
            Succeed: true
            SucceedRate: 1
            TargetMethod: GET
            TargetUrl: http://172.40.1.5:80
            MeanDelay: 61.564358
          - TargetName: AgentIngress_http://172.18.0.50/kdoctoragent
            Metrics:
              ...
            Succeed: true
            SucceedRate: 1
            TargetMethod: GET
            TargetUrl: http://172.18.0.50/kdoctoragent
            MeanDelay: 65.47059
          Succeed: true
          TargetNumber: 6
          TargetType: NetReach
          MaxCPU: 26.203%
          MaxMemory: 101.00MB
        NetReachTaskSpec:
        ...
        PodName: kdoctor-agent-ntp9l
        ReportType: agent test report
        RoundDuration: 11.178657432s
        RoundNumber: 1
        RoundResult: succeed
        StartTimeStamp: "2023-08-01T08:36:56Z"
        EndTimeStamp: "2023-08-01T08:37:07Z"
        TaskName: netreach.task
        TaskType: NetReach
      - NodeName: kdoctor-worker
        NetReachTask:
          Detail:
          - TargetName: AgentClusterV4IP_172.41.249.6:80
            Metrics:
              ...
            Succeed: true
            SucceedRate: 1
            TargetMethod: GET
            TargetUrl: http://172.41.249.6:80
            MeanDelay: 47.25
          - TargetName: AgentNodePortV4IP_172.18.0.2_32713
            Metrics:
              ...
            Succeed: true
            SucceedRate: 1
            TargetMethod: GET
            TargetUrl: http://172.18.0.2:32713
            MeanDelay: 13.480392
          - TargetName: AgentPodV4IP_kdoctor-agent-krrnp_172.40.1.5
            Metrics:
              ...
            Succeed: true
            SucceedRate: 1
            TargetMethod: GET
            TargetUrl: http://172.40.1.5:80
            MeanDelay: 39.637257
          - TargetName: AgentPodV4IP_kdoctor-agent-ntp9l_172.40.0.6
            Metrics:
              ...
            Succeed: true
            SucceedRate: 1
            TargetMethod: GET
            TargetUrl: http://172.40.0.6:80
            MeanDelay: 51.38614
          - TargetName: AgentLoadbalancerV4IP_172.18.0.51:80
            Metrics:
              ...
            Succeed: true
            SucceedRate: 1
            TargetMethod: GET
            TargetUrl: http://172.18.0.51:80
            MeanDelay: 41.735847
            - TargetName: AgentIngress_http://172.18.0.50/kdoctoragent
            Metrics:
              ...
            Succeed: true
            SucceedRate: 1
            TargetMethod: GET
            TargetUrl: http://172.18.0.50/kdoctoragent
            MeanDelay: 60.463634
          Succeed: true
          TargetNumber: 6
          TargetType: NetReach
          MaxCPU: 30.651%
          MaxMemory: 97.00MB
        NetReachTaskSpec:
        ...
        PodName: kdoctor-agent-krrnp
        ReportType: agent test report
        RoundDuration: 11.180813761s
        RoundNumber: 1
        RoundResult: succeed
        StartTimeStamp: "2023-08-01T08:36:56Z"
        EndTimeStamp: "2023-08-01T08:37:07Z"
        TaskName: netreach.task
        TaskType: NetReach
      ReportRoundNumber: 1
      RoundNumber: 1
      Status: Finished
      TaskName: task
      TaskType: NetReach
    ```

> If the reports do not align with the expected results, check the MaxCPU and MaxMemory fields in the report to verify if there are available resources of the agents and adjust the resource limits for the agents accordingly.

## Environment Cleanup

```shell
kubectl delete netreach task
```
