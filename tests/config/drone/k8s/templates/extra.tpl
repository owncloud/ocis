{{- define "ocis.extraVolMounts" -}}
- name: logdir
  mountPath: /logs
{{- if eq .appName "thumbnails" }}
- name: ocis-fonts-ttf
  mountPath: /etc/ocis/fonts
- name: ocis-fonts-map
  mountPath: /etc/ocis/fontsMap.json
  subPath: fontsMap.json
{{- end -}}
{{- end -}}

{{- define "ocis.extraVolumes" -}}
- name: logdir
  hostPath:
    path: /logs
    type: Directory
{{- if eq .appName "thumbnails" }}
- name: ocis-fonts-ttf
  configMap:
    name: ocis-fonts-ttf
- name: ocis-fonts-map
  configMap:
    name: ocis-fonts-map
{{- end -}}
{{- end -}}

{{- define "ocis.extraEnvs" -}}
- name: OCIS_LOG_FILE
  value: /logs/ocis.log
{{- if eq .appName "idm" }}
- name: IDM_ADMIN_PASSWORD
  value: admin
{{- end -}}
{{- if eq .appName "proxy" }}
- name: PROXY_ENABLE_BASIC_AUTH
  value: "true"
{{- end -}}
{{- if eq .appName "audit" }}
- name: AUDIT_LOG_TO_CONSOLE
  value: "false"
{{- end -}}
{{- if eq .appName "thumbnails" }}
- name: THUMBNAILS_TXT_FONTMAP_FILE
  value: /etc/ocis/fontsMap.json
{{- end -}}
{{- if eq .appName "antivirus" }}
- name: ANTIVIRUS_SCANNER_TYPE
  value: clamav
- name: ANTIVIRUS_CLAMAV_SOCKET
  value: "tcp://clamav:3310"
{{- end -}}
{{- end -}}
