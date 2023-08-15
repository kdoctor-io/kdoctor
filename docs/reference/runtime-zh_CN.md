## runtime

当下发任务 CR 后，kdoctor-controller 会根据 CR 中的 AgentSpec 生成对应的任务载体（DaemonSet 或 Deployment）执行任务，每一个任务独立使用一个载体。

### 载体资源

当任务 CR 下发后，kdocotr-controller 会创建如下资源进行任务。

### 工作负载

工作负载为 DaemonSet 或 Deployment，默认为 Daemonset，负载中的每一个 Pod 根据任务配置进行的请求，并将执行结果落盘到 Pod 中，可通过 AgentSpec 中设置
工作负载的销毁时间，默认任务执行完 60 分钟后，销毁工作负载，当删除 CR 任务时，工作负载会一并被删除。

### Service

在创建工作负载时，kdoctor-controller 同时会根据 IP Family 的配置，创建对应的 service 并于工作负载的 pod 绑定。用于测试 service 网络连通性。与工作负载
的销毁逻辑相同。

### Ingress

当任务为 NetReach 时，若测试目标包含 Ingress 时，会创建一个 Ingress，用于测试 Ingress 的网络联通性，与工作负载的销毁逻辑相同。

### 报告收取

当任务 CR 下发后，kdoctor-controller 会将任务注册进 ReportManager，ReportManager 会定期去每一个任务负载中通过 GRPC 接口获取报告，并聚合
在 kdoctor-controller 中，聚合后可通过命令 `kubectl get kdoctorreport` 获取报告结果，因此，若报告未收集完成就将工作负载删除将影响报告聚合结果。
