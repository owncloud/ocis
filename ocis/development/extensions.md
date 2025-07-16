---
title: "Extensions"
date: 2020-02-27T20:35:00+01:00
weight: 40
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/development
geekdocFilePath: extensions.md
---

oCIS is all about files, sync and share - but most of the time there is more you want to do with your files, e.g. having a different view on your photo collection or editing your offices files in an online file editor. ownCloud 10 faced the same problem and solved it with `applications`, which can extend the functionality of ownCloud 10 in a wide range. Since oCIS is different in its architecture compared to ownCloud 10, we had to come up with a similar (yet slightly different) solution. To extend the functionality of oCIS, you can write or install `extensions`. An extension is basically any running code which integrates into oCIS and provides functionality to oCIS and its users. Because extensions are just microservices providing an API, you can technically choose any programming language you like - a huge improvement to ownCloud 10, where it was nearly impossible to use a different programming language than PHP.

We will now introduce you to the oCIS extension system and show you how you can create a custom extension yourself.

## Extension examples

Technically every service in oCIS is an extension, even if oCIS would not really work without some of them. Therefore, you can draw inspiration from any of the plenty of extensions in the [oCIS monorepo](https://github.com/owncloud/ocis).

Besides these "default" extensions in the oCIS monorepo, there are two more extensions you should be aware of:

- [Hello](https://github.com/owncloud/ocis-hello)
- [WOPI server](https://github.com/owncloud/ocis-wopiserver)

Differences between the extensions maintained inside the oCIS monorepo and the ones maintained in their own repository are:

- extensions inside the [oCIS monorepo](https://github.com/owncloud/ocis) are all written in Go, whereas other extensions may choose the programming language freely
- extensions inside the oCIS monorepo heavily share tooling to reduce maintenance efforts, whereas other extensions may use different tooling (e.g. a different CI system)
- extensions inside the oCIS monorepo will be all build into one binary and started with the `ocis server` command. All other extensions must be started individually besides oCIS.


For quickstart purposes we also offer a [template project](https://github.com/owncloud/boilr-ocis-extension) which can be used to generate all the boilerplate code for you. But you also can decide to use your own project layout or even a different programming language.


## Integration into oCIS

Depending on the functionality of your extension, you might need to integrate with one or multiple of the components of oCIS mentioned below.

### ownCloud Web

If your extension is not just doing something in the background, you will need a UI in order to allow the user to interact with your extension. You could just provide your own web frontend for that purpose, but for a better user experience you can easily integrate into the web frontend of oCIS, the new [ownCloud Web](https://github.com/owncloud/web).

ownCloud Web allows you to write an extension for itself and therefore offers a seamless user experience. Upon login, the user will be able to use the application switcher to switch between the files view, settings and other available and installed extensions, yours included. Furthermore it is also possible to register your extension for different file actions. As an example, you could offer your extension to the user for creating and editing office documents. The user will then be able to create or open a file with your application directly from the files view. How to provide create an extension for ownCloud Web can be seen best in [the Hello extension](https://github.com/owncloud/ocis-hello/blob/master/ui/app.js), whereas plain file handling without any web frontend is available in the [WOPI server extension](https://github.com/owncloud/ocis-wopiserver/blob/master/ui/app.js).

To make ownCloud Web pick up your extension, you need to activate it in the configuration like seen in the [Hello extension](https://owncloud.dev/extensions/ocis_hello/running/#configure-and-start-ocis).

For a consistent look and feel, ownCloud Web uses an external design library, the [ownCloud design system](https://github.com/owncloud/owncloud-design-system). Since its classes and components are available through the wrapping `web runtime`, we highly recommend you to leverage it in your extension as well.

### Settings

An extension likely has some behaviour which the user can configure. Fundamental configuration will often be done by administrators during deployment, via configuration files or by setting environment variables. But for other settings, which are supposed to change more often or which are even user specific, this is not a viable way. Therefore you need to offer the users a UI where they can configure your extension to their liking. Because implementing something like this is a repetitive task among extensions, oCIS already offers the settings extensions which does that for your extension. Your extension just needs to register settings bundles, respective permissions and finally read the current values from the settings service. You can read more on that on the [settings extension]({{< ref "../../services/settings" >}}) and see how [oCIS Hello uses these settings](https://owncloud.dev/extensions/ocis_hello/settings/).

### Proxy

The Proxy is an API gateway and acts as the single connection point where all external request from users and devices need to pass through.

To make sure that requests can reach your extension's API, you need to register one or multiple endpoints at the proxy. The registration is an easy task and can be seen best on the [oCIS Hello example](https://owncloud.dev/extensions/ocis_hello/running/#configure-and-start-ocis).

As files in ownCloud must always stay private (unless you share them with your friends or coworkers), requests to oCIS have an authenticated user context. This user context is also available to your extension and can be used to interact with the user's files. How to get the user context and authentication can be seen on the [oCIS Hello example](https://owncloud.dev/extensions/ocis_hello/settings/#account-uuid).

### Storage

oCIS leverages the CS3 APIs and [CS3 REVA](https://github.com/cs3org/reva) as a storage system because it offers a very flexible setup and supports a variety of storage backends like EOS, S3 and of course your local hard drive. REVA makes it easy to support more storage backends as needed.

If you need to interact with files directly, you have the full power of the [CS3 APIs](https://cs3org.github.io/cs3apis/) in your hand. With the user context and the user's authentication token, which your extensions gets from the proxy, your extension can make these request in behalf of the user.

If your extension needs to store persistent data which is not supposed to live in the user's home folder, there is also a so-called metadata storage, intended for exactly that purpose. You should always use the metadata storage in favor of the local filesystem for persistent files, because your extension will then automatically use the storage backend the oCIS admin decides to use. For a temporary cache it is perfectly fine to use the local filesystem.
