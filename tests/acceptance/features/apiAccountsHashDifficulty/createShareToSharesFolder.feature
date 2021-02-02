@api @files_sharing-app-required
Feature: sharing
  As a user
  I want to be able to share files when passwords are stored with the full hash difficulty
  So that I can give people secure controlled access to my data

  Note - this feature is run in CI with ACCOUNTS_HASH_DIFFICULTY set to the default for production
  See https://github.com/owncloud/ocis/issues/1542 and https://github.com/owncloud/ocis/pull/839

  Scenario Outline: Creating a share of a file with a user
    Given the administrator has set the default folder for received shares to "Shares"
    And auto-accept shares has been disabled
    And using OCS API version "<ocs_api_version>"
    And user "Alice" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file with content "ownCloud test text file 0" to "/textfile0.txt"
    And user "Brian" has been created with default attributes and without skeleton files
    When user "Alice" shares file "textfile0.txt" with user "Brian" using the sharing API
    And user "Brian" accepts share "/textfile0.txt" offered by user "Alice" using the sharing API
    Then the HTTP status code should be "200"
    And the content of file "/Shares/textfile0.txt" for user "Brian" should be "ownCloud test text file 0"
    Examples:
      | ocs_api_version |
      | 1               |
      | 2               |
