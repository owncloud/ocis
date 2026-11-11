Bugfix: Preserve accessibility focus in the file versions sidebar

Opening the file versions panel with a keyboard or screen reader could move
focus to the main page, and restoring a version could briefly move focus to the
main page before selecting the sidebar's back button. The sidebar now manages
focus throughout panel transitions and loading states, restores focus to the
control that opened it when closed, and returns focus to the restored version's
Restore button after the operation completes. Version dates and action labels
also provide clearer context for screen reader users without duplicate
announcements or additional tab stops.

https://github.com/owncloud/ocis/pull/12630
