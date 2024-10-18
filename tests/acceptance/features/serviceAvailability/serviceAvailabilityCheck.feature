Feature: service health check


  Scenario: health check defauts service
    When a user requests these endpoints with "GET"
      | endpoint                      | service            | comment |
      # | http://127.0.0.1:9197/healthz | activitylog        | #get 500                      |
      # | http://127.0.0.1:9297/healthz | antivirus          | checked in apiAntivirus suite |
      | http://127.0.0.1:9165/healthz | app-provider       |         |
      | http://127.0.0.1:9243/healthz | app-registry       |         |
      # | http://127.0.0.1:9229/healthz | audit              | extra                         |
      # | http://127.0.0.1:9245/healthz | auth-app           | extra                         |
      | http://127.0.0.1:9147/healthz | auth-basic         |         |
      # | http://127.0.0.1:9149/healthz | auth-bearer        | extra                         |
      | http://127.0.0.1:9167/healthz | auth-machine       |         |
      | http://127.0.0.1:9198/healthz | auth-service       |         |
      | http://127.0.0.1:9260/healthz | clientlog          |         |
      # | http://127.0.0.1:9980/healthz | collaboration      | checked in apiColaboration suite |
      | http://127.0.0.1:9270/healthz | eventhistory       |         |
      | http://127.0.0.1:9141/healthz | frontend           |         |
      | http://127.0.0.1:9143/healthz | gateway            |         |
      # | http://127.0.0.1:9124/healthz | graph              | #get 500                      |
      | http://127.0.0.1:9161/healthz | groups             |         |
      | http://127.0.0.1:9239/healthz | idm                |         |
      # | http://127.0.0.1:9134/healthz | idp                | #get 500                      |
      | http://127.0.0.1:9234/healthz | invitations        |         |
      | http://127.0.0.1:9234/healthz | nats               |         |
      # | http://127.0.0.1:9174/healthz | notifications      | checked in apiNotification suite |
      | http://127.0.0.1:9163/healthz | ocdav              |         |
      # | http://127.0.0.1:9281/healthz | ocm                | should not be used by default |
      # | http://127.0.0.1:9114/healthz | ocs                | #get 500                      |
      # | http://127.0.0.1:9129/healthz | policies           | extra                         |
      | http://127.0.0.1:9255/healthz | postprocessing     |         |
      # | http://127.0.0.1:9205/healthz | proxy              | #get 500                      |
      | http://127.0.0.1:9224/healthz | search             |         |
      # | http://127.0.0.1:9194/healthz | settings           | #get 500                      |
      | http://127.0.0.1:9151/healthz | sharing            |         |
      # | http://127.0.0.1:9135/healthz | sse                | #get 500                      |
      | http://127.0.0.1:9179/healthz | storage-publiclink |         |
      | http://127.0.0.1:9156/healthz | storage-shares     |         |
      | http://127.0.0.1:9217/healthz | storage-system     |         |
      | http://127.0.0.1:9159/healthz | storage-users      |         |
      # | http://127.0.0.1:/healthz     | store              | # deleted                     |
      # | http://127.0.0.1:9189/healthz | thumbnails         | #get 500                      |
      # | http://127.0.0.1:9210/healthz | userlog            | #get 500                      |
      | http://127.0.0.1:9145/healthz | users              |         |
    # | http://127.0.0.1:9104/healthz | web                | #get 500                      |
    # | http://127.0.0.1:9119/healthz | webdav             | #get 500                      |
    # | http://127.0.0.1:/healthz | webfinger          | #fix me need a proper port    |
    Then the HTTP status code of responses on all endpoints should be "200"

  @env-config
  Scenario: health check extra services
    Given the config "OCIS_ADD_RUN_SERVICES" has been set to "audit,auth-app,auth-bearer,policies"
    When a user requests these endpoints with "GET"
      | endpoint                      | service  | comment |
      | http://127.0.0.1:9229/healthz | audit    |         |
      | http://127.0.0.1:9245/healthz | auth-app |         |
      # | http://127.0.0.1:9149/healthz | auth-bearer | donot work |
      | http://127.0.0.1:9129/healthz | policies |         |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: service ready check
    When a user requests these endpoints with "GET"
      | endpoint                     | service            | comment  |
      # | http://127.0.0.1:9197/readyz | activitylog        | #get 500                      |
      # | http://127.0.0.1:9297/readyz | antivirus          | checked in apiAntivirus suite |
      | http://127.0.0.1:9165/readyz | app-provider       |          |
      | http://127.0.0.1:9243/readyz | app-registry       |          |
      # | http://127.0.0.1:9229/readyz | audit              | extra                         |
      # | http://127.0.0.1:9245/readyz | auth-app           | extra                         |
      | http://127.0.0.1:9147/readyz | auth-basic         |          |
      # | http://127.0.0.1:9149/readyz | auth-bearer        | extra                         |
      | http://127.0.0.1:9167/readyz | auth-machine       |          |
      | http://127.0.0.1:9198/readyz | auth-service       |          |
      # | http://127.0.0.1:9260/readyz | clientlog          | #get 500 |
      # | http://127.0.0.1:9980/readyz | collaboration      | checked in apiColaboration suite |
      # | http://127.0.0.1:9270/readyz | eventhistory       | #get 500 |
      | http://127.0.0.1:9141/readyz | frontend           |          |
      | http://127.0.0.1:9143/readyz | gateway            |          |
      # | http://127.0.0.1:9124/readyz | graph              | #get 500                      |
      | http://127.0.0.1:9161/readyz | groups             |          |
      | http://127.0.0.1:9239/readyz | idm                |          |
      # | http://127.0.0.1:9134/readyz | idp                | #get 500                      |
      | http://127.0.0.1:9234/readyz | invitations        |          |
      | http://127.0.0.1:9234/readyz | nats               |          |
      # | http://127.0.0.1:9174/readyz | notifications      | checked in apiNotification suite |
      | http://127.0.0.1:9163/readyz | ocdav              |          |
      # | http://127.0.0.1:9281/readyz | ocm                | should not be used by default |
      # | http://127.0.0.1:9114/readyz | ocs                | #get 500                      |
      # | http://127.0.0.1:9129/readyz | policies           | extra                         |
      | http://127.0.0.1:9255/readyz | postprocessing     |          |
      # | http://127.0.0.1:9205/readyz | proxy              | #get 500                      |
      # | http://127.0.0.1:9224/readyz | search             | #get 500 |
      # | http://127.0.0.1:9194/readyz | settings           | #get 500                      |
      | http://127.0.0.1:9151/readyz | sharing            |          |
      # | http://127.0.0.1:9135/readyz | sse                | #get 500                      |
      | http://127.0.0.1:9179/readyz | storage-publiclink |          |
      | http://127.0.0.1:9156/readyz | storage-shares     |          |
      | http://127.0.0.1:9217/readyz | storage-system     |          |
      | http://127.0.0.1:9159/readyz | storage-users      |          |
      # | http://127.0.0.1:/readyz     | store              | # deleted                     |
      # | http://127.0.0.1:9189/readyz | thumbnails         | #get 500                      |
      # | http://127.0.0.1:9210/readyz | userlog            | #get 500                      |
      | http://127.0.0.1:9145/readyz | users              |          |
    # | http://127.0.0.1:9104/readyz | web                | #get 500                      |
    # | http://127.0.0.1:9119/readyz | webdav             | #get 500                      |
    # | http://127.0.0.1:/readyz | webfinger          | #fix me need a proper port    |
    Then the HTTP status code of responses on all endpoints should be "200"

  @env-config
  Scenario: health check extra services
    Given the config "OCIS_ADD_RUN_SERVICES" has been set to "audit,auth-app,auth-bearer,policies"
    When a user requests these endpoints with "GET"
      | endpoint                     | service  | comment |
      # | http://127.0.0.1:9229/readyz | audit    | #get 500 |
      | http://127.0.0.1:9245/readyz | auth-app |         |
      # | http://127.0.0.1:9149/readyz | auth-bearer | donot work |
      | http://127.0.0.1:9129/readyz | policies |         |
    Then the HTTP status code of responses on all endpoints should be "200"
