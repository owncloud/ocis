@api @issue-ocis-187
Feature: previews of files downloaded through the webdav API

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files

  @issue-ocis-2069
  Scenario Outline: download different sizes of previews of file on ocis
    Given user "Alice" has uploaded file "filesForUpload/lorem.txt" to "/parent.txt"
    When user "Alice" downloads the preview of "/parent.txt" with width <width> and height <height> using the WebDAV API
    Then the HTTP status code should be "200"
    And the downloaded image should be <expected-width> pixels wide and <expected-height> pixels high
    Examples:
      | width | height | expected-width | expected-height |
      | 1     | 1      | 16             | 16              |
      | 32    | 32     | 32             | 32              |
      | 1024  | 1024   | 640            | 480             |
      | 1     | 1024   | 16             | 16              |
      | 1024  | 1      | 640            | 480             |

  @issue-ocis-2071
  Scenario: download previews of other users files in ocis
    Given user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file "filesForUpload/lorem.txt" to "/parent.txt"
    When user "Brian" downloads the preview of "/parent.txt" of "Alice" with width "32" and height "32" using the WebDAV API
    Then the HTTP status code should be "404"
    And the value of the item "/d:error/s:message" in the response about user "Alice" should be "File with name parent.txt could not be located"
    And the value of the item "/d:error/s:exception" in the response about user "Alice" should be "Sabre\DAV\Exception\NotFound"

  @issue-ocis-2070
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario: Download file previews when it is disabled by the administrator
    Given the administrator has updated system config key "enable_previews" with value "false" and type "boolean"
    And user "Alice" has uploaded file "filesForUpload/lorem.txt" to "/parent.txt"
    When user "Alice" downloads the preview of "/parent.txt" with width "32" and height "32" using the WebDAV API
    Then the HTTP status code should be "200"

  @issue-ocis-2070
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario: unset maximum size of previews
    Given user "Alice" has uploaded file "filesForUpload/lorem.txt" to "/parent.txt"
    And the administrator has updated system config key "preview_max_x" with value "null"
    And the administrator has updated system config key "preview_max_y" with value "null"
    When user "Alice" downloads the preview of "/parent.txt" with width "32" and height "32" using the WebDAV API
    Then the HTTP status code should be "200"

  @issue-ocis-2070
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario: set maximum size of previews
    Given user "Alice" has uploaded file "filesForUpload/lorem.txt" to "/parent.txt"
    When the administrator updates system config key "preview_max_x" with value "null" using the occ command
    And the administrator updates system config key "preview_max_y" with value "null" using the occ command
    Then the HTTP status code should be "201"
    When user "Alice" downloads the preview of "/parent.txt" with width "null" and height "null" using the WebDAV API
    Then the HTTP status code should be "400"
