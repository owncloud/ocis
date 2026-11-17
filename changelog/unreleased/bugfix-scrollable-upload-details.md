Bugfix: Keep the upload details list scrollable

After uploading or extracting many files at once, expanding "Show details" in
the upload panel grew the list past the top of the screen. The list could not
be scrolled and "Hide details" was out of reach until the page was reloaded.
The list now has a maximum height and scrolls, so the panel's title bar and
buttons stay visible.
