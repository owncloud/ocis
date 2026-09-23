Bugfix: fixed low color contrast on nav and primary button hover states

Non-active sidebar navigation items, as well as buttons using the
"primary" styling, showed their text and icons in black on top of a
dark hover background in some themes, making them hard to read while
hovering.

The hover state for these elements now switches to a light color
instead of always using black, restoring readable contrast between
the text and the background while hovering.

https://github.com/owncloud/ocis/pull/TODO
