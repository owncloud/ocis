@api @files_sharing-app-required @skipOnOcis-OCIS-Storage
Feature: sharing

  Background:
    Given these users have been created with default attributes and skeleton files:
      | username |
      | Alice    |
      | Brian    |

  @issue-ocis-reva-260
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario Outline: Sharee can't see the share that is filtered out
    Given using OCS API version "<ocs_api_version>"
    And user "Alice" has shared file "textfile0.txt" with user "Brian"
    And user "Alice" has shared file "textfile1.txt" with user "Brian"
    When user "Brian" gets all the shares shared with him that are received as file "textfile0 (2).txt" using the provisioning API
    Then the OCS status code should be "<ocs_status_code>"
    And the HTTP status code should be "200"
    And the last share_id should be included in the response
    Examples:
      | ocs_api_version | ocs_status_code |
      | 1               | 100             |
      | 2               | 200             |
