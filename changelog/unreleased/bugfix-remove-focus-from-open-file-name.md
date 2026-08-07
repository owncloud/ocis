Bugfix: Remove keyboard focus from open file name in top bar

When opening a file such as a text or markdown document, the file name shown
in the top bar could be reached via keyboard tabbing even though it was not
clickable and triggered no action. This made keyboard navigation confusing.

The file name is no longer part of the tab order, since it is a purely
informational, non-interactive element.

https://github.com/owncloud/ocis/pull/12752
