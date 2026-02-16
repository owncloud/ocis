Bugfix: Replace obsolete docker image in the deployment example

In the ocis_ldap deployment example, we were using the bitnami/openldap docker
image. This image isn't available any longer, so the example couldn't be
deployed as intended.

We've replaced the docker image with the osixia/openldap image and we've
adjusted some of the configuration of the openldap image.

https://github.com/owncloud/ocis/pull/11828 
