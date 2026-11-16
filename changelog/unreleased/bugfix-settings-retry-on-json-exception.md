Bugfix: Retry settings values-list fetch on truncated JSON response

Under CI load the settings API occasionally returns a truncated HTTP
response body. When `getAutoAcceptSharesDefaultValue` attempted to
parse that body, `json_decode` threw a `JsonException` with
"Control character error, possibly incorrectly encoded", causing the
`apiCollaboration/wopi.feature` acceptance test to fail spuriously on
the share-invitation step.

The method now retries the full HTTP call up to ten times with a 500 ms
interval when a `JsonException` is caught, obtaining a fresh response
each time instead of re-parsing the same corrupt bytes.
