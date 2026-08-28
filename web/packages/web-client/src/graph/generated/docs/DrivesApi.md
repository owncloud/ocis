# DrivesApi

All URIs are relative to *https://ocis.ocis.rolling.owncloud.works/graph*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**createDrive**](DrivesApi.md#createdrive) | **POST** /v1.0/drives | Create a new drive of a specific type |
| [**createDriveBeta**](DrivesApi.md#createdrivebeta) | **POST** /v1beta1/drives | Create a new drive of a specific type. Alias for \&#39;/v1.0/drives\&#39;, the difference is that grantedtoV2 is used and roles contain unified roles instead of cs3 roles. |
| [**deleteDrive**](DrivesApi.md#deletedrive) | **DELETE** /v1.0/drives/{drive-id} | Delete a specific space |
| [**deleteDriveBeta**](DrivesApi.md#deletedrivebeta) | **DELETE** /v1beta1/drives/{drive-id} | Delete a specific space. Alias for \&#39;/v1.0/drives\&#39;. |
| [**getDrive**](DrivesApi.md#getdrive) | **GET** /v1.0/drives/{drive-id} | Get drive by id |
| [**getDriveBeta**](DrivesApi.md#getdrivebeta) | **GET** /v1beta1/drives/{drive-id} | Get drive by id. Alias for \&#39;/v1.0/drives\&#39;, the difference is that grantedtoV2 is used and roles contain unified roles instead of cs3 roles |
| [**updateDrive**](DrivesApi.md#updatedrive) | **PATCH** /v1.0/drives/{drive-id} | Update the drive |
| [**updateDriveBeta**](DrivesApi.md#updatedrivebeta) | **PATCH** /v1beta1/drives/{drive-id} | Update the drive. Alias for \&#39;/v1.0/drives\&#39;, the difference is that grantedtoV2 is used and roles contain unified roles instead of cs3 roles |



## createDrive

> Drive createDrive(drive)

Create a new drive of a specific type

### Example

```ts
import {
  Configuration,
  DrivesApi,
} from '';
import type { CreateDriveRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new DrivesApi(config);

  const body = {
    // Drive | New space property values
    drive: {"name":"Mars","quota":{"total":1000000000},"description":"Team space mars project"},
  } satisfies CreateDriveRequest;

  try {
    const data = await api.createDrive(body);
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
| **drive** | [Drive](Drive.md) | New space property values | |

### Return type

[**Drive**](Drive.md)

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **201** | Created |  -  |
| **0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## createDriveBeta

> Drive createDriveBeta(drive)

Create a new drive of a specific type. Alias for \&#39;/v1.0/drives\&#39;, the difference is that grantedtoV2 is used and roles contain unified roles instead of cs3 roles.

### Example

```ts
import {
  Configuration,
  DrivesApi,
} from '';
import type { CreateDriveBetaRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new DrivesApi(config);

  const body = {
    // Drive | New space property values
    drive: {"name":"Mars","quota":{"total":1000000000},"description":"Team space mars project"},
  } satisfies CreateDriveBetaRequest;

  try {
    const data = await api.createDriveBeta(body);
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
| **drive** | [Drive](Drive.md) | New space property values | |

### Return type

[**Drive**](Drive.md)

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **201** | Created |  -  |
| **0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## deleteDrive

> deleteDrive(driveId, ifMatch)

Delete a specific space

### Example

```ts
import {
  Configuration,
  DrivesApi,
} from '';
import type { DeleteDriveRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new DrivesApi(config);

  const body = {
    // string | key: id of drive
    driveId: driveId_example,
    // string | ETag (optional)
    ifMatch: ifMatch_example,
  } satisfies DeleteDriveRequest;

  try {
    const data = await api.deleteDrive(body);
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
| **driveId** | `string` | key: id of drive | [Defaults to `undefined`] |
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


## deleteDriveBeta

> deleteDriveBeta(driveId, ifMatch)

Delete a specific space. Alias for \&#39;/v1.0/drives\&#39;.

### Example

```ts
import {
  Configuration,
  DrivesApi,
} from '';
import type { DeleteDriveBetaRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new DrivesApi(config);

  const body = {
    // string | key: id of drive
    driveId: driveId_example,
    // string | ETag (optional)
    ifMatch: ifMatch_example,
  } satisfies DeleteDriveBetaRequest;

  try {
    const data = await api.deleteDriveBeta(body);
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
| **driveId** | `string` | key: id of drive | [Defaults to `undefined`] |
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


## getDrive

> Drive getDrive(driveId)

Get drive by id

### Example

```ts
import {
  Configuration,
  DrivesApi,
} from '';
import type { GetDriveRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new DrivesApi(config);

  const body = {
    // string | key: id of drive
    driveId: driveId_example,
  } satisfies GetDriveRequest;

  try {
    const data = await api.getDrive(body);
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
| **driveId** | `string` | key: id of drive | [Defaults to `undefined`] |

### Return type

[**Drive**](Drive.md)

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Retrieved drive |  -  |
| **0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## getDriveBeta

> Drive getDriveBeta(driveId)

Get drive by id. Alias for \&#39;/v1.0/drives\&#39;, the difference is that grantedtoV2 is used and roles contain unified roles instead of cs3 roles

### Example

```ts
import {
  Configuration,
  DrivesApi,
} from '';
import type { GetDriveBetaRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new DrivesApi(config);

  const body = {
    // string | key: id of drive
    driveId: driveId_example,
  } satisfies GetDriveBetaRequest;

  try {
    const data = await api.getDriveBeta(body);
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
| **driveId** | `string` | key: id of drive | [Defaults to `undefined`] |

### Return type

[**Drive**](Drive.md)

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Retrieved drive |  -  |
| **0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## updateDrive

> Drive updateDrive(driveId, driveUpdate)

Update the drive

### Example

```ts
import {
  Configuration,
  DrivesApi,
} from '';
import type { UpdateDriveRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new DrivesApi(config);

  const body = {
    // string | key: id of drive
    driveId: driveId_example,
    // DriveUpdate | New space values
    driveUpdate: {"quota":{"total":1000000000}},
  } satisfies UpdateDriveRequest;

  try {
    const data = await api.updateDrive(body);
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
| **driveId** | `string` | key: id of drive | [Defaults to `undefined`] |
| **driveUpdate** | [DriveUpdate](DriveUpdate.md) | New space values | |

### Return type

[**Drive**](Drive.md)

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Success |  -  |
| **0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## updateDriveBeta

> Drive updateDriveBeta(driveId, driveUpdate)

Update the drive. Alias for \&#39;/v1.0/drives\&#39;, the difference is that grantedtoV2 is used and roles contain unified roles instead of cs3 roles

### Example

```ts
import {
  Configuration,
  DrivesApi,
} from '';
import type { UpdateDriveBetaRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new DrivesApi(config);

  const body = {
    // string | key: id of drive
    driveId: driveId_example,
    // DriveUpdate | New space values
    driveUpdate: {"quota":{"total":1000000000}},
  } satisfies UpdateDriveBetaRequest;

  try {
    const data = await api.updateDriveBeta(body);
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
| **driveId** | `string` | key: id of drive | [Defaults to `undefined`] |
| **driveUpdate** | [DriveUpdate](DriveUpdate.md) | New space values | |

### Return type

[**Drive**](Drive.md)

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Success |  -  |
| **0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

