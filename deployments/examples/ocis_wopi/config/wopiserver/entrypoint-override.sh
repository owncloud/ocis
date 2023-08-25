#!/bin/sh
set -e

echo "${WOPISECRET}" > /etc/wopi/wopisecret

cp /etc/wopi/wopiserver.conf.dist /etc/wopi/wopiserver.conf
sed -i 's/wopiserver.owncloud.test/'${WOPISERVER_DOMAIN}'/g' /etc/wopi/wopiserver.conf

if [ "$WOPISERVER_INSECURE" = "true" ]; then
    sed -i 's/sslverify\s=\sTrue/sslverify = False/g' /etc/wopi/wopiserver.conf
fi

/app/wopiserver.py
