Enhancement: Theme Processing and Logo Customization

We have made significant improvements to the theme processing in Infinite Scale.
The changes include:

- Enhanced the way themes are composed. Now, the final theme is a combination of the built-in theme and the custom theme provided by the administrator via `WEB_ASSET_THEMES_PATH` and `WEB_UI_THEME_PATH`.
- Introduced a new mechanism to load custom assets. This is particularly useful when a single asset, such as a logo, needs to be overwritten.
- Fixed the logo customization option. Previously, small theme changes would copy the entire theme. Now, only the changed keys are considered, making the process more efficient.
- Default themes are now part of ocis. This change simplifies the theme management process for web.

These changes enhance the robustness of the theme handling in Infinite Scale and provide a better user experience.


https://github.com/owncloud/ocis/pull/9133
https://github.com/owncloud/ocis/issues/8966
