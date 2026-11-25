Bugfix: Bump vips to 8.18.7-r1 in the Docker images

The Docker images now install vips 8.18.7-r1 from the Alpine edge
repository. The Alpine edge repository removed vips 8.18.5-r0. Thus the
Docker image builds failed with "unable to select packages".

https://github.com/owncloud/ocis/pull/13452
