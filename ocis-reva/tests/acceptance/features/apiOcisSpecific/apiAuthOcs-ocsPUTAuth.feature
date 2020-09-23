@api
Feature: auth

  @issue-ocis-reva-30
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario: send PUT request to OCS endpoints as admin with wrong password
    When the administrator requests these endpoints with "PUT" with body "doesnotmatter" using password "invalid" about user "Alice"
      | endpoint                                         |
      | /ocs/v1.php/cloud/users/%username%               |
      | /ocs/v2.php/cloud/users/%username%               |
      | /ocs/v1.php/cloud/users/%username%/disable       |
      | /ocs/v2.php/cloud/users/%username%/disable       |
      | /ocs/v1.php/cloud/users/%username%/enable        |
      | /ocs/v2.php/cloud/users/%username%/enable        |
      | /ocs/v1.php/apps/files_sharing/api/v1/shares/123 |
      | /ocs/v2.php/apps/files_sharing/api/v1/shares/123 |
    Then the HTTP status code of responses on all endpoints should be "401"
    And the OCS status code of responses on all endpoints should be "notset"
