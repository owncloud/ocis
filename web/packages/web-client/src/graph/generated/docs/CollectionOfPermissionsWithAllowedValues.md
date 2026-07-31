# CollectionOfPermissionsWithAllowedValues


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**libre_graph_permissions_roles_allowedValues** | [**Array&lt;UnifiedRoleDefinition&gt;**](UnifiedRoleDefinition.md) | A list of role definitions that can be chosen for the resource. | [optional] [default to undefined]
**libre_graph_permissions_actions_allowedValues** | **Array&lt;string&gt;** | A list of actions that can be chosen for a custom role.  Following the CS3 API we can represent the CS3 permissions by mapping them to driveItem properties or relations like this: | [CS3 ResourcePermission](https://cs3org.github.io/cs3apis/#cs3.storage.provider.v1beta1.ResourcePermissions) | action | comment | | ------------------------------------------------------------------------------------------------------------ | ------ | ------- | | &#x60;stat&#x60; | &#x60;libre.graph/driveItem/basic/read&#x60; | &#x60;basic&#x60; because it does not include versions or trashed items | | &#x60;get_quota&#x60; | &#x60;libre.graph/driveItem/quota/read&#x60; | read only the &#x60;quota&#x60; property | | &#x60;get_path&#x60; | &#x60;libre.graph/driveItem/path/read&#x60; | read only the &#x60;path&#x60; property | | &#x60;move&#x60; | &#x60;libre.graph/driveItem/path/update&#x60; | allows updating the &#x60;path&#x60; property of a CS3 resource | | &#x60;delete&#x60; | &#x60;libre.graph/driveItem/standard/delete&#x60; | &#x60;standard&#x60; because deleting is a common update operation | | &#x60;list_container&#x60; | &#x60;libre.graph/driveItem/children/read&#x60; | | | &#x60;create_container&#x60; | &#x60;libre.graph/driveItem/children/create&#x60; | | | &#x60;initiate_file_download&#x60; | &#x60;libre.graph/driveItem/content/read&#x60; | &#x60;content&#x60; is the property read when initiating a download | | &#x60;initiate_file_upload&#x60; | &#x60;libre.graph/driveItem/upload/create&#x60; | &#x60;uploads&#x60; are a separate property. postprocessing creates the &#x60;content&#x60; | | &#x60;add_grant&#x60; | &#x60;libre.graph/driveItem/permissions/create&#x60; | | | &#x60;list_grant&#x60; | &#x60;libre.graph/driveItem/permissions/read&#x60; | | | &#x60;update_grant&#x60; | &#x60;libre.graph/driveItem/permissions/update&#x60; | | | &#x60;remove_grant&#x60; | &#x60;libre.graph/driveItem/permissions/delete&#x60; | | | &#x60;deny_grant&#x60; | &#x60;libre.graph/driveItem/permissions/deny&#x60; | uses a non CRUD action &#x60;deny&#x60; | | &#x60;list_file_versions&#x60; | &#x60;libre.graph/driveItem/versions/read&#x60; | &#x60;versions&#x60; is a &#x60;driveItemVersion&#x60; collection | | &#x60;restore_file_version&#x60; | &#x60;libre.graph/driveItem/versions/update&#x60; | the only &#x60;update&#x60; action is restore | | &#x60;list_recycle&#x60; | &#x60;libre.graph/driveItem/deleted/read&#x60; | reading a driveItem &#x60;deleted&#x60; property implies listing | | &#x60;restore_recycle_item&#x60; | &#x60;libre.graph/driveItem/deleted/update&#x60; | the only &#x60;update&#x60; action is restore | | &#x60;purge_recycle&#x60; | &#x60;libre.graph/driveItem/deleted/delete&#x60; | allows purging deleted &#x60;driveItems&#x60; |  | [optional] [default to undefined]
**value** | [**Array&lt;Permission&gt;**](Permission.md) |  | [optional] [default to undefined]

## Example

```typescript
import { CollectionOfPermissionsWithAllowedValues } from './api';

const instance: CollectionOfPermissionsWithAllowedValues = {
    libre_graph_permissions_roles_allowedValues,
    libre_graph_permissions_actions_allowedValues,
    value,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
