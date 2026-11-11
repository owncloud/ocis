Bugfix: Show mermaid diagrams in the PDF export

Exporting a markdown file as PDF replaced every mermaid diagram with "Failed
to render Mermaid chart." when the theme defined its colors with `oklch()`. The
diagrams are now included as images in the PDF.

https://github.com/owncloud/ocis/pull/13427
