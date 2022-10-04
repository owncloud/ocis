# Thumbnails service

The thumbnails service provides methods to generate thumbnails for various files.
It retrieves the sources at the location where the user files are stored and saves the tumbnails at the location where system files are stored. Those locations have defaults but can be manually defined via environment variables. The relevant environment variables are:

-   `OCIS_BASE_DATA_PATH` which will contain system relevant data and
-   `STORAGE_USERS_OCIS_ROOT` which will contain the user source data.
-   `THUMBNAILS_FILESYSTEMSTORAGE_ROOT` used if thumnails should be separated from system files.

For details and defaults see the documentation.

Thumbnails can be generated for the following file types:

-   png
-   jpg
-   gif
-   tiff
-   bmp
-   txt

Thumbnails can either be generated as `png`, `jpg` or `gif` files, various resolutions can be defined via `THUMBNAILS_RESOLUTIONS`.

---

**NOTE**

Since source files need to be loaded into memory when generating thumbnails, large source files could potentially crash this service if there is insufficient memory available. For bigger instances, this service can be dedicated to own servers with more memory when using container orchestration deployment methods.

---
