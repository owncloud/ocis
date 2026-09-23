Bugfix: fixed low color contrast on dropdown menu item hover text

Items in dropdown menus, such as the "New" file/folder menu, changed
their text to a color meant for use as a background rather than as
text when hovered. On top of the generic hover background used
across the app, this produced poor contrast in some themes.

Hovered menu items now keep the same text color already used before
hovering, which reads correctly against the hover background in
every theme.

https://github.com/owncloud/ocis/pull/12994
