# CHANGELOG

## Unreleased



## v0.5.0 (2024-04-17)

- Bump github.com/go-ldap/ldap/v3 from 3.4.7 to 3.4.8
- Bump github.com/go-ldap/ldap/v3 from 3.4.6 to 3.4.7
- Bump google.golang.org/protobuf from 1.32.0 to 1.33.0
- Bump golangci-lint action to v4
- Remove deprecated linters from config
- Bump required go version to 1.21 + go mod tidy
- Bump github.com/prometheus/client_golang from 1.18.0 to 1.19.0
- Bump go.etcd.io/bbolt from 1.3.8 to 1.3.9
- Bump github.com/prometheus/client_golang from 1.17.0 to 1.18.0
- Bump github.com/go-logr/logr from 1.3.0 to 1.4.1
- Bump golang.org/x/crypto from 0.14.0 to 0.17.0
- Fix the DN comparison in a ServerFilterScope
- Bump golang.org/x/text from 0.13.0 to 0.14.0
- Bump github.com/spf13/cobra from 1.7.0 to 1.8.0
- Bump github.com/go-logr/logr from 1.2.4 to 1.3.0
- Bump go.etcd.io/bbolt from 1.3.7 to 1.3.8
- Bump github.com/alexedwards/argon2id
- Bump github.com/prometheus/client_golang from 1.16.0 to 1.17.0
- Bump github.com/prometheus/client_golang from 1.15.1 to 1.16.0
- Bump github.com/go-ldap/ldap/v3 from 3.4.5 to 3.4.6
- Bump github.com/go-asn1-ber/asn1-ber from 1.5.4 to 1.5.5
- Bump golang.org/x/text from 0.12.0 to 0.13.0
- Bump golang.org/x/text from 0.11.0 to 0.12.0
- Bump golang.org/x/text from 0.10.0 to 0.11.0
- Bump golang.org/x/text from 0.9.0 to 0.10.0
- Bump github.com/go-ldap/ldap/v3 from 3.4.4 to 3.4.5
- Bump github.com/sirupsen/logrus from 1.9.2 to 1.9.3
- Bump github.com/sirupsen/logrus from 1.9.1 to 1.9.2
- Bump github.com/sirupsen/logrus from 1.9.0 to 1.9.1
- Bump github.com/prometheus/client_golang from 1.15.0 to 1.15.1
- Bump github.com/prometheus/client_golang from 1.14.0 to 1.15.0
- Bump golang.org/x/text from 0.8.0 to 0.9.0
- Bump github.com/spf13/cobra from 1.6.1 to 1.7.0
- Bump github.com/go-logr/logr from 1.2.3 to 1.2.4
- Bump golang.org/x/text from 0.7.0 to 0.8.0
- Fix ModifyDN operation in boltdb backend
- Add support for ModifyDN
- Fix error behaviour when receiving unsupported operation
- Bump golang.org/x/text from 0.6.0 to 0.7.0
- Bump required go version to 1.18 + go mod tidy
- Bump go.etcd.io/bbolt from 1.3.6 to 1.3.7
- Bump golang.org/x/text from 0.5.0 to 0.6.0
- Bump golang.org/x/text from 0.4.0 to 0.5.0


## v0.4.0 (2022-12-01)

