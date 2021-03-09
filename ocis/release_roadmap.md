---
title: "Release Roadmap"
date: 2020-12-16T20:35:00+01:00
weight: 0
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis
geekdocFilePath: release_roadmap.md
---

You may have asked yourself why there are major version 1 tags in our GitHub repository but the Readme still states `ownCloud Infinite Scale is currently in a technical preview state. It will be subject to a lot of changes and is not yet ready for general production deployments.`. How can that be if its a major version 1?

Our initial and also our current plan is to stick to SemVer as versioning scheme. But sometimes there are other factors which cross your plans. Therefore we started releasing oCIS with version `1.0.0 Tech Preview`.

## ownCloud Infinite Scale 1.x technology preview releases

All oCIS releases within major version 1 will be handled as technology previews. There will be no supported releases in terms of us guaranteeing production readiness. We will do releases every 3 weeks, and they will sometimes only include bugfixes, but also introduce new features or optimizations.

We will be fixing bugs if you report them and truly appreciate every report and contribution. Depending on the individual case, we will publish bugfix releases or add the fix to the next minor release.

We are going to stick to major version 1 until we feel confident about running oCIS in production environments. As a consequence of this we cannot raise the major version, like SemVer requires it, even if we need to introduce breaking changes. We will do our best to avoid breaking changes. If there is no way to circumvent this, we will add an automatic migration or at least point out manual migration steps, since we as oCIS developers are already using oCIS on a personal basis. The best place to see if a breaking change happens is our changelog which is available for every release. If things are not working out for you please contact us immediately. We want to know about this and solve it for you.

It isn't our intention to scare you with our addendum "Tech Preview". We want you to have a clear picture of what you can expect from oCIS. You could take it as a disclaimer or even compare it to running an Linux kernel in alpha stage. It can be very pleasing to be on the latest codebase but you could also find yourself with a lot of problems arising because of that.

You clearly can expect a totally new experience of file-sync and share with oCIS and we want you to use it now - but with understanding and caution.

## ownCloud Infinite Scale 2.x general availability releases

Starting with oCIS 2.0.0 release we will strictly stick to SemVer, just as we do right now for ownCloud Server. The general availability release will also mean that we can recommend oCIS warmly to everyone. Use it to store your precious family pictures or you confidential company data!
