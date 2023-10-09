Feature: copying file using file id
  As a user
  I want to copy the file using file id
  So that I can manage my resource

  Background:
    Given using spaces DAV path
    And user "Alice" has been created with default attributes and without skeleton files


  Scenario Outline: copy a file into a folder in personal space
    Given user "Alice" has created folder "/folder"
    And user "Alice" has uploaded file with content "some data" to "/textfile.txt"
    And we save it into "FILEID"
    When user "Alice" copies a file "/textfile.txt" into "/folder" inside space "Personal" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Alice" folder "/" of the space "Personal" should contain these files:
      | textfile.txt |
    And for user "Alice" folder "folder" of the space "Personal" should contain these files:
      | textfile.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: copy a file into a sub-folder in personal space
    Given user "Alice" has created folder "/folder"
    And user "Alice" has created folder "folder/sub-folder"
    And user "Alice" has uploaded file with content "some data" to "/textfile.txt"
    And we save it into "FILEID"
    When user "Alice" copies a file "/textfile.txt" into "/folder/sub-folder" inside space "Personal" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Alice" folder "/" of the space "Personal" should contain these files:
      | textfile.txt |
    And for user "Alice" folder "folder/sub-folder" of the space "Personal" should contain these files:
      | textfile.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: copy a file from a folder into root of personal space
    Given user "Alice" has created folder "/folder"
    And user "Alice" has uploaded file with content "some data" to "folder/textfile.txt"
    And we save it into "FILEID"
    When user "Alice" copies a file "folder/textfile.txt" into "/" inside space "Personal" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Alice" folder "/" of the space "Personal" should contain these files:
      | textfile.txt |
    And for user "Alice" folder "folder" of the space "Personal" should contain these files:
      | textfile.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: copy a file from sub-folder into root of personal space
    Given user "Alice" has created folder "/folder"
    And user "Alice" has created folder "folder/sub-folder"
    And user "Alice" has uploaded file with content "some data" to "folder/sub-folder/textfile.txt"
    And we save it into "FILEID"
    When user "Alice" copies a file "folder/sub-folder/textfile.txt" into "/" inside space "Personal" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Alice" folder "/" of the space "Personal" should contain these files:
      | textfile.txt |
    And for user "Alice" folder "folder/sub-folder" of the space "Personal" should contain these files:
      | textfile.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: copy a file into a folder in project space
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "project-space" with the default quota using the GraphApi
    And user "Alice" has created a folder "/folder" in space "project-space"
    And user "Alice" has uploaded a file inside space "project-space" with content "some data" to "textfile.txt"
    And we save it into "FILEID"
    When user "Alice" copies a file "/textfile.txt" into "/folder" inside space "project-space" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Alice" folder "/" of the space "project-space" should contain these files:
      | textfile.txt |
    And for user "Alice" folder "folder" of the space "project-space" should contain these files:
      | textfile.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: copy a file into a sub-folder in project space
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "project-space" with the default quota using the GraphApi
    And user "Alice" has created a folder "folder/sub-folder" in space "project-space"
    And user "Alice" has uploaded a file inside space "project-space" with content "some data" to "textfile.txt"
    And we save it into "FILEID"
    When user "Alice" copies a file "/textfile.txt" into "/folder/sub-folder" inside space "project-space" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Alice" folder "/" of the space "project-space" should contain these files:
      | textfile.txt |
    And for user "Alice" folder "folder/sub-folder" of the space "project-space" should contain these files:
      | textfile.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: copy a file from a folder into root of project space
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "project-space" with the default quota using the GraphApi
    And user "Alice" has created a folder "folder" in space "project-space"
    And user "Alice" has uploaded a file inside space "project-space" with content "some data" to "folder/textfile.txt"
    And we save it into "FILEID"
    When user "Alice" copies a file "folder/textfile.txt" into "/" inside space "project-space" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Alice" folder "/" of the space "project-space" should contain these files:
      | textfile.txt |
    And for user "Alice" folder "folder" of the space "project-space" should contain these files:
      | textfile.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: copy a file from sub-folder into root of project space
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "project-space" with the default quota using the GraphApi
    And user "Alice" has created a folder "folder/sub-folder" in space "project-space"
    And user "Alice" has uploaded a file inside space "project-space" with content "some data" to "folder/sub-folder/textfile.txt"
    And we save it into "FILEID"
    When user "Alice" copies a file "folder/sub-folder/textfile.txt" into "/" inside space "project-space" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Alice" folder "/" of the space "project-space" should contain these files:
      | textfile.txt |
    And for user "Alice" folder "folder/sub-folder" of the space "project-space" should contain these files:
      | textfile.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |
