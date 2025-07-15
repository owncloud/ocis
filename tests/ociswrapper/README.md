## oCIS Wrapper

The oCIS wrapper is a tool that wraps the oCIS binary and allows to dynamically reconfigure or extend the running instance. This is done by sending JSON payloads with updated environment variables.

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
     --bin string              Full oCIS binary path (default "/usr/bin/ocis")
     --url string              oCIS server url (default "https://localhost:9200")
     --retry string            Number of retries to start oCIS server (default "5")
     -p, --port string         Wrapper API server port (default "5200")
     --admin-username string   admin username for oCIS server
     --admin-password string   admin password for oCIS server
     --skip-ocis-run           Skip running oCIS server
    ```

Access the API server at `http://localhost:5200`.

Also, see `./bin/ociswrapper help` for more information.

### API

**ociswrapper** exposes the following endpoints:

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

    - `200 OK` - oCIS is successfully reconfigured
    - `400 Bad Request` - request body is not a valid JSON object
    - `500 Internal Server Error` - oCIS server is not running

2.  `DELETE /rollback`

    Rolls back the configuration to the starting point.

    Returns:

    - `200 OK` - rollback is successful
    - `500 Internal Server Error` - oCIS server is not running

3.  `POST /command`

    Executes the provided command on the oCIS server. The body of the request should be a JSON object with the following structure:

    ```yml
    {
      "command": "<ocis-command>", # without the ocis binary. e.g. "list"
    }
    ```

    If the command requires user input, the body of the request should be a JSON object with the following structure:

    ```json
    {
      "command": "<ocis-command>",
      "inputs": ["value1"]
    }
    ```

    Returns:

    ```json
    {
      "status": "OK",
      "exitCode": 0,
      "message": "<command output>"
    }
    OR
    {
      "status": "ERROR",
      "exitCode": <error-exit-code>,
      "message": "<command output>"
    }
    ```

    - `200 OK` - command is successfully executed
    - `400 Bad Request` - request body is not a valid JSON object
    - `500 Internal Server Error`

4.  `POST /start`

    Starts the oCIS server.

    Returns:

    - `200 OK` - oCIS server is started
    - `409 Conflict` - oCIS server is already running
    - `500 Internal Server Error` - Unable to start oCIS server

5.  `POST /stop`

    Stops the oCIS server.

    Returns:

    - `200 OK` - oCIS server is stopped
    - `500 Internal Server Error` - Unable to stop oCIS server

6. `POST /services/{service-name}`

   Restart oCIS instances without specified service and start that service independently (not covered by the oCIS supervisor).

    Body of the request should be a JSON object with the following structure:

    ```json
    {
      "ENV_KEY1": "value1",
      "ENV_KEY2": "value2"
    }
    ```

    > **⚠️ Note:**
    >
    > You need to set the proper addresses to access the service from other steps in the CI pipeline.
    >
    > `{SERVICE-NAME}_DEBUG_ADDR=0.0.0.0:{DEBUG_PORT}`
    >
    > `{SERVICE-NAME}_HTTP_ADDR=0.0.0.0:{HTTP_PORT}`

   Returns:

    - `200 OK` - oCIS service started successfully
    - `400 Bad Request` - request body is not a valid JSON object
    - `500 Internal Server Error` - Failed to start oCIS service audit

7. `DELETE /services/{service-name}`

   Stop individually running oCIS service

   Returns:

    - `200 OK` - oCIS service stopped successfully
    - `500 Internal Server Error` - Unable to stop oCIS service

8. `PATCH /services/{service-name}`

    Updates the configuration of the running service instance.
    Body of the request should be a JSON object with the following structure:

    ```json
    {
      "ENV_KEY1": "value1",
      "ENV_KEY2": "value2"
    }
    ```

    Returns:

    - `200 OK` - service is successfully reconfigured
    - `500 Internal Server Error` - service is not running

9. `DELETE /services/rollback`

    Stop and rollback all service configurations to the starting point.

    Returns:

    - `200 OK` - rollback is successful
    - `500 Internal Server Error` - oCIS server is not running
