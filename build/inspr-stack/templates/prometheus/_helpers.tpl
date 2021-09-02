{{- define "inspr-stack.prometheus.name" -}}
{{- if .Values.prometheus.fullnameOverride }}
{{- .Values.prometheus.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default "prometheus" .Values.prometheus.name }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end -}}

{{/*
Common labels
*/}}
{{- define "inspr-stack.prometheus.labels" -}}
{{ include "inspr-stack.prometheus.selector-labels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
release: {{ .Release.Name }}
{{- if .Values.prometheus.extraLabels }}
{{ toYaml .Values.prometheus.extraLabels }}
{{- end }}
{{- end -}}

{{/*
Selector labels
*/}}
{{- define "inspr-stack.prometheus.selector-labels" -}}
app: {{ include "inspr-stack.prometheus.name" .}}
app.kubernetes.io/name: {{ include "inspr-stack.prometheus.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

