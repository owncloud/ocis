@issue-1328 @skipOnReva
Feature: a subfolder of a received share can be reshared
  As a user
  I want to re-share a resource
  So that other users can have access to it

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
      | Carol    |
    And user "Alice" has created folder "/TMP"
    And user "Alice" has created folder "/TMP/SUB"

  @smokeTest @issue-2214
  Scenario Outline: user is allowed to reshare a sub-folder with the same permissions
    Given using OCS API version "<ocs-api-version>"
    And user "Alice" has shared folder "/TMP" with user "Brian" with permissions "share,read"
    When user "Brian" shares folder "/Shares/TMP/SUB" with user "Carol" with permissions "share,read" using the sharing API
    Then the OCS status code should be "<ocs-status-code>"
    And the HTTP status code should be "200"
    And as "Carol" folder "/Shares/SUB" should exist
    And as "Brian" folder "/Shares/TMP/SUB" should exist
    Examples:
      | ocs-api-version | ocs-status-code |
      | 1               | 100             |
      | 2               | 200             |

  @skipOnRevaMaster
  Scenario Outline: user is not allowed to reshare a sub-folder with more permissions
    Given using OCS API version "<ocs-api-version>"
    And user "Alice" has shared folder "/TMP" with user "Brian" with permissions <received-permissions>
    When user "Brian" shares folder "/Shares/TMP/SUB" with user "Carol" with permissions <reshare-permissions> using the sharing API
    Then the OCS status code should be "403"
    And the HTTP status code should be "<http-status-code>"
    And as "Carol" folder "/Shares/SUB" should not exist
    And the sharing API should report to user "Carol" that no shares are in the pending state
    And as "Brian" folder "/Shares/TMP/SUB" should exist
    Examples:
      | ocs-api-version | http-status-code | received-permissions | reshare-permissions |
      # try to pass on more bits including reshare
      | 1               | 200              | 17                   | 19                  |
      | 2               | 403              | 17                   | 19                  |
      | 1               | 200              | 17                   | 21                  |
      | 2               | 403              | 17                   | 21                  |
      | 1               | 200              | 17                   | 23                  |
      | 2               | 403              | 17                   | 23                  |
      | 1               | 200              | 17                   | 31                  |
      | 2               | 403              | 17                   | 31                  |
      | 1               | 200              | 19                   | 23                  |
      | 2               | 403              | 19                   | 23                  |
      | 1               | 200              | 19                   | 31                  |
      | 2               | 403              | 19                   | 31                  |
      # try to pass on more bits but not reshare
      | 1               | 200              | 17                   | 3                   |
      | 2               | 403              | 17                   | 3                   |
      | 1               | 200              | 17                   | 5                   |
      | 2               | 403              | 17                   | 5                   |
      | 1               | 200              | 17                   | 7                   |
      | 2               | 403              | 17                   | 7                   |
      | 1               | 200              | 17                   | 15                  |
      | 2               | 403              | 17                   | 15                  |
      | 1               | 200              | 19                   | 7                   |
      | 2               | 403              | 19                   | 7                   |
      | 1               | 200              | 19                   | 15                  |
      | 2               | 403              | 19                   | 15                  |
      # try to pass on extra delete (including reshare)
      | 1               | 200              | 17                   | 25                  |
      | 2               | 403              | 17                   | 25                  |
      | 1               | 200              | 19                   | 27                  |
      | 2               | 403              | 19                   | 27                  |
      | 1               | 200              | 23                   | 31                  |
      | 2               | 403              | 23                   | 31                  |
      # try to pass on extra delete (but not reshare)
      | 1               | 200              | 17                   | 9                   |
      | 2               | 403              | 17                   | 9                   |
      | 1               | 200              | 19                   | 11                  |
      | 2               | 403              | 19                   | 11                  |
      | 1               | 200              | 23                   | 15                  |
      | 2               | 403              | 23                   | 15                  |

  @issue-2214
  Scenario Outline: user is allowed to update reshare of a sub-folder with less permissions
    Given using OCS API version "<ocs-api-version>"
    And user "Alice" has shared folder "/TMP" with user "Brian" with permissions "share,create,update,read"
    And user "Brian" has shared folder "/Shares/TMP/SUB" with user "Carol" with permissions "share,create,update,read"
    When user "Brian" updates the last share using the sharing API with
      | permissions | share,read |
    Then the OCS status code should be "<ocs-status-code>"
    And the HTTP status code should be "200"
    And as "Carol" folder "/Shares/SUB" should exist
    But user "Carol" should not be able to upload file "filesForUpload/textfile.txt" to "/Shares/SUB/textfile.txt"
    And as "Brian" folder "/Shares/TMP/SUB" should exist
    And user "Brian" should be able to upload file "filesForUpload/textfile.txt" to "/Shares/TMP/SUB/textfile.txt"
    Examples:
      | ocs-api-version | ocs-status-code |
      | 1               | 100             |
      | 2               | 200             |

  @issue-2214
  Scenario Outline: user is allowed to update reshare of a sub-folder to the maximum allowed permissions
    Given using OCS API version "<ocs-api-version>"
    And user "Alice" has shared folder "/TMP" with user "Brian" with permissions "share,create,update,read"
    And user "Brian" has shared folder "/Shares/TMP/SUB" with user "Carol" with permissions "share,read"
    When user "Brian" updates the last share using the sharing API with
      | permissions | share,create,update,read |
    Then the OCS status code should be "<ocs-status-code>"
    And the HTTP status code should be "200"
    And as "Carol" folder "/Shares/SUB" should exist
    And user "Carol" should be able to upload file "filesForUpload/textfile.txt" to "/Shares/SUB/textfile.txt"
    And as "Brian" folder "/Shares/TMP/SUB" should exist
    And user "Brian" should be able to upload file "filesForUpload/textfile.txt" to "/Shares/TMP/SUB/textfile.txt"
    Examples:
      | ocs-api-version | ocs-status-code |
      | 1               | 100             |
      | 2               | 200             |

  @issue-2214 @skipOnRevaMaster
  Scenario Outline: user is not allowed to update reshare of a sub-folder with more permissions
    Given using OCS API version "<ocs-api-version>"
    And user "Alice" has shared folder "/TMP" with user "Brian" with permissions "share,read"
    And user "Brian" has shared folder "/Shares/TMP/SUB" with user "Carol" with permissions "share,read"
    When user "Brian" updates the last share using the sharing API with
      | permissions | all |
    Then the OCS status code should be "403"
    And the HTTP status code should be "<http-status-code>"
    And as "Carol" folder "/Shares/SUB" should exist
    But user "Carol" should not be able to upload file "filesForUpload/textfile.txt" to "/Shares/SUB/textfile.txt"
    And as "Brian" folder "/Shares/TMP/SUB" should exist
    But user "Brian" should not be able to upload file "filesForUpload/textfile.txt" to "/Shares/TMP/SUB/textfile.txt"
    Examples:
      | ocs-api-version | http-status-code |
      | 1               | 200              |
      | 2               | 403              |