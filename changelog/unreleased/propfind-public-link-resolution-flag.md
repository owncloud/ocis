Enhancement: The public link resolution flag added to PROPFIND

The public link resolution flag and the context timeout added to PROPFIND request
The ocdav PROPFIND handler resolves public link shares to populate the oc:share-type property. 
A new toggle for skipping that lookup. The toggle is now exposed as
OCDAV_DISABLE_PROPFIND_PUBLIC_LINK_RESOLUTION, which can be set to reduce load
on services for large collections.
Also, the bounded context was added for least public share request. Since this property is needed only
to display an icon next to files in the list, indicating that a public link exists, 
it can be omitted if the ListPublicShares request is slow.

https://github.com/owncloud/ocis/pull/13031
