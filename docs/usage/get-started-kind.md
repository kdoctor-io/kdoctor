# Kind Quick Start

[**简体中文**](./get-started-kind-zh_CN.md) | **English**

Kind is a tool running local Kubernetes clusters using Docker container nodes. kdoctor provides scripts for installing Kind clusters to quickly build a Kind cluster with nginx, ingress, and loadbalancer, which you can use for kdoctor testing and experience.

## Prerequisites

* Execute ``make checkBin`` to check if the development tools on the localhost meet the requirements for deploying a Kind cluster with kdoctor, and if the components are missing, they will be installed.

## Deploy kdoctor on a Kind Cluster

1. Clone the code for the stable version of kdoctor to the localhost and go to the root directory of the kdoctor project.

    ```bash
    ~# LATEST_RELEASE_VERISON=$(curl -s https://api.github.com/repos/kdoctor-io/kdoctor/releases | grep '"tag_name":' | grep -v rc | grep -Eo "([0-9]+\.[0-9]+\.[0-9])" | sort -r | head -n 1)
    ~# curl -Lo /tmp/$LATEST_RELEASE_VERISON.tar.gz https://github.com/kdoctor-io/kdoctor/archive/refs/tags/v$LATEST_RELEASE_VERISON.tar.gz
    ~# tar -xvf /tmp/$LATEST_RELEASE_VERISON.tar.gz -C /tmp/
    ~# cd /tmp/kdoctor-$LATEST_RELEASE_VERISON
    ```

2. Get the latest image of kdoctor by using the following method.

    ```bash
    ~# KDOCTOR_LATEST_IMAGE_TAG=$(curl -s https://api.github.com/repos/kdoctor-io/kdoctor/releases | jq -r '. [].tag_name | select(("^v1.[0-9]*. [0-9]*$"))' | head -n 1)
    ```

3. Execute the following command to create a Kind cluster and install metallb, contour, nginx, and kdoctor.

    ```bash
    ~# make e2e_init -e PROJECT_IMAGE_VERSION=KDOCTOR_LATEST_IMAGE_TAG
    ```

    note: If you are a domestic user, you can use the following command to avoid pulling the image failures.

    ```bash
    ~# make e2e_init -e E2E_SPIDERPOOL_TAG=$SPIDERPOOL_LATEST_IMAGE_TAG -e E2E_CHINA_IMAGE_REGISTRY=true
    ```

## Verify the Installation

Configure KUBECONFIG of the Kind cluster for kubectl by executing the following command in the root directory of the kdoctor project.

   ```bash
   ~# export KUBECONFIG=$(pwd)/test/runtime/kubeconfig_kdoctor.config
   ```

You can see the following output:

```bash
~# kubectl get nodes 
NAME                    STATUS   ROLES           AGE     VERSION
kdoctor-control-plane   Ready    control-plane   3h50m    v1.27.1
kdoctor-worker          Ready    <none>          3h50m   v1.27.1

~# kubectl get pod -n kdoctor -owide
NAME                                     READY   STATUS    RESTARTS   AGE     IP            NODE                    NOMINATED NODE   READINESS GATES
kdoctor-agent-5n4nb                      1/1     Running   0          3h46m   172.40.1.29   kdoctor-worker          <none>           <none>
kdoctor-agent-zm4tn                      1/1     Running   0          3h46m   172.40.0.83   kdoctor-control-plane   <none>           <none>
kdoctor-controller-78589d96c8-lgcw9      1/1     Running   0          3h46m   172.40.1.28   kdoctor-worker          <none>           <none>
kdoctor-test-server-6bf7f9df47-dq8th     1/1     Running   0          3h46m   172.40.0.82   kdoctor-control-plane   <none>           <none>
kdoctor-test-server-6bf7f9df47-mhsml     1/1     Running   0          3h46m   172.40.1.27   kdoctor-worker          <none>           <none>

~# kubectl get svc -n kdoctor
NAME                    TYPE           CLUSTER-IP       EXTERNAL-IP              PORT(S)                                     AGE
kdoctor-agent-ipv4      LoadBalancer   172.41.217.12    172.18.0.51              5711:30778/TCP,80:31835/TCP,443:30675/TCP   3h46m
kdoctor-agent-ipv6      LoadBalancer   fd41::2274       fc00:f853:ccd:e793::50   5711:30022/TCP,80:30761/TCP,443:30516/TCP   3h46m
kdoctor-controller      ClusterIP      172.41.210.120   <none>                   5721/TCP,5722/TCP,443/TCP                   3h46m
kdoctor-test-server     ClusterIP      172.41.95.144    <none>                   80/TCP,443/TCP,53/UDP,53/TCP,853/TCP        3h46m
```

> `kdoctor-test-server` is kdoctor's test server, which contains http server, https server, dns udp server, dns tcp server, for testing kdocotr functions.

## 配置任务

