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