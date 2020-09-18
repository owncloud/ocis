Change: Use predefined resolutions for thumbnail generation 

We implemented predefined resolutions to prevent attacker from flooding the service with a large number of thumbnails.
The requested resolution gets mapped to the closest matching predefined resolution.

https://github.com/owncloud/ocis-thumbnails/issues/7
