# Nethttp

[**简体中文**](./apphttphealthy-zh_CN.md) | **English**

## Introduction

kdoctor-controller creates the necessary resources, including [agent](../concepts/runtime.md), based on the agentSpec. Each agent Pod sends DNS requests to a specified DNS server. By default, the concurrency level is set to 50, which can handle scenarios with multiple replicas. The concurrency level can be configured in the kdoctor configmap. The success rate and average latency are measured, and the results are evaluated based on predefined success criteria. Detailed reports can be obtained using the aggregate API.

1. Use cases:

    * Test connectivity to ensure that a specific application can be accessed from every corner of the cluster.
    * Conduct large-scale cluster testing by simulating a higher number of clients to generate increased pressure and assess the application's resilience. Simulate more source IPs to create additional application sessions and test resource limitations.
    * Inject pressure into a specific application for purposes such as gray release, chaos tests, bug reproduction, etc.
    * Test external services of the cluster to verify the proper functioning of cluster egress.

2. For a more detailed description of the AppHttpHealthy CRD, please refer to[AppHttpHealthy](../reference/apphttphealthy.md)

3. Features

    * Support HTTP, HTTPS, and HTTP2 protocols, allowing customization of headers and bodies.

## Steps

The following example demonstrates how to use `AppHttpHealthy`.

### Install kdoctor 

Follow the [installation guide](./install.md) to install kdoctor.

### Install Test Server (Optional)

The official kdoctor repository includes an application called "server" that contains an HTTP server, HTTPS server, and DNS server. This server can be employed to test the functionality of kdoctor. If you have other test servers available, you can skip this installation step.

```shell
helm repo add kdoctor https://kdoctor-io.github.io/kdoctor
helm repo update kdoctor
helm install server kdoctor/server -n kdoctor --wait --debug --create-namespace 
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

### Create AppHttpHealthy 

Create an `AppHttpHealthy` task for HTTP that will run continuously for 10 seconds. The task will send GET requests to the specified server at a rate of 10 QPS and be executed immediately.

We are using the service address of the test server. If you have a different server address available, feel free to use it instead.

```shell
SERVER="172.41.71.0"
cat <<EOF | kubectl apply -f -
apiVersion: kdoctor.io/v1beta1
kind: AppHttpHealthy
metadata:
  name: http1
spec:
  request:
    durationInSecond: 10
    perRequestTimeoutInMS: 1000
    qps: 10
  schedule:
    roundNumber: 1
    roundTimeoutMinute: 1
    schedule: 0 1
  expect:
    meanAccessDelayInMs: 1000
    successRate: 1
  target:
    host: http://${SERVER}
    method: GET
