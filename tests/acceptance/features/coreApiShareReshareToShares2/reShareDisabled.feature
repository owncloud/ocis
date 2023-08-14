@issue-1328
Feature: resharing can be disabled
  As a user
  I want to share a resource without reshare permission
  So that the resource won't be accessible to unwanted individuals


  Scenario Outline: ordinary sharing is allowed when allow resharing has been disabled
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "/textfile0.txt"
    And using OCS API version "<ocs_api_version>"
    When user "Alice" shares file "/textfile0.txt" with user "Brian" with permissions "share,update,read" using the sharing API
    Then the OCS status code should be "<ocs_status_code>"
    And the HTTP status code should be "200"
    And user "Brian" should be able to accept pending share "/textfile0.txt" offered by user "Alice"
    And as "Brian" file "/Shares/textfile0.txt" should exist
    Examples:
      | ocs_api_version | ocs_status_code |
      | 1               | 100             |
      | 2               | 200             |
