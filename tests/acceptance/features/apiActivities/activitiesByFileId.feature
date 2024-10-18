Feature: check activities
  As a user
  I want to check who made which changes to files using file-id
  So that I can track modifications

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files


  Scenario: check add activity for a file using file-id
    Given user "Alice" has uploaded file with content "ownCloud test text file 0" to "/textfile.txt"
    And we save it into "FILEID"
    And user "Alice" has created folder "/New Folder"
    And user "Alice" has copied a file "/textfile.txt" into "New Folder" inside space "Personal" using file-id path "/dav/spaces/<<FILEID>>"
    When user "Alice" lists the activities of file "New Folder" from space "Personal" using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["value"],
        "properties": {
          "value": {
            "type": "array",
            "minItems": 2,
            "maxItems": 2,
            "uniqueItems": true,
            "items": {
              "oneOf": [
                {
                  "type": "object",
                  "required": ["id","template","times"],
                  "properties": {
                    "id": {
                      "type": "string",
                      "pattern": "^%user_id_pattern%$"
                    },
                    "template": {
                      "type": "object",
                      "required": ["message","variables"],
                      "properties": {
                        "message": {
                          "const": "{user} added {resource} to {folder}"
                        },
                        "variables": {
                          "type": "object",
                          "required": ["folder", "resource", "user"],
                          "properties": {
                            "folder": {
                              "type": "object",
                              "required": ["name"],
                              "properties": {
                                "name": {
                                  "const": "Alice Hansen"
                                }
                              }
                            },
                            "resource": {
                              "type": "object",
                              "required": ["id","name"],
                              "properties": {
                                "id": {
                                  "type": "string",
                                  "pattern": "^%file_id_pattern%$"
                                },
                                "name": {
                                  "const": "New Folder"
                                }
                              }
                            },
                            "user": {
                              "type": "object",
                              "required": ["id","displayName"],
                              "properties": {
                                "id": {
                                  "type": "string",
                                  "pattern": "%user_id_pattern%"
                                },
                                "displayName": {
                                  "const": "Alice Hansen"
                                }
                              }
                            }
                          }
                        }
                      }
                    }
                  }
                },
                {
                  "type": "object",
                  "required": ["id","template","times"],
                  "properties": {
                    "id": {
                      "type": "string",
                      "pattern": "^%user_id_pattern%$"
                    },
                    "template": {
                      "type": "object",
                      "required": ["message","variables"],
                      "properties": {
                        "message": {
                          "const": "{user} added {resource} to {folder}"
                        },
                        "variables": {
                          "type": "object",
                          "required": ["folder", "resource", "user"],
                          "properties": {
                            "folder": {
                              "type": "object",
                              "required": ["name"],
                              "properties": {
                                "name": {
                                  "const": "New Folder"
                                }
                              }
                            },
                            "resource": {
                              "type": "object",
                              "required": ["id","name"],
                              "properties": {
                                "id": {
                                  "type": "string",
                                  "pattern": "^%file_id_pattern%$"
                                },
                                "name": {
                                  "const": "textfile.txt"
                                }
                              }
                            },
                            "user": {
                              "type": "object",
                              "required": ["id","displayName"],
                              "properties": {
                                "id": {
                                  "type": "string",
                                  "pattern": "%user_id_pattern%"
                                },
                                "displayName": {
                                  "const": "Alice Hansen"
                                }
                              }
                            }
                          }
                        }
                      }
                    },
                    "times": {
                      "type": "object",
                      "required": ["recordedTime"],
                      "properties": {
                        "recordedTime": {
                          "type": "string",
                          "format": "date-time"
                        }
                      }
                    }
                  }
                }
              ]
            }
          }
        }
      }
      """


  Scenario: check edit activity for a file using file-id
    Given user "Alice" has uploaded file with content "ownCloud test text file 0" to "/textfile.txt"
    And we save it into "FILEID"
    And user "Alice" has sent HTTP method "PUT" to URL "/dav/spaces/<<FILEID>>" with content "updated content"
    When user "Alice" lists the activities of file "/textfile.txt" from space "Personal" using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["value"],
        "properties": {
          "value": {
            "type": "array",
            "minItems": 2,
            "maxItems": 2,
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
                          "const": "{user} added {resource} to {folder}"
                        }
                      }
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
                          "const": "{user} updated {resource} in {folder}"
                        },
                        "variables": {
                          "type": "object",
                          "required": ["folder", "resource", "user"],
                          "properties": {
                            "folder": {
                              "type": "object",
                              "required": ["name"],
                              "properties": {
                                "id": {
                                  "type": "string",
                                  "pattern": "%user_id_pattern%"
                                },
                                "name": {
                                  "const": "Alice Hansen"
                                }
                              }
                            },
                            "resource": {
                              "type": "object",
                              "required": ["id", "name"],
                              "properties": {
                                "id": {
                                  "type": "string",
                                  "pattern": "%file_id_pattern%"
                                },
                                "name": {
                                  "const": "textfile.txt"
                                }
                              }
                            },
                            "user": {
                              "type": "object",
                              "required": ["id","displayName"],
                              "properties": {
                                "id": {
                                  "type": "string",
                                  "pattern": "%user_id_pattern%"
                                },
                                "displayName": {
                                  "const": "Alice Hansen"
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

  @issue-9744
  Scenario: check rename activity for a file using file-id
    Given user "Alice" has uploaded file with content "ownCloud test text file 0" to "/textfile.txt"
    And we save it into "FILEID"
    And user "Alice" has renamed a file "/textfile.txt" into "renamed.txt" inside space "Personal" using file-id path "/dav/spaces/<<FILEID>>"
    When user "Alice" lists the activities of file "renamed.txt" from space "Personal" using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["value"],
        "properties": {
          "value": {
            "type": "array",
            "minItems": 2,
            "maxItems": 2,
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
                          "const": "{user} added {resource} to {folder}"
                        }
                      }
                    }
                  }
                },
                {
                  "type": "object",
                  "required": ["id", "template", "times"],
                  "properties": {
                    "id": {
                      "type": "string",
                      "pattern": "^%user_id_pattern%$"
                    },
                    "template": {
                      "type": "object",
                      "required": ["message", "variables"],
                      "properties": {
                        "message": {
                          "const": "{user} renamed {oldResource} to {resource}"
                        },
                        "variables": {
                          "type": "object",
                          "required": ["oldResource", "resource", "user"],
                          "properties": {
                            "oldResource": {
                              "type": "object",
                              "required": ["id", "name"],
                              "properties": {
                                "id": {
                                  "const": ""
                                },
                                "name": {
                                  "const": "textfile.txt"
                                }
                              }
                            },
                            "resource": {
                              "type": "object",
                              "required": ["id", "name"],
                              "properties": {
                                "id": {
                                  "type": "string",
                                  "pattern": "^%file_id_pattern%$"
                                },
                                "name": {
                                  "const": "renamed.txt"
                                }
                              }
                            },
                            "user": {
                              "type": "object",
                              "required": ["id", "displayName"],
                              "properties": {
                                "id": {
                                  "type": "string",
                                  "pattern": "^%user_id_pattern%$"
                                },
                                "displayName": {
                                  "const": "Alice Hansen"
                                }
                              }
                            }
                          }
                        }
                      }
                    },
                    "times": {
                      "type": "object",
                      "required": ["recordedTime"],
                      "properties": {
                        "recordedTime": {
                          "type": "string",
                          "format": "date-time"
                        }
                      }
                    }
                  }
                }
              ]
            }
          }
        }
      }
      """


  Scenario: check move activity for a file using file-id
    Given user "Alice" has uploaded file with content "ownCloud test text file 0" to "/textfile.txt"
    And we save it into "FILEID"
    And user "Alice" has created folder "/New Folder"
    And user "Alice" has moved a file "/textfile.txt" into "New Folder" inside space "Personal" using file-id path "/dav/spaces/<<FILEID>>"
    When user "Alice" lists the activities of file "New Folder/textfile.txt" from space "Personal" using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["value"],
        "properties": {
          "value": {
            "type": "array",
            "minItems": 2,
            "maxItems": 2,
            "uniqueItems": true,
            "items": {
              "oneOf": [
                {
                  "type": "object",
                  "required": ["id","template","times"],
                  "properties": {
                    "template": {
                      "type": "object",
                      "required": ["message","variables"],
                      "properties": {
                        "message": {
                          "const": "{user} added {resource} to {folder}"
                        }
                      }
                    }
                  }
                },
                {
                  "type": "object",
                  "required": ["id","template","times"],
                  "properties": {
                    "id": {
                      "type": "string",
                      "pattern": "^%user_id_pattern%$"
                    },
                    "template": {
                      "type": "object",
                      "required": ["message","variables"],
                      "properties": {
                        "message": {
                          "const": "{user} moved {resource} to {folder}"
                        },
                        "variables": {
                          "type": "object",
                          "required": ["folder", "resource", "user"],
                          "properties": {
                            "folder": {
                              "type": "object",
                              "required": ["name"],
                              "properties": {
                                "name": {
                                  "const": "New Folder"
                                }
                              }
                            },
                            "resource": {
                              "type": "object",
                              "required": ["id","name"],
                              "properties": {
                                "id": {
                                  "type": "string",
                                  "pattern": "^%file_id_pattern%$"
                                },
                                "name": {
                                  "const": "textfile.txt"
                                }
                              }
                            },
                            "user": {
                              "type": "object",
                              "required": ["id","displayName"],
                              "properties": {
                                "id": {
                                  "type": "string",
                                  "pattern": "%user_id_pattern%"
                                },
                                "displayName": {
                                  "const": "Alice Hansen"
                                }
                              }
                            }
                          }
                        }
                      }
                    },
                    "times": {
                      "type": "object",
                      "required": ["recordedTime"],
                      "properties": {
                        "recordedTime": {
                          "type": "string",
                          "format": "date-time"
                        }
                      }
                    }
                  }
                }
              ]
            }
          }
        }
      }
      """
