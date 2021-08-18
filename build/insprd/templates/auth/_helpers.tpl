{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "auth.fullname" -}}
{{ if .Values.auth.fullNameOverride -}}
{{ .Values.auth.fullNameOverride }}
{{ else -}}
{{ printf "%s-%s" .Release.Name .Values.auth.name }}
{{- end }}
{{- end }}

{{- define "auth.labels"}}
{{- include "common.labels" $ }}
app: {{ include "auth.fullname" $ }}
{{- end -}}
{{- define "auth.healthcheck" -}}
{{ include "common.healthcheck" .Values.auth }}
{{- end -}}
