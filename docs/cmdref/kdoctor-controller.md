# kdoctor-controller

This page describes CLI options and ENV of kdoctor-controller.

## kdoctor-controller daemon

Run the kdoctor controller daemon.

### Options

```
    --config-dir        string    config file path (default /tmp/config-map/conf.yml) 
    --tls-ca-cert       string    The CA certificate path, The CA is used to validate the certificate (default /etc/tls/ca.crt) 
    --tls-server-cert   string    The server tls cert path  (default /etc/tls/tls.crt) 
    --tls-server-key    string    The server tls key path   (default /etc/tls/tls.key) 
```

### ENV

```
    ENV_LOG_LEVEL                                  log level (DEBUG|INFO|ERROR)
    ENV_ENABLED_METRIC                             enable metrics (true|false)
    ENV_METRIC_HTTP_PORT                           metric port (default to 5711)
    ENV_HTTP_PORT                                  http port  (default to 80)
    ENV_GOPS_LISTEN_PORT                           Gops port 
    ENV_WEBHOOK_PORT                               controller webhook port 
    ENV_ENABLE_AGGREGATE_AGENT_REPORT              enable aggregate report  (default to false)
    ENV_CLEAN_AGED_REPORT_INTERVAL_IN_MINUTE       clean aggregate report interval in minute (default to 10)
    ENV_AGENT_REPORT_STORAGE_PATH                  aggregate agent report storage path (default to /report)
    ENV_CONTROLLER_REPORT_STORAGE_PATH             controller report storage path (default to /report)
    ENV_COLLECT_AGENT_REPORT_INTERVAL_IN_SECOND    collect agent report interval time (default to 600)
    ENV_PYROSCOPE_PUSH_SERVER_ADDRESS              pyroscope addr   
    ENV_POD_NAME                                   pod name
    ENV_POD_NAMESPACE                              pod namespace
    ENV_GOLANG_MAXPROCS                            golang runtime max procs
    ENV_AGENT_GRPC_LISTEN_PORT                     agent grpc port (default to 3000)
    ENV_CONTROLLER_REPORT_AGE_IN_DAY               controller report age in ady (default to 30)
```
