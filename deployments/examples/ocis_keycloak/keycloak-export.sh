#! /bin/bash
docker-compose exec keycloak \
    sh -c "cd /opt/jboss/keycloak && \
    timeout 60 bin/standalone.sh \
    -Djboss.socket.binding.port-offset=100 \
    -Dkeycloak.migration.action=export \
    -Dkeycloak.migration.provider=singleFile \
    -Dkeycloak.migration.file=keycloak-export.json \
    -Djboss.httin/standalone.sh -Dkeycloak.migration.action=export \
    -Dkeycloak.migration.provider=singleFile \
    -Dkeycloak.migration.file=keycloak-export.json"

docker-compose exec keycloak cat /opt/jboss/keycloak/keycloak-export.json > keycloak-export.json
