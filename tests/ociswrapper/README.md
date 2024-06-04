## oCIS Wrapper

A tool that wraps the oCIS binary and provides a way to re-configure the running oCIS instance.

When run, **ociswrapper** starts an API server that exposes some endpoints to re-configure the oCIS server.

### Usage

1.  Build

    ```bash
    make build
    ```

2.  Run

    ```bash
    ./bin/ociswrapper serve --bin=<path-to-ocis-binary>
    ```
    To check other available options:
    ```bash
    ./bin/ociswrapper serve --help
    ```
    ```bash
     --url string              oCIS server url (default "https://localhost:9200")
     --retry string            Number of retries to start oCIS server (default "5")
     -p, --port string         Wrapper API server port (default "5200")
     --admin-username string   admin username for oCIS server
     --admin-password string   admin password for oCIS server
    ```

Access the API server at `http://localhost:5200`.

Also, see `./bin/ociswrapper help` for more information.

### API

**ociswrapper** exposes two endpoints:

1.  `PUT /config`

    Updates the configuration of the running oCIS instance.
    Body of the request should be a JSON object with the following structure:

    ```json
    {
      "ENV_KEY1": "value1",
      "ENV_KEY2": "value2"
    }
    ```

    Returns:

    *   `200 OK` - oCIS is successfully reconfigured
    *   `400 Bad Request` - request body is not a valid JSON object
    *   `500 Internal Server Error` - oCIS server is not running

2.  `DELETE /rollback`

    Rolls back the configuration to the starting point.

    Returns:

    *   `200 OK` - rollback is successful
    *   `500 Internal Server Error` - oCIS server is not running
