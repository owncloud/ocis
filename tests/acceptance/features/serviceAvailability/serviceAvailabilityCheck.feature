Feature: service health check


  Scenario: health check defauts service
    When a user requests these endpoints:
      | endpoint                                        | service            | comment |
      # | %base_url_without_scheme_and_port%:9197/healthz | activitylog        | #get 500 |
      # | %base_url_without_scheme_and_port%:9297/healthz | antivirus          | checked in apiAntivirus suite |
      | %base_url_without_scheme_and_port%:9165/healthz | app-provider       |         |
      | %base_url_without_scheme_and_port%:9243/healthz | app-registry       |         |
      # | %base_url_without_scheme_and_port%:9229/healthz | audit              | extra   |
      # | %base_url_without_scheme_and_port%:9245/healthz | auth-app           | extra   |
      | %base_url_without_scheme_and_port%:9147/healthz | auth-basic         |         |
      # | %base_url_without_scheme_and_port%:9149/healthz | auth-bearer        | extra                         |
      | %base_url_without_scheme_and_port%:9167/healthz | auth-machine       |         |
      | %base_url_without_scheme_and_port%:9198/healthz | auth-service       |         |
      | %base_url_without_scheme_and_port%:9260/healthz | clientlog          |         |
      # | %base_url_without_scheme_and_port%:9980/healthz | collaboration      | checked in apiColaboration suite |
      | %base_url_without_scheme_and_port%:9270/healthz | eventhistory       |         |
      # | %base_url_without_scheme_and_port%:9141/healthz | frontend           | #get 500 |
      | %base_url_without_scheme_and_port%:9143/healthz | gateway            |         |
      # | %base_url_without_scheme_and_port%:9124/healthz | graph              | #get 500 |
      | %base_url_without_scheme_and_port%:9161/healthz | groups             |         |
      | %base_url_without_scheme_and_port%:9239/healthz | idm                |         |
      # | %base_url_without_scheme_and_port%:9134/healthz | idp                | #get 500 |
      # | %base_url_without_scheme_and_port%:/healthz | invitations        | #fix me need a proper port|
      | %base_url_without_scheme_and_port%:9234/healthz | nats               |         |
      # | %base_url_without_scheme_and_port%:9174/healthz | notifications      | checked in apiNotification suite |
      | %base_url_without_scheme_and_port%:9163/healthz | ocdav              |         |
      # | %base_url_without_scheme_and_port%:9281/healthz | ocm                | #get 500 should not be used by default |
      # | %base_url_without_scheme_and_port%:9114/healthz | ocs                | #get 500                      |
      # | %base_url_without_scheme_and_port%:9129/healthz | policies           | extra                         |
      | %base_url_without_scheme_and_port%:9255/healthz | postprocessing     |         |
      # | %base_url_without_scheme_and_port%:9205/healthz | proxy              | #get 500                      |
      | %base_url_without_scheme_and_port%:9224/healthz | search             |         |
      # | %base_url_without_scheme_and_port%:9194/healthz | settings           | #get 500                      |
      | %base_url_without_scheme_and_port%:9151/healthz | sharing            |         |
      # | %base_url_without_scheme_and_port%:9135/healthz | sse                | #get 500                      |
      | %base_url_without_scheme_and_port%:9179/healthz | storage-publiclink |         |
      | %base_url_without_scheme_and_port%:9156/healthz | storage-shares     |         |
      | %base_url_without_scheme_and_port%:9217/healthz | storage-system     |         |
      | %base_url_without_scheme_and_port%:9159/healthz | storage-users      |         |
      # | %base_url_without_scheme_and_port%:/healthz     | store              | # deleted                     |
      # | %base_url_without_scheme_and_port%:9189/healthz | thumbnails         | #get 500                      |
      # | %base_url_without_scheme_and_port%:9210/healthz | userlog            | #get 500                      |
      | %base_url_without_scheme_and_port%:9145/healthz | users              |         |
    # | %base_url_without_scheme_and_port%:9104/healthz | web                | #get 500                      |
    # | %base_url_without_scheme_and_port%:9119/healthz | webdav             | #get 500                      |
    # | %base_url_without_scheme_and_port%:/healthz | webfinger          | #fix me need a proper port    |
    Then the HTTP status code of responses on all endpoints should be "200"

  @env-config
  Scenario: health check extra services
    Given the following configs have been set:
      | config                 | value                               |
      | OCIS_ADD_RUN_SERVICES  | audit,auth-app,auth-bearer,policies |
      | AUDIT_DEBUG_ADDR       | 0.0.0.0:9229                        |
      | AUTH_APP_DEBUG_ADDR    | 0.0.0.0:9245                        |
      | AUTH_BEARER_DEBUG_ADDR | 0.0.0.0:9149                        |
      | POLICIES_DEBUG_ADDR    | 0.0.0.0:9129                        |
    When a user requests these endpoints:
      | endpoint                                        | service  | comment |
      | %base_url_without_scheme_and_port%:9229/healthz | audit    |         |
      | %base_url_without_scheme_and_port%:9245/healthz | auth-app |         |
      # | %base_url_without_scheme_and_port%:9149/healthz | auth-bearer | donot work |
      | %base_url_without_scheme_and_port%:9129/healthz | policies |         |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: service ready check
    When a user requests these endpoints:
      | endpoint                                       | service            | comment |
      # | %base_url_without_scheme_and_port%:9197/readyz | activitylog        | #get 500                      |
      # | %base_url_without_scheme_and_port%:9297/readyz | antivirus          | checked in apiAntivirus suite |
      | %base_url_without_scheme_and_port%:9165/readyz | app-provider       |         |
      | %base_url_without_scheme_and_port%:9243/readyz | app-registry       |         |
      # | %base_url_without_scheme_and_port%:9229/readyz | audit              | extra                         |
      # | %base_url_without_scheme_and_port%:9245/readyz | auth-app           | extra                         |
      | %base_url_without_scheme_and_port%:9147/readyz | auth-basic         |         |
      # | %base_url_without_scheme_and_port%:9149/readyz | auth-bearer        | extra                         |
      | %base_url_without_scheme_and_port%:9167/readyz | auth-machine       |         |
      | %base_url_without_scheme_and_port%:9198/readyz | auth-service       |         |
      # | %base_url_without_scheme_and_port%:9260/readyz | clientlog          | #get 500 |
      # | %base_url_without_scheme_and_port%:9980/readyz | collaboration      | checked in apiColaboration suite |
      # | %base_url_without_scheme_and_port%:9270/readyz | eventhistory       | #get 500 |
      # | %base_url_without_scheme_and_port%:9141/readyz | frontend           | #get 500 |
      | %base_url_without_scheme_and_port%:9143/readyz | gateway            |         |
      # | %base_url_without_scheme_and_port%:9124/readyz | graph              | #get 500                      |
      | %base_url_without_scheme_and_port%:9161/readyz | groups             |         |
      | %base_url_without_scheme_and_port%:9239/readyz | idm                |         |
      # | %base_url_without_scheme_and_port%:9134/readyz | idp                | #get 500                      |
      # | %base_url_without_scheme_and_port%:/readyz | invitations        |#fix me need a proper port|
      | %base_url_without_scheme_and_port%:9234/readyz | nats               |         |
      # | %base_url_without_scheme_and_port%:9174/readyz | notifications      | checked in apiNotification suite |
      | %base_url_without_scheme_and_port%:9163/readyz | ocdav              |         |
      # | %base_url_without_scheme_and_port%:9281/readyz | ocm                | should not be used by default |
      # | %base_url_without_scheme_and_port%:9114/readyz | ocs                | #get 500                      |
      # | %base_url_without_scheme_and_port%:9129/readyz | policies           | extra                         |
      | %base_url_without_scheme_and_port%:9255/readyz | postprocessing     |         |
      # | %base_url_without_scheme_and_port%:9205/readyz | proxy              | #get 500                      |
      # | %base_url_without_scheme_and_port%:9224/readyz | search             | #get 500 |
      # | %base_url_without_scheme_and_port%:9194/readyz | settings           | #get 500                      |
      | %base_url_without_scheme_and_port%:9151/readyz | sharing            |         |
      # | %base_url_without_scheme_and_port%:9135/readyz | sse                | #get 500                      |
      | %base_url_without_scheme_and_port%:9179/readyz | storage-publiclink |         |
      | %base_url_without_scheme_and_port%:9156/readyz | storage-shares     |         |
      | %base_url_without_scheme_and_port%:9217/readyz | storage-system     |         |
      | %base_url_without_scheme_and_port%:9159/readyz | storage-users      |         |
      # | %base_url_without_scheme_and_port%:/readyz     | store              | # deleted                     |
      # | %base_url_without_scheme_and_port%:9189/readyz | thumbnails         | #get 500                      |
      # | %base_url_without_scheme_and_port%:9210/readyz | userlog            | #get 500                      |
      | %base_url_without_scheme_and_port%:9145/readyz | users              |         |
    # | %base_url_without_scheme_and_port%:9104/readyz | web                | #get 500                      |
    # | %base_url_without_scheme_and_port%:9119/readyz | webdav             | #get 500                      |
    # | %base_url_without_scheme_and_port%:/readyz | webfinger          | #fix me need a proper port    |
    Then the HTTP status code of responses on all endpoints should be "200"

  @env-config
  Scenario: health check extra services
    Given the following configs have been set:
      | config                 | value                               |
      | OCIS_ADD_RUN_SERVICES  | audit,auth-app,auth-bearer,policies |
      | AUDIT_DEBUG_ADDR       | 0.0.0.0:9229                        |
      | AUTH_APP_DEBUG_ADDR    | 0.0.0.0:9245                        |
      | AUTH_BEARER_DEBUG_ADDR | 0.0.0.0:9149                        |
      | POLICIES_DEBUG_ADDR    | 0.0.0.0:9129                        |
    When a user requests these endpoints:
      | endpoint                                       | service  | comment |
      | %base_url_without_scheme_and_port%:9229/readyz | audit    |         |
      | %base_url_without_scheme_and_port%:9245/readyz | auth-app |         |
      # | %base_url_without_scheme_and_port%:9149/readyz | auth-bearer | donot work |
      | %base_url_without_scheme_and_port%:9129/readyz | policies |         |
    Then the HTTP status code of responses on all endpoints should be "200"
