{{- define "secret.name" }}
{{- printf "%s-%s" .Release.Name .secretName }}
{{- end }}

{{- define "secret.labels" }}
labels: {{- include "common.labels.standard" .context | nindent 2 }}
  {{- if .commonLabels }}
  {{- include "common.tplvalues.render" ( dict "value" .commonLabels "context" .context ) | nindent 2 }}
  {{- end }}
{{- end }}

{{- define "secret.annotations" }}
{{- if .commonAnnotations }}
annotations: {{- include "common.tplvalues.render" ( dict "value" .commonAnnotations "context" .context ) | nindent 2 }}
{{- end }}
{{- end }}

{{- define "secret.data" }}
data:
{{- range $key, $val := .data }}
  {{ $key }}: {{ $val | b64enc | quote }}
{{- end }}
{{- end }}
