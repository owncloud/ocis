# TagsApi

All URIs are relative to *https://ocis.ocis.rolling.owncloud.works/graph*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**assignTags**](#assigntags) | **PUT** /v1.0/extensions/org.libregraph/tags | Assign tags to a resource|
|[**getTags**](#gettags) | **GET** /v1.0/extensions/org.libregraph/tags | Get all known tags|
|[**unassignTags**](#unassigntags) | **DELETE** /v1.0/extensions/org.libregraph/tags | Unassign tags from a resource|

# **assignTags**
> assignTags()


### Example

```typescript
import {
    TagsApi,
    Configuration,
    TagAssignment
} from './api';

const configuration = new Configuration();
const apiInstance = new TagsApi(configuration);

let tagAssignment: TagAssignment; // (optional)

const { status, data } = await apiInstance.assignTags(
    tagAssignment
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **tagAssignment** | **TagAssignment**|  | |


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
|**200** | No content |  -  |
|**0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getTags**
> CollectionOfTags getTags()


### Example

```typescript
import {
    TagsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new TagsApi(configuration);

const { status, data } = await apiInstance.getTags();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**CollectionOfTags**

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Retrieved tags |  -  |
|**0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **unassignTags**
> unassignTags()


### Example

```typescript
import {
    TagsApi,
    Configuration,
    TagUnassignment
} from './api';

const configuration = new Configuration();
const apiInstance = new TagsApi(configuration);

let tagUnassignment: TagUnassignment; // (optional)

const { status, data } = await apiInstance.unassignTags(
    tagUnassignment
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **tagUnassignment** | **TagUnassignment**|  | |


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
|**200** | No content |  -  |
|**0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

