# ActivitiesApi

All URIs are relative to *https://ocis.ocis.rolling.owncloud.works/graph*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**getActivities**](ActivitiesApi.md#getactivities) | **GET** /v1beta1/extensions/org.libregraph/activities | Get activities |



## getActivities

> CollectionOfActivities getActivities(kql)

Get activities

### Example

```ts
import {
  Configuration,
  ActivitiesApi,
} from '';
import type { GetActivitiesRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new ActivitiesApi(config);

  const body = {
    // string (optional)
    kql: resourceid:a0ca6a90-a365-4782-871e-d44447bbc668$a0ca6a90-a365-4782-871e-d44447bbc668 depth:2,
  } satisfies GetActivitiesRequest;

  try {
    const data = await api.getActivities(body);
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
| **kql** | `string` |  | [Optional] [Defaults to `undefined`] |

### Return type

[**CollectionOfActivities**](CollectionOfActivities.md)

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Found activities |  -  |
| **0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

