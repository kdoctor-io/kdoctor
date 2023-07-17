{{- define "project.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}
{{/*
Expand the name of project .
*/}}
{{- define "project.name" -}}
{{- .Values.appName | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "project.selectorLabels" -}}
app.kubernetes.io/name: {{ include "project.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/component: {{ .Values.appName | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
app Common labels
*/}}
{{- define "project.labels" -}}
helm.sh/chart: {{ include "project.chart" . }}
{{ include "project.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}


{{/*
return the image
*/}}
{{- define "project.image" -}}
{{- $registryName := .Values.image.registry -}}
{{- $repositoryName := .Values.image.repository -}}
{{ if $registryName }}
    {{- printf "%s/%s" $registryName $repositoryName -}}
{{- else -}}
    {{- printf "%s" $repositoryName -}}
{{- end -}}
{{- if .Values.image.digest }}
    {{- print "@" .Values.image.digest -}}
{{- else if .Values.image.tag -}}
    {{- printf ":%s" .Values.image.tag -}}
{{- else -}}
    {{- printf ":v%s" .Chart.AppVersion -}}
{{- end -}}
{{- end -}}


{{/*
generate the CA cert , 200 years
*/}}
{{- define "generate-ca-certs" }}
    {{- $ca := genCA "kdoctor.io" 73000 -}}
    {{- $_ := set . "ca" $ca -}}
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