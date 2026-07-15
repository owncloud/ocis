Feature: vault
  As a user
  I want to store resource in vault storage
  So that vault resources are isolated with regular drive storage

  Background:
    Given using spaces DAV path
    And these users have been created with default attributes:
      | username |
      | Alice    |


  Scenario: user can create folders and files in personal space in vault
    Given user "Alice" has logged in via web UI
    When user "Alice" creates a folder "vaultFolder" in space "Personal" in vault using the WebDav Api
    Then the HTTP status code should be "201"
    When user "Alice" uploads a file inside space "Personal" with content "some content" to "vaultFile.txt" in vault using the WebDAV API
    Then the HTTP status code should be "201"
    And for user "Alice" the space "Personal" in vault should contain these entries:
      | vaultFolder   |
      | vaultFile.txt |


  Scenario: user can create folders and files in project space in vault
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has logged in via web UI
    And user "Alice" has created a space "vault-space" in vault with the default quota using the Graph API
    When user "Alice" creates a folder "vaultFolder" in space "vault-space" in vault using the WebDav Api
    Then the HTTP status code should be "201"
    When user "Alice" uploads a file inside space "vault-space" with content "some content" to "vaultFile.txt" in vault using the WebDAV API
    Then the HTTP status code should be "201"
    And for user "Alice" the space "vault-space" in vault should contain these entries:
      | vaultFolder   |
      | vaultFile.txt |


  Scenario: resources in drive and vault are isolated
    Given user "Alice" has logged in via web UI
    And user "Alice" has created a folder "driveFolder" in space "Personal"
    And user "Alice" has uploaded a file inside space "Personal" with content "some content" to "driveFile.txt"
    And user "Alice" has created a folder "vaultFolder" in space "Personal" in vault
    When user "Alice" uploads a file inside space "Personal" with content "some content" to "vaultFile.txt" in vault using the WebDAV API
    Then the HTTP status code should be "201"
    And for user "Alice" the space "Personal" in vault should contain these entries:
      | vaultFolder   |
      | vaultFile.txt |
    And for user "Alice" the space "Personal" should contain these entries:
      | driveFolder   |
      | driveFile.txt |
    And for user "Alice" the space "Personal" in vault should not contain these entries:
      | driveFolder   |
      | driveFile.txt |
    And for user "Alice" the space "Personal" should not contain these entries:
      | vaultFolder   |
      | vaultFile.txt |

  @env-config @keycloak-config
  Scenario: user can set custom auth level names
    Given the administrator has set the Keycloak realm attribute "acr.loa.map" to '{"regular":"1","testing":"2"}'
    And the config "OCIS_MFA_AUTH_LEVEL_NAMES" has been set to "testing"
    And user "Alice" has logged in via web UI
    When user "Alice" uploads a file inside space "Personal" with content "some content" to "vaultFile.txt" in vault using the WebDAV API
    Then the HTTP status code should be "201"
    And user "Alice" should have acr value "testing"
