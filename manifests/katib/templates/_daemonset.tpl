{{- define "ds.name" }}
{{- printf "%s-%s" .Release.Name .dsName }}
{{- end }}

{{- define "ds.labels" }}
labels: {{- include "common.labels.standard" .context | nindent 2 }}
  {{- if .commonLabels }}
  {{- include "common.tplvalues.render" ( dict "value" .commonLabels "context" .context ) | nindent 2 }}
  {{- end }}
{{- end }}

{{- define "ds.annotations" }}
{{- if .commonAnnotations }}
annotations: {{- include "common.tplvalues.render" ( dict "value" .commonAnnotations "context" .context ) | nindent 2 }}
{{- end }}
{{- end }}
