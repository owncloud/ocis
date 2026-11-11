Bugfix: Support Collabora Online 26.x in the ocis_full deployment example

collabora/code 26.x ships without a shell and without the `coolconfig` helper the
ocis_full example relied on to start the container and generate its WOPI proof
key, so the example failed to start and, once running, would reject every
document open with a 500. The example now uses the image's default entrypoint
and a one-shot init container that generates and persists the WOPI proof key
before Collabora starts.

https://github.com/owncloud/ocis/pull/12794
