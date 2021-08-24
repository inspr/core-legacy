{{- define "uidp.admin.password" -}}
{{- if .Values.admin.generatePassword -}}
{{ randAscii 10 }}
{{- else -}}
{{ .Values.admin.password }}
{{- end -}}
{{- end -}}

{{- define "uidp.security.refresh-key" -}}
{{ randAscii 32 }}
{{- end -}}
