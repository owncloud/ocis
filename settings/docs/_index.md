* * *

title: "Settings"
date: 2018-05-02T00:00:00+00:00
weight: 10
geekdocRepo: <https://github.com/owncloud/ocis-settings>
geekdocEditPath: edit/master/docs

## geekdocFilePath: \_index.md

## Abstract

When using oCIS, the requirement to store settings arises. This extension provides functionality
for other extensions to register new settings within oCIS. It is responsible for storing the respective
settings values as well.

For ease of use, this extension provides an ocis-web extension which allows users to change their settings values.
Please refer to the [ocis-web extension docs](https://owncloud.github.io/ocis/extensions/#external-phoenix-apps)
for running ocis-web extensions.

{{&lt; mermaid class="text-center">}}
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
{{&lt; /mermaid >}}

The diagram shows how the settings service integrates into oCIS:

**Settings management:**

-   oCIS extensions can register _settings bundles_ with the ocis-settings service.
-   The settings frontend can be plugged into ocis-web, showing forms for changing _settings values_ as a user.
    The forms are generated from the registered _settings bundles_.

**Settings usage:**

-   Extensions can query ocis-settings for _settings values_ of a user.
-   The ownCloud SDK, used as a data abstraction layer for ocis-web, will query ocis-settings for _settings values_ of a user,
    if it's available. The SDK uses sensible defaults when ocis-settings is not part of the setup.

For compatibility with ownCloud 10, a migration of ownCloud 10 settings into the storage of ocis-settings will be available.
