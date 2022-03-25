@api
Feature: add groups
  As an administrator
  I want to be able to create group using the Graph API
  So that I can more easily manage access to resources by groups rather than individual users

  Scenario:
    When the administrator sends a group creation request for the following groups using the graph API
      | group_display_name |
      | simplegroup        |
      | España§àôœ€      |
      | नेपाली               |
    And the HTTP status code of responses on all endpoints should be "200"
    And these groups should exist:
      | groupname     |
      | simplegroup   |
      | España§àôœ€ |
      | नेपाली          |


  Scenario: admin creates a group with special characters
    When the administrator sends a group creation request for the following groups using the graph API
      | group_display_name  | comment            |
      | brand-new-group     | dash               |
      | the.group           | dot                |
      | left,right          | comma              |
      | 0                   | The "false" group  |
      | Finance (NP)        | Space and brackets |
      | Admin&Finance       | Ampersand          |
      | admin:Pokhara@Nepal | Colon and @        |
      | maint+eng           | Plus sign          |
      | $x<=>[y*z^2]!       | Maths symbols      |
      | Mgmt\Middle         | Backslash          |
      | 😅 😆               | emoji              |
      | [group1]            | brackets           |
      | group [ 2 ]         | bracketsAndSpace   |
    And the HTTP status code of responses on all endpoints should be "200"
    And these groups should exist:
      | groupname           |
      | brand-new-group     |
      | the.group           |
      | left,right          |
      | 0                   |
      | Finance (NP)        |
      | Admin&Finance       |
      | admin:Pokhara@Nepal |
      | maint+eng           |
      | $x<=>[y*z^2]!       |
      | Mgmt\Middle         |
      | 😅 😆               |
      | [group1]            |
      | group [ 2 ]         |
