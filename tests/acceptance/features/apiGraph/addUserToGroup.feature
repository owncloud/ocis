@api
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
      | नेपाली        | Unicode group name                    |
    When the administrator adds the following users to the following groups using the Graph API
      | username | groupname   |
      | Alice    | simplegroup |
      | Alice    | España§àôœ€ |
      | Alice    | नेपाली        |
    And the HTTP status code of responses on all endpoints should be "204"


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
    And the HTTP status code of responses on all endpoints should be "204"


  Scenario: adding a user to a group with % and # in its name
    Given these groups have been created:
      | groupname           | comment                                 |
      | maintenance#123     | Hash sign                               |
      | 50%pass             | Percent sign (special escaping happens) |
      | 50%25=0             | %25 literal looks like an escaped "%"   |
      | 50%2Eagle           | %2E literal looks like an escaped "."   |
      | 50%2Fix             | %2F literal looks like an escaped slash |
      | Mgmt\Middle         | Backslash                               |
      | staff?group         | Question mark                           |
    When the administrator adds the following users to the following groups using the Graph API
      | username | groupname       |
      | Alice    | maintenance#123 |
      | Alice    | 50%pass         |
      | Alice    | 50%25=0         |
      | Alice    | 50%2Eagle       |
      | Alice    | 50%2Fix         |
      | Alice    | Mgmt\Middle     |
      | Alice    | staff?group     |
    And the HTTP status code of responses on all endpoints should be "204"


  Scenario: adding a user to a group that has a forward-slash in the group name
    Given these groups have been created:
      | groupname        | comment                            |
      | Mgmt/Sydney      | Slash (special escaping happens)   |
      | Mgmt//NSW/Sydney | Multiple slash                     |
      | priv/subadmins/1 | Subadmins mentioned not at the end |
    When the administrator adds the following users to the following groups using the Graph API
      | username | groupname        |
      | Alice    | Mgmt/Sydney      |
      | Alice    | Mgmt//NSW/Sydney |
      | Alice    | priv/subadmins/1 |
    And the HTTP status code of responses on all endpoints should be "204"


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

  @skipOnLDAP
  Scenario: admin tries to add user to a group which does not exist
    Given user "Alice" has been created with default attributes and without skeleton files
    And group "nonexistentgroup" has been deleted
    When the administrator tries to add user "Alice" to group "nonexistentgroup" using the provisioning API
    Then the OCS status code should be "400"
    And the HTTP status code should be "400"
    And the API should not return any data

  @skipOnLDAP
  Scenario: admin tries to add user to a group without sending the group
    Given user "Alice" has been created with default attributes and without skeleton files
    When the administrator tries to add user "Alice" to group "" using the provisioning API
    Then the OCS status code should be "400"
    And the HTTP status code should be "400"
    And the API should not return any data

  @skipOnLDAP
  Scenario: admin tries to add a user which does not exist to a group
    Given user "nonexistentuser" has been deleted
    And group "brand-new-group" has been created
    When the administrator tries to add user "nonexistentuser" to group "brand-new-group" using the provisioning API
    Then the OCS status code should be "400"
    And the HTTP status code should be "400"
    And the API should not return any data

  @skipOnLDAP @notToImplementOnOCIS
  Scenario: subadmin adds users to groups the subadmin is responsible for
    Given these users have been created with default attributes and without skeleton files:
      | username       |
      | Alice |
      | subadmin       |
    And group "brand-new-group" has been created
    And user "subadmin" has been made a subadmin of group "brand-new-group"
    When user "subadmin" tries to add user "Alice" to group "brand-new-group" using the provisioning API
    Then the OCS status code should be "403"
    And the HTTP status code should be "403"
    And user "Alice" should not belong to group "brand-new-group"

  @skipOnLDAP @notToImplementOnOCIS
  Scenario: subadmin tries to add user to groups the subadmin is not responsible for
    Given these users have been created with default attributes and without skeleton files:
      | username         |
      | Alice   |
      | another-subadmin |
    And group "brand-new-group" has been created
    And group "another-new-group" has been created
    And user "another-subadmin" has been made a subadmin of group "another-new-group"
    When user "another-subadmin" tries to add user "Alice" to group "brand-new-group" using the provisioning API
    Then the OCS status code should be "403"
    And the HTTP status code should be "403"
    And user "Alice" should not belong to group "brand-new-group"

  @skipOnLDAP @skipOnOcV10.6 @skipOnOcV10.7 @skipOnOcV10.8.0 @notToImplementOnOCIS
  Scenario: a subadmin can add users to other groups the subadmin is responsible for
    Given these users have been created with default attributes and without skeleton files:
      | username         |
      | Alice   |
      | another-subadmin |
    And group "brand-new-group" has been created
    And group "another-new-group" has been created
    And user "Alice" has been added to group "brand-new-group"
    And user "another-subadmin" has been made a subadmin of group "brand-new-group"
    And user "another-subadmin" has been made a subadmin of group "another-new-group"
    When user "another-subadmin" tries to add user "Alice" to group "another-new-group" using the provisioning API
    Then the OCS status code should be "200"
    And the HTTP status code should be "200"
    And user "Alice" should belong to group "brand-new-group"

  # merge this with scenario on line 62 once the issue is fixed
  @issue-31015 @skipOnLDAP @toImplementOnOCIS @issue-product-284
  Scenario Outline: adding a user to a group that has a forward-slash and dot in the group name
    Given user "Alice" has been created with default attributes and without skeleton files
    And the administrator sends a group creation request for group "<group_id>" using the provisioning API
    When the administrator adds user "Alice" to group "<group_id>" using the provisioning API
    Then the OCS status code should be "200"
    And the HTTP status code should be "200"
    Examples:
      | group_id         | comment                            |
      | var/../etc       | using slash-dot-dot                |
