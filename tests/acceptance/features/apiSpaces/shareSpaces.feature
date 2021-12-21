@api @skipOnOcV10
Feature: Share spaces
  As the owner of a space
  I want to be able to add members to a space, and to remove access for them

  Note - this feature is run in CI with ACCOUNTS_HASH_DIFFICULTY set to the default for production
  See https://github.com/owncloud/ocis/issues/1542 and https://github.com/owncloud/ocis/pull/839

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files
    And user "Brian" has been created with default attributes and without skeleton files
    And the administrator has given "Alice" the role "Admin" using the settings api


  Scenario: A user can share a space to another user
    Given user "Alice" has created a space "Space to share" of type "project" with quota "10"
    When user "Alice" shares a space "Space to share" to user "Brian"
    Then the HTTP status code should be "200"
    And the OCS status code should be "200"
    And the OCS status message should be "OK"


  Scenario: A user can see that a received shared space is available
    Given user "Alice" has created a space "Share space to Brian" of type "project" with quota "10"
    And user "Alice" has shared a space "Share space to Brian" to user "Brian"
    When user "Brian" lists all available spaces via the GraphApi
    Then the json responded should contain a space "Share space to Brian" with these key and value pairs:
      | key              | value                            |
      | driveType        | share                            |
      | id               | %space_id%                       |
      | name             | Share space to Brian             |
      | quota@@@state    | normal                           |
      | root@@@webDavUrl | %base_url%/dav/spaces/%space_id% |


  Scenario: A user can see a file in a received shared space
    Given user "Alice" has created a space "Share space with file" of type "project" with quota "10"
    And user "Alice" has uploaded a file inside space "Share space with file" with content "Test" to "test.txt"
    When user "Alice" has shared a space "Share space with file" to user "Brian"
    Then for user "Brian" the space "Share space with file" should contain these entries:
      | test.txt |


  Scenario: A user can see a folder in received shared space
    Given user "Alice" has created a space "Share space with folder" of type "project" with quota "10"
    And user "Alice" has created a folder "Folder Main" in space "Share space with folder"
    When user "Alice" has shared a space "Share space with folder" to user "Brian"
    Then for user "Brian" the space "Share space with folder" should contain these entries:
      | Folder Main |


  Scenario: When a user unshares a space, the space becomes unavailable to the receiver
    Given user "Alice" has created a space "Unshare space" of type "project" with quota "10"
    And user "Alice" has shared a space "Unshare space" to user "Brian"
    When user "Brian" lists all available spaces via the GraphApi
    Then the json responded should contain a space "Unshare space" with these key and value pairs:
      | key       | value         |
      | driveType | share         |
      | id        | %space_id%    |
      | name      | Unshare space |
    When user "Alice" unshares a space "Unshare space" to user "Brian"
    Then the HTTP status code should be "200"
    And user "Brian" lists all available spaces via the GraphApi
    And the json responded should not contain a space "Unshare space"
