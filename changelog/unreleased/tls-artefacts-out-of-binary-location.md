Enhancement: Use system default location to store TLS artefacts.

This used to default to the current location of the binary, which is not ideal after a first run as it leaves traces behind. It now uses the system's location for artefacts with the help of https://golang.org/pkg/os/#UserConfigDir.

https://github.com/owncloud/ocis/pull/2129
