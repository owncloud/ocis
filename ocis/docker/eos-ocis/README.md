# Docker image for ocis running on eos

Image is based on [owncloud/eos-base](https://hub.docker.com/r/owncloud/eos-base) from [eos-stack](https://github.com/owncloud-docker/eos-stack)

## Build
Build owncloud/ocis master branch
```
docker build -t owncloud/eos-ocis:latest .
```

Or build a certain branch / tag
```
docker build -t owncloud/eos-ocis:1.0.0 --build-arg BRANCH=v1.0.0./eos-ocis
```

## Publish
```
docker push owncloud/eos-ocis:latest
```

## Maintainer 

- [Felix Böhm](https://github.com/felixboehm)

## Disclaimer 
Use only for development or testing. Setup is not secured nor tested.

## Example Usage

See https://github.com/owncloud-docker/compose-playground/tree/master/examples/eos-compose


