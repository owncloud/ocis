#!/bin/bash
sleep 10
IFS='
';
INITIAL_FILE_DESCRIPTOR_STORE=${INITIAL_FILE_DESCRIPTOR_STORE:-"initial.txt"}
CURRENT_FILE_DESCRIPTOR_STORE=${CURRENT_FILE_DESCRIPTOR_STORE:-"final.txt"}

for i in `ps xao pid,cmd | grep ocis | grep -v grep`
do 
  echo $i
  trimmed=$(echo $i | xargs)
  echo $trimmed
  PID=$(echo $trimmed | cut -d" " -f 1)
  echo $PID
  SERVICE=$(echo $trimmed | cut -d" " -f 5)
  COUNT_DESCRIPTORS=`lsof -p $PID 2> /dev/null | wc -l`
  echo $SERVICE $COUNT_DESCRIPTORS
  echo $SERVICE $COUNT_DESCRIPTORS >> $INITIAL_FILE_DESCRIPTOR_STORE
done
cat $INITIAL_FILE_DESCRIPTOR_STORE
cat $INITIAL_FILE_DESCRIPTOR_STORE > $CURRENT_FILE_DESCRIPTOR_STORE
cat $CURRENT_FILE_DESCRIPTOR_STORE
while :
do
  sleep 5
  cat /dev/null > $CURRENT_FILE_DESCRIPTOR_STORE
  for i in `ps xao pid,cmd | grep ocis | grep -v grep`
  do
    trimmed=$(echo $i | xargs)
    PID=$(echo $trimmed | cut -d" " -f 1)
    SERVICE=$(echo $trimmed | cut -d" " -f 3)
    COUNT_DESCRIPTORS=`lsof -p $PID 2> /dev/null | wc -l`
    echo $SERVICE $COUNT_DESCRIPTORS >> $CURRENT_FILE_DESCRIPTOR_STORE
  done
done

