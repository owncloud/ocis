@api @skipOnOcV10
Feature: get users
  As an admin
  I want to be able to retrieve user information
  So that I can see the information

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
    And the administrator has given "Alice" the role "Admin" using the settings api

  @skipOnStable2.0
  Scenario: admin user gets the information of a user
    When user "Alice" gets information of user "Brian" using Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
    """
    {
      "type": "object",
      "required": [
        "displayName",
        "id",
        "mail",
        "onPremisesSamAccountName",
        "accountEnabled"
      ],
      "properties": {
        "displayName": {
          "type": "string",
          "enum": ["Brian Murphy"]
        },
        "id" : {
          "type": "string",
          "pattern": "^%user_id_pattern%$"
        },
        "mail": {
          "type": "string",
          "enum": ["brian@example.org"]
        },
        "onPremisesSamAccountName": {
          "type": "string",
          "enum": ["Brian"]
        },
        "accountEnabled": {
          "type": "boolean",
          "enum": [true]
        }
      }
    }
    """


  Scenario: non-admin user tries to get the information of a user
    When user "Brian" tries to get information of user "Alice" using Graph API
    Then the HTTP status code should be "401"
    And the JSON data of the response should match
    """
    {
      "type": "object",
      "required": [
        "error"
      ],
      "properties": {
        "error": {
          "type": "object",
          "required": [
            "message"
          ],
          "properties": {
            "type": "string",
            "enum": ["Unauthorized"]
          }
        }
      }
    }
    """

  @skipOnStable2.0
  Scenario: admin user gets all users
    When user "Alice" gets all users using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should contain the user "Alice Hansen" in the item 'value', the user-details should match
    """
    {
      "type": "object",
      "required": [
        "id",
        "mail",
        "onPremisesSamAccountName",
        "accountEnabled"
      ],
      "properties": {
        "id" : {
          "type": "string",
          "pattern": "^%user_id_pattern%$"
        },
        "mail": {
          "type": "string",
          "enum": ["alice@example.org"]
        },
        "onPremisesSamAccountName": {
          "type": "string",
          "enum": ["Alice"]
        },
        "accountEnabled": {
          "type": "boolean",
          "enum": [true]
        }
      }
    }
    """
    And the JSON data of the response should contain the user "Brian Murphy" in the item 'value', the user-details should match
    """
    {
      "type": "object",
      "required": [
        "id",
        "mail",
        "onPremisesSamAccountName",
        "accountEnabled"
      ],
      "properties": {
        "id" : {
          "type": "string",
          "pattern": "^%user_id_pattern%$"
        },
        "mail": {
          "type": "string",
          "enum": ["brian@example.org"]
        },
        "onPremisesSamAccountName": {
          "type": "string",
          "enum": ["Brian"]
        },
        "accountEnabled": {
          "type": "boolean",
          "enum": [true]
        }
      }
    }
    """

  @skipOnStable2.0
  Scenario: admin user gets all users include disabled users
    Given the user "Alice" has disabled user "Brian" using the Graph API
    When user "Alice" gets all users using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should contain the user "Alice Hansen" in the item 'value', the user-details should match
    """
    {
      "type": "object",
      "required": [
        "id",
        "mail",
        "onPremisesSamAccountName",
        "accountEnabled"
      ],
      "properties": {
        "id" : {
          "type": "string",
          "pattern": "^%user_id_pattern%$"
        },
        "mail": {
          "type": "string",
          "enum": ["alice@example.org"]
        },
        "onPremisesSamAccountName": {
          "type": "string",
          "enum": ["Alice"]
        },
        "accountEnabled": {
          "type": "boolean",
          "enum": [true]
        }
      }
    }
    """
    And the JSON data of the response should contain the user "Brian Murphy" in the item 'value', the user-details should match
    """
    {
      "type": "object",
      "required": [
        "id",
        "mail",
        "onPremisesSamAccountName",
        "accountEnabled"
      ],
      "properties": {
        "id" : {
          "type": "string",
          "pattern": "^%user_id_pattern%$"
        },
        "mail": {
          "type": "string",
          "enum": ["brian@example.org"]
        },
        "onPremisesSamAccountName": {
          "type": "string",
          "enum": ["Brian"]
        },
        "accountEnabled": {
          "type": "boolean",
          "enum": [false]
        }
      }
    }
    """


  Scenario: non-admin user tries to get all users
    When user "Brian" tries to get all users using the Graph API
    Then the HTTP status code should be "401"
    And the JSON data of the response should match
    """
    {
      "type": "object",
      "required": [
        "error"
      ],
      "properties": {
        "error": {
          "type": "object",
          "required": [
            "message"
          ],
          "properties": {
            "type": "string",
            "enum": ["Unauthorized"]
          }
        }
      }
    }
    """

  @skipOnStable2.0
  Scenario: admin user gets the drive information of a user
    When the user "Alice" gets user "Brian" along with his drive information using Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
    """
      {
        "type": "object",
        "required": [
          "displayName",
          "id",
          "mail",
          "onPremisesSamAccountName",
          "drive",
          "accountEnabled"
        ],
        "properties": {
          "displayName": {
            "type": "string",
            "enum": ["Brian Murphy"]
          },
          "id" : {
            "type": "string",
            "pattern": "^%user_id_pattern%$"
          },
          "mail": {
            "type": "string",
            "enum": ["brian@example.org"]
          },
          "onPremisesSamAccountName": {
            "type": "string",
            "enum": ["Brian"]
          },
          "accountEnabled": {
            "type": "boolean",
            "enum": [true]
          },
          "drive": {
            "type": "object",
            "required": [
              "driveAlias",
              "id",
              "name",
              "owner",
              "quota",
              "root",
              "webUrl"
            ],
            "properties": {
              "driveType" : {
                "type": "string",
                "enum": ["personal"]
              },
              "driveAlias" : {
                "type": "string",
                "enum": ["personal/brian"]
              },
              "id" : {
                "type": "string",
                "pattern": "^%space_id_pattern%$"
              },
              "name": {
                "type": "string",
                "enum": ["Brian Murphy"]
              },
              "owner": {
                "type": "object",
                "required": [
                  "user"
                ],
                "properties": {
                  "user": "string",
                  "required": [
                    "id"
                  ],
                  "properties": {
                    "id": {
                      "type": "string",
                      "enum": ["%user_id_pattern%"]
                    }
                  }
                }
              },
              "quota": {
                "type": "object",
                "required": [
                  "state"
                ],
                "properties": {
                  "state": {
                    "type": "string",
                    "enum": ["normal"]
                  }
                }
              },
              "root": {
                "type": "object",
                "required": [
                  "id",
                  "webDavUrl"
                ],
                "properties": {
                  "state": {
                    "type": "string",
                    "enum": ["normal"]
                  },
                  "webDavUrl": {
                    "type": "string",
                    "pattern": "^%base_url%/dav/spaces/%space_id_pattern%$"
                  }
                }
              },
              "webUrl": {
                "type": "string",
                "pattern": "^%base_url%/f/%space_id_pattern%$"
              }
            }
          }
        }
      }
    """

  @skipOnStable2.0
  Scenario: normal user gets his/her own drive information
    When the user "Brian" gets his drive information using Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
    """
      {
        "type": "object",
        "required": [
          "displayName",
          "id",
          "mail",
          "onPremisesSamAccountName",
          "drive",
          "accountEnabled"
        ],
        "properties": {
          "displayName": {
            "type": "string",
            "enum": ["Brian Murphy"]
          },
          "id" : {
            "type": "string",
            "pattern": "^%user_id_pattern%$"
          },
          "mail": {
            "type": "string",
            "enum": ["brian@example.org"]
          },
          "onPremisesSamAccountName": {
            "type": "string",
            "enum": ["Brian"]
          },
          "accountEnabled": {
            "type": "boolean",
            "enum": [true]
          },
          "drive": {
            "type": "object",
            "required": [
              "driveAlias",
              "id",
              "name",
              "owner",
              "quota",
              "root",
              "webUrl"
            ],
            "properties": {
              "driveType" : {
                "type": "string",
                "enum": ["personal"]
              },
              "driveAlias" : {
                "type": "string",
                "enum": ["personal/brian"]
              },
              "id" : {
                "type": "string",
                "pattern": "^%space_id_pattern%$"
              },
              "name": {
                "type": "string",
                "enum": ["Brian Murphy"]
              },
              "owner": {
                "type": "object",
                "required": [
                  "user"
                ],
                "properties": {
                  "user": "string",
                  "required": [
                    "id"
                  ],
                  "properties": {
                    "id": {
                      "type": "string",
                      "enum": ["%user_id_pattern%"]
                    }
                  }
                }
              },
              "quota": {
                "type": "object",
                "required": [
                  "state"
                ],
                "properties": {
                  "state": {
                    "type": "string",
                    "enum": ["normal"]
                  }
                }
              },
              "root": {
                "type": "object",
                "required": [
                  "id",
                  "webDavUrl"
                ],
                "properties": {
                  "state": {
                    "type": "string",
                    "enum": ["normal"]
                  },
                  "webDavUrl": {
                    "type": "string",
                    "pattern": "^%base_url%/dav/spaces/%space_id_pattern%$"
                  }
                }
              },
              "webUrl": {
                "type": "string",
                "pattern": "^%base_url%/f/%space_id_pattern%$"
              }
            }
          }
        }
      }
    """

  @skipOnStable2.0
  Scenario: admin user gets the group information of a user
    Given group "tea-lover" has been created
    And group "coffee-lover" has been created
    And user "Brian" has been added to group "tea-lover"
    And user "Brian" has been added to group "coffee-lover"
    When the user "Alice" gets user "Brian" along with his group information using Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
    """
    {
      "type": "object",
      "required": [
        "id",
        "mail",
        "onPremisesSamAccountName"
      ],
      "properties": {
        "id" : {
          "type": "string",
          "pattern": "^%user_id_pattern%$"
        },
        "mail": {
          "type": "string",
          "enum": ["brian@example.org"]
        },
        "onPremisesSamAccountName": {
          "type": "string",
          "enum": ["Brian"]
        },
        "memberOf": {
          "type": "array",
          "items": [
            {
              "type": "object",
              "required": [
                "displayName"
              ],
              "properties": {
                "displayName": {
                  "type": "string",
                  "enum": ["tea-lover"]
                }
              }
            },
            {
              "type": "object",
              "required": [
                "displayName"
              ],
              "properties": {
                "displayName": {
                  "type": "string",
                  "enum": ["coffee-lover"]
                }
              }
            }
          ]
        }
      }
    }
    """


  Scenario: non-admin user tries to get the group information of a user
    Given user "Carol" has been created with default attributes and without skeleton files
    And group "coffee-lover" has been created
    And user "Brian" has been added to group "coffee-lover"
    When the user "Carol" gets user "Brian" along with his group information using Graph API
    Then the HTTP status code should be "401"
    And the JSON data of the response should match
    """
    {
      "type": "object",
      "required": [
        "error"
      ],
      "properties": {
        "error": {
          "type": "object",
          "required": [
            "message"
          ],
          "properties": {
            "type": "string",
            "enum": ["Unauthorized"]
          }
        }
      }
    }
    """

  @skipOnStable2.0
  Scenario: admin user gets all users of certain groups
    Given user "Carol" has been created with default attributes and without skeleton files
    And the user "Alice" has disabled user "Carol" using the Graph API
    And group "tea-lover" has been created
    And group "coffee-lover" has been created
    And user "Alice" has been added to group "tea-lover"
    And user "Brian" has been added to group "tea-lover"
    And user "Brian" has been added to group "coffee-lover"
    When the user "Alice" gets all users of the group "tea-lover" using the Graph API
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
          "items": [
            {
              "type": "object",
              "required": [
                "id",
                "mail",
                "onPremisesSamAccountName",
                "accountEnabled"
              ],
              "properties": {
                "id" : {
                  "type": "string",
                  "pattern": "^%user_id_pattern%$"
                },
                "mail": {
                  "type": "string",
                  "enum": ["alice@example.org"]
                },
                "onPremisesSamAccountName": {
                  "type": "string",
                  "enum": ["Alice"]
                },
                "accountEnabled": {
                  "type": "boolean",
                  "enum": [true]
                }
              }
            },
            {
              "type": "object",
              "required": [
                "id",
                "mail",
                "onPremisesSamAccountName",
                "accountEnabled"
              ],
              "properties": {
                "id" : {
                  "type": "string",
                  "pattern": "^%user_id_pattern%$"
                },
                "mail": {
                  "type": "string",
                  "enum": ["brian@example.org"]
                },
                "onPremisesSamAccountName": {
                  "type": "string",
                  "enum": ["Brian"]
                },
                "accountEnabled": {
                  "type": "boolean",
                  "enum": [true]
                }
              }
            }
          ],
          "additionalItems": false
        }
      }
    }
    """
    And the JSON data of the response should not contain the user "Carol King" in the item 'value'
    When the user "Alice" gets all users of two groups "tea-lover,coffee-lover" using the Graph API
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
          "items": [
            {
              "type": "object",
              "required": [
                "id",
                "mail",
                "onPremisesSamAccountName",
                "accountEnabled"
              ],
              "properties": {
                "id" : {
                  "type": "string",
                  "pattern": "^%user_id_pattern%$"
                },
                "mail": {
                  "type": "string",
                  "enum": ["brian@example.org"]
                },
                "onPremisesSamAccountName": {
                  "type": "string",
                  "enum": ["Brian"]
                },
                "accountEnabled": {
                  "type": "boolean",
                  "enum": [true]
                }
              }
            }
          ],
          "additionalItems": false
        }
      }
    }
    """
    And the JSON data of the response should not contain the user "Carol King" in the item 'value'
    And the JSON data of the response should not contain the user "Alice Hansen" in the item 'value'



  Scenario Outline: non admin user tries to get users of certain groups
    Given the administrator has given "Brian" the role "<role>" using the settings api
    And group "tea-lover" has been created
    And user "Alice" has been added to group "tea-lover"
    When the user "Brian" gets all users of the group "tea-lover" using the Graph API
    Then the HTTP status code should be "401"
    And the JSON data of the response should match
    """
    {
      "type": "object",
      "required": [
        "error"
      ],
      "properties": {
        "error": {
          "type": "object",
          "required": [
            "message"
          ],
          "properties": {
            "type": "string",
            "enum": ["Unauthorized"]
          }
        }
      }
    }
    """
    Examples:
      | role        |
      | Space Admin |
      | User        |
      | Guest       |

  @skipOnStable2.0
  Scenario: admin user gets all users with certain roles and members of a certain group
    Given user "Carol" has been created with default attributes and without skeleton files
    And the administrator has given "Brian" the role "Space Admin" using the settings api
    And the administrator has given "Carol" the role "Space Admin" using the settings api
    And group "tea-lover" has been created
    And user "Brian" has been added to group "tea-lover"
    When the user "Alice" gets all users with role "Space Admin" using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should contain the user "Brian Murphy" in the item 'value', the user-details should match
    """
    {
      "type": "object",
      "required": [
        "id",
        "mail",
        "onPremisesSamAccountName",
        "accountEnabled"
      ],
      "properties": {
        "id" : {
          "type": "string",
          "pattern": "^%user_id_pattern%$"
        },
        "mail": {
          "type": "string",
          "enum": ["brian@example.org"]
        },
        "onPremisesSamAccountName": {
          "type": "string",
          "enum": ["Brian"]
        },
        "accountEnabled": {
          "type": "boolean",
          "enum": [true]
        }
      }
    }
    """
    And the JSON data of the response should contain the user "Carol King" in the item 'value', the user-details should match
    """
    {
      "type": "object",
      "required": [
        "id",
        "mail",
        "onPremisesSamAccountName",
        "accountEnabled"
      ],
      "properties": {
        "id" : {
          "type": "string",
          "pattern": "^%user_id_pattern%$"
        },
        "mail": {
          "type": "string",
          "enum": ["carol@example.org"]
        },
        "onPremisesSamAccountName": {
          "type": "string",
          "enum": ["Carol"]
        },
        "accountEnabled": {
          "type": "boolean",
          "enum": [true]
        }
      }
    }
    """
    But the JSON data of the response should not contain the user "Alice Hansen" in the item 'value'
    When the user "Alice" gets all users with role "Space Admin" and member of the group "tea-lover" using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should contain the user "Brian Murphy" in the item 'value', the user-details should match
    """
    {
      "type": "object",
      "required": [
        "id",
        "mail",
        "onPremisesSamAccountName",
        "accountEnabled"
      ],
      "properties": {
        "id" : {
          "type": "string",
          "pattern": "^%user_id_pattern%$"
        },
        "mail": {
          "type": "string",
          "enum": ["brian@example.org"]
        },
        "onPremisesSamAccountName": {
          "type": "string",
          "enum": ["Brian"]
        },
        "accountEnabled": {
          "type": "boolean",
          "enum": [true]
        }
      }
    }
    """
    But the JSON data of the response should not contain the user "Carol King" in the item 'value'

  @skipOnStable2.0
  Scenario Outline: non-admin user tries to get users with a certain role
    Given the administrator has given "Brian" the role "<role>" using the settings api
    When the user "Brian" gets all users with role "Admin" using the Graph API
    Then the HTTP status code should be "401"
    And the JSON data of the response should match
    """
    {
      "type": "object",
      "required": [
        "error"
      ],
      "properties": {
        "error": {
          "type": "object",
          "required": [
            "message"
          ],
          "properties": {
            "type": "string",
            "enum": ["Unauthorized"]
          }
        }
      }
    }
    """
    Examples:
      | role        |
      | Space Admin |
      | User        |
      | Guest       |
