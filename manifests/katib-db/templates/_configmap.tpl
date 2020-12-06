{{- define "cm.name" }}
{{- printf "%s-%s" .Release.Name .cmName }}
{{- end }}

{{- define "cm.labels" }}
labels: {{- include "common.labels.standard" .context | nindent 2 }}
  {{- if .commonLabels }}
  {{- include "common.tplvalues.render" ( dict "value" .commonLabels "context" .context ) | nindent 2 }}
  {{- end }}
{{- end }}

{{- define "cm.annotations" }}
{{- if .commonAnnotations }}
annotations: {{- include "common.tplvalues.render" ( dict "value" .commonAnnotations "context" .context ) | nindent 2 }}
{{- end }}
{{- end }}
