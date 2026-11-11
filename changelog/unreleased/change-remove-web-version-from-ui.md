Change: Remove web version from the Web UI

The Web UI no longer displays its own version. It was previously shown in the
bottom of the left sidebar and on the public file drop page, and was also
announced in the browser console and the page metadata. This avoids exposing
internal version details in the interface. The backend version information
remains unchanged.

https://github.com/owncloud/ocis/pull/13346
