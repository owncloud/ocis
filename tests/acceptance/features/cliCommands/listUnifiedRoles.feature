@env-config
Feature: List unified roles
  As an administrator
  I want to list all available unified roles
  So that I can verify which roles exist, their permissions, and their status

  Scenario: List unified roles with expected fields
    And the following headers should not be set
      | header                        |
      | Access-Control-Allow-Headers  |
      | Access-Control-Expose-Headers |
      | Access-Control-Allow-Origin   |
      | Access-Control-Allow-Methods  |

    When the administrator lists unified roles using the CLI
    Then the command should be successful
    And the command output should contain "💚 No inconsistency found. The backup in '%storage_path%' seems to be valid."
    And the command output should include the following fields for each role:
      | Label                  |
      | uid                    |
      | Description            |
      | Enabled                |
      | Condition              |
      | Allowed Resource Action |

