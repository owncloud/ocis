---
title: Postprocessing
date: 2023-04-19T07:33:00.967407009Z
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/services/postprocessing
geekdocFilePath: _index.md
geekdocCollapseSection: true
---

## Abstract

The `postprocessing` service handles the coordination of asynchronous postprocessing steps. 

## Table of Contents

* [General Prerequisites](#general-prerequisites)
* [Postprocessing Functionality](#postprocessing-functionality)
* [Additional Prerequisites for the Postprocessing Service](#additional-prerequisites-for-the-postprocessing-service)
* [Postprocessing Steps](#postprocessing-steps)
  * [Virus Scanning](#virus-scanning)
  * [Delay](#delay)
  * [Custom Postprocessing Steps](#custom-postprocessing-steps)
    * [Prerequisites](#prerequisites)
    * [Workflow](#workflow)
* [Example Yaml Config](#example-yaml-config)

## General Prerequisites

To use the postprocessing service, an event system needs to be configured for all services. By default, `ocis` ships with a preconfigured `nats` service.

## Postprocessing Functionality

The storageprovider service (`storage-users`) can be configured to initiate asynchronous postprocessing by setting the `STORAGE_USERS_OCIS_ASYNC_UPLOADS` environment variable to `true`. If this is the case, postprocessing will get initiated *after* uploading a file and all bytes have been received.
The `postprocessing` service will then coordinate configured postprocessing steps like scanning the file for viruses. During postprocessing, the file will be in a `processing state` where only a limited set of actions are available. Note that this processing state excludes file accessability by users.
When all postprocessing steps have completed successfully, the file will be made accessible for users.

## Additional Prerequisites for the Postprocessing Service

When postprocessing has been enabled, configuring any postprocessing step will require the requested services to be enabled and pre-configured. For example, to use the `virusscan` step, one needs to have an enabled and configured `antivirus` service. 

## Postprocessing Steps

The postporcessing service is individually configurable. This is achieved by allowing a list of postprocessing steps that are processed in order of their appearance in the `POSTPROCESSING_STEPS` envvar. This envvar expects a comma separated list of steps that will be executed. Currently known steps to the system are `virusscan` and `delay`. Custom steps can be added but need an existing target for processing.

### Virus Scanning

To enable virus scanning as a postprocessing step after uploading a file, the environment variable `POSTPROCESSING_STEPS` needs to contain the word `virusscan` at one location in the list of steps. As a result, each uploaded file gets virus scanned as part of the postprocessing steps. Note that the `antivirus` service is required to be enabled and configured for this to work.

### Delay

Though this is for development purposes only and NOT RECOMMENDED on production systems, setting the environment variable `POSTPROCESSING_DELAY` to a duration not equal to zero will add a delay step with the configured amount of time. ocis will continue postprocessing the file after the configured delay. Use the enviroment variable `POSTPROCESSING_STEPS` and the keyword `delay` if you have multiple postprocessing steps and want to define their order. If `POSTPROCESSING_DELAY` is set but the keyword `delay` is not contained in `POSTPROCESSING_STEPS`, it will be processed as last postprocessing step without being listed there. In this case, a log entry will be written on service startup to notify the admin about that situation. That log entry can be avoided by adding the keyword `delay` to `POSTPROCESSING_STEPS`.

### Custom Postprocessing Steps

By using the envvar `POSTPROCESSING_STEPS`, custom postprocessing steps can be added. Any word can be used as step name but be careful not to conflict with exising keywords like `virusscan` and `delay`. In addition, if a keyword is misspelled or the corresponding service does either not exist or does not follow the necessary event communication, the postprocessing service will wait forever getting the required response to proceed and does not continue any other processing.

#### Prerequisites

For using custom postprocessing steps you need a custom service listening to the configured event system (see `General Prerequisites`)

#### Workflow

When setting a custom postprocessing step (eg. `"customstep"`) the postprocessing service will eventually sent an event during postprocessing. The event will be of type `StartPostprocessingStep` with its field `StepToStart` set to `"customstep"`. When the custom service receives this event it can savely execute its actions, postprocessing service will wait until it has finished its work. The event contains further information (filename, executing user, size, ...) and also required tokens and urls to download the file in case byte inspection is necessary.
Once the custom service has finished its work, it should sent an event of type `PostprocessingFinished` via the configured events system. This event needs to contain a `FinishedStep` field set to `"customstep"`. It also must contain the outcome of the step, which can be one of "delete" (abort postprocessing, delete the file), "abort" (abort postprocessing, keep the file) and "continue" (continue postprocessing, this is the success case).
See the [cs3 org](https://github.com/cs3org/reva/blob/edge/pkg/events/postprocessing.go) for up-to-date information of reserved step names and event definitons.

## Example Yaml Config

{{< include file="services/_includes/postprocessing-config-example.yaml"  language="yaml" >}}

{{< include file="services/_includes/postprocessing_configvars.md" >}}

