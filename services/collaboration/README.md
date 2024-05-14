# Collaboration

The collaboration service connects ocis with document servers such as Collabora and ONLYOFFICE using the WOPI protocol.

Since this service requires an external document server, it won't start by default when using `ocis server`. You must start it manually with the `ocis collaboration server` command.

## Requirements

The collaboration service requires the target document server (ONLYOFFICE, Collabora, etc.) to be up and running. Additionally, some Infinite Scale services are also required to be running in order to register the GRPC service for the `open in app` action in the webUI. The following internal and external services need to be available:

* External document server.
* The gateway service.
* The app-registry service.

If any of the named services above have not been started or are not reachable, the collaboration service won't start. For the binary or the docker release of Infinite Scale, check with the `ocis list` command if they have been started. If not, you must start them manually upfront before starting the collaboration service.

## WOPI Configuration

There are a few variables that you need to set:

* `COLLABORATION_WOPIAPP_ADDR`:\
  The URL of the WOPI app (onlyoffice, collabora, etc).\
  For example: `https://office.example.com`.

* `COLLABORATION_HTTP_ADDR`:\
  The external address of the collaboration service. The target app (onlyoffice, collabora, etc) will use this address to read and write files from Infinite Scale.\
  For example: `https://wopiserver.example.com`.

* `COLLABORATION_HTTP_SCHEME`:\
  The scheme to be used when accessing the collaboration service. Either `http` or `https`. This will be used to finally build the URL that the WOPI app needs in order to contact the collaboration service.

The rest of the configuration options available can be left with the default values.
