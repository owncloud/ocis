# Auth-App

The auth-app service provides authentication for 3rd party apps.

## The `auth` Service Family

ocis uses serveral authentication services for different use cases. All services that start with `auth-` are part of the authentication service family. Each member authenticates requests with different scopes. As of now, these services exist:
  -   `auth-app` handles authentication of external 3rd party apps
  -   `auth-basic` handles basic authentication
  -   `auth-bearer` handles oidc authentication
  -   `auth-machine` handles interservice authentication when a user is impersonated
  -   `auth-service` handles interservice authentication when using service accounts

## Service Startup

Because this service is not started automatically, a manual start needs to be initiated which can be done in several ways. To configure the service usage, an environment variable for the proxy service needs to be set to allow app authentication.
```bash
OCIS_ADD_RUN_SERVICES=auth-app  # deployment specific. Add the service to the manual startup list, use with binary deployments. Alternatively you can start the service explicitly via the command line.
PROXY_ENABLE_APP_AUTH=true      # mandatory, allow app authentication. In case of a distributed environment, this envvar needs to be set in the proxy service.
```

## App Tokens

App Tokens are used to authenticate 3rd party access via https like when using curl (apps) to access an API endpoint. These apps need to authenticate themselves as no logged in user authenticates the request. To be able to use an app token, one must first create a token. There are different options of creating a token.

### Via CLI (dev only)

Replace the `user-name` with an existing user. For the `token-expiration`, you can use any time abbreviation from the following list: `h, m, s`. Examples: `72h` or `1h` or `1m` or `1s.` Default is `72h`.

```bash
ocis auth-app create --user-name={user-name} --expiration={token-expiration}
```

Once generated, these tokens can be used to authenticate requests to ocis. They are passed as part of the request as `Basic Auth` header.

### Via API

The `auth-app` service provides an API to create (POST), list (GET) and delete (DELETE) tokens at the `/auth-app/tokens` endpoint.

When using curl for the respective command, you need to authenticate with a header. To do so, get from the browsers developer console the currently active bearer token. Consider that this token has a short lifetime. In any example, replace `<your host[:port]>` with the URL:port of your Infinite Scale instance, and `{token}`  `{value}` accordingly. Note that the active bearer token authenticates the user the token was issued for.

* **Create a token**\
  The POST request requires:
  * A `expiry` key/value pair in the form of `expiry=<number><h|m|s>`\
    Example: `expiry=72h`
  * An active bearer token
  ```bash
  curl --request POST 'https://<your host:9200>/auth-app/tokens?expiry={value}' \
       --header 'accept: application/json' \
       --header 'authorization: Bearer {token}'
  ```
  Example output:
  ```
  {
  "token": "3s2K7816M4vuSpd5",
  "expiration_date": "2024-08-08T13:42:42.796888022+02:00",
  "created_date": "2024-08-07T13:42:42+02:00",
  "label": "Generated via API"
  }
  ```

* **List tokens**\
  The GET request only requires an active bearer token for authentication:\
  Note that `--request GET` is technically not required because it is curl default. 
  ```bash
  curl --request GET 'https://<your host:9200>/auth-app/tokens' \
       --header 'accept: application/json' \
       --header 'authorization: Bearer {token}'
  ```
  Example output:
  ```
  [
    {
      "token": "$2a$11$EyudDGAJ18bBf5NG6PL9Ru9gygZAu0oPyLawdieNjGozcbXyyuUhG",
      "expiration_date": "2024-08-08T13:44:31.025199075+02:00",
      "created_date": "2024-08-07T13:44:31+02:00",
      "label": "Generated via Impersonation API"
    },
    {
      "token": "$2a$11$dfRBQrxRMPg8fvyvkFwaX.IPoIUiokvhzK.YNI/pCafk0us3MyPzy",
      "expiration_date": "2024-08-08T13:46:41.936052281+02:00",
      "created_date": "2024-08-07T13:46:42+02:00",
      "label": "Generated via Impersonation API"
    }
  ]
  ```

* **Delete a token**\
  The DELETE request requires:
  * A `token` key/value pair in the form of `token=<token_issued>`\
    Example: `token=Z3s2K7816M4vuSpd5`
  * An active bearer token
  ```bash
  curl --request DELETE 'https://<your host:9200>/auth-app/tokens?token={value}' \
       --header 'accept: application/json' \
       --header 'authorization: Bearer {token}'
  ```

### Via Impersonation API

When setting the environment variable `AUTH_APP_ENABLE_IMPERSONATION` to `true`, admins will be able to use the `/auth-app/tokens` endpoint to create tokens for other users but using their own bearer token for authentication. This can be important for migration scenarios, but should not be considered for regular tasks on a production system for security reasons.

To impersonate, the respective requests from the CLI commands above extend with the following parameters, where you can use one or the other:

* The `userID` in the form of: `userID={value}`\
  Example:\
  `userID=4c510ada- ... -42cdf82c3d51`

* The `userName` in the form of: `userName={value}`\
  Example:\
  `userName=einstein`

Example:\
A final create request would then look like:
```bash
curl --request POST 'https://<your host:9200>/auth-app/tokens?expiry={value}&userName={value}' \
     --header 'accept: application/json' \
     --header 'authorization: Bearer {token}'
```
