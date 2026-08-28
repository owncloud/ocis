# GroupsApi

All URIs are relative to *https://ocis.ocis.rolling.owncloud.works/graph*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**createGroup**](GroupsApi.md#creategroup) | **POST** /v1.0/groups | Add new entity to groups |
| [**listGroups**](GroupsApi.md#listgroups) | **GET** /v1.0/groups | Get entities from groups |



## createGroup

> Group createGroup(group)

Add new entity to groups

### Example

```ts
import {
  Configuration,
  GroupsApi,
} from '';
import type { CreateGroupRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new GroupsApi(config);

  const body = {
    // Group | New entity
    group: ...,
  } satisfies CreateGroupRequest;

  try {
    const data = await api.createGroup(body);
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
| **group** | [Group](Group.md) | New entity | |

### Return type

[**Group**](Group.md)

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


## listGroups

> CollectionOfGroup listGroups($search, $orderby, $select, $expand)

Get entities from groups

### Example

```ts
import {
  Configuration,
  GroupsApi,
} from '';
import type { ListGroupsRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new GroupsApi(config);

  const body = {
    // string | Search items by search phrases (optional)
    $search: $search_example,
    // Set<'displayName' | 'displayName desc'> | Order items by property values (optional)
    $orderby: ...,
    // Set<'id' | 'description' | 'displayName' | 'mail' | 'members'> | Select properties to be returned (optional)
    $select: ...,
    // Set<'members'> | Expand related entities (optional)
    $expand: ...,
  } satisfies ListGroupsRequest;

  try {
    const data = await api.listGroups(body);
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
| **$orderby** | `displayName`, `displayName desc` | Order items by property values | [Optional] [Enum: displayName, displayName desc] |
| **$select** | `id`, `description`, `displayName`, `mail`, `members` | Select properties to be returned | [Optional] [Enum: id, description, displayName, mail, members] |
| **$expand** | `members` | Expand related entities | [Optional] [Enum: members] |

### Return type

[**CollectionOfGroup**](CollectionOfGroup.md)

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

