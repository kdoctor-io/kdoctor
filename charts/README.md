# kdoctor

## Introduction

## Features

## Parameters

### Global parameters

| Name                           | Description                                | Value                         |
| ------------------------------ | ------------------------------------------ | ----------------------------- |
| `global.imageRegistryOverride` | Global Docker image registry               | `""`                          |
| `global.imageTagOverride`      | Global Docker image tag                    | `""`                          |
| `global.name`                  | instance name                              | `kdoctor`                     |
| `global.clusterDnsDomain`      | cluster dns domain                         | `cluster.local`               |
| `global.commonAnnotations`     | Annotations to add to all deployed objects | `{}`                          |
| `global.commonLabels`          | Labels to add to all deployed objects      | `{}`                          |
| `global.configName`            | the configmap name                         | `kdoctor`                     |
| `global.configAppTemplate`     | the configmap name of agent                | `kdoctor-app-config-template` |

### feature parameters

| Name                                                                    | Description                                                             | Value                                |
| ----------------------------------------------------------------------- | ----------------------------------------------------------------------- | ------------------------------------ |
| `feature.enableIPv4`                                                    | enable ipv4                                                             | `true`                               |
| `feature.enableIPv6`                                                    | enable ipv6                                                             | `true`                               |
| `feature.netReachRequestMaxQPS`                                         | qps for kind NetReach                                                   | `20`                                 |
| `feature.appHttpHealthyRequestMaxQPS`                                   | qps for kind AppHttpHealthy                                             | `100`                                |
| `feature.netHttpDefaultRequestQPS`                                      | qps for kind NetHttp                                                    | `10`                                 |
| `feature.netHttpDefaultRequestDurationInSecond`                         | Duration In Second for kind NetHttp                                     | `2`                                  |
| `feature.netHttpDefaultRequestPerRequestTimeoutInMS`                    | PerRequest Timeout In MS for kind NetHttp                               | `500`                                |
| `feature.netDnsRequestMaxQPS`                                           | qps for kind NetDns                                                     | `100`                                |
| `feature.agentDefaultTerminationGracePeriodMinutes`                     | agent termination after minutes                                         | `60`                                 |
| `feature.taskPollIntervalInSecond`                                      | the interval to poll the task in controller and agent pod               | `5`                                  |
| `feature.multusPodAnnotationKey`                                        | the multus annotation key for ip status                                 | `k8s.v1.cni.cncf.io/networks-status` |
| `feature.crdMaxHistory`                                                 | max history items inf CRD status                                        | `10`                                 |
| `feature.aggregateReport.enabled`                                       | aggregate report from agent for each crd                                | `true`                               |
| `feature.aggregateReport.cleanAgedReportIntervalInMinute`               | the interval in minute for removing aged report                         | `10`                                 |
| `feature.aggregateReport.agent.reportPath`                              | the path where the agent pod temporarily store task report.             | `/report`                            |
| `feature.aggregateReport.controller.reportHostPath`                     | storage path when pvc is disabled                                       | `/var/run/kdoctor/reports`           |
| `feature.aggregateReport.controller.maxAgeInDay`                        | report file maximum age in days                                         | `30`                                 |
| `feature.aggregateReport.controller.collectAgentReportIntervalInSecond` | how long the controller collects all agent report at interval in second | `600`                                |
| `feature.aggregateReport.controller.pvc.enabled`                        | store report to pvc                                                     | `false`                              |
| `feature.aggregateReport.controller.pvc.storageClass`                   | storage class name                                                      | `""`                                 |
| `feature.aggregateReport.controller.pvc.storageRequests`                | storage request                                                         | `100Mi`                              |
| `feature.aggregateReport.controller.pvc.storageLimits`                  | storage limit                                                           | `1024Mi`                             |

### kdoctorAgent parameters

