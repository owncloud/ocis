Feature: accessing files using file id
  As a user
  I want to access the files using file id
  So that I can get the content of a file

  Background:
    Given using spaces DAV path
    And user "Alice" has been created with default attributes and without skeleton files


  Scenario Outline: get content of a file
    Given user "Alice" has uploaded file with content "some data" to "/textfile.txt"
    And we save it into "FILEID"
    When user "Alice" sends HTTP method "GET" to URL "<dav-path>"
    Then the HTTP status code should be "200"
    And the downloaded content should be "some data"
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: get content of a file inside a folder
    Given user "Alice" has created folder "uploadFolder"
    And user "Alice" has uploaded file with content "some data" to "uploadFolder/textfile.txt"
    And we save it into "FILEID"
    When user "Alice" sends HTTP method "GET" to URL "<dav-path>"
    Then the HTTP status code should be "200"
    And the downloaded content should be "some data"
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: get content of a file inside a project space
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "new-space" with content "some data" to "textfile.txt"
    And we save it into "FILEID"
    When user "Alice" sends HTTP method "GET" to URL "<dav-path>"
    Then the HTTP status code should be "200"
    And the downloaded content should be "some data"
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: sharee gets content of a shared file
    Given user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file with content "some data" to "/textfile.txt"
    And we save it into "FILEID"
    And user "Alice" has shared file "/textfile.txt" with user "Brian"
    When user "Brian" sends HTTP method "GET" to URL "<dav-path>"
    Then the HTTP status code should be "200"
    And the downloaded content should be "some data"
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: sharee gets content of a file inside a shared folder
    Given user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has created folder "uploadFolder"
    And user "Alice" has uploaded file with content "some data" to "uploadFolder/textfile.txt"
    And we save it into "FILEID"
    And user "Alice" has shared folder "/uploadFolder" with user "Brian"
    When user "Brian" sends HTTP method "GET" to URL "<dav-path>"
    Then the HTTP status code should be "200"
    And the downloaded content should be "some data"
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: sharee gets content of a file inside a shared space
    Given user "Brian" has been created with default attributes and without skeleton files
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "new-space" with content "some data" to "textfile.txt"
    And we save it into "FILEID"
    And user "Alice" has shared a space "new-space" with settings:
      | shareWith | Brian  |
      | role      | viewer |
    When user "Brian" sends HTTP method "GET" to URL "<dav-path>"
    Then the HTTP status code should be "200"
    And the downloaded content should be "some data"
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: user tries to get content of file owned by others
    Given user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file with content "some data" to "/textfile.txt"
    And we save it into "FILEID"
    When user "Brian" sends HTTP method "GET" to URL "<dav-path>"
    Then the HTTP status code should be "404"
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |
