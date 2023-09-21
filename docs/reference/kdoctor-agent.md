# kdoctor-agent

This page describes CLI options and ENV of kdoctor-agent.

## kdoctor-agent daemon

Run the kdoctor agent daemon.

### Options

| options             | type   | default                  | description                                                                                             |
|---------------------|--------|--------------------------|---------------------------------------------------------------------------------------------------------|
| --config-dir        | string | /tmp/config-map/conf.yml | config file path.                                                                                       |
| --app-mode          | bool   | false                    | agent running mode ,when using app mode, the agent only provides an HTTP and HTTPS server for testing.  |
| --tls-insecure      | bool   | true                     | The HTTPS server skips TLS authentication.                                                              |
| --tls-ca-cert       | string | /etc/tls/ca.crt          | The CA certificate path, which is used by the agent to generate the signing certificate.                |
| --tls-ca-key        | string | /etc/tls/ca.key          | The CA key path, which is used by the agent to generate the signing certificate.                        |
| --task-kind         | string | ""                       | The kind of task. values AppHttpHealthy、NetReach and Netdns.                                            |
| --task-name         | string | ""                       | The name of task.                                                                                       |
| --service-ipv4-name | string | ""                       | The ipv4 service name of the task workload.                                                             |
| --service-ipv6-name | string | ""                       | The ipv6 service name of the task workload.                                                             |

### ENV

| env                                            | default       | description                                                                        |
|------------------------------------------------|---------------|------------------------------------------------------------------------------------|
| ENV_LOG_LEVEL                                  | info          | Log level, optional values are "debug", "info", "warn", "error", "fatal", "panic". |
| ENV_ENABLED_METRIC                             | false         | Enable/disable metrics.                                                            |
| ENV_METRIC_HTTP_PORT                           | 5711          | Metric HTTP server port.                                                           |
| ENV_AGENT_HEALTH_HTTP_PORT                     | 5710          | kdoctor-agent health backend HTTP server port.                                     |
| ENV_AGENT_APP_HTTP_PORT                        | 80            | kdoctor-agent app backend HTTP server port.                                        |
| ENV_AGENT_APP_HTTPS_PORT                       | 443           | kdoctor-agent app backend HTTP server port.                                        |
| ENV_ENABLE_AGGREGATE_AGENT_REPORT              | false         | enable aggregate report                                                            |
| ENV_CLEAN_AGED_REPORT_INTERVAL_IN_MINUTE       | 10            | clean aggregate report interval in minute                                          |
| ENV_AGENT_REPORT_STORAGE_PATH                  | /report       | aggregate report storage path                                                      |
| ENV_GOPS_LISTEN_PORT                           | 5712          | Gops port                                                                          |
| ENV_PYROSCOPE_PUSH_SERVER_ADDRESS              | ""            | pyroscope addr                                                                     |
| ENV_POD_NAME                                   | ""            | agent pod name                                                                     |
| ENV_POD_NAMESPACE                              | ""            | agent pod namespace                                                                |
| ENV_GOLANG_MAXPROCS                            | 8             | golang runtime max procs                                                           |
| ENV_AGENT_GRPC_LISTEN_PORT                     | 3000          | agent grpc port                                                                    |
| ENV_CLUSTER_DNS_DOMAIN                         | cluster.local | cluster domian                                                                     |
| ENV_LOCAL_NODE_IP                              | ""            | loacl node ip                                                                      |
| ENV_LOCAL_NODE_NAME                            | ""            | loacl node name                                                                    |
| ENV_AGENT_RESOURCE_COLLECT_INTERVAL_IN_SECOND  | "1"           | agent CPU and memory usage collection interval time                                |



