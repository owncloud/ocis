Feature: check activities
  As a user
  I want to check who made which changes to files
  So that I can track modifications

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files

  @issue-9712
  Scenario: check activities after uploading a file and a folder
    Given user "Alice" has uploaded file with content "ownCloud test text file 0" to "/textfile0.txt"
    And user "Alice" has created folder "/FOLDER"
    When user "Alice" lists the activities for file "textfile0.txt" of space "Personal" using the Graph API
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
                      "const": "{user} added {resource} to {space}"
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
                              "const": "textfile0.txt"
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
                              "const": "Alice"
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
    When user "Alice" lists the activities for folder "FOLDER" of space "Personal" using the Graph API
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
                      "const": "{user} added {resource} to {space}"
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
                              "const": "FOLDER"
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
                              "const": "Alice"
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
    Given user "Alice" has uploaded file with content "ownCloud test text file 0" to "/textfile0.txt"
    And user "Alice" has created folder "/FOLDER"
    And user "Alice" has deleted file "textfile0.txt"
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
                          "const": "{user} deleted {resource} from {space}"
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
                                  "const": "textfile0.txt"
                                }
                              }
                            },
                            "space": {
                              "type": "object",
                              "required": ["id","name"],
                              "properties": {
                                "id": {
                                  "type": "string",
                                  "pattern": "^%file_id_pattern%$"
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
                                  "const": "Alice"
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
                          "const": "{user} deleted {resource} from {space}"
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
                                  "const": "FOLDER"
                                }
                              }
                            },
                            "space": {
                              "type": "object",
                              "required": ["id","name"],
                              "properties": {
                                "id": {
                                  "type": "string",
                                  "pattern": "^%file_id_pattern%$"
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
                                  "const": "Alice"
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
    Given user "Alice" has uploaded file with content "ownCloud test text file 0" to "/textfile0.txt"
    And user "Alice" has created folder "/FOLDER"
    And user "Alice" has created folder "/New Folder"
    And user "Alice" has moved file "textfile0.txt" to "New Folder/textfile0.txt"
    And user "Alice" has moved folder "FOLDER" to "New Folder/FOLDER"
    When user "Alice" lists the activities for file "New Folder/textfile0.txt" of space "Personal" using the Graph API
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
                          "const": "{user} added {resource} to {space}"
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
                          "const": "{user} moved {resource} to {space}"
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
                                  "pattern": "^%file_id_pattern%$"
                                },
                                "name": {
                                  "const": "textfile0.txt"
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
                                  "const": "Alice"
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
    When user "Alice" lists the activities for folder "New Folder/FOLDER" of space "Personal" using the Graph API
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
                          "const": "{user} added {resource} to {space}"
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
                          "const": "{user} moved {resource} to {space}"
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
                                  "pattern": "^%file_id_pattern%$"
                                },
                                "name": {
                                  "const": "FOLDER"
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
                                  "const": "Alice"
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
    Given user "Alice" has uploaded file with content "ownCloud test text file 0" to "/textfile0.txt"
    And user "Alice" has created folder "/FOLDER"
    And user "Alice" has moved file "textfile0.txt" to "renamed.txt"
    And user "Alice" has moved folder "/FOLDER" to "RENAMED FOLDER"
    When user "Alice" lists the activities for file "renamed.txt" of space "Personal" using the Graph API
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
                          "const": "{user} added {resource} to {space}"
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
                                  "const": "textfile0.txt"
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
                                  "const": "Alice"
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
    When user "Alice" lists the activities for folder "RENAMED FOLDER" of space "Personal" using the Graph API
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
                          "const": "{user} added {resource} to {space}"
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
                                  "const": "Alice"
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
    And user "Alice" has uploaded file with content "ownCloud test text file 0" to "/New Folder/textfile0.txt"
    And user "Alice" has moved file "/New Folder/textfile0.txt" to "/New Folder/Sub Folder/textfile0.txt"
    And user "Alice" has moved folder "/New Folder/Folder" to "/New Folder/Sub Folder/Folder"
    And user "Alice" has moved file "/New Folder/Sub Folder/textfile0.txt" to "/New Folder/Sub Folder/renamed.txt"
    And user "Alice" has moved folder "/New Folder/Sub Folder/Folder" to "/New Folder/Sub Folder/Renamed Folder"
    And user "Alice" has deleted file "/New Folder/Sub Folder/renamed.txt"
    And user "Alice" has deleted folder "/New Folder/Sub Folder/Renamed Folder"
    When user "Alice" lists the activities for folder "/New Folder" of space "Personal" using the Graph API
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
                                  "const": "Folder"
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
                          "const": "{user} moved {resource} to {space}"
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
                          "const": "{user} moved {resource} to {space}"
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
                                  "const": "Folder"
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
                                  "const": "textfile0.txt"
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
                          "const": "{user} deleted {resource} from {space}"
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
                          "const": "{user} deleted {resource} from {space}"
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
                                  "const": "Renamed Folder"
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
