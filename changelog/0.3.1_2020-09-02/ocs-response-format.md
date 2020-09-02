Bugfix: Add the top level response structure to json responses 

Probably during moving the ocs code into the ocis-ocs repo the response format was changed.
This change adds the top level response to json responses. Doing that the reponse should be compatible to the responses from OC10.

https://github.com/owncloud/product/issues/181
https://github.com/owncloud/product/issues/181#issuecomment-683604168

