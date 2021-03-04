@api
Feature: auth

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files

  @issue-ocis-reva-29
  @issue-ocis-reva-30
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario: using OCS anonymously
    When a user requests these endpoints with "GET" and no authentication
      | endpoint                                                    |
      | /ocs/v1.php/apps/files_external/api/v1/mounts               |
      | /ocs/v2.php/apps/files_external/api/v1/mounts               |
      | /ocs/v1.php/apps/files_sharing/api/v1/remote_shares         |
      | /ocs/v2.php/apps/files_sharing/api/v1/remote_shares         |
      | /ocs/v1.php/apps/files_sharing/api/v1/remote_shares/pending |
      | /ocs/v2.php/apps/files_sharing/api/v1/remote_shares/pending |
      | /ocs/v1.php/apps/files_sharing/api/v1/shares                |
      | /ocs/v2.php/apps/files_sharing/api/v1/shares                |
      | /ocs/v1.php/cloud/apps                                      |
      | /ocs/v2.php/cloud/apps                                      |
      | /ocs/v1.php/config                                          |
      | /ocs/v2.php/config                                          |
      | /ocs/v1.php/privatedata/getattribute                        |
      | /ocs/v2.php/privatedata/getattribute                        |
    Then the HTTP status code of responses on all endpoints should be "401"
    And the OCS status code of responses on all endpoints should be "notset"

  @issue-ocis-ocs-26
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario: using OCS anonymously
    When a user requests these endpoints with "GET" and no authentication
      | endpoint                 |
      | /ocs/v1.php/cloud/users  |
      | /ocs/v2.php/cloud/users  |
      | /ocs/v1.php/cloud/groups |
      | /ocs/v2.php/cloud/groups |
    Then the HTTP status code of responses on all endpoints should be "401"
    And the OCS status code of responses on all endpoints should be "997"


  @issue-ocis-reva-11
  @issue-ocis-reva-30
  @issue-ocis-reva-31
  @issue-ocis-reva-32
  @issue-ocis-reva-33
  @issue-ocis-reva-34
  @issue-ocis-reva-35
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario: using OCS with non-admin basic auth
    When the user "Alice" requests these endpoints with "GET" with basic auth
      | endpoint                                                    |
      | /ocs/v1.php/apps/files_external/api/v1/mounts               |
      | /ocs/v1.php/apps/files_sharing/api/v1/remote_shares         |
      | /ocs/v1.php/apps/files_sharing/api/v1/remote_shares/pending |
      | /ocs/v1.php/privatedata/getattribute                        |
      | /ocs/v1.php/cloud/apps                                      |
    Then the HTTP status code of responses on all endpoints should be "200"
    And the OCS status code of responses on all endpoints should be "998"
    When the user "Alice" requests these endpoints with "GET" with basic auth
      | endpoint           |
      | /ocs/v1.php/config |
    Then the HTTP status code of responses on all endpoints should be "200"
    And the OCS status code of responses on all endpoints should be "100"
    When the user "Alice" requests these endpoints with "GET" with basic auth
      | endpoint                                                    |
      | /ocs/v2.php/apps/files_external/api/v1/mounts               |
      | /ocs/v2.php/apps/files_sharing/api/v1/remote_shares         |
      | /ocs/v2.php/apps/files_sharing/api/v1/remote_shares/pending |
     # | /ocs/v1.php/apps/files_sharing/api/v1/shares                | 100      | 200       |
     # | /ocs/v2.php/apps/files_sharing/api/v1/shares                | 100      | 200       |

      | /ocs/v2.php/cloud/apps                                      |
      | /ocs/v2.php/privatedata/getattribute                        |
    Then the HTTP status code of responses on all endpoints should be "404"
    And the OCS status code of responses on all endpoints should be "998"
    When the user "Alice" requests these endpoints with "GET" with basic auth
      | endpoint                 |
      | /ocs/v1.php/cloud/users  |
      | /ocs/v2.php/cloud/users  |
      | /ocs/v1.php/cloud/groups |
      | /ocs/v2.php/cloud/groups |
    Then the HTTP status code of responses on all endpoints should be "401"
    And the OCS status code of responses on all endpoints should be "997"
    When the user "Alice" requests these endpoints with "GET" with basic auth
      | endpoint           |
      | /ocs/v2.php/config |
    Then the HTTP status code of responses on all endpoints should be "200"
    And the OCS status code of responses on all endpoints should be "200"

  @issue-ocis-reva-29
  @issue-ocis-reva-30
  @issue-ocis-accounts-73
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario: using OCS as normal user (username has a capital letter) with wrong password
    When user "Alice" requests these endpoints with "GET" using password "invalid"
      | endpoint                                                    |
      | /ocs/v1.php/apps/files_external/api/v1/mounts               |
      | /ocs/v2.php/apps/files_external/api/v1/mounts               |
      | /ocs/v1.php/apps/files_sharing/api/v1/remote_shares         |
      | /ocs/v2.php/apps/files_sharing/api/v1/remote_shares         |
      | /ocs/v1.php/apps/files_sharing/api/v1/remote_shares/pending |
      | /ocs/v2.php/apps/files_sharing/api/v1/remote_shares/pending |
      | /ocs/v1.php/apps/files_sharing/api/v1/shares                |
      | /ocs/v2.php/apps/files_sharing/api/v1/shares                |
      | /ocs/v1.php/cloud/apps                                      |
      | /ocs/v2.php/cloud/apps                                      |
      | /ocs/v1.php/cloud/groups                                    |
      | /ocs/v2.php/cloud/groups                                    |
      | /ocs/v1.php/config                                          |
      | /ocs/v2.php/config                                          |
      | /ocs/v1.php/privatedata/getattribute                        |
      | /ocs/v2.php/privatedata/getattribute                        |
    Then the HTTP status code of responses on all endpoints should be "401"
    And the OCS status code of responses on all endpoints should be "notset"

  @issue-ocis-reva-29
  @issue-ocis-reva-30
  @issue-ocis-accounts-73
  @issue-ocis-ocs-26
  @smokeTest
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario: using OCS as normal user (username has a capital letter) with wrong password
    When user "Alice" requests these endpoints with "GET" using password "invalid"
      | endpoint                |
      | /ocs/v1.php/cloud/users |
      | /ocs/v2.php/cloud/users |
    Then the HTTP status code of responses on all endpoints should be "401"
    And the OCS status code of responses on all endpoints should be "notset"

  @skipOnOcV10
  @issue-ocis-reva-29
  @issue-ocis-reva-30
  @issue-ocis-accounts-73
  @issue-ocis-ocs-26
  @smokeTest
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario: using OCS as normal user (username does not have a capital letter) with wrong password
    Given user "brian" has been created with default attributes and without skeleton files
    When user "brian" requests these endpoints with "GET" using password "invalid"
      | endpoint                                                    |
      | /ocs/v1.php/apps/files_external/api/v1/mounts               |
      | /ocs/v2.php/apps/files_external/api/v1/mounts               |
      | /ocs/v1.php/apps/files_sharing/api/v1/remote_shares         |
      | /ocs/v2.php/apps/files_sharing/api/v1/remote_shares         |
      | /ocs/v1.php/apps/files_sharing/api/v1/remote_shares/pending |
      | /ocs/v2.php/apps/files_sharing/api/v1/remote_shares/pending |
      | /ocs/v1.php/apps/files_sharing/api/v1/shares                |
      | /ocs/v2.php/apps/files_sharing/api/v1/shares                |
      | /ocs/v1.php/cloud/apps                                      |
      | /ocs/v2.php/cloud/apps                                      |
      | /ocs/v1.php/cloud/groups                                    |
      | /ocs/v2.php/cloud/groups                                    |
      | /ocs/v1.php/config                                          |
      | /ocs/v2.php/config                                          |
      | /ocs/v1.php/privatedata/getattribute                        |
      | /ocs/v2.php/privatedata/getattribute                        |
    Then the HTTP status code of responses on all endpoints should be "401"
    And the OCS status code of responses on all endpoints should be "notset"

  @skipOnOcV10
  @issue-ocis-reva-29
  @issue-ocis-reva-30
  @issue-ocis-accounts-73
  @issue-ocis-ocs-26
  @smokeTest
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario: using OCS as normal user (username does not have a capital letter) with wrong password
    Given user "brian" has been created with default attributes and without skeleton files
    When user "brian" requests these endpoints with "GET" using password "invalid"
      | endpoint                |
      | /ocs/v1.php/cloud/users |
      | /ocs/v2.php/cloud/users |
    Then the HTTP status code of responses on all endpoints should be "401"
    And the OCS status code of responses on all endpoints should be "notset"
