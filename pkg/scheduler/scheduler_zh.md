# 任务资源调度器

在Kdoctor旧版本中，每个任务的运行载体是Kdoctor-Agent这个DaemonSet。当部署多个任务的时候，该DaemonSet的负载较高且会影响到各个任务执行的效果与准确性。
因此，为了实现执行每个任务的载体之间的资源隔离，我们设计出一个任务可拥有一个对应的Kubernetes控制器载体，如DaemonSet,Deployment...(默认为DaemonSet)。

## 载体

### 资源定义

我们可在 `Netdns`, `NetReach`, `AppHttpHealthy` 等任务实例下发之前，为其Spec中补充`载体`的自定义。

```yaml
type AgentSpec struct {
	// +kubebuilder:validation:Optional
	Annotation map[string]string `json:"annotation,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=DaemonSet
	// +kubebuilder:validation:Enum=Deployment;DaemonSet
	Kind *string `json:"kind,omitempty"`

	// +kubebuilder:validation:Optional
	DeploymentReplicas *int32 `json:"deploymentReplicas,omitempty"`

	// +kubebuilder:validation:Optional
	Affinity *v1.Affinity `json:"affinity,omitempty"`

	// +kubebuilder:validation:Optional
	Env []v1.EnvVar `json:"env,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=false
	HostNetwork bool `json:"hostNetwork,omitempty"`

	// +kubebuilder:validation:Optional
	Resources *v1.ResourceRequirements `json:"resources,omitempty"`

	// +kubebuilder:default=60
	// +kubebuilder:validation:Optional
	TerminationGracePeriodMinutes *int64 `json:"terminationGracePeriodMinutes,omitempty"`
}
```

### 载体调度

#### 载体控制器

当任务下发后，我们会为该任务资源实例的 `status` 字段补充一个 `resource` 的属性，用来展示该任务对应载体的运行情况。
`runtimeStatus`表明该载体资源的当前情况，如 `creating` 表示该载体正在创建中，`created` 表明该载体已创建完毕，对应的任务可开始执行, `deleted` 表明该载体已删除。

```yaml
type TaskResource struct {
  // +kubebuilder:validation:Required
  RuntimeName string `json:"runtimeName,omitempty"`

  // +kubebuilder:validation:Required
  RuntimeType string `json:"runtimeType,omitempty"`

  // +kubebuilder:validation:Optional
  ServiceNameV4 *string `json:"serviceNameV4,omitempty"`

  // +kubebuilder:validation:Optional
  ServiceNameV6 *string `json:"serviceNameV6,omitempty"`

  // +kubebuilder:validation:Required
  // +kubebuilder:validation:Enum=creating;created;deleted
  RuntimeStatus string `json:"runtimeStatus,omitempty"`
}
```

#### 载体Service

我们会为每个载体创建对应的Service资源，且使用OwnerReference将该Service资源与载体控制器资源绑定在一起。因此载体控制器资源被删除后，其Service资源会被Kubernetes GC掉。

### 代码逻辑

#### scheduler

我们在controllerReconcile中根据某个任务的`status.resource`属性状态，选择调用一个scheduler模块来为其创建载体控制器与否，并将其记录在一个cacheDB的内存数据库中。

#### tracker

我们会启动tracker协程异步的监控cacheDB内存数据库中记录的 `载体与任务` 的数据，并选择更新该任务的`status.resource`属性数据。另外，该tracker也会追踪选择在何时清理掉该载体控制器资源，以减轻空闲残余资源。

### Notice

1. 载体会有一个OwnerReference与任务进行绑定。
2. 载体Service会有一个OwnerReference与载体控制器资源绑定。
