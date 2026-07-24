Bugfix: Stop breadcrumb items from being announced twice

Each breadcrumb list item had `tabindex="0"`, putting the item itself in the
keyboard/screen reader focus order in addition to the link or button nested
inside it. This caused every breadcrumb segment to be announced twice in a
row. The list item is no longer a separate focus stop; only its inner link
or button is, while the item remains focusable enough for existing
drag-and-drop styling to keep working.

https://github.com/owncloud/ocis/pull/12643
