#!/bin/bash
DEPTH=${DEPTH:-3}
WIDTH=${WIDTH:-10}
OCIS_URL=${OCIS_URL:-https://localhost:9200}
ENDPOINT=${ENDPOINT:-/webdav}
FOLDER=${FOLDER:-w$WIDTH x d$DEPTH folders}
USER=${USER:-einstein}
PASSWORD=${PASSWORD:-relativity}
CURL_OPTS=${CURL_OPTS:--k}

COUNT=0
MAX=0

calc_max()
{
    for i in $(seq 1 $2);
    do {
        MAX=$(( MAX + ($1 ** i) ))
    }
    done
}

calc_max $WIDTH $DEPTH

create_tree()
{
    if (( $2 >= 1 )); then
        # first create w dirs
        for w in $(seq 1 $1);
        do {
            p="$3/w${w}d$2"
            COUNT=$(( COUNT + 1 ))
            echo "creating $COUNT/$MAX $OCIS_URL$ENDPOINT/$FOLDER$p"
            curl -X MKCOL "$OCIS_URL$ENDPOINT/$FOLDER$p" -u $USER:$PASSWORD -w "%{http_code}" $CURL_OPTS || { echo "could not create collection '$OCIS_URL$ENDPOINT/$FOLDER$p'" >&2; exit 1; } &
            create_tree $1 $(( $2 - 1 )) $p
        }
        done
    fi
}

#     creating 20/20 https://cloud.ocis.test/webdav/w20 x d1 folders/w20d1
#   creating 420/400 https://cloud.ocis.test/webdav/w20 x d2 folders/w20d2/w20d1
# creating 8420/8000 https://cloud.ocis.test/webdav/w20 x d3 folders/w20d3/w20d2/w20d1

#      creating 10/10 https://cloud.ocis.test/webdav/w10 x d1 folders/w10d1
#    creating 110/100 https://cloud.ocis.test/webdav/w10 x d2 folders/w10d2/w10d1
#  creating 1110/1000 https://cloud.ocis.test/webdav/w10 x d3 folders/w10d3/w10d2/w10d1
#creating 11110/10000 https://cloud.ocis.test/webdav/w10 x d4 folders/w10d4/w10d3/w10d2/w10d1  

# w^d + 

curl -X MKCOL "$OCIS_URL$ENDPOINT/$FOLDER" -u $USER:$PASSWORD -w "%{http_code}" $CURL_OPTS || { echo "could not create collection '$OCIS_URL$ENDPOINT/$FOLDER/'" >&2; exit 1; }

create_tree $WIDTH $DEPTH
