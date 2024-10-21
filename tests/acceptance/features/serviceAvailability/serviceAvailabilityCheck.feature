Feature: service health check


  Scenario: health check defauts service
    When a user requests these endpoints with "GET"
      | endpoint                                        | service            | comment |
      # | http://localhost/healthz      | activitylog        | #get 500                      |
      # | http://localhost:9297/healthz | antivirus          | checked in apiAntivirus suite |
      | http://localhost:9165/healthz | app-provider       |         |
      | http://localhost:9243/healthz | app-registry       |         |
      # | http://localhost:9229/healthz | audit              | extra   |
      # | http://localhost:9245/healthz | auth-app           | extra   |
      | http://localhost:9147/healthz | auth-basic         |         |
      # | http://localhost:9149/healthz | auth-bearer        | extra                         |
      | http://localhost:9167/healthz | auth-machine       |         |
      | http://localhost:9198/healthz | auth-service       |         |
      | http://localhost:9260/healthz | clientlog          |         |
      # | http://localhost:9980/healthz | collaboration      | checked in apiColaboration suite |
      | http://localhost:9270/healthz | eventhistory       |         |
      | http://localhost:9141/healthz | frontend           |         |
      | http://localhost:9143/healthz | gateway            |         |
      # | http://localhost:9124/healthz | graph              | #get 500                      |
      | http://localhost:9161/healthz | groups             |         |
      | http://localhost:9239/healthz | idm                |         |
      # | http://localhost:9134/healthz | idp                | #get 500                      |
      | http://localhost:9234/healthz | invitations        |         |
      | http://localhost:9234/healthz | nats               |         |
      # | http://localhost:9174/healthz | notifications      | checked in apiNotification suite |
      | http://localhost:9163/healthz | ocdav              |         |
      # | http://localhost:9281/healthz | ocm                | should not be used by default |
      # | http://localhost:9114/healthz | ocs                | #get 500                      |
      # | http://localhost:9129/healthz | policies           | extra                         |
      | http://localhost:9255/healthz | postprocessing     |         |
      # | http://localhost:9205/healthz | proxy              | #get 500                      |
      | http://localhost:9224/healthz | search             |         |
      # | http://localhost:9194/healthz | settings           | #get 500                      |
      | http://localhost:9151/healthz | sharing            |         |
      # | http://localhost:9135/healthz | sse                | #get 500                      |
      | http://localhost:9179/healthz | storage-publiclink |         |
      | http://localhost:9156/healthz | storage-shares     |         |
      | http://localhost:9217/healthz | storage-system     |         |
      | http://localhost:9159/healthz | storage-users      |         |
      # | http://localhost:/healthz     | store              | # deleted                     |
      # | http://localhost:9189/healthz | thumbnails         | #get 500                      |
      # | http://localhost:9210/healthz | userlog            | #get 500                      |
      | http://localhost:9145/healthz | users              |         |
    # | http://localhost:9104/healthz | web                | #get 500                      |
    # | http://localhost:9119/healthz | webdav             | #get 500                      |
    # | http://localhost:/healthz | webfinger          | #fix me need a proper port    |

  @env-config
  Scenario: health check extra services
    Given the config "OCIS_ADD_RUN_SERVICES" has been set to "audit,auth-app,auth-bearer,policies"
    When a user requests these endpoints with "GET"
      | endpoint                      | service  | comment |
      | http://localhost:9229/healthz | audit    |         |
      | http://localhost:9245/healthz | auth-app |         |
      # | http://localhost:9149/healthz | auth-bearer | donot work |
      | http://localhost:9129/healthz | policies |         |


  Scenario: service ready check
    When a user requests these endpoints with "GET"
      | endpoint                                       | service            | comment |
      # | http://localhost:9197/readyz | activitylog        | #get 500                      |
      # | http://localhost:9297/readyz | antivirus          | checked in apiAntivirus suite |
      | http://localhost:9165/readyz | app-provider       |         |
      | http://localhost:9243/readyz | app-registry       |         |
      # | http://localhost:9229/readyz | audit              | extra                         |
      # | http://localhost:9245/readyz | auth-app           | extra                         |
      | http://localhost:9147/readyz | auth-basic         |         |
      # | http://localhost:9149/readyz | auth-bearer        | extra                         |
      | http://localhost:9167/readyz | auth-machine       |         |
      | http://localhost:9198/readyz | auth-service       |         |
      # | http://localhost:9260/readyz | clientlog          | #get 500 |
      # | http://localhost:9980/readyz | collaboration      | checked in apiColaboration suite |
      # | http://localhost:9270/readyz | eventhistory       | #get 500 |
      | http://localhost:9141/readyz | frontend           |         |
      | http://localhost:9143/readyz | gateway            |         |
      # | http://localhost:9124/readyz | graph              | #get 500                      |
      | http://localhost:9161/readyz | groups             |         |
      | http://localhost:9239/readyz | idm                |         |
      # | http://localhost:9134/readyz | idp                | #get 500                      |
      | http://localhost:9234/readyz | invitations        |         |
      | http://localhost:9234/readyz | nats               |         |
      # | http://localhost:9174/readyz | notifications      | checked in apiNotification suite |
      | http://localhost:9163/readyz | ocdav              |         |
      # | http://localhost:9281/readyz | ocm                | should not be used by default |
      # | http://localhost:9114/readyz | ocs                | #get 500                      |
      # | http://localhost:9129/readyz | policies           | extra                         |
      | http://localhost:9255/readyz | postprocessing     |         |
      # | http://localhost:9205/readyz | proxy              | #get 500                      |
      # | http://localhost:9224/readyz | search             | #get 500 |
      # | http://localhost:9194/readyz | settings           | #get 500                      |
      | http://localhost:9151/readyz | sharing            |         |
      # | http://localhost:9135/readyz | sse                | #get 500                      |
      | http://localhost:9179/readyz | storage-publiclink |         |
      | http://localhost:9156/readyz | storage-shares     |         |
      | http://localhost:9217/readyz | storage-system     |         |
      | http://localhost:9159/readyz | storage-users      |         |
      # | http://localhost:/readyz     | store              | # deleted                     |
      # | http://localhost:9189/readyz | thumbnails         | #get 500                      |
      # | http://localhost:9210/readyz | userlog            | #get 500                      |
      | http://localhost:9145/readyz | users              |         |
    # | http://localhost:9104/readyz | web                | #get 500                      |
    # | http://localhost:9119/readyz | webdav             | #get 500                      |
    # | http://localhost:/readyz | webfinger          | #fix me need a proper port    |

  @env-config
  Scenario: health check extra services
    Given the config "OCIS_ADD_RUN_SERVICES" has been set to "audit,auth-app,auth-bearer,policies"
    When a user requests these endpoints with "GET"
      | endpoint                                       | service  | comment |
      # | http://localhost:9229/readyz | audit    | #get 500 |
      | http://localhost:9245/readyz | auth-app |         |
      # | http://localhost:9149/readyz | auth-bearer | donot work |
      | http://localhost:9129/readyz | policies |         |
