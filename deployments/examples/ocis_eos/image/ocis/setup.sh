#!/bin/bash

set -x

echo "----- [ocis] LDAP setup -----"
authconfig --enableldap --enableldapauth --ldapserver=${EOS_LDAP_HOST} --ldapbasedn="dc=ocis,dc=test" --update
sed -i "s/#binddn cn=.*/binddn ${LDAP_BINDDN}/" /etc/nslcd.conf
sed -i "s/#bindpw .*/bindpw ${LDAP_BINDPW}/" /etc/nslcd.conf
# start in debug mode
nslcd -d &

# echo "----- [ocis] eos setup -----"
# eos -r 0 0 -b vid set membership daemon -uids adm
# eos -r 0 0 -b vid set membership daemon -gids adm
# eos -r 0 0 -b vid set membership daemon +sudo
# eos -r 0 0 -b vid set map -unix "<pwd>" vuid:0 vgid:0
# # eos -r 0 0 -b vid add gateway ocis.testnet
# eos -r 0 0 -b vid add gateway ocis

# todo start ocis as daemon not as root
# eos -r 0 0 -b vid set membership root -uids adm
# eos -r 0 0 -b vid set membership root -gids adm
# eos -r 0 0 -b vid set membership root +sudo
