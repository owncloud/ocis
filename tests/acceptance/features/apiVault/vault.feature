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


  Scenario: check capabilities endpoint for vault
    Given using OCS API version "2"
    And user "Alice" has logged in via web UI
    When user "Alice" retrieves the vault mode capabilities using the capabilities API
    Then the OCS status code should be "200"
    And the HTTP status code should be "200"
    And the ocs JSON data of the response should match
      """
      {
        "type": "object",
        "required": [ "capabilities" ],
        "properties": {
          "capabilities": {
            "type": "object",
            "required": [
              "core",
              "files",
              "files_sharing",
              "auth",
              "vault"
            ],
            "properties": {
              "files_sharing": {
                "type": "object",
                "required": [
                  "api_enabled",
                  "default_permissions",
                  "public",
                  "resharing",
                  "federation",
                  "group_sharing",
                  "share_with_group_members_only",
                  "share_with_membership_groups_only",
                  "auto_accept_share",
                  "user_enumeration"
                ],
                "properties": {
                  "federation": {
                    "type": "object",
                    "required": [
                      "outgoing",
                      "incoming"
                    ],
                    "properties": {
                      "outgoing": {
                        "const": false
                      },
                      "incoming": {
                        "const": false
                      }
                    }
                  },
                  "public": {
                    "type": "object",
                    "required": [
                      "enabled",
                      "multiple",
                      "upload",
                      "supports_upload_only",
                      "send_mail",
                      "social_share"
                    ],
                    "properties": {
                      "enabled": {
                        "const": false
                      }
                    }
                  }
                }
              },
              "auth": {
                "type": "object",
                "required": [
                  "mfa"
                ],
                "properties": {
                  "mfa": {
                    "type": "object",
                    "required": [
                      "enabled",
                      "levelnames"
                    ],
                    "properties": {
                      "enabled": {
                        "const": true
                      },
                      "levelnames": {
                        "type": "array",
                        "minItems": 1,
                        "maxItems": 1,
                        "items": {
                          "const": "advanced"
                        }
                      }
                    }
                  }
                }
              },
              "vault": {
                "type": "object",
                "required": [
                  "enabled",
                  "vault_storage_provider"
                ],
                "properties": {
                  "enabled": {
                    "const": true
                  },
                  "vault_storage_provider": {
                    "pattern": "%uuidv4_pattern%"
                  }
                }
              }
            }
          }
        }
      }
      """


  Scenario: user copies folder from drive to vault
    Given user "Alice" has logged in via web UI
    And user "Alice" has created a folder "driveFolder" in space "Personal"
    When user "Alice" copies folder "driveFolder" from space "Personal" to "driveFolder" inside space "Personal" in vault using the WebDAV API
    Then the HTTP status code should be "201"
    And for user "Alice" the space "Personal" in vault should contain these entries:
      | driveFolder |
    And for user "Alice" the space "Personal" should contain these entries:
      | driveFolder |


  Scenario: user copies file from drive to vault
    Given user "Alice" has logged in via web UI
    And user "Alice" has uploaded a file inside space "Personal" with content "some content" to "testfile.txt"
    When user "Alice" copies file "testfile.txt" from space "Personal" to "testfile.txt" inside space "Personal" in vault using the WebDAV API
    Then the HTTP status code should be "201"
    And for user "Alice" the space "Personal" in vault should contain these entries:
      | testfile.txt |
    And for user "Alice" the content of the file "testfile.txt" of the space "Personal" in vault should be "some content"
    And for user "Alice" the space "Personal" should contain these entries:
      | testfile.txt |


  Scenario: user tries to copy folder from vault to drive
    Given user "Alice" has logged in via web UI
    And user "Alice" has created a folder "vaultFolder" in space "Personal" in vault
    When user "Alice" copies folder "vaultFolder" from space "Personal" in vault to "vaultFolder" inside space "Personal" using the WebDAV API
    Then the HTTP status code should be "409"
    And for user "Alice" the space "Personal" should not contain these entries:
      | vaultFolder |
    And for user "Alice" the space "Personal" in vault should contain these entries:
      | vaultFolder |


  Scenario: user tries to copy file from vault to drive
    Given user "Alice" has logged in via web UI
    And user "Alice" has uploaded a file inside space "Personal" with content "some content" to "testfile.txt" in vault
    When user "Alice" copies file "testfile.txt" from space "Personal" in vault to "testfile.txt" inside space "Personal" using the WebDAV API
    Then the HTTP status code should be "409"
    And for user "Alice" the space "Personal" should not contain these entries:
      | testfile.txt |
    And for user "Alice" the space "Personal" in vault should contain these entries:
      | testfile.txt |


  Scenario: user copies sub-folder from drive to vault
    Given user "Alice" has logged in via web UI
    And user "Alice" has created a folder "driveFolder" in space "Personal"
    And user "Alice" has created a folder "driveFolder/subFolder" in space "Personal"
    When user "Alice" copies folder "driveFolder/subFolder" from space "Personal" to "subFolder" inside space "Personal" in vault using the WebDAV API
    Then the HTTP status code should be "201"
    And for user "Alice" the space "Personal" in vault should contain these entries:
      | subFolder |
    And for user "Alice" folder "driveFolder" of the space "Personal" should contain these entries:
      | subFolder |


  Scenario: user copies file inside folder from drive to vault
    Given user "Alice" has logged in via web UI
    And user "Alice" has created a folder "driveFolder" in space "Personal"
    And user "Alice" has uploaded a file inside space "Personal" with content "some content" to "driveFolder/testfile.txt"
    When user "Alice" copies file "driveFolder/testfile.txt" from space "Personal" to "testfile.txt" inside space "Personal" in vault using the WebDAV API
    Then the HTTP status code should be "201"
    And for user "Alice" the space "Personal" in vault should contain these entries:
      | testfile.txt |
    And for user "Alice" the content of the file "testfile.txt" of the space "Personal" in vault should be "some content"
    And for user "Alice" folder "driveFolder" of the space "Personal" should contain these entries:
      | testfile.txt |


  Scenario: user copies sub-folder from drive to a folder in vault
    Given user "Alice" has logged in via web UI
    And user "Alice" has created a folder "driveFolder" in space "Personal"
    And user "Alice" has created a folder "driveFolder/subFolder" in space "Personal"
    And user "Alice" has created a folder "vaultFolder" in space "Personal" in vault
    When user "Alice" copies folder "driveFolder/subFolder" from space "Personal" to "vaultFolder/subFolder" inside space "Personal" in vault using the WebDAV API
    Then the HTTP status code should be "201"
    And for user "Alice" folder "vaultFolder" of the space "Personal" in vault should contain these entries:
      | subFolder |
    And for user "Alice" folder "driveFolder" of the space "Personal" should contain these entries:
      | subFolder |


  Scenario: user copies file inside folder from drive to a folder in vault
    Given user "Alice" has logged in via web UI
    And user "Alice" has created a folder "driveFolder" in space "Personal"
    And user "Alice" has uploaded a file inside space "Personal" with content "some content" to "driveFolder/testfile.txt"
    And user "Alice" has created a folder "vaultFolder" in space "Personal" in vault
    When user "Alice" copies file "driveFolder/testfile.txt" from space "Personal" to "vaultFolder/testfile.txt" inside space "Personal" in vault using the WebDAV API
    Then the HTTP status code should be "201"
    And for user "Alice" folder "vaultFolder" of the space "Personal" in vault should contain these entries:
      | testfile.txt |
    And for user "Alice" the content of the file "vaultFolder/testfile.txt" of the space "Personal" in vault should be "some content"
    And for user "Alice" folder "driveFolder" of the space "Personal" should contain these entries:
      | testfile.txt |
