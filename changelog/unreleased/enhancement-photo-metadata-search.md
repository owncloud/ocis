Enhancement: Add photo EXIF metadata and Tika object detection to search

We've added support for photo metadata fields and AI-generated object detection
labels and captions in the Bleve search index and WebDAV REPORT responses.

The following photo metadata fields are now indexed and searchable:
- `photo.takenDateTime` - When the photo was taken (supports date range queries)
- `photo.cameraMake` - Camera manufacturer (e.g., Canon, Nikon, Samsung)
- `photo.cameraModel` - Camera model name
- `photo.fNumber` - Aperture f-stop value
- `photo.focalLength` - Focal length in millimeters
- `photo.iso` - ISO sensitivity
- `photo.orientation` - Image orientation
- `photo.exposureNumerator` - Exposure time numerator (for shutter speed calculation)
- `photo.exposureDenominator` - Exposure time denominator (for shutter speed calculation)

GPS location data is also included when available:
- `photo.location.latitude` - GPS latitude
- `photo.location.longitude` - GPS longitude
- `photo.location.altitude` - GPS altitude

Object detection labels and captions from Tika's ObjectRecognitionParser
(Inception V3 / Show and Tell models) are now extracted and searchable:
- `objectLabel:` - Search by detected object labels (e.g., `objectLabel:dog`)
- `objectCaption:` - Search by generated image captions (e.g., `objectCaption:beach`)
- Results are exposed as `oc:object-labels` and `oc:object-captions` WebDAV properties

These fields are returned in WebDAV search results, allowing web extensions
to build photo timeline views, filter by camera, show photos on a map, or
search photos by their visual content.

https://github.com/owncloud/ocis/pull/11912
https://github.com/owncloud/ocis/pull/12072
