Feature: cannot share resources with invalid permissions
  As a user
  I want to share resources with invalid permission
  So that I can make sure it doesn't work

  Background:
    Given user "Alice" has been created with default attributes
    And user "Alice" has uploaded file with content "some data" to "/textfile0.txt"
    And user "Alice" has created folder "/PARENT"


  Scenario Outline: cannot create a share of a file or folder with invalid permissions
    Given using OCS API version "<ocs-api-version>"
    And user "Brian" has been created with default attributes
    When user "Alice" creates a share using the sharing API with settings
      | path        | <resource>    |
      | shareWith   | Brian         |
      | shareType   | user          |
      | permissions | <permissions> |
    Then the OCS status code should be "<ocs-status-code>"
    And the HTTP status code should be "<http-status-code>"
    And as "Brian" entry "<resource>" should not exist
    And as "Brian" entry "/Shares/<resource>" should not exist
    And the sharing API should report to user "Brian" that no shares are in the pending state
    Examples:
      | ocs-api-version | ocs-status-code | http-status-code | resource      | permissions |
      | 1               | 400             | 200              | textfile0.txt | 0           |
      | 2               | 400             | 400              | textfile0.txt | 0           |
      | 1               | 400             | 200              | PARENT        | 0           |
      | 2               | 400             | 400              | PARENT        | 0           |
      | 1               | 404             | 200              | textfile0.txt | 32          |
      | 2               | 404             | 404              | textfile0.txt | 32          |
      | 1               | 404             | 200              | PARENT        | 32          |
      | 2               | 404             | 404              | PARENT        | 32          |


  Scenario Outline: cannot create a share of a file with a user with only create permission
    Given using OCS API version "<ocs-api-version>"
    And user "Brian" has been created with default attributes
    When user "Alice" creates a share using the sharing API with settings
      | path        | textfile0.txt |
      | shareWith   | Brian         |
      | shareType   | user          |
      | permissions | create        |
    Then the OCS status code should be "400"
    And the HTTP status code should be "<http-status-code>"
    And as "Brian" entry "textfile0.txt" should not exist
    And as "Brian" entry "/Shares/textfile0.txt" should not exist
    And the sharing API should report to user "Brian" that no shares are in the pending state
    Examples:
      | ocs-api-version | http-status-code |
      | 1               | 200              |
      | 2               | 400              |


  Scenario Outline: cannot create a share of a file with a user with only (create,delete) permission
    Given using OCS API version "<ocs-api-version>"
    And user "Brian" has been created with default attributes
    When user "Alice" creates a share using the sharing API with settings
      | path        | textfile0.txt |
      | shareWith   | Brian         |
      | shareType   | user          |
      | permissions | <permissions> |
    Then the OCS status code should be "400"
    And the HTTP status code should be "<http-status-code>"
    And as "Brian" entry "textfile0.txt" should not exist
    And as "Brian" entry "/Shares/textfile0.txt" should not exist
    And the sharing API should report to user "Brian" that no shares are in the pending state
    Examples:
      | ocs-api-version | http-status-code | permissions   |
      | 1               | 200              | delete        |
      | 2               | 400              | delete        |
      | 1               | 200              | create,delete |
      | 2               | 400              | create,delete |


  Scenario Outline: cannot create a share of a file with a group with only create permission
    Given using OCS API version "<ocs-api-version>"
    And user "Brian" has been created with default attributes
    And group "grp1" has been created
    And user "Brian" has been added to group "grp1"
    When user "Alice" creates a share using the sharing API with settings
      | path        | textfile0.txt |
      | shareWith   | grp1          |
      | shareType   | group         |
      | permissions | create        |
    Then the OCS status code should be "400"
    And the HTTP status code should be "<http-status-code>"
    And as "Brian" entry "textfile0.txt" should not exist
    And as "Brian" entry "/Shares/textfile0.txt" should not exist
    And the sharing API should report to user "Brian" that no shares are in the pending state
    Examples:
      | ocs-api-version | http-status-code |
      | 1               | 200              |
      | 2               | 400              |


  Scenario Outline: cannot create a share of a file with a group with only (create,delete) permission
    Given using OCS API version "<ocs-api-version>"
    And user "Brian" has been created with default attributes
    And group "grp1" has been created
    And user "Brian" has been added to group "grp1"
    When user "Alice" creates a share using the sharing API with settings
      | path        | textfile0.txt |
      | shareWith   | grp1          |
      | shareType   | group         |
      | permissions | <permissions> |
    Then the OCS status code should be "400"
    And the HTTP status code should be "<http-status-code>"
    And as "Brian" entry "textfile0.txt" should not exist
    And as "Brian" entry "/Shares/textfile0.txt" should not exist
    And the sharing API should report to user "Brian" that no shares are in the pending state
    Examples:
      | ocs-api-version | http-status-code | permissions   |
      | 1               | 200              | delete        |
      | 2               | 400              | delete        |
      | 1               | 200              | create,delete |
      | 2               | 400              | create,delete |
