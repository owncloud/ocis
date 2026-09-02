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
{{/* Without these overrides this instance falls back to the shared OCIS_CACHE_STORE=nats-js-kv
     used by the regular storage-users service. Since personal-space IDs are just the user's
     opaque ID (identical in both instances), the vault instance's existence check on
     CreateStorageSpace sees a cache hit from whatever the regular instance already created for
     that user and returns AlreadyExists without ever writing its own space - the vault personal
     space then never actually exists on this instance, even though every check claims it does. */}}
- name: STORAGE_USERS_FILEMETADATA_CACHE_STORE
  value: memory
- name: STORAGE_USERS_ID_CACHE_STORE
  value: memory
{{/* services/storage-users/pkg/config/parser/parse.go calls EnsureDefaults() (which forces
     MountID to the vault constant when EnableVaultMode is set) BEFORE envdecode.Decode()
     applies env vars - so EnableVaultMode is still false at that point and the override never
     fires. Without this, MountID falls back to whatever STORAGE_USERS_MOUNT_ID resolves to
     (the "storage-uuid" ConfigMap, shared with the regular storage-users instance), so spaces
     created here come back with the wrong storage id embedded in their space id. That id still
     lists fine (gateway's registry matches on its own static rule, not this value), but any
     later ID-based lookup (e.g. creating a folder inside a newly created project space) fails
     to route back to this instance. Setting it directly here sidesteps the ordering bug
     entirely, since env vars are applied regardless of EnsureDefaults. */}}
- name: STORAGE_USERS_MOUNT_ID
  value: "1a01c2c4-4309-4483-a845-842fd56d8622"
{{- end -}}
{{- end -}}
{{- end -}}
