Bugfix: Wrap long values in share and user details views

In the "Access details" popup of a share, and in the admin settings user
details panel, a long user name, display name, email or domain would
overflow the popup or sidebar instead of wrapping, hiding parts of the
value and breaking the layout.

Long values are now wrapped onto multiple lines so these views stay
within their bounds regardless of value length.

https://github.com/owncloud/ocis/pull/12840
