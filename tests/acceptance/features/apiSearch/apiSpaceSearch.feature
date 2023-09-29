Feature: Search
  As a user
  I want to search for resources in the space
  So that I can get them quickly

  Note - this feature is run in CI with ACCOUNTS_HASH_DIFFICULTY set to the default for production
  See https://github.com/owncloud/ocis/issues/1542 and https://github.com/owncloud/ocis/pull/839

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
    And using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "find data" with the default quota using the Graph API
    And user "Alice" has created a folder "folderMain/SubFolder1/subFOLDER2" in space "find data"
    And user "Alice" has uploaded a file inside space "find data" with content "some content" to "folderMain/SubFolder1/subFOLDER2/insideTheFolder.txt"
    And using new DAV path


  Scenario: user can find data from the project space
    When user "Alice" searches for "fol" using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result should contain "4" entries
    And the search result of user "Alice" should contain these entries:
      | /folderMain                                           |
      | /folderMain/SubFolder1                                |
      | /folderMain/SubFolder1/subFOLDER2                     |
      | /folderMain/SubFolder1/subFOLDER2/insideTheFolder.txt |


  Scenario: user can only find data that they searched for from the project space
    When user "Alice" searches for "SUB" using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result should contain "2" entries
    And the search result of user "Alice" should contain these entries:
      | /folderMain/SubFolder1            |
      | /folderMain/SubFolder1/subFOLDER2 |
    But the search result of user "Alice" should not contain these entries:
      | /folderMain                                           |
      | /folderMain/SubFolder1/subFOLDER2/insideTheFolder.txt |


  Scenario: user can find data from the shares
    Given user "Alice" has created a share inside of space "find data" with settings:
      | path      | folderMain |
      | shareWith | Brian      |
      | role      | viewer     |
    And user "Brian" has accepted share "/folderMain" offered by user "Alice"
    When user "Brian" searches for "folder" using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result should contain "4" entries
    And the search result of user "Brian" should contain these entries:
      | folderMain/SubFolder1                                |
      | folderMain/SubFolder1/subFOLDER2                     |
      | folderMain/SubFolder1/subFOLDER2/insideTheFolder.txt |
    And for user "Brian" the search result should contain space "mountpoint/folderMain"


  Scenario: user can find hidden file
    Given user "Alice" has created a folder ".space" in space "find data"
    When user "Alice" searches for ".sp" using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result should contain "1" entries
    And the search result of user "Alice" should contain these entries:
      | /.space |


  Scenario: user cannot find pending share
    Given user "Alice" has created a share inside of space "find data" with settings:
      | path      | folderMain |
      | shareWith | Brian      |
      | role      | viewer     |
    When user "Brian" searches for "folder" using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result should contain "0" entries
    And the search result of user "Brian" should not contain these entries:
      | /SubFolder1                                |
      | /SubFolder1/subFOLDER2                     |
      | /SubFolder1/subFOLDER2/insideTheFolder.txt |


  Scenario: user cannot find declined share
    Given user "Alice" has created a share inside of space "find data" with settings:
      | path      | folderMain |
      | shareWith | Brian      |
      | role      | viewer     |
    And user "Brian" has declined share "/folderMain" offered by user "Alice"
    When user "Brian" searches for "folder" using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result should contain "0" entries
    And the search result of user "Brian" should not contain these entries:
      | /SubFolder1                                |
      | /SubFolder1/subFOLDER2                     |
      | /SubFolder1/subFOLDER2/insideTheFolder.txt |


  Scenario: user cannot find deleted folder
    Given user "Alice" has removed the folder "folderMain" from space "find data"
    When user "Alice" searches for "folderMain" using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result should contain "0" entries


  Scenario: user can find project space by name
    When user "Alice" searches for "find data" using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result should contain "1" entries
    And for user "Alice" the search result should contain space "find data"


  Scenario: user can search inside folder in space
    When user "Alice" searches for "folder" inside folder "/folderMain" in space "find data" using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result should contain "3" entries
    And the search result of user "Alice" should contain only these entries:
      | folderMain/SubFolder1                                |
      | folderMain/SubFolder1/subFOLDER2                     |
      | folderMain/SubFolder1/subFOLDER2/insideTheFolder.txt |
    But the search result of user "Alice" should not contain these entries:
      | /folderMain |


  Scenario: search inside folder in shares
    Given user "Alice" has created a share inside of space "find data" with settings:
      | path      | folderMain |
      | shareWith | Brian      |
      | role      | viewer     |
    And user "Brian" has accepted share "/folderMain" offered by user "Alice"
    When user "Brian" searches for "folder" inside folder "/folderMain" in space "Shares" using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result of user "Brian" should contain only these entries:
      | folderMain/SubFolder1                                |
      | folderMain/SubFolder1/subFOLDER2                     |
      | folderMain/SubFolder1/subFOLDER2/insideTheFolder.txt |
    But the search result of user "Brian" should not contain these entries:
      | /folderMain |
