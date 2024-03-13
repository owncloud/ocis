Feature: accept/decline incoming share
  As a user
  I want to have control over the share received
  So that I can filter out the files and folder shared to Me

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
    And using spaces DAV path


  Scenario: decline shared file when auto accept is enabled
    And user "Alice" has created folder "FolderToShare"
    And user "Alice" has uploaded file with content "hello world" to "/textfile0.txt"
    And user "Alice" has sent the following share invitation:
      | resource        | textfile0.txt |
      | space           | Personal      |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | Viewer        |
    When user "Brian" declines share "/textfile0.txt" using the Graph API
    And user "Brian" lists the shares shared with him using the Graph API
    Then the HTTP status code of responses on all endpoints should be "200"
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
          "minItems": 1,
          "maxItems": 1,
          "items": {
            "type": "object",
            "required": [
              "@client.synchronize"
            ],
            "properties": {
              "@client.synchronize": {
                "const": false
              }
            }
          }
        }
      }
    }
    """


  Scenario: decline shared folder when auto accept is enabled
    Given user "Alice" has created folder "folder"
    And user "Alice" has sent the following share invitation:
      | resource        | folder   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Viewer   |
    When user "Brian" declines share "folder" using the Graph API
    And user "Brian" lists the shares shared with him using the Graph API
    Then the HTTP status code of responses on all endpoints should be "200"
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
          "minItems": 1,
          "maxItems": 1,
          "items": {
            "type": "object",
            "required": [
              "@client.synchronize"
            ],
            "properties": {
              "@client.synchronize": {
                "const": false
              }
            }
          }
        }
      }
    }
    """

