@skipOnReva
Feature: sharing
  As a user
  I want to get all the shares
  So that I can know I have proper access to them

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |

  @smokeTest @issue-1258
  Scenario Outline: getting all shares from a user
    Given using OCS API version "<ocs_api_version>"
    And user "Alice" has uploaded file with content "some data" to "/file_to_share.txt"
    And user "Alice" has shared file "file_to_share.txt" with user "Brian"
    When user "Alice" gets all shares shared by her using the sharing API
    Then the OCS status code should be "<ocs_status_code>"
    And the HTTP status code should be "200"
    And file "/Shares/file_to_share.txt" should be included in the response
    Examples:
      | ocs_api_version | ocs_status_code |
      | 1               | 100             |
      | 2               | 200             |

  @issue-1319
  Scenario Outline: getting all shares of a user using another user
    Given using OCS API version "<ocs_api_version>"
    And user "Alice" has uploaded file with content "some data" to "/textfile0.txt"
    And user "Alice" has shared file "textfile0.txt" with user "Brian"
    When the administrator gets all shares shared by him using the sharing API
    Then the OCS status code should be "<ocs_status_code>"
    And the HTTP status code should be "200"
    And file "/Shares/textfile0.txt" should not be included in the response
    Examples:
      | ocs_api_version | ocs_status_code |
      | 1               | 100             |
      | 2               | 200             |

  @smokeTest
  Scenario Outline: getting all shares of a file
    Given using OCS API version "<ocs_api_version>"
    And these users have been created with default attributes and without skeleton files:
      | username |
      | Carol    |
      | David    |
    And user "Alice" has uploaded file with content "some data" to "/textfile0.txt"
    And user "Alice" has shared file "textfile0.txt" with user "Brian"
    And user "Alice" has shared file "textfile0.txt" with user "Carol"
    When user "Alice" gets all the shares of the file "textfile0.txt" using the sharing API
    Then the OCS status code should be "<ocs_status_code>"
    And the HTTP status code should be "200"
    And user "Brian" should be included in the response
    And user "Carol" should be included in the response
    And user "David" should not be included in the response
    Examples:
      | ocs_api_version | ocs_status_code |
      | 1               | 100             |
      | 2               | 200             |

  @smokeTest
  Scenario Outline: getting all shares of a file with reshares
    Given using OCS API version "<ocs_api_version>"
    And these users have been created with default attributes and without skeleton files:
      | username |
      | Carol    |
      | David    |
    And user "Alice" has uploaded file with content "some data" to "/textfile0.txt"
    And user "Alice" has shared file "textfile0.txt" with user "Brian"
    And user "Brian" has shared file "/Shares/textfile0.txt" with user "Carol"
    When user "Alice" gets all the shares with reshares of the file "textfile0.txt" using the sharing API
    Then the OCS status code should be "<ocs_status_code>"
    And the HTTP status code should be "200"
    And user "Brian" should be included in the response
    And user "Carol" should be included in the response
    And user "David" should not be included in the response
    Examples:
      | ocs_api_version | ocs_status_code |
      | 1               | 100             |
      | 2               | 200             |

  @smokeTest
  Scenario Outline: resource can be reshared to resource owner
    Given using OCS API version "<ocs_api_version>"
    And group "grp1" has been created
    And user "Carol" has been created with default attributes and without skeleton files
    And user "Carol" has been added to group "grp1"
    And user "Carol" has created folder "/shared"
    And user "Carol" has uploaded file with content "some data" to "/shared/shared_file.txt"
    And user "Carol" has shared folder "/shared" with user "Brian"
    And user "Brian" has shared folder "/Shares/shared" with group "grp1"
    # no need to accept this share as it is Carol's file
    When user "Carol" gets all the shares shared with her using the sharing API
    Then the OCS status code should be "<ocs_status_code>"
    And the HTTP status code should be "200"
    And file "/Shares/shared" should be included in the response
    Examples:
      | ocs_api_version | ocs_status_code |
      | 1               | 100             |
      | 2               | 200             |

  @smokeTest @issue-1226 @issue-1270 @issue-1271
  Scenario Outline: getting share info of a share
    Given using OCS API version "<ocs_api_version>"
    And user "Alice" has uploaded file with content "some data" to "/file_to_share.txt"
    And user "Alice" has shared file "file_to_share.txt" with user "Brian"
    When user "Alice" gets the info of the last share using the sharing API
    Then the OCS status code should be "<ocs_status_code>"
    And the HTTP status code should be "200"
    And the fields of the last response to user "Alice" sharing with user "Brian" should include
      | id                     | A_STRING                  |
      | item_type              | file                      |
      | item_source            | A_STRING                  |
      | share_type             | user                      |
      | share_with             | %username%                |
      | file_source            | A_STRING                  |
      | file_target            | /Shares/file_to_share.txt |
      | path                   | /file_to_share.txt        |
      | permissions            | share,read,update         |
      | stime                  | A_NUMBER                  |
      | storage                | A_STRING                  |
      | mail_send              | 0                         |
      | uid_owner              | %username%                |
      | share_with_displayname | %displayname%             |
      | displayname_owner      | %displayname%             |
      | mimetype               | text/plain                |
    Examples:
      | ocs_api_version | ocs_status_code |
      | 1               | 100             |
      | 2               | 200             |

  @issue-1233
  Scenario Outline: get a share with a user that didn't receive the share
    Given using OCS API version "<ocs_api_version>"
    And user "Carol" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file with content "some data" to "/textfile0.txt"
    And user "Alice" has shared file "textfile0.txt" with user "Brian"
    When user "Carol" gets the info of the last share using the sharing API
    Then the OCS status code should be "404"
    And the HTTP status code should be "<http_status_code>"
    Examples:
      | ocs_api_version | http_status_code |
      | 1               | 200              |
      | 2               | 404              |

  @issue-1289
  Scenario: share a folder to a group, and remove user from that group
    Given using OCS API version "1"
    And user "Carol" has been created with default attributes and without skeleton files
    And group "group0" has been created
    And user "Brian" has been added to group "group0"
    And user "Carol" has been added to group "group0"
    And user "Alice" has created folder "/PARENT"
    And user "Alice" has uploaded file with content "some data" to "/PARENT/parent.txt"
    And user "Alice" has shared folder "/PARENT" with group "group0"
    When the administrator removes user "Carol" from group "group0" using the provisioning API
    Then the HTTP status code should be "204"
    And user "Brian" should see the following elements
      | /Shares/PARENT/           |
      | /Shares/PARENT/parent.txt |
    But user "Carol" should not see the following elements
      | /Shares/PARENT/           |
      | /Shares/PARENT/parent.txt |

  @issue-1231
  Scenario Outline: getting all the shares inside the folder
    Given using OCS API version "<ocs_api_version>"
    And user "Alice" has created folder "/PARENT"
    And user "Alice" has uploaded file with content "some data" to "/PARENT/parent.txt"
    And user "Alice" has shared file "PARENT/parent.txt" with user "Brian"
    When user "Alice" gets all the shares inside the folder "PARENT" using the sharing API
    Then the OCS status code should be "<ocs_status_code>"
    And the HTTP status code should be "200"
    And file "/Shares/parent.txt" should be included in the response
    Examples:
      | ocs_api_version | ocs_status_code |
      | 1               | 100             |
      | 2               | 200             |
