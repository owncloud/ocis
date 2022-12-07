Feature: Create user
    As a admin
    I want to create a user
    So that I can use the user in other tests

  Scenario: the admin creates a user
    When the administrator creates user using the Graph API with the following settings:
      # | key         | value             |
      # | userName    | Alice             |
      # | displayName | Alice Hansen      |
      # | email       | alice@example.com |
      # | password    | 123456            |
    Then the HTTP status code should be "200"
    And user "Alice" should exist
