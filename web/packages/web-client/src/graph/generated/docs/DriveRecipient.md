# DriveRecipient

Represents a person, group, or other recipient to share a drive item with using the invite action.  When using invite to add permissions, the `driveRecipient` object would specify the `email`, `alias`, or `objectId` of the recipient. Only one of these values is required; multiple values are not accepted. 

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**objectId** | **string** | The unique identifier for the recipient in the directory. | [optional] [default to undefined]
**libre_graph_recipient_type** | **string** | When the recipient is referenced by objectId this annotation is used to differentiate &#x60;user&#x60; and &#x60;group&#x60; recipients. | [optional] [default to 'user']

## Example

```typescript
import { DriveRecipient } from './api';

const instance: DriveRecipient = {
    objectId,
    libre_graph_recipient_type,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
