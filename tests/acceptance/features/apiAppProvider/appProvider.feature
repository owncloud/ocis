Feature: App Provider


  Scenario: open file with Collabora
    Given user "Alice" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file "filesForUpload/simple.odt" to "simple.odt"
    When user "Alice" opens file "simple.odt" with Collabora service
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "app_url",
          "method",
          "form_parameters"
        ],
        "properties": {
          "app_url": {
            "type": "string",
            "pattern": "^https:\\/\\/host\\.docker\\.internal:9980\\/browser\\/8ec5fda\\/cool\\.html\\?WOPISrc=https%3A%2F%2Fhost\\.docker\\.internal%3A9300%2Fwopi%2Ffiles%2F[a-fA-F0-9]{64}$"
          },
          "method": {
            "const": "POST"
          },
          "form_parameters": {
            "type": "object",
            "required": [
              "access_token",
              "access_token_ttl"
            ],
            "properties": {
              "access_token": {
                "type": "string"
              },
              "access_token_ttl": {
                "type": "string"
              }
            }
          }
        }
      }
      """
