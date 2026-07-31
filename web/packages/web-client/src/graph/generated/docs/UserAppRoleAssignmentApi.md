# UserAppRoleAssignmentApi

All URIs are relative to *https://ocis.ocis.rolling.owncloud.works/graph*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**userCreateAppRoleAssignments**](#usercreateapproleassignments) | **POST** /v1.0/users/{user-id}/appRoleAssignments | Grant an appRoleAssignment to a user|
|[**userDeleteAppRoleAssignments**](#userdeleteapproleassignments) | **DELETE** /v1.0/users/{user-id}/appRoleAssignments/{appRoleAssignment-id} | Delete the appRoleAssignment from a user|
|[**userListAppRoleAssignments**](#userlistapproleassignments) | **GET** /v1.0/users/{user-id}/appRoleAssignments | Get appRoleAssignments from a user|

# **userCreateAppRoleAssignments**
> AppRoleAssignment userCreateAppRoleAssignments(appRoleAssignment)

Use this API to assign a global role to a user. To grant an app role assignment to a user, you need three identifiers: * `principalId`: The `id` of the user to whom you are assigning the app role. * `resourceId`: The `id` of the resource `servicePrincipal` or `application` that has defined the app role. * `appRoleId`: The `id` of the `appRole` (defined on the resource service principal or application) to assign to the user. 

### Example

```typescript
import {
    UserAppRoleAssignmentApi,
    Configuration,
    AppRoleAssignment
} from './api';

const configuration = new Configuration();
const apiInstance = new UserAppRoleAssignmentApi(configuration);

let userId: string; //key: id of user (default to undefined)
let appRoleAssignment: AppRoleAssignment; //New app role assignment value

const { status, data } = await apiInstance.userCreateAppRoleAssignments(
    userId,
    appRoleAssignment
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **appRoleAssignment** | **AppRoleAssignment**| New app role assignment value | |
| **userId** | [**string**] | key: id of user | defaults to undefined|


### Return type

**AppRoleAssignment**

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Created new app role assignment. |  -  |
|**0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **userDeleteAppRoleAssignments**
> userDeleteAppRoleAssignments()


### Example

```typescript
import {
    UserAppRoleAssignmentApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new UserAppRoleAssignmentApi(configuration);

let userId: string; //key: id of user (default to undefined)
let appRoleAssignmentId: string; //key: id of appRoleAssignment. This is the concatenated {user-id}:{appRole-id} separated by a colon. (default to undefined)
let ifMatch: string; //ETag (optional) (default to undefined)

const { status, data } = await apiInstance.userDeleteAppRoleAssignments(
    userId,
    appRoleAssignmentId,
    ifMatch
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **userId** | [**string**] | key: id of user | defaults to undefined|
| **appRoleAssignmentId** | [**string**] | key: id of appRoleAssignment. This is the concatenated {user-id}:{appRole-id} separated by a colon. | defaults to undefined|
| **ifMatch** | [**string**] | ETag | (optional) defaults to undefined|


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

# **userListAppRoleAssignments**
> CollectionOfAppRoleAssignments userListAppRoleAssignments()

Represents the global roles a user has been granted for an application.

### Example

```typescript
import {
    UserAppRoleAssignmentApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new UserAppRoleAssignmentApi(configuration);

let userId: string; //key: id of user (default to undefined)

const { status, data } = await apiInstance.userListAppRoleAssignments(
    userId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **userId** | [**string**] | key: id of user | defaults to undefined|


### Return type

**CollectionOfAppRoleAssignments**

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Retrieved appRoleAssignments |  -  |
|**0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

