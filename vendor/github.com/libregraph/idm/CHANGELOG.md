# CHANGELOG

## Unreleased



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

