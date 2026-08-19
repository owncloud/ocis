Bugfix: Fix double HTML-escaping of notification emails for multiple recipients

The notifications service reuses the same template variables map for every
recipient of an event (for example when a group is invited to a space). The
helper that HTML-escapes those variables for the HTML email body escaped the
map in place, so every recipient after the first received a plain-text body
containing HTML entities and an HTML body with one additional layer of escaping
per recipient (a space named `R&D` became `R&amp;amp;D`). The helper now returns
a new map and leaves the shared variables untouched, so every recipient renders
from the original values.

https://github.com/owncloud/ocis/pull/12413
