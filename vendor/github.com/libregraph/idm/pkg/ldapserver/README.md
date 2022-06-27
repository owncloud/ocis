# LDAP server library for Golang

This library provides LDAP server v3 functionality for the GO programming
language.

The server implementation is based on github.com/nmcclain/ldap and is enhanced
so it can be used together with github.com/go-ldap/ldap/v3.

From the server perspective, all of RFC4510 is implemented except:

4.5.1.3. SearchRequest.derefAliases
4.5.1.5. SearchRequest.timeLimit
4.5.1.6. SearchRequest.typesOnly
4.14. StartTLS Operation

The purpose of this library is not a general LDAP server implementation but to
provide enough of an LDAP server for Kopano compatible identity management.

## License

See `LICENSE.txt` for licensing information of this module.
