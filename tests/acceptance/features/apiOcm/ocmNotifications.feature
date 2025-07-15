@ocm @notification @email
Feature: ocm notifications
  As a user
  I want to manage my notification settings
  So that I do not get notified of unimportant events

  Background:
    Given user "Alice" has been created with default attributes


  Scenario: federated user disables mail and in-app notification for "Share Received" event
    Given user "Alice" has uploaded file with content "ocm test" to "textfile.txt"
    And "Alice" has created the federation share invitation
    And using server "REMOTE"
    And user "Brian" has been created with default attributes
    And "Brian" has accepted invitation
    When user "Brian" disables notification for the following event using the settings API:
      | event             | Share Received |
      | notificationTypes | mail,in-app    |
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

  @issue-10059
  Scenario: federated user gets an email notification if their email was specified when creating the federation share invitation
    Given using server "REMOTE"
    And user "David" has been created with default attributes
    And using server "LOCAL"
    When "Alice" has created the federation share invitation with email "david@example.org" and description "a share invitation from Alice"
    And user "David" should have received the following email from user "Alice" ignoring whitespaces
      """
      Hi,

      Alice Hansen (alice@example.org) wants to start sharing collaboration resources with you.

      Please visit your federation settings and use the following details:
        Token: %fed_invitation_token%
        ProviderDomain: %local_base_url%
      """

  @issue-11042
  Scenario: no in-app notification should pop-up for unshared resource when "Share Removed" event is disabled
    Given using server "REMOTE"
    And user "Brian" has been created with default attributes
    And "Brian" has created the federation share invitation
    And using server "LOCAL"
    And "Alice" has accepted invitation
    And user "Alice" has uploaded file with content "ocm test" to "textfile.txt"
    And user "Alice" has sent the following resource share invitation to federated user:
      | resource        | textfile.txt |
      | space           | Personal     |
      | sharee          | Brian        |
      | shareType       | user         |
      | permissionsRole | Viewer       |
    And using server "REMOTE"
    When user "Brian" disables notification for the following event using the settings API:
      | event             | Share Removed |
      | notificationTypes | in-app        |
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
                        "maxItems": 1,
                        "minItems": 1,
                        "items": {
                          "oneOf": [
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
    And user "Brian" should get a notification with subject "Resource shared" and message:
      | message                                   |
      | Alice Hansen shared textfile.txt with you |
    And user "Brian" should not have a notification related to resource "textfile.txt" with subject "Resource unshared"

  @issue-10718
  Scenario: federation user gets an in-app notification for share received from local user
    Given using server "REMOTE"
    And user "Brian" has been created with default attributes
    And "Brian" has created the federation share invitation
    And using server "LOCAL"
    And user "Alice" has uploaded file with content "ocm test" to "textfile.txt"
    And "Alice" has accepted invitation
    When user "Alice" sends the following resource share invitation to federated user using the Graph API:
      | resource        | textfile.txt |
      | space           | Personal     |
      | sharee          | Brian        |
      | shareType       | user         |
      | permissionsRole | Viewer       |
    Then the HTTP status code should be "200"
    And using server "REMOTE"
    And user "Brian" should get a notification with subject "Resource shared" and message:
      | message                                   |
      | Alice Hansen shared textfile.txt with you |

  @issue-11042
  Scenario: federation user gets an in-app notification for share removed from local user
    Given using server "REMOTE"
    And user "Brian" has been created with default attributes
    And "Brian" has created the federation share invitation
    And using server "LOCAL"
    And user "Alice" has uploaded file with content "ocm test" to "textfile.txt"
    And "Alice" has accepted invitation
    And user "Alice" has sent the following resource share invitation to federated user:
      | resource        | textfile.txt |
      | space           | Personal     |
      | sharee          | Brian        |
      | shareType       | user         |
      | permissionsRole | Viewer       |
    When user "Alice" removes the access of user "Brian" from resource "textfile.txt" of space "Personal" using the Graph API
    Then the HTTP status code should be "204"
    And using server "REMOTE"
    And user "Brian" should get a notification with subject "Resource shared" and message:
      | message                                   |
      | Alice Hansen shared textfile.txt with you |
    And user "Brian" should get a notification with subject "Resource unshared" and message:
      | message                                   |
      | Alice Hansen unshared textfile.txt with you |
