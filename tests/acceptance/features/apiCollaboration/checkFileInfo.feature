Feature: check file info with different wopi apps
  As a user
  I want to request file information on wopi server
  So that I can get all the information of a file

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files


  Scenario: check file info with fakeOffice
    Given user "Alice" has uploaded file with content "hello world" to "/textfile0.txt"
    When user "Alice" checks the information of file "textfile0.txt" of space "Personal" using office "FakeOffice"
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "BaseFileName",
          "OwnerId",
          "Size",
          "UserId",
          "Version",
          "SupportsCobalt",
          "SupportsContainers",
          "SupportsDeleteFile",
          "SupportsEcosystem",
          "SupportsExtendedLockLength",
          "SupportsFolders",
          "SupportsGetLock",
          "SupportsLocks",
          "SupportsRename",
          "SupportsUpdate",
          "SupportsUserInfo",
          "UserFriendlyName",
          "ReadOnly",
          "RestrictedWebViewOnly",
          "UserCanAttend",
          "UserCanNotWriteRelative",
          "UserCanPresent",
          "UserCanRename",
          "UserCanWrite",
          "AllowAdditionalMicrosoftServices",
          "AllowExternalMarketplace",
          "DisablePrint",
          "DisableTranslation",
          "BreadcrumbDocName"
        ],
        "properties": {
          "BaseFileName": {
            "const": "textfile0.txt"
          },
          "OwnerId": {
            "type": "string"
          },
          "Size": {
            "const": 11
          },
          "UserId": {
            "type": "string"
          },
          "Version": {
            "type": "string"
          },
          "SupportsCobalt": {
            "const": false
          },
          "SupportsContainers": {
            "const": false
          },
          "SupportsDeleteFile": {
            "const": true
          },
          "SupportsEcosystem": {
            "const": false
          },
          "SupportsExtendedLockLength": {
            "const": true
          },
          "SupportsFolders": {
            "const": false
          },
          "SupportsGetLock": {
            "const": true
          },
          "SupportsLocks": {
            "const": true
          },
          "SupportsRename": {
            "const": true
          },
          "SupportsUpdate": {
            "const": true
          },
          "SupportsUserInfo": {
            "const": false
          },
          "UserFriendlyName": {
            "const": "Alice Hansen"
          },
          "ReadOnly": {
            "const": false
          },
          "RestrictedWebViewOnly": {
            "const": false
          },
          "UserCanAttend": {
            "const": false
          },
          "UserCanNotWriteRelative": {
            "const": false
          },
          "UserCanPresent": {
            "const": false
          },
          "UserCanRename": {
            "const": true
          },
          "UserCanWrite": {
            "const": true
          },
          "AllowAdditionalMicrosoftServices": {
            "const": false
          },
          "AllowExternalMarketplace": {
            "const": false
          },
          "DisablePrint": {
            "const": false
          },
          "DisableTranslation": {
            "const": false
          },
          "BreadcrumbDocName": {
            "const": "textfile0.txt"
          }
        }
      }
      """


  Scenario: check file info with collabora
    Given user "Alice" has uploaded file with content "hello world" to "/textfile0.txt"
    When user "Alice" checks the information of file "textfile0.txt" of space "Personal" using office "Collabora"
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "BaseFileName",
          "DisablePrint",
          "OwnerId",
          "PostMessageOrigin",
          "Size",
          "UserCanWrite",
          "UserCanNotWriteRelative",
          "UserId",
          "UserFriendlyName",
          "EnableOwnerTermination",
          "SupportsLocks",
          "SupportsRename",
          "UserCanRename",
          "BreadcrumbDocName"
        ],
        "properties": {
          "BaseFileName": {
            "const": "textfile0.txt"
          },
          "PostMessageOrigin": {
            "const": "https://localhost:9200"
          },
          "DisablePrint": {
            "const": false
          },
          "OwnerId": {
            "type": "string"
          },
          "Size": {
            "const": 11
          },
          "UserCanWrite": {
            "const": true
          },
          "UserCanNotWriteRelative": {
            "const": false
          },
          "EnableOwnerTermination": {
            "const": true
          },
          "UserId": {
            "type": "string"
          },
          "SupportsLocks": {
            "const": true
          },
          "SupportsRename": {
            "const": true
          },
          "UserFriendlyName": {
            "const": "Alice Hansen"
          },
          "UserCanRename": {
            "const": true
          },
          "BreadcrumbDocName": {
            "const": "textfile0.txt"
          }
        }
      }
      """


  Scenario: check file info with onlyOffice
    Given user "Alice" has uploaded file with content "hello world" to "/textfile0.txt"
    When user "Alice" checks the information of file "textfile0.txt" of space "Personal" using office "OnlyOffice"
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "BaseFileName",
          "Version",
          "BreadcrumbDocName",
          "BreadcrumbFolderName",
          "BreadcrumbFolderUrl",
          "PostMessageOrigin",
          "DisablePrint",
          "UserFriendlyName",
          "UserId",
          "ReadOnly",
          "UserCanNotWriteRelative",
          "UserCanRename",
          "UserCanWrite",
          "SupportsLocks",
          "SupportsRename",
          "SupportsUpdate"
        ],
        "properties": {
          "BaseFileName": {
            "const": "textfile0.txt"
          },
          "UserId": {
            "type": "string"
          },
          "Version": {
            "type": "string"
          },
          "SupportsLocks": {
            "const": true
          },
          "SupportsRename": {
            "const": true
          },
          "SupportsUpdate": {
            "const": true
          },
          "UserFriendlyName": {
            "const": "Alice Hansen"
          },
          "ReadOnly": {
            "const": false
          },
          "UserCanNotWriteRelative": {
            "const": false
          },
          "UserCanRename": {
            "const": true
          },
          "UserCanWrite": {
            "const": true
          },
          "DisablePrint": {
            "const": false
          },
          "BreadcrumbDocName": {
            "const": "textfile0.txt"
          },
          "BreadcrumbFolderName": {
            "const": "Alice Hansen"
          },
          "BreadcrumbFolderUrl": {
            "type": "string"
          },
          "PostMessageOrigin": {
            "type": "string"
          }
        }
      }
      """
