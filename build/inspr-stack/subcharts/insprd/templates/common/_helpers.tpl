{{- define "common.healthcheck" -}}
livenessProbe:
  httpGet:
    path: /healthz
    port: {{ .service.port }}
    initialDelaySeconds: 3
    periodSeconds: 3
readinessProbe:
  httpGet:
    path: /healthz
    port: {{ .service.port }}
    initialDelaySeconds: 5
    periodSeconds: 3
{{- end -}}
