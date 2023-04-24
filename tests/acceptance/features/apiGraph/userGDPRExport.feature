@api @skipOnOcV10
Feature: user GDPR (General Data Protection Regulation) report
  As a user
  I want to get or generate GDPR report of my own data
  So that i can see the report of my own data at any time

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files
    And using spaces DAV path


  Scenario: generate a GDPR report and check user data in the downloaded report
    When user "Alice" generates GDPR reports of his own data to "/.personal_data_export.json"
    And user "Alice" downloads the content of generated GDPR report of file ".personal_data_export.json" using password "123456"
    Then the HTTP status code of responses on each endpoint should be "201, 200" respectively
    And the downloaded JSON should contain key "user" and match
    """
    {
      "type": "object",
      "required": [
        "id",
        "username",
        "mail",
        "display_name",
        "uid_number",
        "gid_number"
      ],
      "properties": {
        "id": {
          "type": "object",
          "required": [
            "idp",
            "opaque_id",
            "type"
          ],
          "properties": {
            "idp": {
              "type": "string",
              "pattern": "^%base_url%$"
            },
            "opaque_id": {
              "type": "string",
              "pattern": "^%user_id_pattern%$"
            },
            "type": {
              "type": "number",
              "enum": [1]
            }
          }
        },
        "username": {
          "type": "string",
          "enum": ["Alice"]
        },
        "mail": {
          "type": "string",
          "enum": ["alice@example.org"]
        },
        "display_name": {
          "type": "string",
          "enum": ["Alice Hansen"]
        },
        "uid_number": {
          "type": "number",
          "enum": [99]
        },
        "gid_number": {
          "type": "number",
          "enum": [99]
        }
      }
    }
    """

  Scenario: generate a GDPR report and check events when a user is created
    When user "Alice" generates GDPR reports of his own data to "/.personal_data_export.json"
    And user "Alice" downloads the content of generated GDPR report of file ".personal_data_export.json" using password "123456"
    Then the HTTP status code of responses on each endpoint should be "201, 200" respectively
    And the downloaded JSON content should contain event type "events.UserCreated" in item 'events' and should match
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
            "UserID"
          ],
          "properties": {
            "Executant": {
              "type": "object",
              "required": [
                "idp",
                "opaque_id",
                "type"
              ],
              "properties": {
                "idp": {
                  "type": "string",
                  "pattern": "^%base_url%$"
                },
                "opaque_id": {
                  "type": "string",
                  "pattern": "^%user_id_pattern%$"
                },
                "type": {
                  "type": "number",
                  "enum": [1]
                }
              }
            },
            "UserID": {
              "type": "string",
              "pattern": "^%user_id_pattern%$"
            }
          }
        }
      }
    }
    """
    And the downloaded JSON content should contain event type "events.SpaceCreated" in item 'events' and should match
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
            "Name",
            "Type"
          ],
          "properties": {
            "Executant": {
              "type": "object",
              "required": [
                "idp",
                "opaque_id",
                "type"
              ],
              "properties": {
                "idp": {
                  "type": "string",
                  "pattern": "^%base_url%$"
                },
                "opaque_id": {
                  "type": "string",
                  "pattern": "^%user_id_pattern%$"
                },
                "type": {
                  "type": "number",
                  "enum": [1]
                }
              }
            },
            "Name": {
              "type": "string",
              "enum": ["Alice Hansen"]
            },
            "Type": {
              "type": "string",
              "enum": ["personal"]
            }
          }
        }
      }
    }
    """
