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


  Scenario Outline: user can find data from the project space
    Given using <dav-path-version> DAV path
    When user "Alice" searches for "*fol*" using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result should contain "4" entries
    And the search result of user "Alice" should contain these entries:
      | /folderMain                                           |
      | /folderMain/SubFolder1                                |
      | /folderMain/SubFolder1/subFOLDER2                     |
      | /folderMain/SubFolder1/subFOLDER2/insideTheFolder.txt |
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: user can only find data that they searched for from the project space
    Given using <dav-path-version> DAV path
    When user "Alice" searches for "*SUB*" using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result should contain "2" entries
    And the search result of user "Alice" should contain these entries:
      | /folderMain/SubFolder1            |
      | /folderMain/SubFolder1/subFOLDER2 |
    But the search result of user "Alice" should not contain these entries:
      | /folderMain                                           |
      | /folderMain/SubFolder1/subFOLDER2/insideTheFolder.txt |
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: user can find data from the shares
    Given using <dav-path-version> DAV path
    And user "Alice" has created a share inside of space "find data" with settings:
      | path      | folderMain |
      | shareWith | Brian      |
      | role      | viewer     |
    When user "Brian" searches for "*folder*" using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result should contain "4" entries
    And the search result of user "Brian" should contain these entries:
      | folderMain/SubFolder1                                |
      | folderMain/SubFolder1/subFOLDER2                     |
      | folderMain/SubFolder1/subFOLDER2/insideTheFolder.txt |
    And for user "Brian" the search result should contain space "mountpoint/folderMain"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: user can find hidden file
    Given using <dav-path-version> DAV path
    And user "Alice" has created a folder ".space" in space "find data"
    When user "Alice" searches for "*.sp*" using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result should contain "1" entries
    And the search result of user "Alice" should contain these entries:
      | /.space |
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: user cannot find pending share
    Given user "Brian" has disabled auto-accepting
    And using <dav-path-version> DAV path
    And user "Alice" has created a share inside of space "find data" with settings:
      | path      | folderMain |
      | shareWith | Brian      |
      | role      | viewer     |
    When user "Brian" searches for "*folder*" using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result should contain "0" entries
    And the search result of user "Brian" should not contain these entries:
      | /SubFolder1                                |
      | /SubFolder1/subFOLDER2                     |
      | /SubFolder1/subFOLDER2/insideTheFolder.txt |
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: user cannot find declined share
    Given using <dav-path-version> DAV path
    And user "Alice" has created a share inside of space "find data" with settings:
      | path      | folderMain |
      | shareWith | Brian      |
      | role      | viewer     |
    And user "Brian" has declined share "/Shares/folderMain" offered by user "Alice"
    When user "Brian" searches for "*folder*" using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result should contain "0" entries
    And the search result of user "Brian" should not contain these entries:
      | /SubFolder1                                |
      | /SubFolder1/subFOLDER2                     |
      | /SubFolder1/subFOLDER2/insideTheFolder.txt |
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: user cannot find deleted folder
    Given using <dav-path-version> DAV path
    And user "Alice" has removed the folder "folderMain" from space "find data"
    When user "Alice" searches for "*folderMain*" using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result should contain "0" entries
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: user can find project space by name
    Given using <dav-path-version> DAV path
    When user "Alice" searches for '"*find data*"' using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result should contain "1" entries
    And for user "Alice" the search result should contain space "find data"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: user can search inside folder in space
    Given using <dav-path-version> DAV path
    When user "Alice" searches for "*folder*" inside folder "/folderMain" in space "find data" using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result should contain "3" entries
    And the search result of user "Alice" should contain only these entries:
      | folderMain/SubFolder1                                |
      | folderMain/SubFolder1/subFOLDER2                     |
      | folderMain/SubFolder1/subFOLDER2/insideTheFolder.txt |
    But the search result of user "Alice" should not contain these entries:
      | /folderMain |
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: search inside folder in shares
    Given using <dav-path-version> DAV path
    And user "Alice" has created a share inside of space "find data" with settings:
      | path      | folderMain |
      | shareWith | Brian      |
      | role      | viewer     |
    When user "Brian" searches for "*folder*" inside folder "/folderMain" in space "Shares" using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result of user "Brian" should contain only these entries:
      | folderMain/SubFolder1                                |
      | folderMain/SubFolder1/subFOLDER2                     |
      | folderMain/SubFolder1/subFOLDER2/insideTheFolder.txt |
    But the search result of user "Brian" should not contain these entries:
      | /folderMain |
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  @skipOnStable3.0
  Scenario Outline: search files inside the folder
    Given using <dav-path-version> DAV path
    And user "Alice" has uploaded file with content "hello world inside root" to "file1.txt"
    And user "Alice" has created folder "/Folder"
    And user "Alice" has uploaded file with content "hello world inside folder" to "/Folder/file2.txt"
    And user "Alice" has created folder "/Folder/SubFolder"
    And user "Alice" has uploaded file with content "hello world inside sub-folder" to "/Folder/SubFolder/file3.txt"
    When user "Alice" searches for "*file*" inside folder "/Folder" using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result of user "Alice" should contain only these entries:
      | /Folder/file2.txt           |
      | /Folder/SubFolder/file3.txt |
    But the search result of user "Alice" should not contain these entries:
      | file1.txt |
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |

  @issue-7114
  Scenario Outline: search files inside the folder with white space character in its name
    Given using <dav-path-version> DAV path
    And user "Alice" has created folder "/New Folder"
    And user "Alice" has uploaded file with content "hello world inside folder" to "/New Folder/file.txt"
    And user "Alice" has created folder "/New Folder/Sub Folder"
    And user "Alice" has uploaded file with content "hello world inside sub folder" to "/New Folder/Sub Folder/file1.txt"
    When user "Alice" searches for "*file*" inside folder "/New Folder" using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result of user "Alice" should contain only these entries:
      | /New Folder/file.txt             |
      | /New Folder/Sub Folder/file1.txt |
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |

  @issue-7114
  Scenario Outline: search files with white space character in its name
    Given using <dav-path-version> DAV path
    And user "Alice" has created folder "/New Folder"
    And user "Alice" has uploaded file with content "hello world" to "/new file.txt"
    And user "Alice" has created folder "/New Folder/New Sub Folder"
    When user "Alice" searches for "*new*" using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result of user "Alice" should contain only these entries:
      | /New Folder                |
      | /New Folder/New Sub Folder |
      | /new file.txt              |
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |

  @issue-enterprise-6000 @issue-7028 @issue-7092
  Scenario Outline: sharee cannot find resources that are not shared
    Given using <dav-path-version> DAV path
    And user "Alice" has created a folder "foo/sharedToBrian" in space "Alice Hansen"
    And user "Alice" has created a folder "sharedToCarol" in space "Alice Hansen"
    And user "Alice" has created a share inside of space "Alice Hansen" with settings:
      | path      | foo    |
      | shareWith | Brian  |
      | role      | viewer |
    When user "Brian" searches for "shared*" using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result of user "Brian" should contain these entries:
      | foo/sharedToBrian |
    But the search result of user "Brian" should not contain these entries:
      | /sharedToCarol |
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: search resources using different search patterns (KQL feature)
    Given using spaces DAV path
    And user "Alice" has created a folder "subfolder" in space "find data"
    When user "Alice" searches for '<pattern>' using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result should contain "1" entries
    And the search result of user "Alice" should contain these entries:
      | <search-result> |
    Examples:
      | description                                            | pattern      | search-result                     |
      | starts with                                            | fold*        | /folderMain                       |
      | ends with                                              | *der1        | /folderMain/SubFolder1            |
      | strict search                                          | subfolder    | /subfolder                        |
      | using patern "name:"                                   | name:*der2   | /folderMain/SubFolder1/subFOLDER2 |
      | using the pattern "name:" where the value is in quotes | name:"*der2" | /folderMain/SubFolder1/subFOLDER2 |
