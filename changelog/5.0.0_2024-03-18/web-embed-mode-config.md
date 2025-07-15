Enhancement: Configs for Web embed mode

New configs for the Web embed mode have been added:

* `enabled` Defines if embed mode is enabled.
* `target` Defines how Web is being integrated when running in embed mode.
* `messagesOrigin` Defines a URL under which Web can be integrated via iFrame.
* `delegateAuthentication` Defines whether Web should require authentication to be done by the parent application.
* `delegateAuthenticationOrigin` Defines the host to validate the message event origin against when running Web in 'embed' mode.

https://github.com/owncloud/ocis/pull/7670
https://github.com/owncloud/web/issues/9768
