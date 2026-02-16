Change: generate cryptographically secure state token 

Replaced Math.random with a cryptographically secure way to generate the oidc state token using the javascript crypto api. 

https://developer.mozilla.org/en-US/docs/Web/API/Crypto/getRandomValues
https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Math/random
https://github.com/owncloud/ocis/pull/1203
