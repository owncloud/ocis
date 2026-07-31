# DrivesPermissionsApi

All URIs are relative to *https://ocis.ocis.rolling.owncloud.works/graph*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**createLink**](#createlink) | **POST** /v1beta1/drives/{drive-id}/items/{item-id}/createLink | Create a sharing link for a DriveItem|
|[**deletePermission**](#deletepermission) | **DELETE** /v1beta1/drives/{drive-id}/items/{item-id}/permissions/{perm-id} | Remove access to a DriveItem|
|[**getPermission**](#getpermission) | **GET** /v1beta1/drives/{drive-id}/items/{item-id}/permissions/{perm-id} | Get sharing permission for a file or folder|
|[**invite**](#invite) | **POST** /v1beta1/drives/{drive-id}/items/{item-id}/invite | Send a sharing invitation|
|[**listPermissions**](#listpermissions) | **GET** /v1beta1/drives/{drive-id}/items/{item-id}/permissions | List the effective sharing permissions on a driveItem.|
|[**setPermissionPassword**](#setpermissionpassword) | **POST** /v1beta1/drives/{drive-id}/items/{item-id}/permissions/{perm-id}/setPassword | Set sharing link password|
|[**updatePermission**](#updatepermission) | **PATCH** /v1beta1/drives/{drive-id}/items/{item-id}/permissions/{perm-id} | Update sharing permission|

# **createLink**
> Permission createLink()

You can use the createLink action to share a driveItem via a sharing link.  The response will be a permission object with the link facet containing the created link details.  ## Link types  For now, The following values are allowed for the type parameter.  | Value          | Display name      | Description                                                     | | -------------- | ----------------- | --------------------------------------------------------------- | | view           | View              | Creates a read-only link to the driveItem.                      | | upload         | Upload            | Creates a read-write link to the folder driveItem.              | | edit           | Edit              | Creates a read-write link to the driveItem.                     | | createOnly     | File Drop         | Creates an upload-only link to the folder driveItem.            | | blocksDownload | Secure View       | Creates a read-only link that blocks download to the driveItem. | 

### Example

```typescript
import {
    DrivesPermissionsApi,
    Configuration,
    DriveItemCreateLink
} from './api';

const configuration = new Configuration();
const apiInstance = new DrivesPermissionsApi(configuration);

let driveId: string; //key: id of drive (default to undefined)
let itemId: string; //key: id of item (default to undefined)
let driveItemCreateLink: DriveItemCreateLink; //In the request body, provide a JSON object with the following parameters. (optional)

const { status, data } = await apiInstance.createLink(
    driveId,
    itemId,
    driveItemCreateLink
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **driveItemCreateLink** | **DriveItemCreateLink**| In the request body, provide a JSON object with the following parameters. | |
| **driveId** | [**string**] | key: id of drive | defaults to undefined|
| **itemId** | [**string**] | key: id of item | defaults to undefined|


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

# **deletePermission**
> deletePermission()

Remove access to a DriveItem.  Only sharing permissions that are not inherited can be deleted. The `inheritedFrom` property must be `null`. 

### Example

```typescript
import {
    DrivesPermissionsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DrivesPermissionsApi(configuration);

let driveId: string; //key: id of drive (default to undefined)
let itemId: string; //key: id of item (default to undefined)
let permId: string; //key: id of permission (default to undefined)

const { status, data } = await apiInstance.deletePermission(
    driveId,
    itemId,
    permId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **driveId** | [**string**] | key: id of drive | defaults to undefined|
| **itemId** | [**string**] | key: id of item | defaults to undefined|
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

# **getPermission**
> Permission getPermission()

Return the effective sharing permission for a particular permission resource. 

### Example

```typescript
import {
    DrivesPermissionsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DrivesPermissionsApi(configuration);

let driveId: string; //key: id of drive (default to undefined)
let itemId: string; //key: id of item (default to undefined)
let permId: string; //key: id of permission (default to undefined)

const { status, data } = await apiInstance.getPermission(
    driveId,
    itemId,
    permId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **driveId** | [**string**] | key: id of drive | defaults to undefined|
| **itemId** | [**string**] | key: id of item | defaults to undefined|
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

# **invite**
> CollectionOfPermissions invite()

Sends a sharing invitation for a `driveItem`. A sharing invitation provides permissions to the recipients and optionally sends them an email with a sharing link.  The response will be a permission object with the grantedToV2 property containing the created grant details.  ## Roles property values For now, roles are only identified by a uuid. There are no hardcoded aliases like `read` or `write` because role actions can be completely customized. 

### Example

```typescript
import {
    DrivesPermissionsApi,
    Configuration,
    DriveItemInvite
} from './api';

const configuration = new Configuration();
const apiInstance = new DrivesPermissionsApi(configuration);

let driveId: string; //key: id of drive (default to undefined)
let itemId: string; //key: id of item (default to undefined)
let driveItemInvite: DriveItemInvite; //In the request body, provide a JSON object with the following parameters. To create a custom role submit a list of actions instead of roles. (optional)

const { status, data } = await apiInstance.invite(
    driveId,
    itemId,
    driveItemInvite
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **driveItemInvite** | **DriveItemInvite**| In the request body, provide a JSON object with the following parameters. To create a custom role submit a list of actions instead of roles. | |
| **driveId** | [**string**] | key: id of drive | defaults to undefined|
| **itemId** | [**string**] | key: id of item | defaults to undefined|


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

# **listPermissions**
> CollectionOfPermissionsWithAllowedValues listPermissions()

The permissions collection includes potentially sensitive information and may not be available for every caller.  * For the owner of the item, all sharing permissions will be returned. This includes co-owners. * For a non-owner caller, only the sharing permissions that apply to the caller are returned. * Sharing permission properties that contain secrets (e.g. `webUrl`) are only returned for callers that are able to create the sharing permission.  All permission objects have an `id`. A permission representing * a link has the `link` facet filled with details. * a share has the `roles` property set and the `grantedToV2` property filled with the grant recipient details. 

### Example

```typescript
import {
    DrivesPermissionsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DrivesPermissionsApi(configuration);

let driveId: string; //key: id of drive (default to undefined)
let itemId: string; //key: id of item (default to undefined)
let $filter: string; //Filter items by property values. By default all permissions are returned and the avalable sharing roles are limited to normal users. To get a list of sharing roles applicable to federated users use the example $select query and combine it with $filter to omit the list of permissions. (optional) (default to undefined)
let $select: Set<'@libre.graph.permissions.actions.allowedValues' | '@libre.graph.permissions.roles.allowedValues' | 'value'>; //Select properties to be returned. By default all properties are returned. Select the roles property to fetch the available sharing roles without resolving all the permissions. Combine this with the $filter parameter to fetch the actions applicable to federated users. (optional) (default to undefined)

const { status, data } = await apiInstance.listPermissions(
    driveId,
    itemId,
    $filter,
    $select
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **driveId** | [**string**] | key: id of drive | defaults to undefined|
| **itemId** | [**string**] | key: id of item | defaults to undefined|
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

# **setPermissionPassword**
> Permission setPermissionPassword(sharingLinkPassword)

Set the password of a sharing permission.  Only the `password` property can be modified this way. 

### Example

```typescript
import {
    DrivesPermissionsApi,
    Configuration,
    SharingLinkPassword
} from './api';

const configuration = new Configuration();
const apiInstance = new DrivesPermissionsApi(configuration);

let driveId: string; //key: id of drive (default to undefined)
let itemId: string; //key: id of item (default to undefined)
let permId: string; //key: id of permission (default to undefined)
let sharingLinkPassword: SharingLinkPassword; //New password value

const { status, data } = await apiInstance.setPermissionPassword(
    driveId,
    itemId,
    permId,
    sharingLinkPassword
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **sharingLinkPassword** | **SharingLinkPassword**| New password value | |
| **driveId** | [**string**] | key: id of drive | defaults to undefined|
| **itemId** | [**string**] | key: id of item | defaults to undefined|
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

# **updatePermission**
> Permission updatePermission(permission)

Update the properties of a sharing permission by patching the permission resource.  Only the `roles`, `expirationDateTime` and `password` properties can be modified this way. 

### Example

```typescript
import {
    DrivesPermissionsApi,
    Configuration,
    Permission
} from './api';

const configuration = new Configuration();
const apiInstance = new DrivesPermissionsApi(configuration);

let driveId: string; //key: id of drive (default to undefined)
let itemId: string; //key: id of item (default to undefined)
let permId: string; //key: id of permission (default to undefined)
let permission: Permission; //New property values

const { status, data } = await apiInstance.updatePermission(
    driveId,
    itemId,
    permId,
    permission
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **permission** | **Permission**| New property values | |
| **driveId** | [**string**] | key: id of drive | defaults to undefined|
| **itemId** | [**string**] | key: id of item | defaults to undefined|
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

