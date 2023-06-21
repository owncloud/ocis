---
title: "Release Roadmap"
date: 2020-12-16T20:35:00+01:00
weight: 0
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis
geekdocFilePath: release_roadmap.md
---

## Semantic Versioning
Since oCIS 2.0.0 release we will strictly stick to SemVer, just as we do right now for ownCloud Server. The general availability release will also mean that we can recommend oCIS warmly to everyone. Use it to store your precious family pictures or you confidential company data!

## Community Release
Releases of the type "Community" contain the latest greatest features, but are not as stable as an Enterprise release. Community Releases are characterized as followed:
- Cycle: After each sprint (usually every 3 weeks).
- No known blockers. Data loss is not expected and is a bug.
- Can be used in production for specific purposes approved by ownCloud.
- Issues are fixed based on "best effort".
- Supported by community (no support by ownCloud).

## Enterprise Releases
Releases of the type "Enterprise" are Releases that received our full quality assurance cycle but are released less often compared to Community Releases. Enterprise Releases are characterized as followed:
- Cycle: About 2 times per year
- The release is available for all customers to use under all documented circumstances.
- Issues will be fixed according to our SLAs.
- Full QA cycle incl. all clients

## General Note
The release types "Enterprise" and "Community" are not dependent to the major number of the release version. This means if the major version number is counted up (2.x.x -> 3.0.0) this does not necessarily imply an Enterprise Release. Sticking to the rules of Semantic Versioning, the increment of the major number means that there are incompatibe API changes, but is says nothing about the release type "Enterprise" or "Community".
