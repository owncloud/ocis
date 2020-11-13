# ocis with keycloak as identity provider

## set up DNS / hostnames

If you only want to start the example on your local computer, please add following lines to your hosts file (`/etc/hosts):

```
127.0.0.1 traefik.owncloud.test
127.0.0.1 ocis.owncloud.test
127.0.0.1 keycloak.owncloud.test
```

For a deployment of the example on a internet facing server please set up dns entries pointing to your server.

## configure the deployment example

If you are deploying the example on your local computer, you don't have to change anything in the `.env` file. Please continue with starting the stack.

For a internet facing server you must change at least the domains according to your dns setup:
- `TRAEFIK_DOMAIN`
- `OCIS_DOMAIN`
- `KEYCLOAK_DOMAIN`

In order to receive ssl certificates with letsencrypt change `TRAEFIK_ACME_MAIL` to a valid email address

For security reasons please change credentials
- `TRAEFIK_BASIC_AUTH_USERS`
- `KEYCLOAK_ADMIN_USER` and `KEYCLOAK_ADMIN_PASSWORD`

Please also note that there might be default users in ocis, depending on the version you are using.

## starting the stack

just run `docker-compose up`

## set up keycloak

1. go to http://keycloak.owncloud.test (https://KEYCLOAK_DOMAIN) and log in as `admin` (password is `admin`).
2. go to clients settings and add a client. The client id is `ocis-phoenix` ($OCIS_OIDC_CLIENT_ID). The client protocol is openid-connect. Insert `https://ocis.owncloud.test` (https://OCIS_DOMAIN) as root url. Then save the client.
3. you can now add users in the users section.

## test ocis
you now can login with your users configured in keycloak
