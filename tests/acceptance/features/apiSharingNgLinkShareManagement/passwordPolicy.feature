@env-config
Feature: Password policy for public links password

  Password requirements. set by default:
  | OCIS_SHARING_PUBLIC_SHARE_MUST_HAVE_PASSWORD  | true |
  | OCIS_PASSWORD_POLICY_MIN_CHARACTERS           | 8    |
  | OCIS_PASSWORD_POLICY_MIN_LOWERCASE_CHARACTERS | 1    |
  | OCIS_PASSWORD_POLICY_MIN_UPPERCASE_CHARACTERS | 1    |
  | OCIS_PASSWORD_POLICY_MIN_DIGITS               | 1    |
  | OCIS_PASSWORD_POLICY_MIN_SPECIAL_CHARACTERS   | 1    |

  Background:
    Given the following configs have been set:
      | service | config                                           | value |
      | sharing | SHARING_PASSWORD_POLICY_MIN_CHARACTERS           | 13    |
      | sharing | SHARING_PASSWORD_POLICY_MIN_LOWERCASE_CHARACTERS | 3     |
      | sharing | SHARING_PASSWORD_POLICY_MIN_UPPERCASE_CHARACTERS | 2     |
      | sharing | SHARING_PASSWORD_POLICY_MIN_DIGITS               | 2     |
      | sharing | SHARING_PASSWORD_POLICY_MIN_SPECIAL_CHARACTERS   | 2     |
    And user "Alice" has been created with default attributes
    And using SharingNG


  Scenario: user should be allowed to create a public link with a password in accordance with the password policy
    Given user "Alice" has created folder "folder"
    And user "Alice" has uploaded file with content "test file" to "folder/testfile.txt"
    When user "Alice" creates the following resource link share using the Graph API:
      | resource        | folder        |
      | space           | Personal      |
      | permissionsRole | Edit          |
      | password        | 3s:5WW9uE5h=A |
    Then the HTTP status code should be "200"
    And the public should not be able to download file "/testfile.txt" from inside the last public link shared folder using the public WebDAV API without a password
    And the public should not be able to download file "/testfile.txt" from inside the last public link shared folder using the public WebDAV API with password "wrong pass"
    But the public should be able to download file "/testfile.txt" from inside the last public link shared folder using the public WebDAV API with password "3s:5WW9uE5h=A"


  Scenario: user tries to create a public link with a password that does not comply with the password policy
    Given user "Alice" has uploaded file with content "test file" to "/testfile.txt"
    And user "Alice" has created folder "folder"
    And user "Alice" has uploaded file with content "test file" to "folder/testfile.txt"
    When user "Alice" creates the following resource link share using the Graph API:
      | resource        | folder   |
      | space           | Personal |
      | permissionsRole | View     |
      | password        | Pas1     |
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


  Scenario: user update a public link with a password in accordance with the password policy
    Given user "Alice" has created folder "folder"
    And user "Alice" has uploaded file with content "test file" to "folder/testfile.txt"
    And using SharingNG
    And user "Alice" has created the following resource link share:
      | resource        | folder        |
      | space           | Personal      |
      | permissionsRole | Upload        |
      | password        | 6a0Q;A3 +i^m[ |
    When user "Alice" sets the following password for the last link share using the Graph API:
      | resource | folder           |
      | space    | Personal         |
      | password | 6afsa0Q;A3 +i^m[ |
    Then the HTTP status code should be "200"
    And the public should not be able to download file "/testfile.txt" from inside the last public link shared folder using the public WebDAV API without a password
    And the public should not be able to download file "/testfile.txt" from inside the last public link shared folder using the public WebDAV API with password "wrong pass"
    But the public should be able to download file "/testfile.txt" from inside the last public link shared folder using the public WebDAV API with password "6afsa0Q;A3 +i^m["


  Scenario: user tries to update a public link with a password that does not comply with the password policy
    Given user "Alice" has created folder "folder"
    And user "Alice" has uploaded file with content "test file" to "folder/testfile.txt"
    And using SharingNG
    And user "Alice" has created the following resource link share:
      | resource        | folder        |
      | space           | Personal      |
      | permissionsRole | File Drop     |
      | password        | 6a0Q;A3 +i^m[ |
    When user "Alice" sets the following password for the last link share using the Graph API:
      | resource | folder   |
      | space    | Personal |
      | password | Pws^     |
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
