@issue-1328
Feature: resharing a resource with an expiration date
  As a user
  I want to reshare resources with expiration date
  So that other users will have access to the resources only for a limited amount of time

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "/textfile0.txt"
    And user "Carol" has been created with default attributes and without skeleton files


  Scenario Outline: user should be able to set expiration while resharing a file with user
    Given using OCS API version "<ocs_api_version>"
    And user "Alice" has shared file "/textfile0.txt" with user "Brian" with permissions "read,update,share"
    And user "Brian" has accepted share "/textfile0.txt" offered by user "Alice"
    When user "Brian" creates a share using the sharing API with settings
      | path        | /Shares/textfile0.txt |
      | shareType   | user                  |
      | permissions | change                |
      | shareWith   | Carol                 |
      | expireDate  | +3 days               |
    Then the HTTP status code should be "200"
    And the OCS status code should be "<ocs_status_code>"
    And user "Carol" should be able to accept pending share "/textfile0.txt" offered by user "Brian"
    And the information of the last share of user "Brian" should include
      | expiration | +3 days |
    And the response when user "Carol" gets the info of the last share should include
      | expiration | +3 days |
    Examples:
      | ocs_api_version | ocs_status_code |
      | 1               | 100             |
      | 2               | 200             |

  @issue-1289
  Scenario Outline: user should be able to set expiration while resharing a file with group
    Given using OCS API version "<ocs_api_version>"
    And group "grp1" has been created
    And user "Carol" has been added to group "grp1"
    And user "Alice" has shared file "/textfile0.txt" with user "Brian" with permissions "read,update,share"
    And user "Brian" has accepted share "/textfile0.txt" offered by user "Alice"
    When user "Brian" creates a share using the sharing API with settings
      | path        | /Shares/textfile0.txt |
      | shareType   | group                 |
      | permissions | change                |
      | shareWith   | grp1                  |
      | expireDate  | +3 days               |
    Then the HTTP status code should be "200"
    And the OCS status code should be "<ocs_status_code>"
    And user "Carol" should be able to accept pending share "/textfile0.txt" offered by user "Brian"
    And the information of the last share of user "Brian" should include
      | expiration | +3 days |
    And the response when user "Carol" gets the info of the last share should include
      | expiration | +3 days |
    Examples:
      | ocs_api_version | ocs_status_code |
      | 1               | 100             |
      | 2               | 200             |


  Scenario Outline: resharing using the sharing API with default expire date set but not enforced
    Given using OCS API version "<ocs_api_version>"
    And user "Alice" has shared file "/textfile0.txt" with user "Brian" with permissions "read,update,share"
    And user "Brian" has accepted share "/textfile0.txt" offered by user "Alice"
    When user "Brian" creates a share using the sharing API with settings
      | path        | /Shares/textfile0.txt |
      | shareType   | user                  |
      | permissions | change                |
      | shareWith   | Carol                 |
    Then the HTTP status code should be "200"
    And the OCS status code should be "<ocs_status_code>"
    And user "Carol" should be able to accept pending share "/textfile0.txt" offered by user "Brian"
    And the information of the last share of user "Brian" should include
      | expiration |  |
    And the response when user "Carol" gets the info of the last share should include
      | expiration |  |
    Examples:
      | ocs_api_version | ocs_status_code |
      | 1               | 100             |
      | 2               | 200             |
