Feature: check activities
  As a user
  I want to check who made which changes to files
  So that I can track modifications

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files

  @issue-9712
  Scenario: check activities after uploading a file and a folder
    Given user "Alice" has uploaded file with content "ownCloud test text file 0" to "/textfile.txt"
    And user "Alice" has created folder "/FOLDER"
    When user "Alice" lists the activities of file "textfile.txt" from space "Personal" using the Graph API
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
          }
        }
      }
      """
    When user "Alice" lists the activities of folder "FOLDER" from space "Personal" using the Graph API
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
                      "required": ["folder","resource","user"],
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
                              "pattern": "%file_id_pattern%"
                            },
                            "name": {
                              "const": "FOLDER"
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
          }
        }
      }
      """

  @issue-9712
  Scenario: check activities after deleting a file and a folder
    Given user "Alice" has uploaded file with content "ownCloud test text file 0" to "/textfile.txt"
    And user "Alice" has created folder "/FOLDER"
    And user "Alice" has deleted file "textfile.txt"
    And user "Alice" has deleted folder "FOLDER"
    When user "Alice" lists the activities of space "Personal" using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["value"],
        "properties": {
          "value": {
            "type": "array",
            "minItems": 4,
            "maxItems": 4,
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
                        },
                        "variables": {
                          "type": "object",
                          "required": ["resource", "folder", "user"],
                          "properties": {
                            "resource": {
                              "type": "object",
                              "required": ["id", "name"],
                              "properties": {
                                "name": {
                                  "const": "textfile.txt"
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
                    "template": {
                      "type": "object",
                      "required": ["message","variables"],
                      "properties": {
                        "message": {
                          "const": "{user} added {resource} to {folder}"
                        },
                        "variables": {
                          "type": "object",
                          "required": ["resource", "folder", "user"],
                          "properties": {
                            "resource": {
                              "type": "object",
                              "required": ["id", "name"],
                              "properties": {
                                "name": {
                                  "const": "FOLDER"
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
                          "const": "{user} deleted {resource} from {folder}"
                        },
                        "variables": {
                          "type": "object",
                          "required": ["resource","folder","user"],
                          "properties": {
                            "resource": {
                              "type": "object",
                              "required": ["id","name"],
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
                            "folder": {
                              "type": "object",
                              "required": ["name"],
                              "properties": {
                                "name": {
                                  "const": "Alice Hansen"
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
                          "const": "{user} deleted {resource} from {folder}"
                        },
                        "variables": {
                          "type": "object",
                          "required": ["resource","folder","user"],
                          "properties": {
                            "resource": {
                              "type": "object",
                              "required": ["id","name"],
                              "properties": {
                                "id": {
                                  "type": "string",
                                  "pattern": "%file_id_pattern%"
                                },
                                "name": {
                                  "const": "FOLDER"
                                }
                              }
                            },
                            "folder": {
                              "type": "object",
                              "required": ["name"],
                              "properties": {
                                "name": {
                                  "const": "Alice Hansen"
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

  @issue-9712
  Scenario: check move activity for a file and a folder
    Given user "Alice" has uploaded file with content "ownCloud test text file 0" to "/textfile.txt"
    And user "Alice" has created folder "/FOLDER"
    And user "Alice" has created folder "/New Folder"
    And user "Alice" has moved file "textfile.txt" to "New Folder/textfile.txt"
    And user "Alice" has moved folder "FOLDER" to "New Folder/FOLDER"
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
    When user "Alice" lists the activities of folder "New Folder/FOLDER" from space "Personal" using the Graph API
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
                                "name": {
                                  "const": "FOLDER"
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
                                  "const": "FOLDER"
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

  @issue-9712
  Scenario: check rename activity for a file and a folder
    Given user "Alice" has uploaded file with content "ownCloud test text file 0" to "/textfile.txt"
    And user "Alice" has created folder "/FOLDER"
    And user "Alice" has moved file "textfile.txt" to "renamed.txt"
    And user "Alice" has moved folder "/FOLDER" to "RENAMED FOLDER"
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
    When user "Alice" lists the activities of folder "RENAMED FOLDER" from space "Personal" using the Graph API
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
                                  "const": "FOLDER"
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
                                  "const": "RENAMED FOLDER"
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

  @issue-9712
  Scenario: check activities of a folder
    Given user "Alice" has created folder "/New Folder"
    And user "Alice" has created folder "/New Folder/Folder"
    And user "Alice" has created folder "/New Folder/Sub Folder"
    And user "Alice" has uploaded file with content "ownCloud test text file 0" to "/New Folder/textfile.txt"
    And user "Alice" has moved file "/New Folder/textfile.txt" to "/New Folder/Sub Folder/textfile.txt"
    And user "Alice" has moved folder "/New Folder/Folder" to "/New Folder/Sub Folder/Folder"
    And user "Alice" has moved file "/New Folder/Sub Folder/textfile.txt" to "/New Folder/Sub Folder/renamed.txt"
    And user "Alice" has moved folder "/New Folder/Sub Folder/Folder" to "/New Folder/Sub Folder/Renamed Folder"
    And user "Alice" has deleted file "/New Folder/Sub Folder/renamed.txt"
    And user "Alice" has deleted folder "/New Folder/Sub Folder/Renamed Folder"
    When user "Alice" lists the activities of folder "/New Folder" from space "Personal" using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["value"],
        "properties": {
          "value": {
            "type": "array",
            "minItems": 10,
            "maxItems": 10,
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
                              "required": ["id", "name"],
                              "properties": {
                                "id": {
                                  "type": "string"
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
                  "required": ["id", "template", "times"],
                  "properties": {
                    "template": {
                      "type": "object",
                      "required": ["message", "variables"],
                      "properties": {
                        "message": {
                          "const": "{user} added {resource} to {folder}"
                        },
                        "variables": {
                          "type": "object",
                          "required": ["resource", "folder", "user"],
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
                              "required": ["name"],
                              "properties": {
                                "name": {
                                  "const": "Folder"
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
                  "required": ["id", "template", "times"],
                  "properties": {
                    "template": {
                      "type": "object",
                      "required": ["message", "variables"],
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
                              "required": ["id", "name"],
                              "properties": {
                                "name": {
                                  "const": "Sub Folder"
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
                  "required": ["id", "template", "times"],
                  "properties": {
                    "template": {
                      "type": "object",
                      "required": ["message", "variables"],
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
                              "required": ["id", "name"],
                              "properties": {
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
                                  "const": "Sub Folder"
                                }
                              }
                            },
                            "resource": {
                              "type": "object",
                              "required": ["id", "name"],
                              "properties": {
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
                                  "const": "Sub Folder"
                                }
                              }
                            },
                            "resource": {
                              "type": "object",
                              "required": ["id", "name"],
                              "properties": {
                                "name": {
                                  "const": "Folder"
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
                  "required": ["id", "template", "times"],
                  "properties": {
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
                                "name": {
                                  "const": "textfile.txt"
                                }
                              }
                            },
                            "resource": {
                              "type": "object",
                              "required": ["id", "name"],
                              "properties": {
                                "name": {
                                  "const": "renamed.txt"
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
                  "required": ["id", "template", "times"],
                  "properties": {
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
                                "name": {
                                  "const": "Folder"
                                }
                              }
                            },
                            "resource": {
                              "type": "object",
                              "required": ["id", "name"],
                              "properties": {
                                "name": {
                                  "const": "Renamed Folder"
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
                  "required": ["id", "template", "times"],
                  "properties": {
                    "template": {
                      "type": "object",
                      "required": ["message", "variables"],
                      "properties": {
                        "message": {
                          "const": "{user} deleted {resource} from {folder}"
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
                                  "const": "Sub Folder"
                                }
                              }
                            },
                            "resource": {
                              "type": "object",
                              "required": ["id", "name"],
                              "properties": {
                                "name": {
                                  "const": "renamed.txt"
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
                  "required": ["id", "template", "times"],
                  "properties": {
                    "template": {
                      "type": "object",
                      "required": ["message", "variables"],
                      "properties": {
                        "message": {
                          "const": "{user} deleted {resource} from {folder}"
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
                                  "const": "Sub Folder"
                                }
                              }
                            },
                            "resource": {
                              "type": "object",
                              "required": ["id", "name"],
                              "properties": {
                                "name": {
                                  "const": "Renamed Folder"
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
                }
              ]
            }
          }
        }
      }
      """

  @issue-9856 @issue-10127 @skip
  Scenario: check activity message with different language
    Given user "Alice" has uploaded file with content "ownCloud test text file" to "textfile.txt"
    And user "Alice" has switched the system language to "de" using the Graph API
    When user "Alice" lists the activities of file "textfile.txt" from space "Personal" using the Graph API
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
                      "const": "{user} hat {resource} zu {folder} hinzugefügt"
                    },
                    "variables": {
                      "type": "object",
                      "required": ["resource","space","user"],
                      "properties": {
                        "resource": {
                          "type": "object",
                          "required": ["id","name"],
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
                        "space": {
                          "type": "object",
                          "required": ["id","name"],
                          "properties": {
                            "id": {
                              "type": "string",
                              "pattern": "^%user_id_pattern%!%user_id_pattern%$"
                            },
                            "name": {
                              "const": "Alice Hansen"
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
          }
        }
      }
      """

  @issue-9850
  Scenario: check activity with -1 depth filter
    Given user "Alice" has created folder "/New Folder"
    And user "Alice" has created folder "/New Folder/Sub Folder"
    And user "Alice" has uploaded file with content "ownCloud test text file 0" to "/New Folder/Sub Folder/textfile.txt"
    When user "Alice" lists the activities of folder "New Folder" from space "Personal" with depth "-1" using the Graph API
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
                              "required": ["id", "name"],
                              "properties": {
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
                              "required": ["id", "name"],
                              "properties": {
                                "name": {
                                  "const": "Sub Folder"
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
                                  "const": "Sub Folder"
                                }
                              }
                            },
                            "resource": {
                              "type": "object",
                              "required": ["id", "name"],
                              "properties": {
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

  @issue-9850
  Scenario: check activity with depth filter
    Given user "Alice" has created folder "/New Folder"
    And user "Alice" has created folder "/New Folder/Sub Folder"
    And user "Alice" has uploaded file with content "ownCloud test text file 0" to "/New Folder/Sub Folder/textfile.txt"
    When user "Alice" lists the activities of folder "New Folder" from space "Personal" with depth "1" using the Graph API
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
                              "required": ["id", "name"],
                              "properties": {
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
                              "required": ["id", "name"],
                              "properties": {
                                "name": {
                                  "const": "Sub Folder"
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

  @issue-9880
  Scenario: check activity with limit filter
    Given user "Alice" has created folder "/New Folder"
    And user "Alice" has created folder "/New Folder/Sub Folder"
    And user "Alice" has uploaded file with content "ownCloud test text file 0" to "/New Folder/Sub Folder/textfile.txt"
    When user "Alice" lists the activities of folder "New Folder" from space "Personal" with limit "2" using the Graph API
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
                              "required": ["id", "name"],
                              "properties": {
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
                              "required": ["id", "name"],
                              "properties": {
                                "name": {
                                  "const": "Sub Folder"
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

  @issue-9860
  Scenario: user tries to check activities of another user's file
    Given user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file with content "ownCloud test text file" to "textfile.txt"
    And user "Alice" has uploaded file with content "updated ownCloud test text file" to "textfile.txt"
    When user "Brian" tries to list the activities of file "textfile.txt" from space "Personal" owned by user "Alice" using the Graph API
    Then the HTTP status code should be "403"
