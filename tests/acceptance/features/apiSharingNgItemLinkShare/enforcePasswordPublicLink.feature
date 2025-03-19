@env-config
Feature: enforce password on shareNg public link
  As a user
  I want to enforce passwords on public links shared with upload, edit, or contribute permission
  So that the password is required to access the contents of the link

  Password requirements. set by default:
  | OCIS_SHARING_PUBLIC_SHARE_MUST_HAVE_PASSWORD  | true |
  | OCIS_PASSWORD_POLICY_MIN_CHARACTERS           | 8    |
  | OCIS_PASSWORD_POLICY_MIN_LOWERCASE_CHARACTERS | 1    |
  | OCIS_PASSWORD_POLICY_MIN_UPPERCASE_CHARACTERS | 1    |
  | OCIS_PASSWORD_POLICY_MIN_DIGITS               | 1    |
  | OCIS_PASSWORD_POLICY_MIN_SPECIAL_CHARACTERS   | 1    |


  Scenario Outline: user should be able to create public link with permission view and internal without a password when enforce-password is enabled
    Given the following configs have been set:
      | config                                                 | value |
      | OCIS_SHARING_PUBLIC_SHARE_MUST_HAVE_PASSWORD           | false |
      | OCIS_SHARING_PUBLIC_WRITEABLE_SHARE_MUST_HAVE_PASSWORD | true  |
    And user "Alice" has been created with default attributes
    And user "Alice" has created folder "folder"
    When user "Alice" creates the following resource link share using the Graph API:
      | resource        | folder             |
      | space           | Personal           |
      | permissionsRole | <permissions-role> |
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["hasPassword","id","link"],
        "properties": {
          "hasPassword": {"const": false},
          "id": {"pattern": "^[a-zA-Z]{15}$"},
          "link": {
            "type": "object",
            "required": ["@libre.graph.displayName","@libre.graph.quickLink","preventsDownload","type","webUrl"],
            "properties": {
              "@libre.graph.displayName": {"const": ""},
              "@libre.graph.quickLink": {"const": false},
              "preventsDownload": {"const": false},
              "type": {"const": "<permissions-role-value>"},
              "webUrl": {"pattern": "^%base_url%/s/[a-zA-Z]{15}$"}
            }
          }
        }
      }
      """
    Examples:
      | permissions-role | permissions-role-value |
      | View             | view                   |
      | Internal         | internal               |


  Scenario Outline: user shouldn't be able to create public link with permission edit, upload, file drop and secure viewer without a password when enforce-password is enabled
    Given the following configs have been set:
      | config                                                 | value |
      | OCIS_SHARING_PUBLIC_SHARE_MUST_HAVE_PASSWORD           | false |
      | OCIS_SHARING_PUBLIC_WRITEABLE_SHARE_MUST_HAVE_PASSWORD | true  |
    And user "Alice" has been created with default attributes
    And user "Alice" has created folder "folder"
    When user "Alice" creates the following resource link share using the Graph API:
      | resource        | folder             |
      | space           | Personal           |
      | permissionsRole | <permissions-role> |
    Then the HTTP status code should be "400"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["error"],
        "properties": {
          "error": {
            "type": "object",
            "required": ["code","innererror","message"],
            "properties": {
              "code": {"const": "invalidRequest"},
              "innererror": {
                "type": "object",
                "required": ["date","request-id"]
              },
              "message": {"const": "password protection is enforced"}
            }
          }
        }
      }
      """
    Examples:
      | permissions-role |
      | Edit             |
      | Upload           |
      | File Drop        |


  Scenario Outline: user should allowed to updated public link to edit, upload file drop and secure viewer permission when enforce-password is enabled
    Given the following configs have been set:
      | config                                                 | value |
      | OCIS_SHARING_PUBLIC_SHARE_MUST_HAVE_PASSWORD           | false |
      | OCIS_SHARING_PUBLIC_WRITEABLE_SHARE_MUST_HAVE_PASSWORD | true  |
    And using spaces DAV path
    And user "Alice" has been created with default attributes
    And user "Alice" has created folder "folder"
    And user "Alice" has uploaded file with content "test file" to "folder/testfile.txt"
    And using SharingNG
    And user "Alice" has created the following resource link share:
      | resource        | folder   |
      | space           | Personal |
      | permissionsRole | View     |
      | password        | %public% |
    And user "Alice" updates the last public link share using the permissions endpoint of the Graph API:
      | resource        | folder                 |
      | space           | Personal               |
      | permissionsRole | <new-permissions-role> |
    Then the HTTP status code should be "200"
    And the public should not be able to download file "/testfile.txt" from inside the last public link shared folder using the public WebDAV API without a password
    And the public should not be able to download file "/testfile.txt" from inside the last public link shared folder using the public WebDAV API with password "wrong pass"
    But the public should be able to download file "/testfile.txt" from inside the last public link shared folder using the public WebDAV API with password "%public%"
    Examples:
      | new-permissions-role |
      | Edit                 |
      | Upload               |


  Scenario Outline: user shouldn't be allowed to updated public link without password to edit, upload file drop and secure viewer permission when enforce-password is enabled
    Given the following configs have been set:
      | config                                                 | value |
      | OCIS_SHARING_PUBLIC_SHARE_MUST_HAVE_PASSWORD           | false |
      | OCIS_SHARING_PUBLIC_WRITEABLE_SHARE_MUST_HAVE_PASSWORD | true  |
    And using spaces DAV path
    And user "Alice" has been created with default attributes
    And user "Alice" has created folder "folder"
    And user "Alice" has uploaded file with content "test file" to "folder/testfile.txt"
    And user "Alice" has created the following resource link share:
      | resource        | folder   |
      | space           | Personal |
      | permissionsRole | View     |
    And user "Alice" updates the last public link share using the permissions endpoint of the Graph API:
      | resource        | folder                 |
      | space           | Personal               |
      | permissionsRole | <new-permissions-role> |
    Then the HTTP status code should be "400"
    Examples:
      | new-permissions-role |
      | Edit                 |
      | Upload               |
      | File Drop            |


  Scenario Outline: user should be allowed to create a public link with a password in accordance with the password policy
    Given the following configs have been set:
      | config                                                 | value |
      | OCIS_SHARING_PUBLIC_WRITEABLE_SHARE_MUST_HAVE_PASSWORD | true  |
      | OCIS_PASSWORD_POLICY_MIN_CHARACTERS                    | 13    |
      | OCIS_PASSWORD_POLICY_MIN_LOWERCASE_CHARACTERS          | 3     |
      | OCIS_PASSWORD_POLICY_MIN_UPPERCASE_CHARACTERS          | 2     |
      | OCIS_PASSWORD_POLICY_MIN_DIGITS                        | 2     |
      | OCIS_PASSWORD_POLICY_MIN_SPECIAL_CHARACTERS            | 2     |
    And user "Alice" has been created with default attributes
    And using SharingNG
    And user "Alice" has created folder "folder"
    And user "Alice" has uploaded file with content "test file" to "folder/testfile.txt"
    When user "Alice" creates the following resource link share using the Graph API:
      | resource        | folder             |
      | space           | Personal           |
      | permissionsRole | <permissions-role> |
      | password        | 3s:5WW9uE5h=A      |
    Then the HTTP status code should be "200"
    And the public should not be able to download file "/testfile.txt" from inside the last public link shared folder using the public WebDAV API without a password
    And the public should not be able to download file "/testfile.txt" from inside the last public link shared folder using the public WebDAV API with password "wrong pass"
    But the public should be able to download file "/testfile.txt" from inside the last public link shared folder using the public WebDAV API with password "3s:5WW9uE5h=A"
    Examples:
      | permissions-role |
      | Edit             |
      | Upload           |


  Scenario Outline: user tries to create a public link with a password that does not comply with the password policy
    Given the following configs have been set:
      | config                                        | value |
      | OCIS_PASSWORD_POLICY_MIN_CHARACTERS           | 13    |
      | OCIS_PASSWORD_POLICY_MIN_LOWERCASE_CHARACTERS | 3     |
      | OCIS_PASSWORD_POLICY_MIN_UPPERCASE_CHARACTERS | 2     |
      | OCIS_PASSWORD_POLICY_MIN_DIGITS               | 2     |
      | OCIS_PASSWORD_POLICY_MIN_SPECIAL_CHARACTERS   | 2     |
    And user "Alice" has been created with default attributes
    And user "Alice" has uploaded file with content "test file" to "/testfile.txt"
    And user "Alice" has created folder "folder"
    And user "Alice" has uploaded file with content "test file" to "folder/testfile.txt"
    When user "Alice" creates the following resource link share using the Graph API:
      | resource        | folder             |
      | space           | Personal           |
      | permissionsRole | <permissions-role> |
      | password        | Pas1               |
    Then the HTTP status code should be "400"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["error"],
        "properties": {
          "error": {
            "type": "object",
            "required": ["code","innererror","message"],
            "properties": {
              "code": {"const": "invalidRequest"},
              "innererror": {
                "type": "object",
                "required": ["date","request-id"]
              },
              "message": {"const":  "at least 13 characters are required\nat least 3 lowercase letters are required\nat least 2 uppercase letters are required\nat least 2 numbers are required\nat least 2 special characters are required  !\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"}
            }
          }
        }
      }
      """
    Examples:
      | permissions-role |
      | View             |
      | Edit             |
      | Upload           |
      | File Drop        |


  Scenario Outline: user update a public link with a password in accordance with the password policy
    Given the following configs have been set:
      | config                                                 | value |
      | OCIS_SHARING_PUBLIC_SHARE_MUST_HAVE_PASSWORD           | false |
      | OCIS_SHARING_PUBLIC_WRITEABLE_SHARE_MUST_HAVE_PASSWORD | true  |
      | OCIS_PASSWORD_POLICY_MIN_CHARACTERS                    | 13    |
      | OCIS_PASSWORD_POLICY_MIN_LOWERCASE_CHARACTERS          | 3     |
      | OCIS_PASSWORD_POLICY_MIN_UPPERCASE_CHARACTERS          | 2     |
      | OCIS_PASSWORD_POLICY_MIN_DIGITS                        | 1     |
      | OCIS_PASSWORD_POLICY_MIN_SPECIAL_CHARACTERS            | 2     |
    And user "Alice" has been created with default attributes
    And user "Alice" has created folder "folder"
    And user "Alice" has uploaded file with content "test file" to "folder/testfile.txt"
    And using SharingNG
    And user "Alice" has created the following resource link share:
      | resource        | folder             |
      | space           | Personal           |
      | permissionsRole | <permissions-role> |
      | password        | 6a0Q;A3 +i^m[      |
    When user "Alice" sets the following password for the last link share using the Graph API:
      | resource | folder           |
      | space    | Personal         |
      | password | 6afsa0Q;A3 +i^m[ |
    Then the HTTP status code should be "200"
    And the public should not be able to download file "/testfile.txt" from inside the last public link shared folder using the public WebDAV API without a password
    And the public should not be able to download file "/testfile.txt" from inside the last public link shared folder using the public WebDAV API with password "wrong pass"
    But the public should be able to download file "/testfile.txt" from inside the last public link shared folder using the public WebDAV API with password "6afsa0Q;A3 +i^m["
    Examples:
      | permissions-role |
      | Edit             |
      | Upload           |


  Scenario Outline: user tries to update a public link with a password that does not comply with the password policy
    Given the following configs have been set:
      | config                                                 | value |
      | OCIS_SHARING_PUBLIC_SHARE_MUST_HAVE_PASSWORD           | false |
      | OCIS_SHARING_PUBLIC_WRITEABLE_SHARE_MUST_HAVE_PASSWORD | true  |
      | OCIS_PASSWORD_POLICY_MIN_CHARACTERS                    | 13    |
      | OCIS_PASSWORD_POLICY_MIN_LOWERCASE_CHARACTERS          | 3     |
      | OCIS_PASSWORD_POLICY_MIN_UPPERCASE_CHARACTERS          | 2     |
      | OCIS_PASSWORD_POLICY_MIN_DIGITS                        | 1     |
      | OCIS_PASSWORD_POLICY_MIN_SPECIAL_CHARACTERS            | 2     |
    And user "Alice" has been created with default attributes
    And user "Alice" has created folder "folder"
    And user "Alice" has uploaded file with content "test file" to "folder/testfile.txt"
    And using SharingNG
    And user "Alice" has created the following resource link share:
      | resource        | folder             |
      | space           | Personal           |
      | permissionsRole | <permissions-role> |
      | password        | 6a0Q;A3 +i^m[      |
    When user "Alice" sets the following password for the last link share using the Graph API:
      | resource | folder   |
      | space    | Personal |
      | password | Pws^     |
    Then the HTTP status code should be "400"
    Examples:
      | permissions-role |
      | Edit             |
      | Upload           |
      | File Drop        |


  Scenario Outline: user creates a public link with a password in accordance with the password policy (valid cases)
    Given the config "<config>" has been set to "<config-value>"
    And user "Alice" has been created with default attributes
    And using SharingNG
    And user "Alice" has uploaded file with content "test file" to "/testfile.txt"
    When user "Alice" creates the following resource link share using the Graph API:
      | resource        | testfile.txt |
      | space           | Personal     |
      | permissionsRole | View         |
      | password        | <password>   |
    Then the HTTP status code should be "200"
    And the public should not be able to download file "/testfile.txt" from inside the last public link shared folder using the public WebDAV API without a password
    And the public should not be able to download file "/testfile.txt" from inside the last public link shared folder using the public WebDAV API with password "wrong pass"
    But the public should be able to download file "/testfile.txt" from inside the last public link shared folder using the public WebDAV API with password "<password>"
    Examples:
      | config                                        | config-value | password                             |
      | OCIS_PASSWORD_POLICY_MIN_CHARACTERS           | 4            | Ps-1                                 |
      | OCIS_PASSWORD_POLICY_MIN_CHARACTERS           | 14           | Ps1:with space                       |
      | OCIS_PASSWORD_POLICY_MIN_LOWERCASE_CHARACTERS | 4            | PS1:test                             |
      | OCIS_PASSWORD_POLICY_MIN_UPPERCASE_CHARACTERS | 3            | PS1:TeƒsT                            |
      | OCIS_PASSWORD_POLICY_MIN_DIGITS               | 2            | PS1:test2                            |
      | OCIS_PASSWORD_POLICY_MIN_SPECIAL_CHARACTERS   | 2            | PS1:test pass                        |
      | OCIS_PASSWORD_POLICY_MIN_SPECIAL_CHARACTERS   | 33           | pS1! #$%&'()*+,-./:;<=>?@[\]^_`{  }~ |
      | OCIS_PASSWORD_POLICY_MIN_SPECIAL_CHARACTERS   | 5            | 1sameCharacterShouldWork!!!!!        |


  Scenario Outline: user tries to create a public link with a password that does not comply with the password policy (invalid cases)
    Given user "Alice" has been created with default attributes
    And user "Alice" has uploaded file with content "test file" to "/testfile.txt"
    When user "Alice" creates the following resource link share using the Graph API:
      | resource        | testfile.txt |
      | space           | Personal     |
      | permissionsRole | View         |
      | password        | <password>   |
    Then the HTTP status code should be "400"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["error"],
        "properties": {
          "error": {
            "type": "object",
            "required": ["code","innererror","message"],
            "properties": {
              "code": {"const": "invalidRequest"},
              "innererror": {
                "type": "object",
                "required": ["date","request-id"]
              },
              "message": {"const":  "<message>"}
            }
          }
        }
      }
      """
    Examples:
      | password | message                                   |
      | 1Pw:     | at least 8 characters are required        |
      | 1P:12345 | at least 1 lowercase letters are required |
      | test-123 | at least 1 uppercase letters are required |
      | Test-psw | at least 1 numbers are required           |
