Feature: get user's own information
  As user
  I want to be able to retrieve my own information
  So that I can see my information

  Background:
    Given user "Alice" has been created with default attributes


  Scenario: user gets his/her own information with no group involvement
    When the user "Alice" retrieves her information using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "id",
          "mail",
          "onPremisesSamAccountName"
        ],
        "properties": {
          "id" : {
            "type": "string",
            "pattern": "^%user_id_pattern%$"
          },
          "mail": {
            "type": "string",
            "enum": ["alice@example.org"]
          },
          "onPremisesSamAccountName": {
            "type": "string",
            "enum": ["Alice"]
          }
        }
      }
      """


  Scenario: user gets his/her own information with group involvement
    Given group "tea-lover" has been created
    And group "coffee-lover" has been created
    And user "Alice" has been added to group "tea-lover"
    And user "Alice" has been added to group "coffee-lover"
    When the user "Alice" retrieves her information using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "id",
          "mail",
          "onPremisesSamAccountName"
        ],
        "properties": {
          "id" : {
            "type": "string",
            "pattern": "^%user_id_pattern%$"
          },
          "mail": {
            "type": "string",
            "enum": ["alice@example.org"]
          },
          "onPremisesSamAccountName": {
            "type": "string",
            "enum": ["Alice"]
          },
          "memberOf": {
            "type": "array",
            "minItems": 2,
            "maxItems": 2,
            "uniqueItems": true,
            "items": {
              "oneOf": [
                {
                  "type": "object",
                  "required": [
                    "displayName"
                  ],
                  "properties": {
                    "displayName": {
                      "type": "string",
                      "enum": ["tea-lover"]
                    }
                  }
                },
                {
                  "type": "object",
                  "required": [
                    "displayName"
                  ],
                  "properties": {
                    "displayName": {
                      "type": "string",
                      "enum": ["coffee-lover"]
                    }
                  }
                }
              ]
            }
          }
        }
      }
      """
