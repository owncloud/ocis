Bugfix: Don't report network failures as wrong credentials on login

The login page showed "Invalid username or password" for every kind of
failure: when the IdP rejected the credentials, when the server answered with
an error, and when the request never reached the server at all (no network,
server down, connection blocked by a firewall). Users with a correct password
kept retrying it because the page told them it was wrong. An aborted request
was also mapped to the wrong-credentials message, and a request that got no
answer left the form waiting with no time limit.

Network failures now show "The server could not be reached", a request with
no answer within 30 seconds shows "The server did not answer in time", and
server errors and any other unexpected error show a generic "Something went
wrong". Only a rejection by the IdP (204/401) reports invalid credentials.

https://github.com/owncloud/ocis/pull/13072
