Feature: create a public link share when share_folder is set to Shares
  As a user
  I want to create public links
  So that I can share resources to people who aren't owncloud users

  Background:
    Given user "Alice" has been created with default attributes


  Scenario Outline: creating a new public link share of a file gives the correct response
    Given using OCS API version "<ocs-api-version>"
    And user "Alice" has uploaded file with content "Random data" to "/randomfile.txt"
    When user "Alice" creates a public link share using the sharing API with settings
      | path     | randomfile.txt |
      | password | %public%       |
    Then the OCS status code should be "<ocs-status-code>"
    And the HTTP status code should be "200"
    And the fields of the last response to user "Alice" should include
      | item_type              | file            |
      | mimetype               | text/plain      |
      | file_target            | /randomfile.txt |
      | path                   | /randomfile.txt |
      | permissions            | read            |
      | share_type             | public_link     |
      | displayname_file_owner | %displayname%   |
      | displayname_owner      | %displayname%   |
      | uid_file_owner         | %username%      |
      | uid_owner              | %username%      |
      | name                   |                 |
    Examples:
      | ocs-api-version | ocs-status-code |
      | 1               | 100             |
      | 2               | 200             |
