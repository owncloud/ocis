@env-config @cli-stale-uploads
Feature: stale upload via CLI command
  As an administrator
  I want to manage stale uploads
  So that I clean up stale uploads from storage

  Background:
    Given user "Alice" has been created with default attributes
    And user "Brian" has been created with default attributes


  Scenario: list and delete all stale uploads
    Given the config "POSTPROCESSING_DELAY" has been set to "10s" for "postprocessing" service
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "staleuploads" with the default quota using the Graph API
    And user "Alice" has uploaded a file "filesForUpload/testavatar.jpg" to "/testavatar.jpg" in space "staleuploads"
    And user "Brian" has uploaded file "filesForUpload/textfile.txt" to "textfile.txt"
    And the administrator has stopped the server
    And the administrator has created stale upload
    And the administrator has started the server
    When the administrator lists all the stale uploads
    Then the command should be successful
    And the CLI response should contain the following message:
      """
      Scanning all spaces for stale processing nodes...
      Total stale nodes: 2
      """
    When the administrator deletes all the stale uploads
    Then the command should be successful
    And the CLI response should contain the following message:
      """
      Scanning all spaces for stale processing nodes...
      Total stale nodes: 2
      """
    And there should be "0" stale uploads


  Scenario: list and delete all stale uploads of a specific space
    Given user "Alice" has created folder "FolderToShare"
    And using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "staleuploads" with the default quota using the Graph API
    And using SharingNG
    And the config "POSTPROCESSING_DELAY" has been set to "10s" for "postprocessing" service
    And user "Alice" has sent the following resource share invitation:
      | resource        | FolderToShare |
      | space           | Personal      |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | Uploader      |
    And user "Brian" has uploaded a file "filesForUpload/testavatar.png" to "FolderToShare/testavatar.png" in space "Shares"
    And user "Alice" has uploaded file "filesForUpload/testavatar.png" to "testavatar.png"
    And user "Alice" has uploaded a file "filesForUpload/testavatar.jpg" to "/testavatar.jpg" in space "staleuploads"
    And the administrator has stopped the server
    And the administrator has created stale upload
    And the administrator has started the server
    When the administrator lists all the stale uploads of space "Personal" owned by user "Alice"
    Then the command should be successful
    And the CLI response should contain the following message:
      """
      Total stale nodes: 2
      """
    When the administrator lists all the stale uploads of space "staleuploads" owned by user "Alice"
    Then the command should be successful
    And the CLI response should contain the following message:
      """
      Total stale nodes: 1
      """
    When the administrator deletes all the stale uploads of space "Personal" owned by user "Alice"
    Then the command should be successful
    And the CLI response should contain the following message:
      """
      Total stale nodes: 2
      """
    And there should be "0" stale uploads of space "Personal" owned by user "Alice"
    When the administrator deletes all the stale uploads of space "staleuploads" owned by user "Alice"
    Then the command should be successful
    And the CLI response should contain the following message:
      """
      Total stale nodes: 1
      """
    And there should be "0" stale uploads of space "Personal" owned by user "Alice"
