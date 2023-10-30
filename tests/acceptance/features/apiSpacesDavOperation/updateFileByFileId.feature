Feature: update files using file id
  As a user
  I want to update the files using file id
  So that I can make changes on the content of a file

  Background:
    Given using spaces DAV path
    And user "Alice" has been created with default attributes and without skeleton files


  Scenario Outline: update content of a file
    Given user "Alice" has uploaded file with content "some data" to "/textfile.txt"
    And we save it into "FILEID"
    When user "Alice" sends HTTP method "PUT" to URL "<dav-path>" with content "updated content"
    Then the HTTP status code should be "204"
    And for user "Alice" the content of the file "/textfile.txt" of the space "Personal" should be "updated content"
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: update content of a file inside a folder
    Given user "Alice" has created folder "uploadFolder"
    And user "Alice" has uploaded file with content "some data" to "uploadFolder/textfile.txt"
    And we save it into "FILEID"
    When user "Alice" sends HTTP method "PUT" to URL "<dav-path>" with content "updated content"
    Then the HTTP status code should be "204"
    And for user "Alice" the content of the file "/uploadFolder/textfile.txt" of the space "Personal" should be "updated content"
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: update content of a file inside a project space
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "new-space" with content "some data" to "/textfile.txt"
    And we save it into "FILEID"
    When user "Alice" sends HTTP method "PUT" to URL "<dav-path>" with content "updated content"
    Then the HTTP status code should be "204"
    And for user "Alice" the content of the file "/textfile.txt" of the space "new-space" should be "updated content"
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: sharee updates content of a shared file
    Given user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file with content "some data" to "/textfile.txt"
    And we save it into "FILEID"
    And user "Alice" has shared file "/textfile.txt" with user "Brian"
    When user "Brian" sends HTTP method "PUT" to URL "<dav-path>" with content "updated content"
    Then the HTTP status code should be "204"
    And for user "Alice" the content of the file "/textfile.txt" of the space "Personal" should be "updated content"
    And for user "Brian" the content of the file "textfile.txt" of the space "Shares" should be "updated content"
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: sharee updates content of a file inside a shared folder
    Given user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has created folder "uploadFolder"
    And user "Alice" has uploaded file with content "some data" to "uploadFolder/textfile.txt"
    And we save it into "FILEID"
    And user "Alice" has shared folder "/uploadFolder" with user "Brian"
    When user "Brian" sends HTTP method "PUT" to URL "<dav-path>" with content "updated content"
    Then the HTTP status code should be "204"
    And for user "Alice" the content of the file "uploadFolder/textfile.txt" of the space "Personal" should be "updated content"
    And for user "Brian" the content of the file "uploadFolder/textfile.txt" of the space "Shares" should be "updated content"
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: sharee with different role tries to update content of a file inside a shared space
    Given user "Brian" has been created with default attributes and without skeleton files
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "new-space" with content "some data" to "/textfile.txt"
    And we save it into "FILEID"
    And user "Alice" has shared a space "new-space" with settings:
      | shareWith | Brian  |
      | role      | <role> |
    When user "Brian" sends HTTP method "PUT" to URL "<dav-path>" with content "updated content"
    Then the HTTP status code should be "<http_status_code>"
    And for user "Alice" the content of the file "/textfile.txt" of the space "new-space" should be "<file-content>"
    And for user "Brian" the content of the file "/textfile.txt" of the space "new-space" should be "<file-content>"
    Examples:
      | dav-path                          | role   | http_status_code | file-content    |
      | /remote.php/dav/spaces/<<FILEID>> | viewer | 403              | some data       |
      | /dav/spaces/<<FILEID>>            | editor | 204              | updated content |


  Scenario Outline: user tries to update content of a file owned by others
    Given user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file with content "some data" to "/textfile.txt"
    And we save it into "FILEID"
    When user "Brian" sends HTTP method "PUT" to URL "<dav-path>" with content "updated content"
    Then the HTTP status code should be "404"
    And for user "Alice" the content of the file "/textfile.txt" of the space "Personal" should be "some data"
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spcaes/<<FILEID>>            |
