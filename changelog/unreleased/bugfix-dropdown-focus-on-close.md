Bugfix: Keep keyboard focus on the trigger when a dropdown closes

Dropdown and context menus that close on item click (e.g. the per-file "..."
menu) removed the clicked item from the page as they closed, so keyboard
focus fell back to the page body. Screen readers then re-announced the
whole page, drowning out any status message the action had just triggered.

The dropdown now returns focus to the button that opened it when it closes,
unless the clicked action already moved focus somewhere else on its own.

https://github.com/owncloud/ocis/pull/13025
