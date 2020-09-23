@api
Feature: auth

  @issue-ocis-reva-30 @issue-ocis-reva-65
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario: send DELETE requests to OCS endpoints as admin with wrong password
    When the administrator requests these endpoints with "DELETE" using password "invalid" about user "Alice"
      | endpoint                                                        |
      | /ocs/v1.php/apps/files_sharing/api/v1/remote_shares/pending/123 |
      | /ocs/v2.php/apps/files_sharing/api/v1/remote_shares/pending/123 |
      | /ocs/v1.php/apps/files_sharing/api/v1/remote_shares/123         |
      | /ocs/v2.php/apps/files_sharing/api/v1/remote_shares/123         |
      | /ocs/v2.php/apps/files_sharing/api/v1/shares/123                |
      | /ocs/v1.php/apps/files_sharing/api/v1/shares/pending/123        |
      | /ocs/v2.php/apps/files_sharing/api/v1/shares/pending/123        |
      | /ocs/v1.php/cloud/apps/testing                                  |
      | /ocs/v2.php/cloud/apps/testing                                  |
      | /ocs/v1.php/cloud/groups/group1                                 |
      | /ocs/v2.php/cloud/groups/group1                                 |
      | /ocs/v1.php/cloud/users/%username%                              |
      | /ocs/v2.php/cloud/users/%username%                              |
      | /ocs/v1.php/cloud/users/%username%/groups                       |
      | /ocs/v2.php/cloud/users/%username%/groups                       |
      | /ocs/v1.php/cloud/users/%username%/subadmins                    |
      | /ocs/v2.php/cloud/users/%username%/subadmins                    |
    Then the HTTP status code of responses on all endpoints should be "401"
    And the OCS status code of responses on all endpoints should be "notset"
