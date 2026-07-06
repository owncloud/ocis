@env-config
Feature: blobstore check and get via CLI
  As an administrator
  I want to verify blobstore connectivity
  So that I can ensure my storage backend is working correctly


  Scenario: administrator checks the blobstore with default blob size
    When the administrator checks the blobstore using the CLI
    Then the command should be successful
    And the command output should contain "Upload: OK"
    And the command output should contain "Download and verify: OK"
    And the command output should contain "Delete: OK"
    And the command output should contain "Blobstore check successful."


  Scenario: administrator checks the blobstore with a custom blob size
    When the administrator checks the blobstore with blob size "1KB" using the CLI
    Then the command should be successful
    And the command output should contain "Blobstore check successful."


  Scenario: administrator checks the blobstore with an invalid blob size
    When the administrator checks the blobstore with blob size "abc" using the CLI
    Then the command should be unsuccessful
    And the command output should contain "invalid --blob-size"


  Scenario: administrator tries to get a blob without providing required flags
    When the administrator gets a blob from the blobstore using the CLI
    Then the command should be unsuccessful
    And the command output should contain "either --path or both --blob-id and --space-id must be set"


  Scenario: administrator tries to get a non-existent blob from the blobstore
    When the administrator gets a non-existent blob from the blobstore using the CLI
    Then the command should be unsuccessful
    And the command output should contain "download failed: could not read blob"
