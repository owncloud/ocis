Feature: checking file versions using file id
  As a user
  I want to share file outside of the space
  So that other users can access the file

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
    And using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "Project1" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "Project1" with content "hello world version 1" to "text.txt"
    And we save it into "FILEID"
    And user "Alice" has uploaded a file inside space "Project1" with content "hello world version 1.1" to "text.txt"


  Scenario Outline: check the file versions of a file shared from project space
    Given user "Alice" has sent the following resource share invitation:
      | resource        | text.txt |
      | space           | Project1 |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | <role>   |
    And user "Brian" has a share "text.txt" synced
    When user "Alice" gets the number of versions of file "/text.txt" using file-id path "/meta/<<FILEID>>/v"
    Then the HTTP status code should be "207"
    And the number of versions should be "1"
    When user "Brian" tries to get the number of versions of file "/text.txt" using file-id path "/meta/<<FILEID>>/v"
    Then the HTTP status code should be "403"
    Examples:
      | role        |
      | File Editor |
      | Viewer      |


  Scenario Outline: check the versions of a file in a shared space as editor/manager
    Given user "Alice" has sent the following space share invitation:
      | space           | Project1     |
      | sharee          | Brian        |
      | shareType       | user         |
      | permissionsRole | <space-role> |
    When user "Alice" gets the number of versions of file "/text.txt" using file-id path "/meta/<<FILEID>>/v"
    Then the HTTP status code should be "207"
    And the number of versions should be "1"
    When user "Brian" gets the number of versions of file "/text.txt" using file-id path "/meta/<<FILEID>>/v"
    Then the HTTP status code should be "207"
    And the number of versions should be "1"
    Examples:
      | space-role   |
      | Space Editor |
      | Manager      |


  Scenario: check the versions of a file in a shared space as viewer
    Given user "Alice" has sent the following space share invitation:
      | space           | Project1     |
      | sharee          | Brian        |
      | shareType       | user         |
      | permissionsRole | Space Viewer |
    When user "Brian" tries to get the number of versions of file "/text.txt" using file-id path "/meta/<<FILEID>>/v"
    Then the HTTP status code should be "403"

  @issue-7738
  Scenario Outline: check the versions of a file after moving to a shared folder inside a project space as editor/viewer
    Given user "Alice" has created a folder "testFolder" in space "Project1"
    And user "Alice" has sent the following resource share invitation:
      | resource        | text.txt |
      | space           | Project1 |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | <role>   |
    And user "Brian" has a share "text.txt" synced
    And user "Alice" has moved file "text.txt" to "/testFolder/movedText.txt" in space "Project1"
    When user "Alice" gets the number of versions of file "/testFolder/movedText.txt" using file-id path "/meta/<<FILEID>>/v"
    Then the HTTP status code should be "207"
    And the number of versions should be "1"
    When user "Brian" tries to get the number of versions of file "/Shares/testFolder/movedText.txt" using file-id path "/meta/<<FILEID>>/v"
    Then the HTTP status code should be "403"
    Examples:
      | role        |
      | File Editor |
      | Viewer      |

  @issue-7738
  Scenario: check the versions of a file after moving it to a shared folder inside a project space as manager
    Given user "Alice" has created a folder "testFolder" in space "Project1"
    And user "Alice" has sent the following resource share invitation:
      | resource        | testFolder |
      | space           | Project1   |
      | sharee          | Brian      |
      | shareType       | user       |
      | permissionsRole | Viewer     |
    And user "Brian" has a share "testFolder" synced
    And user "Alice" has sent the following space share invitation:
      | space           | Project1 |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Manager  |
    And user "Alice" has moved file "text.txt" to "/testFolder/movedText.txt" in space "Project1"
    When user "Brian" gets the number of versions of file "/text.txt" using file-id path "/meta/<<FILEID>>/v"
    Then the HTTP status code should be "207"
    And the number of versions should be "1"

  @issue-777
  Scenario Outline: check file versions after moving to-and-from folder in personal space
    Given user "Alice" has created folder "folder"
    And user "Alice" has uploaded file with content "some data" to "<source>textfile.txt"
    And user "Alice" has uploaded file with content "some data - edited" to "<source>textfile.txt"
    And we save it into "FILEID"
    And user "Alice" has moved file "<source>textfile.txt" to "<destination>textfile.txt" in space "Personal"
    When user "Alice" gets the number of versions of file "<destination>textfile.txt" using file-id path "/meta/<<FILEID>>/v"
    Then the HTTP status code should be "207"
    And the number of versions should be "1"
    Examples:
      | source  | destination |
      | /       | folder/     |
      | folder/ | /           |

  @issue-777
  Scenario Outline: check file versions after moving to-and-from folder in personal space (MOVE using file-id)
    Given user "Alice" has created folder "folder"
    And user "Alice" has uploaded file with content "some data" to "<source>textfile.txt"
    And user "Alice" has uploaded file with content "some data - edited" to "<source>textfile.txt"
    And we save it into "FILEID"
    When user "Alice" moves a file "<source>textfile.txt" into "<destination>" inside space "Personal" using file-id path "/dav/spaces/<<FILEID>>"
    Then the HTTP status code should be "201"
    When user "Alice" gets the number of versions of file "<destination>textfile.txt" using file-id path "/meta/<<FILEID>>/v"
    Then the HTTP status code should be "207"
    And the number of versions should be "1"
    Examples:
      | source  | destination |
      | /       | folder/     |
      | folder/ | /           |

  @issue-777
  Scenario Outline: check file versions after moving to-and-from folder in project space
    Given user "Alice" has created a folder "folder" in space "Project1"
    And user "Alice" has uploaded a file inside space "Project1" with content "some data" to "<source>textfile.txt"
    And user "Alice" has uploaded a file inside space "Project1" with content "some data - edited" to "<source>textfile.txt"
    And we save it into "FILEID"
    And user "Alice" has moved file "<source>textfile.txt" to "<destination>textfile.txt" in space "Project1"
    When user "Alice" gets the number of versions of file "<source>textfile.txt" using file-id path "/meta/<<FILEID>>/v"
    Then the HTTP status code should be "207"
    And the number of versions should be "1"
    Examples:
      | source  | destination |
      | /       | folder/     |
      | folder/ | /           |

  @issue-777
  Scenario Outline: check file versions after moving to-and-from folder in project space (MOVE using file-id)
    And user "Alice" has created a folder "folder" in space "Project1"
    And user "Alice" has uploaded a file inside space "Project1" with content "some data" to "<source>textfile.txt"
    And user "Alice" has uploaded a file inside space "Project1" with content "some data - edited" to "<source>textfile.txt"
    And we save it into "FILEID"
    When user "Alice" moves a file "<source>textfile.txt" into "<destination>" inside space "Project1" using file-id path "/dav/spaces/<<FILEID>>"
    Then the HTTP status code should be "201"
    When user "Alice" gets the number of versions of file "<destination>textfile.txt" using file-id path "/meta/<<FILEID>>/v"
    Then the HTTP status code should be "207"
    And the number of versions should be "1"
    Examples:
      | source  | destination |
      | /       | folder/     |
      | folder/ | /           |
