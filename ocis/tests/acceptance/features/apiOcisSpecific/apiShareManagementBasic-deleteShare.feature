@api @files_sharing-app-required @issue-ocis-reva-243
Feature: sharing

  @issue-ocis-reva-356
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario Outline: delete a share
    Given user "Alice" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file with content "ownCloud test text file 0" to "/textfile0.txt"
    And user "Brian" has been created with default attributes and without skeleton files
    And using OCS API version "<ocs_api_version>"
    And user "Alice" has shared file "textfile0.txt" with user "Brian"
    When user "Alice" deletes the last share using the sharing API
    Then the OCS status code should be "<ocs_status_code>"
    And the HTTP status code should be "200"
    Examples:
      | ocs_api_version | ocs_status_code |
      | 1               | 100             |
      | 2               | 200             |
