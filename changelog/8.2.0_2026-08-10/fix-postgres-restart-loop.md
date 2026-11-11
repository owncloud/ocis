Bugfix: Fix PostgreSQL container restart loop in Keycloak deployments

The PostgreSQL volume was mounted directly at `/var/lib/postgresql/data`.
On ext4 storage backends, Docker creates a `lost+found` directory at the volume root,
causing PostgreSQL's `initdb` to fail because the data directory is not empty.
This resulted in the container entering a restart loop.

The volume mount path has been changed to `/var/lib/postgresql` so that
PostgreSQL creates the `data/` subdirectory itself, avoiding the conflict.

https://github.com/owncloud/ocis/pull/12359
