Feature: Remove access to a drive
  https://owncloud.dev/libre-graph-api/#/drives.root/DeletePermissionSpaceRoot
  
  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
    And using spaces DAV path


  Scenario Outline: user removes user member from project space
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has sent the following share invitation:
      | space           | NewSpace           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    When user "Alice" removes the access of user "Brian" from space "NewSpace" using root endpoint of the Graph API
    Then the HTTP status code should be "204"
    And the user "Brian" should not have a space called "NewSpace"
    Examples:
      | permissions-role |
      | Space Viewer     |
      | Space Editor     |
      | Manager          |

  @issue-8768
  Scenario Outline: user removes group from project space
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And group "group1" has been created
    And user "Brian" has been added to group "group1"
    And user "Alice" has sent the following share invitation:
      | space           | NewSpace           |
      | sharee          | group1             |
      | shareType       | group              |
      | permissionsRole | <permissions-role> |
    When user "Alice" removes the access of group "group1" from space "NewSpace" using root endpoint of the Graph API
    Then the HTTP status code should be "204"
    And the user "Brian" should not have a space called "NewSpace"
    Examples:
      | permissions-role |
      | Space Viewer     |
      | Space Editor     |
      | Manager          |
