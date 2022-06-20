---
title: "Settings"
date: 2018-05-02T00:00:00+00:00
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/extensions/settings
geekdocFilePath: _index.md
geekdocCollapseSection: true
---

## Abstract

When using oCIS, the requirement to store settings arises. This extension provides functionality
for other extensions to register new settings within oCIS. It is responsible for storing the respective
settings values as well.

For ease of use, this extension provides an ocis-web extension which allows users to change their settings values.
Please refer to the [ocis-web extension docs]({{< ref "../../ocis/development/extensions/#external-ownCloud-Web-apps" >}})
for running ocis-web extensions.

{{< mermaid class="text-center">}}
graph TD
    subgraph ow[ocis-web]
        ows[ocis-web-settings]
        owc[ocis-web-core]
    end
    ows ---|"listSettingsBundles(),<br>saveSettingsValue(value)"| os[ocis-settings]
    owc ---|"listSettingsValues()"| sdk[oC SDK]
    sdk --- sdks{ocis-settings<br>available?}
    sdks ---|"yes"| os
    sdks ---|"no"| defaults[Use set of<br>default values]
    oa[oCIS extensions<br>e.g. ocis-accounts] ---|"saveSettingsBundle(bundle)"| os
{{< /mermaid >}}

The diagram shows how the settings service integrates into oCIS:

**Settings management:**
- oCIS extensions can register *settings bundles* with the ocis-settings service.
- The settings frontend can be plugged into ocis-web, showing forms for changing *settings values* as a user.
The forms are generated from the registered *settings bundles*.

**Settings usage:**
- Extensions can query ocis-settings for *settings values* of a user.
- The ownCloud SDK, used as a data abstraction layer for ocis-web, will query ocis-settings for *settings values* of a user,
if it's available. The SDK uses sensible defaults when ocis-settings is not part of the setup.

For compatibility with ownCloud 10, a migration of ownCloud 10 settings into the storage of ocis-settings will be available.
