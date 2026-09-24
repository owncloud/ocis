Bugfix: fixed low color contrast on share permission dropdown hover text

The share permission dropdown's role options, such as "Can view" or
"Can edit", showed their hovered text in a color too close in
lightness to the hover background in the dark theme, making it hard
to read.

The hover text color has been lightened slightly, restoring
readable contrast between the text and the background while
hovering.

https://github.com/owncloud/ocis/pull/12995
