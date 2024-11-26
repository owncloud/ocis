Feature: create auth token
  As a admin
  I want to create App Tokens
  So that I can use 3rd party apps


  Scenario: admin creates app token
    When the administrator creates app token with expiration time "72h" using the API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "token",
          "expiration_date",
          "created_date",
          "label"
        ],
        "properties": {
          "token": {
            "type": "string",
            "pattern": "^[a-zA-Z0-9]{16}$"
          },
          "label": {
            "const": "Generated via API"
          }
        }
      }
      """


  Scenario: admin lists app token
    Given the administrator has created app token with expiration time "72h" using the API
    When admin lists all created tokens
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
            "token",
            "expiration_date",
            "created_date",
            "label"
          ],
          "properties": {
            "token": {
              "type": "string",
              "pattern": "^\\$2a\\$11\\$[A-Za-z0-9./]{53}$"
            },
            "label": {
              "const": "Generated via API"
            }
          }
        }
      }
      """
