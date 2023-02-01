@api @issue-ocis-reva-172
Feature: actions on a locked item are possible if the token is sent with the request

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files

  @files_sharing-app-required
  Scenario Outline: two users having both a shared lock can use the resource
    Given using <dav-path> DAV path
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file with content "some data" to "textfile0.txt"
    And user "Brian" has uploaded file with content "some data" to "textfile0.txt"
    And user "Alice" has shared file "/textfile0.txt" with user "Brian"
    And user "Brian" has accepted share "/textfile0.txt" offered by user "Alice"
    And user "Alice" has locked file "textfile0.txt" setting the following properties
      | lockscope | shared |
    And user "Brian" has locked file "Shares/textfile0.txt" setting the following properties
      | lockscope | shared |
    When user "Alice" uploads file with content "from user 0" to "textfile0.txt" sending the locktoken of file "textfile0.txt" using the WebDAV API
    Then the HTTP status code should be "204"
    And the content of file "textfile0.txt" for user "Alice" should be "from user 0"
    And the content of file "Shares/textfile0.txt" for user "Brian" should be "from user 0"
    When user "Brian" uploads file with content "from user 1" to "Shares/textfile0.txt" sending the locktoken of file "Shares/textfile0.txt" using the WebDAV API
    Then the HTTP status code should be "204"
    And the content of file "textfile0.txt" for user "Alice" should be "from user 1"
    And the content of file "Shares/textfile0.txt" for user "Brian" should be "from user 1"
    Examples:
      | dav-path |
      | old      |
      | new      |

    @personalSpace
    Examples:
      | dav-path |
      | spaces   |
