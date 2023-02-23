@api @skipOnOcV10
Feature: edit user

  Note - this feature is run in CI with ACCOUNTS_HASH_DIFFICULTY set to the default for production
  See https://github.com/owncloud/ocis/issues/1542 and https://github.com/owncloud/ocis/pull/839

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files
    And the administrator has given "Alice" the role "Admin" using the settings api
    And the user "Alice" has created a new user using the Graph API with the following settings:
      | userName    | Brian             |
      | displayName | Brian Murphy      |
      | email       | brian@example.com |
      | password    | 1234              |


  Scenario: the admin user can edit another user email
    When the user "Alice" changes the email of user "Brian" to "newemail@example.com" using the Graph API
    Then the HTTP status code should be "200"
    And the user "Brian" should have information with these key and value pairs:
      | key  | value                |
      | mail | newemail@example.com |


  Scenario: the admin user can override an existing user email of another user
    When the user "Alice" changes the email of user "Brian" to "brian@example.com" using the Graph API
    Then the HTTP status code should be "200"
    And the user "Brian" should have information with these key and value pairs:
      | key  | value             |
      | mail | brian@example.com |


  Scenario: the admin user cannot clear an existing user email
    When the user "Alice" tries to change the email of user "Brian" to "" using the Graph API
    Then the HTTP status code should be "400"
    And the user "Brian" should have information with these key and value pairs:
      | key  | value             |
      | mail | brian@example.com |


  Scenario Outline: a normal user should not be able to change their email address
    Given the administrator has given "Brian" the role "<role>" using the settings api
    When the user "Brian" tries to change the email of user "Brian" to "newemail@example.com" using the Graph API
    Then the HTTP status code should be "401"
    And the user "Brian" should have information with these key and value pairs:
      | key  | value             |
      | mail | brian@example.com |
    Examples:
      | role        |
      | Space Admin |
      | User        |


  Scenario Outline: a normal user should not be able to edit another user's email
    Given the administrator has given "Brian" the role "<role>" using the settings api
    And the user "Alice" has created a new user using the Graph API with the following settings:
      | userName    | Carol             |
      | displayName | Carol King        |
      | email       | carol@example.com |
      | password    | 1234              |
    When the user "Brian" tries to change the email of user "Carol" to "newemail@example.com" using the Graph API
    Then the HTTP status code should be "401"
    And the user "Carol" should have information with these key and value pairs:
      | key  | value             |
      | mail | carol@example.com |
    Examples:
      | role        |
      | Space Admin |
      | User        |


  Scenario: the admin user can edit another user display name
    When the user "Alice" changes the display name of user "Brian" to "Carol King" using the Graph API
    Then the HTTP status code should be "200"
    And the user "Brian" should have information with these key and value pairs:
      | key         | value      |
      | displayName | Carol King |


  Scenario: the admin user cannot clear another user display name
    When the user "Alice" tries to change the display name of user "Brian" to "" using the Graph API
    Then the HTTP status code should be "400"
    And the user "Brian" should have information with these key and value pairs:
      | key         | value        |
      | displayName | Brian Murphy |


  Scenario Outline: a normal user should not be able to change his/her own display name
    Given the administrator has given "Brian" the role "<role>" using the settings api
    When the user "Brian" tries to change the display name of user "Brian" to "Brian Murphy" using the Graph API
    Then the HTTP status code should be "401"
    And the user "Alice" should have information with these key and value pairs:
      | key         | value        |
      | displayName | Alice Hansen |
    Examples:
      | role        |
      | Space Admin |
      | User        |


  Scenario Outline: a normal user should not be able to edit another user's display name
    Given the administrator has given "Brian" the role "<role>" using the settings api
    And the user "Alice" has created a new user using the Graph API with the following settings:
      | userName    | Carol             |
      | displayName | Carol King        |
      | email       | carol@example.com |
      | password    | 1234              |
    When the user "Brian" tries to change the display name of user "Carol" to "Alice Hansen" using the Graph API
    Then the HTTP status code should be "401"
    And the user "Carol" should have information with these key and value pairs:
      | key         | value      |
      | displayName | Carol King |
    Examples:
      | role        |
      | Space Admin |
      | User        |


  Scenario: the admin user resets password of another user
    Given user "Brian" has uploaded file with content "test file for reset password" to "/resetpassword.txt"
    When the user "Alice" resets the password of user "Brian" to "newpassword" using the Graph API
    Then the HTTP status code should be "200"
    And the content of file "resetpassword.txt" for user "Brian" using password "newpassword" should be "test file for reset password"


  Scenario Outline: a normal user should not be able to reset the password of another user
    Given the administrator has given "Brian" the role "<role>" using the settings api
    And the user "Alice" has created a new user using the Graph API with the following settings:
      | userName    | Carol             |
      | displayName | Carol King        |
      | email       | carol@example.com |
      | password    | 1234              |
    And user "Carol" has uploaded file with content "test file for reset password" to "/resetpassword.txt"
    When the user "Brian" resets the password of user "Carol" to "newpassword" using the Graph API
    Then the HTTP status code should be "401"
    And the content of file "resetpassword.txt" for user "Carol" using password "1234" should be "test file for reset password"
    But user "Carol" using password "newpassword" should not be able to download file "resetpassword.txt"
    Examples:
      | role        |
      | Space Admin |
      | User        |


  Scenario: the admin user disables another user
    When the user "Alice" disables user "Brian" using the Graph API
    Then the HTTP status code should be "200"
    When user "Alice" gets information of user "Brian" using Graph API
    Then the HTTP status code should be "200"
    And the user retrieve API response should contain the following information:
      | displayName  | id        | mail              | onPremisesSamAccountName | accountEnabled |
      | Brian Murphy | %uuid_v4% | brian@example.com | Brian                    | false          |


  Scenario Outline: a normal user should not be able to disable another user
    Given user "Carol" has been created with default attributes and without skeleton files
    And the administrator has given "Brian" the role "<role>" using the settings api
    When the user "Brian" tries to disable user "Carol" using the Graph API
    Then the HTTP status code should be "401"
    When user "Alice" gets information of user "Carol" using Graph API
    Then the HTTP status code should be "200"
    And the user retrieve API response should contain the following information:
      | displayName | id        | mail              | onPremisesSamAccountName | accountEnabled |
      | Carol King  | %uuid_v4% | carol@example.org | Carol                    | true           |
    Examples:
      | role        |
      | Space Admin |
      | User        |
      | Guest       |


  Scenario: the admin user enables disabled user
    Given the user "Alice" has disabled user "Brian" using the Graph API
    When the user "Alice" enables user "Brian" using the Graph API
    Then the HTTP status code should be "200"
    When user "Alice" gets information of user "Brian" using Graph API
    Then the HTTP status code should be "200"
    And the user retrieve API response should contain the following information:
      | displayName  | id        | mail              | onPremisesSamAccountName | accountEnabled |
      | Brian Murphy | %uuid_v4% | brian@example.com | Brian                    | true           |


  Scenario Outline: a normal user should not be able to enable another user
    Given user "Carol" has been created with default attributes and without skeleton files
    And the user "Alice" has disabled user "Carol" using the Graph API
    And the administrator has given "Brian" the role "<role>" using the settings api
    When the user "Brian" tries to enable user "Carol" using the Graph API
    Then the HTTP status code should be "401"
    When user "Alice" gets information of user "Carol" using Graph API
    Then the HTTP status code should be "200"
    And the user retrieve API response should contain the following information:
      | displayName | id        | mail              | onPremisesSamAccountName | accountEnabled |
      | Carol King  | %uuid_v4% | carol@example.org | Carol                    | false          |
    Examples:
      | role        |
      | Space Admin |
      | User        |
      | Guest       |
