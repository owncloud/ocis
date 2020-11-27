## Requirements
*  [K6](https://k6.io/)
*  [YARN](https://yarnpkg.com/)
*  [OCIS](https://github.com/owncloud/ocis)

## How to build
```console
$ yarn
$ yarn build
```

## How to run
```console
k6 run ./dist/test/NAME_OF_TEST.js
```

## Environment variables
```console
$ OC_OCIS_HOST=URL k6 run ...
$ OC_OIDC_HOST=URL k6 run ...
$ OC_OIDC=BOOL k6 run ...
$ OC_TEST_FILE=STRING k6 run ...
```