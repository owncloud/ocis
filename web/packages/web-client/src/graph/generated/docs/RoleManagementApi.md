# RoleManagementApi

All URIs are relative to *https://ocis.ocis.rolling.owncloud.works/graph*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**getPermissionRoleDefinition**](#getpermissionroledefinition) | **GET** /v1beta1/roleManagement/permissions/roleDefinitions/{role-id} | Get unifiedRoleDefinition|
|[**listPermissionRoleDefinitions**](#listpermissionroledefinitions) | **GET** /v1beta1/roleManagement/permissions/roleDefinitions | List roleDefinitions|

# **getPermissionRoleDefinition**
> UnifiedRoleDefinition getPermissionRoleDefinition()

Read the properties and relationships of a `unifiedRoleDefinition` object. 

### Example

```typescript
import {
    RoleManagementApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RoleManagementApi(configuration);

let roleId: string; //key: id of roleDefinition (default to undefined)

const { status, data } = await apiInstance.getPermissionRoleDefinition(
    roleId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **roleId** | [**string**] | key: id of roleDefinition | defaults to undefined|


### Return type

**UnifiedRoleDefinition**

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **listPermissionRoleDefinitions**
> UnifiedRoleDefinition listPermissionRoleDefinitions()

Get a list of `unifiedRoleDefinition` objects for the permissions provider. This list determines the roles that can be selected when creating sharing invites. 

### Example

```typescript
import {
    RoleManagementApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new RoleManagementApi(configuration);

const { status, data } = await apiInstance.listPermissionRoleDefinitions();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**UnifiedRoleDefinition**

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | A list of permission roles than can be used when sharing with users or groups. |  -  |
|**0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

