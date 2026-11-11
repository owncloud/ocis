# ApplicationsApi

All URIs are relative to *https://ocis.ocis.rolling.owncloud.works/graph*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**getApplication**](ApplicationsApi.md#getapplication) | **GET** /v1.0/applications/{application-id} | Get application by id |
| [**listApplications**](ApplicationsApi.md#listapplications) | **GET** /v1.0/applications | Get all applications |



## getApplication

> Application getApplication(applicationId)

Get application by id

### Example

```ts
import {
  Configuration,
  ApplicationsApi,
} from '';
import type { GetApplicationRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new ApplicationsApi(config);

  const body = {
    // string | key: id of application
    applicationId: applicationId_example,
  } satisfies GetApplicationRequest;

  try {
    const data = await api.getApplication(body);
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
| **applicationId** | `string` | key: id of application | [Defaults to `undefined`] |

### Return type

[**Application**](Application.md)

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | OK |  -  |
| **0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## listApplications

> CollectionOfApplications listApplications()

Get all applications

### Example

```ts
import {
  Configuration,
  ApplicationsApi,
} from '';
import type { ListApplicationsRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new ApplicationsApi(config);

  try {
    const data = await api.listApplications();
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

[**CollectionOfApplications**](CollectionOfApplications.md)

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

