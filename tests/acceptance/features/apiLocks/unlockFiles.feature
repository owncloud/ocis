Feature: unlock locked items
  As a user
  I want to unlock the resources previously locked by myself
  So that other users can make changes to the resources

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files

  @issue-7696
  Scenario Outline: unlock a locked file in project space
    Given using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "project-space" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "project-space" with content "some data" to "textfile.txt"
    And user "Alice" has locked file "textfile.txt" inside space "project-space" setting the following properties
      | lockscope | <lock-scope> |
    When user "Alice" unlocks the last created lock of file "textfile.txt" inside space "project-space" using the WebDAV API
    Then the HTTP status code should be "204"
    Examples:
      | lock-scope |
      | shared     |
      | exclusive  |
