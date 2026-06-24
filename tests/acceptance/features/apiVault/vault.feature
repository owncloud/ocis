Feature: vault

  Scenario:
    Given these users have been created with default attributes:
      | username |
      | Alice    |
      | Brian    |
#    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
#    And using spaces DAV path
    When user "Alice" creates folder "test-folder" using the WebDAV API
#    When user "Alice" creates a folder "vaultFolder" in vault space "Personal" using the WebDav Api
