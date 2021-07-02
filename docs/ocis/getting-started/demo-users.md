---
title: "Demo Users"
date: 2020-02-27T20:35:00+01:00
weight: 0
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/getting-started
geekdocFilePath: demo-users.md
---

As long as oCIS is released as [technology preview]({{< ref "../release_roadmap#release_roadmap" >}}) it will come with default demo users. These enable you to do quick testing and developing.

Following users are available in the demo set:

| username | password      | email                | role  | groups                                                                  |
| -------- | ------------- | -------------------- | ----- | ----------------------------------------------------------------------- |
| admin    | admin         | admin@example.org    | admin | users                                                                   |
| einstein | relativity    | einstein@example.org | user  | users, philosophy-haters, physics-lovers, sailing-lovers, violin-haters |
| marie    | radioactivity | marie@example.org    | user  | users, physics-lovers, polonium-lovers, radium-lovers                   |
| moss     | vista         | moss@example.org     | admin | users                                                                   |
| richard  | superfluidity | richard@example.org  | user  | users, philosophy-haters, physics-lovers, quantum-lovers                |

You may also want to run oCIS with only your custom users by [deleting the demo users]({{< ref "../deployment#delete-demo-users" >}}).
