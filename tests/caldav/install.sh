#!/usr/bin/env bash
set -e
SCRIPT=`realpath $0`
SCRIPTPATH=`dirname $SCRIPT`

cd "$SCRIPTPATH"
if [ ! -f CalDAVTester/testcaldav.py ]; then
    git clone https://github.com/DeepDiver1975/ccs-caldavtester.git -b python3 CalDAVTester
    cd CalDAVTester
    python3 -m venv venv
    source venv/bin/activate
    pip install -r requirements.txt
fi

# create test user
# TODO: call ocis calendar create command
#cd "$SCRIPTPATH/../../../../../"
#OC_PASS=user01 php occ user:add --password-from-env user01
#php occ dav:create-calendar user01 calendar
#php occ dav:create-calendar user01 shared
#OC_PASS=user02 php occ user:add --password-from-env user02
#php occ dav:create-calendar user02 calendar
#cd "$SCRIPTPATH/../../../../../"
