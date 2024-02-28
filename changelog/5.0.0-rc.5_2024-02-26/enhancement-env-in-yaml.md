Enhancement: Use environment variables in yaml config files

We added the ability to use environment variables in yaml config files. This allows to use environment variables in the config files of the ocis services which will be replaced by the actual value of the environment variable at runtime.

Example:

```
web:
  http:
    addr: ${SOME_HTTP_ADDR}
```

This makes it possible to use the same config file for different environments without the need to change the config file itself. This is especially useful when using docker-compose to run the ocis services. It is a common pattern to create an .env file which contains the environment variables for the docker-compose file. Now you can use the same .env file to configure the ocis services.

https://github.com/owncloud/ocis/pull/8339
