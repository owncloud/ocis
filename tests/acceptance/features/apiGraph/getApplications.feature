@api @skipOnOcV10
Feature: get applications
  As an user
  I want to be able to get applications information with existings roles

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files


  Scenario Outline: admin user lists all the groups
    Given the administrator has given "Alice" the role "<role>" using the settings api
    When user "Alice" gets all applications using the Graph API
    Then the HTTP status code should be "200"
    Then the user retrieve API response should contain the following applications information:
      | key                        | value                   |
      | displayName                | ownCloud Infinite Scale |
      | id                         | %uuid_v4%               |
      | appRoles@@@0@@@displayName | Admin                   |
      | appRoles@@@0@@@id          | %uuid_v4%               |
      | appRoles@@@1@@@displayName | User                    |
      | appRoles@@@1@@@id          | %uuid_v4%               |
      | appRoles@@@2@@@displayName | Space Admin             |
      | appRoles@@@2@@@id          | %uuid_v4%               |
      | appRoles@@@3@@@displayName | Guest                   |
      | appRoles@@@3@@@id          | %uuid_v4%               |
    Examples:
      | role        |
      | Admin       |
      | Space Admin |
      | User        |
      | Guest       |