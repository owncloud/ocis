# RoleManagementApi

All URIs are relative to *https://ocis.ocis.rolling.owncloud.works/graph*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**getPermissionRoleDefinition**](RoleManagementApi.md#getpermissionroledefinition) | **GET** /v1beta1/roleManagement/permissions/roleDefinitions/{role-id} | Get unifiedRoleDefinition |
| [**listPermissionRoleDefinitions**](RoleManagementApi.md#listpermissionroledefinitions) | **GET** /v1beta1/roleManagement/permissions/roleDefinitions | List roleDefinitions |



## getPermissionRoleDefinition

> UnifiedRoleDefinition getPermissionRoleDefinition(roleId)

Get unifiedRoleDefinition

Read the properties and relationships of a &#x60;unifiedRoleDefinition&#x60; object. 

### Example

```ts
import {
  Configuration,
  RoleManagementApi,
} from '';
import type { GetPermissionRoleDefinitionRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new RoleManagementApi(config);

  const body = {
    // string | key: id of roleDefinition
    roleId: roleId_example,
  } satisfies GetPermissionRoleDefinitionRequest;

  try {
    const data = await api.getPermissionRoleDefinition(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **roleId** | `string` | key: id of roleDefinition | [Defaults to `undefined`] |

### Return type

[**UnifiedRoleDefinition**](UnifiedRoleDefinition.md)

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## listPermissionRoleDefinitions

> Array&lt;UnifiedRoleDefinition&gt; listPermissionRoleDefinitions()

List roleDefinitions

Get a list of &#x60;unifiedRoleDefinition&#x60; objects for the permissions provider. This list determines the roles that can be selected when creating sharing invites. 

### Example

```ts
import {
  Configuration,
  RoleManagementApi,
} from '';
import type { ListPermissionRoleDefinitionsRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new RoleManagementApi(config);

  try {
    const data = await api.listPermissionRoleDefinitions();
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

This endpoint does not need any parameter.

### Return type

[**Array&lt;UnifiedRoleDefinition&gt;**](UnifiedRoleDefinition.md)

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | A list of permission roles than can be used when sharing with users or groups. |  -  |
| **0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

