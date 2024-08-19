@ocm
Feature: create invitation
  As a user
  I can create an invitations and send it to the person I want to share with

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files


  Scenario: user creates invitation
    Given using server "LOCAL"
    When "Alice" creates the federation share invitation
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "expiration",
          "token"
        ],
        "properties": {
          "expiration": {
            "type": "integer",
            "pattern": "^[0-9]{10}$"
          },
          "token": {
            "type": "string",
            "pattern": "^%fed_invitation_token_pattern%$"
          }
        }
      }
      """
    When "Alice" lists the created invitations
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "array",
        "minItems": 1,
        "maxItems": 1,
        "items": {
          "type": "object",
          "required": [
            "expiration",
            "token"
          ],
          "properties": {
            "expiration": {
              "type": "integer",
              "pattern": "^[0-9]{10}$"
            },
            "token": {
              "type": "string",
              "pattern": "^%fed_invitation_token_pattern%$"
            }
          }
        }
      }
      """

  @issue-9591
  Scenario: user creates invitation with email and description
    Given using server "LOCAL"
    When "Alice" creates the federation share invitation with email "brian@example.com" and description "a share invitation from Alice"
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "expiration",
          "token",
          "description",
          "recipient"
        ],
        "properties": {
          "expiration": {
            "type": "integer",
            "pattern": "^[0-9]{10}$"
          },
          "token": {
            "type": "string",
            "pattern": "^%fed_invitation_token_pattern%$"
          },
          "description": {
            "const": "a share invitation from Alice"
          },
          "recipient": {
            "const": "brian@example.com"
          }
        }
      }
      """
    And the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "array",
        "minItems": 1,
        "maxItems": 1,
        "items": {
          "type": "object",
          "required": [
            "expiration",
            "token"
          ],
          "properties": {
            "expiration": {
              "type": "integer",
              "pattern": "^[0-9]{10}$"
            },
            "token": {
              "type": "string",
              "pattern": "^%fed_invitation_token_pattern%$"
            },
            "description": {
              "const": "a share invitation from Alice"
            },
            "recipient": {
              "const": "brian@example.com"
            }
          }
        }
      }
      """

  @env-config
  Scenario: user cannot see expired invitation tokens
    Given using server "LOCAL"
    And the config "OCM_OCM_INVITE_MANAGER_TOKEN_EXPIRATION" has been set to "1s"
    When "Alice" creates the federation share invitation
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "expiration",
          "token"
        ],
        "properties": {
          "expiration": {
            "type": "integer",
            "pattern": "^[0-9]{10}$"
          },
          "token": {
            "type": "string",
            "pattern": "^%fed_invitation_token_pattern%$"
          }
        }
      }
      """
    And the user waits "2" seconds for the token to expire
    When "Alice" lists the created invitations
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "array",
        "minItems": 0,
        "maxItems": 0
      }
      """