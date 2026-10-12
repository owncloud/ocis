Bugfix: Purge revisions stored with whole-second timestamps

The revisions purge command reported a clean storage when revisions were still
present. Revision filenames are the previous mtime formatted with
time.RFC3339Nano, which omits fractional seconds on a whole second. Names such
as `.REV.1601-01-01T00:00:00Z` never matched the purge pattern, so those
revisions stayed on disk.

The matcher now treats the decimal point and fractional digits as optional.
Whole-second revisions and their `.mpk` and `.mlock` companions are discovered
and purged along with revisions that still carry a fractional timestamp.

https://github.com/owncloud/ocis/issues/11167
