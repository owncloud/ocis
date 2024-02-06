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


  Scenario: user lists the file shared with them when auto-sync is disabled
    Given user "Alice" has uploaded file with content "to share" to "textfile.txt"
    And user "Brian" has disabled the auto-sync share
    And user "Alice" has sent the following share invitation:
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
                  "enum": [false]
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
                  "enum": ["textfile.txt"]
                },
                "parentReference": {
                  "type": "object",
                  "required": ["driveId", "driveType", "id"],
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
                      "required": ["user"],
                      "properties": {
                        "user": {
                          "type": "object",
                          "required": ["displayName","id"],
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
                      "pattern": "^%file_id_pattern%$"
                    },
                    "name": {
                      "type": "string",
                      "enum": ["textfile.txt"]
                    },
                    "parentReference": {
                      "type": "object",
                      "required": ["driveId", "driveType"],
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
                                    },
                                    "id": {
                                      "type": "string",
                                      "pattern": "^%user_id_pattern%$"
                                    }
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
                                      },
                                      "id": {
                                        "type": "string",
                                        "pattern": "^%user_id_pattern%$"
                                      }
                                    }
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
                      ]
                    },
                    "size": {
                      "type": "number",
                      "enum": [8]
                    }
                  }
                },
                "size": {
                  "type": "number",
                  "enum": [8]
                }
              }
            }
          }
        }
      }
      """


  Scenario: user lists the folder shared with them when auto-sync is disabled
    Given user "Alice" has created folder "folderToShare"
    And user "Brian" has disabled the auto-sync share
    And user "Alice" has sent the following share invitation:
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
                  "enum": [false]
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
                  "enum": ["folderToShare"]
                },
                "parentReference": {
                  "type": "object",
                  "required": ["driveId", "driveType", "id"],
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
                      "required": ["user"],
                      "properties": {
                        "user": {
                          "type": "object",
                          "required": ["displayName","id"],
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
                      "pattern": "^%file_id_pattern%$"
                    },
                    "name": {
                      "type": "string",
                      "enum": ["folderToShare"]
                    },
                    "parentReference": {
                      "type": "object",
                      "required": ["driveId", "driveType"],
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
                                    },
                                    "id": {
                                      "type": "string",
                                      "pattern": "^%user_id_pattern%$"
                                    }
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
                                      },
                                      "id": {
                                        "type": "string",
                                        "pattern": "^%user_id_pattern%$"
                                      }
                                    }
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


  Scenario: group member lists the file shared with them when auto-sync is disabled
    And user "Carol" has been created with default attributes and without skeleton files
    And group "grp1" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
      | Carol    | grp1      |
    Given user "Alice" has uploaded file with content "to share" to "textfile.txt"
    And user "Brian" has disabled the auto-sync share
    And user "Alice" has sent the following share invitation:
      | resource        | textfile.txt |
      | space           | Personal     |
      | sharee          | grp1         |
      | shareType       | group        |
      | permissionsRole | Viewer       |
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
                  "enum": [false]
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
                  "enum": ["textfile.txt"]
                },
                "parentReference": {
                  "type": "object",
                  "required": ["driveId", "driveType", "id"],
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
                      "required": ["user"],
                      "properties": {
                        "user": {
                          "type": "object",
                          "required": ["displayName","id"],
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
                      "pattern": "^%file_id_pattern%$"
                    },
                    "name": {
                      "type": "string",
                      "enum": ["textfile.txt"]
                    },
                    "parentReference": {
                      "type": "object",
                      "required": ["driveId", "driveType"],
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
                            "grantedToV2": {
                              "type": "object",
                              "required": ["group"],
                              "properties": {
                                "group": {
                                  "type": "object",
                                  "required": ["displayName", "id"],
                                  "properties": {
                                    "displayName": {
                                      "type": "string",
                                      "enum": ["grp1"]
                                    },
                                    "id": {
                                      "type": "string",
                                      "pattern": "^%group_id_pattern%$"
                                    }
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
                                      },
                                      "id": {
                                        "type": "string",
                                        "pattern": "^%user_id_pattern%$"
                                      }
                                    }
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
                      ]
                    },
                    "size": {
                      "type": "number",
                      "enum": [8]
                    }
                  }
                },
                "size": {
                  "type": "number",
                  "enum": [8]
                }
              }
            }
          }
        }
      }
      """


  Scenario: group member lists the folder shared with them when auto-sync is disabled
    And user "Carol" has been created with default attributes and without skeleton files
    And group "grp1" has been created
    And the following users have been added to the following groups
      | username | groupname |
      | Brian    | grp1      |
      | Carol    | grp1      |
    Given user "Alice" has created folder "folderToShare"
    And user "Brian" has disabled the auto-sync share
    And user "Alice" has sent the following share invitation:
      | resource        | folderToShare |
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
                  "enum": [false]
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
                  "enum": ["folderToShare"]
                },
                "parentReference": {
                  "type": "object",
                  "required": ["driveId", "driveType", "id"],
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
                      "required": ["user"],
                      "properties": {
                        "user": {
                          "type": "object",
                          "required": ["displayName","id"],
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
                      "pattern": "^%file_id_pattern%$"
                    },
                    "name": {
                      "type": "string",
                      "enum": ["folderToShare"]
                    },
                    "parentReference": {
                      "type": "object",
                      "required": ["driveId", "driveType"],
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
                            "grantedToV2": {
                              "type": "object",
                              "required": ["group"],
                              "properties": {
                                "user": {
                                  "type": "object",
                                  "required": ["displayName", "id"],
                                  "properties": {
                                    "displayName": {
                                      "type": "string",
                                      "enum": ["grp1"]
                                    },
                                    "id": {
                                      "type": "string",
                                      "pattern": "^%group_id_pattern%$"
                                    }
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
                                      },
                                      "id": {
                                        "type": "string",
                                        "pattern": "^%user_id_pattern%$"
                                      }
                                    }
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
