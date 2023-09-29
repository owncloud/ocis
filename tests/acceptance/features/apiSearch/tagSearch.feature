Feature: tag search
  As a user
  I want to do search resources by tag
  So that I can find the files with the tag I am looking for

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files


  Scenario Outline: search files by tag
    Given using <dav-path-version> DAV path
    And user "Alice" has created the following folders
      | path                      |
      | folderWithFile            |
      | folderWithFile/subFolder/ |
    And user "Alice" has uploaded the following files with content "some data"
      | path                                             |
      | fileInRootLevel.txt                              |
      | folderWithFile/fileInsideFolder.txt              |
      | folderWithFile/subFolder/fileInsideSubFolder.txt |
    And user "Alice" has tagged the following files of the space "Personal":
      | path                                             | tagName |
      | fileInRootLevel.txt                              | tag1    |
      | folderWithFile/fileInsideFolder.txt              | tag1    |
      | folderWithFile/subFolder/fileInsideSubFolder.txt | tag1    |
    When user "Alice" searches for "Tags:tag1" using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result of user "Alice" should contain only these files:
      | fileInRootLevel.txt                              |
      | folderWithFile/fileInsideFolder.txt              |
      | folderWithFile/subFolder/fileInsideSubFolder.txt |
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: search project space files by tag
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "tag-space" with the default quota using the Graph API
    And user "Alice" has created a folder "spacesFolderWithFile/spacesSubFolder" in space "tag-space"
    And user "Alice" has uploaded a file inside space "tag-space" with content "tagged file" to "spacesFile.txt"
    And user "Alice" has uploaded a file inside space "tag-space" with content "untagged file" to "spacesFileWithoutTag.txt"
    And user "Alice" has uploaded a file inside space "tag-space" with content "tagged file in folder" to "spacesFolderWithFile/spacesFileInsideFolder.txt"
    And user "Alice" has uploaded a file inside space "tag-space" with content "tagged file in subfolder" to "spacesFolderWithFile/spacesSubFolder/spacesFileInsideSubFolder.txt"
    And user "Alice" has tagged the following files of the space "tag-space":
      | path                                                               | tagName |
      | spacesFile.txt                                                     | tag1    |
      | spacesFolderWithFile/spacesFileInsideFolder.txt                    | tag1    |
      | spacesFolderWithFile/spacesSubFolder/spacesFileInsideSubFolder.txt | tag1    |
    And using <dav-path-version> DAV path
    When user "Alice" searches for "Tags:tag1" using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result of user "Alice" should contain only these files:
      | spacesFile.txt                                                     |
      | spacesFolderWithFile/spacesFileInsideFolder.txt                    |
      | spacesFolderWithFile/spacesSubFolder/spacesFileInsideSubFolder.txt |
    Examples:
      | dav-path-version |
      | old              |
      | new              |

    @skipOnStable3.0
    Examples:
      | dav-path-version |
      | spaces           |


  Scenario Outline: search folders using a tag
    Given using <dav-path-version> DAV path
    And user "Alice" has created folder "uploadFolder1"
    And user "Alice" has created folder "uploadFolder2"
    And user "Alice" has created folder "uploadFolder3"
    And user "Alice" has tagged the following folders of the space "Personal":
      | path          | tagName |
      | uploadFolder1 | tag1    |
      | uploadFolder2 | tag1    |
    When user "Alice" searches for "Tags:tag1" using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result of user "Alice" should contain only these entries:
      | uploadFolder1 |
      | uploadFolder2 |
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: search project space folders by tag
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "tag-space" with the default quota using the Graph API
    And user "Alice" has created a folder "spacesFolder/spacesSubFolder" in space "tag-space"
    And user "Alice" has created a folder "unTagSpacesFolder/unTagSpacesSubFolder" in space "tag-space"
    And user "Alice" has tagged the following folders of the space "tag-space":
      | path                         | tagName |
      | spacesFolder                 | tag1    |
      | spacesFolder/spacesSubFolder | tag1    |
    And using <dav-path-version> DAV path
    When user "Alice" searches for "Tags:tag1" using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result of user "Alice" should contain only these files:
      | spacesFolder                 |
      | spacesFolder/spacesSubFolder |
    Examples:
      | dav-path-version |
      | old              |
      | new              |

    @skipOnStable3.0
    Examples:
      | dav-path-version |
      | spaces           |


  Scenario Outline: sharee searches shared files using a tag
    Given using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has created folder "uploadFolder"
    And user "Alice" has uploaded file with content "hello world" to "uploadFolder/file1.txt"
    And user "Alice" has uploaded file with content "Namaste nepal" to "uploadFolder/file2.txt"
    And user "Alice" has uploaded file with content "hello nepal" to "uploadFolder/file3.txt"
    And user "Alice" has created the following tags for file "uploadFolder/file1.txt" of the space "Personal":
      | tag1 |
    And user "Alice" has shared folder "/uploadFolder" with user "Brian"
    And user "Brian" has accepted share "/uploadFolder" offered by user "Alice"
    And user "Brian" has created the following tags for file "uploadFolder/file2.txt" of the space "Shares":
      | tag1 |
    When user "Brian" searches for "Tags:tag1" using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result of user "Brian" should contain only these files:
      | uploadFolder/file1.txt |
      | uploadFolder/file2.txt |
    Examples:
      | dav-path-version |
      | old              |
      | new              |

    @skipOnStable3.0
    Examples:
      | dav-path-version |
      | spaces           |


  Scenario Outline: sharee searches shared project space files using a tag
    Given using spaces DAV path
    And user "Brian" has been created with default attributes and without skeleton files
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "tag-space" with the default quota using the Graph API
    And user "Alice" has shared a space "tag-space" with settings:
      | shareWith | Brian  |
      | role      | viewer |
    And user "Alice" has created a folder "spacesFolderWithFile/spacesSubFolder" in space "tag-space"
    And user "Alice" has uploaded a file inside space "tag-space" with content "tagged file" to "spacesFile.txt"
    And user "Alice" has uploaded a file inside space "tag-space" with content "untagged file" to "spacesFileWithoutTag.txt"
    And user "Alice" has uploaded a file inside space "tag-space" with content "tagged file in folder" to "spacesFolderWithFile/spacesFileInsideFolder.txt"
    And user "Alice" has uploaded a file inside space "tag-space" with content "tagged file in subfolder" to "spacesFolderWithFile/spacesSubFolder/spacesFileInsideSubFolder.txt"
    And user "Alice" has tagged the following files of the space "tag-space":
      | path                                                               | tagName |
      | spacesFile.txt                                                     | tag1    |
      | spacesFolderWithFile/spacesFileInsideFolder.txt                    | tag1    |
      | spacesFolderWithFile/spacesSubFolder/spacesFileInsideSubFolder.txt | tag1    |
    And using <dav-path-version> DAV path
    When user "Brian" searches for "Tags:tag1" using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result of user "Alice" should contain only these files:
      | spacesFile.txt                                                     |
      | spacesFolderWithFile/spacesFileInsideFolder.txt                    |
      | spacesFolderWithFile/spacesSubFolder/spacesFileInsideSubFolder.txt |
    Examples:
      | dav-path-version |
      | old              |
      | new              |

    @skipOnStable3.0
    Examples:
      | dav-path-version |
      | spaces           |


  Scenario Outline: search files using a deleted tag
    Given using <dav-path-version> DAV path
    And user "Alice" has uploaded file with content "hello world" to "file1.txt"
    And user "Alice" has created the following tags for file "file1.txt" of the space "Personal":
      | tag1 |
    And user "Alice" has removed the following tags for file "file1.txt" of space "Personal":
      | tag1 |
    When user "Alice" searches for "Tags:tag1" using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result should contain "0" entries
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: search restored files using a tag
    Given using <dav-path-version> DAV path
    And user "Alice" has uploaded file with content "hello world" to "file1.txt"
    And user "Alice" has uploaded file with content "Namaste nepal" to "file2.txt"
    And user "Alice" has created the following tags for file "file1.txt" of the space "Personal":
      | tag1 |
    And user "Alice" has deleted file "/file1.txt"
    And user "Alice" has restored the file with original path "/file1.txt"
    When user "Alice" searches for "Tags:tag1" using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result of user "Alice" should contain only these entries:
      | /file1.txt |
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: search restored version of a file using a tag
    Given using <dav-path-version> DAV path
    And user "Alice" has uploaded file with content "version one file" to "file.txt"
    And user "Alice" has created the following tags for file "file.txt" of the space "Personal":
      | tag1 |
    And user "Alice" has uploaded file with content "version two file" to "file.txt"
    And user "Alice" has restored version index "1" of file "file.txt"
    When user "Alice" searches for "Tags:tag1" using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result of user "Alice" should contain only these entries:
      | /file.txt |
    And the content of file "file.txt" for user "Alice" should be "version one file"
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
    When user "Alice" searches for "file" inside folder "/Folder" using the WebDAV API
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