| Name                                                           | Description                                                                                                                     | Value                           |
| -------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------- | ------------------------------- |
| `kdoctorAgent.uniqueMatchLabelKey`                             | the unique match label key for Agent                                                                                            | `""`                            |
| `kdoctorAgent.name`                                            | the kdoctorAgent name                                                                                                           | `kdoctor-agent`                 |
| `kdoctorAgent.cmdBinName`                                      | the binary name of kdoctorAgent                                                                                                 | `/usr/bin/agent`                |
| `kdoctorAgent.hostnetwork`                                     | enable hostnetwork mode of kdoctorAgent pod                                                                                     | `false`                         |
| `kdoctorAgent.image.registry`                                  | the image registry of kdoctorAgent                                                                                              | `ghcr.io`                       |
| `kdoctorAgent.image.repository`                                | the image repository of kdoctorAgent                                                                                            | `kdoctor-io/kdoctor-agent`      |
| `kdoctorAgent.image.pullPolicy`                                | the image pullPolicy of kdoctorAgent                                                                                            | `IfNotPresent`                  |
| `kdoctorAgent.image.digest`                                    | the image digest of kdoctorAgent, which takes preference over tag                                                               | `""`                            |
| `kdoctorAgent.image.tag`                                       | the image tag of kdoctorAgent, overrides the image tag whose default is the chart appVersion.                                   | `""`                            |
| `kdoctorAgent.image.imagePullSecrets`                          | the image imagePullSecrets of kdoctorAgent                                                                                      | `[]`                            |
| `kdoctorAgent.serviceAccount.create`                           | create the service account for the kdoctorAgent                                                                                 | `true`                          |
| `kdoctorAgent.serviceAccount.annotations`                      | the annotations of kdoctorAgent service account                                                                                 | `{}`                            |
| `kdoctorAgent.service.annotations`                             | the annotations for kdoctorAgent service                                                                                        | `{}`                            |
| `kdoctorAgent.service.type`                                    | the type for kdoctorAgent service                                                                                               | `LoadBalancer`                  |
| `kdoctorAgent.ingress.enable`                                  | install ingress                                                                                                                 | `true`                          |
| `kdoctorAgent.ingress.ingressClass`                            | ingress class name                                                                                                              | `""`                            |
| `kdoctorAgent.ingress.route`                                   | the route of agent ingress. Default to "/kdoctoragent", if it changes, ingress please re-write url forwarded to "/kdoctoragent" | `/kdoctoragent`                 |
| `kdoctorAgent.priorityClassName`                               | the priority Class Name for kdoctorAgent                                                                                        | `system-node-critical`          |
| `kdoctorAgent.reportHostPath`                                  | storage path when pvc is disabled                                                                                               | `/var/run/kdoctor/agent`        |
| `kdoctorAgent.affinity`                                        | the affinity of kdoctorAgent                                                                                                    | `{}`                            |
| `kdoctorAgent.extraArgs`                                       | the additional arguments of kdoctorAgent container                                                                              | `[]`                            |
| `kdoctorAgent.extraEnv`                                        | the additional environment variables of kdoctorAgent container                                                                  | `[]`                            |
| `kdoctorAgent.extraVolumes`                                    | the additional volumes of kdoctorAgent container                                                                                | `[]`                            |
| `kdoctorAgent.extraVolumeMounts`                               | the additional hostPath mounts of kdoctorAgent container                                                                        | `[]`                            |
| `kdoctorAgent.podAnnotations`                                  | the additional annotations of kdoctorAgent pod                                                                                  | `{}`                            |
| `kdoctorAgent.podLabels`                                       | the additional label of kdoctorAgent pod                                                                                        | `{}`                            |
| `kdoctorAgent.resources.limits.cpu`                            | the cpu limit of kdoctorAgent pod                                                                                               | `1000m`                         |
| `kdoctorAgent.resources.limits.memory`                         | the memory limit of kdoctorAgent pod                                                                                            | `1024Mi`                        |
| `kdoctorAgent.securityContext`                                 | the security Context of kdoctorAgent pod                                                                                        | `{}`                            |
| `kdoctorAgent.grpcServer.port`                                 | the Port for grpc server                                                                                                        | `3000`                          |
| `kdoctorAgent.httpServer.healthPort`                           | the http Port for kdoctorAgent, for health checking                                                                             | `5710`                          |
| `kdoctorAgent.httpServer.appHttpPort`                          | the http Port for kdoctorAgent, testing connect                                                                                 | `80`                            |
| `kdoctorAgent.httpServer.appHttpsPort`                         | the https Port for kdoctorAgent, testing connect                                                                                | `443`                           |
| `kdoctorAgent.httpServer.startupProbe.failureThreshold`        | the failure threshold of startup probe for kdoctorAgent health checking                                                         | `60`                            |
| `kdoctorAgent.httpServer.startupProbe.periodSeconds`           | the period seconds of startup probe for kdoctorAgent health checking                                                            | `2`                             |
| `kdoctorAgent.httpServer.livenessProbe.failureThreshold`       | the failure threshold of startup probe for kdoctorAgent health checking                                                         | `6`                             |
| `kdoctorAgent.httpServer.livenessProbe.periodSeconds`          | the period seconds of startup probe for kdoctorAgent health checking                                                            | `10`                            |
| `kdoctorAgent.httpServer.readinessProbe.failureThreshold`      | the failure threshold of startup probe for kdoctorAgent health checking                                                         | `3`                             |
| `kdoctorAgent.httpServer.readinessProbe.periodSeconds`         | the period seconds of startup probe for kdoctorAgent health checking                                                            | `10`                            |
| `kdoctorAgent.prometheus.enabled`                              | enable template agent to collect metrics                                                                                        | `false`                         |
| `kdoctorAgent.prometheus.port`                                 | the metrics port of template agent                                                                                              | `5711`                          |
| `kdoctorAgent.prometheus.serviceMonitor.install`               | install serviceMonitor for template agent. This requires the prometheus CRDs to be available                                    | `false`                         |
| `kdoctorAgent.prometheus.serviceMonitor.namespace`             | the serviceMonitor namespace. Default to the namespace of helm instance                                                         | `""`                            |
| `kdoctorAgent.prometheus.serviceMonitor.annotations`           | the additional annotations of kdoctorAgent serviceMonitor                                                                       | `{}`                            |
| `kdoctorAgent.prometheus.serviceMonitor.labels`                | the additional label of kdoctorAgent serviceMonitor                                                                             | `{}`                            |
| `kdoctorAgent.prometheus.prometheusRule.install`               | install prometheusRule for template agent. This requires the prometheus CRDs to be available                                    | `false`                         |
| `kdoctorAgent.prometheus.prometheusRule.namespace`             | the prometheusRule namespace. Default to the namespace of helm instance                                                         | `""`                            |
| `kdoctorAgent.prometheus.prometheusRule.annotations`           | the additional annotations of kdoctorAgent prometheusRule                                                                       | `{}`                            |
| `kdoctorAgent.prometheus.prometheusRule.labels`                | the additional label of kdoctorAgent prometheusRule                                                                             | `{}`                            |
| `kdoctorAgent.prometheus.grafanaDashboard.install`             | install grafanaDashboard for template agent. This requires the prometheus CRDs to be available                                  | `false`                         |
| `kdoctorAgent.prometheus.grafanaDashboard.namespace`           | the grafanaDashboard namespace. Default to the namespace of helm instance                                                       | `""`                            |
| `kdoctorAgent.prometheus.grafanaDashboard.annotations`         | the additional annotations of kdoctorAgent grafanaDashboard                                                                     | `{}`                            |
| `kdoctorAgent.prometheus.grafanaDashboard.labels`              | the additional label of kdoctorAgent grafanaDashboard                                                                           | `{}`                            |
| `kdoctorAgent.debug.logLevel`                                  | the log level of template agent [debug, info, warn, error, fatal, panic]                                                        | `info`                          |
| `kdoctorAgent.debug.gopsPort`                                  | the gops port of template agent                                                                                                 | `5712`                          |
| `kdoctorController.name`                                       | the kdoctorController name                                                                                                      | `kdoctor-controller`            |
| `kdoctorController.replicas`                                   | the replicas number of kdoctorController pod                                                                                    | `1`                             |
| `kdoctorController.cmdBinName`                                 | the binName name of kdoctorController                                                                                           | `/usr/bin/controller`           |
| `kdoctorController.hostnetwork`                                | enable hostnetwork mode of kdoctorController pod. Notice, if no CNI available before template installation, must enable this    | `false`                         |
| `kdoctorController.image.registry`                             | the image registry of kdoctorController                                                                                         | `ghcr.io`                       |
| `kdoctorController.image.repository`                           | the image repository of kdoctorController                                                                                       | `kdoctor-io/kdoctor-controller` |
| `kdoctorController.image.pullPolicy`                           | the image pullPolicy of kdoctorController                                                                                       | `IfNotPresent`                  |
| `kdoctorController.image.digest`                               | the image digest of kdoctorController, which takes preference over tag                                                          | `""`                            |
| `kdoctorController.image.tag`                                  | the image tag of kdoctorController, overrides the image tag whose default is the chart appVersion.                              | `""`                            |
| `kdoctorController.image.imagePullSecrets`                     | the image imagePullSecrets of kdoctorController                                                                                 | `[]`                            |
| `kdoctorController.serviceAccount.create`                      | create the service account for the kdoctorController                                                                            | `true`                          |
| `kdoctorController.serviceAccount.annotations`                 | the annotations of kdoctorController service account                                                                            | `{}`                            |
| `kdoctorController.service.annotations`                        | the annotations for kdoctorController service                                                                                   | `{}`                            |
| `kdoctorController.service.type`                               | the type for kdoctorController service                                                                                          | `ClusterIP`                     |
| `kdoctorController.priorityClassName`                          | the priority Class Name for kdoctorController                                                                                   | `system-node-critical`          |
| `kdoctorController.affinity`                                   | the affinity of kdoctorController                                                                                               | `{}`                            |
| `kdoctorController.extraArgs`                                  | the additional arguments of kdoctorController container                                                                         | `[]`                            |
| `kdoctorController.extraEnv`                                   | the additional environment variables of kdoctorController container                                                             | `[]`                            |
| `kdoctorController.extraVolumes`                               | the additional volumes of kdoctorController container                                                                           | `[]`                            |
| `kdoctorController.extraVolumeMounts`                          | the additional hostPath mounts of kdoctorController container                                                                   | `[]`                            |
| `kdoctorController.podAnnotations`                             | the additional annotations of kdoctorController pod                                                                             | `{}`                            |
| `kdoctorController.podLabels`                                  | the additional label of kdoctorController pod                                                                                   | `{}`                            |
| `kdoctorController.securityContext`                            | the security Context of kdoctorController pod                                                                                   | `{}`                            |
| `kdoctorController.resources.limits.cpu`                       | the cpu limit of kdoctorController pod                                                                                          | `500m`                          |
| `kdoctorController.resources.limits.memory`                    | the memory limit of kdoctorController pod                                                                                       | `1024Mi`                        |
| `kdoctorController.resources.requests.cpu`                     | the cpu requests of kdoctorController pod                                                                                       | `100m`                          |
| `kdoctorController.resources.requests.memory`                  | the memory requests of kdoctorController pod                                                                                    | `128Mi`                         |
| `kdoctorController.podDisruptionBudget.enabled`                | enable podDisruptionBudget for kdoctorController pod                                                                            | `false`                         |
| `kdoctorController.podDisruptionBudget.minAvailable`           | minimum number/percentage of pods that should remain scheduled.                                                                 | `1`                             |
| `kdoctorController.httpServer.port`                            | the http Port for kdoctorController, for health checking and http service                                                       | `80`                            |
| `kdoctorController.httpServer.startupProbe.failureThreshold`   | the failure threshold of startup probe for kdoctorController health checking                                                    | `30`                            |
| `kdoctorController.httpServer.startupProbe.periodSeconds`      | the period seconds of startup probe for kdoctorController health checking                                                       | `2`                             |
| `kdoctorController.httpServer.livenessProbe.failureThreshold`  | the failure threshold of startup probe for kdoctorController health checking                                                    | `6`                             |
| `kdoctorController.httpServer.livenessProbe.periodSeconds`     | the period seconds of startup probe for kdoctorController health checking                                                       | `10`                            |
| `kdoctorController.httpServer.readinessProbe.failureThreshold` | the failure threshold of startup probe for kdoctorController health checking                                                    | `3`                             |
| `kdoctorController.httpServer.readinessProbe.periodSeconds`    | the period seconds of startup probe for kdoctorController health checking                                                       | `10`                            |
| `kdoctorController.webhookPort`                                | the http port for kdoctorController webhook                                                                                     | `5722`                          |
| `kdoctorController.prometheus.enabled`                         | enable template Controller to collect metrics                                                                                   | `false`                         |
| `kdoctorController.prometheus.port`                            | the metrics port of template Controller                                                                                         | `5721`                          |
| `kdoctorController.prometheus.serviceMonitor.install`          | install serviceMonitor for template agent. This requires the prometheus CRDs to be available                                    | `false`                         |
| `kdoctorController.prometheus.serviceMonitor.namespace`        | the serviceMonitor namespace. Default to the namespace of helm instance                                                         | `""`                            |
| `kdoctorController.prometheus.serviceMonitor.annotations`      | the additional annotations of kdoctorController serviceMonitor                                                                  | `{}`                            |
| `kdoctorController.prometheus.serviceMonitor.labels`           | the additional label of kdoctorController serviceMonitor                                                                        | `{}`                            |
| `kdoctorController.prometheus.prometheusRule.install`          | install prometheusRule for template agent. This requires the prometheus CRDs to be available                                    | `false`                         |
| `kdoctorController.prometheus.prometheusRule.namespace`        | the prometheusRule namespace. Default to the namespace of helm instance                                                         | `""`                            |
| `kdoctorController.prometheus.prometheusRule.annotations`      | the additional annotations of kdoctorController prometheusRule                                                                  | `{}`                            |
| `kdoctorController.prometheus.prometheusRule.labels`           | the additional label of kdoctorController prometheusRule                                                                        | `{}`                            |
| `kdoctorController.prometheus.grafanaDashboard.install`        | install grafanaDashboard for template agent. This requires the prometheus CRDs to be available                                  | `false`                         |
| `kdoctorController.prometheus.grafanaDashboard.namespace`      | the grafanaDashboard namespace. Default to the namespace of helm instance                                                       | `""`                            |
| `kdoctorController.prometheus.grafanaDashboard.annotations`    | the additional annotations of kdoctorController grafanaDashboard                                                                | `{}`                            |
| `kdoctorController.prometheus.grafanaDashboard.labels`         | the additional label of kdoctorController grafanaDashboard                                                                      | `{}`                            |
| `kdoctorController.debug.logLevel`                             | the log level of template Controller [debug, info, warn, error, fatal, panic]                                                   | `info`                          |
| `kdoctorController.debug.gopsPort`                             | the gops port of template Controller                                                                                            | `5724`                          |
| `kdoctorController.apiserver.name`                             | the kdoctorApiserver name                                                                                                       | `kdoctor-apiserver`             |
| `tls.ca.secretName`                                            | the secret name for storing TLS certificates                                                                                    | `kdoctor-ca`                    |
| `tls.client.secretName`                                        | the secret name for storing TLS certificates                                                                                    | `kdoctor-client-cert`           |
| `tls.server.method`                                            | the method for generating TLS certificates. [ provided , certmanager , auto]                                                    | `auto`                          |
| `tls.server.secretName`                                        | the secret name for storing TLS certificates                                                                                    | `kdoctor-controller-cert`       |
| `tls.server.certmanager.certValidityDuration`                  | generated certificates validity duration in days for 'certmanager' method                                                       | `365`                           |
| `tls.server.certmanager.issuerName`                            | issuer name of cert manager 'certmanager'. If not specified, a CA issuer will be created.                                       | `""`                            |
| `tls.server.certmanager.extraDnsNames`                         | extra DNS names added to certificate when it's auto generated                                                                   | `[]`                            |
| `tls.server.certmanager.extraIPAddresses`                      | extra IP addresses added to certificate when it's auto generated                                                                | `[]`                            |
| `tls.server.provided.tlsCert`                                  | encoded tls certificate for provided method                                                                                     | `""`                            |
| `tls.server.provided.tlsKey`                                   | encoded tls key for provided method                                                                                             | `""`                            |
| `tls.server.provided.tlsCa`                                    | encoded tls CA for provided method                                                                                              | `""`                            |
| `tls.server.auto.caExpiration`                                 | ca expiration for auto method                                                                                                   | `73000`                         |
| `tls.server.auto.certExpiration`                               | server cert expiration for auto method                                                                                          | `73000`                         |
| `tls.server.auto.extraIpAddresses`                             | extra IP addresses of server certificate for auto method                                                                        | `[]`                            |
| `tls.server.auto.extraDnsNames`                                | extra DNS names of server cert for auto method                                                                                  | `[]`                            |
