Bugfix: Rate-limit the exportPersonalData endpoint

We've added a rate limit to the `exportPersonalData`
endpoint to mitigate an authenticated application-level denial-of-service.
Rate-limit per endpoint path carries the userID, so effectively per user.
The endpoint is now limited to 5 requests per minute.

https://github.com/owncloud/ocis/issues/12516
