Feature: Send a sharing invitations
  As the owner of a resource
  I want to be able to send invitations to other users
  So that they can have access to it

  https://owncloud.dev/libre-graph-api/#/drives.permissions/Invite

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |


  Scenario Outline: send share invitation to user with different roles
    Given user "Alice" has uploaded file with content "to share" to "/textfile1.txt"
    And user "Alice" has created folder "FolderToShare"
    When user "Alice" sends the following share invitation using the Graph API:
      | resourceType | <resource-type> |
      | resource     | <path>          |
      | space        | Personal        |
      | sharee       | Brian           |
      | shareType    | user            |
      | role         | <role>          |
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
                "id",
                "roles",
                "grantedToV2"
              ],
              "properties": {
                "id": {
                  "type": "string",
                  "pattern": "^%share_id_pattern%$"
                },
                "roles": {
                  "type": "array",
                  "items": {
                    "type": "string",
                    "pattern": "^%role_id_pattern%$"
                  }
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
                            "Brian Murphy"
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
    Examples:
      | role        | resource-type | path           |
      | Viewer      | file          | /textfile1.txt |
      | File Editor | file          | /textfile1.txt |
      | Co Owner    | file          | /textfile1.txt |
      | Manager     | file          | /textfile1.txt |
      | Viewer      | folder        | FolderToShare  |
      | Editor      | folder        | FolderToShare  |
      | Co Owner    | folder        | FolderToShare  |
      | Uploader    | folder        | FolderToShare  |
      | Manager     | folder        | FolderToShare  |


  Scenario Outline: send share invitation to group with different roles
    Given user "Carol" has been created with default attributes and without skeleton files
    And group "grp1" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
      | Carol    | grp1      |
    And user "Alice" has uploaded file with content "to share" to "/textfile1.txt"
    And user "Alice" has created folder "FolderToShare"
    When user "Alice" sends the following share invitation using the Graph API:
      | resourceType | <resource-type> |
      | resource     | <path>          |
      | space        | Personal        |
      | sharee       | grp1            |
      | shareType    | group           |
      | role         | <role>          |
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
                "id",
                "roles",
                "grantedToV2"
              ],
              "properties": {
                "id": {
                  "type": "string",
                  "pattern": "^%share_id_pattern%$"
                },
                "roles": {
                  "type": "array",
                  "items": {
                    "type": "string",
                    "pattern": "^%role_id_pattern%$"
                  }
                },
                "grantedToV2": {
                  "type": "object",
                  "required": [
                    "group"
                  ],
                  "properties": {
                    "group": {
                      "type": "object",
                      "required": [
                        "id",
                        "displayName"
                      ],
                      "properties": {
                        "id": {
                          "type": "string",
                          "pattern": "^%group_id_pattern%$"
                        },
                        "displayName": {
                          "type": "string",
                          "enum": [
                            "grp1"
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
    Examples:
      | role        | resource-type | path           |
      | Viewer      | file          | /textfile1.txt |
      | File Editor | file          | /textfile1.txt |
      | Co Owner    | file          | /textfile1.txt |
      | Manager     | file          | /textfile1.txt |
      | Viewer      | folder        | FolderToShare  |
      | Editor      | folder        | FolderToShare  |
      | Co Owner    | folder        | FolderToShare  |
      | Uploader    | folder        | FolderToShare  |
      | Manager     | folder        | FolderToShare  |


  Scenario Outline: send share invitation for a file to user with different permissions
    Given user "Alice" has uploaded file with content "to share" to "textfile1.txt"
    When user "Alice" sends the following share invitation using the Graph API:
      | resourceType | file          |
      | resource     | textfile1.txt |
      | space        | Personal      |
      | sharee       | Brian         |
      | shareType    | user          |
      | permission   | <permission>  |
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
                "id",
                "@libre.graph.permissions.actions",
                "grantedToV2"
              ],
              "properties": {
                "id": {
                  "type": "string",
                  "pattern": "^%share_id_pattern%$"
                },
                "@libre.graph.permissions.actions": {
                  "type": "array",
                  "items": {
                    "type": "string",
                    "pattern": "^libre\\.graph\\/driveItem\\/<permission>$"
                  }
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
                            "Brian Murphy"
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
    Examples:
      | permission         |
      | permissions/create |
      | children/create    |
      | upload/create      |
      | path/read          |
      | quota/read         |
      | content/read       |
      | permissions/read   |
      | children/read      |
      | versions/read      |
      | deleted/read       |
      | basic/read         |
      | path/update        |
      | versions/update    |
      | deleted/update     |
      | permissions/update |
      | standard/delete    |
      | permissions/delete |
      | deleted/delete     |
      | permissions/deny   |


  Scenario Outline: send share invitation for a folder to user with different permissions
    Given user "Alice" has created folder "FolderToShare"
    When user "Alice" sends the following share invitation using the Graph API:
      | resourceType | folder        |
      | resource     | FolderToShare |
      | space        | Personal      |
      | sharee       | Brian         |
      | shareType    | user          |
      | permission   | <permission>  |
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
                "id",
                "@libre.graph.permissions.actions",
                "grantedToV2"
              ],
              "properties": {
                "id": {
                  "type": "string",
                  "pattern": "^%share_id_pattern%$"
                },
                "@libre.graph.permissions.actions": {
                  "type": "array",
                  "items": {
                    "type": "string",
                    "pattern": "^libre\\.graph\\/driveItem\\/<permission>$"
                  }
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
                            "Brian Murphy"
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
    Examples:
      | permission         |
      | permissions/create |
      | children/create    |
      | upload/create      |
      | path/read          |
      | quota/read         |
      | content/read       |
      | permissions/read   |
      | children/read      |
      | versions/read      |
      | deleted/read       |
      | basic/read         |
      | path/update        |
      | versions/update    |
      | deleted/update     |
      | permissions/update |
      | standard/delete    |
      | permissions/delete |
      | deleted/delete     |
      | permissions/deny   |


  Scenario Outline: send share invitation for a file to group with different permissions
    Given user "Carol" has been created with default attributes and without skeleton files
    And group "grp1" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
      | Carol    | grp1      |
    And user "Alice" has uploaded file with content "to share" to "textfile1.txt"
    When user "Alice" sends the following share invitation using the Graph API:
      | resourceType | file          |
      | resource     | textfile1.txt |
      | space        | Personal      |
      | sharee       | grp1          |
      | shareType    | group         |
      | permission   | <permission>  |
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
                "id",
                "@libre.graph.permissions.actions",
                "grantedToV2"
              ],
              "properties": {
                "id": {
                  "type": "string",
                  "pattern": "^%share_id_pattern%$"
                },
                "@libre.graph.permissions.actions": {
                  "type": "array",
                  "items": {
                    "type": "string",
                    "pattern": "^libre\\.graph\\/driveItem\\/<permission>$"
                  }
                },
                "grantedToV2": {
                  "type": "object",
                  "required": [
                    "group"
                  ],
                  "properties": {
                    "group": {
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
                            "grp1"
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
    Examples:
      | permission         |
      | permissions/create |
      | children/create    |
      | upload/create      |
      | path/read          |
      | quota/read         |
      | content/read       |
      | permissions/read   |
      | children/read      |
      | versions/read      |
      | deleted/read       |
      | basic/read         |
      | path/update        |
      | versions/update    |
      | deleted/update     |
      | permissions/update |
      | standard/delete    |
      | permissions/delete |
      | deleted/delete     |
      | permissions/deny   |


  Scenario Outline: send share invitation for a folder to group with different permissions
    Given user "Carol" has been created with default attributes and without skeleton files
    And group "grp1" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
      | Carol    | grp1      |
    And user "Alice" has created folder "FolderToShare"
    When user "Alice" sends the following share invitation using the Graph API:
      | resourceType | folder        |
      | resource     | FolderToShare |
      | space        | Personal      |
      | sharee       | grp1          |
      | shareType    | group         |
      | permission   | <permission>  |
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
                "id",
                "@libre.graph.permissions.actions",
                "grantedToV2"
              ],
              "properties": {
                "id": {
                  "type": "string",
                  "pattern": "^%share_id_pattern%$"
                },
                "@libre.graph.permissions.actions": {
                  "type": "array",
                  "items": {
                    "type": "string",
                    "pattern": "^libre\\.graph\\/driveItem\\/<permission>$"
                  }
                },
                "grantedToV2": {
                  "type": "object",
                  "required": [
                    "group"
                  ],
                  "properties": {
                    "group": {
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
                            "grp1"
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
    Examples:
      | permission         |
      | permissions/create |
      | children/create    |
      | upload/create      |
      | path/read          |
      | quota/read         |
      | content/read       |
      | permissions/read   |
      | children/read      |
      | versions/read      |
      | deleted/read       |
      | basic/read         |
      | path/update        |
      | versions/update    |
      | deleted/update     |
      | permissions/update |
      | standard/delete    |
      | permissions/delete |
      | deleted/delete     |
      | permissions/deny   |


  Scenario Outline: send share invitation for file to group with different roles
    Given user "Carol" has been created with default attributes and without skeleton files
    And group "grp1" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
      | Carol    | grp1      |
    And user "Alice" has uploaded file with content "to share" to "/textfile1.txt"
    And user "Alice" has sent the following share invitation using the Graph API:
      | resourceType | file           |
      | resource     | /textfile1.txt |
      | space        | Personal       |
      | sharee       | grp1           |
      | shareType    | group          |
      | role         | <role>         |
    When user "Brian" lists the resources shared with him using the Graph API
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
                "id",
                "name",
                "parentReference",
                "remoteItem"
              ],
              "properties": {
                "id": {
                  "type": "string",
                  "pattern": "^%share_id_pattern%$"
                },
                "name": {
                  "type": "string",
                  "enum": [
                    "textfile1.txt"
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
                      "pattern": "^%share_id_pattern%$"
                    },
                    "driveType": {
                      "type": "string",
                      "enum": ["personal"]
                    }
                  }
                },
                "remoteItem": {
                  "type": "object",
                  "required": [
                    "createdBy",
                    "file",
                    "id",
                    "name",
                    "shared",
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
                      "pattern": "^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}\\$[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}![a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$"
                    },
                    "name": {
                      "type": "string",
                      "enum": [
                        "textfile1.txt"
                      ]
                    },
                    "shared": {
                      "type": "object",
                      "required": [
                        "sharedBy",
                        "owner"
                      ],
                      "properties": {
                        "owner": {
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
                        "sharedBy": {
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
                        }
                      }
                    },
                    "size": {
                      "type": "number",
                      "enum": [
                        8
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
    When user "Carol" lists the resources shared with her using the Graph API
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
                "id",
                "name",
                "parentReference",
                "remoteItem"
              ],
              "properties": {
                "id": {
                  "type": "string",
                  "pattern": "^%share_id_pattern%$"
                },
                "name": {
                  "type": "string",
                  "enum": [
                    "textfile1.txt"
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
                      "pattern": "^%share_id_pattern%$"
                    },
                    "driveType": {
                      "type": "string",
                      "enum": ["personal"]
                    }
                  }
                },
                "remoteItem": {
                  "type": "object",
                  "required": [
                    "createdBy",
                    "file",
                    "id",
                    "name",
                    "shared",
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
                      "pattern": "^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}\\$[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}![a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$"
                    },
                    "name": {
                      "type": "string",
                      "enum": [
                        "textfile1.txt"
                      ]
                    },
                    "shared": {
                      "type": "object",
                      "required": [
                        "sharedBy",
                        "owner"
                      ],
                      "properties": {
                        "owner": {
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
                        "sharedBy": {
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
                        }
                      }
                    },
                    "size": {
                      "type": "number",
                      "enum": [
                        8
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
    Examples:
      | role        |
      | Viewer      |
      | File Editor |
      | Co Owner    |
      | Manager     |


  Scenario Outline: send share invitation for folder to group with different roles
    Given user "Carol" has been created with default attributes and without skeleton files
    And group "grp1" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
      | Carol    | grp1      |
    And user "Alice" has created folder "FolderToShare"
    And user "Alice" has sent the following share invitation using the Graph API:
      | resourceType | folder        |
      | resource     | FolderToShare |
      | space        | Personal      |
      | sharee       | grp1          |
      | shareType    | group         |
      | role         | <role>        |
    When user "Brian" lists the resources shared with him using the Graph API
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
                "name",
                "id",
                "remoteItem"
              ],
              "properties": {
                "id": {
                  "type": "string",
                  "pattern": "^%share_id_pattern%$"
                },
                "name": {
                  "type": "string",
                  "enum": [
                    "FolderToShare"
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
                      "pattern": "^%share_id_pattern%$"
                    },
                    "driveType": {
                      "type": "string",
                      "enum": ["personal"]
                    }
                  }
                },
                "remoteItem": {
                  "type": "object",
                  "required": [
                    "createdBy",
                    "folder",
                    "id",
                    "name",
                    "shared",
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
                    "id": {
                      "type": "string",
                      "pattern": "^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}\\$[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}![a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$"
                    },
                    "name": {
                      "type": "string",
                      "enum": [
                        "FolderToShare"
                      ]
                    },
                    "shared": {
                      "type": "object",
                      "required": [
                        "sharedBy",
                        "owner"
                      ],
                      "properties": {
                        "owner": {
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
                        "sharedBy": {
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
                        }
                      }
                    },
                    "size": {
                      "type": "number",
                      "enum": [
                        0
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
    When user "Carol" lists the resources shared with her using the Graph API
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
                "name",
                "id",
                "remoteItem"
              ],
              "properties": {
                "id": {
                  "type": "string",
                  "pattern": "^%share_id_pattern%$"
                },
                "name": {
                  "type": "string",
                  "enum": [
                    "FolderToShare"
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
                      "pattern": "^%share_id_pattern%$"
                    },
                    "driveType": {
                      "type": "string",
                      "enum": ["personal"]
                    }
                  }
                },
                "remoteItem": {
                  "type": "object",
                  "required": [
                    "createdBy",
                    "folder",
                    "id",
                    "name",
                    "shared",
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
                    "id": {
                      "type": "string",
                      "pattern": "^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}\\$[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}![a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$"
                    },
                    "name": {
                      "type": "string",
                      "enum": [
                        "FolderToShare"
                      ]
                    },
                    "shared": {
                      "type": "object",
                      "required": [
                        "sharedBy",
                        "owner"
                      ],
                      "properties": {
                        "owner": {
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
                        "sharedBy": {
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
                        }
                      }
                    },
                    "size": {
                      "type": "number",
                      "enum": [
                        0
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
    Examples:
      | role        |
      | Viewer      |
      | Editor      |
      | Co Owner    |
      | Uploader    |
      | Manager     |
