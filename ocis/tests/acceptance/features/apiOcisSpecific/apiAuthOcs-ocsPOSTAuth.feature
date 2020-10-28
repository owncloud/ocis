@api
Feature: auth

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files

  @issue-ocis-ocs-26 @issue-ocis-reva-30
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario: send POST requests to OCS endpoints as normal user with wrong password
    When user "Alice" requests these endpoints with "POST" including body "doesnotmatter" using password "invalid" about user "Alice"
      | endpoint                                                        |
      | /ocs/v1.php/apps/files_sharing/api/v1/remote_shares/pending/123 |
      | /ocs/v2.php/apps/files_sharing/api/v1/remote_shares/pending/123 |
      | /ocs/v1.php/apps/files_sharing/api/v1/shares                    |
      | /ocs/v2.php/apps/files_sharing/api/v1/shares                    |
      | /ocs/v1.php/apps/files_sharing/api/v1/shares/pending/123        |
      | /ocs/v2.php/apps/files_sharing/api/v1/shares/pending/123        |
      | /ocs/v1.php/cloud/apps/testing                                  |
      | /ocs/v2.php/cloud/apps/testing                                  |
      | /ocs/v1.php/cloud/groups                                        |
      | /ocs/v2.php/cloud/groups                                        |
      | /ocs/v1.php/person/check                                        |
      | /ocs/v2.php/person/check                                        |
      | /ocs/v1.php/privatedata/deleteattribute/testing/test            |
      | /ocs/v2.php/privatedata/deleteattribute/testing/test            |
      | /ocs/v1.php/privatedata/setattribute/testing/test               |
      | /ocs/v2.php/privatedata/setattribute/testing/test               |
    Then the HTTP status code of responses on all endpoints should be "401"
    And the OCS status code of responses on all endpoints should be "notset"

  @issue-ocis-reva-30
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario: send POST requests to OCS endpoints as normal user with wrong password
    When user "Alice" requests these endpoints with "POST" including body "doesnotmatter" using password "invalid" about user "Alice"
      | endpoint                                     |
      | /ocs/v1.php/cloud/users                      |
      | /ocs/v2.php/cloud/users                      |
      | /ocs/v1.php/cloud/users/%username%/groups    |
      | /ocs/v2.php/cloud/users/%username%/groups    |
      | /ocs/v1.php/cloud/users/%username%/subadmins |
      | /ocs/v2.php/cloud/users/%username%/subadmins |
    Then the HTTP status code of responses on all endpoints should be "401"
    And the OCS status code of responses on all endpoints should be "notset"
