@api
Feature:
  As an administrator
  I want to be able to create user using the Graph API
  So that I can manage users more easily


  @smokeTest
  Scenario: admin creates a user
    Given user "brand-new-user" has been deleted
    When the administrator sends a user creation request for user "brand-new-user" password "%alt1%" using the graph API
    Then the HTTP status code should be "200"
    And user "brand-new-user" should exist
    And user "brand-new-user" should be able to upload file "filesForUpload/textfile.txt" to "/textfile.txt"


  Scenario Outline: admin creates a user with special characters in the username
    Given user "<username>" has been deleted
    When the administrator sends a user creation request for user "<username>" password "%alt1%" using the graph API
    Then the HTTP status code of responses on all endpoints should be "400"
    And the graph API response should return the following error
      | code    | invalidRequest                                                    |
      | message | username '<username>' must be at least the local part of an email |
    And user "<username>" should not exist
    Examples:
      | username |
      | a@-+_.b  |
      | a space  |

  Scenario: admin creates a user and specifies a password with special characters
    When the administrator sends a user creation request for the following users with password using the graph API
      | username        | password                     |
      | brand-new-user1 | !@#$%^&*()-_+=[]{}:;,.<>?~   |
      | brand-new-user2 | España§àôœ€                |
      | brand-new-user3 | नेपाली                         |
    And the HTTP status code of responses on all endpoints should be "200"
    And the following users should exist
      | username        |
      | brand-new-user1 |
      | brand-new-user2 |
      | brand-new-user3 |
    And the following users should be able to upload file "filesForUpload/textfile.txt" to "/textfile.txt"
      | username        |
      | brand-new-user1 |
      | brand-new-user2 |
      | brand-new-user3 |

  Scenario: admin tries to create an existing user
    And user "brand-new-user" has been created with default attributes and without skeleton files
    When the administrator sends a user creation request for user "brand-new-user" password "%alt1%" using the graph API
    And the HTTP status code should be "500"
    Then the graph API response should return the following error
      | code    | generalException                                   |
      | message | LDAP Result Code 68 "Entry Already Exists":{space} |


  Scenario: admin creates a user and specifies password containing just space
    Given user "brand-new-user" has been deleted
    When the administrator sends a user creation request for user "brand-new-user" password " " using the graph API
    And the HTTP status code should be "200"
    And user "brand-new-user" should exist
    And user "brand-new-user" should be able to upload file "filesForUpload/textfile.txt" to "/textfile.txt"

