#!/usr/bin/env bash
SCRIPT=`realpath $0`
SCRIPTPATH=`dirname $SCRIPT`

# start the server
#OCIS_LOG_LEVEL=error PROXY_ENABLE_BASIC_AUTH=true ocis server
#sleep 30

# run the tests
python3 CalDAVTester/testcaldav.py --ssl --print-details-onfail --basedir "$SCRIPTPATH/caldavtest/" \
 "CalDAV/caldavIOP.xml"

RESULT=$?

exit $RESULT
