Enhancement: Make static GIF thumbnails configurable

Set THUMBNAILS_GIF_FIRST_FRAME_ONLY=true to generate static GIF thumbnails by
decoding only the first frame, reducing CPU time and memory use for animations
with many frames. Animated GIF thumbnails remain the default. The original file
is unchanged for the media viewer.

https://github.com/owncloud/ocis/pull/13113
https://github.com/owncloud/ocis/issues/13028
