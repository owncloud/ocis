#!/usr/bin/env bash

SHARE_ENDPOINT="ocs/v2.php/apps/files_sharing/api/v1/shares"

# envs
ENV="SPACE_ID="
# get space id
SPACE_ID=$(curl -ks -uadmin:admin "${TEST_SERVER_URL}/graph/v1.0/me/drives" | grep -Po '(?<="webDavUrl":").*?(?="})' | cut -d"/" -f6 | sed "s/\\$/\\\\$/g")
ENV+=${SPACE_ID}

# create a folder
curl -ks -ueinstein:relativity -X MKCOL "${TEST_SERVER_URL}/remote.php/webdav/new_folder"

SHARE_ID=$(curl -ks -ueinstein:relativity "${TEST_SERVER_URL}/${SHARE_ENDPOINT}" -d "path=/new_folder&shareType=0&permissions=15&name=new_folder&shareWith=admin" | grep -oP "(?<=<id>).*(?=</id>)")
# accept share
if [ ! -z "${SHARE_ID}" ];
then
  curl -XPOST -ks -uadmin:admin "${TEST_SERVER_URL}/${SHARE_ENDPOINT}/pending/${SHARE_ID}"
fi

# create public share
PUBLIC_TOKEN=$(curl -ks -ueinstein:relativity "${TEST_SERVER_URL}/${SHARE_ENDPOINT}" -d "path=/new_folder&shareType=3&permissions=15&name=new_folder" | grep -oP "(?<=<token>).*(?=</token>)")
ENV+="\nPUBLIC_TOKEN="
ENV+=${PUBLIC_TOKEN}

# create an .env file in the repo root dir
echo -e $ENV >> .env
