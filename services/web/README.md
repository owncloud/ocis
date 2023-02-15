# Web Service

The web service embeds and serves the static files for the [Infinite Scale Web Client](https://github.com/owncloud/web).  
Note that clients will respond with a connection error if the web service is not available.

The web service also provides a minimal API for branding functionality like changing the logo shown.

## Custom Compiled Web Assets

If you want to use your custom compiled web client assets instead of the embedded ones, then you can do that by setting the `WEB_ASSET_PATH` variable to point to your compiled files. See [ownCloud Web / Getting Started](https://owncloud.dev/clients/web/getting-started/) and [ownCloud Web / Setup with oCIS](https://owncloud.dev/clients/web/backend-ocis/) for more details.
