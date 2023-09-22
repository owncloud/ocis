Feature: Propfind by file/folder using file id
  As a user
  I want to check the PROPFIND response of file/folder using their file id
  So that I can make sure that the response contains all the relevant values

  Background:
    Given using spaces DAV path
    And user "Alice" has been created with default attributes and without skeleton files


  Scenario Outline: send a PROPFIND requests to a file with its FILEID in dav-path url inside root of personal space
    Given user "Alice" has uploaded file with content "some data" to "/textfile.txt"
    And we save it into "FILEID"
    When user "Alice" sends HTTP method "PROPFIND" to URL "<dav-path>"
    Then the HTTP status code should be "207"
    And the "PROPFIND" response to user "Alice" should contain a mountpoint "Alice Hansen" with these key and value pairs:
      | key            | value               |
      | oc:name        | textfile.txt        |
      | oc:permissions | RDNVWZP             |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: send a PROPFIND requests to a file with its FILEID in dav-path url inside a folder of personal space
    Given user "Alice" has created folder "folder"
    And user "Alice" has uploaded file with content "some data" to "folder/textfile.txt"
    And we save it into "FILEID"
    When user "Alice" sends HTTP method "PROPFIND" to URL "<dav-path>"
    Then the HTTP status code should be "207"
    And the "PROPFIND" response to user "Alice" should contain a mountpoint "Alice Hansen" with these key and value pairs:
      | key            | value               |
      | oc:name        | textfile.txt        |
      | oc:permissions | RDNVWZP             |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: sends a PROPFIND request a file in personal space with its FILEID in dav-path url of another user
    Given user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file with content "some data" to "textfile.txt"
    And we save it into "FILEID"
    When user "Brian" sends HTTP method "PROPFIND" to URL "<dav-path>"
    Then the HTTP status code should be "404"
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: send a PROPFIND requests to a file with its FILEID in dav-path url inside project space
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the GraphApi
    And user "Alice" has uploaded a file inside space "new-space" with content "some content" to "textfile.txt"
    And we save it into "FILEID"
    When user "Alice" sends HTTP method "PROPFIND" to URL "<dav-path>"
    Then the HTTP status code should be "207"
    And the "PROPFIND" response should contain a space "new-space" with these key and value pairs:
      | key            | value        |
      | oc:name        | textfile.txt |
      | oc:permissions | RDNVWZP     |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: send a PROPFIND requests to a file with its FILEID in dav-path url inside a folder of project space
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the GraphApi
    And user "Alice" has created a folder "folder" in space "new-space"
    And user "Alice" has uploaded a file inside space "new-space" with content "some content" to "/folder/textfile.txt"
    And we save it into "FILEID"
    When user "Alice" sends HTTP method "PROPFIND" to URL "<dav-path>"
    Then the HTTP status code should be "207"
    And the "PROPFIND" response should contain a space "new-space" with these key and value pairs:
      | key            | value        |
      | oc:name        | textfile.txt |
      | oc:permissions | RDNVWZP     |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: sends a PROPFIND request a file in project space with its FILEID in dav-path url of another user
    Given user "Brian" has been created with default attributes and without skeleton files
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the GraphApi
    And user "Alice" has uploaded a file inside space "new-space" with content "some content" to "textfile.txt"
    And we save it into "FILEID"
    When user "Brian" sends HTTP method "PROPFIND" to URL "<dav-path>"
    Then the HTTP status code should be "404"
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: send a PROPFIND requests to a file with its FILEID in dav-path of a shared file
    Given user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file with content "some data" to "/textfile.txt"
    And we save it into "FILEID"
    And user "Alice" has shared file "/textfile.txt" with user "Brian"
    And user "Brian" has accepted share "/textfile.txt" offered by user "Alice"
    When user "Brian" sends HTTP method "PROPFIND" to URL "<dav-path>"
    Then the HTTP status code should be "207"
    And the "PROPFIND" response to user "Alice" should contain a mountpoint "Brian Murphy" with these key and value pairs:
      | key            | value        |
      | oc:name        | textfile.txt |
      | oc:permissions | SRNVW        |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: sharee send a PROPFIND requests to a file with its FILEID in dav-path inside of a shared folder
    Given user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has created folder "/folder"
    And user "Alice" has shared folder "/folder" with user "Brian"
    And user "Brian" has accepted share "/folder" offered by user "Alice"
    And user "Alice" has uploaded file with content "some data" to "/folder/textfile.txt"
    And we save it into "FILEID"
    When user "Brian" sends HTTP method "PROPFIND" to URL "<dav-path>"
    Then the HTTP status code should be "207"
    And the "PROPFIND" response to user "Alice" should contain a mountpoint "Brian Murphy" with these key and value pairs:
      | key            | value        |
      | oc:name        | textfile.txt |
      | oc:permissions | RDNVWZP      |
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |
