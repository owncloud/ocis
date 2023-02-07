@api @files_sharing-app-required @issue-ocis-1328 @issue-ocis-1250
Feature: a default expiration date can be specified for shares with users or groups

  Background:
    Given the administrator has set the default folder for received shares to "Shares"
    And auto-accept shares has been disabled
    And these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |


  Scenario Outline: sharing with default expiration date enabled but not enforced for users, user shares without specifying expireDate
    Given using OCS API version "<ocs_api_version>"
    And parameter "shareapi_default_expire_date_user_share" of app "core" has been set to "yes"
    And user "Alice" has created folder "/FOLDER"
    When user "Alice" shares folder "/FOLDER" with user "Brian" using the sharing API
    And user "Brian" accepts share "/FOLDER" offered by user "Alice" using the sharing API
    Then the OCS status code of responses on all endpoints should be "<ocs_status_code>"
    And the HTTP status code of responses on all endpoints should be "<http_status_code>"
    And the fields of the last response to user "Alice" should include
      | expiration |  |
    And the response when user "Brian" gets the info of the last share should include
      | expiration |  |
    Examples:
      | ocs_api_version | ocs_status_code | http_status_code |
      | 1               | 100             | 200              |
      | 2               | 200             | 200              |


  Scenario Outline: sharing with default expiration date enabled but not enforced for groups, user shares without specifying expireDate
    Given using OCS API version "<ocs_api_version>"
    And parameter "shareapi_default_expire_date_group_share" of app "core" has been set to "yes"
    And group "grp1" has been created
    And user "Brian" has been added to group "grp1"
    And user "Alice" has created folder "/FOLDER"
    And user "Alice" has shared folder "/FOLDER" with group "grp1"
    When user "Brian" accepts share "/FOLDER" offered by user "Alice" using the sharing API
    Then the OCS status code should be "<ocs_status_code>"
    And the HTTP status code should be "<http_status_code>"
    And the fields of the last response to user "Alice" should include
      | expiration |  |
    And the response when user "Brian" gets the info of the last share should include
      | expiration |  |
    Examples:
      | ocs_api_version | ocs_status_code | http_status_code |
      | 1               | 100             | 200              |
      | 2               | 200             | 200              |


  Scenario Outline: sharing with default expiration date not enabled for groups, user shares with expiration date set
    Given using OCS API version "<ocs_api_version>"
    And group "grp1" has been created
    And user "Brian" has been added to group "grp1"
    And user "Alice" has created folder "/FOLDER"
    And user "Alice" has created a share with settings
      | path        | /FOLDER    |
      | shareType   | group      |
      | shareWith   | grp1       |
      | permissions | read,share |
      | expireDate  | +15 days   |
    When user "Brian" accepts share "/FOLDER" offered by user "Alice" using the sharing API
    Then the info about the last share by user "Alice" with user "Brian" should include
      | share_type  | group          |
      | file_target | /Shares/FOLDER |
      | uid_owner   | %username%     |
      | expiration  | +15 days       |
      | share_with  | grp1           |
    And the response when user "Brian" gets the info of the last share should include
      | expiration | +15 days |
    Examples:
      | ocs_api_version |
      | 1               |
      | 2               |


  Scenario Outline: sharing with default expiration date enforced for users, user shares to a group without setting an expiration date
    Given using OCS API version "<ocs_api_version>"
    And parameter "shareapi_default_expire_date_user_share" of app "core" has been set to "yes"
    And parameter "shareapi_enforce_expire_date_user_share" of app "core" has been set to "yes"
    And group "grp1" has been created
    And user "Brian" has been added to group "grp1"
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has shared folder "FOLDER" with group "grp1" with permissions "read,share"
    When user "Brian" accepts share "/FOLDER" offered by user "Alice" using the sharing API
    Then the info about the last share by user "Alice" with user "Brian" should include
      | expiration |  |
    And the response when user "Brian" gets the info of the last share should include
      | expiration |  |
    Examples:
      | ocs_api_version |
      | 1               |
      | 2               |


  Scenario Outline: sharing with default expiration date enforced for groups, user shares to another user
    Given using OCS API version "<ocs_api_version>"
    And parameter "shareapi_default_expire_date_group_share" of app "core" has been set to "yes"
    And parameter "shareapi_enforce_expire_date_group_share" of app "core" has been set to "yes"
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has shared folder "/FOLDER" with user "Brian" with permissions "read,share"
    When user "Brian" accepts share "/FOLDER" offered by user "Alice" using the sharing API
    Then the info about the last share by user "Alice" with user "Brian" should include
      | expiration |  |
    And the response when user "Brian" gets the info of the last share should include
      | expiration |  |
    Examples:
      | ocs_api_version |
      | 1               |
      | 2               |
