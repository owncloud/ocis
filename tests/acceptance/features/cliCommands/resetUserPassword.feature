@env-config
Feature: reset user password via CLI command


  Scenario: reset user password
    Given the user "Admin" has created a new user with the following attributes:
      | userName    | Alice        |
      | displayName | Alice Hansen |
      | password    | %alt1%       |
    And the administrator has stopped the server
    When the administrator resets the password of user "Alice" to "newpass" using the CLI
    Then the command should be successful
    And the command output should contain "Password for user 'uid=Alice,ou=users,o=libregraph-idm' updated."
    But the command output should not contain "Failed to update user password: entry does not exist"
    And the administrator has started the server
    And user "Alice" should be able to create folder "newFolder" using password "newpass"
    But user "Alice" should not be able to create folder "anotherFolder" using password "%alt1%"
