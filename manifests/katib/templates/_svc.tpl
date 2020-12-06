{{- define "svc.name" }}
{{- printf "%s-%s" .Release.Name .svcName }}
{{- end }}

{{- define "svc.labels" }}
labels: {{- include "common.labels.standard" .context | nindent 2 }}
  {{- if .commonLabels }}
  {{- include "common.tplvalues.render" ( dict "value" .commonLabels "context" .context ) | nindent 2 }}
  {{- end }}
{{- end}}

{{- define "svc.annotations" }}
{{- if .commonAnnotations }}
annotations: {{- include "common.tplvalues.render" ( dict "value" .commonAnnotations "context" .context ) | nindent 2 }}
{{- end }}
{{- end }}

{{- define "svc.selector" }}
selector: {{- include "common.labels.matchLabels" .context | nindent 2 }}
  {{- if .svc.selector }}
  {{- include "common.tplvalues.render" ( dict "value" .svc.selector "context" .context ) | nindent 2 }}
  {{- end }}
{{- end}}

{{- define "svc.clusterIP" }}
{{- if and .svc.clusterIP (eq .svc.type "ClusterIP") }}
clusterIP: {{ .svc.clusterIP }}
{{- end }}
{{- end }}

{{- define "svc.externalTrafficPolicy" }}
{{- if (and (or (eq .svc.type "LoadBalancer") (eq .svc.type "NodePort")) .svc.externalTrafficPolicy ) }}
externalTrafficPolicy: {{ .svc.externalTrafficPolicy | quote }}
{{- end }}
{{- end }}

{{- define "svc.loadBalancerSourceRanges" }}
{{- if (and (eq .svc.type "LoadBalancer") .svc.loadBalancerSourceRanges) }}
loadBalancerSourceRanges: {{- toYaml .svc.loadBalancerSourceRanges | nindent 0 }}
{{- end }}
{{- end }}

{{- define "svc.loadBalancerIP" }}
{{- if (and (eq .svc.type "LoadBalancer") .svc.loadBalancerIP) }}
loadBalancerIP: {{ .svc.loadBalancerIP }}
{{- end }}
{{- end }}

{{- define "svc.ports" }}
{{- if .svc.ports }}
ports:
{{- range $portName, $port := .svc.ports }}
- name: {{ $portName }}
  port: {{ $port.port }}
  protocol: {{ $port.protocol }}
  targetPort: {{ $port.targetPort }}
  {{- if (and (or (eq $.svc.type "NodePort") (eq $.svc.type "LoadBalancer")) $port.nodePort) }}
  nodePort: {{ $port.nodePort }}
  {{- end }}
{{- end }}
{{- end }}
{{- end }}
