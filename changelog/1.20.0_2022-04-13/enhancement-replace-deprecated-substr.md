Enhancement: Replace deprecated String.prototype.substr()

We've replaced all occurrences of the deprecated String.prototype.substr()
function with String.prototype.slice() which works similarly but isn't
deprecated.

https://github.com/owncloud/ocis/pull/3448
