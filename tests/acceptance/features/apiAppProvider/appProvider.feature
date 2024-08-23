Feature: App Provider
  As a user
  I want to access files with collaboration service apps
  So that I can get the content of the file


  Scenario Outline: open file with .odt extension with collaboration app
    Given user "Alice" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file "filesForUpload/simple.odt" to "simple.odt"
    And we save it into "FILEID"
    When user "Alice" sends HTTP method "POST" to URL "<endpoint>"
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
            "pattern": "^https:\\/\\/fakeoffice\\.owncloud\\.test\\/not\\/relevant\\?WOPISrc=http%3A%2F%2Fwopiserver%3A9300%2Fwopi%2Ffiles%2F[a-fA-F0-9]{64}$"
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
      | endpoint                                         |
      | /app/open?file_id=<<FILEID>>&app_name=FakeOffice |
      | /app/open?file_id=<<FILEID>>                     |
