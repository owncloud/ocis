Enhancement: Colour tag chips by tag name

Tag chips were all the same grey, so a file's tags were hard to tell apart at
a glance and the same tag looked different to nothing across the app.

Tag chips are now filled with a colour derived from the tag name, so a given
tag always looks the same — for every user, on every device, without the
colour being stored anywhere. Themes supply the palette through a new
`designTokens.tagColorsList`, and the label colour is chosen per fill so it
stays legible in both the light and dark themes.

The tag overflow indicator on the files list also became a real, keyboard
reachable button, and it now shows the hidden tags in a popover instead of
opening the sidebar.

https://github.com/owncloud/ocis/pull/12892
