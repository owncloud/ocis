# UserAppRoleAssignmentApi

All URIs are relative to *https://ocis.ocis.rolling.owncloud.works/graph*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**userCreateAppRoleAssignments**](UserAppRoleAssignmentApi.md#usercreateapproleassignments) | **POST** /v1.0/users/{user-id}/appRoleAssignments | Grant an appRoleAssignment to a user |
| [**userDeleteAppRoleAssignments**](UserAppRoleAssignmentApi.md#userdeleteapproleassignments) | **DELETE** /v1.0/users/{user-id}/appRoleAssignments/{appRoleAssignment-id} | Delete the appRoleAssignment from a user |
| [**userListAppRoleAssignments**](UserAppRoleAssignmentApi.md#userlistapproleassignments) | **GET** /v1.0/users/{user-id}/appRoleAssignments | Get appRoleAssignments from a user |



## userCreateAppRoleAssignments

> AppRoleAssignment userCreateAppRoleAssignments(userId, appRoleAssignment)

Grant an appRoleAssignment to a user

Use this API to assign a global role to a user. To grant an app role assignment to a user, you need three identifiers: * &#x60;principalId&#x60;: The &#x60;id&#x60; of the user to whom you are assigning the app role. * &#x60;resourceId&#x60;: The &#x60;id&#x60; of the resource &#x60;servicePrincipal&#x60; or &#x60;application&#x60; that has defined the app role. * &#x60;appRoleId&#x60;: The &#x60;id&#x60; of the &#x60;appRole&#x60; (defined on the resource service principal or application) to assign to the user. 

### Example

```ts
import {
  Configuration,
  UserAppRoleAssignmentApi,
} from '';
import type { UserCreateAppRoleAssignmentsRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new UserAppRoleAssignmentApi(config);

  const body = {
    // string | key: id of user
    userId: userId_example,
    // AppRoleAssignment | New app role assignment value
    appRoleAssignment: ...,
  } satisfies UserCreateAppRoleAssignmentsRequest;

  try {
    const data = await api.userCreateAppRoleAssignments(body);
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
| **userId** | `string` | key: id of user | [Defaults to `undefined`] |
| **appRoleAssignment** | [AppRoleAssignment](AppRoleAssignment.md) | New app role assignment value | |

### Return type

[**AppRoleAssignment**](AppRoleAssignment.md)

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Created new app role assignment. |  -  |
| **0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## userDeleteAppRoleAssignments

> userDeleteAppRoleAssignments(userId, appRoleAssignmentId, ifMatch)

Delete the appRoleAssignment from a user

### Example

```ts
import {
  Configuration,
  UserAppRoleAssignmentApi,
} from '';
import type { UserDeleteAppRoleAssignmentsRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new UserAppRoleAssignmentApi(config);

  const body = {
    // string | key: id of user
    userId: userId_example,
    // string | key: id of appRoleAssignment. This is the concatenated {user-id}:{appRole-id} separated by a colon.
    appRoleAssignmentId: appRoleAssignmentId_example,
    // string | ETag (optional)
    ifMatch: ifMatch_example,
  } satisfies UserDeleteAppRoleAssignmentsRequest;

  try {
    const data = await api.userDeleteAppRoleAssignments(body);
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
| **userId** | `string` | key: id of user | [Defaults to `undefined`] |
| **appRoleAssignmentId** | `string` | key: id of appRoleAssignment. This is the concatenated {user-id}:{appRole-id} separated by a colon. | [Defaults to `undefined`] |
| **ifMatch** | `string` | ETag | [Optional] [Defaults to `undefined`] |

### Return type

`void` (Empty response body)

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **204** | Success |  -  |
| **0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## userListAppRoleAssignments

> CollectionOfAppRoleAssignments userListAppRoleAssignments(userId)

Get appRoleAssignments from a user

Represents the global roles a user has been granted for an application.

### Example

```ts
import {
  Configuration,
  UserAppRoleAssignmentApi,
} from '';
import type { UserListAppRoleAssignmentsRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new UserAppRoleAssignmentApi(config);

  const body = {
    // string | key: id of user
    userId: userId_example,
  } satisfies UserListAppRoleAssignmentsRequest;

  try {
    const data = await api.userListAppRoleAssignments(body);
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
| **userId** | `string` | key: id of user | [Defaults to `undefined`] |

### Return type

[**CollectionOfAppRoleAssignments**](CollectionOfAppRoleAssignments.md)

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Retrieved appRoleAssignments |  -  |
| **0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

