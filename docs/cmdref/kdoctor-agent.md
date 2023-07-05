# kdoctor-agent

This page describes CLI options and ENV of kdoctor-agent.

## kdoctor-agent daemon

Run the kdoctor agent daemon.

### Options

```
    --config-dir    string    config file path (default /tmp/config-map/conf.yml)
    --app-mode      bool      agent running mode ,when using app mode, the agent only provides an HTTP and HTTPS server for testing (default false) 
    --tls-insecure  bool      The HTTP server skips TLS authentication (default true) 
    --tls-ca-cert   string    The CA certificate path, which is used by the agent to generate the signing certificate  (default /etc/tls/ca.crt) 
    --tls-ca-key    string    The CA key path, which is used by the agent to generate the signing certificate  (default /etc/tls/ca.key) 
```

### ENV

```
    ENV_LOG_LEVEL                              log level (DEBUG|INFO|ERROR)
    ENV_ENABLED_METRIC                         enable metrics (true|false)
    ENV_METRIC_HTTP_PORT                       metric port (default to 5711)
    ENV_AGENT_HEALTH_HTTP_PORT                 http health port  (default to 5710)
    ENV_AGENT_APP_HTTP_PORT                    http app port  (default to 80)
    ENV_AGENT_APP_HTTPS_PORT                   https app port  (default to 443)
    ENV_ENABLE_AGGREGATE_AGENT_REPORT          enable aggregate report  (default to false)
    ENV_CLEAN_AGED_REPORT_INTERVAL_IN_MINUTE   clean aggregate report interval in minute (default to 10)
    ENV_AGENT_REPORT_STORAGE_PATH              aggregate report storage path (default to /report)
    ENV_GOPS_LISTEN_PORT                       Gops port 
    ENV_WEBHOOK_PORT                           webhook port 
    ENV_PYROSCOPE_PUSH_SERVER_ADDRESS          pyroscope addr       
    ENV_POD_NAME                               pod name
    ENV_POD_NAMESPACE                          pod namespace
    ENV_GOLANG_MAXPROCS                        golang runtime max procs
    ENV_AGENT_GRPC_LISTEN_PORT                 agent grpc port (default to 3000)
    ENV_CLUSTER_DNS_DOMAIN                     cluster domian
    ENV_LOCAL_NODE_IP                          loacl node ip
    ENV_LOCAL_NODE_NAME                        loacl node name
```


