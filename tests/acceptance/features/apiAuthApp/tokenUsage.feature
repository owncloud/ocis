Feature: create auth-app token
  As a user
  I want to create auth-app Tokens
  So that I can let 3rd party apps use my account

  Background:
    Given user "Alice" has been created with default attributes
    And using auth-app token


  Scenario: admin creates user using auth-app token
    Given user "Admin" has created auth-app token with expiration time "1h" using the auth-app API
    When the user "Admin" creates a new user with the following attributes using the Graph API:
      | userName       | Brian                 |
      | displayName    | This is another Brian |
      | email          | brian@example.com     |
      | password       | brian123              |
      | accountEnabled | true                  |
    Then the HTTP status code should be "201"
    And the JSON response should contain space called "Alice Hansen" and match
      """
      {
        "type": "object",
        "required": [
          "accountEnabled",
          "displayName",
          "id",
          "mail",
          "onPremisesSamAccountName",
          "surname",
          "userType"
        ],
        "properties": {
          "accountEnabled": { "const": true },
          "displayName": { "const": "This is another Brian" },
          "id": { "pattern": "%user_id_pattern%" },
          "mail": { "const": "brian@example.com" },
          "onPremisesSamAccountName": { "const": "Brian" },
          "surname": { "const":"Brian" },
          "userType": { "const":"Member" }
        }
      }
      """
    And user "Brian" should be able to upload file "filesForUpload/lorem.txt" to "lorem.txt"


  Scenario: user lists their drives using auth-app token
    Given user "Alice" has created auth-app token with expiration time "1h" using the auth-app API
    When user "Alice" lists all available spaces via the Graph API
    Then the HTTP status code should be "200"
    And the JSON response should contain space called "Alice Hansen" and match
      """
      {
        "type": "object",
        "required": [
          "driveType",
          "driveAlias",
          "name",
          "id",
          "quota",
          "root",
          "webUrl"
        ],
        "properties": {
          "name": { "const": "Alice Hansen" },
          "driveType": { "const": "personal" },
          "driveAlias": { "const": "personal/alice" },
          "id": { "pattern": "%space_id_pattern%" },
          "quota": {
             "type": "object",
             "required": ["state"],
             "properties": {
                "state": { "const": "normal" }
             }
          },
          "root": {
            "type": "object",
            "required": ["webDavUrl"],
            "properties": {
                "webDavUrl": { "const": "%base_url%/dav/spaces/%space_id%" }
             }
          },
          "webUrl": {
            "const": "%base_url%/f/%space_id%"
          }
        }
      }
      """