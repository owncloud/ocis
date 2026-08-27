Bugfix: Apply renewed access token in Web embed mode with delegated authentication

In embed mode with delegated authentication the host application renews the
access token by posting an "owncloud-embed:update-token" message into the
iframe. The message listener was registered as an unbound method reference, so
it crashed on every incoming message before the renewed token could be applied.
It also passed the whole message payload to the user context update instead of
the access token string. The embedded Web instance therefore kept using the
initial access token until it expired and the user got logged out, even though
the host had delivered fresh tokens in time.

https://github.com/owncloud/ocis/pull/12856
