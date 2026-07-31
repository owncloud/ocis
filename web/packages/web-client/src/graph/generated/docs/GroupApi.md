# GroupApi

All URIs are relative to *https://ocis.ocis.rolling.owncloud.works/graph*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**addMember**](#addmember) | **POST** /v1.0/groups/{group-id}/members/$ref | Add a member to a group|
|[**deleteGroup**](#deletegroup) | **DELETE** /v1.0/groups/{group-id} | Delete entity from groups|
|[**deleteMember**](#deletemember) | **DELETE** /v1.0/groups/{group-id}/members/{directory-object-id}/$ref | Delete member from a group|
|[**getGroup**](#getgroup) | **GET** /v1.0/groups/{group-id} | Get entity from groups by key|
|[**listMembers**](#listmembers) | **GET** /v1.0/groups/{group-id}/members | Get a list of the group\&#39;s direct members|
|[**updateGroup**](#updategroup) | **PATCH** /v1.0/groups/{group-id} | Update entity in groups|

# **addMember**
> addMember(memberReference)


### Example

```typescript
import {
    GroupApi,
    Configuration,
    MemberReference
} from './api';

const configuration = new Configuration();
const apiInstance = new GroupApi(configuration);

let groupId: string; //key: id of group (default to undefined)
let memberReference: MemberReference; //Object to be added as member

const { status, data } = await apiInstance.addMember(
    groupId,
    memberReference
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **memberReference** | **MemberReference**| Object to be added as member | |
| **groupId** | [**string**] | key: id of group | defaults to undefined|


### Return type

void (empty response body)

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**204** | Success |  -  |
|**0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deleteGroup**
> deleteGroup()


### Example

```typescript
import {
    GroupApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new GroupApi(configuration);

let groupId: string; //key: id of group (default to undefined)
let ifMatch: string; //ETag (optional) (default to undefined)

const { status, data } = await apiInstance.deleteGroup(
    groupId,
    ifMatch
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **groupId** | [**string**] | key: id of group | defaults to undefined|
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

# **deleteMember**
> deleteMember()


### Example

```typescript
import {
    GroupApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new GroupApi(configuration);

let groupId: string; //key: id of group (default to undefined)
let directoryObjectId: string; //key: id of group member to remove (default to undefined)
let ifMatch: string; //ETag (optional) (default to undefined)

const { status, data } = await apiInstance.deleteMember(
    groupId,
    directoryObjectId,
    ifMatch
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **groupId** | [**string**] | key: id of group | defaults to undefined|
| **directoryObjectId** | [**string**] | key: id of group member to remove | defaults to undefined|
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

# **getGroup**
> Group getGroup()


### Example

```typescript
import {
    GroupApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new GroupApi(configuration);

let groupId: string; //key: id or name of group (default to undefined)
let $select: Set<'id' | 'description' | 'displayName' | 'members'>; //Select properties to be returned (optional) (default to undefined)
let $expand: Set<'members'>; //Expand related entities (optional) (default to undefined)

const { status, data } = await apiInstance.getGroup(
    groupId,
    $select,
    $expand
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **groupId** | [**string**] | key: id or name of group | defaults to undefined|
| **$select** | **Array<&#39;id&#39; &#124; &#39;description&#39; &#124; &#39;displayName&#39; &#124; &#39;members&#39;>** | Select properties to be returned | (optional) defaults to undefined|
| **$expand** | **Array<&#39;members&#39;>** | Expand related entities | (optional) defaults to undefined|


### Return type

**Group**

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Retrieved entity |  -  |
|**0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **listMembers**
> CollectionOfUsers listMembers()


### Example

```typescript
import {
    GroupApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new GroupApi(configuration);

let groupId: string; //key: id or name of group (default to undefined)

const { status, data } = await apiInstance.listMembers(
    groupId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **groupId** | [**string**] | key: id or name of group | defaults to undefined|


### Return type

**CollectionOfUsers**

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Retrieved group members |  -  |
|**0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateGroup**
> updateGroup(group)


### Example

```typescript
import {
    GroupApi,
    Configuration,
    Group
} from './api';

const configuration = new Configuration();
const apiInstance = new GroupApi(configuration);

let groupId: string; //key: id of group (default to undefined)
let group: Group; //New property values

const { status, data } = await apiInstance.updateGroup(
    groupId,
    group
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **group** | **Group**| New property values | |
| **groupId** | [**string**] | key: id of group | defaults to undefined|


### Return type

void (empty response body)

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**204** | Success |  -  |
|**0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

