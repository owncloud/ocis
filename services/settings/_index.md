---
title: "Settings"
date: 2018-05-02T00:00:00+00:00
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/services/settings
geekdocFilePath: _index.md
geekdocCollapseSection: true
---

## Abstract

When using oCIS, the requirement to store settings arises. This extension provides functionality
for other extensions to register new settings within oCIS. It is responsible for storing the respective
settings values as well.

{{< mermaid class="text-center">}}
graph TD
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

**Settings usage:**
- Extensions can query ocis-settings for *settings values* of a user.

## Table of Contents

{{< toc-tree >}}
