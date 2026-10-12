Bugfix: Open files whose name contains a hash or question mark

Cmd+Clicking or middle-clicking a file with a `#` or `?` in its name opened the
space root instead of the file, and the top bar showed the space name rather than
the file name. The file's path is carried in a route parameter which the router
percent-encodes when it builds the URL, but the editor route was resolved twice
and resolving twice does not encode twice: the second pass rebuilt the URL from
the still-decoded parameter with the encoding step skipped. The browser then read
everything from the `#` onward as a fragment, discarding the query string along
with the `fileId` needed to identify the file. Left clicks were unaffected because
they never leave the app.

Editor routes are now handed to links unresolved, so the file name is encoded
exactly once and survives both a new tab and an in-app navigation. Hash and
question mark remain valid characters in file names.

https://github.com/owncloud/ocis/pull/12740
