@api @skipOnOcV10
Feature: TUS upload resources using TUS protocol with mtime
  As a user
  I want the mtime of an uploaded file to be the creation date on upload source not the upload date
  So that I can find files by their real creation date

  Background:
    Given using spaces DAV path
    And these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
    And user "Alice" has created folder "/toShare"
    And user "Alice" has shared folder "/toShare" with user "Brian"
    And user "Brian" has accepted share "/toShare" offered by user "Alice"


  Scenario: upload file with mtime to a received share
    When user "Brian" uploads a file "filesForUpload/textfile.txt" to "toShare/file.txt" with mtime "Thu, 08 Aug 2012 04:18:13 GMT" via TUS inside of the space "Shares Jail" using the WebDAV API
    Then the HTTP status code should be "200"
    And for user "Brian" folder "toShare" of the space "Shares Jail" should contain these entries:
      | file.txt |
    And as "Brian" the mtime of the file "/toShare/file.txt" in space "Shares Jail" should be "Thu, 08 Aug 2012 04:18:13 GMT"
    And as "Alice" the mtime of the file "/toShare/file.txt" in space "Personal" should be "Thu, 08 Aug 2012 04:18:13 GMT"


  Scenario: upload file with mtime to a send share
    When user "Alice" uploads a file "filesForUpload/textfile.txt" to "toShare/file.txt" with mtime "Thu, 08 Aug 2012 04:18:13 GMT" via TUS inside of the space "Personal" using the WebDAV API
    Then the HTTP status code should be "200"
    And for user "Alice" folder "toShare" of the space "Personal" should contain these entries:
      | file.txt |
    And as "Alice" the mtime of the file "/toShare/file.txt" in space "Personal" should be "Thu, 08 Aug 2012 04:18:13 GMT"
    And as "Brian" the mtime of the file "/toShare/file.txt" in space "Shares Jail" should be "Thu, 08 Aug 2012 04:18:13 GMT"


  Scenario: overwriting a file with mtime in a received share
    Given user "Alice" has uploaded file with content "uploaded content" to "/toShare/file.txt"
    When user "Brian" uploads a file "filesForUpload/textfile.txt" to "toShare/file.txt" with mtime "Thu, 08 Aug 2012 04:18:13 GMT" via TUS inside of the space "Shares Jail" using the WebDAV API
    Then the HTTP status code should be "200"
    And for user "Brian" folder "toShare" of the space "Shares Jail" should contain these entries:
      | file.txt |
    And as "Brian" the mtime of the file "/toShare/file.txt" in space "Shares Jail" should be "Thu, 08 Aug 2012 04:18:13 GMT"
    And as "Alice" the mtime of the file "/toShare/file.txt" in space "Personal" should be "Thu, 08 Aug 2012 04:18:13 GMT"


  Scenario: overwriting a file with mtime in a send share
    Given user "Brian" has uploaded a file inside space "Shares Jail" with content "uploaded content" to "toShare/file.txt"
    When user "Alice" uploads a file "filesForUpload/textfile.txt" to "toShare/file.txt" with mtime "Thu, 08 Aug 2012 04:18:13 GMT" via TUS inside of the space "Personal" using the WebDAV API
    Then the HTTP status code should be "200"
    And for user "Alice" folder "toShare" of the space "Personal" should contain these entries:
      | file.txt |
    And as "Alice" the mtime of the file "/toShare/file.txt" in space "Personal" should be "Thu, 08 Aug 2012 04:18:13 GMT"
    And as "Brian" the mtime of the file "/toShare/file.txt" in space "Shares Jail" should be "Thu, 08 Aug 2012 04:18:13 GMT"

