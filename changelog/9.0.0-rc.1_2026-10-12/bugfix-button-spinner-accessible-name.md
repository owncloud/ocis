Bugfix: fixed missing accessible name on button loading spinner

`OcButton`'s loading spinner is a `role="img"` element, but it had no
accessible name. This tripped an axe `role-img-alt` violation on any
page whose e2e a11y check ran while a button was in its loading state,
such as the public link password page's "Continue" button.

The spinner now gets a real accessible name ("Loading") via
`OcSpinner`'s `ariaLabel` prop, marked `aria-live="polite"` so it is
actively announced to screen reader users instead of only being
described if focus happens to already be on the button.

https://github.com/owncloud/ocis/pull/13018
