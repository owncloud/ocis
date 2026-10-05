Bugfix: Don't report network failures as wrong credentials on login

The login page showed "Invalid username or password" for every kind of
failure: when the IdP rejected the credentials, when the server answered with
an error, and when the request never reached the server at all (no network,
server down, connection blocked by a firewall). Users with a correct password
kept retrying it because the page told them it was wrong. An aborted request
was also mapped to the wrong-credentials message, and a request that got no
answer left the form waiting with no time limit.

Network failures and timeouts now show "The server could not be reached",
server errors show a generic "Something went wrong", and only a rejection by
the IdP (204/401) reports invalid credentials. Login requests now give up
after 30 seconds without an answer.

https://github.com/owncloud/ocis/pull/13072
