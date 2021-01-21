{{/*
Return the proper image name
{{ include "common.images.image" ( dict "context" $ "repository" .Values.platformConfig.imageRepositoryRelease "imageRoot" .Values.image ) }}
*/}}
{{- define "common.images.image" -}}
{{- printf "%s/%s/%s:%s" .context.Values.platformConfig.imageRegistry .repository .imageRoot.name ( default .context.Chart.AppVersion .imageRoot.tag ) -}}
{{- end -}}
