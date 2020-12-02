# Docker image for oCIS development with eos

Image is based on [owncloud/eos-base](https://hub.docker.com/r/owncloud/eos-base) from [eos-stack](https://github.com/owncloud-docker/eos-stack)

{{< hint info >}}
On MacOS do not mount a local folder to the `/var/tmp/reva/` path. The fuse driver used by docker [does not support extended attributes](https://docs.docker.com/v18.09/docker-for-mac/osxfs/). See [#182](https://github.com/owncloud/ocis/issues/182) for more details.
{{< /hint >}}

## Build
```shell
docker build -t owncloud/eos-ocis-dev:latest .
```

## Publish
```shell
docker push owncloud/eos-ocis-dev:latest
```

## Maintainer

* [Felix Böhm](https://github.com/felixboehm)

## Disclaimer
Use only for development or testing. Setup is not secured nor tested.

## oCIS development on eos

### Setup oCIS

To build and run your local ocis code with default storage driver

```shell
docker run --rm -ti --name ocis -v $PWD:/ocis -p 9200:9200 owncloud/eos-ocis-dev
```

ocis will use the owncloud storage driver and store files in the container at /var/tmp/reva/data/<username>/files

Data is here: `docker exec -it ocis ll /var/tmp/reva/`

Alternative: With the [docker-compose.yml file in ocis repo](https://github.com/owncloud/ocis/blob/master/docker-compose.yml) you can also start ocis via compose:

```shell
docker-compose up -d ocis
```

Now try to list the running services

```shell
docker-compose exec ocis ./bin/ocis list
```

## Setup eos storage

1. Start the eos cluster and ocis via the compose stack

```shell
docker-compose up -d
```

2. Configure to use eos storage driver instead of default storage driver

* kill the home storage and data providers. we need to switch them to the eoshome driver:

```shell
docker-compose exec ocis ./bin/ocis kill reva-storage-home
docker-compose exec ocis ./bin/ocis kill reva-storage-home-data
```

* restart them with the eoshome driver and a new layout:

```shell
docker-compose exec -e REVA_STORAGE_EOS_LAYOUT="{{substr 0 1 .Username}}/{{.Username}}" -e REVA_STORAGE_HOME_DRIVER=eoshome -d ocis ./bin/ocis run reva-storage-home
docker-compose exec -e REVA_STORAGE_EOS_LAYOUT="{{substr 0 1 .Username}}/{{.Username}}" -e REVA_STORAGE_HOME_DATA_DRIVER=eoshome -d ocis ./bin/ocis run reva-storage-home-data
```

* restart the reva frontend with a new namespace (pointing to the eos storage provider) for the dav files endpoint

```shell
docker-compose exec ocis ./bin/ocis kill reva-frontend
docker-compose exec -e DAV_FILES_NAMESPACE="/eos/" -d ocis ./bin/ocis run reva-frontend
```

* login with `einstein / relativity`, upload a file to einsteins home and verify the file is there using

```shell
docker-compose exec ocis eos ls -l /eos/dockertest/reva/users/e/einstein/
-rw-r--r--   1 einstein users              10 Jul  1 15:24 newfile.txt
```
