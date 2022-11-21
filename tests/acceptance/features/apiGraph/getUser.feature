@api
Feature: get user information
  As an admin, subadmin or as myself
  I want to be able to retrieve user information
  So that I can see the information


  Scenario: admin gets an existing user
    Given these users have been created with default attributes and without skeleton files:
      | username       | displayname    |
      | brand-new-user | Brand New User |
