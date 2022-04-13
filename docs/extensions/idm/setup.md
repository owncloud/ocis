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
export LDAP_URI=ldaps://localhost:9235
export LDAP_INSECURE="true"
export LDAP_USER_BASE_DN="ou=users,o=libregraph-idm"
export LDAP_USER_SCHEMA_ID="ownclouduuid"
export LDAP_USER_SCHEMA_MAIL="mail"
export LDAP_USER_SCHEMA_USERNAME="uid"
export LDAP_USER_OBJECTCLASS="inetOrgPerson"
export LDAP_GROUP_BASE_DN="ou=groups,o=libregraph-idm"
export LDAP_GROUP_SCHEMA_ID="ownclouduuid"
export LDAP_GROUP_SCHEMA_MAIL="mail"
export LDAP_GROUP_SCHEMA_GROUPNAME="member"
export LDAP_GROUP_OBJECTCLASS="groupOfNames"
export GRAPH_LDAP_BIND_DN="uid=libregraph,ou=sysusers,o=libregraph-idm"
export GRAPH_LDAP_BIND_PASSWORD=idm
export GRAPH_LDAP_SERVER_WRITE_ENABLED="true"
export IDP_INSECURE="true"
export IDP_LDAP_BIND_DN="uid=idp,ou=sysusers,o=libregraph-idm"
export IDP_LDAP_BIND_PASSWORD="idp"
export IDP_LDAP_LOGIN_ATTRIBUTE=uid
export PROXY_ACCOUNT_BACKEND_TYPE=cs3
export OCS_ACCOUNT_BACKEND_TYPE=cs3
export STORAGE_LDAP_BIND_DN="uid=reva,ou=sysusers,o=libregraph-idm"
export STORAGE_LDAP_BIND_PASSWORD=reva
export OCIS_RUN_EXTENSIONS=settings,storage-metadata,graph,graph-explorer,ocs,store,thumbnails,web,webdav,storage-frontend,storage-gateway,storage-userprovider,storage-groupprovider,storage-authbasic,storage-authbearer,storage-authmachine,storage-users,storage-shares,storage-public-link,storage-appprovider,storage-sharing,proxy,idp,nats,idm,ocdav
export OCIS_INSECURE=true
bin/ocis server
```

