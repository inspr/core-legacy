{{/*
Expand the name of the chart.
*/}}
{{- define "insprd.name" -}}
{{- default .Chart.Name .Values.name | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "insprd.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.name }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{- define "insprd.labels" }}
{{- include "common.labels" $ }}
app: {{ include "insprd.fullname" $ }}
app.kubernetes.io/name: {{ include "insprd.name" $ }}
{{- end }}

{{- define "insprd.selectorLabels" }}
{{- include "common.labels" $ }}
app: {{ include "insprd.fullname" $ }}
{{- end }}

{{- define "insprd.healthcheck" -}}
{{ include "common.healthcheck" .Values }}
{{- end -}}
