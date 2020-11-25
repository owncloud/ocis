@api @files_sharing-app-required @issue-ocis-reva-243
Feature: cannot share resources with invalid permissions

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file with content "some data" to "/textfile0.txt"
    And user "Alice" has created folder "/PARENT"

  @issue-ocis-reva-45 @issue-ocis-reva-243
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario Outline: Cannot create a share of a file with a user with only create permission
    Given using OCS API version "<ocs_api_version>"
    And user "Brian" has been created with default attributes and without skeleton files
    When user "Alice" creates a share using the sharing API with settings
      | path        | textfile0.txt |
      | shareWith   | Brian         |
      | shareType   | user          |
      | permissions | create        |
    Then the OCS status code should be "<ocs_status_code>" or "<eos_status_code>"
    And the HTTP status code should be "<http_status_code_ocs>" or "<http_status_code_eos>"
    And as "Brian" entry "textfile0.txt" should not exist
    Examples:
      | ocs_api_version | ocs_status_code | eos_status_code | http_status_code_ocs | http_status_code_eos |
      | 1               | 100             | 996             | 200                  | 500                  |
      | 2               | 200             | 996             | 200                  | 500                  |

  @issue-ocis-reva-45 @issue-ocis-reva-243
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario Outline: Cannot create a share of a file with a user with only (create,delete) permission
    Given using OCS API version "<ocs_api_version>"
    And user "Brian" has been created with default attributes and without skeleton files
    When user "Alice" creates a share using the sharing API with settings
      | path        | textfile0.txt |
      | shareWith   | Brian         |
      | shareType   | user          |
      | permissions | <permissions> |
    Then the OCS status code should be "<ocs_status_code>" or "<eos_status_code>"
    And the HTTP status code should be "<http_status_code_ocs>" or "<http_status_code_eos>"
    And as "Brian" entry "textfile0.txt" should not exist
    Examples:
      | ocs_api_version | eos_status_code | ocs_status_code | http_status_code_ocs | http_status_code_eos | permissions   |
      | 1               | 100             | 996             | 200                  | 500                  | delete        |
      | 2               | 200             | 996             | 200                  | 500                  | delete        |
      | 1               | 100             | 996             | 200                  | 500                  | create,delete |
      | 2               | 200             | 996             | 200                  | 500                  | create,delete |
