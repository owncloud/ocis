Bugfix: fixed low color contrast on disabled danger buttons

Disabled buttons using the "danger" styling, such as the "Empty trash bin"
button when there is nothing to delete, showed their label text with very
low contrast against the button background, making the text hard to read.

The muted background color used for these buttons in a disabled state has
been darkened, and the fading effect applied to disabled buttons no longer
washes out the danger button background, restoring readable contrast
between the label and the background.

https://github.com/owncloud/ocis/pull/TODO
