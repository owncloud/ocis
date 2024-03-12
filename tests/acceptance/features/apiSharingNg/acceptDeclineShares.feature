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


    Scenario: accept incoming file share when auto accept is disabled
      Given user "Brian" has disabled auto-accepting
      And user "Alice" has uploaded file with content "hello world" to "/textfile0.txt"
      And user "Alice" has sent the following share invitation:
        | resource        | textfile0.txt |
        | space           | Personal      |
        | sharee          | Brian         |
        | shareType       | user          |
        | permissionsRole | Viewer        |
      When user "Brian" accepts share "/textfile0.txt" using the Graph API
      Then the HTTP status code should be "201"
      And the JSON data of the response should match
        """
        {
          "type": "object",
          "required": [
            "@client.synchronize"
          ],
          "properties": {
            "@client.synchronize": {
              "const": true
            }
          }
        }
        """


  Scenario: accept incoming folder share when auto accept is disabled
    Given user "Brian" has disabled auto-accepting
    Given user "Alice" has created folder "folder"
    And user "Alice" has sent the following share invitation:
      | resource        | folder   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Viewer   |
    When user "Brian" accepts share "folder" using the Graph API
    Then the HTTP status code should be "201"
    And the JSON data of the response should match
      """
        {
          "type": "object",
          "required": [
            "@client.synchronize"
          ],
          "properties": {
            "@client.synchronize": {
              "const": true
            }
          }
        }
        """




