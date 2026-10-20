Bugfix: Keep /me working when the preferred language lookup fails

The graph service returned HTTP 500 for the whole `GET /graph/v1.0/me` response
whenever the settings lookup for the user's preferred language failed, including
on a transient timeout. Since `/me` is the first request the Web UI makes after
login, that locked affected users out of oCIS entirely. The lookup failure is
now logged at warn level and `preferredLanguage` falls back to
`OCIS_DEFAULT_LANGUAGE`.

https://github.com/owncloud/ocis/pull/13188
