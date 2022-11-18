@api @skipOnOcV10
Feature: delete groups
  As an admin
  I want to be able to delete groups
  So that I can remove unnecessary groups

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files
    And the administrator has given "Alice" the role "Admin" using the settings api


  Scenario Outline: admin user deletes a group
    Given group "<group_id>" has been created
    When user "Alice" deletes group "<group_id>" using the Graph API
    And the HTTP status code should be "204"
    And group "<group_id>" should not exist
    Examples:
      | group_id            | comment                               |
      | simplegroup         | nothing special here                  |
      | España§àôœ€         | special European and other characters |
      | नेपाली                | Unicode group name                    |
      | brand-new-group     | dash                                  |
      | the.group           | dot                                   |
      | left,right          | comma                                 |
      | 0                   | The "false" group                     |
      | Finance (NP)        | Space and brackets                    |
      | Admin&Finance       | Ampersand                             |
      | admin:Pokhara@Nepal | Colon and @                           |
      | maint+eng           | Plus sign                             |
      | $x<=>[y*z^2]!       | Maths symbols                         |
      | Mgmt\Middle         | Backslash                             |
      | 😁 😂               | emoji                                 |
      | maintenance#123     | Hash sign                             |
      | 50%25=0             | %25 literal looks like an escaped "%" |
      | staff?group         | Question mark                         |

  @cannot-create-group
  Scenario Outline: admin deletes a group
    Given group "<group_id>" has been created
    When user "Alice" deletes group "<group_id>" using the Graph API
    And the HTTP status code should be "204"
    And group "<group_id>" should not exist
    Examples:
      | group_id            | comment                                 |
      | 50%pass             | Percent sign (special escaping happens) |
      | 50%2Eagle           | %2E literal looks like an escaped "."   |
      | 50%2Fix             | %2F literal looks like an escaped slash |


  Scenario Outline: admin deletes a group that has a forward-slash in the group name
    Given group "<group_id>" has been created
    When user "Alice" deletes group "<group_id>" using the Graph API
    And the HTTP status code should be "204"
    And group "<group_id>" should not exist
    Examples:
      | group_id         | comment                            |
      | Mgmt/Sydney      | Slash (special escaping happens)   |
      | Mgmt//NSW/Sydney | Multiple slash                     |
      | priv/subadmins/1 | Subadmins mentioned not at the end |
      | var/../etc       | using slash-dot-dot                |


  Scenario: normal user tries to delete a group
    Given user "Brian" has been created with default attributes and without skeleton files
    And group "new-group" has been created
    When user "Brian" tries to delete group "new-group" using the Graph API
    And the HTTP status code should be "401"
    And group "new-group" should exist
