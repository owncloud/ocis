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
***Thest tests will have absolute control over your server. Running them in your production server might result in the permanent loss of data***

Use the following command to run the tests
```console
k6 run ./dist/TESTNAME.js
```

## Running with different backends
### 1. Running with OCIS backend
The tests run by default on the oCIS backend. They use the address https://localhost:9200 to run the tests.
If your oCIS instance is running on different address use `OC_HOST_NAME` env variable to specify the address of the server.

### 2. Running with OC10 (classic) backend
To run the tests with oc10 classic backend set the address of oc10 server on `OC_HOST_NAME` env variable and also set `TEST_OC10` to `true`

eg.
```
OC_HOST_NAME=http://owncloud-server.com TEST_OC10=true k6 run ./dist/TESTNAME.js
```