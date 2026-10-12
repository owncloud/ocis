Bugfix: added a keyboard shortcut to move focus out of the text editor

Opening the text/markdown editor moved keyboard focus straight into the
edit area, but Tab is bound to indenting rather than moving focus, so
there was no reliable way to reach the toolbar without leaving the file
entirely. The editor's own built-in escape hatch for this doesn't work
on a standard Mac keyboard, since the key combo it relies on is consumed
by macOS as a diacritic shortcut before it reaches the page.

Pressing Ctrl+M now moves focus from the edit area to the toolbar, the
same way on every platform. The editor's accessible description also
announces this shortcut for assistive technology users.

https://github.com/owncloud/ocis/pull/13004
