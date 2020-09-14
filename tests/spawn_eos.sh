#!/bin/bash
# ME=$DRONE_HCLOUD_USER
# SERVER_NAME=droneci-eos-test-${DRONE_COMMIT_ID}-${RUN_PART}

# # setup ssh keys for hcloud
# ssh-keygen -b 2048 -t rsa -f /root/.ssh/id_rsa -q -N ""
# hcloud ssh-key create --name drone-${DRONE_COMMIT_ID}-${RUN_PART} --public-key-from-file /root/.ssh/id_rsa.pub

# # Create a new machine on hcloud for eos
# hcloud server create --type cx21 --image ubuntu-20.04 --ssh-key $SERVER_NAME --name $SERVER_NAME --label owner=$ME --label for=test --label from=eos-compose
# # time for the server to start up
# sleep 15

# IPADDR=$(hcloud server ip $SERVER_NAME)
# export IPADDR=$IPADDR
# export TEST_SERVER_URL=https://${IPADDR}:9200

# ssh -o StrictHostKeyChecking=no root@$IPADDR


# Setup system and clone ocis
ssh -tt root@$IPADDR apt-get update -y
ssh -tt root@$IPADDR apt-get install -y git screen docker.io docker-compose ldap-utils
ssh -tt root@$IPADDR git clone https://github.com/owncloud/ocis.git /ocis
ssh -tt root@$IPADDR "cd /ocis && git checkout $DRONE_COMMIT_ID"

# Create necessary files
ssh -tt root@$IPADDR "mkdir -p /ocis/tests/eos-config"
ssh -tt root@$IPADDR "cd /ocis/config && OCIS_DOMAIN=${IPADDR} bash /ocis/tests/config/drone/create-config.json.sh"
ssh -tt root@$IPADDR "cd /ocis/config && OCIS_DOMAIN=${IPADDR} bash /ocis/tests/config/drone/create-identifier-registration.sh"

# run ocis with eos
ssh -tt root@$IPADDR "cd /ocis && OCIS_DOMAIN=${IPADDR} docker-compose up -d"

# Some necessary configuration for eos
ssh -tt root@$IPADDR "cd /ocis && OCIS_DOMAIN=${IPADDR} bash /ocis/tests/config/drone/setup-eos.sh"
