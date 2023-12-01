#!/bin/sh

apk add --no-cache openssh-client sshpass

sshpass -p "$SSH_PASSWORD" ssh -o StrictHostKeyChecking=no "$SSH_USERNAME@$SSH_SERVER" "$SSH_COMMAND"
