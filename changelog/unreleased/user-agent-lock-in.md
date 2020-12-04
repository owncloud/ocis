Enhancement: Add www-authenticate based on user agent

We now comply with HTTP spec by adding Www-Authenticate headers on every `401` request. Furthermore we not only take care of such thing at the Proxy but also Reva will take care of it. In addition we now are able to lock-in user-agents to specific challenges.

https://github.com/owncloud/ocis/pull/1009
