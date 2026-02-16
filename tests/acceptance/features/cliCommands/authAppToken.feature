@env-config
Feature: create auth-app token
  As an administrator
  I want to create auth-app Tokens
  So that I can use 3rd party apps

  Background:
    Given user "Alice" has been created with default attributes
    And the following configs have been set:
      | service | config                | value    |
      | authapp | OCIS_ADD_RUN_SERVICES | auth-app |
      | authapp | PROXY_ENABLE_APP_AUTH | true     |


  Scenario: creates auth-app token via CLI
    When the administrator creates auth-app token for user "Alice" with expiration time "72h" using the auth-app CLI
    Then the command should be successful
    And the command output should contain "App token created for Alice"


  Scenario: user deletes the created auth-app token
    Given the config "AUTH_APP_ENABLE_IMPERSONATION" has been set to "true" for "auth-app" service
    And the administrator has created app token for user "Alice" with expiration time "72h" using the auth-app CLI
    When user "Alice" deletes all the created auth-app tokens using the auth-app API
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


  Scenario: try to create auth-app token without expiry
    When the administrator tries to create auth-app token for user "Alice" with expiration time "" using the auth-app CLI
    Then the command should be unsuccessful
    And the command output should contain "time: invalid duration"
