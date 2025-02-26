## LibreGraph Identity Management

The LibreGraph Identity Management provides a LDAP server, which is easy to configure, does not have external dependencies and is tailored to work perfectly with other LibreGraph software.

The goal is that everyone who does not already have or needs an LDAP server, uses IDM.

Thus, IDM is a (currently read-only) drop in replacement for an existing LDAP server and does provide an LDAP interface if none is there already. IDM uses hard coded indexes and supports LDAP search, bind and unbind operations.

### Running idmd from a source build

Until packages and containers for more environments are available it is the easiest to just create a local build of `idmd`. For this just run `make`.

IDM uses a mixture of environment variables and parameters for configuration and needs to be at least passed a the location of an individual ldif file or a directory containing multiple ldif files.

```bash
$ ./idmd serve --ldif-main ./export.ldif
INFO[0000] LDAP listener started                         listen_addr="127.0.0.1:10389"
INFO[0000] ready
```

### Configuration

The default base DN of IDM is `dc=lg,dc=local`. There is usually no need to change, it if you don't use the LDAP data for anything else. The value needs to match what the clients have configured. Similarly, the default mail domain is `lg.local`.

Both values can be changed by passing `--ldap-base-dn` or `--ldif-template-default-mail-domain` respectively.

IDM uses ldif files for its data source and those files, the location of these files needs to be passed at startup using the `--ldif-main` parameter.

#### Adding a service user for LDAP access

By default IDM does not have any users and anonymous bind is disabled. You can enable anonymous bind support for local requests by passing `--ldap-allow-local-anonymous` when running `idmd`. Alternatively a service user can be specified in the following way:

```bash
cat <<EOF > ./config.ldif
dn: cn=readonly,{{.BaseDN}}
cn: readonly
description: LDAP read only service user
objectClass: simpleSecurityObject
objectClass: organizationalRole
userPassword: readonly
EOF
```

And then passed as an additional parameter when starting `idmd` by passing `--ldif-config ./config.ldif`. The `config.ldif` is for service users only and the data in there is used for bind requests only, but never returned for search requests.

#### Add users to the ldap service

`idmd` serves all ldif files from the folder specified by `--ldif-main` (loaded in lexical order and parsed as templates). Whenever any of the ldif files are changed, added or removed, make sure to restart `idmd`.

`idmd` listens on `127.0.0.1:10389` by default and does not ship with any default users. Example configuration can be found in the [scripts directory](https://github.com/libregraph/idm/tree/master/scripts) of this repository.

##### Add new users using the `gen newusers` command

IDM provides a way to create ldif data for new users using batch mode similar to the unix `newusers` command using the following standard password file format:

```bash
uid:userPassword:uidNumber:gidNumber:cn,[mail][,mailAlternateAddress...]:ignored:ignored
```

For example, like this:

```bash
cat << EOF | ./idmd gen newusers - --min-password-strength=4 > ./ldif/50-users.ldif
jonas:passwordOfJonas123:::Jonas Brekke,jonas@lg.local::
timmothy:passwordOfTimmothy456:::Timmothy Schöwalter::
EOF
```

This outputs an LDIF template file which you can modify as needed. When done run restart `idmd` to make the new users available. Keep in mind that some of the attributes must be unique.

##### Replace existing OpenLDAP with IDM

On the LDAP server export all its data using `slapcat` and write the resulting ldif to for example `./ldif/10-main.ldif`. This is a drop in replacement and all what was in OpenLDAP is now also in IDM.

Either stop `slapd` and change the IDM configuration to listen where `slapd` used to listen or change the clients to connect to where `idmd` listens to migrate.

### Extra goodies

#### Template support

All ldif files loaded by IDM support template syntax as defined in https://golang.org/pkg/text/template to allow auto generation and replacement of various values. You can find example templates in the [scripts directory](https://github.com/libregraph/idm/tree/master/scripts) as well. All the `gen` commands output template syntax if applicable.

#### Generate secure password hash using the `gen passwd` command

IDM supports secure password hashing using ARGON2. To create such password hashes either use `gen newusers` or the interactive `gen passwd` which is very similar to `slappasswd` from OpenLDAP.

```bash
./idmd gen passwd
New password:
Re-enter new password:
{ARGON2}$argon2id$v=19$m=65536,t=1,p=2$MaB5gX2BI484dATbGFyEIg$h2X8rbPowzZ/Exsz4W20Z/Zk54C30YnY+YbivSIRpcI
```

#### Test IDM

Since `idmd` provides a standard LDAP interface, also standard LDAP tools can be used to interact with it for testing. Run `apt install ldap-utils` to install LDAP commandline tools.

```bash
ldapsearch -x -H ldap://127.0.0.1:10389 -b "dc=lg,dc=local" -D "cn=readonly,dc=lg,dc=local" -w 'readonly'
```
