#! /bin/bash
docker-compose exec keycloak \
    sh -c "cd /opt/jboss/keycloak && \
    timeout 60 bin/standalone.sh \
    -Djboss.httin/standalone.sh \
    -Djboss.socket.binding.port-offset=100 \
    -Dkeycloak.migration.action=export \
    -Dkeycloak.migration.provider=singleFile \
    -Dkeycloak.migration.realmName=owncloud \
    -Dkeycloak.migration.file=owncloud-realm.json"

docker-compose exec keycloak \
    cp /opt/jboss/keycloak/owncloud-realm.json /opt/jboss/keycloak/owncloud-realm.dist.json
