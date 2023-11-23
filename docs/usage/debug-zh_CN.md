# Debug

[**English**](./debug.md) | **简体中文**

**Q: 想要使用更高的 QPS 应该如何设置？**
* A: 当 QPS 设置过大，会导致服务器资源占用过高，影响业务。为了防止在生产坏境出现误操作。kdoctor 在 webhook 中添加了 QPS 的检查。如果您想使用更高的 QPS
可通过参数设置 QPS 限制 `--set feature.nethttp_defaultRequest_MaxQps=1000`，也可以通过 kdoctor 的 configmap 中去更改 `nethttp_defaultRequest_MaxQps` ，
并重启 kdoctor 的相关 pod 重新加载 configmap。

**Q: 为什么我的任务无法达到期望的 QPS ？**
* A：无法达到 QPS 的期望原因有很多主要分为以下几种原因：
  * 并发 worker 设置过低，kdoctor 可通过设置参数调整并发数 `--set feature.nethttp_defaultConcurrency=50`，`--set feature.netdns_defaultConcurrency=50`。
  * kdoctor agent 分配资源不充足，可通过 kdoctor 的聚合报告`kubectl get kdoctorreport `查看任务消耗的 cpu 与 内存使用量，确定 kdoctor agent 资源分配是否充足。
     ```shell
      ~kubectl get kdoctorreport test-task -oyaml
      ...
      SystemResource:
        MaxCPU: 52.951%
        MaxMemory: 120.00MB
        MeanCPU: 32.645%
      ...
    ```
  * kdoctor agent 中是否在同时执行其他任务将资源占满。可通过 kdoctor 的聚合报告`kubectl get kdoctorreport 查看同时执行的其他任务 QPS 数量。
    错开任务执行时间或通过定义 agentSpec 指定 kdoctor agent 执行任务将任务进行隔离。因 QPS 统计具有时效性，所以可搭配日志一起作为参考，在任务执行开始前，会将当前在执行的 QPS 输出到日志中。
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
**Q: 为什么 Kdoctor agent 会 OOM？**
* A: kdoctor agent 作为默认的执行任务的 agent，在任务中没有指定 agent 时，默认使用 kdoctor agent 执行，目前的 agent 还不支持根据任务负载情况，拒绝执行任务或延迟执行任务功能。
     因此当 kdoctor agent 同时执行大量任务时，由于请求量过大，内存限制过低，将会导致 kdoctor agent 内存过载，导致 OOM，我们可以根据任务情况，错开使用 kdoctor agent，调整内存限制，或使用指定的 agent 隔离任务。
