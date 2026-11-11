Bugfix: fixed low color contrast on disabled danger button text

Disabled filled danger buttons, such as the "Empty trash bin" button
in the deleted files view, displayed their text and icon in a color
too close to the background fill, making them hard to read.

The disabled text and icon color is now forced to white, restoring
readable contrast against the button's background. The vault light
theme's muted danger color was also darkened, since it was too light
to meet contrast requirements against white text.

https://github.com/owncloud/ocis/pull/12998
