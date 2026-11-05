@env-config
Feature: reset user password via CLI command


  Scenario: reset user password
    Given the user "Admin" has created a new user with the following attributes:
      | userName    | Alice        |
      | displayName | Alice Hansen |
      | password    | %alt1%       |
    And the administrator has stopped the server
    When the administrator resets the password of existing user "Alice" to "newpass" using the CLI
    Then the command should be successful
    And the command output should contain "Password for user 'uid=Alice,ou=users,o=libregraph-idm' updated."
    But the command output should not contain "Failed to update user password: entry does not exist"
    And the administrator has started the server
    And user "Alice" should be able to create folder "newFolder" using password "newpass"
    But user "Alice" should not be able to create folder "anotherFolder" using password "%alt1%"


  Scenario: try to reset password of non-existing user
    Given the administrator has stopped the server
    When the administrator resets the password of non-existing user "Alice" to "newpass" using the CLI
    Then the command should be successful
    But the command output should contain "Failed to update user password: entry does not exist"


  Scenario: reset password of admin user
    Given the user "Admin" has created a new user with the following attributes:
      | userName    | Alice        |
      | displayName | Alice Hansen |
      | password    | %alt1%       |
    And the administrator has assigned the role "Admin" to user "Alice" using the Graph API
    And the administrator has stopped the server
    When the administrator resets the password of existing user "Alice" to "newpass" using the CLI
    Then the command should be successful
    And the command output should contain "Password for user 'uid=Alice,ou=users,o=libregraph-idm' updated."
    But the command output should not contain "Failed to update user password: entry does not exist"
    And the administrator starts the server
    And user "Alice" using password "newpass" should be able to create a new user "Brian" with default attributes


  Scenario: reset password after renaming the admin user
    Given the user "Admin" has created a new user with the following attributes:
      | userName    | Alice        |
      | displayName | Alice Hansen |
      | password    | %alt1%       |
    And the administrator has assigned the role "Admin" to user "Alice" using the Graph API
    And user "Alice" has changed the username to "superUser"
    And the administrator has stopped the server
    When the administrator resets the password of existing user "superUser" to "newpass" using the CLI
    Then the command should be successful
    And the command output should contain "Password for user 'uid=superUser,ou=users,o=libregraph-idm' updated."
    But the command output should not contain "Failed to update user password: entry does not exist"
    And the administrator starts the server
    And user "superUser" using password "newpass" should be able to create a new user "Brian" with default attributes


  Scenario: try to reset password of regular user using service user type
    Given the user "Admin" has created a new user with the following attributes:
      | userName    | Alice        |
      | displayName | Alice Hansen |
      | password    | %alt1%       |
    And the administrator has stopped the server
    When the administrator resets the password of existing user "Alice" to "newpass" with user type "service" using the CLI
    Then the command should be successful
    But the command output should contain "Failed to update user password: entry does not exist"


  Scenario: reset password of service user
    Given the administrator has stopped the server
    When the administrator resets the password of existing user "reva" to "newpass" with user type "service" using the CLI
    Then the command should be successful
    And the command output should contain "Password for user 'uid=reva,ou=sysusers,o=libregraph-idm' updated."
    And the command output should not contain "Failed to update user password"
    # resetpassword changes the IDM entry; the basic auth check below catches a
    # reva bind password that was not synced into ocis.yaml.
    And the administrator has started the server
    And the admin user should be able to access the Graph API
    # resetpassword needs the IDM database offline, so stop again to restore the original password.
    And the administrator has stopped the server
    When the administrator rollsback the password of service user "reva" using the CLI
    Then the command should be successful
    And the command output should not contain "Failed to update user password"
