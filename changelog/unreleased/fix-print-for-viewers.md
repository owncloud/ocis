Bugfix: Allow viewers to print office files

Users with the "Can view" role could not print, copy or export in
OnlyOffice and Collabora, although they may download the file. They now
open files in read mode. Only secure viewers keep these restrictions.

https://github.com/owncloud/ocis/pull/13525
