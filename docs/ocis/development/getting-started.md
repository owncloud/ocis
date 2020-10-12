---
title: "Getting Started with Development"
date: 2020-07-07T20:35:00+01:00
weight: 15
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/development
geekdocFilePath: getting-started.md
---

{{< toc >}}

## Docker dev environment

### Option 1: Plain docker

To build and run your local ocis code with default storage driver

```
docker run --rm -ti --name ocis -v $PWD:/ocis -p 9200:9200 owncloud/eos-ocis-dev
```

The eos-ocis-dev container will build and run ocis using the owncloud storage driver and store files in the container at `/var/tmp/reva/data/<username>/files`

To check the uploaded files start digging with: `docker exec -it ocis ls -l /var/tmp/reva/`

{{< hint info >}}
On MacOS do not mount a local folder to the `/var/tmp/reva/` path. The fuse driver used by docker [does not support extended attributes](https://docs.docker.com/v18.09/docker-for-mac/osxfs/). See [#182](https://github.com/owncloud/ocis/issues/182) for more details.
{{< /hint >}}


### Option 2: Docker compose

With the `docker-compose.yml` file in ocis repo you can also start ocis via compose:

```
docker-compose up -d ocis
```

{{< hint info >}}
We are only starting the `ocis` container here.
{{< /hint >}}

## Verification

Check the services are running

```
$ docker-compose exec ocis ./bin/ocis list
+--------------------------+-----+
|        EXTENSION         | PID |
+--------------------------+-----+
| accounts                 | 172 |
| api                      | 204 |
| glauth                   | 187 |
| graph                    |  41 |
| graph-explorer           |  55 |
| konnectd                 | 196 |
| ocs                      |  59 |
| phoenix                  |  29 |
| proxy                    |  22 |
| registry                 | 226 |
| reva-auth-basic          |  96 |
| reva-auth-bearer         | 104 |
| reva-frontend            | 485 |
| reva-gateway             |  78 |
| reva-sharing             | 286 |
| reva-storage-eos         | 129 |
| reva-storage-eos-data    | 134 |
| reva-storage-home        | 442 |
| reva-storage-home-data   | 464 |
| reva-storage-oc          | 149 |
| reva-storage-oc-data     | 155 |
| reva-storage-public-link | 168 |
| reva-users               | 420 |
| settings                 |  23 |
| thumbnails               | 201 |
| web                      | 218 |
| webdav                   |  63 |
+--------------------------+-----+
```
