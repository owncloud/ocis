Feature: an user gets the resources shared to them
  As a user
  I want to get resources shared with me
  So that I can know about what resources I have access to

  https://owncloud.dev/libre-graph-api/#/me.drive/ListSharedWithMe

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |


  Scenario: user lists the file shared with them
    Given user "Alice" has uploaded file with content "hello world" to "/textfile0.txt"
    And user "Alice" has sent the following share invitation:
      | resource        | textfile0.txt |
      | space           | Personal      |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | Viewer        |
    When user "Brian" lists the shares shared with him using the Graph API
    Then the HTTP status code should be "200"
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
          "items": {
            "type": "object",
            "required": [
              "createdBy",
              "eTag",
              "file",
              "id",
              "lastModifiedDateTime",
              "name",
              "parentReference",
              "remoteItem"
            ],
            "properties": {
              "createdBy": {
                "type": "object",
                "required": [
                  "user"
                ],
                "properties": {
                  "user": {
                    "type": "object",
                    "required": ["displayName", "id"],
                    "properties": {
                      "displayName": {
                        "type": "string",
                        "enum": ["Alice Hansen"]
                      },
                      "id": {
                        "type": "string",
                        "pattern": "^%user_id_pattern%$"
                      }
                    }
                  }
                }
              },
              "eTag": {
                "type": "string",
                "pattern": "%eTag%"
              },
              "file": {
                "type": "object",
                "required": ["mimeType"],
                "properties": {
                  "mimeType": {
                    "type": "string",
                    "enum": ["text/plain"]
                  }
                }
              },
              "id": {
                "type": "string",
                "pattern": "^%share_id_pattern%$"
              },
              "name": {
                "type": "string",
                "enum": [
                  "textfile0.txt"
                ]
              },
              "parentReference": {
                "type": "object",
                "required": [
                  "driveId",
                  "driveType"
                ],
                "properties": {
                  "driveId": {
                    "type": "string",
                    "pattern": "^%space_id_pattern%$"
                  },
                  "driveType" : {
                    "type": "string",
                    "enum": ["virtual"]
                  }
                }
              },
              "remoteItem": {
                "type": "object",
                "required": [
                  "createdBy",
                  "eTag",
                  "file",
                  "id",
                  "lastModifiedDateTime",
                  "name",
                  "parentReference",
                  "permissions",
                  "size"
                ],
                "properties": {
                  "createdBy": {
                    "type": "object",
                    "required": [
                      "user"
                    ],
                    "properties": {
                      "user": {
                        "type": "object",
                        "required": [
                          "id",
                          "displayName"
                        ],
                        "properties": {
                          "id": {
                            "type": "string",
                            "pattern": "^%user_id_pattern%$"
                          },
                          "displayName": {
                            "type": "string",
                            "enum": [
                              "Alice Hansen"
                            ]
                          }
                        }
                      }
                    }
                  },
                  "eTag": {
                    "type": "string",
                    "pattern": "%eTag%"
                  },
                  "file": {
                    "type": "object",
                    "required": [
                      "mimeType"
                    ],
                    "properties": {
                      "mimeType": {
                        "type": "string",
                        "pattern": "text/plain"
                      }
                    }
                  },
                  "id": {
                    "type": "string",
                    "pattern": "^%file_id_pattern%$"
                  },
                  "name": {
                    "type": "string",
                    "enum": [
                      "textfile0.txt"
                    ]
                  },
                  "permissions": {
                    "type": "array",
                    "items": [
                      {
                        "type": "object",
                        "required": [
                           "grantedToV2",
                           "id",
                           "invitation"
                         ],
                         "properties": {
                           "id": {
                             "type": "string"
                           },
                           "grantedToV2": {
                             "type": "object",
                             "required": [
                               "user"
                             ],
                             "properties": {
                               "user": {
                                 "type": "object",
                                 "properties": {
                                   "displayName": {
                                     "type": "string"
                                   },
                                   "id": {
                                     "type": "string"
                                   }
                                 },
                                 "required": [
                                   "displayName",
                                   "id"
                                 ]
                               }
                             }
                           },
                           "invitation": {
                             "type": "object",
                             "properties": {
                               "invitedBy": {
                                 "type": "object",
                                 "properties": {
                                   "user": {
                                     "type": "object",
                                     "properties": {
                                       "displayName": {
                                         "type": "string"
                                       },
                                       "id": {
                                         "type": "string"
                                       }
                                     },
                                     "required": [
                                       "displayName",
                                       "id"
                                     ]
                                   }
                                 },
                                 "required": [
                                   "user"
                                 ]
                               }
                             },
                             "required": [
                               "invitedBy"
                             ]
                           },
                           "@libre.graph.permissions.actions": {
                             "type": "array",
                             "items": [
                               {
                                 "type": "string"
                               }
                             ]
                           },
                           "roles": {
                             "type": "array",
                             "items": [
                               {
                                 "type": "string"
                               }
                             ]
                           }
                         }
                       }
                     ]
                   },
                   "size": {
                     "type": "number",
                     "enum": [
                       11
                     ]
                   }
                 }
               }
             }
           }
         }
       }
     }
    """


  Scenario: user lists the folder shared with them
    Given user "Alice" has created folder "folder"
    And user "Alice" has sent the following share invitation:
      | resource        | folder   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Viewer   |
    When user "Brian" lists the shares shared with him using the Graph API
    Then the HTTP status code should be "200"
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
          "items": {
            "type": "object",
            "required": [
              "createdBy",
              "eTag",
              "folder",
              "id",
              "lastModifiedDateTime",
              "name",
              "parentReference",
              "remoteItem"
            ],
            "properties": {
              "createdBy": {
                "type": "object",
                "required": [
                  "user"
                ],
                "properties": {
                  "user": {
                    "type": "object",
                    "required": ["displayName", "id"],
                    "properties": {
                      "displayName": {
                        "type": "string",
                        "enum": ["Alice Hansen"]
                      },
                      "id": {
                        "type": "string",
                        "pattern": "^%user_id_pattern%$"
                      }
                    }
                  }
                }
              },
              "eTag": {
                "type": "string",
                "pattern": "%eTag%"
              },
              "folder": {
                "type": "object",
                "required": [],
                "properties": {}
              },
              "id": {
                "type": "string",
                "pattern": "^%share_id_pattern%$"
              },
              "name": {
                "type": "string",
                "enum": [
                  "folder"
                ]
              },
              "parentReference": {
                "type": "object",
                "required": [
                  "driveId",
                  "driveType"
                ],
                "properties": {
                  "driveId": {
                    "type": "string",
                    "pattern": "^%space_id_pattern%$"
                  },
                  "driveType" : {
                    "type": "string",
                    "enum": ["virtual"]
                  }
                }
              },
              "remoteItem": {
                "type": "object",
                "required": [
                  "createdBy",
                  "eTag",
                  "folder",
                  "id",
                  "lastModifiedDateTime",
                  "name",
                  "parentReference",
                  "permissions"
                ],
                "properties": {
                  "createdBy": {
                    "type": "object",
                    "required": [
                      "user"
                    ],
                    "properties": {
                      "user": {
                        "type": "object",
                        "required": [
                          "id",
                          "displayName"
                        ],
                        "properties": {
                          "id": {
                            "type": "string",
                            "pattern": "^%user_id_pattern%$"
                          },
                          "displayName": {
                            "type": "string",
                            "enum": [
                              "Alice Hansen"
                            ]
                          }
                        }
                      }
                    }
                  },
                  "eTag": {
                    "type": "string",
                    "pattern": "%eTag%"
                  },
                  "file": {
                    "type": "object",
                    "required": [
                      "mimeType"
                    ],
                    "properties": {
                      "mimeType": {
                        "type": "string",
                        "pattern": "text/plain"
                      }
                    }
                  },
                  "id": {
                    "type": "string",
                    "pattern": "^%file_id_pattern%$"
                  },
                  "name": {
                    "type": "string",
                    "enum": [
                      "folder"
                    ]
                  },
                  "permissions": {
                    "type": "array",
                    "items": [
                      {
                        "type": "object",
                        "required": [
                           "grantedToV2",
                           "id",
                           "invitation"
                         ],
                         "properties": {
                           "id": {
                             "type": "string"
                           },
                           "grantedToV2": {
                             "type": "object",
                             "required": [
                               "user"
                             ],
                             "properties": {
                               "user": {
                                 "type": "object",
                                 "properties": {
                                   "displayName": {
                                     "type": "string"
                                   },
                                   "id": {
                                     "type": "string"
                                   }
                                 },
                                 "required": [
                                   "displayName",
                                   "id"
                                 ]
                               }
                             }
                           },
                           "invitation": {
                             "type": "object",
                             "properties": {
                               "invitedBy": {
                                 "type": "object",
                                 "properties": {
                                   "user": {
                                     "type": "object",
                                     "properties": {
                                       "displayName": {
                                         "type": "string"
                                       },
                                       "id": {
                                         "type": "string"
                                       }
                                     },
                                     "required": [
                                       "displayName",
                                       "id"
                                     ]
                                   }
                                 },
                                 "required": [
                                   "user"
                                 ]
                               }
                             },
                             "required": [
                               "invitedBy"
                             ]
                           },
                           "@libre.graph.permissions.actions": {
                             "type": "array",
                             "items": [
                               {
                                 "type": "string"
                               }
                             ]
                           },
                           "roles": {
                             "type": "array",
                             "items": [
                               {
                                 "type": "string"
                               }
                             ]
                           }
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
    """
