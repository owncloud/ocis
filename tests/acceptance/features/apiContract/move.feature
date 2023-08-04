@api
Feature: move test
  As a user
  I want to check the MOVE request
  So that I can make sure that the resource moved successfully


  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
    And using spaces DAV path

    Scenario: move file
