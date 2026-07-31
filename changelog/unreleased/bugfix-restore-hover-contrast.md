Bugfix: Restore hover feedback in high-contrast and dark themes

The Light and Dark High-Contrast themes set the background-highlight color
token equal to the default page background, making hover backgrounds
indistinguishable from the resting background wherever background-highlight
was used for interactive feedback, such as the global search results
dropdown. The non-high-contrast Dark theme had the same collision.

background-highlight now matches background-hover in the affected theme
definitions, and the search results dropdown and create-shortcut context
menu now use the background-hover color directly for their hover and active
states, since background-highlight is otherwise used for static surfaces
like cards, modals and form fields rather than interaction feedback.

https://github.com/owncloud/ocis/pull/12613
