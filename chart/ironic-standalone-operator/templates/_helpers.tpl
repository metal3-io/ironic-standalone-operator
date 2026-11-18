{{/*
Chart name.
*/}}
{{- define "irso.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Fully qualified app name, used as a prefix for all resources.
*/}}
{{- define "irso.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{- define "irso.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "irso.selectorLabels" -}}
app.kubernetes.io/name: {{ include "irso.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{- define "irso.labels" -}}
helm.sh/chart: {{ include "irso.chart" . }}
{{ include "irso.selectorLabels" . }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/part-of: ironic-standalone-operator
{{- with .Values.commonLabels }}
{{ include "irso.stringMap" . | trim }}
{{- end }}
{{- end }}

{{/*
Render a map with all values quoted as strings (labels and annotations must
be strings, but --set turns e.g. "true" into a boolean).
*/}}
{{- define "irso.stringMap" -}}
{{- range $key, $value := . }}
{{ $key }}: {{ $value | toString | quote }}
{{- end }}
{{- end }}

{{/*
Common metadata.annotations block (indented for metadata), merged with optional
extra annotations passed as "extra". Renders nothing when empty.
Usage: include "irso.annotations" (dict "ctx" $ "extra" (dict ...))
*/}}
{{- define "irso.annotations" -}}
{{- $annotations := merge (deepCopy (default dict .extra)) .ctx.Values.commonAnnotations }}
{{- with $annotations }}
  annotations:
    {{- include "irso.stringMap" . | trim | nindent 4 }}
{{- end }}
{{- end }}

{{- define "irso.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (printf "%s-controller-manager" (include "irso.fullname" .)) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{- define "irso.image" -}}
{{- if .Values.image.digest }}
{{- printf "%s@%s" .Values.image.repository .Values.image.digest }}
{{- else }}
{{- printf "%s:%s" .Values.image.repository (default .Chart.AppVersion .Values.image.tag) }}
{{- end }}
{{- end }}

{{/*
Namespace holding the leader election lease. The operator uses the watch
namespace when set, otherwise the namespace it runs in.
*/}}
{{- define "irso.leaderElectionNamespace" -}}
{{- default .Release.Namespace .Values.watchNamespace }}
{{- end }}

{{- define "irso.webhookServiceName" -}}
{{- printf "%s-webhook-service" (include "irso.fullname" .) }}
{{- end }}

{{- define "irso.certificateName" -}}
{{- printf "%s-serving-cert" (include "irso.fullname" .) }}
{{- end }}

{{- define "irso.webhookSecretName" -}}
{{- if .Values.webhook.certManager.enabled }}
{{- printf "%s-webhook-server-cert" (include "irso.fullname" .) }}
{{- else if .Values.webhook.existingSecret.name }}
{{- .Values.webhook.existingSecret.name }}
{{- else }}
{{- fail "webhook.existingSecret.name is required when webhook.enabled=true and webhook.certManager.enabled=false" }}
{{- end }}
{{- end }}

{{/*
Annotation used by cert-manager's cainjector to fill in caBundle fields.
*/}}
{{- define "irso.caInjectionAnnotation" -}}
{{- if and .Values.webhook.enabled .Values.webhook.certManager.enabled }}
{{- dict "cert-manager.io/inject-ca-from" (printf "%s/%s" .Release.Namespace (include "irso.certificateName" .)) | toYaml }}
{{- end }}
{{- end }}

{{/*
Manager RBAC rules, read from the controller-gen output (config/rbac/role.yaml)
that `make helm-sync` copies into the chart.
*/}}
{{- define "irso.managerRules" -}}
{{- $role := .Files.Get "files/rbac/role.yaml" | fromYaml }}
{{- if not $role.rules }}
{{- fail "files/rbac/role.yaml is missing or has no rules; run `make helm-sync`" }}
{{- end }}
{{- toYaml $role.rules }}
{{- end }}
