Bugfix: Vault share search could show ineligible users

Searching for a share recipient by a cross-instance reference (e.g.
`user@school`) could return a user without vault access, letting them be
selected in the Safe share picker before failing with an error afterward.
Vault eligibility is now always enforced for this kind of search too.

https://github.com/owncloud/ocis/pull/13154
