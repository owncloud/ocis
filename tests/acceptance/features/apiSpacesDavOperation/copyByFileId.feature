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
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "project-space" with the default quota using the Graph API
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
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "project-space" with the default quota using the Graph API
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
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "project-space" with the default quota using the Graph API
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
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "project-space" with the default quota using the Graph API
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


  Scenario Outline: copy a file from personal to project space
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "project-space" with the default quota using the Graph API
    And user "Alice" has uploaded file with content "some data" to "textfile.txt"
    And we save it into "FILEID"
    When user "Alice" copies a file "/textfile.txt" into "/" inside space "project-space" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Alice" folder "/" of the space "project-space" should contain these files:
      | textfile.txt |
    And for user "Alice" folder "/" of the space "Personal" should contain these files:
      | textfile.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: copy a file from project to personal space
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "project-space" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "project-space" with content "some data" to "textfile.txt"
    And we save it into "FILEID"
    When user "Alice" copies a file "/textfile.txt" into "/" inside space "Personal" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Alice" folder "/" of the space "project-space" should contain these files:
      | textfile.txt |
    And for user "Alice" folder "/" of the space "Personal" should contain these files:
      | textfile.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: copy a file from sub-folder to root folder inside Shares space
    Given user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has created folder "/folder"
    And user "Alice" has created folder "folder/sub-folder"
    And user "Alice" has uploaded file with content "some data" to "/folder/sub-folder/test.txt"
    And we save it into "FILEID"
    And user "Alice" has sent the following resource share invitation:
      | resource        | folder   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "folder" synced
    When user "Brian" copies a file "Shares/folder/sub-folder/test.txt" into "Shares/folder" inside space "Shares" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Brian" folder "folder" of the space "Shares" should contain these files:
      | test.txt |
    And for user "Brian" folder "folder/sub-folder" of the space "Shares" should contain these files:
      | test.txt |
    And for user "Alice" folder "folder" of the space "Personal" should contain these files:
      | test.txt |
    And for user "Alice" folder "folder/sub-folder" of the space "Personal" should contain these files:
      | test.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: copy a file from personal to share space
    Given user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has created folder "/folder"
    And user "Alice" has sent the following resource share invitation:
      | resource        | folder   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "folder" synced
    And user "Brian" has uploaded file with content "some data" to "/test.txt"
    And we save it into "FILEID"
    And user "Brian" has a share "folder" synced
    When user "Brian" copies a file "/test.txt" into "Shares/folder" inside space "Shares" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Brian" folder "folder" of the space "Shares" should contain these files:
      | test.txt |
    And for user "Brian" folder "/" of the space "Personal" should contain these files:
      | test.txt |
    And for user "Alice" folder "folder" of the space "Personal" should contain these files:
      | test.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: copy a file from share to personal space
    Given user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has created folder "/folder"
    And user "Alice" has uploaded file with content "some data" to "/folder/test.txt"
    And we save it into "FILEID"
    And user "Alice" has sent the following resource share invitation:
      | resource        | folder            |
      | space           | Personal          |
      | sharee          | Brian             |
      | shareType       | user              |
      | permissionsRole | <permission-role> |
    And user "Brian" has a share "folder" synced
    When user "Brian" copies a file "/test.txt" into "/" inside space "Personal" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Brian" folder "folder" of the space "Shares" should contain these files:
      | test.txt |
    And for user "Brian" folder "/" of the space "Personal" should contain these files:
      | test.txt |
    And for user "Alice" folder "folder" of the space "Personal" should contain these files:
      | test.txt |
    Examples:
      | permission-role | dav-path                          |
      | Editor          | /remote.php/dav/spaces/<<FILEID>> |
      | Viewer          | /remote.php/dav/spaces/<<FILEID>> |
      | Uploader        | /remote.php/dav/spaces/<<FILEID>> |
      | Editor          | /dav/spaces/<<FILEID>>            |
      | Viewer          | /dav/spaces/<<FILEID>>            |
      | Uploader        | /dav/spaces/<<FILEID>>            |


  Scenario Outline: copy a file between two project spaces
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "first-project-space" with the default quota using the Graph API
    And user "Alice" has created a space "second-project-space" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "first-project-space" with content "data from first project space" to "firstProjectSpacetextfile.txt"
    And user "Alice" has uploaded a file inside space "second-project-space" with content "data from second project space" to "secondProjectSpacetextfile.txt"
    And we save it into "FILEID"
    When user "Alice" copies a file "/secondProjectSpacetextfile.txt" into "/" inside space "first-project-space" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Alice" folder "/" of the space "first-project-space" should contain these files:
      | firstProjectSpacetextfile.txt  |
      | secondProjectSpacetextfile.txt |
    And for user "Alice" folder "/" of the space "second-project-space" should contain these files:
      | secondProjectSpacetextfile.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: sharee tries to copy a file from shares space with secure viewer to personal space
    Given user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has created folder "/folder"
    And user "Alice" has uploaded file with content "some data" to "/folder/test.txt"
    And we save it into "FILEID"
    And user "Alice" has sent the following resource share invitation:
      | resource        | folder        |
      | space           | Personal      |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | Secure viewer |
    And user "Brian" has a share "folder" synced
    When user "Brian" copies a file "/test.txt" into "/" inside space "Personal" using file-id path "<dav-path>"
    Then the HTTP status code should be "403"
    And for user "Brian" folder "folder" of the space "Shares" should contain these files:
      | test.txt |
    And for user "Brian" folder "/" of the space "Personal" should not contain these files:
      | test.txt |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: sharee copies a file from shares to project space
    Given user "Brian" has been created with default attributes and without skeleton files
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created folder "/folder"
    And user "Alice" has uploaded file with content "some data" to "/folder/test.txt"
    And we save it into "FILEID"
    And user "Alice" has sent the following resource share invitation:
      | resource        | folder            |
      | space           | Personal          |
      | sharee          | Brian             |
      | shareType       | user              |
      | permissionsRole | <permission-role> |
    And user "Brian" has a share "folder" synced
    And user "Alice" has created a space "project-space" with the default quota using the Graph API
    And user "Alice" has sent the following space share invitation:
      | space           | project-space |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | <space-role>  |
    When user "Brian" copies a file "Shares/folder/test.txt" into "/" inside space "project-space" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Brian" folder "folder" of the space "Shares" should contain these files:
      | test.txt |
    And for user "Brian" folder "/" of the space "project-space" should contain these files:
      | test.txt |
    And for user "Alice" folder "/" of the space "project-space" should contain these files:
      | test.txt |
    Examples:
      | permission-role | space-role    | dav-path                          |
      | Viewer          | Manager       | /remote.php/dav/spaces/<<FILEID>> |
      | Viewer          | Space Editor  | /remote.php/dav/spaces/<<FILEID>> |
      | Editor          | Manager       | /remote.php/dav/spaces/<<FILEID>> |
      | Editor          | Space Editor  | /remote.php/dav/spaces/<<FILEID>> |
      | Uploader        | Manager       | /remote.php/dav/spaces/<<FILEID>> |
      | Uploader        | Space Editor  | /remote.php/dav/spaces/<<FILEID>> |
      | Viewer          | Manager       | /dav/spaces/<<FILEID>>            |
      | Viewer          | Space Editor  | /dav/spaces/<<FILEID>>            |
      | Editor          | Manager       | /dav/spaces/<<FILEID>>            |
      | Editor          | Space Editor  | /dav/spaces/<<FILEID>>            |
      | Uploader        | Manager       | /dav/spaces/<<FILEID>>            |
      | Uploader        | Space Editor  | /dav/spaces/<<FILEID>>            |


  Scenario Outline: sharee tries to copy a file from shares to project space
    Given user "Brian" has been created with default attributes and without skeleton files
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created folder "/folder"
    And user "Alice" has uploaded file with content "some data" to "/folder/test.txt"
    And we save it into "FILEID"
    And user "Alice" has sent the following resource share invitation:
      | resource        | folder            |
      | space           | Personal          |
      | sharee          | Brian             |
      | shareType       | user              |
      | permissionsRole | <permission-role> |
    And user "Brian" has a share "folder" synced
    And user "Alice" has created a space "project-space" with the default quota using the Graph API
    And user "Alice" has sent the following space share invitation:
      | space           | project-space |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | <space-role>  |
    When user "Brian" copies a file "Shares/folder/test.txt" into "/" inside space "project-space" using file-id path "<dav-path>"
    Then the HTTP status code should be "403"
    And for user "Brian" folder "folder" of the space "Shares" should contain these files:
      | test.txt |
    And for user "Brian" folder "/" of the space "project-space" should not contain these files:
      | test.txt |
    And for user "Alice" folder "/" of the space "project-space" should not contain these files:
      | test.txt |
    Examples:
      | permission-role | space-role    | dav-path                          |
      | Secure viewer   | Manager       | /remote.php/dav/spaces/<<FILEID>> |
      | Secure viewer   | Space Viewer  | /remote.php/dav/spaces/<<FILEID>> |
      | Secure viewer   | Space Editor  | /remote.php/dav/spaces/<<FILEID>> |
      | Editor          | Space Viewer  | /remote.php/dav/spaces/<<FILEID>> |
      | Viewer          | Space Viewer  | /remote.php/dav/spaces/<<FILEID>> |
      | Uploader        | Space Viewer  | /remote.php/dav/spaces/<<FILEID>> |
      | Secure viewer   | Manager       | /dav/spaces/<<FILEID>>            |
      | Secure viewer   | Space Viewer  | /dav/spaces/<<FILEID>>            |
      | Secure viewer   | Space Editor  | /dav/spaces/<<FILEID>>            |
      | Editor          | Space Viewer  | /dav/spaces/<<FILEID>>            |
      | Viewer          | Space Viewer  | /dav/spaces/<<FILEID>>            |
      | Uploader        | Space Viewer  | /dav/spaces/<<FILEID>>            |


  Scenario Outline: sharee copies a file between shares spaces
    Given user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has created folder "/share1"
    And user "Alice" has created folder "/share2"
    And user "Alice" has uploaded file with content "some data" to "/share1/test.txt"
    And we save it into "FILEID"
    And user "Alice" has sent the following resource share invitation:
      | resource        | share1            |
      | space           | Personal          |
      | sharee          | Brian             |
      | shareType       | user              |
      | permissionsRole | <from-share-role> |
    And user "Brian" has a share "share1" synced
    And user "Alice" has sent the following resource share invitation:
      | resource        | share2          |
      | space           | Personal        |
      | sharee          | Brian           |
      | shareType       | user            |
      | permissionsRole | <to-share-role> |
    And user "Brian" has a share "share2" synced
    When user "Brian" copies a file "Shares/share1/test.txt" into "share2" inside space "Shares" using file-id path "<dav-path>"
    Then the HTTP status code should be "201"
    And for user "Brian" folder "share1" of the space "Shares" should contain these files:
      | test.txt |
    And for user "Brian" folder "share2" of the space "Shares" should contain these files:
      | test.txt |
    And for user "Alice" folder "share1" of the space "Personal" should contain these files:
      | test.txt |
    And for user "Alice" folder "share2" of the space "Personal" should contain these files:
      | test.txt |
    Examples:
      | from-share-role | to-share-role | dav-path                          |
      | Viewer          | Editor        | /remote.php/dav/spaces/<<FILEID>> |
      | Viewer          | Uploader      | /remote.php/dav/spaces/<<FILEID>> |
      | Editor          | Editor        | /remote.php/dav/spaces/<<FILEID>> |
      | Editor          | Uploader      | /remote.php/dav/spaces/<<FILEID>> |
      | Uploader        | Editor        | /remote.php/dav/spaces/<<FILEID>> |
      | Uploader        | Uploader      | /remote.php/dav/spaces/<<FILEID>> |
      | Viewer          | Editor        | /dav/spaces/<<FILEID>>            |
      | Viewer          | Uploader      | /dav/spaces/<<FILEID>>            |
      | Editor          | Editor        | /dav/spaces/<<FILEID>>            |
      | Editor          | Uploader      | /dav/spaces/<<FILEID>>            |
      | Uploader        | Editor        | /dav/spaces/<<FILEID>>            |
      | Uploader        | Uploader      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: sharee tries to copy a file between shares space
    Given user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has created folder "/share1"
    And user "Alice" has created folder "/share2"
    And user "Alice" has uploaded file with content "some data" to "/share1/test.txt"
    And we save it into "FILEID"
    And user "Alice" has sent the following resource share invitation:
      | resource        | share1            |
      | space           | Personal          |
      | sharee          | Brian             |
      | shareType       | user              |
      | permissionsRole | <from-share-role> |
    And user "Brian" has a share "share1" synced
    And user "Alice" has sent the following resource share invitation:
      | resource        | share2          |
      | space           | Personal        |
      | sharee          | Brian           |
      | shareType       | user            |
      | permissionsRole | <to-share-role> |
    And user "Brian" has a share "share2" synced
    When user "Brian" copies a file "Shares/share1/test.txt" into "share2" inside space "Shares" using file-id path "<dav-path>"
    Then the HTTP status code should be "403"
    And for user "Brian" folder "share1" of the space "Shares" should contain these files:
      | test.txt |
    And for user "Brian" folder "share2" of the space "Shares" should not contain these files:
      | test.txt |
    And for user "Alice" folder "share1" of the space "Personal" should contain these files:
      | test.txt |
    And for user "Alice" folder "share2" of the space "Personal" should not contain these files:
      | test.txt |
    Examples:
      | from-share-role | to-share-role | dav-path                          |
      | Secure viewer   | Viewer        | /remote.php/dav/spaces/<<FILEID>> |
      | Secure viewer   | Editor        | /remote.php/dav/spaces/<<FILEID>> |
      | Secure viewer   | Uploader      | /remote.php/dav/spaces/<<FILEID>> |
      | Secure viewer   | Secure viewer | /remote.php/dav/spaces/<<FILEID>> |
      | Viewer          | Viewer        | /remote.php/dav/spaces/<<FILEID>> |
      | Editor          | Viewer        | /remote.php/dav/spaces/<<FILEID>> |
      | Uploader        | Viewer        | /remote.php/dav/spaces/<<FILEID>> |
      | Viewer          | Secure viewer | /remote.php/dav/spaces/<<FILEID>> |
      | Editor          | Secure viewer | /remote.php/dav/spaces/<<FILEID>> |
      | Uploader        | Secure viewer | /remote.php/dav/spaces/<<FILEID>> |
      | Secure viewer   | Viewer        | /dav/spaces/<<FILEID>>            |
      | Secure viewer   | Editor        | /dav/spaces/<<FILEID>>            |
      | Secure viewer   | Uploader      | /dav/spaces/<<FILEID>>            |
      | Secure viewer   | Secure viewer | /dav/spaces/<<FILEID>>            |
      | Viewer          | Viewer        | /dav/spaces/<<FILEID>>            |
      | Editor          | Viewer        | /dav/spaces/<<FILEID>>            |
      | Uploader        | Viewer        | /dav/spaces/<<FILEID>>            |
      | Viewer          | Secure viewer | /dav/spaces/<<FILEID>>            |
      | Editor          | Secure viewer | /dav/spaces/<<FILEID>>            |
      | Uploader        | Secure viewer | /dav/spaces/<<FILEID>>            |
