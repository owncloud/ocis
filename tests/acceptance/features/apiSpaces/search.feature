@api @skipOnOcV10
Feature: Search
  It is possible to search files in the shares jail and the project space

  Note - this feature is run in CI with ACCOUNTS_HASH_DIFFICULTY set to the default for production
  See https://github.com/owncloud/ocis/issues/1542 and https://github.com/owncloud/ocis/pull/839

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
    And using new DAV path
    And the administrator has given "Alice" the role "Space Admin" using the settings api
    And user "Alice" has created a space "find data" with the default quota using the GraphApi
    And user "Alice" has created a folder "folder/SubFolder1/subFOLDER2" in space "find data"
    And user "Alice" has uploaded a file inside space "find data" with content "some content" to "folder/SubFolder1/subFOLDER2/insideTheFolder.txt"


  Scenario: Alice can find data from the project space
    When user "Alice" searches for "fol" using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result should contain "4" entries
    And the search result of user "Alice" should contain these entries:
      | /folder                                           |
      | /folder/SubFolder1                                |
      | /folder/SubFolder1/subFOLDER2                     |
      | /folder/SubFolder1/subFOLDER2/insideTheFolder.txt |


  Scenario: Alice can find data from the project space
    When user "Alice" searches for "SUB" using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result should contain "2" entries
    And the search result of user "Alice" should contain these entries:
      | /folder/SubFolder1             |
      | /folder/SubFolder1/subFOLDER2  |
    But the search result of user "Alice" should not contain these entries:
      | /folder                                           |
      | /folder/SubFolder1/subFOLDER2/insideTheFolder.txt |


  Scenario: Brian can find data from the shares jail
    Given user "Alice" shares the following entity "folder" inside of space "find data" with user "Brian" with role "viewer"
    And user "Brian" has accepted share "/folder" offered by user "Alice"
    When user "Brian" searches for "folder" using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result should contain "4" entries
    And the search result of user "Brian" should contain these entries:
      | /folder                                           |
      | /folder/SubFolder1                                |
      | /folder/SubFolder1/subFOLDER2                     |
      | /folder/SubFolder1/subFOLDER2/insideTheFolder.txt |

