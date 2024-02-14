# Collaboration

The collaboration service connects ocis with document servers such as collabora and onlyoffice using the WOPI protocol.

Since this service requires an external service (onlyoffice, for example), it won't run by default with the general `ocis server` command. You need to run it manually with the `ocis collaboration server` command.

## Requirements

The collaboration service requires the target document server (onlyoffice, collabora, etc) to be up and running.
We also need reva's gateway and app provider services to be running in order to register the GRPC service for the "open in app" action.

If any of those services are down, the collaboration service won't start.

## Configuration

There are a few variables that you need to set:

* `COLLABORATION_WOPIAPP_ADDR`: The URL of the WOPI app (onlyoffice, collabora, etc). For example: "https://office.mycloud.prv".
* `COLLABORATION_HTTP_ADDR`: The external address of the collaboration service. The target app (onlyoffice, collabora) will use this address to read and write files from ocis. For example: "wopiserver.mycloud.prv"
* `COLLABORATION_HTTP_SCHEME`: The scheme to be used when accessing the collaboration service. Either "http" or "https". This will be used to build the URL that the WOPI app needs in order to contact this service.

The rest of the configuration options available can be left with the default values.
