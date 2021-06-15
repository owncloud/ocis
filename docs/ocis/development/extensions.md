---
title: "Extensions"
date: 2020-02-27T20:35:00+01:00
weight: 40
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/development
geekdocFilePath: extensions.md
---

oCIS is all about files, sync and share - but most of the time there is some more you want to do with your files, eg. have a different view on your photo collection or edit your office file in an online file editor. ownCloud 10 faced the same problems and tries to solve them with so called "apps". These can extend the functionality of ownCloud 10 in a wide range. oCIS has a similar concept to be extended in its functionality: extensions. Because oCIS is very different in its architecture compared to ownCloud 10 this also applies to the extensions. An extension is basically any running code which satisfies basic interfaces and provides functionality to oCIS and its users. Because extensions are micro services you can choose almost any programming language you like. (even if for some languages the task might be a lot easier - but for ownCloud 10 it was nearly impossible to use a different programming language than php).

We will now introduce you to concepts, possibilities and ... of the oCIS extension system.



### Extensions outside of the oCIS Monorepo
Technically every service in oCIS is an extension, even if oCIS would not really work without them. So there are plenty of extensions which you can have a look at in the oCIS Monorepo.

Besides these "default" extensions there are also some more extensions:

- Hello
- WOPI server

Difference between extensions living inside the oCIS monorepo and in its own repo are:
- extensions inside the oCIS monorepo are all written in Go, whereas other extensions may choose the programming language freely
- extension inside the oCIS monorepo share tooling, whereas other extensions may use different tooling (eg. a different CI system)
- extensions inside the oCIS monorepo will be all build into one binary and started with the `ocis server` command. Other extensions must be started individually.



### Web
- App
- Design System

### Settings
An extension likely has some behaviour which the user can configure. Fundamental configuration will often be done by admins during deploy time in configuration files or by environment variables. But for other settings - which are supposed to change more often or which are even user specific - this is not a viable way. Therefore you need to offer the users a UI where they can configure your extension to their liking. Because implementing something like this is a repetitive task among extensions, oCIS already offers the settings extensions which does that for your extension. Your extension just needs to register settings bundles and permissions and read the current values from the settings service. You can read more on that on [Settings Extensions]({{< ref "../../extensions/settings" >}}) and see how [oCIS Hello uses settings]().

### Proxy
The Proxy acts as an API gateway and is the single connection point where all request from users and devices are sent to.

In order that requests can reach your extensions' api you need to register one or multiple endpoints at the proxy. This is currently done by manually adding routes to a config file but will be done dynamically by service discovery in the future. The registration is a easy task and can be seen best on the [oCIS Hello example]().

As files you store in your ownCloud must always stay private unless you share them with your friends or coworkers, requests to oCIS have always a user context. This user context is also available to your extension and can be used to interact with the users' files. How to get the user context and authentication can be seen bet on the [oCIS Hello example]() or the [WOPI server example]().



### Storage
oCIS leverages the CS3 APIs and REVA as a storage system because it offers a very flexible setup and supports a variety of storage backends like EOS, S3 and of course your local hard disk. REVA makes it easy to support more storage backends as needed.

If you need to interact with files, you have the full power of the [CS3 APIs]() in your hand. With the user context and authorization your extensions gets from the proxy you can make these request in behalf of the user.

If your extension needs to store data which is not supposed to life in the users home folder, there is also so called metadata storage which can be used for that purpose.

They main point you should get about storage in an oCIS extension is that you should never use the filesystem, but always use the CS3 APIs.

### Deployment


### Outstanding development

- dynamic registration web / proxy
- entitle service to act in behalf of users without request
- access to metadata storage
- Events (eg. callbacks on file create/change/delete events)
