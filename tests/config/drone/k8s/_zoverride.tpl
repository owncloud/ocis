{{/*

DO NOT RENAME THIS FILE

Filename: _zoverride.tpl

Override filename starts with 'z'
to make sure it is loaded after all other templates.

Using the existing variable definitions
which are already included in most of the templates.
*/}}

{{- define "ocis.caPath" -}}
- name: logdir
  mountPath: /logs
{{- end -}}

{{- define "ocis.caVolume" -}}
- name: logdir
  hostPath:
    path: /logs
    type: Directory
{{- end -}}

{{- define "ocis.caEnv" -}}
- name: OCIS_LOG_FILE
  value: /logs/ocis.log
{{- end -}}
