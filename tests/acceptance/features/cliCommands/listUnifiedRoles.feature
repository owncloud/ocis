@env-config
Feature: List unified roles
  As an administrator
  I want to list all available unified roles
  So that I can check which roles exist, their permissions, and their status

  @issue-11254
  Scenario: List unified roles with expected fields
    When the administrator lists all the unified roles using the CLI
    Then the command should be successful
    And the command output should include the following roles:
      | LABEL                            | ENABLED  | DESCRIPTION                                                                          |
      | Viewer                           | enabled  | View and download.                                                                   |
      | ViewerListGrants                 | disabled | View, download and show all invited people.                                          |
      | SpaceViewer                      | enabled  | View and download.                                                                   |
      | Editor                           | enabled  | View, download, upload, edit, add and delete.                                        |
      | EditorListGrants                 | disabled | View, download, upload, edit, add, delete and show all invited people.               |
      | EditorListGrantsWithVersions     | disabled | View, download, upload, edit, delete, show all versions and all invited people.      |
      | SpaceEditor                      | enabled  | View, download, upload, edit, add, show all versions and delete.                     |
      | SpaceEditorWithoutVersions       | disabled | View, download, upload, edit and add.                                                |
      | SpaceEditorWithoutTrashbin       | disabled | View, download, upload, edit, add and show all versions.                             |
      | FileEditor                       | enabled  | View, download, upload and edit.                                                     |
      | FileEditorListGrants             | disabled | View, download, upload, edit and show all invited people.                            |
      | FileEditorListGrantsWithVersions | disabled | View, download, upload, edit, show all versions and all invited people.              |
      | EditorLite                       | enabled  | View, download, upload, edit and add.                                                |
      | Manager                          | enabled  | View, download, upload, edit, add, show all versions, delete and manage members.     |
      | SecureViewer                     | disabled | View only documents, images and PDFs. Watermarks will be applied.                    |
      | Denied                           | disabled | Deny all access.                                                                     |
