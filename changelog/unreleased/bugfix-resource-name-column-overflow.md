Bugfix: Truncate long resource names instead of overlapping the next column

In resource tables, a name that was too long for its column was not truncated.
Instead it kept its full width and was painted over the value of the adjacent
column, making both unreadable. This was most visible in the Spaces overview,
where a long space name overlapped the "Manager" column, but it affected every
resource table using the name column, including the files list and search
results, as well as the parent folder path shown below a search result.

The name is now truncated with an ellipsis at the column boundary, and the
rename button stays inside the name column. The full name remains available via
the title tooltip on hover.

https://github.com/owncloud/ocis/pull/12754
