@api @files_sharing-app-required
Feature: sharing

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file with content "ownCloud test text file 0" to "/textfile0.txt"

  @skipOnOcis-OC-Storage @skipOnOcis-OCIS-Storage @issue-ocis-reva-301 @issue-ocis-reva-302
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario Outline: Creating a share of a file with a user and asking for various permission combinations
    Given using OCS API version "<ocs_api_version>"
    And user "Brian" has been created with default attributes and without skeleton files
    When user "Alice" shares file "textfile0.txt" with user "Brian" with permissions <requested_permissions> using the sharing API
    Then the OCS status code should be "<ocs_status_code>"
    And the HTTP status code should be "200"
    And the fields of the last response to user "Alice" sharing with user "Brian" should include
      | share_with        | %username%               |
      | file_target       | /textfile0.txt           |
      | path              | /textfile0.txt           |
      | permissions       | <granted_permissions>    |
      | uid_owner         | %username%               |
      | displayname_owner |                          |
      | item_type         | file                     |
      | mimetype          | application/octet-stream |
      | storage_id        | ANY_VALUE                |
      | share_type        | user                     |
    And the fields of the last response should not include
      | share_with_displayname | %displayname% |
    Examples:
      | ocs_api_version | requested_permissions | granted_permissions | ocs_status_code |
      # Ask for full permissions. You get share plus read plus update. create and delete do not apply to shares of a file
      | 1               | 31                    | 19                  | 100             |
      | 2               | 31                    | 19                  | 200             |
      # Ask for read, share (17), create and delete. You get share plus read
      | 1               | 29                    | 17                  | 100             |
      | 2               | 29                    | 17                  | 200             |
      # Ask for read, update, create, delete. You get read plus update.
      | 1               | 15                    | 3                   | 100             |
      | 2               | 15                    | 3                   | 200             |
      # Ask for just update. You get exactly update (you do not get read or anything else)
      | 1               | 2                     | 2                   | 100             |
      | 2               | 2                     | 2                   | 200             |

  @issue-ocis-reva-243
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario Outline: more tests to demonstrate different ocis-reva issue 243 behaviours
    Given using OCS API version "<ocs_api_version>"
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has created folder "/home"
    And user "Alice" has uploaded file with content "Random data" to "/home/randomfile.txt"
    When user "Alice" shares file "/home/randomfile.txt" with user "Brian" using the sharing API
    And the HTTP status code should be "<http_status_code_ocs>" or "<http_status_code_eos>"
    And as "Brian" file "randomfile.txt" should not exist
    Examples:
      | ocs_api_version | http_status_code_ocs | http_status_code_eos |
      | 1               | 200                  | 500                  |
      | 2               | 200                  | 500                  |

  @skipOnOcis-OC-Storage @skipOnOcis-OCIS-Storage @issue-ocis-reva-301 @issue-ocis-reva-302
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario Outline: Creating a share of a folder with a user, the default permissions are all permissions(31)
    Given using OCS API version "<ocs_api_version>"
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has created folder "/FOLDER"
    When user "Alice" shares folder "/FOLDER" with user "Brian" using the sharing API
    Then the OCS status code should be "<ocs_status_code>"
    And the HTTP status code should be "200"
    And the fields of the last response to user "Alice" sharing with user "Brian" should include
      | share_with        | %username%           |
      | file_target       | /FOLDER              |
      | path              | /FOLDER              |
      | permissions       | all                  |
      | uid_owner         | %username%           |
      | displayname_owner |                      |
      | item_type         | folder               |
      | mimetype          | httpd/unix-directory |
      | storage_id        | ANY_VALUE            |
      | share_type        | user                 |
    And the fields of the last response should not include
      | share_with_displayname | %displayname% |
    Examples:
      | ocs_api_version | ocs_status_code |
      | 1               | 100             |
      | 2               | 200             |

  @issue-ocis-reva-372 @issue-ocis-reva-243 @skipOnOcis-OCIS-Storage
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario Outline: sharing subfolder of already shared folder, GET result is correct
    Given using OCS API version "<ocs_api_version>"
    And these users have been created with default attributes and without skeleton files:
      | username |
      | Brian    |
      | Carol    |
      | David    |
      | Emily    |
    And user "Alice" has created folder "/folder1"
    And user "Alice" has shared folder "/folder1" with user "Brian"
    And user "Alice" has shared folder "/folder1" with user "Carol"
    And user "Alice" has created folder "/folder1/folder2"
    And user "Alice" has shared folder "/folder1/folder2" with user "David"
    And user "Alice" has shared folder "/folder1/folder2" with user "Emily"
    When user "Alice" sends HTTP method "GET" to OCS API endpoint "/apps/files_sharing/api/v1/shares"
    Then the OCS status code should be "<ocs_status_code>"
    And the HTTP status code should be "200"
    And the response should contain 4 entries
    And folder "/folder1" should be included as path in the response
    # And folder "/folder1/folder2" should be included as path in the response
    And folder "/folder2" should be included as path in the response
    And user "Alice" sends HTTP method "GET" to OCS API endpoint "/apps/files_sharing/api/v1/shares?path=/folder1/folder2"
    And the response should contain 2 entries
    And folder "/folder1" should not be included as path in the response
    # And folder "/folder1/folder2" should be included as path in the response
    And folder "/folder2" should be included as path in the response
    Examples:
      | ocs_api_version | ocs_status_code |
      | 1               | 100             |
      | 2               | 200             |
