Feature: collaboration (wopi)
  As a user
  I want to access files with collaboration service apps
  So that I can collaborate with other users


  Scenario Outline: open file with .odt extension
    Given user "Alice" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file "filesForUpload/simple.odt" to "simple.odt"
    And we save it into "FILEID"
    When user "Alice" sends HTTP method "POST" to URL "<app-endpoint>"
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
            "pattern": "^.*\\?WOPISrc=.*wopi%2Ffiles%2F[a-fA-F0-9]{64}$"
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
    Examples:
      | app-endpoint                                     |
      | /app/open?file_id=<<FILEID>>&app_name=FakeOffice |
      | /app/open?file_id=<<FILEID>>                     |
