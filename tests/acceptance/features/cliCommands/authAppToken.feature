@env-config
Feature: create auth-app token
  As an admin
  I want to create auth-app Tokens
  So that I can use 3rd party apps


  Scenario: creates auth-app token via CLI
    Given the following configs have been set:
      | config                | value    |
      | OCIS_ADD_RUN_SERVICES | auth-app |
      | PROXY_ENABLE_APP_AUTH | true     |
    And user "Alice" has been created with default attributes
    When the administrator creates app token for user "Alice" with expiration time "72h" using the auth-app CLI
    Then the command should be successful
    And the command output should contain "App token created for Alice"

  @env-config
  Scenario: user deletes the created auth-app token
    Given the config "AUTH_APP_ENABLE_IMPERSONATION" has been set to "true"
    And the administrator creates app token for user "Alice" with expiration time "72h" using the auth-app CLI
    When user "Admin" deletes all the created auth-app tokens using the auth-app API
    Then the HTTP status code should be "200"
    And user "Admin" should have "0" auth-app tokens
    When user "Admin" lists all created tokens using the auth-app API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "array",
        "minItems": 0,
        "maxItems": 0
      }
      """
