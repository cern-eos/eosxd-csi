{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "eosxd-csi.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "eosxd-csi.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "eosxd-csi.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create fully qualified app name for the node plugin.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "eosxd-csi.nodeplugin.fullname" -}}
{{- if .Values.nodeplugin.fullnameOverride -}}
{{- .Values.nodeplugin.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- printf "%s-%s" .Release.Name .Values.nodeplugin.name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s-%s" .Release.Name $name .Values.nodeplugin.name  | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create fully qualified app name for the controller plugin.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "eosxd-csi.controllerplugin.fullname" -}}
{{- if .Values.controllerplugin.fullnameOverride -}}
{{- .Values.controllerplugin.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- printf "%s-%s" .Release.Name .Values.controllerplugin.name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s-%s" .Release.Name $name .Values.controllerplugin.name  | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create the name of the service account to use for the node plugin.
*/}}
{{- define "eosxd-csi.serviceAccountName.nodeplugin" -}}
{{- .Values.nodeplugin.serviceAccount.serviceAccountName | default (include "eosxd-csi.nodeplugin.fullname" .) -}}
{{- end -}}

{{/*
Create the name of the service account to use for the node plugin.
*/}}
{{- define "eosxd-csi.serviceAccountName.controllerplugin" -}}
{{- .Values.controllerplugin.serviceAccount.serviceAccountName | default (include "eosxd-csi.controllerplugin.fullname" .) -}}
{{- end -}}

{{/*
Create unified labels for eosxd-csi components
*/}}

{{- define "eosxd-csi.common.metaLabels" -}}
app.kubernetes.io/name: {{ template "eosxd-csi.name" . }}
helm.sh/chart: {{ .Chart.Name }}-{{ .Chart.Version | replace "+" "_" }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/version: {{ .Chart.AppVersion }}
{{- if .Values.extraMetaLabels -}}
{{ toYaml .Values.extraMetaLabels }}
{{- end -}}
{{- end -}}

{{- define "eosxd-csi.common.matchLabels" -}}
app.kubernetes.io/name: {{ template "eosxd-csi.name" . }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{- define "eosxd-csi.nodeplugin.metaLabels" -}}
{{- if .Values.nodeplugin.metaLabelsOverride -}}
{{ toYaml .Values.nodeplugin.metaLabelsOverride }}
{{- else -}}
app.kubernetes.io/component: nodeplugin
{{ include "eosxd-csi.common.metaLabels" . }}
{{- end -}}
{{- end -}}

{{- define "eosxd-csi.nodeplugin.matchLabels" -}}
{{- if .Values.nodeplugin.matchLabelsOverride -}}
{{ toYaml .Values.nodeplugin.matchLabelsOverride }}
{{- else -}}
app.kubernetes.io/component: nodeplugin
{{ include "eosxd-csi.common.matchLabels" . }}
{{- end -}}
{{- end -}}

{{- define "eosxd-csi.controllerplugin.metaLabels" -}}
{{- if .Values.controllerplugin.metaLabelsOverride -}}
{{ toYaml .Values.controllerplugin.metaLabelsOverride }}
{{- else -}}
app.kubernetes.io/component: controllerplugin
{{ include "eosxd-csi.common.metaLabels" . }}
{{- end -}}
{{- end -}}

{{- define "eosxd-csi.controllerplugin.matchLabels" -}}
{{- if .Values.controllerplugin.matchLabelsOverride -}}
{{ toYaml .Values.controllerplugin.matchLabelsOverride }}
{{- else -}}
app.kubernetes.io/component: controllerplugin
{{ include "eosxd-csi.common.matchLabels" . }}
{{- end -}}
{{- end -}}
