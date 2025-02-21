@ocm @notification @email
Feature: ocm notifications
  As a user
  I want to manage my notification settings
  So that I do not get notified of unimportant events

  Background:
    Given user "Alice" has been created with default attributes
    And "Alice" has created the federation share invitation
    And user "Alice" has uploaded file with content "ocm test" to "textfile.txt"
    And using server "REMOTE"
    And user "Brian" has been created with default attributes
    And "Brian" has accepted invitation


  Scenario: federated user disables mail and in-app notification for "Share Received" event
    When user "Brian" disables notification for the following events using the settings API:
      | Share Received | mail,in-app |
    Then the HTTP status code should be "201"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["value"],
        "properties": {
          "value": {
            "type": "object",
            "required": ["identifier","value"],
            "properties": {
              "identifier":{
                "type": "object",
                "required": ["extension","bundle","setting"],
                "properties": {
                  "extension":{ "const": "ocis-accounts" },
                  "bundle":{ "const": "profile" },
                  "setting":{ "const": "event-share-created-options" }
                }
              },
              "value":{
                "type": "object",
                "required": ["id","bundleId","settingId","accountUuid","resource","collectionValue"],
                "properties":{
                  "id":{ "pattern":"%uuidv4_pattern%" },
                  "bundleId":{ "pattern":"%uuidv4_pattern%" },
                  "settingId":{ "pattern":"%uuidv4_pattern%" },
                  "accountUuid":{ "pattern":"%uuidv4_pattern%" },
                  "resource":{
                    "type": "object",
                    "required":["type"],
                    "properties": {
                      "type":{ "const": "TYPE_USER" }
                    }
                  },
                  "collectionValue":{
                    "type": "object",
                    "required":["values"],
                    "properties": {
                      "values":{
                        "type": "array",
                        "maxItems": 2,
                        "minItems": 2,
                        "uniqueItems": true,
                        "items": {
                          "oneOf": [
                            {
                              "type": "object",
                              "required": ["key","boolValue"],
                              "properties": {
                                "key":{ "const": "mail" },
                                "boolValue":{ "const": false }
                              }
                            },
                            {
                              "type": "object",
                              "required": ["key","boolValue"],
                              "properties": {
                                "key":{ "const": "in-app" },
                                "boolValue":{ "const": false }
                              }
                            }
                          ]
                        }
                      }
                    }
                  }
                }
              }
            }
          }
        }
      }
      """
    And using server "LOCAL"
    And user "Alice" has sent the following resource share invitation to federated user:
      | resource        | textfile.txt |
      | space           | Personal     |
      | sharee          | Brian        |
      | shareType       | user         |
      | permissionsRole | Viewer       |
    And using server "REMOTE"
    When user "Brian" lists all notifications
    Then the HTTP status code should be "200"
    And the notifications should be empty
    And user "Brian" should have "0" emails

  @issue-10718 @issue-11042
  Scenario: federated user disables mail and in-app notification for "Share Removed" event
    Given using server "LOCAL"
    And user "Alice" has sent the following resource share invitation to federated user:
      | resource        | textfile.txt |
      | space           | Personal     |
      | sharee          | Brian        |
      | shareType       | user         |
      | permissionsRole | Viewer       |
    And using server "REMOTE"
    When user "Brian" disables notification for the following events using the settings API:
      | Share Removed | mail, in-app |
    Then the HTTP status code should be "201"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["value"],
        "properties": {
          "value": {
            "type": "object",
            "required": ["identifier","value"],
            "properties": {
              "identifier":{
                "type": "object",
                "required": ["extension","bundle","setting"],
                "properties": {
                  "extension":{ "const": "ocis-accounts" },
                  "bundle":{ "const": "profile" },
                  "setting":{ "const": "event-share-removed-options" }
                }
              },
              "value":{
                "type": "object",
                "required": ["id","bundleId","settingId","accountUuid","resource","collectionValue"],
                "properties":{
                  "id":{ "pattern":"%uuidv4_pattern%" },
                  "bundleId":{ "pattern":"%uuidv4_pattern%" },
                  "settingId":{ "pattern":"%uuidv4_pattern%" },
                  "accountUuid":{ "pattern":"%uuidv4_pattern%" },
                  "resource":{
                    "type": "object",
                    "required":["type"],
                    "properties": {
                      "type":{ "const": "TYPE_USER" }
                    }
                  },
                  "collectionValue":{
                    "type": "object",
                    "required":["values"],
                    "properties": {
                      "values":{
                        "type": "array",
                        "maxItems": 2,
                        "minItems": 2,
                        "uniqueItems": true,
                        "items": {
                          "oneOf": [
                            {
                              "type": "object",
                              "required": ["key","boolValue"],
                              "properties": {
                                "key":{ "const": "mail" },
                                "boolValue":{ "const": false }
                              }
                            },
                            {
                              "type": "object",
                              "required": ["key","boolValue"],
                              "properties": {
                                "key":{ "const": "in-app" },
                                "boolValue":{ "const": false }
                              }
                            }
                          ]
                        }
                      }
                    }
                  }
                }
              }
            }
          }
        }
      }
      """
    And using server "LOCAL"
    And user "Alice" has removed the access of user "Brian" from resource "textfile.txt" of space "Personal"
    And using server "REMOTE"
    And user "Brian" should have "1" emails
    And user "Brian" should get a notification with subject "Resource shared" and message:
      | message                                   |
      | Alice Hansen shared textfile.txt with you |
    And user "Brian" should not have a notification related to resource "textfile.txt" with subject "Resource unshared"
