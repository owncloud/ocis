@skipOnReva @issue-1327
Feature: shares are received in the default folder for received shares
  As a user
  I want to share the default Shares folder
  So that I can make sure it does not work

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files


  Scenario Outline: do not allow sharing of the entire share folder
    Given using OCS API version "<ocs_api_version>"
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has created folder "FOLDER"
    When user "Alice" shares folder "/FOLDER" with user "Brian" using the sharing API
    And user "Brian" declines share "/Shares/FOLDER" offered by user "Alice" using the sharing API
    And user "Brian" shares folder "/Shares" with user "Alice" using the sharing API
    Then the OCS status code of responses on each endpoint should be "<ocs_status_code>" respectively
    And the HTTP status code of responses on each endpoint should be "<http_status_code>" respectively
    Examples:
      | ocs_api_version | ocs_status_code    | http_status_code   |
      | 1               | 100, 100, 100, 400 | 200, 200, 200, 200 |
      | 2               | 200, 200, 200, 400 | 200, 200, 200, 400 |
