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
    And user "Alice" has created folder "folder"
    And user "Alice" has shared folder "folder" with user "Brian" with permissions "31"
    And user "Brian" accepts share "/folder" offered by user "Alice" using the sharing API
    And as "Brian" folder "Shares/folder" should exist
    And user "Brian" has shared folder "Shares/folder" with user "Carol" with permissions "31"
    And user "Carol" accepts share "/folder" offered by user "Brian" using the sharing API
    And as "Carol" folder "Shares/folder" should exist
    And user "Carol" has shared folder "Shares/folder" with user "Damian" with permissions "17"
    And user "Damian" accepts share "/folder" offered by user "Carol" using the sharing API
    And as "Damian" folder "Shares/folder" should exist

  Scenario Outline: You should only be able to see direct incoming and outgoing shares not all the chain:
    When user "<user>" gets all the shares inside the folder "Shares/folder" using the sharing API
    Then the OCS status code should be "100"
    And the HTTP status code should be "200"
    #And user "Alice" <AliceVisible> included in the response $TODO: How to check alice is owner every time?
    And the response should contain <numVisibleShares> entries
    And user "Brian" <BrianVisible> included in the response
    And user "Carol" <CarolVisible> included in the response
    And user "Damian" <DamianVisible> included in the response
    Examples:
      | user   | numVisibleShares | BrianVisible  | CarolVisible   | DamianVisible  |
      | Brian  | 2                | should be     | should be      | should not be  |
      | Carol  | 2                | should not be | should be      | should be      |
      | Damian | 1                | should not be | should not be  | should be      |

  Scenario: Owners can see all the chain:
    When user "Alice" gets all the shares inside the folder "folder" using the sharing API
    Then the OCS status code should be "100"
    And the HTTP status code should be "200"
    And the response should contain 3 entries
    And user "Brian" should be included in the response
    And user "Carol" should be included in the response
    And user "Damian" should be included in the response

  Scenario: You can't share with more permissions than you have
    When user "Ember" has been created with default attributes and without skeleton files    
    And user "Damian" shares folder "Shares/folder" with user "Ember" with permissions "31" using the sharing API
    Then the OCS status code should be "404"
    And the OCS status message should be "Cannot set the requested share permissions"

  Scenario Outline: Editing reshares
    When user "Fred" has been created with default attributes and without skeleton files    
    And user "Carol" has shared folder "Shares/folder" with user "Fred" with permissions "17"
    And user "Fred" accepts share "/folder" offered by user "Carol" using the sharing API
    And as "Fred" folder "Shares/folder" should exist
    Then user "<user>" updates the last share using the sharing API with
      | permissions | 31 |
    And the OCS status code should be "<code>"
    And user "Fred" <canUpload> able to upload file "filesForUpload/textfile.txt" to "/Shares/folder/textfile.txt"
    Examples:
      | user  | code | canUpload     |
      | Alice | 100  | should be     |
      | Brian | 404  | should not be |
      | Carol | 100  | should be     |

  Scenario Outline: Deleting reshares
    When user "Gina" has been created with default attributes and without skeleton files    
    And user "Carol" has shared folder "Shares/folder" with user "Gina" with permissions "17"
    And user "Gina" accepts share "/folder" offered by user "Carol" using the sharing API
    And as "Gina" folder "Shares/folder" should exist
    Then user "<user>" deletes the last share using the sharing API
    And the OCS status code should be "<code>"
    And as "Gina" folder "Shares/folder" <exists>
    Examples:
      | user  | code | exists           |
      | Alice | 100  | should not exist |
      | Brian | 400  | should exist     |
      | Carol | 100  | should not exist |

