Bugfix: Generate animated GIF thumbnails from the first frame

Animated GIF thumbnails decoded and resized every frame even though thumbnails
do not need to preserve animation. The thumbnail preprocessor now stops after
the first frame, reducing CPU time and memory use for animations with many
frames while leaving the original file unchanged for the media viewer.

https://github.com/owncloud/ocis/issues/13028
