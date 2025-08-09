{{/*
Expand the name of the chart.
*/}}
{{- define "titanic-go-service.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "titanic-go-service.fullname" -}}
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
{{- define "titanic-go-service.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "titanic-go-service.labels" -}}
helm.sh/chart: {{ include "titanic-go-service.chart" . }}
{{ include "titanic-go-service.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "titanic-go-service.selectorLabels" -}}
app.kubernetes.io/name: {{ include "titanic-go-service.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}
{{- define "titanic-go-service.configmapdata" -}}
server:
  port: "8080"
data:
  source: "{{ .Values.config.dataSource }}"
  csv_file: "/data/titanic.csv"
  db_file: "/data/titanic.db"
{{- end -}}

{{- define "titanic-go-service.validateValues" -}}
{{- $allowedDataSources := list "csv" "sqlite" -}}
{{- if not (has .Values.config.dataSource $allowedDataSources) -}}
{{- $message := printf "Invalid config.dataSource: '%s'. Allowed values are 'csv' or 'sqlite'." .Values.config.dataSource -}}
{{- fail $message -}}
{{- end -}}
{{- end -}}
