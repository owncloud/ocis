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


  Scenario Outline: the admin user can edit another user's email
    When the user "Alice" changes the email of user "Brian" to "<newEmail>" using the Graph API
    Then the HTTP status code should be "<code>"
    And the user "Brian" should have information with these key and value pairs:
      | key  | value           |
      | mail | <emailAsResult> |
    Examples:
      | action description        | newEmail             | code | emailAsResult        |
      | change to a valid email   | newemail@example.com | 200  | newemail@example.com |
      | override existing mail    | brian@example.com    | 200  | brian@example.com    |
      | two users with same mail  | alice@example.org    | 200  | alice@example.org    |
      | empty mail                |                      | 400  | brian@example.com    |
      | change to a invalid email | invalidEmail         | 400  | brian@example.com    |


  Scenario Outline: the admin user can edit another user's name
    Given user "Carol" has been created with default attributes and without skeleton files
    When the user "Alice" changes the user name of user "Carol" to "<userName>" using the Graph API
    Then the HTTP status code should be "<code>"
    And the user "<userNameAsResult>" should have information with these key and value pairs:
      | key                      | value              |
      | onPremisesSamAccountName | <userNameAsResult> |
    Examples:
      | action description           | userName | code | userNameAsResult |
      | change to a valid user name  | Lionel   | 200  | Lionel           |
      | user name characters         | *:!;_+-& | 200  | *:!;_+-&         |
      | change to existing user name | Brian    | 409  | Brian            |
      | empty user name              |          | 400  | Brian            |


  Scenario: the admin user changes the name of a user to the name of an existing disabled user
    Given the user "Alice" has created a new user using the Graph API with the following settings:
      | userName    | sam             |
      | displayName | sam             |
      | email       | sam@example.com |
      | password    | 1234            |
    And the user "Alice" has disabled user "Brian" using the Graph API
    When the user "Alice" changes the user name of user "sam" to "Brian" using the Graph API
    Then the HTTP status code should be "409"
    And the user "sam" should have information with these key and value pairs:
      | key                      | value |
      | onPremisesSamAccountName | sam   |


  Scenario: the admin user changes the name of a user to the name of a previously deleted user
    Given the user "Alice" has created a new user using the Graph API with the following settings:
      | userName    | sam             |
      | displayName | sam             |
      | email       | sam@example.com |
      | password    | 1234            |
    And the user "Alice" has deleted a user "sam" using the Graph API
    When the user "Alice" changes the user name of user "Brian" to "sam" using the Graph API
    Then the HTTP status code should be "200"
    And the user "sam" should have information with these key and value pairs:
      | key                      | value |
      | onPremisesSamAccountName | sam   |


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


  Scenario Outline: the admin user can edit another user display name
    When the user "Alice" changes the display name of user "Brian" to "<newDisplayName>" using the Graph API
    Then the HTTP status code should be "200"
    And the user "Brian" should have information with these key and value pairs:
      | key         | value                 |
      | displayName | <displayNameAsResult> |
    Examples:
      | action description                | newDisplayName | code | displayNameAsResult |
      | change to a display name          | Olaf Scholz    | 200  | Olaf Scholz         |
      | override to existing display name | Carol King     | 200  | Carol King          |
      | change to an empty display name   |                | 400  | Brian Murphy        |
      | displayName with characters       | *:!;_+-&#(?)   | 200  | *:!;_+-&#(?)        |


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

  @skipOnStable2.0
  Scenario: the admin user disables another user
    When the user "Alice" disables user "Brian" using the Graph API
    Then the HTTP status code should be "200"
    When user "Alice" gets information of user "Brian" using Graph API
    Then the HTTP status code should be "200"
    And the user retrieve API response should contain the following information:
      | displayName  | id        | mail              | onPremisesSamAccountName | accountEnabled |
      | Brian Murphy | %uuid_v4% | brian@example.com | Brian                    | false          |

  @skipOnStable2.0
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

  @skipOnStable2.0
  Scenario: the admin user enables disabled user
    Given the user "Alice" has disabled user "Brian" using the Graph API
    When the user "Alice" enables user "Brian" using the Graph API
    Then the HTTP status code should be "200"
    When user "Alice" gets information of user "Brian" using Graph API
    Then the HTTP status code should be "200"
    And the user retrieve API response should contain the following information:
      | displayName  | id        | mail              | onPremisesSamAccountName | accountEnabled |
      | Brian Murphy | %uuid_v4% | brian@example.com | Brian                    | true           |

  @skipOnStable2.0
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
