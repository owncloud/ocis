---
title: "Discover oCIS with Docker"
date: 2022-06-14T16:00:00+02:00
weight: 8
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/guides
geekdocFilePath: ocis-local-docker.md
geekdocCollapseSection: true
---

{{< toc >}}

## Prerequisites

- Local docker installation (e.g. Docker for Desktop)
- Check [oCIS and Containers]({{< ref "ocis-and-containers" >}})

## Start oCIS with docker compose

### Create the project

Use the following skeleton as a docker-compose.yml:

```bash
mkdir simple-ocis && \
cd simple-ocis && \
touch docker-compose.yml
```

Copy the following file content into `docker-compose.yml` and save it.

```yaml
version: "3.7"

services:
  ocis:
    image: owncloud/ocis:latest
    environment:
      # INSECURE: needed if oCIS / Traefik is using self generated certificates
      OCIS_INSECURE: "true"

      # OCIS_URL: the external domain / ip address of oCIS (with protocol, must always be https)
      OCIS_URL: "https://localhost:9200"

      # OCIS_LOG_LEVEL: error / info / ... / debug
      OCIS_LOG_LEVEL: info
```

### Prepare Paths

Create directories if not exists:

```bash
mkdir -p $(pwd)/ocis-config \
mkdir -p $(pwd)/ocis-data
```


Set the user for the directories to be the same as the user inside the container:

```bash
sudo chown -Rfv 1000:1000 $(pwd)/ocis-config/ \
sudo chown -Rfv 1000:1000 $(pwd)/ocis-data
```

### Initialize

```bash
docker run --rm -it -v $(pwd):/etc/ocis/ owncloud/ocis:latest init
```

You will get the following output:

```txt {hl_lines=[9]}
Do you want to configure Infinite Scale with certificate checking disabled?
 This is not recommended for public instances! [yes | no = default] yes

=========================================
 generated OCIS Config
=========================================
 configpath : /etc/ocis/ocis.yaml
 user       : admin
 password   : t3p4N0jJ47LbhpQ04s9W%u1$d2uE3Y.3
```

Check your local folder. We just generated a default ocis configuration file with random passwords and secrets.

```bash
ls # list the current folder
docker-compose.yml                    ocis.yaml # ocis.yaml has been generated
```

Run `cat ocis.yaml`

```yaml {linenos=table,hl_lines=[19]}
token_manager:
  jwt_secret: X35rffWpS9BR.=^#LDt&z3ykYOd7h@w*
machine_auth_api_key: -0$4ieu5+t6HD6Ui^0PpKU6B0qxisv.m
system_user_api_key: oVxICwMR9YcKXTau+@pqKZ0EO-OHz8sF
transfer_secret: e%3Sda=WFBuy&ztBUmriAbBR$i2CmaDv
system_user_id: b7d976a1-7300-4db7-82df-13502d6b5e18
admin_user_id: c59a6ae9-5f5e-4eef-b82e-0e5c34f93e52
graph:
  spaces:
    insecure: false
  identity:
    ldap:
      bind_password: wElKpGjeH0d.E4moXh=.dc@s2CtB0vy%
idp:
  ldap:
    bind_password: Ft2$2%#=6Mi22@.YPkhh-c6Kj=3xBZAb
idm:
  service_user_passwords:
    admin_password: t3p4N0jJ47LbhpQ04s9W%u1$d2uE3Y.3
    idm_password: wElKpGjeH0d.E4moXh=.dc@s2CtB0vy%
    reva_password: pJAdZ2fU!SFKgcdDPRW%ruIiNM6GnN1D
    idp_password: Ft2$2%#=6Mi22@.YPkhh-c6Kj=3xBZAb
proxy:
  insecure_backends: false
frontend:
  archiver:
    insecure: false
auth_basic:
  auth_providers:
    ldap:
      bind_password: pJAdZ2fU!SFKgcdDPRW%ruIiNM6GnN1D
auth_bearer:
  auth_providers:
    oidc:
      insecure: false
users:
  drivers:
    ldap:
      bind_password: pJAdZ2fU!SFKgcdDPRW%ruIiNM6GnN1D
groups:
  drivers:
    ldap:
      bind_password: pJAdZ2fU!SFKgcdDPRW%ruIiNM6GnN1D
storage_system:
  data_provider_insecure: false
storage_users:
  data_provider_insecure: false
ocdav:
  insecure: false
thumbnails:
  thumbnail:
    transfer_secret: z-E%G8MTeFpuT-ez2o8BjfnG1Jl2yLLm
    webdav_allow_insecure: false
    cs3_allow_insecure: false
```

