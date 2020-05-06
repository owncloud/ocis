---
title: "EOS"
date: 2020-02-27T20:35:00+01:00
weight: 30
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs
geekdocFilePath: eos.md
---

{{< toc >}}

OCIS can be configured to run on top of [eos](https://eos.web.cern.ch/). While the [eos documentation](http://eos-docs.web.cern.ch/) does cover a lot of topics it leaves out some details that you may have to either pull from various [docker containers](https://gitlab.cern.ch/eos/eos-docker), the [forums](https://eos-community.web.cern.ch/) or even the [source](https://github.com/cern-eos/eos) itself.

This document is a work in progress of the current setup.

## Current status

Using ocis and eos it is possible today to manage folders. Sharing is [heavily](https://github.com/cs3org/reva/pull/523) [under](https://github.com/cs3org/reva/pull/585) [development](https://github.com/cs3org/reva/pull/482). FIle up and download needs proper configuration of the dataprovider to also use eos.

## How to do it

### Grab it!

```
$ git clone git@github.com:owncloud/ocis.git
$ cd ocis
```


### Run it!

Preconditions
* `go` (from golang.org/dl) and `gcc` (via e.g. `apt install build-essential`) are installed
* No eos components are running. If in doubt, begin with `make eos-stop`

We poured the nitty gritty details of setting up ocis into Makefile targets. After running

```
$ make eos-start
```

the eos related docker containers will be created, started and setup to authenticate a gainst the ocis-glauth service.

It will also copy the ocis binary tho the `eos-cli1` container and start `ocis reva-storage-home` with the necessary environment variables to use the eos storage driver.

For details have a look at the `Makefile`.


### Test it!

You should now be able to point your browser to https://localhost:9200 and login using the demo user credentials, eg `einstein:relativity`.

{{< hint info >}}
If you encounter an error when the IdP redirects you back to phoenix, just reload the page and it should be gone ... or debug it. PR welcome!
{{< /hint >}}

Create a folder in the ui. Then check it was created in eos:

```
$ docker exec -it eos-mgm1 eos ls /eos/dockertest/einstein
```

Now create a new folder in eos (using eos-mgm1 you will be logged in as admin, see the `whoami`, which is why we `chown` the folder to the uid and gid of einstein afterwards):

```
$ docker exec -it eos-mgm1 eos whoami
$ docker exec -it eos-mgm1 eos mkdir /eos/dockertest/einstein/rocks
$ docker exec -it eos-mgm1 eos chown 20000:30000 /eos/dockertest/einstein/rocks
```

Check that the folder exists in the web ui.

## Next steps

- configure storage-home-data to enable file upload, PRs against `ocis-reva` welcome
- get sharing implemented, PRs against `reva` welcome
- simplify home logic, see https://github.com/cs3org/reva/issues/601 and https://github.com/cs3org/reva/issues/578
