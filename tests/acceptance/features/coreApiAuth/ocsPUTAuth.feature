Feature: auth
  As a user
  I want to send PUT request to various endpoints
  So that I can make sure the endpoints need proper authentication

  @issue-1337 @smokeTest
  Scenario: send PUT request to OCS endpoints as admin with wrong password
    When user "admin" requests these endpoints with "PUT" including body "doesnotmatter" using password "invalid" about user "Alice"
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
    And the OCS status code of responses on all endpoints should be "997"


  Scenario: request to edit nonexistent user by authorized admin gets unauthorized in http response
    When user "admin" requests these endpoints with "PUT" including body "doesnotmatter" about user "nonexistent"
      | endpoint                                         |
      | /ocs/v1.php/cloud/users/%username%               |
      | /ocs/v2.php/cloud/users/%username%               |
    Then the HTTP status code of responses on all endpoints should be "200"
    And the OCS status code of responses on all endpoints should be "101"
