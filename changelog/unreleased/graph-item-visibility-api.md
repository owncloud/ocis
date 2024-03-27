Enhancement: Ability to Change Share Item Visibility in Graph API

Introduce the `PATCH /graph/v1beta1/drives/{driveID}/items/{itemID}` Graph API endpoint which allows updating individual Drive Items.

At the moment, only the share visibility is considered changeable, but in the future, more properties can be added to this endpoint.

This enhancement is needed for the user interface, allowing specific shares to be hidden or unhidden as needed,
thereby improving the user experience.

https://github.com/owncloud/ocis/pull/8750
https://github.com/owncloud/ocis/issues/8654
