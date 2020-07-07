---
title: "Getting Started with Development"
date: 2020-07-07T20:35:00+01:00
weight: 15
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs
geekdocFilePath: development.md
---

{{< toc >}}

## Docker dev environment

To build and run your local ocis code with default storage driver

```
docker run --rm -ti --name ocis -v $PWD:/ocis -p 9200:9200 owncloud/eos-ocis-dev
```

ocis will use the owncloud storage driver and store files in the container at /var/tmp/reva/data/<username>/files

Data is here: `docker exec -it ocis ll /var/tmp/reva/`

Alternative: With the `docker-compose.yml` file in ocis repo you can also start ocis via compose:

```
docker-compose up -d ocis
```

Now try to list the running services

```
docker-compose exec ocis ./bin/ocis list
```

## Docker dev environment for eos storage

1. Start the eos cluster and ocis via the compose stack

```
docker-compose up -d
```

2. Start the ldap authentication

```
docker-compose exec -d ocis /start-ldap
```

3. Configure to use eos storage driver instead of default storage driver

- kill the home storage and data providers. we need to switch them to the eoshome driver:

```
docker-compose exec ocis ./bin/ocis kill reva-storage-home
docker-compose exec ocis ./bin/ocis kill reva-storage-home-data
```

- restart them with the eoshome driver and a new layout:

```
docker-compose exec -e REVA_STORAGE_EOS_LAYOUT="{{substr 0 1 .Username}}/{{.Username}}" -e REVA_STORAGE_HOME_DRIVER=eoshome -d ocis ./bin/ocis run reva-storage-home
docker-compose exec -e REVA_STORAGE_EOS_LAYOUT="{{substr 0 1 .Username}}/{{.Username}}" -e REVA_STORAGE_HOME_DATA_DRIVER=eoshome -d ocis ./bin/ocis run reva-storage-home-data
```

- restart the reva frontend with a new namespace (pointing to the eos storage provider) for the dav files endpoint

```
docker-compose exec ocis ./bin/ocis kill reva-frontend
docker-compose exec -e DAV_FILES_NAMESPACE="/eos/" -d ocis ./bin/ocis run reva-frontend
```

- login with `einstein / relativity`, upload a file to einsteins home and verify the file is there using 

```
docker-compose exec ocis eos ls -l /eos/dockertest/reva/users/e/einstein/
-rw-r--r--   1 einstein users              10 Jul  1 15:24 newfile.txt
```
