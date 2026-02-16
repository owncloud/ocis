Bugfix: Fix extended env parser

The extended envvar parser would be angry if there are two `os.Getenv` in the same line.
We fixed this.

https://github.com/owncloud/ocis/pull/8409
