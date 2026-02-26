Enhancement: Add photo EXIF metadata to search index and WebDAV results

We've added support for photo metadata fields in the Bleve search index and
WebDAV REPORT responses. This enables photo gallery applications to efficiently
query photos by their EXIF metadata and display camera information.

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

These fields are returned in WebDAV search results using the `oc:photo-*`
property namespace, allowing web extensions to build photo timeline views,
filter by camera, or show photos on a map.

https://github.com/owncloud/ocis/pull/11912
