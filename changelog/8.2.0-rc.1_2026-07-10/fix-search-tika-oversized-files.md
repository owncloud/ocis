Bugfix: Extract metadata from oversized files and fix ISO field

Two issues were found in the Tika content extractor:

1. Files exceeding `SEARCH_CONTENT_EXTRACTION_SIZE_LIMIT` (default 20MB)
were skipped entirely — no EXIF, no photo metadata, no image dimensions
were extracted. This particularly affected Pixel Motion Photos (`.MP.jpg`)
which embed an MP4 video making them 3-9MB. Since EXIF metadata lives in
the JPEG header (first few KB), a truncated stream is sufficient. The
extractor now wraps the download in `io.LimitReader` instead of skipping
Tika, sending only the first N bytes for metadata extraction.

2. The ISO speed field was read from `"Base ISO"`, a Canon-specific Tika
field (sensor base sensitivity). Most cameras — Pixel, iPhone, Samsung —
provide ISO via the standard `"exif:IsoSpeedRatings"` field. The extractor
now checks `exif:IsoSpeedRatings` first and falls back to `Base ISO` for
Canon compatibility.

https://github.com/owncloud/ocis/pull/12000
