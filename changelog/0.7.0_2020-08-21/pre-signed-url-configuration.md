Enhancement: add configuration options for the pre-signed url middleware 

Added an option to define allowed http methods for pre-signed url requests.
This is useful since we only want clients to GET resources and don't upload anything with presigned requests.

https://github.com/owncloud/ocis-proxy/issues/91
https://github.com/owncloud/product/issues/150
