@ocm
Feature: search federation users
  As a user
  I can find federation users after accepting an invitation to share resources


  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Carol    |
    And using server "REMOTE"
    And user "Brian" has been created with default attributes and without skeleton files


  Scenario: users search for federation users by display name
    Given using server "LOCAL"
    And "Alice" has created the federation share invitation
    And using server "REMOTE"
    And "Brian" has accepted invitation
    When user "Brian" searches for user "ali" using Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "value"
        ],
        "properties": {
          "value": {
            "type": "array",
            "minItems": 1,
            "maxItems": 1,
            "items": {
              "type": "object",
              "required": [
                "displayName",
                "id"
              ],
              "properties": {
                "displayName": {
                  "const": "Alice Hansen"
                },
                "id": {
                  "type": "string",
                  "pattern": "^%user_id_pattern%$"
                }
              }
            }
          }
        }
      }
      """
    And using server "LOCAL"
    When user "Alice" searches for user "bri" using Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "value"
        ],
        "properties": {
          "value": {
            "type": "array",
            "minItems": 1,
            "maxItems": 1,
            "items": {
              "type": "object",
              "required": [
                "displayName",
                "id"
              ],
              "properties": {
                "displayName": {
                  "const": "Brian Murphy"
                },
                "id": {
                  "type": "string",
                  "pattern": "^%user_id_pattern%$"
                }
              }
            }
          }
        }
      }
      """


  Scenario: sers search for federation users by email
    Given using server "LOCAL"
    And "Alice" has created the federation share invitation
    And using server "REMOTE"
    And "Brian" has accepted invitation
    When user "Brian" searches for user "%22alice@example.org%22" using Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "value"
        ],
        "properties": {
          "value": {
            "type": "array",
            "minItems": 1,
            "maxItems": 1,
            "items": {
              "type": "object",
              "required": [
                "displayName",
                "id"
              ],
              "properties": {
                "displayName": {
                  "const": "Alice Hansen"
                },
                "id": {
                  "type": "string",
                  "pattern": "^%user_id_pattern%$"
                }
              }
            }
          }
        }
      }
      """
    And using server "LOCAL"
    When user "Alice" searches for user "%22brian@example.org%22" using Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "value"
        ],
        "properties": {
          "value": {
            "type": "array",
            "minItems": 1,
            "maxItems": 1,
            "items": {
              "type": "object",
              "required": [
                "displayName",
                "id"
              ],
              "properties": {
                "displayName": {
                  "const": "Brian Murphy"
                },
                "id": {
                  "type": "string",
                  "pattern": "^%user_id_pattern%$"
                }
              }
            }
          }
        }
      }
      """


  Scenario: sers search for federation users without federated connection
    Given using server "LOCAL"
    And "Alice" has created the federation share invitation
    And using server "REMOTE"
    And "Brian" has accepted invitation
    When user "Brian" searches for user "%22carol@example.org%22" using Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "value"
        ],
        "properties": {
          "value": {
            "type": "array",
            "minItems": 0,
            "maxItems": 0
          }
        }
      }
      """
    And using server "LOCAL"
    When user "Carol" searches for user "bria" using Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "value"
        ],
        "properties": {
          "value": {
            "type": "array",
            "minItems": 0,
            "maxItems": 0
          }
        }
      }
      """
    And using server "REMOTE"

# TODO try to find federation users after deleting federated conection