- Migrate to Go rndm module from GitHub
- Bump github.com/prometheus/client_golang from 1.13.0 to 1.14.0
- Bump github.com/coreos/go-systemd/v22 from 22.4.0 to 22.5.0
- Bump github.com/spf13/cobra from 1.6.0 to 1.6.1
- Bump github.com/bombsimon/logrusr/v3 from 3.0.0 to 3.1.0
- Bump golang.org/x/text from 0.3.8 to 0.4.0
- Bump stash.kopano.io/kgol/rndm from 1.1.1 to 1.1.2
- Bump github.com/spf13/cobra from 1.5.0 to 1.6.0
- Bump golang.org/x/text from 0.3.7 to 0.3.8
- Bump github.com/coreos/go-systemd/v22 from 22.3.2 to 22.4.0
- Bump github.com/prometheus/client_golang from 1.12.2 to 1.13.0
- Bump github.com/go-ldap/ldap/v3 from 3.4.3 to 3.4.4
- Switch pkg/ldapserver to logr
- Set custom logger for go-ldap/ldap
- Bump github.com/sirupsen/logrus from 1.8.1 to 1.9.0
- Make substring filter case-insensitve
- Return proper error code when exceeding size limit
- Fix normalized DN attribute escaping
- Switch to go-ldap/ldap for filter (de-)compilation
- Fix DN compoare condition
- Switch github action to use `make test`
- improve DN comparison
- pass through unparsed DN
- Address a few linter complaints
- Implement modify password extended operation for boltdb backend
- Add backend plumbing for password modify extended operation
- pwexop: Add support of generating a random password
- Groundwork for password modify extended operation
- Bump github.com/spf13/cobra from 1.4.0 to 1.5.0
- Bump github.com/Songmu/prompter from 0.5.0 to 0.5.1
- Bump github.com/prometheus/client_golang from 1.12.1 to 1.12.2
- Bump github.com/go-ldap/ldap/v3 from 3.4.2 to 3.4.3
- Bump github.com/go-asn1-ber/asn1-ber from 1.5.3 to 1.5.4
- boltdb: Fix modify replace on RDN Attribute
- Bump github.com/spf13/cobra from 1.3.0 to 1.4.0
- boltdb bind: attributeTypes are case-insensitive
- Tone down debug logging
- encodeSearchDone might be called with nil doneControls
- Bump go-crypt to latest master
- Allow to disable go-crypt related code
- Fix build on Darwin
- Bump github.com/go-ldap/ldap/v3 from 3.4.1 to 3.4.2
- Cleanup logging in boltdb handler
- Bump github.com/prometheus/client_golang from 1.12.0 to 1.12.1
- Bump github.com/prometheus/client_golang from 1.11.0 to 1.12.0
- Introduce new parameter "ldap-admin-dn"
- Normalize BaseDN and BindDN
- LDAP Modify support for boltdb Handler
- Add utils to apply LDAP Modify Request on Entries
- Create ldapentry and ldapdn helper modules
- Add shortcut for normalizing DN string
- Parse and validate incoming LDAP Modify Requests
- fix typo
- boltdb: Add getEntryByID method
- boltdb: Make internal helper methods private
- Bump all unversioned dependencies to their latest code
- Implement Delete Support for boltdb Handler
- Parse and validate incoming LDAP Delete Requests
- Bump github.com/sirupsen/logrus from 1.6.0 to 1.8.1
- Bump github.com/spf13/cobra from 1.2.1 to 1.3.0
- Bump github.com/prometheus/client_golang from 0.9.3 to 1.11.0
- Bump golang.org/x/text from 0.3.5 to 0.3.7
- Initial LDAPAdd Support for the boltdb Handler
- LDAPAdd support for the backend handlers
- boltdb: Disallow adding an already existing Entry
- Bump github.com/spf13/cobra from 1.1.3 to 1.2.1
- Bump github.com/coreos/go-systemd/v22 from 22.3.0 to 22.3.2
- Enable dependabot for go modules
- Don't consider linter failures fatal
- Parse and validate incoming LDAP Add Requests
- Add basic plumbing for LDAP Add support
- Update to latest bbolt release
- Add some initial unit tests for boltdb backend ([#23](https://github.com/libregraph/idm/issues/23/))
- Tone down golangci-lint annotation to warnings ([#24](https://github.com/libregraph/idm/issues/24/))
- Add "boltdb export" subcommand
- Set a default log-level for the boltdb related subcommands
- Add ability to pass bolt.Options on database
- Add SimpleBind support for BoltDB
- Introduce a BoltDB based Database Handler
- Add options to use other backends than 'ldif'
- Add TLS support
- Adjust golangci-lint config
- Add initial Github Action as a starting point for CI
- Bump go-ldap to v3.4.1


## v0.3.0 (2021-09-29)

- Add new contributor/authors
- Fix loading of LDIF directory
- Change license to Apache License 2.0
- review comments
- review comments
- Update readme for usage from compiled binary
- Rewrite readme
- Remove Kopano wording from readme file
- Change copyright headers from Kopano to LibreGraph Authors
- Add A+C files
- Avoid duplicate index entries when using sub and pres
- Index mail pres and sub for mail attribute
- Cure potential panic in search without pagination
- Apply search BaseDN when returning values from index
- Introduce proper way to set defaults with option to override
- Remove Kopano specific defaults and naming for white label rename
- Rename public stuttering API functions
- Make internal ldappasswd package importable
- Make internal ldapserver package importable
- Remove Jenkinsfile to prepare for external CI
- Move project to github.com/libregraph/idm
- Add proper LICENSE file
- Add readme file


## v0.2.7 (2021-05-31)

- Skip loading nil LDIF entries


## v0.2.6 (2021-05-26)

- Use correct parts count for glibc2 CRYPT
- Ignore case when selecting password crypt algo
- Use absolute path for kill command


## v0.2.5 (2021-05-26)

- Fix file loading in newusers sub command


## v0.2.4 (2021-04-29)

- Fix missing variable in default LDIF main config template


## v0.2.3 (2021-04-29)

- Ensure to setup folders with correct permissions


## v0.2.2 (2021-04-29)

- Add setup step for systemd based startup


## v0.2.1 (2021-04-29)

- Fix refactoring error for hash based password checks


## v0.2.0 (2021-04-29)

- Move password hash functionality to internal module
- Add password strength checks
- Add gen passwd subcommand
- Consolidate password hashing functions
- Ignore commented lines when processing templates
- Support relative paths in templates
- Include demo LDIF generator script
- Only load files in templates which are in a base folder
- Unify config and commandline options
- Add binscript, systemd service and config
- Add reload support via SIGHUP
- Enable index and index lookup for objectClass only filters
- Add sub index support
- Add present index support
- Add proper license headers and origin reference
- Add some AD attributres for equality indexing


## v0.1.0 (2021-04-22)

- Improve string comparison performance
- Improve LDIF parse logging
- Prevent duplicates from multiple search equality index matches
- Allow negative search equality index match
- Add support to load LDIF data from folder
- Implement gen newusers sub command with LDIF output
- Add support for argon2 password hashing
- Implement more LDAP server metrics
- Add metrics support
- Fix LDAP server stats support
- Log LDAP close
- Remove unsupported Unbinder
- Fix debug log formatting
- Use better anonymous bind for standard compliance
- Add pprof support
- Implement difference between startup and runtime errors
- Add environment variables to set default config values
- Move serve command into sub folder to prepare for other sub commands
- Use template syntax in demo users generator
- Apply ldif template defaults
- Move LDIF template functionality into its own file
- Improve flexibility of template support
- Support setting current value in AutoIncrement template function
- Improve commandline parameter naming
- Use better names for example ldif
- Allow configuration of LDIF template defaults
- Add support to allow local anonymoys LDAP bind and search
- Load LDIF files with template support
- Actually allow LDIF middleware bind to succeed

