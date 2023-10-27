Feature: moving/renaming file using file id
  As a user
  I want to be able to move or rename files using file id
  So that I can manage my file system

  Background:
    Given using spaces DAV path
    And user "Alice" has been created with default attributes and without skeleton files


  Scenario Outline: move a file into a folder inside personal space
    Given user "Alice" has created folder "/folder"
    And user "Alice" has uploaded file with content "some data" to "/textfile.txt"
    And we save it into "FILEID"
    When user "Alice" moves a file "/textfile.txt" into "/folder" inside space "Personal" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Alice" folder "folder" of the space "Personal" should contain these files:
      | textfile.txt |
    But for user "Alice" the space "Personal" should not contain these entries:
      | textfile.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: move a file into a sub-folder inside personal space
    Given user "Alice" has created folder "/folder"
    And user "Alice" has created folder "folder/sub-folder"
    And user "Alice" has uploaded file with content "some data" to "/textfile.txt"
    And we save it into "FILEID"
    When user "Alice" moves a file "/textfile.txt" into "/folder/sub-folder" inside space "Personal" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Alice" folder "folder/sub-folder/" of the space "Personal" should contain these files:
      | textfile.txt |
    But for user "Alice" the space "Personal" should not contain these entries:
      | textfile.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: move a file from folder to root inside personal space
    Given user "Alice" has created folder "/folder"
    And user "Alice" has uploaded file with content "some data" to "folder/textfile.txt"
    And we save it into "FILEID"
    When user "Alice" moves a file "folder/textfile.txt" into "/" inside space "Personal" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Alice" the space "Personal" should contain these entries:
      | textfile.txt |
    But for user "Alice" folder "folder" of the space "Personal" should not contain these files:
      | textfile.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: move a file from sub-folder to root inside personal space
    Given user "Alice" has created folder "/folder"
    And user "Alice" has created folder "folder/sub-folder"
    And user "Alice" has uploaded file with content "some data" to "folder/sub-folder/textfile.txt"
    And we save it into "FILEID"
    When user "Alice" moves a file "folder/sub-folder/textfile.txt" into "/" inside space "Personal" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Alice" the space "Personal" should contain these entries:
      | textfile.txt |
    But for user "Alice" folder "folder/sub-folder" of the space "Personal" should not contain these files:
      | textfile.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: move a file into a folder inside project space
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "project-space" with the default quota using the Graph API
    And user "Alice" has created a folder "/folder" in space "project-space"
    And user "Alice" has uploaded a file inside space "project-space" with content "some data" to "textfile.txt"
    And we save it into "FILEID"
    When user "Alice" moves a file "/textfile.txt" into "/folder" inside space "project-space" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Alice" folder "folder" of the space "project-space" should contain these files:
      | textfile.txt |
    But for user "Alice" the space "project-space" should not contain these entries:
      | textfile.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: move a file into a sub-folder inside project space
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "project-space" with the default quota using the Graph API
    And user "Alice" has created a folder "folder/sub-folder" in space "project-space"
    And user "Alice" has uploaded a file inside space "project-space" with content "some data" to "/textfile.txt"
    And we save it into "FILEID"
    When user "Alice" moves a file "/textfile.txt" into "/folder/sub-folder" inside space "project-space" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Alice" folder "folder/sub-folder/" of the space "project-space" should contain these files:
      | textfile.txt |
    But for user "Alice" the space "Personal" should not contain these entries:
      | textfile.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: move a file from folder to root inside project space
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "project-space" with the default quota using the Graph API
    And user "Alice" has created a folder "folder" in space "project-space"
    And user "Alice" has uploaded a file inside space "project-space" with content "some data" to "folder/textfile.txt"
    And we save it into "FILEID"
    When user "Alice" moves a file "folder/textfile.txt" into "/" inside space "project-space" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Alice" the space "project-space" should contain these entries:
      | textfile.txt |
    But for user "Alice" folder "folder" of the space "project-space" should not contain these files:
      | textfile.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: move a file from sub-folder to root inside project space
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "project-space" with the default quota using the Graph API
    And user "Alice" has created a folder "folder/sub-folder" in space "project-space"
    And user "Alice" has uploaded a file inside space "project-space" with content "some data" to "folder/sub-folder/textfile.txt"
    And we save it into "FILEID"
    When user "Alice" moves a file "folder/sub-folder/textfile.txt" into "/" inside space "project-space" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Alice" the space "project-space" should contain these entries:
      | textfile.txt |
    But for user "Alice" folder "folder/sub-folder" of the space "project-space" should not contain these files:
      | textfile.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: move a file from sub-folder to root folder inside shares space
    Given user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has created folder "/folder"
    And user "Alice" has created folder "folder/sub-folder"
    And user "Alice" has uploaded file with content "some data" to "/folder/sub-folder/test.txt"
    And we save it into "FILEID"
    And user "Alice" has shared folder "/folder" with user "Brian" with permissions "all"
    And user "Brian" has accepted share "/folder" offered by user "Alice"
    When user "Brian" moves a file "Shares/folder/sub-folder/test.txt" into "Shares/folder" inside space "Shares" using file-id path "<dav-path>"
    Then the HTTP status code should be "502"
    And the value of the item "/d:error/s:message" in the response about user "Brian" should be "gateway does not support cross storage move, use copy and delete"
    And for user "Brian" folder "folder/sub-folder" of the space "Shares" should contain these files:
      | test.txt |
    And for user "Brian" folder "folder" of the space "Shares" should not contain these files:
      | test.txt |
    And for user "Alice" folder "folder/sub-folder" of the space "Personal" should contain these files:
      | test.txt |
    And for user "Alice" folder "folder" of the space "Personal" should not contain these files:
      | test.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: rename a root file inside personal space
    Given user "Alice" has uploaded file with content "some data" to "textfile.txt"
    And we save it into "FILEID"
    When user "Alice" renames a file "textfile.txt" into "renamed.txt" inside space "Personal" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Alice" the space "Personal" should contain these entries:
      | renamed.txt |
    But for user "Alice" the space "Personal" should not contain these entries:
      | textfile.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: rename a file and move into a folder inside personal space
    Given user "Alice" has created folder "/folder"
    And user "Alice" has uploaded file with content "some data" to "/textfile.txt"
    And we save it into "FILEID"
    When user "Alice" renames a file "textfile.txt" into "/folder/renamed.txt" inside space "Personal" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Alice" folder "folder" of the space "Personal" should contain these files:
      | renamed.txt |
    But for user "Alice" the space "Personal" should not contain these entries:
      | textfile.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: rename a file and move into a sub-folder inside personal space
    Given user "Alice" has created folder "/folder"
    And user "Alice" has created folder "folder/sub-folder"
    And user "Alice" has uploaded file with content "some data" to "/textfile.txt"
    And we save it into "FILEID"
    When user "Alice" renames a file "textfile.txt" into "/folder/sub-folder/renamed.txt" inside space "Personal" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Alice" folder "folder/sub-folder" of the space "Personal" should contain these files:
      | renamed.txt |
    But for user "Alice" the space "Personal" should not contain these entries:
      | textfile.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: rename a file and move from a folder to root inside personal space
    Given user "Alice" has created folder "/folder"
    And user "Alice" has uploaded file with content "some data" to "folder/textfile.txt"
    And we save it into "FILEID"
    When user "Alice" renames a file "folder/textfile.txt" into "/renamed.txt" inside space "Personal" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Alice" the space "Personal" should contain these entries:
      | renamed.txt |
    But for user "Alice" folder "folder" of the space "Personal" should not contain these files:
      | renamed.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: rename a file and move from sub-folder to root inside personal space
    Given user "Alice" has created folder "/folder"
    And user "Alice" has created folder "folder/sub-folder"
    And user "Alice" has uploaded file with content "some data" to "folder/sub-folder/textfile.txt"
    And we save it into "FILEID"
    When user "Alice" renames a file "folder/sub-folder/textfile.txt" into "/renamed.txt" inside space "Personal" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Alice" the space "Personal" should contain these entries:
      | renamed.txt |
    But for user "Alice" folder "folder/sub-folder" of the space "Personal" should not contain these files:
      | textfile.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |
