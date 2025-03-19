Feature: Remove access to a drive item
  https://owncloud.dev/libre-graph-api/#/drives.permissions/DeletePermission

  Background:
    Given these users have been created with default attributes:
      | username |
      | Alice    |
      | Brian    |
    And using spaces DAV path


  Scenario Outline: user removes access to resource in the user share
    Given user "Alice" has created folder "FolderToShare"
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "textfile.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | <resource>         |
      | space           | Personal           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    When user "Alice" removes the access of user "Brian" from resource "<resource>" of space "Personal" using the Graph API
    Then the HTTP status code should be "204"
    And for user "Brian" the space "Shares" should not contain these entries:
      | <resource> |
    Examples:
      | permissions-role | resource      |
      | Viewer           | textfile.txt  |
      | File Editor      | textfile.txt  |
      | Viewer           | FolderToShare |
      | Editor           | FolderToShare |
      | Uploader         | FolderToShare |


  Scenario Outline: user removes access to resource inside of a project space in the user share
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has created a folder "FolderToShare" in space "NewSpace"
    And user "Alice" has uploaded a file inside space "NewSpace" with content "some content" to "textfile.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | <resource>         |
      | space           | NewSpace           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    When user "Alice" removes the access of user "Brian" from resource "<resource>" of space "NewSpace" using the Graph API
    Then the HTTP status code should be "204"
    And for user "Brian" the space "Shares" should not contain these entries:
      | <resource> |
    Examples:
      | permissions-role | resource      |
      | Viewer           | textfile.txt  |
      | File Editor      | textfile.txt  |
      | Viewer           | FolderToShare |
      | Editor           | FolderToShare |
      | Uploader         | FolderToShare |


  Scenario Outline: user removes access to a resource in a group share
    Given group "group1" has been created
    And user "Brian" has been added to group "group1"
    And user "Alice" has been added to group "group1"
    And user "Alice" has created folder "FolderToShare"
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "textfile.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | <resource>         |
      | space           | Personal           |
      | sharee          | group1             |
      | shareType       | group              |
      | permissionsRole | <permissions-role> |
    When user "Alice" removes the access of group "group1" from resource "<resource>" of space "Personal" using the Graph API
    Then the HTTP status code should be "204"
    And for user "Brian" the space "Shares" should not contain these entries:
      | <resource> |
    Examples:
      | permissions-role | resource      |
      | Viewer           | textfile.txt  |
      | File Editor      | textfile.txt  |
      | Viewer           | FolderToShare |
      | Editor           | FolderToShare |
      | Uploader         | FolderToShare |


  Scenario Outline: user removes access to a resource inside of a project space in group share
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has created a folder "FolderToShare" in space "NewSpace"
    And user "Alice" has uploaded a file inside space "NewSpace" with content "some content" to "textfile.txt"
    And group "group1" has been created
    And user "Brian" has been added to group "group1"
    And user "Alice" has been added to group "group1"
    And user "Alice" has sent the following resource share invitation:
      | resource        | <resource>         |
      | space           | NewSpace           |
      | sharee          | group1             |
      | shareType       | group              |
      | permissionsRole | <permissions-role> |
    When user "Alice" removes the access of group "group1" from resource "<resource>" of space "NewSpace" using the Graph API
    Then the HTTP status code should be "204"
    And for user "Brian" the space "Shares" should not contain these entries:
      | <resource> |
    Examples:
      | permissions-role | resource      |
      | Viewer           | textfile.txt  |
      | File Editor      | textfile.txt  |
      | Viewer           | FolderToShare |
      | Editor           | FolderToShare |
      | Uploader         | FolderToShare |


  Scenario Outline: user removes user member from project space using permissions endpoint
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And user "Alice" has sent the following space share invitation:
      | space           | NewSpace           |
      | sharee          | Brian              |
      | shareType       | user               |
      | permissionsRole | <permissions-role> |
    When user "Alice" removes the access of user "Brian" from space "NewSpace" using permissions endpoint of the Graph API
    Then the HTTP status code should be "204"
    And the user "Brian" should not have a space called "NewSpace"
    Examples:
      | permissions-role |
      | Space Viewer     |
      | Space Editor     |
      | Manager          |


  Scenario Outline: user removes group from project space using permissions endpoint
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "NewSpace" with the default quota using the Graph API
    And group "group1" has been created
    And user "Brian" has been added to group "group1"
    And user "Alice" has sent the following space share invitation:
      | space           | NewSpace           |
      | sharee          | group1              |
      | shareType       | group               |
      | permissionsRole | <permissions-role> |
    When user "Alice" removes the access of group "group1" from space "NewSpace" using permissions endpoint of the Graph API
    Then the HTTP status code should be "204"
    And the user "Brian" should not have a space called "NewSpace"
    Examples:
      | permissions-role |
      | Space Viewer     |
      | Space Editor     |
      | Manager          |

  @env-config
  Scenario Outline: remove share after the share role Secure Viewer has been disabled (Personal Space)
    Given the administrator has enabled the permissions role "Secure Viewer"
    And user "Alice" has uploaded file with content "some content" to "textfile.txt"
    And user "Alice" has created folder "folderToShare"
    And user "Alice" has sent the following resource share invitation:
      | resource        | <resource>    |
      | space           | Personal      |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | Secure Viewer |
    And the administrator has disabled the permissions role "Secure Viewer"
    When user "Alice" removes the access of user "Brian" from resource "<resource>" of space "Personal" using the Graph API
    Then the HTTP status code should be "204"
    And for user "Brian" the space "Shares" should not contain these entries:
      | <resource> |
    Examples:
      | resource      |
      | textfile.txt  |
      | folderToShare |

  @env-config
  Scenario: remove share after the share role Denied has been disabled (Personal Space)
    Given the administrator has enabled the permissions role "Denied"
    And user "Alice" has created folder "folderToShare"
    And user "Alice" has sent the following resource share invitation:
      | resource        | folderToShare |
      | space           | Personal      |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | Denied        |
    And the administrator has disabled the permissions role "Denied"
    When user "Alice" removes the access of user "Brian" from resource "folderToShare" of space "Personal" using the Graph API
    Then the HTTP status code should be "204"
    And for user "Brian" the space "Shares" should not contain these entries:
      | folderToShare |

  @env-config
  Scenario Outline: remove share after the share role Secure Viewer has been disabled (Project Space)
    Given using spaces DAV path
    And the administrator has enabled the permissions role "Secure Viewer"
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "new-space" with content "some content" to "textfile.txt"
    And user "Alice" has created a folder "folderToShare" in space "new-space"
    And user "Alice" has sent the following resource share invitation:
      | resource        | <resource>    |
      | space           | new-space     |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | Secure Viewer |
    And the administrator has disabled the permissions role "Secure Viewer"
    When user "Alice" removes the access of user "Brian" from resource "<resource>" of space "new-space" using the Graph API
    Then the HTTP status code should be "204"
    And for user "Brian" the space "Shares" should not contain these entries:
      | <resource> |
    Examples:
      | resource      |
      | textfile.txt  |
      | folderToShare |

  @env-config
  Scenario: remove share after the share role Denied has been disabled (Project Space)
    Given using spaces DAV path
    And the administrator has enabled the permissions role "Denied"
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the Graph API
    And user "Alice" has created a folder "folderToShare" in space "new-space"
    And user "Alice" has sent the following resource share invitation:
      | resource        | folderToShare |
      | space           | new-space     |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | Denied        |
    And the administrator has disabled the permissions role "Denied"
    When user "Alice" removes the access of user "Brian" from resource "folderToShare" of space "new-space" using the Graph API
    Then the HTTP status code should be "204"
    And for user "Brian" the space "Shares" should not contain these entries:
      | folderToShare |

  @env-config
  Scenario Outline: remove share from group after the share role Secure Viewer has been disabled (Personal Space)
    Given the administrator has enabled the permissions role "Secure Viewer"
    And group "group1" has been created
    And user "Brian" has been added to group "group1"
    And user "Alice" has uploaded file with content "some content" to "textfile.txt"
    And user "Alice" has created folder "folderToShare"
    And user "Alice" has sent the following resource share invitation:
      | resource        | <resource>    |
      | space           | Personal      |
      | sharee          | group1        |
      | shareType       | group         |
      | permissionsRole | Secure Viewer |
    And the administrator has disabled the permissions role "Secure Viewer"
    When user "Alice" removes the access of group "group1" from resource "<resource>" of space "Personal" using the Graph API
    Then the HTTP status code should be "204"
    And for user "Brian" the space "Shares" should not contain these entries:
      | <resource> |
    Examples:
      | resource      |
      | textfile.txt  |
      | folderToShare |

  @env-config
  Scenario: remove share from group after the share role Denied has been disabled (Personal Space)
    Given the administrator has enabled the permissions role "Denied"
    And group "group1" has been created
    And user "Brian" has been added to group "group1"
    And user "Alice" has created folder "folderToShare"
    And user "Alice" has sent the following resource share invitation:
      | resource        | folderToShare |
      | space           | Personal      |
      | sharee          | group1        |
      | shareType       | group         |
      | permissionsRole | Denied        |
    And the administrator has disabled the permissions role "Denied"
    When user "Alice" removes the access of group "group1" from resource "folderToShare" of space "Personal" using the Graph API
    Then the HTTP status code should be "204"
    And for user "Brian" the space "Shares" should not contain these entries:
      | folderToShare |

  @env-config
  Scenario Outline: remove share from group after the share role Secure Viewer has been disabled (Project Space)
    Given using spaces DAV path
    And the administrator has enabled the permissions role "Secure Viewer"
    And group "group1" has been created
    And user "Brian" has been added to group "group1"
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "new-space" with content "some content" to "textfile.txt"
    And user "Alice" has created a folder "folderToShare" in space "new-space"
    And user "Alice" has sent the following resource share invitation:
      | resource        | <resource>    |
      | space           | new-space     |
      | sharee          | group1        |
      | shareType       | group         |
      | permissionsRole | Secure Viewer |
    And the administrator has disabled the permissions role "Secure Viewer"
    When user "Alice" removes the access of group "group1" from resource "<resource>" of space "new-space" using the Graph API
    Then the HTTP status code should be "204"
    And for user "Brian" the space "Shares" should not contain these entries:
      | <resource> |
    Examples:
      | resource      |
      | textfile.txt  |
      | folderToShare |

  @env-config
  Scenario: remove share from group after the share role Denied has been disabled (Project Space)
    Given using spaces DAV path
    And the administrator has enabled the permissions role "Denied"
    And group "group1" has been created
    And user "Brian" has been added to group "group1"
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the Graph API
    And user "Alice" has created a folder "folderToShare" in space "new-space"
    And user "Alice" has sent the following resource share invitation:
      | resource        | folderToShare |
      | space           | new-space     |
      | sharee          | group1        |
      | shareType       | group         |
      | permissionsRole | Denied        |
    And the administrator has disabled the permissions role "Denied"
    When user "Alice" removes the access of group "group1" from resource "folderToShare" of space "new-space" using the Graph API
    Then the HTTP status code should be "204"
    And for user "Brian" the space "Shares" should not contain these entries:
      | folderToShare |
