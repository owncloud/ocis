@api @files_sharing-app-required
Feature: sharing

  Background:
    Given these users have been created with default attributes and skeleton files:
      | username |
      | Alice    |
      | Brian    |

  @issue-ocis-reva-374
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario Outline: Get a share with a user that didn't receive the share
    Given using OCS API version "<ocs_api_version>"
    And user "Carol" has been created with default attributes and without skeleton files
    And user "Alice" has shared file "textfile0.txt" with user "Brian"
    When user "Carol" gets the info of the last share using the sharing API
    Then the OCS status code should be "998"
#    Then the OCS status code should be "404"
    And the HTTP status code should be "<http_status_code>"
    Examples:
      | ocs_api_version | http_status_code |
      | 1               | 200              |
      | 2               | 404              |

  @issue-ocis-reva-372
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario Outline: getting all the shares inside the folder
    Given using OCS API version "<ocs_api_version>"
    And user "Alice" has shared file "PARENT/parent.txt" with user "Brian"
    When user "Alice" gets all the shares inside the folder "PARENT/parent.txt" using the sharing API
    Then the OCS status code should be "<ocs_status_code>"
    And the HTTP status code should be "200"
    And file "parent.txt" should be included in the response
    Examples:
      | ocs_api_version | ocs_status_code |
      | 1               | 100             |
      | 2               | 200             |
