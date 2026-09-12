# UsersApi

All URIs are relative to *https://ocis.ocis.rolling.owncloud.works/graph*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**createUser**](UsersApi.md#createuser) | **POST** /v1.0/users | Add new entity to users |
| [**listUsers**](UsersApi.md#listusers) | **GET** /v1.0/users | Get entities from users |



## createUser

> User createUser(user)

Add new entity to users

### Example

```ts
import {
  Configuration,
  UsersApi,
} from '';
import type { CreateUserRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new UsersApi(config);

  const body = {
    // User | New entity
    user: ...,
  } satisfies CreateUserRequest;

  try {
    const data = await api.createUser(body);
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
| **user** | [User](User.md) | New entity | |

### Return type

[**User**](User.md)

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **201** | Created entity |  -  |
| **0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## listUsers

> CollectionOfUser listUsers($search, $filter, $orderby, $select, $expand)

Get entities from users

### Example

```ts
import {
  Configuration,
  UsersApi,
} from '';
import type { ListUsersRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new UsersApi(config);

  const body = {
    // string | Search items by search phrases (optional)
    $search: $search_example,
    // string | Filter users by property values and relationship attributes (optional)
    $filter: memberOf/any(x:x/id eq 910367f9-4041-4db1-961b-d1e98f708eaf),
    // Set<'displayName' | 'displayName desc' | 'mail' | 'mail desc' | 'onPremisesSamAccountName' | 'onPremisesSamAccountName desc'> | Order items by property values (optional)
    $orderby: ...,
    // Set<'id' | 'displayName' | 'mail' | 'memberOf' | 'onPremisesSamAccountName' | 'surname'> | Select properties to be returned (optional)
    $select: ...,
    // Set<'drive' | 'drives' | 'memberOf' | 'appRoleAssignments'> | Expand related entities (optional)
    $expand: ...,
  } satisfies ListUsersRequest;

  try {
    const data = await api.listUsers(body);
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
| **$search** | `string` | Search items by search phrases | [Optional] [Defaults to `undefined`] |
| **$filter** | `string` | Filter users by property values and relationship attributes | [Optional] [Defaults to `undefined`] |
| **$orderby** | `displayName`, `displayName desc`, `mail`, `mail desc`, `onPremisesSamAccountName`, `onPremisesSamAccountName desc` | Order items by property values | [Optional] [Enum: displayName, displayName desc, mail, mail desc, onPremisesSamAccountName, onPremisesSamAccountName desc] |
| **$select** | `id`, `displayName`, `mail`, `memberOf`, `onPremisesSamAccountName`, `surname` | Select properties to be returned | [Optional] [Enum: id, displayName, mail, memberOf, onPremisesSamAccountName, surname] |
| **$expand** | `drive`, `drives`, `memberOf`, `appRoleAssignments` | Expand related entities | [Optional] [Enum: drive, drives, memberOf, appRoleAssignments] |

### Return type

[**CollectionOfUser**](CollectionOfUser.md)

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Retrieved entities |  -  |
| **0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

