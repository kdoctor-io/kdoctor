# Debug

[**简体中文**](./debug-zh_CN.md) | **English**


**Q: How to achieve higher QPS?**
* A: When the QPS setting is too high, it can result in excessive server resource utilization, impacting business operations. 
  To prevent accidental misconfiguration in production environments, kdoctor has added QPS checks in the webhook. 
* If you wish to use a higher QPS, you can set the QPS limit using the parameter `--set feature.nethttp_defaultRequest_MaxQps=1000`，You can also modify it through the configmap in kdoctor `nethttp_defaultRequest_MaxQps`,
  And restart the relevant pods of kdoctor to reload the configmap.

**Q: Why is my task unable to achieve the desired QPS ？**
* A：There are several reasons why the expected QPS cannot be achieved, primarily categorized into the following reasons：
    * The concurrency worker setting is too low. kdoctor can adjust the concurrency by setting the parameters `--set feature.nethttp_defaultConcurrency=50` and `--set feature.netdns_defaultConcurrency=50`.
    * The kdoctor agent may have insufficient resource allocation. You can use the kdoctor aggregate report `kubectl get kdoctorreport` to check the CPU and memory usage of the task. This will help you determine if the resource allocation for the kdoctor agent is sufficient.
       ```shell
        ~kubectl get kdoctorreport test-task -oyaml
        ...
        SystemResource:
          MaxCPU: 52.951%
          MaxMemory: 120.00MB
          MeanCPU: 32.645%
        ...
      ```
    * Whether the kdoctor agent is concurrently executing other tasks and occupying resources can be determined by checking the QPS count of other tasks being executed simultaneously. 
      You can use the kdoctor aggregate report `kubectl get kdoctorreport` to view the QPS count of other concurrently running tasks.
      Stagger the task execution time or isolate the task by defining agentSpec to specify the kdoctor agent to execute the task. Because QPS statistics are time-sensitive, they can be used together with the log as a reference. 
      Before the task execution starts, the currently executing QPS will be output to the log。
       ```shell
        ~kubectl logs kdoctor-agent-74rrp  -n kdoctor |grep "Before the current task starts"
        {"level":"DEBUG","ts":"2023-11-07T10:01:02.821Z","agent":"agent.agentController.AppHttpHealthyReconciler.AppHttpHealthy.test-task.round1","caller":"pluginManager/agentTools.go:90","msg":"Before the current task starts, the total QPS of the tasks being executed is AppHttpHealth=100,NetReach=0,NetDNS=0","AppHttpHealthy":"test-task"}
       ```
       ```shell
        ~kubectl get kdoctorreport test-task -oyaml
        ...
        TotalRunningLoad:
          AppHttpHealthyQPS: 100
          NetDnsQPS: 50
          NetReachQPS: 0
        ...
       ```
