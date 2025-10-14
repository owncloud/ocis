Feature: refuse access
  As an administrator
  I want to refuse access to unauthenticated and disabled users
  So that I can secure the system

  Background:
    Given using OCS API version "1"

  @smokeTest @issue-2285
  Scenario Outline: unauthenticated call
    # cannot perform with spaces WebDAV due to the absence of user
    Given using <dav-path-version> DAV path
    When an unauthenticated client connects to the DAV endpoint using the WebDAV API
    Then the HTTP status code should be "401"
    And there should be no duplicate headers
    And the following headers should be set
      | header           | value                                        |
      | WWW-Authenticate | Basic realm="%productname%", charset="UTF-8" |
    Examples:
      | dav-path-version |
      | old              |
      | new              |

  @issue-2285
  Scenario Outline: disabled user cannot use webdav
    Given using <dav-path-version> DAV path
    And user "Alice" has been created with default attributes
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "textfile0.txt"
    And user "Alice" has been disabled
    When user "Alice" downloads file "/textfile0.txt" using the WebDAV API
    Then the HTTP status code should be "401"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |
