{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "project.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Expand the name of project .
*/}}
{{- define "project.name" -}}
{{- .Values.global.name | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
kdoctorAgent Common labels
*/}}
{{- define "project.kdoctorAgent.labels" -}}
helm.sh/chart: {{ include "project.chart" . }}
{{ include "project.kdoctorAgent.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
kdoctorController Common labels
*/}}
{{- define "project.kdoctorController.labels" -}}
helm.sh/chart: {{ include "project.chart" . }}
{{ include "project.kdoctorController.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
kdoctorAgent Selector labels
*/}}
{{- define "project.kdoctorAgent.selectorLabels" -}}
app.kubernetes.io/name: {{ .Values.kdoctorAgent.name | trunc 63 | trimSuffix "-" }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/component: {{ .Values.kdoctorAgent.name | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
kdoctorController Selector labels
*/}}
{{- define "project.kdoctorController.selectorLabels" -}}
app.kubernetes.io/name: {{ .Values.kdoctorController.name | trunc 63 | trimSuffix "-" }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/component: {{ .Values.kdoctorController.name | trunc 63 | trimSuffix "-" }}
{{- end }}


{{/* vim: set filetype=mustache: */}}
{{/*
Renders a value that contains template.
Usage:
{{ include "tplvalues.render" ( dict "value" .Values.path.to.the.Value "context" $) }}
*/}}
{{- define "tplvalues.render" -}}
    {{- if typeIs "string" .value }}
        {{- tpl .value .context }}
    {{- else }}
        {{- tpl (.value | toYaml) .context }}
    {{- end }}
{{- end -}}




{{/*
Return the appropriate apiVersion for poddisruptionbudget.
*/}}
{{- define "capabilities.policy.apiVersion" -}}
{{- if semverCompare "<1.21-0" .Capabilities.KubeVersion.Version -}}
{{- print "policy/v1beta1" -}}
{{- else -}}
{{- print "policy/v1" -}}
{{- end -}}
{{- end -}}

{{/*
Return the appropriate apiVersion for deployment.
*/}}
{{- define "capabilities.deployment.apiVersion" -}}
{{- if semverCompare "<1.14-0" .Capabilities.KubeVersion.Version -}}
{{- print "extensions/v1beta1" -}}
{{- else -}}
{{- print "apps/v1" -}}
{{- end -}}
{{- end -}}


{{/*
Return the appropriate apiVersion for RBAC resources.
*/}}
{{- define "capabilities.rbac.apiVersion" -}}
{{- if semverCompare "<1.17-0" .Capabilities.KubeVersion.Version -}}
{{- print "rbac.authorization.k8s.io/v1beta1" -}}
{{- else -}}
{{- print "rbac.authorization.k8s.io/v1" -}}
{{- end -}}
{{- end -}}

{{/*
return the kdoctorAgent image
*/}}
{{- define "project.kdoctorAgent.image" -}}
{{- $registryName := .Values.kdoctorAgent.image.registry -}}
{{- $repositoryName := .Values.kdoctorAgent.image.repository -}}
{{- if .Values.global.imageRegistryOverride }}
    {{- printf "%s/%s" .Values.global.imageRegistryOverride $repositoryName -}}
{{ else if $registryName }}
    {{- printf "%s/%s" $registryName $repositoryName -}}
{{- else -}}
    {{- printf "%s" $repositoryName -}}
{{- end -}}
{{- if .Values.kdoctorAgent.image.digest }}
    {{- print "@" .Values.kdoctorAgent.image.digest -}}
{{- else if .Values.global.imageTagOverride -}}
    {{- printf ":%s" .Values.global.imageTagOverride -}}
{{- else if .Values.kdoctorAgent.image.tag -}}
    {{- printf ":%s" .Values.kdoctorAgent.image.tag -}}
{{- else -}}
    {{- printf ":v%s" .Chart.AppVersion -}}
{{- end -}}
{{- end -}}


{{/*
return the kdoctorController image
*/}}
{{- define "project.kdoctorController.image" -}}
{{- $registryName := .Values.kdoctorController.image.registry -}}
{{- $repositoryName := .Values.kdoctorController.image.repository -}}
{{- if .Values.global.imageRegistryOverride }}
    {{- printf "%s/%s" .Values.global.imageRegistryOverride $repositoryName -}}
{{ else if $registryName }}
    {{- printf "%s/%s" $registryName $repositoryName -}}
{{- else -}}
    {{- printf "%s" $repositoryName -}}
{{- end -}}
{{- if .Values.kdoctorController.image.digest }}
    {{- print "@" .Values.kdoctorController.image.digest -}}
{{- else if .Values.global.imageTagOverride -}}
    {{- printf ":%s" .Values.global.imageTagOverride -}}
{{- else if .Values.kdoctorController.image.tag -}}
    {{- printf ":%s" .Values.kdoctorController.image.tag -}}
{{- else -}}
    {{- printf ":v%s" .Chart.AppVersion -}}
{{- end -}}
{{- end -}}

{{/*
generate the CA cert
*/}}
{{- define "generate-ca-certs" }}
    {{- $ca := genCA "kdoctor.io" (.Values.tls.server.auto.caExpiration | int) -}}
    {{- $_ := set . "ca" $ca -}}
{{- end }}

{{- define "project.kdoctorAgent.serviceIpv4Name" -}}
{{- printf "%s-ipv4" .Values.kdoctorAgent.name -}}
{{- end -}}

{{- define "project.kdoctorAgent.serviceIpv6Name" -}}
{{- printf "%s-ipv6" .Values.kdoctorAgent.name -}}
{{- end -}}

{{- define "project.kdoctorAgent.ingressName" -}}
{{- printf "%s" .Values.kdoctorAgent.name -}}
{{- end -}}
