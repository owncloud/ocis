<?php

# reference: https://doc.owncloud.com/server/admin_manual/configuration/user/oidc/

function getOIDCConfigFromEnv()
{
    $config = [
        'openid-connect' => [
            'provider-url' => getenv('IDP_OIDC_ISSUER'),
            'client-id' => 'oc10',
            'client-secret' => getenv('IDP_OIDC_CLIENT_SECRET'),
            'loginButtonName' => 'OpenId Connect',
            'search-attribute' => 'preferred_username',
            'mode' => 'userid',
            'autoRedirectOnLoginPage' => true,
            'insecure' => true,
            'post_logout_redirect_uri' => 'https://' . getenv('CLOUD_DOMAIN'),
        ],
    ];
    return $config;
}

$CONFIG = getOIDCConfigFromEnv();
