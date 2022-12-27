# Postprocessing Service

The `postprocessing` service handles the coordination of asynchronous postprocessing steps. 

## General Prerequisites

To use the postprocessing service, an event system needs to be configured for all services. By default, `ocis` ships with a preconfigured `nats` service.

## Postprocessing Functionality

The storageprovider service (`storage-users`) can be configured to initiate asynchronous postprocessing by setting the `STORAGE_USERS_OCIS_ASYNC_UPLOADS` environment variable to `true`. If this is the case, postprocessing will get initiated *after* uploading a file and all bytes have been recieved.

The `postprocessing` service will then coordinate configured postprocessing steps like scanning the file for viruses. During postprocessing, the file will be in a `processing state` where only a limited set of actions are available. Note that this processing state excludes file accessability by users.

When all postprocessing steps have completed successfully, the file will be made accessible for users.

## Additional Prerequisites for the `postprocessing` Service

When postprocessing has been enabled, configuring any postprocessing step will require the requested services to be enabled and pre-configured. For example, to use the `virusscan` step, one needs to have an enabled and configured `antivirus` service. 

## Postprocessing Steps

As of now, `ocis` allows two different postprocessing steps to be enabled via an environment variable.

### Virus Scanning

To enable virus scanning as postprocessing step after uploading a file, the environment variable  `POSTPROCESSING_VIRUSSCAN` needs to be set to ` true`. As a result, each uploaded file gets virus scanned as part of the postprocessing steps. Note that the `antivirus` service is required to be enabled and configured for this to work.

### Delay

Though this is for development purposes only and NOT RECOMMENDED on productive systems, setting the environment variable `POSTPROCESSING_DELAY` to a duration not equal to zero will add a delay step with the configured amount of time. ocis will continue postprocessing the file after the configured delay.
