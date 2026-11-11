
# Permission

The Permission resource provides information about a sharing permission granted for a DriveItem resource.  ### Remarks  The Permission resource uses *facets* to provide information about the kind of permission represented by the resource.  Permissions with a `link` facet represent sharing links created on the item. Sharing links contain a unique token that provides access to the item for anyone with the link.  Permissions with a `invitation` facet represent permissions added by inviting specific users or groups to have access to the file. 

## Properties

Name | Type
------------ | -------------
`id` | string
`hasPassword` | boolean
`expirationDateTime` | string
`createdDateTime` | string
`grantedToV2` | [SharePointIdentitySet](SharePointIdentitySet.md)
`link` | [SharingLink](SharingLink.md)
`roles` | Array&lt;string&gt;
`grantedToIdentities` | [Array&lt;IdentitySet&gt;](IdentitySet.md)
`atLibreGraphPermissionsActions` | Array&lt;string&gt;
`invitation` | [SharingInvitation](SharingInvitation.md)

## Example

```typescript
import type { Permission } from ''

// TODO: Update the object below with actual values
const example = {
  "id": null,
  "hasPassword": null,
  "expirationDateTime": null,
  "createdDateTime": null,
  "grantedToV2": null,
  "link": null,
  "roles": null,
  "grantedToIdentities": null,
  "atLibreGraphPermissionsActions": null,
  "invitation": null,
} satisfies Permission

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as Permission
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


