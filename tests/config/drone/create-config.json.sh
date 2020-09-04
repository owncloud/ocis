#!/bin/bash
echo "{
        \"server\":\"https://${OCIS_DOMAIN}:9200\",
        \"theme\":\"owncloud\",
        \"version\":\"0.1.0\",
        \"openIdConnect\":{
                \"metadata_url\":\"https://${OCIS_DOMAIN}:9200/.well-known/openid-configuration\",
                \"authority\":\"https://${OCIS_DOMAIN}:9200\",
                \"client_id\":\"phoenix\",
                \"response_type\":\"code\",
                \"scope\":\"openid profile email\"
        },
        \"apps\":[\"files\",\"draw-io\",\"pdf-viewer\",\"markdown-editor\",\"media-viewer\"],
        \"external_apps\":[
                {\"id\":\"accounts\",\"path\":\"https://${OCIS_DOMAIN}:9200/accounts.js\"},
                {\"id\":\"settings\",\"path\":\"https://${OCIS_DOMAIN}:9200/settings.js\"}
        ],
        \"options\":{\"hideSearchBar\":true}
}" > $PWD/config.json

