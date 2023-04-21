@api @skipOnOcV10
Feature: user GDPR (General Data Protection Regulation) report
  As a user
  I want to get or generate GDPR report of my own data
  So that i can see the report of my own data at any time

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files
    And using spaces DAV path

  Scenario: request and check a GDPR report when user upload a file
    Given user "Alice" has uploaded file with content "sample text" to "lorem.txt"
    When user "Alice" generates GDPR reports of his own data to "/.personal_data_export.json"
    And user "Alice" downloads the content of generated GDPR report of file "personal_data_export.json" using password "123456"
    Then the HTTP status code of responses on each endpoint should be "201, 200" respectively
    And the downloaded JSON content should contain event type "events.FileUploaded" in item 'events' and should match
    """
    {
      "type": "object",
      "required": [
        "event"
      ],
      "properties": {
        "event" : {
          "type": "object",
          "required": [
            "Executant",
            "Owner",
            "Ref",
            "SpaceOwner"
          ],
          "properties": {
            "Ref": {
              "type": "object",
              "required": [
                "path"
              ],
              "properties": {
                "path" : {
                  "type": "string",
                  "enum": ["./lorem.txt"]
                }
              }
            }
          }
        }
      }
    }
    """
