@tikaServiceNeeded
Feature: content search
  As a user
  I want to do search resources by content
  So that I can find the files with the content I am looking for

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files


  Scenario Outline: search files by content
    Given using <dav-path-version> DAV path
    And user "Alice" has uploaded file with content "hello world from nepal" to "keywordAtStart.txt"
    And user "Alice" has uploaded file with content "saying hello to the world" to "keywordAtMiddle.txt"
    And user "Alice" has uploaded file with content "nepal want to say hello" to "keywordAtLast.txt"
    And user "Alice" has uploaded file with content "namaste from nepal" to "hello.txt"
    When user "Alice" searches for "Content:hello" using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result of user "Alice" should contain only these files:
      | keywordAtStart.txt  |
      | keywordAtMiddle.txt |
      | keywordAtLast.txt   |
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: sharee searches files by content
    Given using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has created folder "uploadFolder"
    And user "Alice" has uploaded file with content "hello world from nepal" to "uploadFolder/keywordAtStart.txt"
    And user "Alice" has uploaded file with content "saying hello to the world" to "uploadFolder/keywordAtMiddle.txt"
    And user "Alice" has uploaded file with content "nepal want to say hello" to "uploadFolder/keywordAtLast.txt"
    And user "Alice" has uploaded file with content "Namaste nepal" to "uploadFolder/hello.txt"
    And user "Alice" has shared folder "/uploadFolder" with user "Brian"
    And user "Brian" has accepted share "/uploadFolder" offered by user "Alice"
    When user "Brian" searches for "Content:hello" using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result of user "Brian" should contain only these files:
      | uploadFolder/keywordAtStart.txt  |
      | uploadFolder/keywordAtMiddle.txt |
      | uploadFolder/keywordAtLast.txt   |
    Examples:
      | dav-path-version |
      | old              |
      | new              |

    @skipOnStable3.0
    Examples:
      | dav-path-version |
      | spaces           |


  Scenario Outline: search deleted files by content
    Given using <dav-path-version> DAV path
    And user "Alice" has created folder "uploadFolder"
    And user "Alice" has uploaded file with content "hello world from nepal" to "uploadFolder/keywordAtStart.txt"
    And user "Alice" has uploaded file with content "saying hello to the world" to "keywordAtMiddle.txt"
    And user "Alice" has uploaded file with content "nepal want to say hello" to "keywordAtLast.txt"
    And user "Alice" has deleted file "keywordAtLast.txt"
    When user "Alice" searches for "Content:hello" using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result of user "Alice" should contain only these files:
      | uploadFolder/keywordAtStart.txt |
      | keywordAtMiddle.txt             |
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: search restored files by content
    Given using <dav-path-version> DAV path
    And user "Alice" has created folder "uploadFolder"
    And user "Alice" has uploaded file with content "hello world from nepal" to "keywordAtStart.txt"
    And user "Alice" has deleted file "keywordAtStart.txt"
    And user "Alice" has restored the file with original path "keywordAtStart.txt"
    When user "Alice" searches for "Content:hello" using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result of user "Alice" should contain only these files:
      | keywordAtStart.txt |
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: search restored version of a file by content
    Given using <dav-path-version> DAV path
    And user "Alice" has uploaded file with content "hello world" to "test.txt"
    And user "Alice" has uploaded file with content "Namaste nepal" to "test.txt"
    And user "Alice" has restored version index "1" of file "test.txt"
    When user "Alice" searches for "Content:hello" using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result of user "Alice" should contain only these files:
      | test.txt |
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: search project space files by content
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "project-space" with the default quota using the Graph API
    And user "Alice" has created a folder "spacesFolderWithFile/spacesSubFolder" in space "project-space"
    And user "Alice" has uploaded a file inside space "project-space" with content "hello world from nepal" to "keywordAtStart.txt"
    And user "Alice" has uploaded a file inside space "project-space" with content "saying hello to the world" to "spacesFolderWithFile/keywordAtMiddle.txt"
    And user "Alice" has uploaded a file inside space "project-space" with content "nepal want to say hello" to "spacesFolderWithFile/spacesSubFolder/keywordAtLast.txt"
    And user "Alice" has uploaded a file inside space "project-space" with content "namaste from nepal" to "hello.txt"
    And using <dav-path-version> DAV path
    When user "Alice" searches for "Content:hello" using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result of user "Alice" should contain only these files:
      | keywordAtStart.txt                                     |
      | spacesFolderWithFile/keywordAtMiddle.txt               |
      | spacesFolderWithFile/spacesSubFolder/keywordAtLast.txt |
    Examples:
      | dav-path-version |
      | old              |
      | new              |

    @skipOnStable3.0
    Examples:
      | dav-path-version |
      | spaces           |


  Scenario Outline: sharee searches shared project space files by content
    Given using spaces DAV path
    And user "Brian" has been created with default attributes and without skeleton files
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "project-space" with the default quota using the Graph API
    And user "Alice" has shared a space "project-space" with settings:
      | shareWith | Brian  |
      | role      | viewer |
    And user "Alice" has created a folder "spacesFolderWithFile/spacesSubFolder" in space "project-space"
    And user "Alice" has uploaded a file inside space "project-space" with content "hello world from nepal" to "keywordAtStart.txt"
    And user "Alice" has uploaded a file inside space "project-space" with content "saying hello to the world" to "spacesFolderWithFile/keywordAtMiddle.txt"
    And user "Alice" has uploaded a file inside space "project-space" with content "nepal wants to say hello" to "spacesFolderWithFile/spacesSubFolder/keywordAtLast.txt"
    And user "Alice" has uploaded a file inside space "project-space" with content "namaste from nepal" to "hello.txt"
    And using <dav-path-version> DAV path
    When user "Brian" searches for "Content:hello" using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result of user "Alice" should contain only these files:
      | keywordAtStart.txt                                     |
      | spacesFolderWithFile/keywordAtMiddle.txt               |
      | spacesFolderWithFile/spacesSubFolder/keywordAtLast.txt |
    Examples:
      | dav-path-version |
      | old              |
      | new              |

    @skipOnStable3.0
    Examples:
      | dav-path-version |
      | spaces           |
