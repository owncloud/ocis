
# oCIS EOS example

`docker-compose up -d`


## setup a space
https://github.com/owncloud-docker/eos-stack/blob/master/eos-mgm/setup


`docker-compose exec mgm-master bash`

```
eos -r 0 0 space define default
eos -r 0 0 space set default on
eos -r 0 0 space config default space.policy.recycle=on
eos -r 0 0 recycle config --add-bin /eos/dockertest/reva/users
eos -r 0 0 recycle config --size 1G

eos -r 0 0  space ls
```


## setup fst

https://github.com/owncloud-docker/eos-stack/blob/master/eos-fst/setup

`docker-compose exec fst bash`


```
for i in {1..4}; do
  mkdir -p /disks/eosfs${i}
  chown daemon:daemon /disks/eosfs${i}
  eos -r 0 0 -b fs add eosfs${i} fst.testnet:1095 /disks/eosfs${i} default rw
done

eos fs ls
```
