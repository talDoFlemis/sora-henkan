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
  value: {{ .Values.database.port | quote }}
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
- name: {{ .prefix }}_DATABASE_DATABASE
  value: {{ .Values.database.name }}
- name: {{ .prefix }}_DATABASE_SSLMODE
  value: {{ .Values.database.sslmode | quote }}
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
{{- if .Values.dynamodb.enabled }}
{{- $prefix := .prefix | upper -}}
- name: {{ .prefix }}_DYNAMODBLOGS_ENABLED
  value: "true"
- name: {{ .prefix }}_DYNAMODBLOGS_TABLE
  value: {{ .Values.dynamodb.table | quote }}
- name: {{ .prefix }}_DYNAMODBLOGS_AWS_ENDPOINT
  value: {{ .Values.dynamodb.aws.endpoint | quote }}
- name: {{ .prefix }}_DYNAMODBLOGS_AWS_ANONYMOUS
  value: {{ .Values.dynamodb.aws.anonymous | quote }}
- name: {{ .prefix }}_DYNAMODBLOGS_AWS_REGION
  value: {{ .Values.dynamodb.aws.region | quote }}
- name: {{ .prefix }}_DYNAMODBLOGS_AWS_ACCESSKEY
{{- if and .Values.dynamodb.aws.secretName .Values.dynamodb.accessKeyKey }}
  valueFrom:
    secretKeyRef: 
      name: {{ .Values.dynamodb.aws.secretName }}
      key: {{ .Values.dynamodb.aws.accessKeyKey }}
{{- else }}
  value: {{ .Values.dynamodb.aws.accessKey }}
{{- end }}
- name: {{ .prefix }}_DYNAMODBLOGS_AWS_SECRETKEY
{{- if and .Values.dynamodb.aws.secretName .Values.dynamodb.secretKeyKey }}
  valueFrom:
    secretKeyRef: 
      name: {{ .Values.dynamodb.aws.secretName }}
      key: {{ .Values.dynamodb.aws.secretKeyKey }}
{{- else }}
  value: {{ .Values.dynamodb.aws.secretKey }}
{{- end }}
{{- end }}
{{- end }}

{{- define "backend.watermill" -}}
{{- $prefix := .prefix | upper -}}
- name: {{ .prefix }}_WATERMILL_IMAGETOPIC
  value: {{ .Values.watermill.imageTopic | quote }}
- name: {{ .prefix }}_WATERMILL_BROKER_KIND
  value: {{ .Values.watermill.broker.kind | quote }}

{{- if .Values.watermill.broker.kind | eq "aws" }}
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
  value: {{ .Values.watermill.broker.aws.accessKey }}
{{- end }}
- name: {{ .prefix }}_WATERMILL_BROKER_AWS_SECRETKEY
{{- if and .Values.watermill.broker.aws.secretName .Values.watermill.broker.aws.secretKeyKey }}
  valueFrom:
    secretKeyRef: 
      name: {{ .Values.watermill.broker.aws.secretName }}
      key: {{ .Values.watermill.broker.aws.secretKeyKey }}
{{- else }}
  value: {{ .Values.watermill.broker.aws.secretKey }}
{{- end }}
- name: {{ .prefix }}_WATERMILL_BROKER_AWS_REGION
  value: {{ .Values.watermill.broker.aws.region }}

{{- else if .Values.watermill.broker.kind | eq "amqp" }}
- name: {{ .prefix }}_WATERMILL_BROKER_AMQP_HOST
  value: {{ .Values.watermill.broker.amqp.host | quote }}
- name: {{ .prefix }}_WATERMILL_BROKER_AMQP_PORT
  value: {{ .Values.watermill.broker.amqp.port | quote }}
- name: {{ .prefix }}_WATERMILL_BROKER_AMQP_USER
{{- if and .Values.watermill.broker.amqp.secretName .Values.watermill.broker.amqp.userKey }}
  valueFrom:
    secretKeyRef: 
      name: {{ .Values.watermill.broker.amqp.secretName }}
      key: {{ .Values.watermill.broker.amqp.userKey }}
{{- else }}
  value: {{ .Values.watermill.broker.amqp.user | quote }}
{{- end }}
- name: {{ .prefix }}_WATERMILL_BROKER_AMQP_PASSWORD
{{- if and .Values.watermill.broker.amqp.secretName .Values.watermill.broker.amqp.passwordKey }}
  valueFrom:
    secretKeyRef: 
      name: {{ .Values.watermill.broker.amqp.secretName }}
      key: {{ .Values.watermill.broker.amqp.passwordKey }}
{{- else }}
  value: {{ .Values.watermill.broker.amqp.password | quote }}
{{- end }}
- name: {{ .prefix }}_WATERMILL_BROKER_AMQP_VHOST
  value: {{ .Values.watermill.broker.amqp.vhost | quote }}
{{- end }}
{{- end }}
