@api @skipOnOcV10
Feature: List and create spaces
  As a user
  I want to be able to work with personal and project spaces to collaborate with individuals and teams

  Note - this feature is run in CI with ACCOUNTS_HASH_DIFFICULTY set to the default for production
  See https://github.com/owncloud/ocis/issues/1542 and https://github.com/owncloud/ocis/pull/839

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files

  Scenario: Alice has an empty space when she has logged in for the first time
    When user "Alice" lists all available spaces via the GraphApi
    Then the HTTP status code should be "200"
    And the quota total of space "Alice Hansen" should be set
    And the quota used of space "Alice Hansen" should be "0"
    And the quota remaining of space "Alice Hansen" should be set
    And the quota state of space "Alice Hansen" should be "normal"
