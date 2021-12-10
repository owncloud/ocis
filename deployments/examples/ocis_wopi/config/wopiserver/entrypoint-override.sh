#/bin/sh!
set -e

echo "${WOPISECRET}" > /etc/wopi/wopisecret
echo "${IOPSECRET}" > /etc/wopi/iopsecret
mkdir -p /var/run/secrets
echo "$CODIMDSECRET" > /var/run/secrets/codimd_apikey

cp /etc/wopi/wopiserver.conf.dist /etc/wopi/wopiserver.conf
sed -i 's/ocis.owncloud.test/'${OCIS_DOMAIN}'/g' /etc/wopi/wopiserver.conf
sed -i 's/collabora.owncloud.test/'${COLLABORA_DOMAIN}'/g' /etc/wopi/wopiserver.conf
sed -i 's/wopiserver.owncloud.test/'${WOPISERVER_DOMAIN}'/g' /etc/wopi/wopiserver.conf

touch /var/log/wopi/wopiserver.log

/app/wopiserver.py &

tail -f /var/log/wopi/wopiserver.log
