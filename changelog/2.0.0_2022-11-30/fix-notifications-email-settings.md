Bugfix: Fix notifications service settings

We've fixed two notifications service setting:
- `NOTIFICATIONS_MACHINE_AUTH_API_KEY` was previously not picked up (only `OCIS_MACHINE_AUTH_API_KEY` was loaded)
- If you used a email sender address in the format of the default value of `NOTIFICATIONS_SMTP_SENDER` no email could be send.

https://github.com/owncloud/ocis/pull/4652
