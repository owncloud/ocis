Bugfix: fixed low color contrast on notification bell badge

The unread notification count badge on the notification bell used a
hardcoded red background with white text, which fell below the
WCAG AA contrast requirement for text that size.

The badge now uses the design system's danger color tokens instead,
which are already tuned to keep sufficient contrast with their
paired text color across all theme variants.

https://github.com/owncloud/ocis/pull/13000
