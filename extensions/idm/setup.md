---
title: Service Setup
date: 2022-03-22T00:00:00+00:00
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/extensions/idm
geekdocFilePath: setup.md
geekdocCollapseSection: true
---

{{< toc >}}

## Using ocis with libregraph/idm

Currently, oCIS still runs the accounts and glauth services to manage users. Until the default is switched
to libregraph/idm, oCIS has to be started with a custom configuration in order to use libregraph/idm as
the users and groups backend (this setup also disables the glauth and accounts service):


```
export GRAPH_IDENTITY_BACKEND=ldap
export GRAPH_LDAP_URI=ldaps://localhost:9235
export GRAPH_LDAP_BIND_DN="uid=libregraph,ou=sysusers,o=libregraph-idm"
export GRAPH_LDAP_BIND_PASSWORD=idm
export GRAPH_LDAP_USER_EMAIL_ATTRIBUTE=mail
export GRAPH_LDAP_USER_NAME_ATTRIBUTE=uid
export GRAPH_LDAP_USER_BASE_DN="ou=users,o=libregraph-idm"
export GRAPH_LDAP_GROUP_BASE_DN="ou=groups,o=libregraph-idm"
export GRAPH_LDAP_SERVER_WRITE_ENABLED="true"
export IDP_INSECURE="true"
export IDP_LDAP_FILTER="(&(objectclass=inetOrgPerson)(objectClass=owncloud))"
export IDP_LDAP_URI=ldaps://localhost:9235
export IDP_LDAP_BIND_DN="uid=idp,ou=sysusers,o=libregraph-idm"
export IDP_LDAP_BIND_PASSWORD="idp"
export IDP_LDAP_BASE_DN="ou=users,o=libregraph-idm"
export IDP_LDAP_LOGIN_ATTRIBUTE=uid
export IDP_LDAP_UUID_ATTRIBUTE="ownclouduuid"
export IDP_LDAP_UUID_ATTRIBUTE_TYPE=binary
export PROXY_ACCOUNT_BACKEND_TYPE=cs3
export OCS_ACCOUNT_BACKEND_TYPE=cs3
export STORAGE_LDAP_HOSTNAME=localhost
export STORAGE_LDAP_PORT=9235
export STORAGE_LDAP_INSECURE="true"
export STORAGE_LDAP_BASE_DN="o=libregraph-idm"
export STORAGE_LDAP_BIND_DN="uid=reva,ou=sysusers,o=libregraph-idm"
export STORAGE_LDAP_BIND_PASSWORD=reva
export STORAGE_LDAP_LOGINFILTER='(&(objectclass=inetOrgPerson)(objectclass=owncloud)(|(uid={{login}})(mail={{login}})))'
export STORAGE_LDAP_USERFILTER='(&(objectclass=inetOrgPerson)(objectclass=owncloud)(|(ownclouduuid={{.OpaqueId}})(uid={{.OpaqueId}})))'
export STORAGE_LDAP_USERATTRIBUTEFILTER='(&(objectclass=owncloud)({{attr}}={{value}}))'
export STORAGE_LDAP_USERFINDFILTER='(&(objectclass=owncloud)(|(uid={{query}}*)(cn={{query}}*)(displayname={{query}}*)(mail={{query}}*)(description={{query}}*)))'
export STORAGE_LDAP_USERGROUPFILER='(&(objectclass=groupOfNames)(member={{query}}*))'
export STORAGE_LDAP_GROUPFILTER='(&(objectclass=groupOfNames)(objectclass=owncloud)(ownclouduuid={{.OpaqueId}}*))'
export OCIS_RUN_EXTENSIONS=settings,storage-metadata,graph,graph-explorer,ocs,store,thumbnails,web,webdav,storage-frontend,storage-gateway,storage-userprovider,storage-groupprovider,storage-authbasic,storage-authbearer,storage-authmachine,storage-users,storage-shares,storage-public-link,storage-appprovider,storage-sharing,proxy,idp,nats,idm
export OCIS_INSECURE=true
bin/ocis server
```

