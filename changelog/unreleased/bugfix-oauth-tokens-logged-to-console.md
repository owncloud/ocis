Bugfix: Stop logging OAuth tokens to the browser console

On every login and every silent token renewal, the web client wrote the
signed-in user's raw OAuth `access_token` and `refresh_token` to the browser
console via a `console.debug` call in the `userLoaded` handler of the auth
service. Anyone with access to that browser's console or log buffer - a
malicious or over-permissioned browser extension, a shared or kiosk terminal,
a screen share or recording, or a session-replay tool that captures console
output - could read and reuse those live credentials. The production build
does not strip `console.debug`, so this shipped in the deployed bundle, not
just local dev.

The debug line now confirms that a new user was loaded without logging the
token, profile, or any other user field.

https://github.com/owncloud/ocis/pull/13233
