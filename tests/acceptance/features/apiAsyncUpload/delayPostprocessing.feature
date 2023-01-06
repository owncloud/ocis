@api
Feature: delay post-processing of uploaded files

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files

  @issue-5326
  Scenario Outline: user sends GET request to the file while it's still being processed
    Given user "Alice" has uploaded file with content "uploaded content" to "/file.txt"
    When user "Alice" requests "<dav_path>" with "GET" without retrying
    Then the HTTP status code should be "425"
    Examples:
      | dav_path                                  |
      | /remote.php/webdav/file.txt               |
      | /remote.php/dav/files/%username%/file.txt |
      | /dav/spaces/%spaceid%/file.txt            |


  Scenario Outline: user sends PROPFIND request to the file while it's still being processed
    Given user "Alice" has uploaded file with content "uploaded content" to "/file.txt"
    When user "Alice" requests "<dav_path>" with "PROPFIND" without retrying
    Then the HTTP status code should be "207"
    And the value of the item "//d:response/d:propstat/d:status" in the response should be "HTTP/1.1 425 TOO EARLY"
    Examples:
      | dav_path                                  |
      | /remote.php/webdav/file.txt               |
      | /remote.php/dav/files/%username%/file.txt |
      | /dav/spaces/%spaceid%/file.txt            |


  Scenario Outline: user sends PROPFIND request to the folder while files in the folder are still being processed
    Given user "Alice" has created folder "my_data"
    And user "Alice" has uploaded file with content "uploaded content" to "/my_data/file.txt"
    When user "Alice" requests "<dav_path>" with "PROPFIND" without retrying
    Then the HTTP status code should be "207"
    And as user "Alice" the value of the item "//d:status" of path "<dav_path>/" in the response should be "HTTP/1.1 200 OK"
    And as user "Alice" the value of the item "//d:status" of path "<dav_path>/file.txt" in the response should be "HTTP/1.1 425 TOO EARLY"
    Examples:
      | dav_path                                 |
      | /remote.php/webdav/my_data               |
      | /remote.php/dav/files/%username%/my_data |
      | /dav/spaces/%spaceid%/my_data            |