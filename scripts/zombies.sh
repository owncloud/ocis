#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

[[ -z "$OFFSET" ]] && OFFSET=1
[[ -z "$COUNT" ]] && "$((COUNT = 1000))"

: "$((i=OFFSET))"
: "$((end=OFFSET+COUNT))"
while [ "$((i <= end))" -ne 0 ]
do
    echo "creating zombie $i"
    curl -X POST 'https://localhost:9200/ocs/v1.php/cloud/users' -k -u admin:admin -d userid="zombie$i" -d password="zombie" -d email="zombie$i@example.org"
    #$DIR/../ocis/bin/ocis accounts add --preferred-name zombie$i --on-premises-sam-account-name zombie$i --mail zombie$i@example.org
    : "$((i = i + 1))"
done
