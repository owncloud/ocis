# DrivesApi

All URIs are relative to *https://ocis.ocis.rolling.owncloud.works/graph*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**createDrive**](#createdrive) | **POST** /v1.0/drives | Create a new drive of a specific type|
|[**createDriveBeta**](#createdrivebeta) | **POST** /v1beta1/drives | Create a new drive of a specific type. Alias for \&#39;/v1.0/drives\&#39;, the difference is that grantedtoV2 is used and roles contain unified roles instead of cs3 roles.|
|[**deleteDrive**](#deletedrive) | **DELETE** /v1.0/drives/{drive-id} | Delete a specific space|
|[**deleteDriveBeta**](#deletedrivebeta) | **DELETE** /v1beta1/drives/{drive-id} | Delete a specific space. Alias for \&#39;/v1.0/drives\&#39;.|
|[**getDrive**](#getdrive) | **GET** /v1.0/drives/{drive-id} | Get drive by id|
|[**getDriveBeta**](#getdrivebeta) | **GET** /v1beta1/drives/{drive-id} | Get drive by id. Alias for \&#39;/v1.0/drives\&#39;, the difference is that grantedtoV2 is used and roles contain unified roles instead of cs3 roles|
|[**updateDrive**](#updatedrive) | **PATCH** /v1.0/drives/{drive-id} | Update the drive|
|[**updateDriveBeta**](#updatedrivebeta) | **PATCH** /v1beta1/drives/{drive-id} | Update the drive. Alias for \&#39;/v1.0/drives\&#39;, the difference is that grantedtoV2 is used and roles contain unified roles instead of cs3 roles|

# **createDrive**
> Drive createDrive(drive)


### Example

```typescript
import {
    DrivesApi,
    Configuration,
    Drive
} from './api';

const configuration = new Configuration();
const apiInstance = new DrivesApi(configuration);

let drive: Drive; //New space property values

const { status, data } = await apiInstance.createDrive(
    drive
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **drive** | **Drive**| New space property values | |


### Return type

**Drive**

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**201** | Created |  -  |
|**0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **createDriveBeta**
> Drive createDriveBeta(drive)


### Example

```typescript
import {
    DrivesApi,
    Configuration,
    Drive
} from './api';

const configuration = new Configuration();
const apiInstance = new DrivesApi(configuration);

let drive: Drive; //New space property values

const { status, data } = await apiInstance.createDriveBeta(
    drive
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **drive** | **Drive**| New space property values | |


### Return type

**Drive**

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**201** | Created |  -  |
|**0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deleteDrive**
> deleteDrive()


### Example

```typescript
import {
    DrivesApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DrivesApi(configuration);

let driveId: string; //key: id of drive (default to undefined)
let ifMatch: string; //ETag (optional) (default to undefined)

const { status, data } = await apiInstance.deleteDrive(
    driveId,
    ifMatch
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **driveId** | [**string**] | key: id of drive | defaults to undefined|
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

# **deleteDriveBeta**
> deleteDriveBeta()


### Example

```typescript
import {
    DrivesApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DrivesApi(configuration);

let driveId: string; //key: id of drive (default to undefined)
let ifMatch: string; //ETag (optional) (default to undefined)

const { status, data } = await apiInstance.deleteDriveBeta(
    driveId,
    ifMatch
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **driveId** | [**string**] | key: id of drive | defaults to undefined|
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

# **getDrive**
> Drive getDrive()


### Example

```typescript
import {
    DrivesApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DrivesApi(configuration);

let driveId: string; //key: id of drive (default to undefined)

const { status, data } = await apiInstance.getDrive(
    driveId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **driveId** | [**string**] | key: id of drive | defaults to undefined|


### Return type

**Drive**

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Retrieved drive |  -  |
|**0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getDriveBeta**
> Drive getDriveBeta()


### Example

```typescript
import {
    DrivesApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DrivesApi(configuration);

let driveId: string; //key: id of drive (default to undefined)

const { status, data } = await apiInstance.getDriveBeta(
    driveId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **driveId** | [**string**] | key: id of drive | defaults to undefined|


### Return type

**Drive**

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Retrieved drive |  -  |
|**0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateDrive**
> Drive updateDrive(driveUpdate)


### Example

```typescript
import {
    DrivesApi,
    Configuration,
    DriveUpdate
} from './api';

const configuration = new Configuration();
const apiInstance = new DrivesApi(configuration);

let driveId: string; //key: id of drive (default to undefined)
let driveUpdate: DriveUpdate; //New space values

const { status, data } = await apiInstance.updateDrive(
    driveId,
    driveUpdate
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **driveUpdate** | **DriveUpdate**| New space values | |
| **driveId** | [**string**] | key: id of drive | defaults to undefined|


### Return type

**Drive**

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Success |  -  |
|**0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateDriveBeta**
> Drive updateDriveBeta(driveUpdate)


### Example

```typescript
import {
    DrivesApi,
    Configuration,
    DriveUpdate
} from './api';

const configuration = new Configuration();
const apiInstance = new DrivesApi(configuration);

let driveId: string; //key: id of drive (default to undefined)
let driveUpdate: DriveUpdate; //New space values

const { status, data } = await apiInstance.updateDriveBeta(
    driveId,
    driveUpdate
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **driveUpdate** | **DriveUpdate**| New space values | |
| **driveId** | [**string**] | key: id of drive | defaults to undefined|


### Return type

**Drive**

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Success |  -  |
|**0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

