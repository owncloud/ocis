#!/bin/bash

git config --global advice.detachedHead false

## GET DEPENDENCIES

if cd $TESTING_DIR > /dev/null 2>&1
then
    git pull
else
    git clone -b master --depth=1 https://github.com/owncloud/testing.git $TESTING_DIR
fi

if cd $PATH_TO_CORE > /dev/null 2>&1
then
    git checkout $CORE_BRANCH
    git pull
    git checkout $CORE_COMMITID
else
    git clone -b $CORE_BRANCH --single-branch --no-tags https://github.com/owncloud/core.git $PATH_TO_CORE
    cd $PATH_TO_CORE
    git checkout $CORE_COMMITID
fi

## CONFIGURE TEST

if [ "$TEST_SOURCE" = "oc10" ]
then
    if [ "$STORAGE" = "owncloud" ]
    then
        export OCIS_REVA_DATA_ROOT='/srv/app/tmp/ocis/owncloud/data/'
        export DELETE_USER_DATA_CMD=''
        export BEHAT_FILTER_TAGS='~@notToImplementOnOCIS&&~@toImplementOnOCIS&&~comments-app-required&&~@federation-app-required&&~@notifications-app-required&&~systemtags-app-required&&~@local_storage&&~@skipOnOcis-OC-Storage'
        export OCIS_SKELETON_STRATEGY='copy'
        export EXPECTED_FAILURES_FILE='/drone/src/tests/acceptance/expected-failures-API-on-OWNCLOUD-storage.txt'
    elif [ "$STORAGE" = "ocis" ]
    then
        export OCIS_REVA_DATA_ROOT=''
        export DELETE_USER_DATA_CMD='rm -rf /srv/app/tmp/ocis/storage/users/nodes/root/* /srv/app/tmp/ocis/storage/users/nodes/*-*-*-*'
        export BEHAT_FILTER_TAGS='~@notToImplementOnOCIS&&~@toImplementOnOCIS&&~comments-app-required&&~@federation-app-required&&~@notifications-app-required&&~systemtags-app-required&&~@local_storage&&~@skipOnOcis-OCIS-Storage'
        export OCIS_SKELETON_STRATEGY='upload'
        export EXPECTED_FAILURES_FILE='/drone/src/tests/acceptance/expected-failures-API-on-OCIS-storage.txt'
    else
        echo "non existing STORAGE selected"
        exit 1
    fi

    unset BEHAT_SUITE

elif [ "$TEST_SOURCE" = "ocis" ]
then

    if [ "$STORAGE" = "owncloud" ]
    then
        export BEHAT_FILTER_TAGS='~@skipOnOcis-OC-Storage'
        export DELETE_USER_DATA_CMD=''
        export OCIS_REVA_DATA_ROOT='/srv/app/tmp/ocis/owncloud/data/'
        export OCIS_SKELETON_STRATEGY='copy'
    elif [ "$STORAGE" = "ocis" ]
    then
        export BEHAT_FILTER_TAGS='~@skipOnOcis-OCIS-Storage'
        export DELETE_USER_DATA_CMD='rm -rf /srv/app/tmp/ocis/storage/users/nodes/root/* /srv/app/tmp/ocis/storage/users/nodes/*-*-*-*'
        export OCIS_REVA_DATA_ROOT=''
        export OCIS_SKELETON_STRATEGY='upload'
    else
        echo "non existing storage selected"
        exit 1
    fi

    unset DIVIDE_INTO_NUM_PARTS
    unset RUN_PART
else
    echo "non existing TEST_SOURCE selected"
    exit 1
fi

if [ ! -z "$BEHAT_FEATURE" ]
then
    echo "feature selected: " + $BEHAT_FEATURE
    # allow to run without filters if its a feature

    unset BEHAT_FILTER_TAGS
    unset DIVIDE_INTO_NUM_PARTS
    unset RUN_PART
    unset EXPECTED_FAILURES_FILE
fi

## RUN TEST

if [ "$TEST_SOURCE" = "oc10" ]
then
    make -C /srv/app/testrunner test-acceptance-api
elif [ "$TEST_SOURCE" = "ocis" ]
then
    cd $OCIS_ROOT
    sleep 10
    make test-acceptance-api
else
    echo "non existing TEST_SOURCE selected"
    exit 1
fi

chmod -R 777 vendor-bin/**/vendor vendor-bin/**/composer.lock tests/acceptance/output
