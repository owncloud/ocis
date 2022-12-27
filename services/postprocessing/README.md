# Postprocessing service

The `postprocessing` service handles coordination of asynchronous postprocessing. 

## Postprocessing functionality

The storageprovider service (`storage-users`) can be configured to do asynchronous postprocessing by setting the `STORAGE_USERS_OCIS_ASYNC_UPLOADS` envvar to true.
If this is the case, the storageprovider will initiate an asynchronous postprocessing after he has reveived all bytes of an upload. The `postprocessing` service will then 
coordinate various postprocessing steps (like e.g. scan the file for viruses). During postprocessing the file will be in a `processing` state during which only limited actions are available.

## Prerequisites for using `postprocessing` service

In the storageprovider (`storage-users`) set `STORAGE_USERS_OCIS_ASYNC_UPLOADS` envvar to `true`. Configuring any postprocessing step will require an additional service to be enabled and configured.
For example to use `virusscan` step one needs to have an enabled and configured `antivirus` service. 

All of this functionality will need an event system to be configured for all services: `ocis` ships with
`nats` enabled by default.

## Postprocessing steps

As of now ocis allows two different postprocessing steps to be enabled via envvar

### Virus scanning

Can be set via envvar `POSTPROCESSING_VIRUSSCAN`. This means that each upload is virus scanned during postprocessing. `antivirus` service is needed for this to work.

### Delay

Can be set via envvar `POSTPROCESSING_DELAY`. This step will just sleep for the configured amount of time. Intended for testing postprocessing functionality. NOT RECOMMENDED on productive systems.
