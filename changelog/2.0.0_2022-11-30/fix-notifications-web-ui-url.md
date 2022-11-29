Bugfix: Fix notifications Web UI url

We've fixed the configuration of the notification service's Web UI url that appears in emails.

Previously it was only configurable via the global "OCIS_URL" and is now also configurable via "NOTIFICATIONS_WEB_UI_URL".

https://github.com/owncloud/ocis/pull/4998
