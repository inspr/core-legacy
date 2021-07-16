{{/*
Expand the name of the chart.
*/}}
{{- define "uidp.name" -}}
{{- .Values.name | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "uidp.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := .Values.name }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "uidp.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "uidp.labels" -}}
helm.sh/chart: {{ include "uidp.chart" . }}
{{ include "uidp.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app: {{ include "uidp.fullname" . }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "uidp.selectorLabels" -}}
app.kubernetes.io/name: {{ include "uidp.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app: {{ include "uidp.fullname" . }}
{{- end }}
{{/*
Return the proper Docker Image Registry Secret Names
{{ include "common.images.pullSecrets" ( dict "images" (list .Values.path.to.the.image1, .Values.path.to.the.image2) "global" .Values.global) }}
*/}}
{{- define "common.images.pullSecrets" -}}
  {{- $pullSecrets := list }}

  {{- if .global }}
    {{- range .global.imagePullSecrets -}}
      {{- $pullSecrets = append $pullSecrets . -}}
    {{- end -}}
  {{- end -}}

  {{- range .images -}}
    {{- range .pullSecrets -}}
      {{- $pullSecrets = append $pullSecrets . -}}
    {{- end -}}
  {{- end -}}

  {{- if (not (empty $pullSecrets)) }}
imagePullSecrets:
    {{- range $pullSecrets }}
  - name: {{ . }}
    {{- end }}
  {{- end }}
{{- end -}}

{{/*
Renders the image value while overriding the image registry
*/}}
{{- define "common.images.image" -}}
{{- $registry := .global.imageRegistry | default .image.registry -}}
{{- if $registry -}}
"{{ $registry }}/{{ .image.repository }}:{{ .image.tag }}"
{{- else -}}
"{{ .image.repository }}:{{ .image.tag }}"
{{- end -}}
{{- end -}}

{{- define "insprd.init-check" -}}
- name: init-insprd
  image: curlimages/curl:7.75.0
  imagePullPolicy: IfNotPresent
  command:
    - "sh"
    - "-c"
    - "until curl -s {{ include "insprd.address" $ }}/heathz; do echo waiting for insprd; sleep 2; done"
{{- end -}}


{{- define "insprd.address" -}}
{{ tpl .Values.insprd.address . }}
{{- end -}}
{{- define "uidp.healthcheck" -}}
livenessProbe:
  httpGet:
    path: /healthz
    port: {{ .service.targetPort }}
  initialDelaySeconds: 3
  periodSeconds: 3
readinessProbe:
  httpGet:
    path: /healthz
    port: {{ .service.targetPort }}
  initialDelaySeconds: 5
  periodSeconds: 3
{{- end -}}
