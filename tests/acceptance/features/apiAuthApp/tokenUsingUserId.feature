Feature: create auth-app token using user-id
  As a user
  I want to create auth-app Tokens using user-id
  So that I can use 3rd party apps

  Background:
    Given user "Alice" has been created with default attributes

  @env-config @issue-11063
  Scenario: admin creates auth-app token for another user using user-id
    Given the config "AUTH_APP_ENABLE_IMPERSONATION" has been set to "true" for "authapp" service
    When user "Admin" creates app token with user-id for user "Alice" with expiration time "72h" using the auth-app API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["token","expiration_date","created_date","label"],
        "properties": {
          "token": { "pattern": "^[a-zA-Z0-9]{16}$" },
          "label": { "const": "Generated via Impersonation API" }
        }
      }
      """
    When user "Alice" lists all created tokens using the auth-app API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "array",
        "minItems": 1,
        "maxItems": 1,
        "items": {
          "oneOf": [
            {
              "type": "object",
              "required": [
                "token",
                "expiration_date",
                "created_date",
                "label"
              ],
              "properties": {
                "token": {
                  "pattern": "^\\$2a\\$11\\$[A-Za-z0-9./]{53}$"
                },
                "label": {
                  "const": "Generated via Impersonation API"
                }
              }
            }
          ]
        }
      }
      """


  Scenario: non-admin user tries to create own auth-app token using user-id with impersonation enabled
    Given the config "AUTH_APP_ENABLE_IMPERSONATION" has been set to "true" for "authapp" service
    When user "Alice" tries to create app token with user-id for user "Alice" with expiration time "72h" using the auth-app API
    Then the HTTP status code should be "403"

  @env-config @issue-11063
  Scenario: non-admin user tries to creates auth-app token for another user using user-id
    Given the config "AUTH_APP_ENABLE_IMPERSONATION" has been set to "true" for "authapp" service
    And user "Brian" has been created with default attributes
    When user "Brian" tries to create app token with user-id for user "Alice" with expiration time "72h" using the auth-app API
    Then the HTTP status code should be "403"

  @issue-11063
  Scenario: non-admin user tries to creates auth-app token for another user using user-id without impersonation enabled
    Given user "Brian" has been created with default attributes
    When user "Brian" tries to create app token with user-id for user "Alice" with expiration time "72h" using the auth-app API
    Then the HTTP status code should be "403"
    And the content in the response should include the following content:
      """
      impersonation is not allowed
      """


  Scenario: non-admin user tries to create own auth-app token using user-id and without expiry
    When user "Alice" tries to create app token with user-id for user "Alice" with expiration time "" using the auth-app API
    Then the HTTP status code should be "400"
    And the content in the response should include the following content:
      """
      error parsing expiry. Use e.g. 30m or 72h
      """

  @env-config
  Scenario: admin tries to create auth-app token for another user with user-id and without expiry
    Given the config "AUTH_APP_ENABLE_IMPERSONATION" has been set to "true" for "authapp" service
    When user "Admin" tries to create app token with user-id for user "Alice" with expiration time "" using the auth-app API
    Then the HTTP status code should be "400"
    And the content in the response should include the following content:
      """
      error parsing expiry. Use e.g. 30m or 72h
      """


  Scenario: non-admin user tries to create auth-app token for another user using user-id and without expiry
    Given the config "AUTH_APP_ENABLE_IMPERSONATION" has been set to "true" for "authapp" service
    And user "Brian" has been created with default attributes
    When user "Brian" tries to create app token with user-id for user "Alice" with expiration time "" using the auth-app API
    Then the HTTP status code should be "400"
    And the content in the response should include the following content:
      """
      error parsing expiry. Use e.g. 30m or 72h
      """


  Scenario: non-existent user tries to create auth-app token for another user using user-id and without expiry
    When user "Brian" tries to create app token with user-id for user "Alice" with expiration time "" using the auth-app API
    Then the HTTP status code should be "401"


  Scenario: non-existent user tries to create auth-app token for another user using user-id
    When user "Brian" tries to create app token with user-id for user "Alice" with expiration time "72h" using the auth-app API
    Then the HTTP status code should be "401"
