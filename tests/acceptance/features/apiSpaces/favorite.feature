@api @skipOnOcV10
Feature: favorite

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
    And using spaces DAV path
    And user "Alice" has created folder "/PARENT"


  Scenario: favorite a received share itself
    Given user "Alice" has shared folder "/PARENT" with user "Brian"
    And user "Brian" has accepted share "/PARENT" offered by user "Alice"
    When user "Brian" favorites element "/PARENT" in space "Shares Jail" using the WebDAV API
    Then the HTTP status code should be "207"
    And as user "Brian" folder "/PARENT" inside space "Shares Jail" should contain a property "oc:favorite" with value "1"


  Scenario: favorite a file inside of a received share
    Given user "Alice" has uploaded file with content "some data" to "/PARENT/parent.txt"
    And user "Alice" has shared folder "/PARENT" with user "Brian"
    And user "Brian" has accepted share "/PARENT" offered by user "Alice"
    When user "Brian" favorites element "/PARENT/parent.txt" in space "Shares Jail" using the WebDAV API
    Then the HTTP status code should be "207"
    And as user "Brian" file "/PARENT/parent.txt" inside space "Shares Jail" should contain a property "oc:favorite" with value "1"


  Scenario: favorite a folder inside of a received share
    Given user "Alice" has created folder "/PARENT/sub-folder"
    And user "Alice" has shared folder "/PARENT" with user "Brian"
    And user "Brian" has accepted share "/PARENT" offered by user "Alice"
    When user "Brian" favorites element "/PARENT/sub-folder" in space "Shares Jail" using the WebDAV API
    Then the HTTP status code should be "207"
    And as user "Brian" folder "/PARENT/sub-folder" inside space "Shares Jail" should contain a property "oc:favorite" with value "1"


  Scenario: sharee file favorite state should not change the favorite state of sharer
    Given user "Alice" has uploaded file with content "some data" to "/PARENT/parent.txt"
    And user "Alice" has shared file "/PARENT/parent.txt" with user "Brian"
    And user "Brian" has accepted share "/parent.txt" offered by user "Alice"
    When user "Brian" favorites element "/parent.txt" in space "Shares Jail" using the WebDAV API
    Then the HTTP status code should be "207"
    And as user "Brian" file "/parent.txt" inside space "Shares Jail" should contain a property "oc:favorite" with value "1"
    And as user "Alice" file "/PARENT/parent.txt" inside space "Personal" should contain a property "oc:favorite" with value "0"
