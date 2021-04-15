---
title: "Ports"
date: 2018-05-02T00:00:00+00:00
weight: 41
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/extensions/storage
geekdocFilePath: ports.md
geekdocCollapseSection: true
---

Currently, every service needs to be configured with a port so oCIS can start them on localhost. We will automate this by using a service registry for more services, until eventually only the proxy has to be configured with a public port.

For now, the storage service uses these ports to preconfigure those services:

| port      | service                                       |
|-----------|-----------------------------------------------|
| 9109      | health, used by cli?                          |
| 9140      | frontend                                      |
| 9141      | frontend debug                                |
| 9142      | gateway                                       |
| 9143      | gateway debug                                 |
| 9144      | users                                         |
| 9145      | users debug                                   |
| 9146      | authbasic                                     |
| 9147      | authbasic debug                               |
| 9148      | authbearer                                    |
| 9149      | authbearer debug                              |
| 9150      | sharing                                       |
| 9151      | sharing debug                                 |
| 9154      | storage home grpc                             |
| 9155      | storage home http                             |
| 9156      | storage home debug                            |
| 9157      | storage users grpc                            |
| 9158      | storage users http                            |
| 9159      | storage users debug                           |
| 9160      | groups                                        |
| 9161      | groups debug                                  |
| 9178      | storage public link                           |
| 9179      | storage public link data                      |
| 9215      | storage meta grpc                             |
| 9216      | storage meta http                             |
| 9217      | storage meta debug                            |
