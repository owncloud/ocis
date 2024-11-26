Feature: service health check


  Scenario: check default services health
    When a user requests these endpoints with "GET" and no authentication
      | endpoint                         | service            |
      | %base_url_hostname%:9197/healthz | activitylog        |
      | %base_url_hostname%:9165/healthz | app-provider       |
      | %base_url_hostname%:9243/healthz | app-registry       |
      | %base_url_hostname%:9147/healthz | auth-basic         |
      | %base_url_hostname%:9167/healthz | auth-machine       |
      | %base_url_hostname%:9198/healthz | auth-service       |
      | %base_url_hostname%:9260/healthz | clientlog          |
      | %base_url_hostname%:9270/healthz | eventhistory       |
      | %base_url_hostname%:9141/healthz | frontend           |
      | %base_url_hostname%:9143/healthz | gateway            |
      | %base_url_hostname%:9124/healthz | graph              |
      | %base_url_hostname%:9161/healthz | groups             |
      | %base_url_hostname%:9239/healthz | idm                |
      | %base_url_hostname%:9134/healthz | idp                |
      | %base_url_hostname%:9234/healthz | nats               |
      | %base_url_hostname%:9163/healthz | ocdav              |
      | %base_url_hostname%:9281/healthz | ocm                |
      | %base_url_hostname%:9114/healthz | ocs                |
      | %base_url_hostname%:9255/healthz | postprocessing     |
      | %base_url_hostname%:9205/healthz | proxy              |
      | %base_url_hostname%:9224/healthz | search             |
      | %base_url_hostname%:9194/healthz | settings           |
      | %base_url_hostname%:9151/healthz | sharing            |
      | %base_url_hostname%:9139/healthz | sse                |
      | %base_url_hostname%:9179/healthz | storage-publiclink |
      | %base_url_hostname%:9156/healthz | storage-shares     |
      | %base_url_hostname%:9217/healthz | storage-system     |
      | %base_url_hostname%:9159/healthz | storage-users      |
      | %base_url_hostname%:9189/healthz | thumbnails         |
      | %base_url_hostname%:9214/healthz | userlog            |
      | %base_url_hostname%:9145/healthz | users              |
      | %base_url_hostname%:9104/healthz | web                |
      | %base_url_hostname%:9119/healthz | webdav             |
      | %base_url_hostname%:9279/healthz | webfinger          |
    Then the HTTP status code of responses on all endpoints should be "200"

  @env-config
  Scenario: check extra services health
    Given the following configs have been set:
      | config                 | value                                           |
      | OCIS_ADD_RUN_SERVICES  | audit,auth-app,auth-bearer,policies,invitations |
      | AUDIT_DEBUG_ADDR       | 0.0.0.0:9229                                    |
      | AUTH_APP_DEBUG_ADDR    | 0.0.0.0:9245                                    |
      | AUTH_BEARER_DEBUG_ADDR | 0.0.0.0:9149                                    |
      | POLICIES_DEBUG_ADDR    | 0.0.0.0:9129                                    |
      | INVITATIONS_DEBUG_ADDR | 0.0.0.0:9269                                    |
    When a user requests these endpoints with "GET" and no authentication
      | endpoint                         | service     |
      | %base_url_hostname%:9229/healthz | audit       |
      | %base_url_hostname%:9245/healthz | auth-app    |
      | %base_url_hostname%:9149/healthz | auth-bearer |
      | %base_url_hostname%:9269/healthz | invitations |
      | %base_url_hostname%:9129/healthz | policies    |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: check default services readiness
    When a user requests these endpoints with "GET" and no authentication
      | endpoint                        | service            |
      | %base_url_hostname%:9197/readyz | activitylog        |
      | %base_url_hostname%:9165/readyz | app-provider       |
      | %base_url_hostname%:9243/readyz | app-registry       |
      | %base_url_hostname%:9147/readyz | auth-basic         |
      | %base_url_hostname%:9167/readyz | auth-machine       |
      | %base_url_hostname%:9198/readyz | auth-service       |
      | %base_url_hostname%:9260/readyz | clientlog          |
      | %base_url_hostname%:9270/readyz | eventhistory       |
      | %base_url_hostname%:9141/readyz | frontend           |
      | %base_url_hostname%:9143/readyz | gateway            |
      | %base_url_hostname%:9124/readyz | graph              |
      | %base_url_hostname%:9161/readyz | groups             |
      | %base_url_hostname%:9239/readyz | idm                |
      | %base_url_hostname%:9134/readyz | idp                |
      | %base_url_hostname%:9234/readyz | nats               |
      | %base_url_hostname%:9163/readyz | ocdav              |
      | %base_url_hostname%:9281/readyz | ocm                |
      | %base_url_hostname%:9114/readyz | ocs                |
      | %base_url_hostname%:9255/readyz | postprocessing     |
      | %base_url_hostname%:9205/readyz | proxy              |
      | %base_url_hostname%:9224/readyz | search             |
      | %base_url_hostname%:9194/readyz | settings           |
      | %base_url_hostname%:9151/readyz | sharing            |
      | %base_url_hostname%:9139/readyz | sse                |
      | %base_url_hostname%:9179/readyz | storage-publiclink |
      | %base_url_hostname%:9156/readyz | storage-shares     |
      | %base_url_hostname%:9217/readyz | storage-system     |
      | %base_url_hostname%:9159/readyz | storage-users      |
      | %base_url_hostname%:9189/readyz | thumbnails         |
      | %base_url_hostname%:9214/readyz | userlog            |
      | %base_url_hostname%:9145/readyz | users              |
      | %base_url_hostname%:9104/readyz | web                |
      | %base_url_hostname%:9119/readyz | webdav             |
      | %base_url_hostname%:9279/readyz | webfinger          |
    Then the HTTP status code of responses on all endpoints should be "200"

  @env-config
  Scenario: check extra services readiness
    Given the following configs have been set:
      | config                 | value                                           |
      | OCIS_ADD_RUN_SERVICES  | audit,auth-app,auth-bearer,policies,invitations |
      | AUDIT_DEBUG_ADDR       | 0.0.0.0:9229                                    |
      | AUTH_APP_DEBUG_ADDR    | 0.0.0.0:9245                                    |
      | AUTH_BEARER_DEBUG_ADDR | 0.0.0.0:9149                                    |
      | POLICIES_DEBUG_ADDR    | 0.0.0.0:9129                                    |
      | INVITATIONS_DEBUG_ADDR | 0.0.0.0:9269                                    |
    When a user requests these endpoints with "GET" and no authentication
      | endpoint                        | service     |
      | %base_url_hostname%:9229/readyz | audit       |
      | %base_url_hostname%:9245/readyz | auth-app    |
      | %base_url_hostname%:9149/readyz | auth-bearer |
      | %base_url_hostname%:9269/readyz | invitations |
      | %base_url_hostname%:9129/readyz | policies    |
    Then the HTTP status code of responses on all endpoints should be "200"
