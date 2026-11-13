Change: Remove the EULA from the release pipeline

Release builds attached an End User License Agreement PDF to the published
assets, and the README stated that stable builds provided by ownCloud GmbH are
subject to that agreement. A non OSS agreement on the distributed binaries
contradicts the Apache 2.0 licensing of Infinite Scale.

We removed the EULA from the release pipeline. The release workflow no longer
copies the PDF into the release assets, the release audit no longer expects it,
the PDF itself was dropped from the repository and the README section was
removed. Releases now ship under Apache 2.0 only.

https://github.com/owncloud/ocis/pull/13262
