# DrivesRootApi

All URIs are relative to *https://ocis.ocis.rolling.owncloud.works/graph*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**createDriveItem**](#createdriveitem) | **POST** /v1beta1/drives/{drive-id}/root/children | Create a drive item|
|[**createLinkSpaceRoot**](#createlinkspaceroot) | **POST** /v1beta1/drives/{drive-id}/root/createLink | Create a sharing link for the root item of a Drive|
|[**deletePermissionSpaceRoot**](#deletepermissionspaceroot) | **DELETE** /v1beta1/drives/{drive-id}/root/permissions/{perm-id} | Remove access to a Drive|
|[**getPermissionSpaceRoot**](#getpermissionspaceroot) | **GET** /v1beta1/drives/{drive-id}/root/permissions/{perm-id} | Get a single sharing permission for the root item of a drive|
|[**getRoot**](#getroot) | **GET** /v1.0/drives/{drive-id}/root | Get root from arbitrary space|
|[**inviteSpaceRoot**](#invitespaceroot) | **POST** /v1beta1/drives/{drive-id}/root/invite | Send a sharing invitation|
|[**listPermissionsSpaceRoot**](#listpermissionsspaceroot) | **GET** /v1beta1/drives/{drive-id}/root/permissions | List the effective permissions on the root item of a drive.|
|[**setPermissionPasswordSpaceRoot**](#setpermissionpasswordspaceroot) | **POST** /v1beta1/drives/{drive-id}/root/permissions/{perm-id}/setPassword | Set sharing link password for the root item of a drive|
|[**updatePermissionSpaceRoot**](#updatepermissionspaceroot) | **PATCH** /v1beta1/drives/{drive-id}/root/permissions/{perm-id} | Update sharing permission|

# **createDriveItem**
> DriveItem createDriveItem()

You can use the root childrens endpoint to mount a remoteItem in the share jail. The `@client.synchronize` property of the `driveItem` in the [sharedWithMe](#/me.drive/ListSharedWithMe) endpoint will change to true. 

### Example

```typescript
import {
    DrivesRootApi,
    Configuration,
    DriveItem
} from './api';

const configuration = new Configuration();
const apiInstance = new DrivesRootApi(configuration);

let driveId: string; //key: id of drive (default to undefined)
let driveItem: DriveItem; //In the request body, provide a JSON object with the following parameters. For mounting a share the necessary remoteItem id and permission id can be taken from the [sharedWithMe](#/me.drive/ListSharedWithMe) endpoint. (optional)

const { status, data } = await apiInstance.createDriveItem(
    driveId,
    driveItem
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **driveItem** | **DriveItem**| In the request body, provide a JSON object with the following parameters. For mounting a share the necessary remoteItem id and permission id can be taken from the [sharedWithMe](#/me.drive/ListSharedWithMe) endpoint. | |
| **driveId** | [**string**] | key: id of drive | defaults to undefined|


### Return type

**DriveItem**

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Response |  -  |
|**0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **createLinkSpaceRoot**
> Permission createLinkSpaceRoot()

You can use the createLink action to share a driveItem via a sharing link.  The response will be a permission object with the link facet containing the created link details.  ## Link types  For now, The following values are allowed for the type parameter.  | Value          | Display name      | Description                                                     | | -------------- | ----------------- | --------------------------------------------------------------- | | view           | View              | Creates a read-only link to the driveItem.                      | | upload         | Upload            | Creates a read-write link to the folder driveItem.              | | edit           | Edit              | Creates a read-write link to the driveItem.                     | | createOnly     | File Drop         | Creates an upload-only link to the folder driveItem.            | | blocksDownload | Secure View       | Creates a read-only link that blocks download to the driveItem. | 

### Example

```typescript
import {
    DrivesRootApi,
    Configuration,
    DriveItemCreateLink
} from './api';

const configuration = new Configuration();
const apiInstance = new DrivesRootApi(configuration);

let driveId: string; //key: id of drive (default to undefined)
let driveItemCreateLink: DriveItemCreateLink; //In the request body, provide a JSON object with the following parameters. (optional)

const { status, data } = await apiInstance.createLinkSpaceRoot(
    driveId,
    driveItemCreateLink
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **driveItemCreateLink** | **DriveItemCreateLink**| In the request body, provide a JSON object with the following parameters. | |
| **driveId** | [**string**] | key: id of drive | defaults to undefined|


### Return type

**Permission**

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Response |  -  |
|**207** | Partial success response TODO |  -  |
|**0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deletePermissionSpaceRoot**
> deletePermissionSpaceRoot()

Remove access to the root item of a drive.  Only sharing permissions that are not inherited can be deleted. The `inheritedFrom` property must be `null`. 

### Example

```typescript
import {
    DrivesRootApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DrivesRootApi(configuration);

let driveId: string; //key: id of drive (default to undefined)
let permId: string; //key: id of permission (default to undefined)

const { status, data } = await apiInstance.deletePermissionSpaceRoot(
    driveId,
    permId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **driveId** | [**string**] | key: id of drive | defaults to undefined|
| **permId** | [**string**] | key: id of permission | defaults to undefined|


### Return type

void (empty response body)

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**204** | Success |  -  |
|**0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getPermissionSpaceRoot**
> Permission getPermissionSpaceRoot()

Return the effective sharing permission for a particular permission resource. 

### Example

```typescript
import {
    DrivesRootApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DrivesRootApi(configuration);

let driveId: string; //key: id of drive (default to undefined)
let permId: string; //key: id of permission (default to undefined)

const { status, data } = await apiInstance.getPermissionSpaceRoot(
    driveId,
    permId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **driveId** | [**string**] | key: id of drive | defaults to undefined|
| **permId** | [**string**] | key: id of permission | defaults to undefined|


### Return type

**Permission**

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Retrieved resource |  -  |
|**0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getRoot**
> DriveItem getRoot()


### Example

```typescript
import {
    DrivesRootApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DrivesRootApi(configuration);

let driveId: string; //key: id of drive (default to undefined)

const { status, data } = await apiInstance.getRoot(
    driveId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **driveId** | [**string**] | key: id of drive | defaults to undefined|


### Return type

**DriveItem**

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Retrieved resource |  -  |
|**0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **inviteSpaceRoot**
> CollectionOfPermissions inviteSpaceRoot()

Sends a sharing invitation for the root of a `drive`. A sharing invitation provides permissions to the recipients and optionally sends them an email with a sharing link.  The response will be a permission object with the grantedToV2 property containing the created grant details.  ## Roles property values For now, roles are only identified by a uuid. There are no hardcoded aliases like `read` or `write` because role actions can be completely customized. 

### Example

```typescript
import {
    DrivesRootApi,
    Configuration,
    DriveItemInvite
} from './api';

const configuration = new Configuration();
const apiInstance = new DrivesRootApi(configuration);

let driveId: string; //key: id of drive (default to undefined)
let driveItemInvite: DriveItemInvite; //In the request body, provide a JSON object with the following parameters. To create a custom role submit a list of actions instead of roles. (optional)

const { status, data } = await apiInstance.inviteSpaceRoot(
    driveId,
    driveItemInvite
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **driveItemInvite** | **DriveItemInvite**| In the request body, provide a JSON object with the following parameters. To create a custom role submit a list of actions instead of roles. | |
| **driveId** | [**string**] | key: id of drive | defaults to undefined|


### Return type

**CollectionOfPermissions**

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Response |  -  |
|**207** | Partial success response TODO |  -  |
|**400** | Bad request |  -  |
|**0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **listPermissionsSpaceRoot**
> CollectionOfPermissionsWithAllowedValues listPermissionsSpaceRoot()

The permissions collection includes potentially sensitive information and may not be available for every caller.  * For the owner of the item, all sharing permissions will be returned. This includes co-owners. * For a non-owner caller, only the sharing permissions that apply to the caller are returned. * Sharing permission properties that contain secrets (e.g. `webUrl`) are only returned for callers that are able to create the sharing permission.  All permission objects have an `id`. A permission representing * a link has the `link` facet filled with details. * a share has the `roles` property set and the `grantedToV2` property filled with the grant recipient details. 

### Example

```typescript
import {
    DrivesRootApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DrivesRootApi(configuration);

let driveId: string; //key: id of drive (default to undefined)
let $filter: string; //Filter items by property values. By default all permissions are returned and the avalable sharing roles are limited to normal users. To get a list of sharing roles applicable to federated users use the example $select query and combine it with $filter to omit the list of permissions. (optional) (default to undefined)
let $select: Set<'@libre.graph.permissions.actions.allowedValues' | '@libre.graph.permissions.roles.allowedValues' | 'value'>; //Select properties to be returned. By default all properties are returned. Select the roles property to fetch the available sharing roles without resolving all the permissions. Combine this with the $filter parameter to fetch the actions applicable to federated users. (optional) (default to undefined)

const { status, data } = await apiInstance.listPermissionsSpaceRoot(
    driveId,
    $filter,
    $select
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **driveId** | [**string**] | key: id of drive | defaults to undefined|
| **$filter** | [**string**] | Filter items by property values. By default all permissions are returned and the avalable sharing roles are limited to normal users. To get a list of sharing roles applicable to federated users use the example $select query and combine it with $filter to omit the list of permissions. | (optional) defaults to undefined|
| **$select** | **Array<&#39;@libre.graph.permissions.actions.allowedValues&#39; &#124; &#39;@libre.graph.permissions.roles.allowedValues&#39; &#124; &#39;value&#39;>** | Select properties to be returned. By default all properties are returned. Select the roles property to fetch the available sharing roles without resolving all the permissions. Combine this with the $filter parameter to fetch the actions applicable to federated users. | (optional) defaults to undefined|


### Return type

**CollectionOfPermissionsWithAllowedValues**

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Retrieved resource |  -  |
|**0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **setPermissionPasswordSpaceRoot**
> Permission setPermissionPasswordSpaceRoot(sharingLinkPassword)

Set the password of a sharing permission.  Only the `password` property can be modified this way. 

### Example

```typescript
import {
    DrivesRootApi,
    Configuration,
    SharingLinkPassword
} from './api';

const configuration = new Configuration();
const apiInstance = new DrivesRootApi(configuration);

let driveId: string; //key: id of drive (default to undefined)
let permId: string; //key: id of permission (default to undefined)
let sharingLinkPassword: SharingLinkPassword; //New password value

const { status, data } = await apiInstance.setPermissionPasswordSpaceRoot(
    driveId,
    permId,
    sharingLinkPassword
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **sharingLinkPassword** | **SharingLinkPassword**| New password value | |
| **driveId** | [**string**] | key: id of drive | defaults to undefined|
| **permId** | [**string**] | key: id of permission | defaults to undefined|


### Return type

**Permission**

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Updated permission |  -  |
|**0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updatePermissionSpaceRoot**
> Permission updatePermissionSpaceRoot(permission)

Update the properties of a sharing permission by patching the permission resource.  Only the `roles`, `expirationDateTime` and `password` properties can be modified this way. 

### Example

```typescript
import {
    DrivesRootApi,
    Configuration,
    Permission
} from './api';

const configuration = new Configuration();
const apiInstance = new DrivesRootApi(configuration);

let driveId: string; //key: id of drive (default to undefined)
let permId: string; //key: id of permission (default to undefined)
let permission: Permission; //New property values

const { status, data } = await apiInstance.updatePermissionSpaceRoot(
    driveId,
    permId,
    permission
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **permission** | **Permission**| New property values | |
| **driveId** | [**string**] | key: id of drive | defaults to undefined|
| **permId** | [**string**] | key: id of permission | defaults to undefined|


### Return type

**Permission**

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Updated permission |  -  |
|**0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

