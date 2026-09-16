Bugfix: Text file previews failed where libvips has no ImageMagick support

Downloading the preview of a text file failed with a `500 could not get thumbnail` error on systems where libvips comes without ImageMagick support, such as minimal container images. The rendered text image was passed through a BMP encoding round-trip first, and libvips can only load BMP images via ImageMagick.

The rendered text image is now handed over to libvips directly from memory, so text previews no longer depend on ImageMagick being installed alongside libvips.

https://github.com/owncloud/ocis/pull/12961
