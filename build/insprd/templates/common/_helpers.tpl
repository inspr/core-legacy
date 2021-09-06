{{- define "common.healthcheck" -}}
livenessProbe:
  httpGet:
    path: /healthz
    port: {{ .service.targetPort }}
  periodSeconds: 10
readinessProbe:
  httpGet:
    path: /healthz
    port: {{ .service.targetPort }}
  periodSeconds: 10
{{- end -}}
