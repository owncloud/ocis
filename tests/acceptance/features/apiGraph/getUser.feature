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


  Scenario: admin user gets the information of a user
    When user "Alice" gets information of user "Brian" using Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should contain display name "Brian Murphy" and match
    """
    {
      "type": "object",
      "required": [
        "displayName",
        "id",
        "mail",
        "onPremisesSamAccountName"
      ],
      "properties": {
        "displayName": {
          "type": "string",
          "enum": ["Brian Murphy"]
        },
        "id" : {
          "type": "string",
          "pattern": "^%uuid_v4%$"
        },
        "mail": {
          "type": "string",
          "enum": ["brian@example.org"]
        },
        "onPremisesSamAccountName": {
          "type": "string",
          "enum": ["Brian"]
        }
      }
    }
    """


  Scenario: non-admin user tries to get the information of a user
    When user "Brian" tries to get information of user "Alice" using Graph API
    Then the HTTP status code should be "401"
    And the last response should be an unauthorized response


  Scenario: admin user gets all users
    When user "Alice" gets all users using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should contain display name "Alice Hansen" and match
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
          "pattern": "^%uuid_v4%$"
        },
        "mail": {
          "type": "string",
          "enum": ["alice@example.org"]
        },
        "onPremisesSamAccountName": {
          "type": "string",
          "enum": ["Alice"]
        }
      }
    }
    """
    And the JSON data of the response should contain display name "Brian Murphy" and match
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
          "pattern": "^%uuid_v4%$"
        },
        "mail": {
          "type": "string",
          "enum": ["brian@example.org"]
        },
        "onPremisesSamAccountName": {
          "type": "string",
          "enum": ["Brian"]
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


  Scenario: admin user gets the drive information of a user
    When the user "Alice" gets user "Brian" along with his drive information using Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should contain display name "Brian Murphy" and match
    """
      {
        "type": "object",
        "required": [
          "displayName",
          "id",
          "mail",
          "onPremisesSamAccountName",
          "drive"
        ],
        "properties": {
          "displayName": {
            "type": "string",
            "enum": ["Brian Murphy"]
          },
          "id" : {
            "type": "string",
            "pattern": "^%uuid_v4%$"
          },
          "mail": {
            "type": "string",
            "enum": ["brian@example.org"]
          },
          "onPremisesSamAccountName": {
            "type": "string",
            "enum": ["Brian"]
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
                "pattern": "^%space_id%$"
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
                      "enum": ["%user_id%"]
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
                    "pattern": "^%base_url%/dav/spaces/%space_id%$"
                  }
                }
              },
              "webUrl": {
                "type": "string",
                "pattern": "^%base_url%/f/%space_id%$"
              }
            }
          }
        }
      }
    """


  Scenario: normal user gets his/her own drive information
    When the user "Brian" gets his drive information using Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should contain display name "Brian Murphy" and match
    """
      {
        "type": "object",
        "required": [
          "displayName",
          "id",
          "mail",
          "onPremisesSamAccountName",
          "drive"
        ],
        "properties": {
          "displayName": {
            "type": "string",
            "enum": ["Brian Murphy"]
          },
          "id" : {
            "type": "string",
            "pattern": "^%uuid_v4%$"
          },
          "mail": {
            "type": "string",
            "enum": ["brian@example.org"]
          },
          "onPremisesSamAccountName": {
            "type": "string",
            "enum": ["Brian"]
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
                "pattern": "^%space_id%$"
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
                      "enum": ["%user_id%"]
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
                    "pattern": "^%base_url%/dav/spaces/%space_id%$"
                  }
                }
              },
              "webUrl": {
                "type": "string",
                "pattern": "^%base_url%/f/%space_id%$"
              }
            }
          }
        }
      }
    """


  Scenario: admin user gets the group information of a user
    Given group "tea-lover" has been created
    And group "coffee-lover" has been created
    And user "Brian" has been added to group "tea-lover"
    And user "Brian" has been added to group "coffee-lover"
    When the user "Alice" gets user "Brian" along with his group information using Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should contain display name "Brian Murphy" and match
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
          "pattern": "^%uuid_v4%$"
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
