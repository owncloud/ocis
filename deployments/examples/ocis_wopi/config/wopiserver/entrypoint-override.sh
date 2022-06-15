#!/bin/bash
set -e

echo "${WOPISECRET}" > /etc/wopi/wopisecret
echo "${IOPSECRET}" > /etc/wopi/iopsecret

cp /etc/wopi/wopiserver.conf.dist /etc/wopi/wopiserver.conf
sed -i 's/wopiserver.owncloud.test/'${WOPISERVER_DOMAIN}'/g' /etc/wopi/wopiserver.conf


if [ "$WOPISERVER_INSECURE" == "true" ]; then
    sed -i 's/sslverify\s=\sTrue/sslverify = False/g' /etc/wopi/wopiserver.conf
fi

touch /var/log/wopi/wopiserver.log

/app/wopiserver.py &

tail -f /var/log/wopi/wopiserver.log
