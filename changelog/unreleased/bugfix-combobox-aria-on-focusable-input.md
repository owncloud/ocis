Bugfix: Announce autocomplete comboboxes correctly to screen readers

Autocomplete fields such as the language and design selectors on the account
page did not announce themselves as comboboxes to assistive technology. The
combobox role and its related ARIA attributes were set on a wrapper element
that could not receive keyboard focus, while the actual focusable input
exposed a generic, hardcoded label instead of the field's own label.

The combobox role, its state attributes and the accessible name are now set
directly on the focusable input, so screen readers announce the field as a
combobox with the correct label.

https://github.com/owncloud/ocis/pull/13041
