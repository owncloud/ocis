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
              "@UI.Hidden",
              "@client.synchronize",
              "createdBy",
              "eTag",
              "file",
              "id",
              "lastModifiedDateTime",
              "name",
              "parentReference",
              "remoteItem",
              "size"
            ],
            "properties": {
              "@UI.Hidden": {
                "type": "boolean",
                "enum": [false]
              },
              "@client.synchronize": {
                "type": "boolean",
                "enum": [true]
              },
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
                "pattern": "%etag_pattern%"
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
                  "driveType",
                  "id"
                ],
                "properties": {
                  "driveId": {
                    "type": "string",
                    "pattern": "^%space_id_pattern%$"
                  },
                  "driveType": {
                    "type": "string",
                    "enum": ["virtual"]
                  },
                  "id": {
                    "type": "string",
                    "pattern": "^%file_id_pattern%$"
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
                    "pattern": "%etag_pattern%"
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
                  "parentReference": {
                    "type": "object",
                    "required": [
                      "driveId",
                      "driveType"
                    ],
                    "properties": {
                      "driveId": {
                        "type": "string",
                        "pattern": "^%file_id_pattern%$"
                      },
                      "driveType": {
                        "type": "string",
                        "enum": ["personal"]
                      }
                    }
                  },
                  "permissions": {
                    "type": "array",
                    "items": [
                      {
                        "type": "object",
                        "required": [
                          "grantedToV2",
                          "id",
                          "invitation",
                          "roles"
                        ],
                        "properties": {
                          "id": {
                            "type": "string",
                            "pattern": "^%permissions_id_pattern%$"
                          },
                          "grantedToV2": {
                            "type": "object",
                            "required": [
                              "user"
                            ],
                            "properties": {
                              "user": {
                                "type": "object",
                                "required": [
                                  "displayName",
                                  "id"
                                ],
                                "properties": {
                                  "displayName": {
                                    "type": "string",
                                    "enum": ["Brian Murphy"]
                                  },
                                  "id": {
                                    "type": "string",
                                    "pattern": "^%user_id_pattern%$"
                                  }
                                }
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
                                        "type": "string",
                                        "enum": ["Alice Hansen"]
                                      },
                                      "id": {
                                        "type": "string",
                                        "pattern": "^%user_id_pattern%$"
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
                          "roles": {
                            "type": "array",
                            "items": [
                              {
                                "type": "string",
                                "pattern": "^%role_id_pattern%$"
                              }
                            ]
                          }
                        }
                      }
                    ]
                  }
                }
              },
              "size": {
                "type": "number",
                "enum": [11]
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
              "@UI.Hidden",
              "@client.synchronize",
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
              "@UI.Hidden": {
                "type": "boolean",
                "enum": [false]
              },
              "@client.synchronize": {
                "type": "boolean",
                "enum": [true]
              },
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
                "pattern": "%etag_pattern%"
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
                  "driveType",
                  "id"
                ],
                "properties": {
                  "driveId": {
                    "type": "string",
                    "pattern": "^%space_id_pattern%$"
                  },
                  "driveType": {
                    "type": "string",
                    "enum": ["virtual"]
                  },
                  "id": {
                    "type": "string",
                    "pattern": "^%file_id_pattern%$"
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
                    "pattern": "%etag_pattern%"
                  },
                  "file": {
                    "type": "object",
                    "required": [
                      "mimeType"
                    ],
                    "properties": {
                      "mimeType": {
                        "type": "string",
                        "enum": ["text/plain"]
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
                  "parentReference": {
                    "type": "object",
                    "required": [
                      "driveId",
                      "driveType"
                    ],
                    "properties": {
                      "driveId": {
                        "type": "string",
                        "pattern": "^%file_id_pattern%$"
                      },
                      "driveType": {
                        "type": "string",
                        "enum": ["personal"]
                      }
                    }
                  },
                  "permissions": {
                    "type": "array",
                    "items": [
                      {
                        "type": "object",
                        "required": [
                          "grantedToV2",
                          "id",
                          "invitation",
                          "roles"
                        ],
                        "properties": {
                          "id": {
                            "type": "string",
                            "pattern": "^%permissions_id_pattern%$"
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
                                    "type": "string",
                                    "enum": ["Brian Murphy"]
                                  },
                                  "id": {
                                    "type": "string",
                                    "pattern": "^%user_id_pattern%$"
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
                                        "type": "string",
                                        "enum": ["Alice Hansen"]
                                      },
                                      "id": {
                                        "type": "string",
                                        "pattern": "^%user_id_pattern%$"
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
                          "roles": {
                            "type": "array",
                            "items": [
                              {
                                "type": "string",
                                "pattern": "^%role_id_pattern%$"
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


  Scenario: sharer shares a file to a group and to a user who is in the shared group
    Given the administrator has assigned the role "Admin" to user "Alice" using the Graph API
    And user "Alice" has uploaded file with content "hello" to "textfile0.txt"
    And user "Alice" has created a group "grp1" using the Graph API
    And user "Brian" has been added to group "grp1"
    And user "Alice" has sent the following share invitation:
      | resource        | textfile0.txt |
      | space           | Personal      |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | Viewer        |
    And user "Alice" has sent the following share invitation:
      | resource        | textfile0.txt |
      | space           | Personal      |
      | sharee          | grp1          |
      | shareType       | group         |
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
              "@UI.Hidden",
              "@client.synchronize",
              "createdBy",
              "eTag",
              "file",
              "id",
              "lastModifiedDateTime",
              "name",
              "parentReference",
              "remoteItem",
              "size"
            ],
            "properties": {
              "@UI.Hidden":{
                "type": "boolean",
                "enum": [false]
              },
              "@client.synchronize":{
                "type": "boolean",
                "enum": [true]
              },
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
                "pattern": "%etag_pattern%"
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
                  "driveType",
                  "id"
                ],
                "properties": {
                  "driveId": {
                    "type": "string",
                    "pattern": "^%space_id_pattern%$"
                  },
                  "driveType" : {
                    "type": "string",
                    "enum": ["virtual"]
                  },
                  "id" : {
                    "type": "string",
                    "pattern": "%space_id_pattern%"
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
                    "pattern": "%etag_pattern%"
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
                    "pattern": "^%file_id_pattern%$"
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
                        "pattern": "%space_id_pattern%"
                      },
                      "driveType" : {
                        "type": "string",
                        "enum": ["personal"]
                      }
                    }
                  },
                  "permissions": {
                    "type": "array",
                    "items": {
                      "oneOf": [
                        {
                          "type": "object",
                          "required": [
                            "grantedToV2",
                            "id",
                            "invitation",
                            "roles"
                          ],
                          "properties": {
                            "grantedToV2": {
                              "type": "object",
                              "required": [
                                "group"
                              ],
                              "properties":{
                                "group": {
                                  "type": "object",
                                  "required": [
                                    "displayName",
                                    "id"
                                  ],
                                  "properties": {
                                    "displayName": {
                                      "type": "string",
                                      "enum": ["grp1"]
                                    },
                                    "id": {
                                      "type": "string",
                                      "pattern": "^%user_id_pattern%$"
                                    }
                                  }
                                }
                              }
                            },
                            "id": {
                              "type": "string",
                              "pattern": "^%permissions_id_pattern%$"
                            },
                            "invitation": {
                              "type": "object",
                              "required": [
                                "invitedBy"
                              ],
                              "properties": {
                                "user":{
                                  "type": "object",
                                  "required": [
                                    "displayName",
                                    "id"
                                  ],
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
                            "roles": {
                              "type": "array",
                              "items": [
                                {
                                  "type": "string",
                                  "pattern": "^%role_id_pattern%$"
                                }
                              ]
                            }
                          }
                        },
                        {
                          "type": "object",
                          "required": [
                            "grantedToV2",
                            "id",
                            "invitation",
                            "roles"
                          ],
                          "properties": {
                            "grantedToV2": {
                              "type": "object",
                              "required": [
                                "user"
                              ],
                              "properties":{
                                "user": {
                                  "type": "object",
                                  "required": [
                                    "displayName",
                                    "id"
                                  ],
                                  "properties": {
                                    "displayName": {
                                      "type": "string",
                                      "enum": ["Brian Murphy"]
                                    },
                                    "id": {
                                      "type": "string",
                                      "pattern": "^%user_id_pattern%$"
                                    }
                                  }
                                }
                              }
                            },
                            "id": {
                              "type": "string",
                              "pattern": "^%permissions_id_pattern%$"
                            },
                            "invitation": {
                              "type": "object",
                              "required": [
                                "invitedBy"
                              ],
                              "properties": {
                                "user":{
                                  "type": "object",
                                  "required": [
                                    "displayName",
                                    "id"
                                  ],
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
                            "roles": {
                              "type": "array",
                              "items": [
                                {
                                  "type": "string",
                                  "pattern": "^%role_id_pattern%$"
                                }
                              ]
                            }
                          }
                        }
                      ]
                    }
                  },
                  "size": {
                    "type": "number",
                    "enum": [
                      5
                    ]
                  }
                }
              },
              "size": {
                "type": "number",
                "enum": [
                  5
                ]
              }
            }
          }
        }
      }
    }
    """

  @issues-8314
  Scenario: sharer shares a file in project space to a group and to a user who is in the shared group file
    Given using spaces DAV path
    And the administrator has assigned the role "Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "new-space" with content "hello world" to "textfile0.txt"
    And user "Alice" has created a group "grp1" using the Graph API
    And user "Brian" has been added to group "grp1"
    And user "Alice" has sent the following share invitation:
      | resource        | textfile0.txt |
      | space           | new-space     |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | Viewer        |
    And user "Alice" has sent the following share invitation:
      | resource        | textfile0.txt |
      | space           | new-space     |
      | sharee          | grp1          |
      | shareType       | group         |
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
              "@UI.Hidden",
              "@client.synchronize",
              "eTag",
              "file",
              "id",
              "lastModifiedDateTime",
              "name",
              "parentReference",
              "remoteItem",
              "size"
            ],
            "properties": {
              "@UI.Hidden":{
                "type": "boolean",
                "enum": [false]
              },
              "@client.synchronize":{
                "type": "boolean",
                "enum": [true]
              },
              "eTag": {
                "type": "string",
                "pattern": "%etag_pattern%"
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
                  "driveType",
                  "id"
                ],
                "properties": {
                  "driveId": {
                    "type": "string",
                    "pattern": "^%space_id_pattern%$"
                  },
                  "driveType" : {
                    "type": "string",
                    "enum": ["virtual"]
                  },
                  "id" : {
                    "type": "string",
                    "pattern": "%space_id_pattern%"
                  }
                }
              },
              "remoteItem": {
                "type": "object",
                "required": [
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
                  "eTag": {
                    "type": "string",
                    "pattern": "%etag_pattern%"
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
                    "pattern": "^%file_id_pattern%$"
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
                        "pattern": "%space_id_pattern%"
                      },
                      "driveType" : {
                        "type": "string",
                        "enum": ["project"]
                      }
                    }
                  },
                  "permissions": {
                    "type": "array",
                    "items": {
                      "oneOf": [
                        {
                          "type": "object",
                          "required": [
                            "grantedToV2",
                            "id",
                            "invitation",
                            "roles"
                          ],
                          "properties": {
                            "grantedToV2": {
                              "type": "object",
                              "required": [
                                "group"
                              ],
                              "properties":{
                                "group": {
                                  "type": "object",
                                  "required": [
                                    "displayName",
                                    "id"
                                  ],
                                  "properties": {
                                    "displayName": {
                                      "type": "string",
                                      "enum": ["grp1"]
                                    },
                                    "id": {
                                      "type": "string",
                                      "pattern": "^%user_id_pattern%$"
                                    }
                                  }
                                }
                              }
                            },
                            "id": {
                              "type": "string",
                              "pattern": "^%permissions_id_pattern%$"
                            },
                            "invitation": {
                              "type": "object",
                              "required": [
                                "invitedBy"
                              ],
                              "properties": {
                                "user":{
                                  "type": "object",
                                  "required": [
                                    "displayName",
                                    "id"
                                  ],
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
                            "roles": {
                              "type": "array",
                              "items": [
                                {
                                  "type": "string",
                                  "pattern": "^%role_id_pattern%$"
                                }
                              ]
                            }
                          }
                        },
                        {
                          "type": "object",
                          "required": [
                            "grantedToV2",
                            "id",
                            "invitation",
                            "roles"
                          ],
                          "properties": {
                            "grantedToV2": {
                              "type": "object",
                              "required": [
                                "user"
                              ],
                              "properties":{
                                "user": {
                                  "type": "object",
                                  "required": [
                                    "displayName",
                                    "id"
                                  ],
                                  "properties": {
                                    "displayName": {
                                      "type": "string",
                                      "enum": ["Brian Murphy"]
                                    },
                                    "id": {
                                      "type": "string",
                                      "pattern": "^%user_id_pattern%$"
                                    }
                                  }
                                }
                              }
                            },
                            "id": {
                              "type": "string",
                              "pattern": "^%permissions_id_pattern%$"
                            },
                            "invitation": {
                              "type": "object",
                              "required": [
                                "invitedBy"
                              ],
                              "properties": {
                                "user":{
                                  "type": "object",
                                  "required": [
                                    "displayName",
                                    "id"
                                  ],
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
                            "roles": {
                              "type": "array",
                              "items": [
                                {
                                  "type": "string",
                                  "pattern": "^%role_id_pattern%$"
                                }
                              ]
                            }
                          }
                        }
                      ]
                    }
                  },
                  "size": {
                    "type": "number",
                    "enum": [
                      11
                    ]
                  }
                }
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
    """


  Scenario: sharer shares a folder to a group and to a user who is in the shared group folder
    Given the administrator has assigned the role "Admin" to user "Alice" using the Graph API
    And  user "Alice" has created folder "folder"
    And user "Alice" has created a group "grp1" using the Graph API
    And user "Brian" has been added to group "grp1"
    And user "Alice" has sent the following share invitation:
      | resource        | folder   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Viewer   |
    And user "Alice" has sent the following share invitation:
      | resource        | folder   |
      | space           | Personal |
      | sharee          | grp1     |
      | shareType       | group    |
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
              "@UI.Hidden",
              "@client.synchronize",
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
              "@UI.Hidden":{
                "type": "boolean",
                "enum": [false]
              },
              "@client.synchronize":{
                "type": "boolean",
                "enum": [true]
              },
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
                "pattern": "%etag_pattern%"
              },
              "folder": {},
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
                  "driveType",
                  "id"
                ],
                "properties": {
                  "driveId": {
                    "type": "string",
                    "pattern": "^%space_id_pattern%$"
                  },
                  "driveType" : {
                    "type": "string",
                    "enum": ["virtual"]
                  },
                  "id" : {
                    "type": "string",
                    "pattern": "%space_id_pattern%"
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
                    "pattern": "%etag_pattern%"
                  },
                  "file": {},
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
                  "parentReference": {
                    "type": "object",
                    "required": [
                      "driveId",
                      "driveType"
                    ],
                    "properties": {
                      "driveId": {
                        "type": "string",
                        "pattern": "%space_id_pattern%"
                      },
                      "driveType" : {
                        "type": "string",
                        "enum": ["personal"]
                      }
                    }
                  },
                  "permissions": {
                    "type": "array",
                    "items": {
                      "oneOf": [
                        {
                          "type": "object",
                          "required": [
                            "grantedToV2",
                            "id",
                            "invitation",
                            "roles"
                          ],
                          "properties": {
                            "grantedToV2": {
                              "type": "object",
                              "required": [
                                "group"
                              ],
                              "properties":{
                                "group": {
                                  "type": "object",
                                  "required": [
                                    "displayName",
                                    "id"
                                  ],
                                  "properties": {
                                    "displayName": {
                                      "type": "string",
                                      "enum": ["grp1"]
                                    },
                                    "id": {
                                      "type": "string",
                                      "pattern": "^%user_id_pattern%$"
                                    }
                                  }
                                }
                              }
                            },
                            "id": {
                              "type": "string",
                              "pattern": "^%permissions_id_pattern%$"
                            },
                            "invitation": {
                              "type": "object",
                              "required": [
                                "invitedBy"
                              ],
                              "properties": {
                                "user":{
                                  "type": "object",
                                  "required": [
                                    "displayName",
                                    "id"
                                  ],
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
                            "roles": {
                              "type": "array",
                              "items": [
                                {
                                  "type": "string",
                                  "pattern": "^%role_id_pattern%$"
                                }
                              ]
                            }
                          }
                        },
                        {
                          "type": "object",
                          "required": [
                            "grantedToV2",
                            "id",
                            "invitation",
                            "roles"
                          ],
                          "properties": {
                            "grantedToV2": {
                              "type": "object",
                              "required": [
                                "user"
                              ],
                              "properties":{
                                "user": {
                                  "type": "object",
                                  "required": [
                                    "displayName",
                                    "id"
                                  ],
                                  "properties": {
                                    "displayName": {
                                      "type": "string",
                                      "enum": ["Brian Murphy"]
                                    },
                                    "id": {
                                      "type": "string",
                                      "pattern": "^%user_id_pattern%$"
                                    }
                                  }
                                }
                              }
                            },
                            "id": {
                              "type": "string",
                              "pattern": "^%permissions_id_pattern%$"
                            },
                            "invitation": {
                              "type": "object",
                              "required": [
                                "invitedBy"
                              ],
                              "properties": {
                                "user":{
                                  "type": "object",
                                  "required": [
                                    "displayName",
                                    "id"
                                  ],
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
                            "roles": {
                              "type": "array",
                              "items": [
                                {
                                  "type": "string",
                                  "pattern": "^%role_id_pattern%$"
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
    }
    """

  @issues-8314
  Scenario: sharer shares a folder in project space to a group and to a user who is in the shared group
    Given using spaces DAV path
    And the administrator has assigned the role "Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the Graph API
    And user "Alice" has created a group "grp1" using the Graph API
    And user "Alice" has created a folder "folder" in space "new-space"
    And user "Brian" has been added to group "grp1"
    And user "Alice" has sent the following share invitation:
      | resource        | folder    |
      | space           | new-space |
      | sharee          | Brian     |
      | shareType       | user      |
      | permissionsRole | Viewer    |
    And user "Alice" has sent the following share invitation:
      | resource        | folder    |
      | space           | new-space |
      | sharee          | grp1      |
      | shareType       | group     |
      | permissionsRole | Viewer    |
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
              "@UI.Hidden",
              "@client.synchronize",
              "eTag",
              "folder",
              "id",
              "lastModifiedDateTime",
              "name",
              "parentReference",
              "remoteItem"
            ],
            "properties": {
              "@UI.Hidden":{
                "type": "boolean",
                "enum": [false]
              },
              "@client.synchronize":{
                "type": "boolean",
                "enum": [true]
              },
              "eTag": {
                "type": "string",
                "pattern": "%etag_pattern%"
              },
              "folder": {},
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
                  "driveType",
                  "id"
                ],
                "properties": {
                  "driveId": {
                    "type": "string",
                    "pattern": "^%space_id_pattern%$"
                  },
                  "driveType" : {
                    "type": "string",
                    "enum": ["virtual"]
                  },
                  "id" : {
                    "type": "string",
                    "pattern": "%space_id_pattern%"
                  }
                }
              },
              "remoteItem": {
                "type": "object",
                "required": [
                  "eTag",
                  "folder",
                  "id",
                  "lastModifiedDateTime",
                  "name",
                  "parentReference",
                  "permissions"
                ],
                "properties": {
                  "eTag": {
                    "type": "string",
                    "pattern": "%etag_pattern%"
                  },
                  "folder": {},
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
                  "parentReference": {
                    "type": "object",
                    "required": [
                      "driveId",
                      "driveType"
                    ],
                    "properties": {
                      "driveId": {
                        "type": "string",
                        "pattern": "%space_id_pattern%"
                      },
                      "driveType" : {
                        "type": "string",
                        "enum": ["project"]
                      }
                    }
                  },
                  "permissions": {
                    "type": "array",
                    "items": {
                      "oneOf": [
                        {
                          "type": "object",
                          "required": [
                            "grantedToV2",
                            "id",
                            "invitation",
                            "roles"
                          ],
                          "properties": {
                            "grantedToV2": {
                              "type": "object",
                              "required": [
                                "group"
                              ],
                              "properties":{
                                "group": {
                                  "type": "object",
                                  "required": [
                                    "displayName",
                                    "id"
                                  ],
                                  "properties": {
                                    "displayName": {
                                      "type": "string",
                                      "enum": ["grp1"]
                                    },
                                    "id": {
                                      "type": "string",
                                      "pattern": "^%user_id_pattern%$"
                                    }
                                  }
                                }
                              }
                            },
                            "id": {
                              "type": "string",
                              "pattern": "^%permissions_id_pattern%$"
                            },
                            "invitation": {
                              "type": "object",
                              "required": [
                                "invitedBy"
                              ],
                              "properties": {
                                "user":{
                                  "type": "object",
                                  "required": [
                                    "displayName",
                                    "id"
                                  ],
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
                            "roles": {
                              "type": "array",
                              "items": [
                                {
                                  "type": "string",
                                  "pattern": "^%role_id_pattern%$"
                                }
                              ]
                            }
                          }
                        },
                        {
                          "type": "object",
                          "required": [
                            "grantedToV2",
                            "id",
                            "invitation",
                            "roles"
                          ],
                          "properties": {
                            "grantedToV2": {
                              "type": "object",
                              "required": [
                                "user"
                              ],
                              "properties":{
                                "user": {
                                  "type": "object",
                                  "required": [
                                    "displayName",
                                    "id"
                                  ],
                                  "properties": {
                                    "displayName": {
                                      "type": "string",
                                      "enum": ["Brian Murphy"]
                                    },
                                    "id": {
                                      "type": "string",
                                      "pattern": "^%user_id_pattern%$"
                                    }
                                  }
                                }
                              }
                            },
                            "id": {
                              "type": "string",
                              "pattern": "^%permissions_id_pattern%$"
                            },
                            "invitation": {
                              "type": "object",
                              "required": [
                                "invitedBy"
                              ],
                              "properties": {
                                "user":{
                                  "type": "object",
                                  "required": [
                                    "displayName",
                                    "id"
                                  ],
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
                            "roles": {
                              "type": "array",
                              "items": [
                                {
                                  "type": "string",
                                  "pattern": "^%role_id_pattern%$"
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
    }
    """


  Scenario: user lists file shared with them in a group from sharer's personal space
    Given the administrator has assigned the role "Admin" to user "Alice" using the Graph API
    And user "Alice" has uploaded file with content "hello" to "textfile0.txt"
    And user "Alice" has created a group "grp1" using the Graph API
    And user "Brian" has been added to group "grp1"
    And user "Alice" has sent the following share invitation:
      | resource        | textfile0.txt |
      | space           | Personal      |
      | sharee          | grp1          |
      | shareType       | group         |
      | permissionsRole | Viewer        |
    When user "Brian" lists the shares shared with him using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
    """
    {
      "type": "object",
      "required": ["value"],
      "properties": {
        "value": {
          "type": "array",
          "items": {
            "type": "object",
            "required": [
              "@UI.Hidden",
              "@client.synchronize",
              "createdBy",
              "eTag",
              "file",
              "id",
              "lastModifiedDateTime",
              "name",
              "parentReference",
              "remoteItem",
              "size"
            ],
            "properties": {
              "@UI.Hidden":{
                "type": "boolean",
                "enum": [false]
              },
              "@client.synchronize":{
                "type": "boolean",
                "enum": [true]
              },
              "createdBy": {
                "type": "object",
                "required": ["user"],
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
                "pattern": "%etag_pattern%"
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
                  "driveType",
                  "id"
                ],
                "properties": {
                  "driveId": {
                    "type": "string",
                    "pattern": "^%space_id_pattern%$"
                  },
                  "driveType" : {
                    "type": "string",
                    "enum": ["virtual"]
                  },
                  "id" : {
                    "type": "string",
                    "pattern": "%space_id_pattern%"
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
                    "required": ["user"],
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
                            "enum": ["Alice Hansen"]
                          }
                        }
                      }
                    }
                  },
                  "eTag": {
                    "type": "string",
                    "pattern": "%etag_pattern%"
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
                    "pattern": "^%file_id_pattern%$"
                  },
                  "name": {
                    "type": "string",
                    "enum": ["textfile0.txt"]
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
                        "pattern": "%space_id_pattern%"
                      },
                      "driveType" : {
                        "type": "string",
                        "enum": ["personal"]
                      }
                    }
                  },
                  "permissions": {
                    "type": "array",
                    "maxItems": 1,
                    "minItems": 1,
                    "items": [
                      {
                        "type": "object",
                        "required": [
                          "grantedToV2",
                          "id",
                          "invitation",
                          "roles"
                        ],
                        "properties": {
                          "grantedToV2": {
                            "type": "object",
                            "required": ["group"],
                            "properties":{
                              "group": {
                                "type": "object",
                                "required": [
                                  "displayName",
                                  "id"
                                ],
                                "properties": {
                                  "displayName": {
                                    "type": "string",
                                    "enum": ["grp1"]
                                  },
                                  "id": {
                                    "type": "string",
                                    "pattern": "^%user_id_pattern%$"
                                  }
                                }
                              }
                            }
                          },
                          "id": {
                            "type": "string",
                            "pattern": "^%permissions_id_pattern%$"
                          },
                          "invitation": {
                            "type": "object",
                            "required": ["invitedBy"],
                            "properties": {
                              "user":{
                                "type": "object",
                                "required": [
                                  "displayName",
                                  "id"
                                ],
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
                          "roles": {
                            "type": "array",
                            "items": [
                              {
                                "type": "string",
                                "pattern": "^%role_id_pattern%$"
                              }
                            ]
                          }
                        }
                      }
                    ]
                  },
                  "size": {
                    "type": "number",
                    "enum": [5]
                  }
                }
              },
              "size": {
                "type": "number",
                "enum": [5]
              }
            }
          }
        }
      }
    }
    """


  Scenario: user lists file shared with them in a group from sharer's project space
    Given using spaces DAV path
    And the administrator has assigned the role "Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "new-space" with content "hello world" to "textfile0.txt"
    And user "Alice" has created a group "grp1" using the Graph API
    And user "Brian" has been added to group "grp1"
    And user "Alice" has sent the following share invitation:
      | resource        | textfile0.txt |
      | space           | new-space     |
      | sharee          | grp1          |
      | shareType       | group         |
      | permissionsRole | Viewer        |
    When user "Brian" lists the shares shared with him using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
    """
    {
      "type": "object",
      "required": ["value"],
      "properties": {
        "value": {
          "type": "array",
          "items": {
            "type": "object",
            "required": [
              "@UI.Hidden",
              "@client.synchronize",
              "eTag",
              "file",
              "id",
              "lastModifiedDateTime",
              "name",
              "parentReference",
              "remoteItem",
              "size"
            ],
            "properties": {
              "@UI.Hidden":{
                "type": "boolean",
                "enum": [false]
              },
              "@client.synchronize":{
                "type": "boolean",
                "enum": [true]
              },
              "eTag": {
                "type": "string",
                "pattern": "%etag_pattern%"
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
                "enum": ["textfile0.txt"]
              },
              "parentReference": {
                "type": "object",
                "required": [
                  "driveId",
                  "driveType",
                  "id"
                ],
                "properties": {
                  "driveId": {
                    "type": "string",
                    "pattern": "^%space_id_pattern%$"
                  },
                  "driveType" : {
                    "type": "string",
                    "enum": ["virtual"]
                  },
                  "id" : {
                    "type": "string",
                    "pattern": "%space_id_pattern%"
                  }
                }
              },
              "remoteItem": {
                "type": "object",
                "required": [
                  "eTag",
                  "file",
                  "id",
                  "lastModifiedDateTime",
                  "name",
                  "parentReference",
                  "permissions"
                ],
                "properties": {
                  "eTag": {
                    "type": "string",
                    "pattern": "%etag_pattern%"
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
                    "pattern": "^%file_id_pattern%$"
                  },
                  "name": {
                    "type": "string",
                    "enum": ["textfile0.txt"]
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
                        "pattern": "%space_id_pattern%"
                      },
                      "driveType" : {
                        "type": "string",
                        "enum": ["project"]
                      }
                    }
                  },
                  "permissions": {
                    "type": "array",
                    "maxItems": 1,
                    "minItems": 1,
                    "items": [
                      {
                        "type": "object",
                        "required": [
                          "grantedToV2",
                          "id",
                          "invitation",
                          "roles"
                        ],
                        "properties": {
                          "grantedToV2": {
                            "type": "object",
                            "required": ["group"],
                            "properties":{
                              "group": {
                                "type": "object",
                                "required": [
                                  "displayName",
                                  "id"
                                ],
                                "properties": {
                                  "displayName": {
                                    "type": "string",
                                    "enum": ["grp1"]
                                  },
                                  "id": {
                                    "type": "string",
                                    "pattern": "^%user_id_pattern%$"
                                  }
                                }
                              }
                            }
                          },
                          "id": {
                            "type": "string",
                            "pattern": "^%permissions_id_pattern%$"
                          },
                          "invitation": {
                            "type": "object",
                            "required": ["invitedBy"],
                            "properties": {
                              "user":{
                                "type": "object",
                                "required": [
                                  "displayName",
                                  "id"
                                ],
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
                          "roles": {
                            "type": "array",
                            "items": [
                              {
                                "type": "string",
                                "pattern": "^%role_id_pattern%$"
                              }
                            ]
                          }
                        }
                      }
                    ]
                  },
                  "size": {
                    "type": "number",
                    "enum": [11]
                  }
                }
              },
              "size": {
                "type": "number",
                "enum": [11]
              }
            }
          }
        }
      }
    }
    """


  Scenario: user lists folder shared with them in a group from sharer's personal space
    Given using spaces DAV path
    And the administrator has assigned the role "Admin" to user "Alice" using the Graph API
    And  user "Alice" has created folder "folder"
    And user "Alice" has created a group "grp1" using the Graph API
    And user "Brian" has been added to group "grp1"
    And user "Alice" has sent the following share invitation:
      | resource        | folder   |
      | space           | Personal |
      | sharee          | grp1     |
      | shareType       | group    |
      | permissionsRole | Viewer   |
    When user "Brian" lists the shares shared with him using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
    """
    {
      "type": "object",
      "required": ["value"],
      "properties": {
        "value": {
          "type": "array",
          "items": {
            "type": "object",
            "required": [
              "@UI.Hidden",
              "@client.synchronize",
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
              "@UI.Hidden":{
                "type": "boolean",
                "enum": [false]
              },
              "@client.synchronize":{
                "type": "boolean",
                "enum": [true]
              },
              "createdBy": {
                "type": "object",
                "required": ["user"],
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
                "pattern": "%etag_pattern%"
              },
              "id": {
                "type": "string",
                "pattern": "^%share_id_pattern%$"
              },
              "name": {
                "type": "string",
                "enum": ["folder"]
              },
              "parentReference": {
                "type": "object",
                "required": [
                  "driveId",
                  "driveType",
                  "id"
                ],
                "properties": {
                  "driveId": {
                    "type": "string",
                    "pattern": "^%space_id_pattern%$"
                  },
                  "driveType" : {
                    "type": "string",
                    "enum": ["virtual"]
                  },
                  "id" : {
                    "type": "string",
                    "pattern": "%space_id_pattern%"
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
                    "required": ["user"],
                    "properties": {
                      "user": {
                        "type": "object",
                        "required": ["id","displayName"],
                        "properties": {
                          "id": {
                            "type": "string",
                            "pattern": "^%user_id_pattern%$"
                          },
                          "displayName": {
                            "type": "string",
                            "enum": ["Alice Hansen"]
                          }
                        }
                      }
                    }
                  },
                  "eTag": {
                    "type": "string",
                    "pattern": "%etag_pattern%"
                  },
                  "id": {
                    "type": "string",
                    "pattern": "^%file_id_pattern%$"
                  },
                  "name": {
                    "type": "string",
                    "enum": ["folder"]
                  },
                  "parentReference": {
                    "type": "object",
                    "required": ["driveId","driveType"],
                    "properties": {
                      "driveId": {
                        "type": "string",
                        "pattern": "%space_id_pattern%"
                      },
                      "driveType" : {
                        "type": "string",
                        "enum": ["personal"]
                      }
                    }
                  },
                  "permissions": {
                    "type": "array",
                    "maxItems": 1,
                    "minItems": 1,
                    "items": [
                      {
                        "type": "object",
                        "required": [
                          "grantedToV2",
                          "id",
                          "invitation",
                          "roles"
                        ],
                        "properties": {
                          "grantedToV2": {
                            "type": "object",
                            "required": ["group"],
                            "properties":{
                              "group": {
                                "type": "object",
                                "required": [
                                  "displayName",
                                  "id"
                                ],
                                "properties": {
                                  "displayName": {
                                    "type": "string",
                                    "enum": ["grp1"]
                                  },
                                  "id": {
                                    "type": "string",
                                    "pattern": "^%user_id_pattern%$"
                                  }
                                }
                              }
                            }
                          },
                          "id": {
                            "type": "string",
                            "pattern": "^%permissions_id_pattern%$"
                          },
                          "invitation": {
                            "type": "object",
                            "required": ["invitedBy"],
                            "properties": {
                              "user":{
                                "type": "object",
                                "required": [
                                  "displayName",
                                  "id"
                                ],
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
                          "roles": {
                            "type": "array",
                            "items": [
                              {
                                "type": "string",
                                "pattern": "^%role_id_pattern%$"
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


  Scenario: user lists folder shared with them in a group from sharer's project space
    Given using spaces DAV path
    And the administrator has assigned the role "Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the Graph API
    And user "Alice" has created a group "grp1" using the Graph API
    And user "Alice" has created a folder "folder" in space "new-space"
    And user "Brian" has been added to group "grp1"
    And user "Alice" has sent the following share invitation:
      | resource        | folder    |
      | space           | new-space |
      | sharee          | grp1      |
      | shareType       | group     |
      | permissionsRole | Viewer    |
    When user "Brian" lists the shares shared with him using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
    """
    {
      "type": "object",
      "required": ["value"],
      "properties": {
        "value": {
          "type": "array",
          "items": {
            "type": "object",
            "required": [
              "@UI.Hidden",
              "@client.synchronize",
              "eTag",
              "folder",
              "id",
              "lastModifiedDateTime",
              "name",
              "parentReference",
              "remoteItem"
            ],
            "properties": {
              "@UI.Hidden":{
                "type": "boolean",
                "enum": [false]
              },
              "@client.synchronize":{
                "type": "boolean",
                "enum": [true]
              },
              "eTag": {
                "type": "string",
                "pattern": "%etag_pattern%"
              },
              "id": {
                "type": "string",
                "pattern": "^%share_id_pattern%$"
              },
              "name": {
                "type": "string",
                "enum": ["folder"]
              },
              "parentReference": {
                "type": "object",
                "required": [
                  "driveId",
                  "driveType",
                  "id"
                ],
                "properties": {
                  "driveId": {
                    "type": "string",
                    "pattern": "^%space_id_pattern%$"
                  },
                  "driveType" : {
                    "type": "string",
                    "enum": ["virtual"]
                  },
                  "id" : {
                    "type": "string",
                    "pattern": "%space_id_pattern%"
                  }
                }
              },
              "remoteItem": {
                "type": "object",
                "required": [
                  "eTag",
                  "folder",
                  "id",
                  "lastModifiedDateTime",
                  "name",
                  "parentReference",
                  "permissions"
                ],
                "properties": {
                  "eTag": {
                    "type": "string",
                    "pattern": "%etag_pattern%"
                  },
                  "id": {
                    "type": "string",
                    "pattern": "^%file_id_pattern%$"
                  },
                  "name": {
                    "type": "string",
                    "enum": ["folder"]
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
                        "pattern": "%space_id_pattern%"
                      },
                      "driveType" : {
                        "type": "string",
                        "enum": ["project"]
                      }
                    }
                  },
                  "permissions": {
                    "type": "array",
                    "maxItems": 1,
                    "minItems": 1,
                    "items": [
                      {
                        "type": "object",
                        "required": [
                          "grantedToV2",
                          "id",
                          "invitation",
                          "roles"
                        ],
                        "properties": {
                          "grantedToV2": {
                            "type": "object",
                            "required": ["group"],
                            "properties":{
                              "group": {
                                "type": "object",
                                "required": [
                                  "displayName",
                                  "id"
                                ],
                                "properties": {
                                  "displayName": {
                                    "type": "string",
                                    "enum": ["grp1"]
                                  },
                                  "id": {
                                    "type": "string",
                                    "pattern": "^%user_id_pattern%$"
                                  }
                                }
                              }
                            }
                          },
                          "id": {
                            "type": "string",
                            "pattern": "^%permissions_id_pattern%$"
                          },
                          "invitation": {
                            "type": "object",
                            "required": ["invitedBy"],
                            "properties": {
                              "user":{
                                "type": "object",
                                "required": [
                                  "displayName",
                                  "id"
                                ],
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
                          "roles": {
                            "type": "array",
                            "items": [
                              {
                                "type": "string",
                                "pattern": "^%role_id_pattern%$"
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


  Scenario: user lists the file with same name shared by two users with him/her
    Given user "Carol" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file with content "to share" to "textfile.txt"
    And user "Carol" has uploaded file with content "to share" to "textfile.txt"
    And user "Alice" has sent the following share invitation:
      | resource        | textfile.txt |
      | space           | Personal     |
      | sharee          | Brian        |
      | shareType       | user         |
      | permissionsRole | Viewer       |
    And user "Carol" has sent the following share invitation:
      | resource        | textfile.txt |
      | space           | Personal     |
      | sharee          | Brian        |
      | shareType       | user         |
      | permissionsRole | Viewer       |
    When user "Brian" lists the shares shared with him using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["value"],
        "properties": {
          "value": {
            "type": "array",
            "items": {
              "oneOf": [
                {
                  "type": "object",
                  "required": [
                    "@UI.Hidden",
                    "@client.synchronize",
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
                    "@UI.Hidden": {
                      "type": "boolean",
                      "enum": [false]
                    },
                    "@client.synchronize": {
                      "type": "boolean",
                      "enum": [true]
                    },
                    "createdBy": {
                      "type": "object",
                      "required": ["user"],
                      "properties": {
                        "user": {
                          "type": "object",
                          "required": ["displayName", "id"],
                          "properties": {
                            "displayName": {
                              "type": "string",
                              "enum": ["Carol King"]
                            }
                          }
                        }
                      }
                    },
                    "name": {
                      "type": "string",
                      "enum": ["textfile (1).txt"]
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
                        "permissions"
                      ],
                      "properties": {
                        "createdBy": {
                          "type": "object",
                          "required": ["user"],
                          "properties": {
                            "user": {
                              "type": "object",
                              "required": ["displayName","id"],
                              "properties": {
                                "displayName": {
                                  "type": "string",
                                  "enum": ["Carol King"]
                                }
                              }
                            }
                          }
                        },
                        "name": {
                          "type": "string",
                          "enum": ["textfile.txt"]
                        },
                        "permissions": {
                          "type": "array",
                          "items": [
                            {
                              "type": "object",
                              "required": ["grantedToV2", "id", "invitation", "roles"],
                              "properties": {
                                "grantedToV2": {
                                  "type": "object",
                                  "required": ["user"],
                                  "properties": {
                                    "user": {
                                      "type": "object",
                                      "required": ["displayName", "id"],
                                      "properties": {
                                        "displayName": {
                                          "type": "string",
                                          "enum": ["Brian Murphy"]
                                        }
                                      }
                                    }
                                  }
                                }
                              },
                              "invitation": {
                                "type": "object",
                                "required": ["invitedBy"],
                                "properties": {
                                  "invitedBy": {
                                    "type": "object",
                                    "required": ["user"],
                                    "properties": {
                                      "user": {
                                        "type": "object",
                                        "required": ["displayName", "id"],
                                        "properties": {
                                          "displayName": {
                                            "type": "string",
                                            "enum": ["Carol King"]
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
                },
                {
                  "type": "object",
                  "required": [
                    "@UI.Hidden",
                    "@client.synchronize",
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
                    "@UI.Hidden": {
                      "type": "boolean",
                      "enum": [false]
                    },
                    "@client.synchronize": {
                      "type": "boolean",
                      "enum": [true]
                    },
                    "createdBy": {
                      "type": "object",
                      "required": ["user"],
                      "properties": {
                        "user": {
                          "type": "object",
                          "required": ["displayName", "id"],
                          "properties": {
                            "displayName": {
                              "type": "string",
                              "enum": ["Alice Hansen"]
                            }
                          }
                        }
                      }
                    },
                    "name": {
                      "type": "string",
                      "enum": ["textfile.txt"]
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
                        "permissions"
                      ],
                      "properties": {
                        "createdBy": {
                          "type": "object",
                          "required": ["user"],
                          "properties": {
                            "user": {
                              "type": "object",
                              "required": ["displayName", "id"],
                              "properties": {
                                "displayName": {
                                  "type": "string",
                                  "enum": ["Alice Hansen"]
                                }
                              }
                            }
                          }
                        },
                        "name": {
                          "type": "string",
                          "enum": ["textfile.txt"]
                        },
                        "permissions": {
                          "type": "array",
                          "items": [
                            {
                              "type": "object",
                              "required": ["grantedToV2", "id", "invitation", "roles"],
                              "properties": {
                                "grantedToV2": {
                                  "type": "object",
                                  "required": ["user"],
                                  "properties": {
                                    "user": {
                                      "type": "object",
                                      "required": ["displayName", "id"],
                                      "properties": {
                                        "displayName": {
                                          "type": "string",
                                          "enum": ["Brian Murphy"]
                                        }
                                      }
                                    }
                                  }
                                }
                              },
                              "invitation": {
                                "type": "object",
                                "required": ["invitedBy"],
                                "properties": {
                                  "invitedBy": {
                                    "type": "object",
                                    "required": ["user"],
                                    "properties": {
                                      "user": {
                                        "type": "object",
                                        "required": ["displayName", "id"],
                                        "properties": {
                                          "displayName": {
                                            "type": "string",
                                            "enum": ["Alice Hansen"]
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
                }
              ]
            }
          }
        }
      }
      """


  Scenario: user lists the folder with same name shared by two users with him/her
    Given user "Carol" has been created with default attributes and without skeleton files
    And user "Alice" has created folder "folderToShare"
    And user "Carol" has created folder "folderToShare"
    And user "Alice" has sent the following share invitation:
      | resource        | folderToShare |
      | space           | Personal      |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | Viewer        |
    And user "Carol" has sent the following share invitation:
      | resource        | folderToShare |
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
        "required": ["value"],
        "properties": {
          "value": {
            "type": "array",
            "items": {
              "oneOf": [
                {
                  "type": "object",
                  "required": [
                    "@UI.Hidden",
                    "@client.synchronize",
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
                    "@UI.Hidden": {
                      "type": "boolean",
                      "enum": [false]
                    },
                    "@client.synchronize": {
                      "type": "boolean",
                      "enum": [true]
                    },
                    "createdBy": {
                      "type": "object",
                      "required": ["user"],
                      "properties": {
                        "user": {
                          "type": "object",
                          "required": ["displayName", "id"],
                          "properties": {
                            "displayName": {
                              "type": "string",
                              "enum": ["Carol King"]
                            }
                          }
                        }
                      }
                    },
                    "name": {
                      "type": "string",
                      "enum": ["folderToShare (1)"]
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
                          "required": ["user"],
                          "properties": {
                            "user": {
                              "type": "object",
                              "required": ["displayName","id"],
                              "properties": {
                                "displayName": {
                                  "type": "string",
                                  "enum": ["Carol King"]
                                }
                              }
                            }
                          }
                        },
                        "name": {
                          "type": "string",
                          "enum": ["folderToShare"]
                        },
                        "permissions": {
                          "type": "array",
                          "items": [
                            {
                              "type": "object",
                              "required": ["grantedToV2", "id", "invitation", "roles"],
                              "properties": {
                                "grantedToV2": {
                                  "type": "object",
                                  "required": ["user"],
                                  "properties": {
                                    "user": {
                                      "type": "object",
                                      "required": ["displayName", "id"],
                                      "properties": {
                                        "displayName": {
                                          "type": "string",
                                          "enum": ["Brian Murphy"]
                                        }
                                      }
                                    }
                                  }
                                }
                              },
                              "invitation": {
                                "type": "object",
                                "required": ["invitedBy"],
                                "properties": {
                                  "invitedBy": {
                                    "type": "object",
                                    "required": ["user"],
                                    "properties": {
                                      "user": {
                                        "type": "object",
                                        "required": ["displayName", "id"],
                                        "properties": {
                                          "displayName": {
                                            "type": "string",
                                            "enum": ["Carol King"]
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
                },
                {
                  "type": "object",
                  "required": [
                    "@UI.Hidden",
                    "@client.synchronize",
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
                    "@UI.Hidden": {
                      "type": "boolean",
                      "enum": [false]
                    },
                    "@client.synchronize": {
                      "type": "boolean",
                      "enum": [true]
                    },
                    "createdBy": {
                      "type": "object",
                      "required": ["user"],
                      "properties": {
                        "user": {
                          "type": "object",
                          "required": ["displayName", "id"],
                          "properties": {
                            "displayName": {
                              "type": "string",
                              "enum": ["Alice Hansen"]
                            }
                          }
                        }
                      }
                    },
                    "name": {
                      "type": "string",
                      "enum": ["folderToShare"]
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
                          "required": ["user"],
                          "properties": {
                            "user": {
                              "type": "object",
                              "required": ["displayName", "id"],
                              "properties": {
                                "displayName": {
                                  "type": "string",
                                  "enum": ["Alice Hansen"]
                                }
                              }
                            }
                          }
                        },
                        "name": {
                          "type": "string",
                          "enum": ["folderToShare"]
                        },
                        "permissions": {
                          "type": "array",
                          "items": [
                            {
                              "type": "object",
                              "required": ["grantedToV2", "id", "invitation", "roles"],
                              "properties": {
                                "grantedToV2": {
                                  "type": "object",
                                  "required": ["user"],
                                  "properties": {
                                    "user": {
                                      "type": "object",
                                      "required": ["displayName", "id"],
                                      "properties": {
                                        "displayName": {
                                          "type": "string",
                                          "enum": ["Brian Murphy"]
                                        }
                                      }
                                    }
                                  }
                                }
                              },
                              "invitation": {
                                "type": "object",
                                "required": ["invitedBy"],
                                "properties": {
                                  "invitedBy": {
                                    "type": "object",
                                    "required": ["user"],
                                    "properties": {
                                      "user": {
                                        "type": "object",
                                        "required": ["displayName", "id"],
                                        "properties": {
                                          "displayName": {
                                            "type": "string",
                                            "enum": ["Alice Hansen"]
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
                }
              ]
            }
          }
        }
      }
      """
