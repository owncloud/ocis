@api
Feature: user GDPR (General Data Protection Regulation) report
  As a user
  I want to generate my GDPR report
  So that I can review what events are stored by the server

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files
    And using spaces DAV path


  Scenario: generate a GDPR report and check user data in the downloaded report
    When user "Alice" exports her GDPR report to "/.personal_data_export.json" using the Graph API
    And user "Alice" downloads the content of GDPR report ".personal_data_export.json"
    Then the HTTP status code of responses on each endpoint should be "201, 200" respectively
    And the downloaded JSON content should contain key 'user' and match
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
    When user "Alice" exports her GDPR report to "/.personal_data_export.json" using the Graph API
    And user "Alice" downloads the content of GDPR report ".personal_data_export.json"
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
            "Type",
            "Quota"
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
            },
            "Quota": {
              "type": ["number", "null"],
              "enum": [null]
            }
          }
        }
      }
    }
    """
