---
title: Additinal Information
date: 2018-05-02T00:00:00+00:00
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/services/general-info
geekdocFilePath: additional-information.md
geekdocCollapseSection: true
---

This section contains information on general topics

## GRPC Maximum message size

Ocis is using grpc for inter service communication. When having a folder with a lot of files (25000+, size doesn't matter) and doing a `PROPFIND` on the folder, the server will run into errors because the
grpc message body gets to big. We introduced the envvar `OCIS_GRPC_MAX_RECEIVED_MESSAGE_SIZE` to raise the max size for the grpc body.

NOTE: With a certain amount of files even raising the grpc message size will not suffice as the requests will run into network timeouts. Also generally the more files are in a folder, the longer it will take to load.

It is recommended to use `OCIS_GRPC_MAX_RECEIVED_MESSAGE_SIZE` only temporary to copy files out of the folder (e.g. via web ui) and use the default value in general.
