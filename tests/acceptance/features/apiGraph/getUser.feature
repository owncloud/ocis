@api @skipOnOcV10
Feature: get users
  As an admin
  I want to be able to retrieve user information
  So that I can see the information

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files
    And the administrator has given "Alice" the role "Admin" using the settings api


  Scenario: admin user lists all the users
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Brian    |
      | Carol    |
    When user "Alice" gets all user using the Graph API
