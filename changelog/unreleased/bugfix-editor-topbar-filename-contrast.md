Bugfix: fixed low color contrast on file name in editor topbar

The file name and extension shown in the topbar when a file is open,
for example in an editor, were displayed in white on top of a light
brand-colored background in some theme variants, making the text
hard to read.

The text now uses a dark color instead, restoring readable contrast
against that background.

https://github.com/owncloud/ocis/pull/12997
