@api @skipOnOcV10
Feature: Resharing
  It is possible to reshare files

  Note - this feature is run in CI with ACCOUNTS_HASH_DIFFICULTY set to the default for production
  See https://github.com/owncloud/ocis/issues/1542 and https://github.com/owncloud/ocis/pull/839

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
      | Carol    |
      | Damian   |
      | Ember    |
      | Fred     |
      | Gina     |
    And user "Alice" has created folder "folder"
    And user "Alice" has shared folder "folder" with user "Brian" with permissions "31"
    And user "Brian" has accepted share "/folder" offered by user "Alice"
    And user "Brian" has shared folder "Shares/folder" with user "Carol" with permissions "31"
    And user "Carol" has accepted share "/folder" offered by user "Brian"
    And user "Carol" has shared folder "Shares/folder" with user "Damian" with permissions "17"
    And user "Damian" has accepted share "/folder" offered by user "Carol"


  Scenario Outline: You should only be able to see direct outgoing shares not all the chain:
    Given user "Brian" has shared folder "Shares/folder" with user "Fred" with permissions "17"
    And user "Fred" has accepted share "/folder" offered by user "Brian"
    When user "<user>" gets all the shares inside the folder "Shares/folder" using the sharing API
    Then the OCS status code should be "100"
    And the HTTP status code should be "200"
    And the response should contain <numVisibleShares> entries
    And user "Brian" should not be included in the response
    And user "Carol" <CarolVisible> included in the response
    And user "Damian" <DamianVisible> included in the response
    And user "Fred" <FredVisible> included in the response
    Examples:
      | user   | numVisibleShares | CarolVisible  | DamianVisible | FredVisible   |
      | Brian  | 2                | should be     | should not be | should be     |
      | Carol  | 1                | should not be | should be     | should not be |
      | Damian | 0                | should not be | should not be | should not be |
      | Fred   | 0                | should not be | should not be | should not be |


  Scenario: Owners can see all the chain:
    When user "Alice" gets all the shares inside the folder "folder" using the sharing API
    Then the OCS status code should be "100"
    And the HTTP status code should be "200"
    And the response should contain 3 entries
    And user "Brian" should be included in the response
    And user "Carol" should be included in the response
    And user "Damian" should be included in the response


  Scenario: You can't share with more permissions than you have
    When user "Damian" shares folder "Shares/folder" with user "Ember" with permissions "31" using the sharing API
    Then the OCS status code should be "404"
    And the OCS status message should be "Cannot set the requested share permissions"


  Scenario Outline: Editing reshares
    Given user "Carol" has shared folder "Shares/folder" with user "Fred" with permissions "17"
    And user "Fred" has accepted share "/folder" offered by user "Carol"
    When user "<user>" updates the last share using the sharing API with
      | permissions | 31 |
    Then the OCS status code should be "<code>"
    And user "Fred" <canUpload> able to upload file "filesForUpload/textfile.txt" to "/Shares/folder/textfile.txt"
    Examples:
      | user  | code | canUpload     |
      | Alice | 100  | should be     |
      | Brian | 998  | should not be |
      | Carol | 100  | should be     |


  Scenario Outline: Deleting reshares
    Given user "Carol" has shared folder "Shares/folder" with user "Gina" with permissions "17"
    And user "Gina" has accepted share "/folder" offered by user "Carol"
    When user "<user>" deletes the last share using the sharing API
    Then the OCS status code should be "<code>"
    And as "Gina" folder "Shares/folder" <exists>
    And as "Carol" folder "Shares/folder" should exist
    Examples:
      | user  | code | exists           |
      | Alice | 100  | should not exist |
      | Brian | 400  | should exist     |
      | Carol | 100  | should not exist |


  Scenario Outline: Resharing with different permissions
    When user "<user>" shares folder "Shares/folder" with user "Ember" with permissions "<permissions>" using the sharing API
    Then the OCS status code should be "<code>"
    Examples:
      | user   | permissions | code |
      | Brian  | 17          | 100  |
      | Carol  | 31          | 100  |
      | Damian | 17          | 100  |
      | Damian | 27          | 404  |
      | Damian | 31          | 404  |


  Scenario Outline: Resharing files with different permissions
    Given user "Alice" has uploaded file with content "Random data" to "/file.txt"
    And user "Alice" has shared file "/file.txt" with user "Brian" with permissions "<shareepermissions>"
    And user "Brian" has accepted share "/file.txt" offered by user "Alice"
    When user "Brian" shares file "Shares/file.txt" with user "Fred" with permissions "<granteepermissions>" using the sharing API
    Then the OCS status code should be "<code>"
    Examples:
      | shareepermissions | granteepermissions | code |
      | 17                | 17                 | 100  |
      | 17                | 19                 | 404  |
      | 19                | 19                 | 100  |


  Scenario Outline: Resharing with group with different permissions
    Given group "security department" has been created
    And the administrator has added a user "Ember" to the group "security department" using GraphApi
    And the administrator has added a user "Fred" to the group "security department" using GraphApi
    When user "Brian" shares folder "Shares/folder" with group "security department" with permissions "<permissions>" using the sharing API
    Then the OCS status code should be "100"
    When user "Ember" accepts share "/folder" offered by user "Brian" using the sharing API
    Then user "Ember" <canUpload> able to upload file "filesForUpload/textfile.txt" to "/Shares/folder/textfile.txt"
    When user "Fred" accepts share "/folder" offered by user "Brian" using the sharing API
    Then user "Fred" <canUpload> able to upload file "filesForUpload/textfile.txt" to "/Shares/folder/textfile.txt"
    Examples:
      | permissions | canUpload     |
      | 17          | should not be |
      | 31          | should be     |
