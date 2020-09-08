#!/bin/sh

sshpass -p "$SSH_PASSWORD_HCLOUD" ssh -tt root@95.217.215.207 "docker exec -it mgm-master eos -r 0 0 rm -r /eos/dockertest/reva/users/$1"
