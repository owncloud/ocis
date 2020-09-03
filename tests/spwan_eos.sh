#!/bin/bash
ME=$DRONE_HCLOUD_USER
SERVER_NAME=droneci-eos-test-${DRONE_COMMIT_ID}

# Create a new machine on hcloud for eos
hcloud server create --type cx21 --image ubuntu-20.04 --ssh-key $ME --name $SERVER_NAME --label owner=$ME --label for=test --label from=eos-compose

IPADDR=$(hcloud server ip $SERVER_NAME)
OCIS_DOMAIN=$(hcloud server ip $SERVER_NAME)

# timeout 180 while [[ \"$(curl -k -v -s -o /dev/null -w ''%{http_code}'' https://:9200)\" != \"200\" ]]; do sleep 2; done
# sleep 15

# Setup system and clone ocis
ssh -t root@$IPADDR apt-get update -y
ssh -t root@$IPADDR apt-get install -y git screen docker.io docker-compose ldap-utils
ssh -t root@$IPADDR git clone https://github.com/owncloud/ocis.git /ocis
ssh -t root@$IPADDR "cd /ocis && git checkout $DRONE_COMMIT_ID"

# Create necessary files
ssh -t root@$IPADDR "cd /ocis/tests/config/drone && OCIS_DOMAIN=${IPADDR} bash /ocis/tests/config/drone/create-config.json.sh"
ssh -t root@$IPADDR "cd /ocis/tests/config/drone && OCIS_DOMAIN=${IPADDR} bash /ocis/tests/config/drone/create-identifier-registration.sh"

# run ocis with eos
ssh -t root@$IPADDR "cd /ocis && OCIS_DOMAIN=${IPADDR} docker-compose -f ./docker-compose-eos-ci.yml up -d"

# Some necessary configuration for eos
ssh -t root@$IPADDR "cd /ocis && bash /ocis/tests/config/drone/setup-eos.sh"
