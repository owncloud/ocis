Bugfix: Text and BMP file previews failed where libvips has no ImageMagick

Downloading the preview of a text or BMP file failed with a `500 could not get thumbnail` error on systems where libvips comes without ImageMagick support, such as minimal container images. The rendered text image was passed through a BMP encoding round-trip first, and uploaded BMP files were handed to libvips as bytes; libvips can only load BMP images via ImageMagick.

The rendered text image is now handed over to libvips directly from memory, and BMP files are decoded in Go before the handover, so both no longer depend on ImageMagick being installed alongside libvips.

https://github.com/owncloud/ocis/pull/12961
