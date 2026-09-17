Bugfix: Truncate long user names in the admin settings users list

In the admin settings users list, a long user name or display name would
overflow its table column, pushing the columns after it, including the
action buttons, out of place or off screen.

Long values in the user name and display name columns are now truncated
with an ellipsis so the table layout stays intact regardless of value
length.

https://github.com/owncloud/ocis/pull/12816
