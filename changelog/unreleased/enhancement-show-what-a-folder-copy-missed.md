Enhancement: Show what an incomplete folder copy missed

When copying or duplicating a folder only copied part of its content, the
message now says how many of the files and folders inside it were copied, e.g.
"87 of 100 items copied.", and its details list the files and folders
that weren't copied. A copy of a folder is now reported as incomplete in that
case too, like a duplication, and the incomplete folder shows up in the files
list right away instead of only after reloading the page. Cutting and pasting
between spaces, which copies, and a successful duplication are no longer
reported as moves.

https://github.com/owncloud/ocis/pull/13499
