@skipOnReva
Feature: sharing
  As a user
  I want to be able to share files when passwords are stored with the full hash difficulty
  So that I can give people secure controlled access to my data

  Note - this feature is run in CI with ACCOUNTS_HASH_DIFFICULTY set to the default for production
  See https://github.com/owncloud/ocis/issues/1542 and https://github.com/owncloud/ocis/pull/839


  Scenario Outline: creating a share of a file with a user
    Given using OCS API version "<ocs-api-version>"
    And user "Alice" has been created with default attributes
    And user "Alice" has uploaded file with content "ownCloud test text file 0" to "/textfile0.txt"
    And user "Brian" has been created with default attributes
    When user "Alice" shares file "textfile0.txt" with user "Brian" using the sharing API
    And the content of file "/Shares/textfile0.txt" for user "Brian" should be "ownCloud test text file 0"
    Examples:
      | ocs-api-version |
      | 1               |
      | 2               |
