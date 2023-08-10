@api
Feature: move test
  As a user
  I want to check the MOVE request
  So that I can make sure that the resource moved successfully


  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
    And using spaces DAV path

    Scenario Outline: move file using fileid
      Given  user "Alice" has uploaded file "filesForUpload/lorem.txt" to "/parent.txt"
      And user "ALice" renames file "/parent.txt" to "/rename.txt" using MOVE request with dav-path "<dav-path>"
      Then the HTTP status code should be "<status-code>"
      And as "Alice" file "/parent.txt" should not exist
      And as "Alice" file "/rename.txt" should exist
      Examples:
        | dav-path                | status-code |
        | /remote.php/dav/spaces/ | 201         |
        | /dav/spaces/            | 400         |


  Scenario Outline: rename file with same name
    Given  user "Alice" has uploaded file "filesForUpload/lorem.txt" to "/parent.txt"
    And user "ALice" renames file "/parent.txt" to "/parent.txt" using MOVE request with dav-path "<dav-path>"
    Then the HTTP status code should be "409"
    Examples:
      | dav-path                        |
      | /remote.php/dav/spaces/dav-path |
      | /remote.php/dav/spaces/         |


  Scenario Outline: move file using fileid (inside subfolder with different name)
    Given user "Alice" has created folder "/folderMain"
    And  user "Alice" has uploaded file "filesForUpload/lorem.txt" to "/parent.txt"
    And user "ALice" renames file "/parent.txt" to "/folderMain/rename.txt" using MOVE request with dav-path "<dav-path>"
    Then the HTTP status code should be "<status-code>"
    And as "Alice" file "/parent.txt" should not exist
    And as "Alice" file "/folderMain/parent.txt" should exist
    Examples:
      | dav-path                | status-code |
      | /remote.php/dav/spaces/ | 201         |
      | /dav/spaces/            | 201         |


  Scenario Outline: move file using fileid (inside subfolder with different name)
    Given user "Alice" has created folder "/folderMain"
    And  user "Alice" has uploaded file "filesForUpload/lorem.txt" to "/parent.txt"
    And user "ALice" renames file "/parent.txt" to "/folderMain/parent.txt" using MOVE request with dav-path "<dav-path>"
    Then the HTTP status code should be "<status-code>"
    And as "Alice" file "/parent.txt" should not exist
    And as "Alice" file "/folderMain/parent.txt" should exist
    Examples:
      | dav-path                | status-code |
      | /remote.php/dav/spaces/ | 201         |
      | /dav/spaces/            | 201         |


  Scenario Outline: move file using other user fileid
    Given user "Brian" has created folder "/folderMain"
    And user "Alice" has uploaded file with content "123" to "/davtest.txt"
    And we save it into "FILEID"
    And user "Brian" makes HTTP request "MOVE" file "<dav-path>/<<FILEID>>" to "<dav-path>/%spaceid%/folderMain/parent.txt"
    Then the HTTP status code should be "<status-code>"
    And as "Alice" file "/davtest.txt" should exist
    And as "Brian" file "/folderMain/davtest.txt" should not exist
    Examples:
      | dav-path                | status-code |
      | /remote.php/dav/spaces/ | 404         |
      | /dav/spaces/            | 404         |


    Scenario Outline: send move request to shared file
      Given user "Brian" has created folder "/folderMain"
      And user "Alice" has uploaded file "filesForUpload/lorem.txt" to "/parent.txt"
      And we save it into "FILEID"
      And user "Alice" has shared entry "/parent.txt" with user "Brian"
      And user "Brian" makes HTTP request "MOVE" file "<dav-path>/<<FILEID>>" to "<dav-path>/%spaceid%/folderMain/parent.txt"
      Then the HTTP status code should be "<status-code>"
      Examples:
        | dav-path                | status-code |
        | /remote.php/dav/spaces/ | 404         |
        | /dav/spaces/            | 404         |
