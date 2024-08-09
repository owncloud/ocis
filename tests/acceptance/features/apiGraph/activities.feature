Feature: check activities
  As a user
  I want to check who made which changes to files
  So that I can track modifications


  Scenario: check activities after uploading a file
    Given user "Alice" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file with content "ownCloud test text file 0" to "/textfile0.txt"
    When user "Alice" checks the activities for file "textfile0.txt" in space "Personal" using the Graph API
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


  Scenario: check move activity for a file
    Given user "Alice" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file with content "ownCloud test text file 0" to "/textfile0.txt"
    And user "Alice" has created folder "/New Folder"
    When user "Alice" moves file "textfile0.txt" to "/New Folder/textfile0.txt" using the WebDAV API
    And user "Alice" checks the activities in space "Personal" using the Graph API
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
                                  "const": ""
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
                                  "const": "New Folder"
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
