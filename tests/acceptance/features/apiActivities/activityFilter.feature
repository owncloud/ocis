@issue-9712
Feature: activity filter
  As a user
  I want to filter activities
  So that I can track modifications of specific resource

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files


  Scenario: check activity with depth filter
    Given user "Alice" has created folder "/New Folder"
    And user "Alice" has created folder "/New Folder/Sub Folder"
    And user "Alice" has uploaded file with content "ownCloud test text file 0" to "/New Folder/Sub Folder/textfile0.txt"
    When user "Alice" lists the activities for folder "New Folder" of space "Personal" with depth "0" using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["value"],
        "properties": {
          "value": {
            "type": "array",
            "minItems": 1,
            "maxItems": 1,
            "items": {
              "type": "object",
              "required": ["id", "template", "times"],
              "properties": {
                "template": {
                  "type": "object",
                  "required": ["message", "variables"],
                  "properties": {
                    "message": {
                      "const": "{user} added {resource} to {space}"
                    },
                    "variables": {
                      "type": "object",
                      "required": ["resource", "space", "user"],
                      "properties": {
                        "resource": {
                          "type": "object",
                          "required": ["id", "name"],
                          "properties": {
                            "name": {
                              "const": "New Folder"
                            }
                          }
                        }
                      }
                    }
                  }
                },
                "times": {
                  "type": "object",
                  "required": ["recordedTime"]
                }
              }
            }
          }
        }
      }
      """
    When user "Alice" lists the activities for folder "New Folder" of space "Personal" with depth "2" using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["value"],
        "properties": {
          "value": {
            "type": "array",
            "minItems": 3,
            "maxItems": 3,
            "uniqueItems": true,
            "items": {
              "oneOf": [
                {
                  "type": "object",
                  "required": ["id", "template", "times"],
                  "properties": {
                    "template": {
                      "type": "object",
                      "required": ["message", "variables"],
                      "properties": {
                        "message": {
                          "const": "{user} added {resource} to {space}"
                        },
                        "variables": {
                          "type": "object",
                          "required": ["resource", "space", "user"],
                          "properties": {
                            "resource": {
                              "type": "object",
                              "required": ["id", "name"],
                              "properties": {
                                "name": {
                                  "const": "New Folder"
                                }
                              }
                            }
                          }
                        }
                      }
                    },
                    "times": {
                      "type": "object",
                      "required": ["recordedTime"]
                    }
                  }
                },
                {
                  "type": "object",
                  "required": ["id", "template", "times"],
                  "properties": {
                    "template": {
                      "type": "object",
                      "required": ["message", "variables"],
                      "properties": {
                        "message": {
                          "const": "{user} added {resource} to {space}"
                        },
                        "variables": {
                          "type": "object",
                          "required": ["resource", "space", "user"],
                          "properties": {
                            "resource": {
                              "type": "object",
                              "required": ["id", "name"],
                              "properties": {
                                "name": {
                                  "const": "Sub Folder"
                                }
                              }
                            }
                          }
                        }
                      }
                    },
                    "times": {
                      "type": "object",
                      "required": ["recordedTime"]
                    }
                  }
                },
                {
                  "type": "object",
                  "required": ["id", "template", "times"],
                  "properties": {
                    "template": {
                      "type": "object",
                      "required": ["message", "variables"],
                      "properties": {
                        "message": {
                          "const": "{user} added {resource} to {space}"
                        },
                        "variables": {
                          "type": "object",
                          "required": ["resource", "space", "user"],
                          "properties": {
                            "resource": {
                              "type": "object",
                              "required": ["id", "name"],
                              "properties": {
                                "name": {
                                  "const": "textfile0.txt"
                                }
                              }
                            }
                          }
                        }
                      }
                    },
                    "times": {
                      "type": "object",
                      "required": ["recordedTime"]
                    }
                  }
                }
              ]
            }
          }
        }
      }
      """
