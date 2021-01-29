#! /bin/bash
docker-compose exec keycloak \
    sh -c "cd /opt/jboss/keycloak && \
    timeout 60 bin/standalone.sh \
    -Djboss.httin/standalone.sh \
    -Djboss.socket.binding.port-offset=100 \
    -Dkeycloak.migration.action=export \
    -Dkeycloak.migration.provider=singleFile \
    -Dkeycloak.migration.realmName=oCIS \
    -Dkeycloak.migration.file=ocis-realm.json"
