@api @skipOnOcV10
Feature:

  Background:
    Given using spaces DAV path
    And these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |

  @issue-4421
  Scenario: Move files between shares by different users
    Given user "Carol" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file with content "some data" to "/textfile0.txt"
    And user "Alice" has created folder "/PARENT"
    And user "Brian" has created folder "/PARENT"
    And user "Alice" has moved file "textfile0.txt" to "PARENT/from_alice.txt" in space "Personal"
    And user "Alice" has shared folder "/PARENT" with user "Carol"
    And user "Brian" has shared folder "/PARENT" with user "Carol"
    And user "Carol" has accepted share "/PARENT" offered by user "Alice"
    And user "Carol" has accepted share "/PARENT" offered by user "Brian"
    When user "Carol" moves file "PARENT/from_alice.txt" to "PARENT (1)/from_alice.txt" in space "Shares Jail" using the WebDAV API
    Then the HTTP status code should be "201"
    And for user "Carol" folder "PARENT" of the space "Shares Jail" should not contain these entries:
      | from_alice.txt |
    And for user "Carol" folder "PARENT (1)" of the space "Shares Jail" should contain these entries:
      | from_alice.txt |


  Scenario: overwrite a received file share
    Given user "Alice" has uploaded file with content "this is the old content" to "/textfile1.txt"
    And user "Alice" has shared file "/textfile1.txt" with user "Brian"
    And user "Brian" has accepted share "/textfile1.txt" offered by user "Alice"
    When user "Brian" uploads a file inside space "Shares Jail" with content "this is a new content" to "textfile1.txt" using the WebDAV API
    Then the HTTP status code should be "204"
    And for user "Brian" the space "Shares Jail" should contain these entries:
      | textfile1.txt |
    And for user "Brian" the content of the file "/textfile1.txt" of the space "Shares Jail" should be "this is a new content"
    And for user "Alice" the content of the file "/textfile1.txt" of the space "Personal" should be "this is a new content"
