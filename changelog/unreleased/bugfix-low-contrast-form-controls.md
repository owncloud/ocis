Bugfix: Fix low-contrast form controls in dark and vault themes

Several form controls were hard to see against dark and vault theme
backgrounds. The view-options range slider had no visible track border and
relied on an opacity dip on hover, making it barely visible at rest. The
switch component's off-state track and thumb, and its on-state thumb, used
colors too close to their surrounding background to read as a toggle. The
select combobox's highlighted/selected option text used a color too close to
the highlight background.

The global search input's text color was forced with `!important`,
overriding the theme-specific search text color some themes set, which made
typed text render in white on a white search input background in the vault
"Dark Theme – High Contrast". The search placeholder color also used an
invalid CSS custom property fallback (a bare property name instead of a
nested `var()`), silently breaking the placeholder color for any theme that
didn't define a search-specific placeholder token, and the vault "Dark
Theme"'s placeholder color was too close to its input's text and border color
to read as muted.

The slider now has a visible border and no longer dims on rest, the switch
and select now use colors with sufficient contrast against their
backgrounds, the `!important` override blocking theme-specific search colors
has been removed, the invalid placeholder fallback has been fixed, and the
vault "Dark Theme" placeholder color has been changed to one with visible
contrast against the input text.

https://github.com/owncloud/ocis/pull/12639
