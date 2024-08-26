#!/bin/bash

#mkdir -p /drone/src/vendor-bin/behat
#cp /tmp/vendor-bin/behat/composer.json /drone/src/vendor-bin/behat/composer.json

git config --global advice.detachedHead false

## CONFIGURE TEST

if [ "$TEST_SOURCE" = "oc10" ]
then
    export ACCEPTANCE_TEST_TYPE='core-api'
    if [ "$STORAGE_DRIVER" = "ocis" ]
    then
        export OCIS_REVA_DATA_ROOT=''
        export BEHAT_FILTER_TAGS='~@skipOnOcis-OCIS-Storage'
        export OCIS_SKELETON_STRATEGY='upload'
        export EXPECTED_FAILURES_FILE='/drone/src/tests/acceptance/expected-failures-API-on-OCIS-storage.md'
    elif [ "$STORAGE_DRIVER" = "s3ng" ]
    then
        export BEHAT_FILTER_TAGS='~@skip&&~@skipOnOcis-S3NG-Storage'
        export OCIS_REVA_DATA_ROOT=''
        export OCIS_SKELETON_STRATEGY='upload'
    else
        echo "non existing STORAGE selected"
        exit 1
    fi

    unset BEHAT_SUITE

elif [ "$TEST_SOURCE" = "ocis" ]
then

    if [ "$STORAGE_DRIVER" = "ocis" ]
    then
        export BEHAT_FILTER_TAGS='~@skip&&~@skipOnOcis-OCIS-Storage'
        export OCIS_REVA_DATA_ROOT=''
        export OCIS_SKELETON_STRATEGY='upload'
    elif [ "$STORAGE_DRIVER" = "s3ng" ]
    then
        export BEHAT_FILTER_TAGS='~@skip&&~@skipOnOcis-S3NG-Storage'
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
    # allow running without filters if its a feature

    unset BEHAT_FILTER_TAGS
    unset DIVIDE_INTO_NUM_PARTS
    unset RUN_PART
    unset EXPECTED_FAILURES_FILE
else
    unset BEHAT_FEATURE
fi

## RUN TEST

if [[ -z "$TEST_SOURCE" ]]
then
    echo "non existing TEST_SOURCE selected"
    exit 1
else
    sleep 10
    make -C $OCIS_ROOT test-acceptance-api
fi

chmod -R 777 vendor-bin/**/vendor vendor-bin/**/composer.lock tests/acceptance/output
