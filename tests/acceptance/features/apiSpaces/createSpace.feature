@api
Feature: create space
  Only user with admin and SpaceAdmin permissions can create new spaces

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files


  Scenario Outline: The user without permissions to create space cannot create a Space via Graph API
    Given the administrator has given "Alice" the role "<role>" using the settings api
    When user "Alice" creates a space "Project Mars" of type "project" with the default quota using the GraphApi
    Then the HTTP status code should be "401"
    And the user "Alice" should not have a space called "share space"
    Examples:
      | role  |
      | User  |
      | Guest |


  Scenario Outline: An admin or space admin user can create a Space via the Graph API with default quota
    Given the administrator has given "Alice" the role "<role>" using the settings api
    When user "Alice" creates a space "Project Mars" of type "project" with the default quota using the GraphApi
    Then the HTTP status code should be "201"
    And the JSON response should contain space called "Project Mars" and match
    """
    {
      "type": "object",
      "required": [
        "driveType",
        "driveAlias",
        "name",
        "id",
        "quota",
        "root",
        "webUrl"
      ],
      "properties": {
        "name": {
          "type": "string",
          "enum": ["Project Mars"]
        },
        "driveType": {
           "type": "string",
          "enum": ["project"]
        },
        "driveAlias": {
           "type": "string",
          "enum": ["project/project-mars"]
        },
        "id": {
           "type": "string",
          "enum": ["%space_id%"]
        },
        "quota": {
           "type": "object",
           "required": [
            "total"
           ],
           "properties": {
              "state": {
                "type": "number",
                "enum": [1000000000]
              }
           }
        },
        "root": {
          "type": "object",
          "required": [
            "webDavUrl"
          ],
          "properties": {
              "webDavUrl": {
                "type": "string",
                "enum": ["%base_url%/dav/spaces/%space_id%"]
              }
           }
        },
        "webUrl": {
          "type": "string",
          "enum": ["%base_url%/f/%space_id%"]
        }
      }
    }
    """
    Examples:
      | role        |
      | Admin       |
      | Space Admin |


  Scenario Outline: An admin or space admin user can create a Space via the Graph API with certain quota
    Given the administrator has given "Alice" the role "<role>" using the settings api
    When user "Alice" creates a space "Project Venus" of type "project" with quota "2000" using the GraphApi
    Then the HTTP status code should be "201"
    And the JSON response should contain space called "Project Venus" and match
    """
    {
      "type": "object",
      "required": [
        "driveType",
        "name",
        "id",
        "quota",
        "root",
        "webUrl"
      ],
      "properties": {
        "name": {
          "type": "string",
          "enum": ["Project Venus"]
        },
        "driveType": {
           "type": "string",
          "enum": ["project"]
        },
        "id": {
           "type": "string",
          "enum": ["%space_id%"]
        },
        "quota": {
           "type": "object",
           "required": [
            "total"
           ],
           "properties": {
              "state": {
                "type": "number",
                "enum": [2000]
              }
           }
        },
        "root": {
          "type": "object",
          "required": [
            "webDavUrl"
          ],
          "properties": {
              "webDavUrl": {
                "type": "string",
                "enum": ["%base_url%/dav/spaces/%space_id%"]
              }
           }
        },
        "webUrl": {
          "type": "string",
          "enum": ["%base_url%/f/%space_id%"]
        }
      }
    }
    """
    Examples:
      | role        |
      | Admin       |
      | Space Admin |
