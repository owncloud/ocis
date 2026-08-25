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
{{- if .Values.features.vault.enabled }}
{{- if eq .appName "proxy" }}
- name: OCIS_ENABLE_VAULT_MODE
  value: "true"
- name: OCIS_MFA_ENABLED
  value: "true"
- name: PROXY_OIDC_ISSUER
  value: "https://keycloak:8443/realms/oCIS"
- name: PROXY_OIDC_REWRITE_WELLKNOWN
  value: "true"
- name: PROXY_AUTOPROVISION_ACCOUNTS
  value: "true"
- name: PROXY_ROLE_ASSIGNMENT_DRIVER
  value: oidc
- name: PROXY_USER_OIDC_CLAIM
  value: preferred_username
- name: PROXY_USER_CS3_CLAIM
  value: username
{{- end -}}
{{- if eq .appName "frontend" }}
- name: OCIS_ENABLE_VAULT_MODE
  value: "true"
- name: OCIS_MFA_ENABLED
  value: "true"
{{- end -}}
{{- if eq .appName "graph" }}
- name: OCIS_ENABLE_VAULT_MODE
  value: "true"
- name: GRAPH_ASSIGN_DEFAULT_USER_ROLE
  value: "false"
- name: GRAPH_USERNAME_MATCH
  value: none
{{- end -}}
{{- if eq .appName "gateway" }}
- name: OCIS_ENABLE_VAULT_MODE
  value: "true"
{{- end -}}
{{- if eq .appName "web" }}
{{/* WEB_OIDC_AUTHORITY intentionally left as the chart default (the oCIS
     domain): PROXY_OIDC_REWRITE_WELLKNOWN transparently proxies the
     well-known document from Keycloak under that same origin, matching
     how the non-k8s vault test setup (run-github.py) configures this. */}}
- name: WEB_OIDC_CLIENT_ID
  value: web
- name: WEB_OIDC_SCOPE
  value: "openid profile email acr"
{{- end -}}
{{- if eq .appName "webfinger" }}
- name: WEBFINGER_OIDC_ISSUER
  value: "https://keycloak:8443/realms/oCIS"
{{- end -}}
{{- if eq .appName "users" }}
- name: USERS_IDP_URL
  value: "https://keycloak:8443/realms/oCIS"
{{- end -}}
{{- if eq .appName "groups" }}
- name: GROUPS_IDP_URL
  value: "https://keycloak:8443/realms/oCIS"
{{- end -}}
{{- if eq .appName "ocs" }}
- name: OCS_IDM_ADDRESS
  value: "https://keycloak:8443/realms/oCIS"
{{- end -}}
{{- if eq .appName "idm" }}
- name: OCIS_OIDC_ISSUER
  value: "https://keycloak:8443/realms/oCIS"
- name: IDM_CREATE_DEMO_USERS
  value: "false"
{{- end -}}
{{- if eq .appName "storageusers-vault" }}
- name: STORAGE_USERS_ENABLE_VAULT_MODE
  value: "true"
- name: STORAGE_USERS_SERVICE_NAME
  value: storage-users-vault
- name: STORAGE_USERS_EVENTS_CONSUMER_GROUP
  value: vault-dcfs
{{- end -}}
{{- end -}}
{{- end -}}
