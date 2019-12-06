Change: Better log level handling within micro

Currently every log message from the micro internals are logged with the info
level, we really need to respect the proper defined log level within our log
wrapper package.

https://github.com/owncloud/ocis-pkg/issues/2
