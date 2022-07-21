Bugfix: Substring search for sharees

We fixed searching for sharees to be no longer case-sensitive. 
With this we introduced two new settings for the users and groups services:
"group_substring_filter_type" for the group services and
"user_substring_filter_type" for the users service.
They allow to set the type of LDAP filter that is used for substring user
searches. Possible values are: "initial", "final" and "any" to do either
prefix, suffix or full substring searches. Both settings default to "initial".

Also a new option "search_min_length" was added for the "frontend" service. It
allows to configure the minimum number of characters to enter before a search
for Sharees is started. This setting is e.g. evaluated by the web ui via the
capabilities endpoint.

https://github.com/owncloud/ocis/issues/547
