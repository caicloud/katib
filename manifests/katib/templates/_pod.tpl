{{- define "pod.labels" }}
labels: {{- include "common.labels.standard" .context | nindent 2 }}
  {{- if .podLabels }}
  {{- include "common.tplvalues.render" ( dict "value" .podLabels "context" .context ) | nindent 2 }}
  {{- end }}
{{- end }}

{{- define "pod.annotations" }}
{{- if .podAnnotations }}
annotations: {{- include "common.tplvalues.render" (dict "value" .podAnnotations "context" .context) | nindent 2 }}
{{- end }}
{{- end }}

{{- define "pod.tolerations" }}
{{- if .tolerations }}
tolerations: {{- include "common.tplvalues.render" ( dict "value" .tolerations "context" .context ) | nindent 0 }}
{{- end }}
{{- end }}

{{- define "pod.affinity" }}
{{- if .affinity }}
affinity: {{- include "common.tplvalues.render" ( dict "value" .affinity "context" .context ) | nindent 2 }}
{{- end }}
{{- end }}

{{- define "pod.nodeSelector" }}
{{- if .nodeSelector }}
nodeSelector: {{- include "common.tplvalues.render" ( dict "value" .nodeSelector "context" .context ) | nindent 2 }}
{{- end }}
{{- end }}

{{- define "pod.volumes" }}
{{- if .volumes }}
volumes:
{{- range $volumeName, $volume := .volumes }}
- name: {{ $volumeName }}
  {{- if hasKey $volume "pvc" }}
  persistenceVolumeClaim:
    claimName: {{ $volume.pvc }}
  {{- else if hasKey $volume "configMap" }}
  configMap:
    name: {{ $volume.configMap }}
    items: {{- include "common.tplvalues.render" ( dict "value" $volume.items "context" $.context ) | nindent 4 }}
  {{- else }}
  emptyDir: {}
  {{- end }}
{{- end }}
{{- end }}
{{- end }}

{{- define "pod.initContainers" }}
{{- $context := .context }}
{{- if .initContainers }}
initContainers:
  {{- range $containerName, $container := .initContainers }}
    {{- include "container.container" ( dict "containerName" $containerName "container" $container "context" $context ) }}
  {{- end }}
{{- end }}
{{- end }}

{{- define "pod.containers" }}
{{- $context := .context }}
{{- if .containers }}
containers:
  {{- range $containerName, $container := .containers }}
    {{- include "container.container" ( dict "containerName" $containerName "container" $container "context" $context ) }}
  {{- end }}
{{- end }}
{{- end }}