===  "AppHttpHealthy"

    We have access to kdocotr-test-server's service ip to get a response from kdocotr-test-server.
    
    ```bash
    SERVER="172.41.95.144"
    cat <<EOF | kubectl apply -f -
    apiVersion: kdoctor.io/v1beta1
    kind: AppHttpHealthy
    metadata:
      name: http-test
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

    Check the status of the task and wait for it to complete.

    ```bash
    ～# kubectl get apphttphealthy
    NAME        FINISH   EXPECTEDROUND   DONEROUND   LASTROUNDSTATUS   SCHEDULE
    http-test   false    1               0                             0 1
    ～# kubectl get apphttphealthy
    NAME        FINISH   EXPECTEDROUND   DONEROUND   LASTROUNDSTATUS   SCHEDULE
    http-test   true     1               1           succeed           0 1
    ```      
      
    Detailed reports on query missions.

    ```bash
    ~# kubectl get kdoctorreport http-test -oyaml
    apiVersion: system.kdoctor.io/v1beta1
    kind: KdoctorReport
    metadata:
      creationTimestamp: "2023-12-11T07:20:01Z"
      name: http-test
    spec:
      FailedRoundNumber: null
      FinishedRoundNumber: 1
      Report:
      - EndTimeStamp: "2023-12-11T07:20:11Z"
        HttpAppHealthyTask:
          Detail:
          - MeanDelay: 10.44
            Metrics:
              Duration: 10.003604664s
              EndTime: "2023-12-11T07:20:11Z"
              Errors: {}
              ExistsNotSendRequests: false
              Latencies:
                MaxInMs: 0
                MeanInMs: 10.44
                MinInMs: 0
                P50InMs: 0
                P90InMs: 0
                P95InMs: 0
                P99InMs: 0
              RequestCounts: 100
              StartTime: "2023-12-11T07:20:01Z"
              StatusCodes:
                "200": 100
              SuccessCounts: 100
              TPS: 9.996396634892049
              TotalDataSize: 37394 byte
            Succeed: true
            SucceedRate: 1
            TargetMethod: GET
            TargetName: HttpAppHealthy target
            TargetUrl: http://172.41.95.144
          Succeed: true
          SystemResource:
            MaxCPU: 12.723%
            MaxMemory: 35.00MB
            MeanCPU: 6.227%
          TargetNumber: 1
          TargetType: HttpAppHealthy
          TotalRunningLoad:
            AppHttpHealthyQPS: 10
            NetDnsQPS: 0
            NetReachQPS: 0
        HttpAppHealthyTaskSpec:
          ...
        NodeName: kdoctor-control-plane
        PodName: kdoctor-agent-zm4tn
        ReportType: agent test report
        RoundDuration: 10.014725655s
        RoundNumber: 1
        RoundResult: succeed
        StartTimeStamp: "2023-12-11T07:20:01Z"
        TaskName: apphttphealthy.http-test
        TaskType: AppHttpHealthy
      - EndTimeStamp: "2023-12-11T07:20:11Z"
        HttpAppHealthyTask:
          Detail:
          - MeanDelay: 11.24
            Metrics:
              Duration: 10.00058331s
              EndTime: "2023-12-11T07:20:11Z"
              Errors: {}
              ExistsNotSendRequests: false
              Latencies:
                MaxInMs: 0
                MeanInMs: 11.24
                MinInMs: 0
                P50InMs: 0
                P90InMs: 0
                P95InMs: 0
                P99InMs: 0
              RequestCounts: 100
              StartTime: "2023-12-11T07:20:01Z"
              StatusCodes:
                "200": 100
              SuccessCounts: 100
              TPS: 9.999416724023071
              TotalDataSize: 37391 byte
            Succeed: true
            SucceedRate: 1
            TargetMethod: GET
            TargetName: HttpAppHealthy target
            TargetUrl: http://172.41.95.144
          Succeed: true
          SystemResource:
            MaxCPU: 12.704%
            MaxMemory: 35.00MB
            MeanCPU: 6.370%
          TargetNumber: 1
          TargetType: HttpAppHealthy
          TotalRunningLoad:
            AppHttpHealthyQPS: 10
            NetDnsQPS: 0
            NetReachQPS: 0
        HttpAppHealthyTaskSpec:
          ...
        NodeName: kdoctor-worker
        PodName: kdoctor-agent-5n4nb
        ReportType: agent test report
        RoundDuration: 10.010301747s
        RoundNumber: 1
        RoundResult: succeed
        StartTimeStamp: "2023-12-11T07:20:01Z"
        TaskName: apphttphealthy.http-test
        TaskType: AppHttpHealthy
      ReportRoundNumber: 1
      RoundNumber: 1
      Status: Finished
      TaskName: http-test
      TaskType: AppHttpHealthy
    ```

===  "NetReach"

    We test the connectivity of the cluster.

    ```bash
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

    Check the status of the task and wait for it to complete.

    ```bash 
    ~# kubectl get netreach
    NAME   FINISH   EXPECTEDROUND   DONEROUND   LASTROUNDSTATUS   SCHEDULE
    task   false    1               0                             0 1
    ~# kubectl get netreach
    NAME   FINISH   EXPECTEDROUND   DONEROUND   LASTROUNDSTATUS   SCHEDULE
    task   true     1               1           succeed           0 1
    ```

    Detailed reports on query missions.

    ```bash
    root@kdoctor-control-plane:/# kubectl get kdoctorreport task -oyaml
    apiVersion: system.kdoctor.io/v1beta1
    kind: KdoctorReport
    metadata:
    creationTimestamp: null
    name: task
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

===  "NetDns"

    We perform a connectivity check on the cluster's dns service.

    ```bash
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

    Check the status of the task and wait for it to complete.

    ```bash
    ~# kubectl get netdns
    NAME             FINISH   EXPECTEDROUND   DONEROUND   LASTROUNDSTATUS   SCHEDULE
    netdns-cluster   false    1               0                             0 1
    ~# kubectl get netdns
    NAME             FINISH   EXPECTEDROUND   DONEROUND   LASTROUNDSTATUS   SCHEDULE
    netdns-cluster   true     1               1           succeed           0 1
    ```

    Detailed reports on query missions.

    ```bash
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

Next you can customize the configuration of the task according to your needs:[AppHttpHealthy](../reference/apphttphealthy.md)、[NetReach](../reference/netreach.md)、[NetDns](../reference/netdns.md)

## Uninstall

* To uninstall the Kind cluster

    Execute `make e2e_clean` to uninstall the Kind cluster.
