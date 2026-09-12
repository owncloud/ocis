# GroupApi

All URIs are relative to *https://ocis.ocis.rolling.owncloud.works/graph*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**addMember**](GroupApi.md#addmember) | **POST** /v1.0/groups/{group-id}/members/$ref | Add a member to a group |
| [**deleteGroup**](GroupApi.md#deletegroup) | **DELETE** /v1.0/groups/{group-id} | Delete entity from groups |
| [**deleteMember**](GroupApi.md#deletemember) | **DELETE** /v1.0/groups/{group-id}/members/{directory-object-id}/$ref | Delete member from a group |
| [**getGroup**](GroupApi.md#getgroup) | **GET** /v1.0/groups/{group-id} | Get entity from groups by key |
| [**listMembers**](GroupApi.md#listmembers) | **GET** /v1.0/groups/{group-id}/members | Get a list of the group\&#39;s direct members |
| [**updateGroup**](GroupApi.md#updategroup) | **PATCH** /v1.0/groups/{group-id} | Update entity in groups |



## addMember

> addMember(groupId, memberReference)

Add a member to a group

### Example

```ts
import {
  Configuration,
  GroupApi,
} from '';
import type { AddMemberRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new GroupApi(config);

  const body = {
    // string | key: id of group
    groupId: groupId_example,
    // MemberReference | Object to be added as member
    memberReference: ...,
  } satisfies AddMemberRequest;

  try {
    const data = await api.addMember(body);
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
| **groupId** | `string` | key: id of group | [Defaults to `undefined`] |
| **memberReference** | [MemberReference](MemberReference.md) | Object to be added as member | |

### Return type

`void` (Empty response body)

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **204** | Success |  -  |
| **0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## deleteGroup

> deleteGroup(groupId, ifMatch)

Delete entity from groups

### Example

```ts
import {
  Configuration,
  GroupApi,
} from '';
import type { DeleteGroupRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new GroupApi(config);

  const body = {
    // string | key: id of group
    groupId: groupId_example,
    // string | ETag (optional)
    ifMatch: ifMatch_example,
  } satisfies DeleteGroupRequest;

  try {
    const data = await api.deleteGroup(body);
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
| **groupId** | `string` | key: id of group | [Defaults to `undefined`] |
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


## deleteMember

> deleteMember(groupId, directoryObjectId, ifMatch)

Delete member from a group

### Example

```ts
import {
  Configuration,
  GroupApi,
} from '';
import type { DeleteMemberRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new GroupApi(config);

  const body = {
    // string | key: id of group
    groupId: groupId_example,
    // string | key: id of group member to remove
    directoryObjectId: directoryObjectId_example,
    // string | ETag (optional)
    ifMatch: ifMatch_example,
  } satisfies DeleteMemberRequest;

  try {
    const data = await api.deleteMember(body);
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
| **groupId** | `string` | key: id of group | [Defaults to `undefined`] |
| **directoryObjectId** | `string` | key: id of group member to remove | [Defaults to `undefined`] |
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


## getGroup

> Group getGroup(groupId, $select, $expand)

Get entity from groups by key

### Example

```ts
import {
  Configuration,
  GroupApi,
} from '';
import type { GetGroupRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new GroupApi(config);

  const body = {
    // string | key: id or name of group
    groupId: groupId_example,
    // Set<'id' | 'description' | 'displayName' | 'members'> | Select properties to be returned (optional)
    $select: ...,
    // Set<'members'> | Expand related entities (optional)
    $expand: ...,
  } satisfies GetGroupRequest;

  try {
    const data = await api.getGroup(body);
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
| **groupId** | `string` | key: id or name of group | [Defaults to `undefined`] |
| **$select** | `id`, `description`, `displayName`, `members` | Select properties to be returned | [Optional] [Enum: id, description, displayName, members] |
| **$expand** | `members` | Expand related entities | [Optional] [Enum: members] |

### Return type

[**Group**](Group.md)

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Retrieved entity |  -  |
| **0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## listMembers

> CollectionOfUsers listMembers(groupId)

Get a list of the group\&#39;s direct members

### Example

```ts
import {
  Configuration,
  GroupApi,
} from '';
import type { ListMembersRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new GroupApi(config);

  const body = {
    // string | key: id or name of group
    groupId: 86948e45-96a6-43df-b83d-46e92afd30de,
  } satisfies ListMembersRequest;

  try {
    const data = await api.listMembers(body);
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
| **groupId** | `string` | key: id or name of group | [Defaults to `undefined`] |

### Return type

[**CollectionOfUsers**](CollectionOfUsers.md)

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Retrieved group members |  -  |
| **0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## updateGroup

> updateGroup(groupId, group)

Update entity in groups

### Example

```ts
import {
  Configuration,
  GroupApi,
} from '';
import type { UpdateGroupRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new GroupApi(config);

  const body = {
    // string | key: id of group
    groupId: groupId_example,
    // Group | New property values
    group: {"displayName":"GroupName"},
  } satisfies UpdateGroupRequest;

  try {
    const data = await api.updateGroup(body);
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
| **groupId** | `string` | key: id of group | [Defaults to `undefined`] |
| **group** | [Group](Group.md) | New property values | |

### Return type

`void` (Empty response body)

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **204** | Success |  -  |
| **0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

