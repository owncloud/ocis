@env-config
Feature: List unified roles
  As an administrator
  I want to list all available unified roles
  So that I can check which roles exist, their permissions, and their status

  @issue-11254
  Scenario: List unified roles with expected fields
    When the administrator lists all the unified roles using the CLI
    Then the command should be successful
    And the command output should include the following fields for each role:
      | field                            |
      | LABEL                            |
      | UID                              |
      | ENABLED                          |
      | DESCRIPTION                      |
      | CONDITION                        |
      | ALLOWED RESOURCE ACTIONS         |
      | CONDITION                        |
      | Viewer                           |
      | ViewerListGrants                 |
      | SpaceViewer                      |
      | Editor                           |
      | EditorListGrants                 |
      | EditorListGrantsWithVersions     |
      | SpaceEditor                      |
      | SpaceEditorWithoutVersions       |
      | FileEditor                       |
      | FileEditorListGrants             |
      | FileEditorListGrantsWithVersions |
      | EditorLite                       |
      | Manager                          |
      | SecureViewer                     |
      | Denied                           |
