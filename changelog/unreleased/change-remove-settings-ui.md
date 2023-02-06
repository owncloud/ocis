Change: Remove the settings ui

With ownCloud Web having transitioned to Vue 3 recently, we would have had to port the settings ui as well. The decision was made
to discontinue the settings ui instead. As a result all traces of the settings ui have been removed.

The only user facing setting that ever existed in the settings service is now integrated into the `account` page of
ownCloud Web (click on top right user menu, then on your username to reach the account page).

https://github.com/owncloud/ocis/pull/5463
