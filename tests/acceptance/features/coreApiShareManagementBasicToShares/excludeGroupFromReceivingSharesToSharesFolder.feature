@api @files_sharing-app-required
Feature: Exclude groups from receiving shares
  As an admin
  I want to exclude groups from receiving shares
  So that users do not mistakenly share with groups they should not e.g. huge meta groups

  Background:
    Given the administrator has set the default folder for received shares to "Shares"
    And auto-accept shares has been disabled
    And these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
      | Carol    |
      | David    |
    And group "grp1" has been created
    And group "grp2" has been created
    And user "Brian" has been added to group "grp1"
    And user "David" has been added to group "grp2"


  Scenario Outline: sharing with a user that is part of a group that is excluded from receiving shares still works
    Given using OCS API version "<ocs_api_version>"
    And user "Alice" has created folder "PARENT"
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "fileToShare.txt"
    And the administrator has added group "grp1" to the exclude groups from receiving shares list
    When user "Alice" shares file "fileToShare.txt" with user "Brian" using the sharing API
    Then the OCS status code should be "<ocs_status_code>"
    And the HTTP status code should be "200"
    When user "Alice" shares folder "PARENT" with user "Brian" using the sharing API
    Then the OCS status code should be "<ocs_status_code>"
    And the HTTP status code should be "200"
    When user "Brian" accepts share "/fileToShare.txt" offered by user "Alice" using the sharing API
    And user "Brian" accepts share "/PARENT" offered by user "Alice" using the sharing API
    Then as "Brian" file "/Shares/fileToShare.txt" should exist
    And as "Brian" folder "/Shares/PARENT" should exist
    Examples:
      | ocs_api_version | ocs_status_code |
      | 1               | 100             |
      | 2               | 200             |


  Scenario Outline: sharing with a user that is part of a group that is excluded from receiving shares using an other group works
    Given using OCS API version "<ocs_api_version>"
    And group "grp3" has been created
    And user "Brian" has been added to group "grp3"
    And user "Alice" has created folder "PARENT"
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "fileToShare.txt"
    And the administrator has added group "grp1" to the exclude groups from receiving shares list
    When user "Alice" shares file "fileToShare.txt" with group "grp3" using the sharing API
    Then the OCS status code should be "<ocs_status_code>"
    And the HTTP status code should be "200"
    When user "Alice" shares folder "PARENT" with group "grp3" using the sharing API
    Then the OCS status code should be "<ocs_status_code>"
    And the HTTP status code should be "200"
    When user "Brian" accepts share "/fileToShare.txt" offered by user "Alice" using the sharing API
    And user "Brian" accepts share "/PARENT" offered by user "Alice" using the sharing API
    Then as "Brian" file "/Shares/fileToShare.txt" should exist
    And as "Brian" folder "/Shares/PARENT" should exist
    Examples:
      | ocs_api_version | ocs_status_code |
      | 1               | 100             |
      | 2               | 200             |


  Scenario Outline: a user that is part of a group that is excluded from receiving shares still can initiate shares
    Given using OCS API version "<ocs_api_version>"
    And user "Brian" has created folder "PARENT"
    And user "Brian" has uploaded file "filesForUpload/textfile.txt" to "fileToShare.txt"
    And the administrator has added group "grp1" to the exclude groups from receiving shares list
    When user "Brian" shares file "fileToShare.txt" with user "Carol" using the sharing API
    Then the OCS status code should be "<ocs_status_code>"
    And the HTTP status code should be "200"
    And user "Brian" shares folder "PARENT" with user "Carol" using the sharing API
    And the OCS status code should be "<ocs_status_code>"
    And the HTTP status code should be "200"
    When user "Carol" accepts share "/fileToShare.txt" offered by user "Brian" using the sharing API
    And user "Carol" accepts share "/PARENT" offered by user "Brian" using the sharing API
    Then as "Carol" file "/Shares/fileToShare.txt" should exist
    And as "Carol" folder "/Shares/PARENT" should exist
    Examples:
      | ocs_api_version | ocs_status_code |
      | 1               | 100             |
      | 2               | 200             |
