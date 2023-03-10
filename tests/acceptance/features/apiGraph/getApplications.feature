@api @skipOnOcV10 @skipOnStable2.0
Feature: get applications
  As an user
  I want to be able to get applications information with existings roles

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files


  Scenario Outline: admin user lists all the groups
    Given the administrator has given "Alice" the role "<role>" using the settings api
    When user "Alice" gets all applications using the Graph API
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
          "items": [
            {
              "type": "object",
              "required": [
                "appRoles",
                "displayName"
              ],
              "properties": {
                "appRoles": {
                  "type": "array",
                  "items": [
                    {
                      "type": "object",
                      "required": [
                        "displayName",
                        "id"
                      ],
                      "properties": {
                        "displayName": {
                          "type": "string",
                          "enum": ["Guest"]
                        },
                        "id": {
                          "type": "string",
                          "pattern": "^%user_id_pattern%$"
                        }
                      }
                    },
                    {
                      "type": "object",
                      "required": [
                        "displayName",
                        "id"
                      ],
                      "properties": {
                        "displayName": {
                          "type": "string",
                          "enum": ["User"]
                        },
                        "id": {
                          "type": "string",
                          "pattern": "^%user_id_pattern%$"
                        }
                      }
                    },
                    {
                      "type": "object",
                      "required": [
                        "displayName",
                        "id"
                      ],
                      "properties": {
                        "displayName": {
                          "type": "string",
                          "enum": ["Space Admin"]
                        },
                        "id": {
                          "type": "string",
                          "pattern": "^%user_id_pattern%$"
                        }
                      }
                    },
                    {
                      "type": "object",
                      "required": [
                        "displayName",
                        "id"
                      ],
                      "properties": {
                        "displayName": {
                          "type": "string",
                          "enum": ["Admin"]
                        },
                        "id": {
                          "type": "string",
                          "pattern": "^%user_id_pattern%$"
                        }
                      }
                    }
                  ]
                },
                "displayName": {
                  "type": "string",
                  "enum": ["ownCloud Infinite Scale"]
                },
                "id": {
                  "type": "string",
                  "pattern": "^%user_id_pattern%$"
                }
              }
            }
          ]
        }
      }
    }
    """
    Examples:
      | role        |
      | Admin       |
      | Space Admin |
      | User        |
      | Guest       |
