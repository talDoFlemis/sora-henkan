{{/*
Expand the name of the chart.
*/}}
{{- define "backend.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "backend.fullname" -}}
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

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "backend.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "backend.labels" -}}
helm.sh/chart: {{ include "backend.chart" . }}
{{ include "backend.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "backend.selectorLabels" -}}
app.kubernetes.io/name: {{ include "backend.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "backend.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "backend.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{- define "backend.database" -}}
{{- $prefix := .prefix | upper -}}
- name: {{ .prefix }}_DATABASE_HOST
  value: {{ .Values.database.host }}
- name: {{ .prefix }}_DATABASE_PORT
  value: {{ .Values.database.port }}
- name: {{ .prefix }}_DATABASE_USER
  value: {{ .Values.database.user }}
- name: {{ .prefix }}_DATABASE_PASSWORD
{{- if and .Values.database.passwordSecret .Values.database.passwordSecretKey }}
  valueFrom:
    secretKeyRef: 
      name: {{ .Values.database.passwordSecret }}
      key: {{ .Values.database.passwordSecretKey }}
{{- else }}
  value: {{ .Values.database.password }}
{{- end }}
- name: {{ .prefix }}_DATABASE_NAME
  value: {{ .Values.database.name }}
- name: {{ .prefix }}_DATABASE_SSLMODE
  value: {{ .Values.database.sslmode }}
{{- end }}

{{- define "backend.objectstorer" -}}
{{- $prefix := .prefix | upper -}}
- name: {{ .prefix }}_OBJECTSTORER_ENDPOINT
  value: {{ .Values.objectstorer.endpoint | quote }}
- name: {{ .prefix }}_OBJECTSTORER_SSL
  value: {{ .Values.objectstorer.ssl | quote }}
- name: {{ .prefix }}_OBJECTSTORER_REGION
  value: {{ .Values.objectstorer.region }}
- name: {{ .prefix }}_OBJECTSTORER_ACCESSKEYID
{{- if and .Values.objectstorer.secretName .Values.objectstorer.accessKeyKey }}
  valueFrom:
    secretKeyRef: 
      name: {{ .Values.objectstorer.secretName }}
      key: {{ .Values.objectstorer.accessKeyKey }}
{{- else }}
  value: {{ .Values.objectstorer.accessKeyId }}
{{- end }}
- name: {{ .prefix }}_OBJECTSTORER_SECRETACCESSKEY
{{- if and .Values.objectstorer.secretName .Values.objectstorer.secretKeyKey }}
  valueFrom:
    secretKeyRef: 
      name: {{ .Values.objectstorer.secretName }}
      key: {{ .Values.objectstorer.secretKeyKey }}
{{- else }}
  value: {{ .Values.objectstorer.secretAccessKey }}
{{- end }}
{{- end }}

{{- define "backend.opentelemetry" -}}
{{- $prefix := .prefix | upper -}}
- name: {{ .prefix }}_OPENTELEMETRY_ENABLED
  value: {{ .Values.opentelemetry.enabled | quote }}
- name: {{ .prefix }}_OPENTELEMETRY_ENDPOINT
  value: {{ .Values.opentelemetry.endpoint | quote }}
{{- end }}

{{- define "backend.dynamodb" -}}
{{- $prefix := .prefix | upper -}}
- name: {{ .prefix }}_DYNAMODB_ENDPOINT
  value: {{ .Values.dynamodb.endpoint | quote }}
- name: {{ .prefix }}_DYNAMODB_TABLE
  value: {{ .Values.dynamodb.table | quote }}
- name: {{ .prefix }}_DYNAMODB_ANONYMOUS
  value: {{ .Values.dynamodb.anonymous | quote }}
- name: {{ .prefix }}_DYNAMODB_REGION
  value: {{ .Values.dynamodb.region | quote }}
- name: {{ .prefix }}_DYNAMODB_ACCESSKEY
{{- if and .Values.dynamodb.secretName .Values.dynamodb.accessKeyKey }}
  valueFrom:
    secretKeyRef: 
      name: {{ .Values.dynamodb.secretName }}
      key: {{ .Values.dynamodb.accessKeyKey }}
{{- else }}
  value: {{ .Values.dynamodb.accessKey }}
{{- end }}
- name: {{ .prefix }}_DYNAMODB_SECRETKEY
{{- if and .Values.dynamodb.secretName .Values.dynamodb.secretKeyKey }}
  valueFrom:
    secretKeyRef: 
      name: {{ .Values.dynamodb.secretName }}
      key: {{ .Values.dynamodb.secretKeyKey }}
{{- else }}
  value: {{ .Values.dynamodb.secretKey }}
{{- end }}
{{- end }}

{{- define "backend.watermill" -}}
{{- $prefix := .prefix | upper -}}
- name: {{ .prefix }}_WATERMILL_IMAGETOPIC
  value: {{ .Values.watermill.image-topic | quote }}
- name: {{ .prefix }}_WATERMILL_BROKER_KIND
  value: {{ .Values.watermill.broker.kind | quote }}
- name: {{ .prefix }}_WATERMILL_BROKER_AWS_ENDPOINT
  value: {{ .Values.watermill.broker.aws.endpoint | quote }}
- name: {{ .prefix }}_WATERMILL_BROKER_AWS_ANONYMOUS
  value: {{ .Values.watermill.broker.aws.anonymous | quote }}
- name: {{ .prefix }}_WATERMILL_BROKER_AWS_ACCESSKEY
{{- if and .Values.watermill.broker.aws.secretName .Values.watermill.broker.aws.accessKeyKey }}
  valueFrom:
    secretKeyRef: 
      name: {{ .Values.watermill.broker.aws.secretName }}
      key: {{ .Values.watermill.broker.aws.accessKeyKey }}
{{- else }}
  value: {{ .Values.watermill.broker.aws.access-key }}
{{- end }}
- name: {{ .prefix }}_WATERMILL_BROKER_AWS_SECRETKEY
{{- if and .Values.watermill.broker.aws.secretName .Values.watermill.broker.aws.secretKeyKey }}
  valueFrom:
    secretKeyRef: 
      name: {{ .Values.watermill.broker.aws.secretName }}
      key: {{ .Values.watermill.broker.aws.secretKeyKey }}
{{- else }}
  value: {{ .Values.watermill.broker.aws.secret-key }}
{{- end }}
- name: {{ .prefix }}_WATERMILL_BROKER_AWS_REGION
  value: {{ .Values.watermill.broker.aws.region }}
{{- end }}