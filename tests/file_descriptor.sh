#!/bin/bash

IFS='
';
INITIAL_FILE_DESCRIPTOR_STORE=${INITIAL_FILE_DESCRIPTOR_STORE:-"initial.txt"}
CURRENT_FILE_DESCRIPTOR_STORE=${CURRENT_FILE_DESCRIPTOR_STORE:-"final.txt"}

declare -A initial

while IFS="" read -r p || [ -n "$p" ]
do
  SERVICE=$(echo $p | cut -d" " -f 1)
  COUNT_DESCRIPTORS=`echo $p | cut -d" " -f 2`
  echo $COUNT_DESCRIPTORS
  initial[$SERVICE]=$COUNT_DESCRIPTORS
done < $INITIAL_FILE_DESCRIPTOR_STORE

while :
do
  while IFS="" read -r p || [ -n "$p" ]
  do
    SERVICE=$(echo $p | cut -d" " -f 1)
    COUNT_DESCRIPTORS=`echo $p | cut -d" " -f 2`

    echo ${initial[$SERVICE]}
    limit=$(( ${initial[$SERVICE]} + 5 ))

    if [ $COUNT_DESCRIPTORS -gt $limit ] 
    then
      echo "File descriptor count exceeded the expected threshold limit"
      echo "Exiting the tests"
      pkill -f behat
      exit 1
    fi
  done < $CURRENT_FILE_DESCRIPTOR_STORE
  sleep 5
done
