Bugfix: Explain why an upload or a duplicate does not fit the space quota

An upload rejected for lack of room named only the shortfall ("You need
additional 1 MB"), and a failed duplicate was reported as a failed move with
"Insufficient quota" underneath. Both messages now name the space and what to do
about it, and a failed duplicate is reported as one.

The upload check also measured only the last file of a batch, and skipped the
check altogether for a space with nothing left, because a remaining quota of 0
was read as "unknown".

https://github.com/owncloud/ocis/pull/12943
