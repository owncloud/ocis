# TagsApi

All URIs are relative to *https://ocis.ocis.rolling.owncloud.works/graph*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**assignTags**](TagsApi.md#assigntags) | **PUT** /v1.0/extensions/org.libregraph/tags | Assign tags to a resource |
| [**getTags**](TagsApi.md#gettags) | **GET** /v1.0/extensions/org.libregraph/tags | Get all known tags |
| [**unassignTags**](TagsApi.md#unassigntags) | **DELETE** /v1.0/extensions/org.libregraph/tags | Unassign tags from a resource |



## assignTags

> assignTags(tagAssignment)

Assign tags to a resource

### Example

```ts
import {
  Configuration,
  TagsApi,
} from '';
import type { AssignTagsRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new TagsApi(config);

  const body = {
    // TagAssignment (optional)
    tagAssignment: ...,
  } satisfies AssignTagsRequest;

  try {
    const data = await api.assignTags(body);
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
| **tagAssignment** | [TagAssignment](TagAssignment.md) |  | [Optional] |

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
| **200** | No content |  -  |
| **0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## getTags

> CollectionOfTags getTags()

Get all known tags

### Example

```ts
import {
  Configuration,
  TagsApi,
} from '';
import type { GetTagsRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new TagsApi(config);

  try {
    const data = await api.getTags();
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

[**CollectionOfTags**](CollectionOfTags.md)

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Retrieved tags |  -  |
| **0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## unassignTags

> unassignTags(tagUnassignment)

Unassign tags from a resource

### Example

```ts
import {
  Configuration,
  TagsApi,
} from '';
import type { UnassignTagsRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new TagsApi(config);

  const body = {
    // TagUnassignment (optional)
    tagUnassignment: ...,
  } satisfies UnassignTagsRequest;

  try {
    const data = await api.unassignTags(body);
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
| **tagUnassignment** | [TagUnassignment](TagUnassignment.md) |  | [Optional] |

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
| **200** | No content |  -  |
| **0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

