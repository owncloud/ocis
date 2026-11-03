Bugfix: Show the quota error when an upload fails

When an upload failed after its data had been transferred, for example because
the quota of the space was exceeded, the upload overlay stayed on "Finalizing
upload..." and the reason was only visible under "Show details". The overlay
now switches to "Upload failed" and names the exceeded quota directly.

https://github.com/owncloud/ocis/pull/13297
