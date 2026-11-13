Enhancement: Update Traefik in the LDAP, Keycloak, Hello and Multi examples

The `ocis_ldap`, `ocis_keycloak` and `ocis_hello` examples pinned
`traefik:v2.9.1`. Its Docker provider can not communicate with Docker Engine 29
and later, so Traefik loads no configuration at all on such hosts and answers
every request with `404`. `ocis_multi` was on `traefik:v3.6.7`.

All of them now run `traefik:v3.6.10`, which is compatible with Docker Engine 29
and later.