EOF
```

### Check Task Status

After completing a round of tasks, you can use the kdoctor aggregate API to view the report for the current round. When the FINISH field is set to true, it indicates that all tasks have been completed, and you can access the overall report.

```shell
kubectl get apphttphealthy
NAME        FINISH   EXPECTEDROUND   DONEROUND   LASTROUNDSTATUS   SCHEDULE
http        true     1               1           succeed           0 1
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
    NAME        CREATED AT
    http        0001-01-01T00:00:00Z
    ```

2. View specific task reports

    The reports are aggregated from the agents running on both the kdoctor-control-plane node and the kdoctor-worker nodes after performing two rounds of stress testing respectively.

    ```shell
    kubectl get kdoctorreport http -oyaml
    apiVersion: system.kdoctor.io/v1beta1
    kind: KdoctorReport
    metadata:
      creationTimestamp: null
      name: http
    spec:
      FailedRoundNumber: null
      FinishedRoundNumber: 1
      Report:
      - NodeName: kdoctor-control-plane
        HttpAppHealthyTask:
          Detail:
          - MeanDelay: 10.317307
            Metrics:
              Duration: 11.022081662s
              EndTime: "2023-07-31T07:25:23Z"
              Errors: {}
              Latencies:
                Max_inMx: 0
                Mean_inMs: 10.317307
                Min_inMs: 0
                P50_inMs: 0
                P90_inMs: 0
                P95_inMs: 0
                P99_inMs: 0
              RequestCounts: 104
              StartTime: "2023-07-31T07:25:12Z"
              StatusCodes:
                "200": 104
              SuccessCounts: 104
              TPS: 9.435604197939574
              TotalDataSize: 40040 byte
            Succeed: true
            SucceedRate: 1
            TargetMethod: GET
            TargetName: HttpAppHealthy target
            TargetUrl: http://172.41.71.0
          Succeed: true
          TargetNumber: 1
          TargetType: HttpAppHealthy
          MaxCPU: 30.651%
          MaxMemory: 97.00MB
        HttpAppHealthyTaskSpec:
        ...
        PodName: kdoctor-agent-fmr9m
        ReportType: agent test report
        RoundDuration: 11.038965547s
        RoundNumber: 1
        RoundResult: succeed
        StartTimeStamp: "2023-07-31T07:25:12Z"
        EndTimeStamp: "2023-07-31T07:25:23Z"
        TaskName: apphttphealthy.http
        TaskType: AppHttpHealthy
      - NodeName: kdoctor-worker
        HttpAppHealthyTask:
          Detail:
          - MeanDelay: 10.548077
            Metrics:
              ...
            Succeed: true
            SucceedRate: 1
            TargetMethod: GET
            TargetName: HttpAppHealthy target
            TargetUrl: http://172.41.71.0
          Succeed: true
          TargetNumber: 1
          TargetType: HttpAppHealthy
        HttpAppHealthyTaskSpec:
        ...
        PodName: kdoctor-agent-s468h
        ReportType: agent test report
        RoundDuration: 11.034140236s
        RoundNumber: 1
        RoundResult: succeed
        StartTimeStamp: "2023-07-31T07:25:12Z"
        EndTimeStamp: "2023-07-31T07:25:23Z"
        TaskName: apphttphealthy.http
        TaskType: AppHttpHealthy
      ReportRoundNumber: 1
      RoundNumber: 1
      Status: Finished
      TaskName: http
      TaskType: AppHttpHealthy
    ```

> If the reports do not align with the expected results, check the MaxCPU and MaxMemory fields in the report to verify if there are available resources of the agents and adjust the resource limits for the agents accordingly.

## Other Common Examples 

Below are examples of HTTP requests with bodies and HTTPS requests:

1. Create an `AppHttpHealthy` task for HTTP with a body. This task will run continuously for 10 seconds. It will send POST requests with the provided body to the specified server at a rate of 10 QPS and be executed immediately.

    We are using the service address of the test server. If you have a different server address available, feel free to use it instead.

    Creating test body data

    ```shell
    cat <<EOF | kubectl apply -f -
    apiVersion: v1
    kind: ConfigMap
    metadata:
      name: kdoctor-test-body
      namespace: kdoctor-test-server
    data:
      test1: test1
      test2: test2
    EOF
    kubectl apply -f http-body.yaml
    ```

    Create `AppHttpHealthy`

    ```shell
    SERVER="172.41.71.0"
    cat <<EOF | kubectl apply -f -
    apiVersion: kdoctor.io/v1beta1
    kind: AppHttpHealthy
    metadata:
      name: http-body
    spec:
      request:
        durationInSecond: 10
        perRequestTimeoutInMS: 1000
        qps: 10
      schedule:
        roundNumber: 1
        roundTimeoutMinute: 1
        schedule: 0 1
      expect:
        meanAccessDelayInMs: 1000
        successRate: 1
        statusCode: 200
      target:
        bodyConfigmapName: kdoctor-test-body
        bodyConfigmapNamespace: kdoctor-test-server
        header:
         - "Content-Type: application/json"
        host: http://${SERVER}
        method: POST
    EOF
    ```

2. Create an `AppHttpHealthy` task for HTTPS. This task will run continuously for 10 seconds. It will send GET requests using the HTTPS protocol with the provided certificate to the specified server at a rate of 10 QPS and be executed immediately.

    The TLS certificate used in this example is generated by the server and is only valid for the Pod's IP. Hence, we are accessing the server using the Pod's IP. If you are using a different server, please create the certificate secret accordingly.

    ```shell
    SERVER="172.40.0.35"
    cat <<EOF | kubectl apply -f -
    apiVersion: kdoctor.io/v1beta1
    kind: AppHttpHealthy
    metadata:
      name: https
    spec:
      request:
        durationInSecond: 10
        perRequestTimeoutInMS: 1000
        qps: 10
      schedule:
        roundNumber: 1
        roundTimeoutMinute: 1
        schedule: 0 1
      expect:
        meanAccessDelayInMs: 1000
        successRate: 1
        statusCode: 200
      target:
        host: https://${SERVER}
        method: GET
        tlsSecretName: https-client-cert
        tlsSecretNamespace: kdoctor-test-server
    EOF
    ```

## Environment Cleanup

```shell
kubectl delete apphttphealthy http https http-body
```
