Security: Disable chart rendering in the markdown editor

The markdown editor and preview, used for text files and for space and app
descriptions, rendered "echarts" code blocks as charts. To do so it evaluated
the block content as script code and could load the charting library from a
third-party CDN. The default content security policy blocks both, but a relaxed
policy would have allowed script code from such a code block to run for everyone
viewing it.

Chart rendering is now disabled, so these code blocks are shown as plain code.

https://github.com/owncloud/ocis/pull/13406
