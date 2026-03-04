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
- name: ocis-fonts-ttf
  mountPath: /etc/ocis/fonts
- name: ocis-fonts-map
  mountPath: /etc/ocis/fontsMap.json
  subPath: fontsMap.json
{{- end -}}

{{- define "ocis.caVolume" -}}
- name: logdir
  hostPath:
    path: /logs
    type: Directory
- name: ocis-fonts-ttf
  configMap:
    name: ocis-fonts-ttf
- name: ocis-fonts-map
  configMap:
    name: ocis-fonts-map
{{- end -}}

{{- define "ocis.caEnv" -}}
- name: IDM_ADMIN_PASSWORD
  value: admin
- name: PROXY_ENABLE_BASIC_AUTH
  value: "true"
- name: OCIS_LOG_FILE
  value: /logs/ocis.log
- name: AUDIT_LOG_TO_CONSOLE
  value: "false"
- name: STORAGE_PUBLICLINK_STORE_STORE
  value: {{ .Values.store.type | quote }}
- name: STORAGE_PUBLICLINK_STORE_NODES
  value: {{ tpl (join "," .Values.store.nodes) . | quote }}
- name: THUMBNAILS_TXT_FONTMAP_FILE
  value: /etc/ocis/fontsMap.json
- name: ANTIVIRUS_SCANNER_TYPE
  value: clamav
- name: ANTIVIRUS_CLAMAV_SOCKET
  value: "tcp://clamav:3310"
{{- end -}}
