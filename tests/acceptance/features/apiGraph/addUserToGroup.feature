@api @skipOnOcV10
Feature: add users to group
  As a admin
  I want to be able to add users to a group
  So that I can give a user access to the resources of the group

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files


  Scenario: adding a user to a group
    Given these groups have been created:
      | groupname   | comment                               |
      | simplegroup | nothing special here                  |
      | España§àôœ€ | special European and other characters |
      | नेपाली      | Unicode group name                    |
    When the administrator adds the following users to the following groups using the Graph API
      | username | groupname   |
      | Alice    | simplegroup |
      | Alice    | España§àôœ€ |
      | Alice    | नेपाली      |
    Then the HTTP status code of responses on all endpoints should be "204"
    And the following users should be listed in the following groups
      | username | groupname   |
      | Alice    | simplegroup |
      | Alice    | España§àôœ€ |
      | Alice    | नेपाली      |


  Scenario: adding a user to a group with special character in its name
    Given these groups have been created:
      | groupname           | comment            |
      | brand-new-group     | dash               |
      | the.group           | dot                |
      | left,right          | comma              |
      | 0                   | The "false" group  |
      | Finance (NP)        | Space and brackets |
      | Admin&Finance       | Ampersand          |
      | maint+eng           | Plus sign          |
      | $x<=>[y*z^2]!       | Maths symbols      |
      | 😁 😂               | emoji              |
      | admin:Pokhara@Nepal | Colon and @        |
    When the administrator adds the following users to the following groups using the Graph API
      | username | groupname           |
      | Alice    | brand-new-group     |
      | Alice    | the.group           |
      | Alice    | left,right          |
      | Alice    | 0                   |
      | Alice    | Finance (NP)        |
      | Alice    | Admin&Finance       |
      | Alice    | maint+eng           |
      | Alice    | $x<=>[y*z^2]!       |
      | Alice    | 😁 😂               |
      | Alice    | admin:Pokhara@Nepal |
    Then the HTTP status code of responses on all endpoints should be "204"
    And the following users should be listed in the following groups
      | username | groupname           |
      | Alice    | brand-new-group     |
      | Alice    | the.group           |
      | Alice    | left,right          |
      | Alice    | 0                   |
      | Alice    | Finance (NP)        |
      | Alice    | Admin&Finance       |
      | Alice    | maint+eng           |
      | Alice    | $x<=>[y*z^2]!       |
      | Alice    | 😁 😂               |
      | Alice    | admin:Pokhara@Nepal |


  Scenario: adding a user to a group with % and # in its name
    Given these groups have been created:
      | groupname       | comment                                 |
      | maintenance#123 | Hash sign                               |
      | 50%pass         | Percent sign (special escaping happens) |
      | 50%25=0         | %25 literal looks like an escaped "%"   |
      | 50%2Eagle       | %2E literal looks like an escaped "."   |
      | 50%2Fix         | %2F literal looks like an escaped slash |
      | Mgmt\Middle     | Backslash                               |
      | staff?group     | Question mark                           |
    When the administrator adds the following users to the following groups using the Graph API
      | username | groupname       |
      | Alice    | maintenance#123 |
      | Alice    | 50%pass         |
      | Alice    | 50%25=0         |
      | Alice    | 50%2Eagle       |
      | Alice    | 50%2Fix         |
      | Alice    | Mgmt\Middle     |
      | Alice    | staff?group     |
    Then the HTTP status code of responses on all endpoints should be "204"
    And the following users should be listed in the following groups
      | username | groupname       |
      | Alice    | maintenance#123 |
      | Alice    | 50%pass         |
      | Alice    | 50%25=0         |
      | Alice    | 50%2Eagle       |
      | Alice    | 50%2Fix         |
      | Alice    | Mgmt\Middle     |
      | Alice    | staff?group     |


  Scenario: adding a user to a group that has a forward-slash in the group name
    Given these groups have been created:
      | groupname        | comment                            |
      | Mgmt/Sydney      | Slash (special escaping happens)   |
      | Mgmt//NSW/Sydney | Multiple slash                     |
      | priv/subadmins/1 | Subadmins mentioned not at the end |
      | var/../etc       | using slash-dot-dot                |
    When the administrator adds the following users to the following groups using the Graph API
      | username | groupname        |
      | Alice    | Mgmt/Sydney      |
      | Alice    | Mgmt//NSW/Sydney |
      | Alice    | priv/subadmins/1 |
      | Alice    | var/../etc       |
    Then the HTTP status code of responses on all endpoints should be "204"
    And the following users should be listed in the following groups
      | username | groupname        |
      | Alice    | Mgmt/Sydney      |
      | Alice    | Mgmt//NSW/Sydney |
      | Alice    | priv/subadmins/1 |
      | Alice    | var/../etc       |


  Scenario: normal user tries to add himself to a group
    Given group "groupA" has been created
    When user "Alice" tries to add himself to group "groupA" using the Graph API
    Then the HTTP status code should be "401"
    And the last response should be an unauthorized response


  Scenario: normal user tries to other user to a group
    Given user "Brian" has been created with default attributes and without skeleton files
    And group "groupA" has been created
    When user "Alice" tries to add user "Brian" to group "groupA" using the Graph API
    Then the HTTP status code should be "401"
    And the last response should be an unauthorized response


  Scenario: admin tries to add user to a non-existing group
    When the administrator tries to add user "Alice" to group "nonexistentgroup" using the Graph API
    Then the HTTP status code should be "404"


  Scenario: admin tries to add a non-existing user to a group
    Given group "groupA" has been created
    When the administrator tries to add user "nonexistentuser" to group "groupA" using the provisioning API
    Then the HTTP status code should be "405"


  Scenario: admin tries to add user to a group without sending the group
    When the administrator tries to add user "Alice" to group "" using the Graph API
    Then the HTTP status code should be "404"


  Scenario: adding a disabled user to a group
    Given these groups have been created:
      | groupname | comment      |
      | sales     | normal group |
    And the user "Admin" has disabled user "Alice" using the Graph API
    When the administrator adds the following users to the following groups using the Graph API
      | username | groupname |
      | Alice    | sales     |
    Then the HTTP status code of responses on all endpoints should be "204"
    And the following users should be listed in the following groups
      | username | groupname |
      | Alice    | sales     |
