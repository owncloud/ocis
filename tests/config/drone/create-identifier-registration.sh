#!/bin/bash
echo "---
# OpenID Connect client registry.
clients:
  - id: phoenix
    name: OCIS
    application_type: web
    insecure: yes
    trusted: yes
    redirect_uris:
      - https://${OCIS_DOMAIN}:9200/oidc-callback.html
      - https://${OCIS_DOMAIN}:9200/
    origins:
      -  https://${OCIS_DOMAIN}:9200
authorities:" > $PWD/identifier-registration.yml
