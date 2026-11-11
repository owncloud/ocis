Enhancement: Add "Undo" action to the delete success notification

Deleting a file or folder now shows an "Undo" action alongside the
"moved to trash bin" notification, letting users restore the deleted
resources directly from the toast instead of navigating to the trash
bin. The action is only shown when the resources can be resolved in
the trash and the space supports restoring from it. The notification
pauses its auto-dismiss timer while hovered or focused, so it does not
disappear before the action can be used. Pressing Ctrl+Z (or Cmd+Z on
macOS) while the notification is visible triggers the same undo.

https://github.com/owncloud/ocis/pull/12650
