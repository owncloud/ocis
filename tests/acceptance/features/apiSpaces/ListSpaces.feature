@api @skipOnOcV10
Feature: List and create spaces
  As a user
  I want to be able to work with personal and project spaces to collaborate with individuals and teams

  Note - this feature is run in CI with ACCOUNTS_HASH_DIFFICULTY set to the default for production
  See https://github.com/owncloud/ocis/issues/1542 and https://github.com/owncloud/ocis/pull/839

  Scenario: list own spaces
    Given user "Alice" has been created with default attributes and without skeleton files
    And user "Alice" lists all available spaces via the GraphApi
    Then the HTTP status code should be "200"
    And the webDavUrl of the personal space has been found
    And user "Alice" lists the content of the personal space root using the WebDav Api
    And the HTTP status code should be "207"