{{< hint type=tip title="Admin password" >}}
**Password initialisation**\
During the run of `./ocis init`, the password for the `admin` user has been set to a random string.

You can override that later by setting `IDM_ADMIN_PASSWORD=secret`. The ENV variable setting always overrides the config file.
{{< /hint >}}

### Mount the config file

Add the config file as a bind mount.

```yaml
    volumes:
      # mount the ocis config file inside the container
      - "./ocis.yaml:/etc/ocis/ocis.yaml"
```

### Apply the changes

```bash
docker compose up -d
```

The service should be running.

```bash
docker compose ps
NAME                 COMMAND                  SERVICE             STATUS              PORTS
simple-ocis-ocis-1   "/usr/bin/ocis server"   ocis                running             9200/tcp
```

### Open the port 9200 to the outside

```yaml
ports:
  - 9200:9200
```

Add a port mapping to your docker compose file and run `docker compose up -d` again. You should now be able to access <https://localhost:9200> and log in. You will get a warning from your browser that the connection is not safe because we are using a self-signed certificate. Accept this warning message to continue. Use the user `admin` and the password which has been initialized before.

Congratulations! You have successfully set up a simple Infinite Scale locally.

{{< hint type=tip title="Docker Compose Helps you" >}}
**docker compose up**\
You do not need to shut down your service to apply changes from the docker-compose.yml file. Running `docker compose up -d` again is enough. Docker compose always tries to bring all services to the desired state.

**docker compose ps**\
This gives you a status of the services of the project.

**docker compose exec \<service name\> \<command\>**\
This command is handy to run specific commands inside your service. Try `docker compose exec ocis ocis version`.
{{< /hint >}}

### Persist data, restart and logging

