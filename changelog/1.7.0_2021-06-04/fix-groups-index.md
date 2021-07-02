Bugfix: Change the groups index to be case sensitive 

Groups are considered to be case sensitive. The index must handle them case sensitive too otherwise we will have undeterministic behavior while editing or deleting groups.

https://github.com/owncloud/ocis/pull/2109
