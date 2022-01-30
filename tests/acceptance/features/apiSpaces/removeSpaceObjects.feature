@api @skipOnOcV10
Feature: Remove files, folder
  As a user
  I want to be able to remove files, folders
  Users with the editor role can also remove objects
  Users with the viewer role cannot remove objects 

  Note - this feature is run in CI with ACCOUNTS_HASH_DIFFICULTY set to the default for production
  See https://github.com/owncloud/ocis/issues/1542 and https://github.com/owncloud/ocis/pull/839

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files
    And user "Brian" has been created with default attributes and without skeleton files
    And the administrator has given "Alice" the role "Admin" using the settings api

#   owner of space (admin permissions)

  Scenario: An owner can delete a folder with some subfolders in a Space via the webDav API
    Given user "Alice" has created a space "Owner deletes folder" of type "project" with quota "10"
    And user "Alice" has created a folder "folderForDeleting/sub1/sub2" in space "Owner deletes folder"
    When user "Alice" removes the object "folderForDeleting" from space "Owner deletes folder"
    Then the HTTP status code should be "204"
    And for user "Alice" the space "Owner deletes folder" should not contain these entries:
      | folderForDeleting |


  Scenario: An owner can delete a subfolder in a Space via the webDav API
    Given user "Alice" has created a space "Owner deletes subfolder" of type "project" with quota "10"
    And user "Alice" has created a subfolder "folder/subFolderForDeleting" in space "Owner deletes subfolder"
    When user "Alice" removes the object "folder/subFolderForDeleting" from space "Owner deletes subfolder"
    Then the HTTP status code should be "204"
    And for user "Alice" the space "Owner deletes subfolder" should contain these entries:
      | folder |
    And for user "Alice" folder "folder/" of the space "Owner deletes subfolder" should not contain these entries:
      | subFolderForDeleting |


  Scenario: An owner can delete a file in a Space via the webDav API
    Given user "Alice" has created a space "Owner deletes file" of type "project" with quota "20"
    And user "Alice" has uploaded a file inside space "Owner deletes file" with content "some content" to "text.txt"
    When user "Alice" removes the object "text.txt" from space "Owner deletes file"
    Then the HTTP status code should be "204"
    And for user "Alice" the space "Owner deletes file" should not contain these entries:
      | text.txt |
    When user "Alice" lists all available spaces via the GraphApi
    Then the json responded should contain a space "Owner deletes file" with these key and value pairs:
      | key          | value             |
      | name         | Owner deletes file |
      | quota@@@used | 0                 |

#   editor role

  Scenario: An editor can delete a folder with some subfolders in a Space via the webDav API
    Given user "Alice" has created a space "Editor deletes folder" of type "project" with quota "10"
    And user "Alice" has created a folder "folderForDeleting/sub1/sub2" in space "Editor deletes folder"
    And user "Alice" has shared a space "Editor deletes folder" to user "Brian" with role "editor"
    When user "Brian" removes the object "folderForDeleting" from space "Editor deletes folder"
    Then the HTTP status code should be "204"
    And for user "Brian" the space "Editor deletes folder" should not contain these entries:
      | folderForDeleting |


  Scenario: An editor can delete a subfolder in a Space via the webDav API
    Given user "Alice" has created a space "Editor deletes subfolder" of type "project" with quota "10"
    And user "Alice" has created a subfolder "folder/subFolderForDeleting" in space "Editor deletes subfolder"
    And user "Alice" has shared a space "Editor deletes subfolder" to user "Brian" with role "editor"
    When user "Brian" removes the object "folder/subFolderForDeleting" from space "Editor deletes subfolder"
    Then the HTTP status code should be "204"
    And for user "Brian" the space "Editor deletes subfolder" should contain these entries:
      | folder |
    And for user "Brian" folder "folder/" of the space "Editor deletes subfolder" should not contain these entries:
      | subFolderForDeleting |


  Scenario: An editor can delete a file in a Space via the webDav API
    Given user "Alice" has created a space "Editor deletes file" of type "project" with quota "20"
    And user "Alice" has uploaded a file inside space "Editor deletes file" with content "some content" to "text.txt"
    And user "Alice" has shared a space "Editor deletes file" to user "Brian" with role "editor"
    When user "Brian" removes the object "text.txt" from space "Editor deletes file"
    Then the HTTP status code should be "204"
    And for user "Brian" the space "Editor deletes file" should not contain these entries:
      | text.txt |
    When user "Brian" lists all available spaces via the GraphApi
    Then the json responded should contain a space "Editor deletes file" with these key and value pairs:
      | key          | value              |
      | name         | Editor deletes file |
      | quota@@@used | 0                  |

#   viewer role 

  Scenario: A viewer cannot delete a folder with some subfolders in a Space via the webDav API
    Given user "Alice" has created a space "Viewer deletes folder" of type "project" with quota "10"
    And user "Alice" has created a folder "folderForDeleting/sub1/sub2" in space "Viewer deletes folder"
    And user "Alice" has shared a space "Viewer deletes folder" to user "Brian" with role "viewer"
    When user "Brian" removes the object "folderForDeleting" from space "Viewer deletes folder"
    Then the HTTP status code should be "403"
    And for user "Brian" the space "Viewer deletes folder" should contain these entries:
      | folderForDeleting |


  Scenario: A viewer cannot delete a subfolder in a Space via the webDav API
    Given user "Alice" has created a space "Viewer deletes subfolder" of type "project" with quota "10"
    And user "Alice" has created a subfolder "folder/subFolderForDeleting" in space "Viewer deletes subfolder"
    And user "Alice" has shared a space "Viewer deletes subfolder" to user "Brian" with role "viewer"
    When user "Brian" removes the object "folder/subFolderForDeleting" from space "Viewer deletes subfolder"
    Then the HTTP status code should be "403"
    And for user "Brian" folder "folder/" of the space "Viewer deletes subfolder" should contain these entries:
      | subFolderForDeleting |


  Scenario: A viewer cannot delete a file in a Space via the webDav API
    Given user "Alice" has created a space "Viewer deletes file" of type "project" with quota "20"
    And user "Alice" has uploaded a file inside space "Viewer deletes file" with content "some content" to "text.txt"
    And user "Alice" has shared a space "Viewer deletes file" to user "Brian" with role "viewer"
    When user "Brian" removes the object "text.txt" from space "Viewer deletes file"
    Then the HTTP status code should be "403"
    And for user "Brian" the space "Viewer deletes file" should contain these entries:
      | text.txt |

  
  Scenario: An user is unable to delete a Space via the webDav API
    Given user "Alice" has created a space "user deletes a space" of type "project" with quota "20"
    When user "Alice" removes the object "" from space "user deletes a space"
    Then the HTTP status code should be "405"
    When user "Alice" lists all available spaces via the GraphApi
    Then the json responded should contain a space "user deletes a space" with these key and value pairs:
      | key          | value                |
      | name         | user deletes a space |
