Security: Don't expose services by default

We've changed the bind behaviour for all non public facing services. Before this PR
all services would listen on all interfaces. After this PR, all services listen on
127.0.0.1 only, except the proxy which is listening on 0.0.0.0:9200.

https://github.com/owncloud/ocis/issues/2612
