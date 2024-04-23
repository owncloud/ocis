Feature: Remove access to a drive
  https://owncloud.dev/libre-graph-api/#/drives.root/DeletePermissionSpaceRoot

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
      | Carol    |
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


  Scenario Outline: user of a group removes another user from project space
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And group "group1" has been created
    And user "Brian" has been added to group "group1"
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has sent the following share invitation:
      | space           | NewSpace           |
      | sharee          | group1             |
      | shareType       | group              |
      | permissionsRole | <permissions-role> |
    And user "Alice" has sent the following share invitation:
      | space           | NewSpace     |
      | sharee          | Carol        |
      | shareType       | user         |
      | permissionsRole | Space Viewer |
    When user "Brian" removes the access of user "Carol" from space "NewSpace" using root endpoint of the Graph API
    Then the HTTP status code should be "<status-code>"
    And the user "Carol" <shouldOrNot> have a space called "NewSpace"
    Examples:
      | permissions-role | status-code | shouldOrNot |
      | Space Viewer     | 403         | should      |
      | Space Editor     | 403         | should      |
      | Manager          | 204         | should not  |


  Scenario Outline: user of a group removes own group from project space
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And group "group1" has been created
    And user "Brian" has been added to group "group1"
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has sent the following share invitation:
      | space           | NewSpace           |
      | sharee          | group1             |
      | shareType       | group              |
      | permissionsRole | <permissions-role> |
    When user "Brian" removes the access of group "group1" from space "NewSpace" using root endpoint of the Graph API
    Then the HTTP status code should be "<status-code>"
    And the user "Brian" <shouldOrNot> have a space called "NewSpace"
    Examples:
      | permissions-role | status-code | shouldOrNot |
      | Space Viewer     | 403         | should      |
      | Space Editor     | 403         | should      |
      | Manager          | 204         | should not  |


  @issue-8819
  Scenario Outline: user removes himself from the project space
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has sent the following share invitation:
      | space           | NewSpace           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    When user "Brian" removes the access of user "Brian" from space "NewSpace" using root endpoint of the Graph API
    Then the HTTP status code should be "<status-code>"
    And the user "Brian" <shouldOrNot> have a space called "NewSpace"
    Examples:
      | permissions-role | status-code | shouldOrNot |
      | Space Viewer     | 403         | should      |
      | Space Editor     | 403         | should      |
      | Manager          | 204         | should not  |

  @issue-8819
  Scenario Outline: user removes another user from project space using permissions endpoint
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has sent the following share invitation:
      | space           | NewSpace           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    And user "Alice" has sent the following share invitation:
      | space           | NewSpace     |
      | sharee          | Carol        |
      | shareType       | user         |
      | permissionsRole | Space Viewer |
    When user "Brian" removes the access of user "Carol" from space "NewSpace" using permissions endpoint of the Graph API
    Then the HTTP status code should be "<status-code>"
    And the user "Carol" <shouldOrNot> have a space called "NewSpace"
     Examples:
      | permissions-role | status-code | shouldOrNot |
      | Space Viewer     | 403         | should      |
      | Space Editor     | 403         | should      |
      | Manager          | 204         | should not  |
