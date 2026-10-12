Bugfix: Fix missing labels and low-contrast elements across the web UI

Several icon-only buttons (search, hide-share, create-user, gallery
navigation) were missing labels, and loading spinners and progress bars did
not expose a name describing what they represent. The global search and
create shortcut dropdowns now use correct roles and states. Tooltips no longer
get caught mid-fade with a low-contrast color when the reduced-motion system
setting is enabled, and low-contrast colors were fixed in the select combobox
search input and in markdown editor links and inline code.

https://github.com/owncloud/ocis/pull/12526
