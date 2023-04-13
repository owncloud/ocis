# Frontend

The frontend service translates various owncloud related HTTP APIs to CS3 requests.

## Endpoints Overview

Currently, the frontend service handles requests for three functionalities, which are `appprovider`, `archiver`, `datagateway` and `ocs`.

### appprovider

The appprovider endpoint, by default `/app`, forwards HTTP requests to the CS3 [App Registry API](https://cs3org.github.io/cs3apis/#cs3.app.registry.v1beta1.RegistryAPI)

### archiver

The archiver endpoint, by default `/archiver`, implements zip and tar download for collections of files. It will internally use the CS3 API to initiate downloads and then stream the individual files as part of a compressed file.

### datagateway

The datagateway endpoint, by default `/data`, forwards file up- and download requests to the correct CS3 data provider. OCIS starts a dataprovider as part of the storage-* services. The routing happens based on the JWT that was created by a storage provider in response to an `InitiateFileDownload` or `InitiateFileUpload` request.

### ocs

The ocs endpoint, by default `/ocs`, implements the ownCloud 10 Open Collaboration Services API by translating it into CS3 API requests. It can handle users, groups, capabilities and also implements the files sharing functionality on top of CS3. The `/ocs/v[12].php/cloud/user/signing-key` is currently handled by the dedicated [ocs](https://github.com/owncloud/ocis/tree/master/services/ocs) service.

## Scalability

While the frontend service does not persist any data it does cache `Stat()` responses and user information. Therefore, multiple instances of this service can be spawned in a bigger deployment like when using container orchestration with Kubernetes, when configuring `FRONTEND_OCS_RESOURCE_INFO_CACHE_TYPE=redis` and the related config options.

## libregraph service interactions

A lot of user management is done via a standardized libregraph API.
Depending on how the system is configured there might be some attributes
for users that an instance admin user can't change because of properties
coming from an external LDAP server, or similar. To make life easier for
admin users there are hints as capabilites telling which attributes are
read-only or not. To configure these hints we have the following
environment variables:

- FRONTEND_READONLY_ATTRIBUTES_ACCOUNT_ENABLED: Default is false
- FRONTEND_READONLY_ATTRIBUTES_DISPLAY_NAME: Default is true
- FRONTEND_READONLY_ATTRIBUTES_GIVEN_NAME: Default is true
- FRONTEND_READONLY_ATTRIBUTES_ID: Default is true
- FRONTEND_READONLY_ATTRIBUTES_MAIL: Default is true
- FRONTEND_READONLY_ATTRIBUTES_ON_PREMISES_SAM_ACCOUNT_NAME: Default is true
- FRONTEND_READONLY_ATTRIBUTES_SURNAME: Default is true
- FRONTEND_READONLY_ATTRIBUTES_QUOTA: Default is false
