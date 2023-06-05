# feature 

测试集群功能、性能巡检工具，批量化巡检，极大降低人工负载

## CRD

### NetReachHealthy and HttpAppHealthy

* 大规模集群部署后，巡检每个节点上的网络情况，具体是：本节点是否能够通过 pod ip（multus多网卡）、cluster ip、
  nodePort、loadbalancer ip、ingress ip 、多网卡 等所有网络渠道，访问到集群其它节点

* 集群所有的 node 上 去 压测一个集群内/外的应用地址，以查看应用的性能、集群每个角落到达该应用的连通性、给应用注入压力复现某类bug

* 给 api server 注入压力，以辅助排查 其他组件（依赖 api server）的高可用

* 生产和开发环境的心跳巡检，以qps=1为压力，每 1m 间隔巡检整个集群内 full mesh 网络的连通性、其它节点到一个应用的可用性

### HttpAppDetect

* 自动发现应用的最大性能

### NetDelayHealthy

* 集群节点间的延时

### DnsHealthy

* 大规模集群部署后，测试集群中每个角落访问 dns 的连通性

* 大规模集群部署后，调试 coredns 的副本数，确认是否满足设计需求

* 测试集群外部的 DNS 服务

### DnsDetect

* 自动发现 dns 的最大性能

### NetTcpHealthy

### NetUdpHealthy

### StorageLocalDisk

* 大规模集群部署后，测试集群中每个主机上的磁盘 吞吐量 和 延时

### CpuPressure ?

* 给每个主机上注入 CPU 压力，以测试应用的稳定性，复现一些 bug

### MemoryPressure ?

* 给每个主机上注入 memory 压力，以测试应用的稳定性，复现一些 bug

### RegistryHealthy

* 检测每个节点到镜像仓库的连通性

### K8sApiHealthy

### MysqlHealthy

## report

支持通过 API 获取报告

支持 pvc、本地磁盘存储

日志吐出

## metric

## 其它

如果有 job 时间重叠了，则只允许运行一个 或者 多个，避免自身 CPU 不足影响 job 的结果

中间件、etcd 等 探测


| kind          | feature                                                               | status |
|---------------|-----------------------------------------------------------------------|--------|
| 任务            | 支持周期和一次性调度                                                            |        |
|               | 多个任务并发时，支持规避，避免 cpu 和 memory 耗尽，使得任务执行准确                              |        |
|               | 所有任务，最大 qps 安全边际设置                                                    |        |
|               | 支持设置发压qps、发压时间                                                        |        |
|               | 支持设置并发压测 worker ， 且支持跟随应用副本数量自动合适并发数，以满足K8S 负载均衡特性                    |        |
|               | 支持select 发压 pod                                                       |        |
| 网络可达性         | 支持 pod ip，cluster ip，nodePort，loadbalancer ip，ingress，多网卡，ipv6 等多样化渠道 |        |
|               | 持续发压，发现偶发丢包                                                           |        |
| 服务 http 巡检    | 支持 pod selector 和 url                                                 |        |
|               | 持续发压，发现偶发丢包                                                           |        |
|               | 支持 http/https/http2                                                   |        |
|               | 支持定制 header、method、body                                               |        |
| dns巡检         | 网络可达性                                                                 |        |
|               | 性能测试                                                                  |        |
| tcp 巡检        | 网络吞吐量                                                                 |        |
| udp 巡检        | 网络丢包率                                                                 |        |
| api server 巡检 | 网络可达性                                                                 |        |
| 存储巡检          | 支持本地磁盘巡检                                                              |        |
|               | IO 吞吐量和延时                                                             |        |
| 镜像仓库巡检        | 网络可达性                                                                 |        |
| 中间件巡检         | 网络可达性                                                                 |        |
|               | mysql  redis                                                          |        |
| 报告            | CR 中状态展示                                                              |        |
|               | 详细报告支持 pv 存储                                                          |        |
|               | 详细报告支持 API 获取                                                         |        |
|               | 详细报告支持 webhook 吐出                                                     |        |
|               | 指标                                                                    |        |
|               | 详细报告的保留时间设置                                                           |        |
|               | 报告轮滚，避免存满 PVC                                                         |        |


