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

{{- define "inspr-stack.prometheus.labels" -}}
{{- include "common.labels" . }}
app: {{ include "inspr-stack.prometheus.name" . }}
release: {{ .Release.Name }}
{{- end -}}
{{- define "inspr-stack.prometheus.selector-labels" -}}
app: {{ include "inspr-stack.prometheus.name" .}}
release: {{ .Release.Name }}
{{- end -}}