The key to a successful container setup is the persistence of the application data to make the data survive a re-boot. Docker normally uses [volumes](https://docs.docker.com/storage/volumes/) for this purpose. A volume can either be a "named volume" which are completely managed by docker and have many advantages (see the linked docker documentation), or "bind mounts" which are using the directory structure and OS of the host system. In our example we already use a bind mount for the config file. We will now add a named volume for the oCIS data directory.

This is the way we should configure the ocis service:

```yaml
    volumes:
      # mount the ocis config file inside the container
      - "./ocis.yaml:/etc/ocis/ocis.yaml"
      # short syntax for using a named volume
      # in the form <volume name>:<mount path in the container>
      # use a named volume for the ocis data directory
      - "ocis-data:/var/lib/ocis"
      # or the more verbose syntax
      #- type: volume
      #  source: ocis-data # name of the volume
      #  target: /var/lib/ocis # the mount path inside the container
```

The docker-compose.yml needs to declare the named volumes globally, add this at the bottom of the file:

```yaml
# https://docs.docker.com/compose/compose-file/compose-file-v3/#volumes
# this declares the named volume with all default settings
# you can also see the volume when running `docker volume list`
volumes:
  ocis-data:
```

Now let us configure the restart policy and the logging settings for the ocis service:

```yaml
    # https://docs.docker.com/compose/compose-file/compose-file-v3/#restart
    restart: always # or on-failure / unless-stopped

    # https://docs.docker.com/config/containers/logging/configure/
    # https://docs.docker.com/compose/compose-file/compose-file-v3/#logging
    # the default log driver does no log rotation
    # you can switch to the "local" log driver which does rotation by default
    logging:
      driver: local
    # otherwise you could specify log rotation explicitly
    #  driver: "json-file" # this is the default driver
    #  options:
    #    max-size: "200k" # limit the size of the log file
    #    max-file: "10" # limit the count of the log files
```

Apply your changes! Just run `docker compose up -d` again.

Now you have an oCIS which will survive reboots, restart automatically and has log rotation by default.

Access the logs via `docker compose logs -f` and do some actions in the frontend to see the effect. Create data by uploading files and adding more users. Then run `docker compose down` to shut oCIS down. Start it again `docker compose up -d`, log in again and check that your data has survived the reboot.

### Pin the oCIS version

Last but not least, it is never a good idea to use the `latest` docker tag. Pin your container image to a released version.

```yaml
    image: owncloud/ocis:latest@sha256:5ce3d5f9da017d6760934448eb207fbaab9ceaf0171b4122e791e292f7c86c97
    # the latest tag is not recommended, because you don't know which version you'll get
    # but even if you use `owncloud/ocis:1.16.0` you cannot be sure that you'll get
    # the exact same image if you pull it at a later point in time (because docker image tags are not immutable).
    # To be 100% that you always get the same image, you can pin the digest (hash) of the
    # image. If you do a `docker pull owncloud/ocis:latest`, it also will also show you the digest.
    # see also https://docs.docker.com/engine/reference/commandline/images/#list-image-digests
```

## Wrapping up

If you have completed this guide, your docker-compose.yml should look like the following example:

{{< expand "Solution" "..." >}}
```yaml
version: "3.7"

services:
  ocis:
    image: owncloud/ocis:latest@sha256:5ce3d5f9da017d6760934448eb207fbaab9ceaf0171b4122e791e292f7c86c97
    # the latest tag is not recommended, because you don't know which version you'll get
    # but even if you use `owncloud/ocis:1.16.0` you cannot be sure that you'll get
    # the exact same image if you pull it at a later point in time (because docker image tags are not immutable).
    # To be 100% that you always get the same image, you can pin the digest (hash) of the
    # image. If you do a `docker pull owncloud/ocis:latest`, it also will also show you the digest.
    # see also https://docs.docker.com/engine/reference/commandline/images/#list-image-digests
    environment:
      # INSECURE: needed if oCIS / Traefik is using self generated certificates
      OCIS_INSECURE: "true"

      # OCIS_URL: the external domain / ip address of oCIS (with protocol, must always be https)
      OCIS_URL: "https://localhost:9200"

      # OCIS_LOG_LEVEL: error / info / ... / debug
      OCIS_LOG_LEVEL: info
    volumes:
      # mount the ocis config file inside the container
      - "./ocis.yaml:/etc/ocis/ocis.yaml"
      # short syntax for using a named volume
      # in the form <volume name>:<mount path in the container>
      # use a named volume for the ocis data directory
      - "ocis-data:/var/lib/ocis"
      # or the more verbose syntax
      #- type: volume
      #  source: ocis-data # name of the volume
      #  target: /var/lib/ocis # the mount path inside the container
    ports:
      - 9200:9200
    # https://docs.docker.com/compose/compose-file/compose-file-v3/#restart
    restart: always # or on-failure / unless-stopped

    # https://docs.docker.com/config/containers/logging/configure/
    # https://docs.docker.com/compose/compose-file/compose-file-v3/#logging
    # the default log driver does no log rotation
    # you can switch to the "local" log driver which does rotation by default
    logging:
      driver: local
    # otherwise you could specify log rotation explicitly
    #  driver: "json-file" # this is the default driver
    #  options:
    #    max-size: "200k" # limit the size of the log file
    #    max-file: "10" # limit the count of the log files

# https://docs.docker.com/compose/compose-file/compose-file-v3/#volumes
# this declares the named volume with all default settings
# you can also see the volume when running `docker volume list`
volumes:
  ocis-data:
```
{{< /expand >}}
