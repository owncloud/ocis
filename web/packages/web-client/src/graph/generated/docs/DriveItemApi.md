# DriveItemApi

All URIs are relative to *https://ocis.ocis.rolling.owncloud.works/graph*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**deleteDriveItem**](DriveItemApi.md#deletedriveitem) | **DELETE** /v1beta1/drives/{drive-id}/items/{item-id} | Delete a DriveItem. |
| [**getDriveItem**](DriveItemApi.md#getdriveitem) | **GET** /v1beta1/drives/{drive-id}/items/{item-id} | Get a DriveItem. |
| [**updateDriveItem**](DriveItemApi.md#updatedriveitem) | **PATCH** /v1beta1/drives/{drive-id}/items/{item-id} | Update a DriveItem. |



## deleteDriveItem

> deleteDriveItem(driveId, itemId)

Delete a DriveItem.

Delete a DriveItem by using its ID.  Deleting items using this method moves the items to the recycle bin instead of permanently deleting the item.  Mounted shares in the share jail are unmounted. The &#x60;@client.synchronize&#x60; property of the &#x60;driveItem&#x60; in the [sharedWithMe](#/me.drive/ListSharedWithMe) endpoint will change to false. 

### Example

```ts
import {
  Configuration,
  DriveItemApi,
} from '';
import type { DeleteDriveItemRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new DriveItemApi(config);

  const body = {
    // string | key: id of drive
    driveId: a0ca6a90-a365-4782-871e-d44447bbc668$a0ca6a90-a365-4782-871e-d44447bbc668,
    // string | key: id of item
    itemId: a0ca6a90-a365-4782-871e-d44447bbc668$a0ca6a90-a365-4782-871e-d44447bbc668!share-id,
  } satisfies DeleteDriveItemRequest;

  try {
    const data = await api.deleteDriveItem(body);
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
| **itemId** | `string` | key: id of item | [Defaults to `undefined`] |

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


## getDriveItem

> DriveItem getDriveItem(driveId, itemId)

Get a DriveItem.

Get a DriveItem by using its ID. 

### Example

```ts
import {
  Configuration,
  DriveItemApi,
} from '';
import type { GetDriveItemRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new DriveItemApi(config);

  const body = {
    // string | key: id of drive
    driveId: a0ca6a90-a365-4782-871e-d44447bbc668$a0ca6a90-a365-4782-871e-d44447bbc668,
    // string | key: id of item
    itemId: a0ca6a90-a365-4782-871e-d44447bbc668$a0ca6a90-a365-4782-871e-d44447bbc668!share-id,
  } satisfies GetDriveItemRequest;

  try {
    const data = await api.getDriveItem(body);
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
| **itemId** | `string` | key: id of item | [Defaults to `undefined`] |

### Return type

[**DriveItem**](DriveItem.md)

### Authorization

[openId](../README.md#openId), [basicAuth](../README.md#basicAuth)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Retrieved driveItem |  -  |
| **0** | error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## updateDriveItem

> DriveItem updateDriveItem(driveId, itemId, driveItem)

Update a DriveItem.

Update a DriveItem.  The request body must include a JSON object with the properties to update. Only the properties that are provided will be updated.  Currently it supports updating the following properties:  * &#x60;@UI.Hidden&#x60; - Hides the item from the UI. 

### Example

```ts
import {
  Configuration,
  DriveItemApi,
} from '';
import type { UpdateDriveItemRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const config = new Configuration({ 
    // To configure HTTP basic authorization: basicAuth
    username: "YOUR USERNAME",
    password: "YOUR PASSWORD",
  });
  const api = new DriveItemApi(config);

  const body = {
    // string | key: id of drive
    driveId: a0ca6a90-a365-4782-871e-d44447bbc668$a0ca6a90-a365-4782-871e-d44447bbc668,
    // string | key: id of item
    itemId: a0ca6a90-a365-4782-871e-d44447bbc668$a0ca6a90-a365-4782-871e-d44447bbc668!share-id,
    // DriveItem | DriveItem properties to update
    driveItem: {"@UI.Hidden":true},
  } satisfies UpdateDriveItemRequest;

  try {
    const data = await api.updateDriveItem(body);
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
| **itemId** | `string` | key: id of item | [Defaults to `undefined`] |
| **driveItem** | [DriveItem](DriveItem.md) | DriveItem properties to update | |

### Return type

[**DriveItem**](DriveItem.md)

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

