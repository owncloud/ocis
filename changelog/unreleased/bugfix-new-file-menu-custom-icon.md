Bugfix: Show the app icon in the new file menu

An app that adds an entry to the "New" file menu without a file extension,
for example to start an assistant, can set its own icon. The menu ignored that
icon and always showed the generic file icon. The menu now shows the icon the
app sets, with the same fill type as the app's "Open in" action.

https://github.com/owncloud/ocis/pull/13566
https://github.com/owncloud/ocis/issues/12659
