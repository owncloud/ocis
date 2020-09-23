Change: environment updates for the username userid split

We updated the owncloud storage driver in reva to properly look up users by userid or username using the userprovider instead of taking the path segment as is. This requires the user service address as well as changing the default layout to the userid instead of the username. The latter is not considered a stable and persistent identifier.

<https://github.com/owncloud/ocis/ocis-revapull/420>
<https://github.com/cs3org/reva/pull/1033>
