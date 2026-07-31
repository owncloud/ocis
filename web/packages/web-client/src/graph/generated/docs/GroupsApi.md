# GroupsApi

All URIs are relative to *https://ocis.ocis.rolling.owncloud.works/graph*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**createGroup**](#creategroup) | **POST** /v1.0/groups | Add new entity to groups|
|[**listGroups**](#listgroups) | **GET** /v1.0/groups | Get entities from groups|

# **createGroup**
> Group createGroup(group)


### Example

```typescript
import {
    GroupsApi,
    Configuration,
    Group
} from './api';

const configuration = new Configuration();
const apiInstance = new GroupsApi(configuration);

let group: Group; //New entity

const { status, data } = await apiInstance.createGroup(
    group
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **group** | **Group**| New entity | |


### Return type

**Group**

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**201** | Created entity |  -  |
|**0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **listGroups**
> CollectionOfGroup listGroups()


### Example

```typescript
import {
    GroupsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new GroupsApi(configuration);

let $search: string; //Search items by search phrases (optional) (default to undefined)
let $orderby: Set<'displayName' | 'displayName desc'>; //Order items by property values (optional) (default to undefined)
let $select: Set<'id' | 'description' | 'displayName' | 'mail' | 'members'>; //Select properties to be returned (optional) (default to undefined)
let $expand: Set<'members'>; //Expand related entities (optional) (default to undefined)

const { status, data } = await apiInstance.listGroups(
    $search,
    $orderby,
    $select,
    $expand
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **$search** | [**string**] | Search items by search phrases | (optional) defaults to undefined|
| **$orderby** | **Array<&#39;displayName&#39; &#124; &#39;displayName desc&#39;>** | Order items by property values | (optional) defaults to undefined|
| **$select** | **Array<&#39;id&#39; &#124; &#39;description&#39; &#124; &#39;displayName&#39; &#124; &#39;mail&#39; &#124; &#39;members&#39;>** | Select properties to be returned | (optional) defaults to undefined|
| **$expand** | **Array<&#39;members&#39;>** | Expand related entities | (optional) defaults to undefined|


### Return type

**CollectionOfGroup**

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Retrieved entities |  -  |
|**0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

