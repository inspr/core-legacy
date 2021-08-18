{{/*
Name of the sidecar
*/}}
{{- define "lbsidecar.fullname" -}}
{{ .Release.Name }}-sidecar
{{- end -}}

