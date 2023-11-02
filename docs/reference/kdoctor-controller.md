# kdoctor-controller

This page describes CLI options and ENV of kdoctor-controller.

## kdoctor-controller daemon

Run the kdoctor controller daemon.

### Options

| options                          | type   | default                                    | description                                                          |
|----------------------------------|--------|--------------------------------------------|----------------------------------------------------------------------|
| --config-dir                     | string | /tmp/config-map/conf.yml                   | config file path.                                                    |
| --tls-ca-cert                    | string | /etc/tls/ca.crt                            | The CA certificate path, The CA is used to validate the certificate. |
| --tls-server-cert                | string | /etc/tls/tls.crt                           | The server tls cert path.                                            |
| --tls-server-key                 | string | /etc/tls/tls.key                           | The server tls key path.                                             |
| --configmap-deployment-template  | string | /tmp/configmap-app-template/deployment.yml | The configmap deployment template file path.                         |
| --configmap-daemonset-template   | string | /tmp/configmap-app-template/daemonset.yml  | The configmap daemonset template file path.                          |
| --configmap-pod-template         | string | /tmp/configmap-app-template/pod.yml        | The configmap pod template file path.                                |
| --configmap-service-template     | string | /tmp/configmap-app-template/service.yml    | The configmap service template file path.                            |
| --configmap-ingress-template     | string | /tmp/configmap-app-template/ingress.yml    | The configmap ingress template file path.                            |

### ENV

| env                                         | default       | description                                                                        |
|---------------------------------------------|---------------|------------------------------------------------------------------------------------|
| ENV_LOG_LEVEL                               | info          | Log level, optional values are "debug", "info", "warn", "error", "fatal", "panic". |
| ENV_ENABLED_METRIC                          | false         | Enable/disable metrics.                                                            |
| ENV_METRIC_HTTP_PORT                        | 5711          | Metric HTTP server port.                                                           |
| ENV_HTTP_PORT                               | 80            | kdoctor-controller backend HTTP server port.                                       |
| ENV_ENABLE_AGGREGATE_AGENT_REPORT           | false         | enable aggregate report                                                            |
| ENV_CLEAN_AGED_REPORT_INTERVAL_IN_MINUTE    | 10            | clean aggregate report interval in minute                                          |
| ENV_COLLECT_AGENT_REPORT_INTERVAL_IN_SECOND | 600           | collect agent report interval time                                                 |
| ENV_CONTROLLER_REPORT_AGE_IN_DAY            | 30            | controller report age in ady                                                       |
| ENV_AGENT_REPORT_STORAGE_PATH               | /report       | aggregate report storage path                                                      |
| ENV_CONTROLLER_REPORT_STORAGE_PATH          | /report       | controller report storage path                                                     |
| ENV_GOPS_LISTEN_PORT                        | 5724          | Gops port                                                                          |
| ENV_WEBHOOK_PORT                            | 5722          | controller webhook port                                                            |
| ENV_PYROSCOPE_PUSH_SERVER_ADDRESS           | ""            | pyroscope addr                                                                     |
| ENV_POD_NAME                                | ""            | controller pod name                                                                |
| ENV_POD_NAMESPACE                           | ""            | controller pod namespace                                                           |
| ENV_GOLANG_MAXPROCS                         | 8             | golang runtime max procs                                                           |
| ENV_DEFAULT_AGENT_NAME                      | kdoctor-agent | default agent name                                                                 |
| ENV_DEFAULT_AGENT_TYPE                      | daemonset     | default agent type                                                                 |
| ENV_DEFAULT_AGENT_SERVICE_V4_NAME           | ""            | default agent server ipv4 name                                                     |
| ENV_DEFAULT_AGENT_SERVICE_V6_NAME           | ""            | default agent server ipv6 name                                                     |



