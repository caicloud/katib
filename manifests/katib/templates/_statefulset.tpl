{{- define "sts.name" }}
{{- printf "%s-%s" .Release.Name .stsName }}
{{- end }}

{{- define "sts.labels" }}
labels: {{- include "common.labels.standard" .context | nindent 2 }}
  {{- if .commonLabels }}
  {{- include "common.tplvalues.render" ( dict "value" .commonLabels "context" .context ) | nindent 2 }}
  {{- end }}
{{- end }}

{{- define "sts.annotations" }}
{{- if .commonAnnotations }}
annotations: {{- include "common.tplvalues.render" ( dict "value" .commonAnnotations "context" .context ) | nindent 2 }}
{{- end }}
{{- end }}

{{- define "sts.serviceName" }}
{{- printf "%s-%s-headless" .Release.Name .stsName }}
{{- end }}

{{- define "sts.vcts" }}
{{- if .vcts }}
volumeClaimTemplates:
{{- range $vctName, $vct := .vcts }}
- metadata:
    name: {{ $vctName }}
    labels: {{- include "common.labels.matchLabels" $.context | nindent 6 }}
  spec:
    accessModes:
    - {{ $vct.accessModes | quote }}
    resources:
      requests:
        storage: {{ $vct.size | quote }}
    storageClassName: {{ $vct.storageClassName | quote }}
{{- end }}
{{- end }}
{{- end }}
