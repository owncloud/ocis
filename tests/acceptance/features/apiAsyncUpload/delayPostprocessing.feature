@api
Feature: delay post-processing of uploaded files

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file with content "uploaded content" to "/file.txt"

  @issue-5326
  Scenario Outline: user sends GET request to the file while it's still being processed
    When user "Alice" requests "<dav_path>" with "GET" without retrying
    Then the HTTP status code should be "425"
    Examples:
      | dav_path                                  |
      | /remote.php/webdav/file.txt               |
      | /remote.php/dav/files/%username%/file.txt |
      | /dav/spaces/%spaceid%/file.txt            |


  Scenario Outline: user sends PROPFIND request to the file while it's still being processed
    When user "Alice" requests "<dav_path>" with "PROPFIND" without retrying
    Then the HTTP status code should be "207"
    And the value of the item "//d:response/d:propstat/d:status" in the response should be "HTTP/1.1 425 TOO EARLY"
    Examples:
      | dav_path                                  |
      | /remote.php/webdav/file.txt               |
      | /remote.php/dav/files/%username%/file.txt |
      | /dav/spaces/%spaceid%/file.txt            |