# IDM

The IDM service provides a minimal LDAP Service, based on [Libregraph idm](https://github.com/libregraph/idm), for oCIS. It is started as part of the default configuration and serves as a central place for storing user and group information.

It is mainly targeted at small oCIS installations. For larger setups it is recommended to replace IDM with a “real” LDAP server or to switch to an external identity management solution.

IDM listens on port 9235 by default. In the default configuration it only accepts TLS-protected connections (LDAPS). The BaseDN of the LDAP tree is `o=libregraph-idm`. IDM gives LDAP write permissions to a single user (DN: `uid=libregraph,ou=sysusers,o=libregraph-idm`). Any other authenticated user has read-only access. IDM stores its data in a boltdb file `idm/ocis.boltdb` inside the oCIS base data directory.

Note: IDM is limited in its functionality. It only supports a subset of the LDAP operations (namely `BIND`, `SEARCH`, `ADD`, `MODIFY`, `DELETE`). Also, IDM currently does not do any schema verification (like. structural vs. auxiliary object classes, require and option attributes, syntax checks, …). Therefore it is not meant as a general purpose LDAP server.
