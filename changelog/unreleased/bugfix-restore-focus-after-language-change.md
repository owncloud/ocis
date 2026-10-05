Bugfix: Restore focus after changing the account language

Selecting a new language in the account page's Language dropdown dropped
keyboard focus to the page body instead of keeping it on the dropdown, and
the confirmation message was not announced to screen readers afterward.

The language change triggers a fair amount of background work (reloading
translations, switching the active language, saving the preference) before
it finishes, which raced against the dropdown's own best-effort focus
handling. Focus is now explicitly restored only once that work has fully
settled, and before the confirmation message is shown, so the announcement
is no longer cancelled by a focus change happening right after it.

https://github.com/owncloud/ocis/pull/13045
