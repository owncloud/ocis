Feature: get users
  As an admin
  I want to be able to retrieve user information
  So that I can see the information

  Background:
    Given these users have been created with default attributes:
      | username |
      | Alice    |
      | Brian    |


  Scenario: admin user gets the information of a user
    Given the administrator has assigned the role "Admin" to user "Alice" using the Graph API
    When user "Alice" gets information of user "Brian" using Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["displayName","id","onPremisesSamAccountName","accountEnabled","userType"],
        "properties": {
          "displayName": {"const": "Brian Murphy"},
          "id": {"pattern": "^%user_id_pattern%$"},
          "onPremisesSamAccountName": {"const": "Brian"},
          "accountEnabled": {"const": true},
          "userType": {"const": "Member"}
        }
      }
      """

  @issue-5125
  Scenario Outline: non-admin user tries to get the information of a user
    Given the administrator has assigned the role "<user-role-2>" to user "Alice" using the Graph API
    And the administrator has assigned the role "<user-role>" to user "Brian" using the Graph API
    When user "Brian" tries to get information of user "Alice" using Graph API
    Then the HTTP status code should be "403"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["error"],
        "properties": {
          "error": {
            "type": "object",
            "required": ["message"],
            "properties": {
              "message": {"const": "Forbidden"}
            }
          }
        }
      }
      """
    Examples:
      | user-role   | user-role-2 |
      | Space Admin | Space Admin |
      | Space Admin | User        |
      | Space Admin | User Light  |
      | Space Admin | Admin       |
      | User        | Space Admin |
      | User        | User        |
      | User        | User Light  |
      | User        | Admin       |
      | User Light  | Space Admin |
      | User Light  | User        |
      | User Light  | User Light  |
      | User Light  | Admin       |


  Scenario: admin user gets all users
    Given the administrator has assigned the role "Admin" to user "Alice" using the Graph API
    When user "Alice" gets all users using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should contain the user "Alice Hansen" in the item 'value', the user-details should match
      """
      {
        "type": "object",
        "required": ["displayName","id","onPremisesSamAccountName","accountEnabled","userType"],
        "properties": {
          "id": {"pattern": "^%user_id_pattern%$"},
          "mail": {"const": "alice@example.org"},
          "onPremisesSamAccountName": {"const": "Alice"},
          "accountEnabled": {"const": true},
          "userType": {"const": "Member"}
        }
      }
      """
    And the JSON data of the response should contain the user "Brian Murphy" in the item 'value', the user-details should match
      """
      {
        "type": "object",
        "required": ["displayName","id","onPremisesSamAccountName","accountEnabled","userType"],
        "properties": {
          "id": {"pattern": "^%user_id_pattern%$"},
          "mail": {"const": "brian@example.org"},
          "onPremisesSamAccountName": {"const": "Brian"},
          "accountEnabled": {"const": true},
          "userType": {"const": "Member"}
        }
      }
      """


  Scenario: admin user gets all users include disabled users
    Given the administrator has assigned the role "Admin" to user "Alice" using the Graph API
    And the user "Alice" has disabled user "Brian"
    When user "Alice" gets all users using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should contain the user "Alice Hansen" in the item 'value', the user-details should match
      """
      {
        "type": "object",
        "required": ["displayName","id","onPremisesSamAccountName","accountEnabled","userType"],
        "properties": {
          "id": {"pattern": "^%user_id_pattern%$"},
          "mail": {"const": "alice@example.org"},
          "onPremisesSamAccountName": {"const": "Alice"},
          "accountEnabled": {"const": true},
          "userType": {"const": "Member"}
        }
      }
      """
    And the JSON data of the response should contain the user "Brian Murphy" in the item 'value', the user-details should match
      """
      {
        "type": "object",
        "required": ["displayName","id","onPremisesSamAccountName","accountEnabled","userType"],
        "properties": {
          "id": {"pattern": "^%user_id_pattern%$"},
          "mail": {"const": "brian@example.org"},
          "onPremisesSamAccountName": {"const": "Brian"},
          "accountEnabled": {"const": false},
          "userType": {"const": "Member"}
        }
      }
      """


  Scenario Outline: non-admin user tries to get all users
    Given the administrator has assigned the role "<user-role>" to user "Alice" using the Graph API
    When user "Brian" tries to get all users using the Graph API
    Then the HTTP status code should be "403"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["error"],
        "properties": {
          "error": {
            "type": "object",
            "required": ["message"],
            "properties": {
              "message": {"const": "search term too short"}
            }
          }
        }
      }
      """
    Examples:
      | user-role   |
      | Space Admin |
      | User        |
      | User Light  |


  Scenario: admin user gets the drive information of a user
    Given the administrator has assigned the role "Admin" to user "Alice" using the Graph API
    When the user "Alice" gets user "Brian" along with his drive information using Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["displayName","id","onPremisesSamAccountName","drive","accountEnabled","userType"],
        "properties": {
          "displayName": {"const": "Brian Murphy"},
          "id": {"pattern": "^%user_id_pattern%$"},
          "onPremisesSamAccountName": {"const": "Brian"},
          "accountEnabled": {"const": true},
          "userType": {"const": "Member"},
          "drive": {
            "type": "object",
            "required": ["driveAlias","id","name","owner","quota","root","webUrl"],
            "properties": {
              "driveType": {"const": "personal"},
              "driveAlias": {"const": "personal/brian"},
              "id": {"pattern": "^%space_id_pattern%$"},
              "name": {"const": "Brian Murphy"},
              "owner": {
                "type": "object",
                "required": ["user"],
                "properties": {
                  "user": {
                    "type": "object",
                    "required": ["id"],
                    "properties": {
                      "id": {"pattern": "%user_id_pattern%"}
                    }
                  }
                }
              },
              "quota": {
                "type": "object",
                "required": ["state"],
                "properties": {
                  "state": {"const": "normal"}
                }
              },
              "root": {
                "type": "object",
                "required": ["id","webDavUrl"],
                "properties": {
                  "state": {"const": "normal"},
                  "webDavUrl": {"pattern": "^%base_url%/dav/spaces/%space_id_pattern%$"}
                }
              },
              "webUrl": {"pattern": "^%base_url%/f/%space_id_pattern%$"}
            }
          }
        }
      }
      """


  Scenario Outline: user gets his/her own information along with drive information
    Given the administrator has assigned the role "<user-role>" to user "Brian" using the Graph API
    When the user "Brian" gets his drive information using Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "displayName","id","onPremisesSamAccountName","drive","accountEnabled","userType"
        ],
        "properties": {
          "displayName": {"const": "Brian Murphy"},
          "id": {"pattern": "^%user_id_pattern%$"},
          "onPremisesSamAccountName": {"const": "Brian"},
          "accountEnabled": {"const": true},
          "userType": {"const": "Member"},
          "drive": {
            "type": "object",
            "required": ["driveAlias","id","name","owner","quota","root","webUrl"],
            "properties": {
              "driveType": {"const": "personal"},
              "driveAlias": {"const": "personal/brian"},
              "id": {"pattern": "^%space_id_pattern%$"},
              "name": {"const": "Brian Murphy"},
              "owner": {
                "type": "object",
                "required": ["user"],
                "properties": {
                  "user": {
                    "type": "object",
                    "required": ["id"],
                    "properties": {
                      "id": {"pattern": "%user_id_pattern%"}
                    }
                  }
                }
              },
              "quota": {
                "type": "object",
                "required": ["state"],
                "properties": {
                  "state": {"const": "normal"}
                }
              },
              "root": {
                "type": "object",
                "required": ["id","webDavUrl"],
                "properties": {
                  "state": {"const": "normal"},
                  "webDavUrl": {"pattern": "^%base_url%/dav/spaces/%space_id_pattern%$"}
                }
              },
              "webUrl": {"pattern": "^%base_url%/f/%space_id_pattern%$"}
            }
          }
        }
      }
      """
    Examples:
      | user-role   |
      | Admin       |
      | Space Admin |
      | User        |
      | User Light  |


  Scenario: admin user gets the group information of a user
    Given the administrator has assigned the role "Admin" to user "Alice" using the Graph API
    And group "tea-lover" has been created
    And group "coffee-lover" has been created
    And user "Brian" has been added to group "tea-lover"
    And user "Brian" has been added to group "coffee-lover"
    When the user "Alice" gets user "Brian" along with his group information using Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["id","onPremisesSamAccountName"],
        "properties": {
          "id": {"pattern": "^%user_id_pattern%$"},
          "onPremisesSamAccountName": {"const": "Brian"},
          "memberOf": {
            "type": "array",
            "minItems": 2,
            "maxItems": 2,
            "uniqueItems": true,
            "items": {
              "oneOf": [
                {
                  "type": "object",
                  "required": ["displayName"],
                  "properties": {
                    "displayName": {"const": "tea-lover"}
                  }
                },
                {
                  "type": "object",
                  "required": ["displayName"],
                  "properties": {
                    "displayName": {"const": "coffee-lover"}
                  }
                }
              ]
            }
          }
        }
      }
      """

  @issue-5125
  Scenario Outline: non-admin user tries to get the group information of a user
    Given the administrator has assigned the role "<user-role>" to user "Alice" using the Graph API
    And the administrator has assigned the role "<user-role-2>" to user "Brian" using the Graph API
    And group "coffee-lover" has been created
    And user "Brian" has been added to group "coffee-lover"
    When the user "Alice" gets user "Brian" along with his group information using Graph API
    Then the HTTP status code should be "403"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["error"],
        "properties": {
          "error": {
            "type": "object",
            "required": ["message"],
            "properties": {
              "message": {"const": "Forbidden"}
            }
          }
        }
      }
      """
    Examples:
      | user-role   | user-role-2 |
      | Space Admin | Space Admin |
      | Space Admin | User        |
      | Space Admin | User Light  |
      | Space Admin | Admin       |
      | User        | Space Admin |
      | User        | User        |
      | User        | User Light  |
      | User        | Admin       |
      | User Light  | Space Admin |
      | User Light  | User        |
      | User Light  | User Light  |
      | User Light  | Admin       |


  Scenario: admin user tries to get the information of nonexistent user
    Given the administrator has assigned the role "Admin" to user "Alice" using the Graph API
    When user "Alice" tries to get information of user "nonexistent" using Graph API
    Then the HTTP status code should be "404"

  @issue-5125
  Scenario Outline: non-admin user tries to get the information of nonexistent user
    Given the administrator has assigned the role "<user-role>" to user "Alice" using the Graph API
    When user "Alice" tries to get information of user "nonexistent" using Graph API
    Then the HTTP status code should be "403"
    Examples:
      | user-role   |
      | Space Admin |
      | User        |
      | User Light  |


  Scenario: admin user gets all users of certain groups
    Given the administrator has assigned the role "Admin" to user "Alice" using the Graph API
    And user "Carol" has been created with default attributes
    And the user "Alice" has disabled user "Carol"
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
                  "required": ["id","mail","onPremisesSamAccountName","accountEnabled","userType"],
                  "properties": {
                    "id": {"pattern": "^%user_id_pattern%$"},
                    "mail": {"const": "alice@example.org"},
                    "onPremisesSamAccountName": {"const": "Alice"},
                    "accountEnabled": {"const": true},
                    "userType": {"const": "Member"}
                  }
                },
                {
                  "type": "object",
                  "required": ["id","mail","onPremisesSamAccountName","accountEnabled"],
                  "properties": {
                    "id": {"pattern": "^%user_id_pattern%$"},
                    "mail": {"const": "brian@example.org"},
                    "onPremisesSamAccountName": {"const": "Brian"},
                    "accountEnabled": {"const": true},
                    "userType": {"const": "Member"}
                  }
                }
              ]
            }
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
        "required": ["value"],
        "properties": {
          "value": {
            "type": "array",
            "minItems": 1,
            "maxItems": 1,
            "items": {
              "type": "object",
              "required": [
                "id","mail","onPremisesSamAccountName","accountEnabled","userType"],
              "properties": {
                "id": {"pattern": "^%user_id_pattern%$"},
                "mail": {"const": "brian@example.org"},
                "onPremisesSamAccountName": {"const": "Brian"},
                "accountEnabled": {"const": true},
                "userType": {"const": "Member"}
              }
            }
          }
        }
      }
      """
    And the JSON data of the response should not contain the user "Carol King" in the item 'value'
    And the JSON data of the response should not contain the user "Alice Hansen" in the item 'value'


  Scenario: admin user gets all users of two groups
    Given the administrator has assigned the role "Admin" to user "Alice" using the Graph API
    And user "Carol" has been created with default attributes
    And group "tea-lover" has been created
    And group "coffee-lover" has been created
    And group "wine-lover" has been created
    And user "Alice" has been added to group "tea-lover"
    And user "Brian" has been added to group "coffee-lover"
    And user "Carol" has been added to group "wine-lover"
    When the user "Alice" gets all users that are members in the group "tea-lover" or the group "coffee-lover" using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should contain the user "Alice Hansen" in the item 'value', the user-details should match
      """
      {
        "type": "object",
        "required": ["displayName","id","onPremisesSamAccountName","accountEnabled","userType"],
        "properties": {
          "id": {"pattern": "^%user_id_pattern%$"},
          "mail": {"const": "alice@example.org"},
          "onPremisesSamAccountName": {"const": "Alice"},
          "accountEnabled": {"const": true},
          "userType": {"const": "Member"}
        }
      }
      """
    And the JSON data of the response should contain the user "Brian Murphy" in the item 'value', the user-details should match
      """
      {
        "type": "object",
        "required": ["displayName","id","onPremisesSamAccountName","accountEnabled","userType"],
        "properties": {
          "id": {"pattern": "^%user_id_pattern%$"},
          "mail": {"const": "brian@example.org"},
          "onPremisesSamAccountName": {"const": "Brian"},
          "accountEnabled": {"const": true},
          "userType": {"const": "Member"}
        }
      }
      """
    But the JSON data of the response should not contain the user "Carol King" in the item 'value'


  Scenario Outline: non admin user tries to get users of certain groups
    Given the administrator has assigned the role "Admin" to user "Alice" using the Graph API
    And the administrator has assigned the role "<user-role>" to user "Brian" using the Graph API
    And group "tea-lover" has been created
    And user "Alice" has been added to group "tea-lover"
    When the user "Brian" gets all users of the group "tea-lover" using the Graph API
    Then the HTTP status code should be "403"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["error"],
        "properties": {
          "error": {
            "type": "object",
            "required": ["message"],
            "properties": {
              "message": {"const": "search term too short"}
            }
          }
        }
      }
      """
    Examples:
      | user-role   |
      | Space Admin |
      | User        |
      | User Light  |


  Scenario: admin user gets all users with certain roles and members of a certain group
    Given the administrator has assigned the role "Admin" to user "Alice" using the Graph API
    And user "Carol" has been created with default attributes
    And the administrator has assigned the role "Space Admin" to user "Brian" using the Graph API
    And the administrator has assigned the role "Space Admin" to user "Carol" using the Graph API
    And group "tea-lover" has been created
    And user "Brian" has been added to group "tea-lover"
    When the user "Alice" gets all users with role "Space Admin" using the Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should contain the user "Brian Murphy" in the item 'value', the user-details should match
      """
      {
        "type": "object",
        "required": ["displayName","id","onPremisesSamAccountName","accountEnabled","userType"],
        "properties": {
          "id": {"pattern": "^%user_id_pattern%$"},
          "mail": {"const": "brian@example.org"},
          "onPremisesSamAccountName": {"const": "Brian"},
          "accountEnabled": {"const": true},
          "userType": {"const": "Member"}
        }
      }
      """
    And the JSON data of the response should contain the user "Carol King" in the item 'value', the user-details should match
      """
      {
        "type": "object",
        "required": ["displayName","id","onPremisesSamAccountName","accountEnabled","userType"],
        "properties": {
          "id": {"pattern": "^%user_id_pattern%$"},
          "mail": {"const": "carol@example.org"},
          "onPremisesSamAccountName": {"const": "Carol"},
          "accountEnabled": {"const": true},
          "userType": {"const": "Member"}
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
        "required": ["displayName","id","onPremisesSamAccountName","accountEnabled","userType"],
        "properties": {
          "id": {"pattern": "^%user_id_pattern%$"},
          "mail": {"const": "brian@example.org"},
          "onPremisesSamAccountName": {"const": "Brian"},
          "accountEnabled": {"const": true},
          "userType": {"const": "Member"}
        }
      }
      """
    But the JSON data of the response should not contain the user "Carol King" in the item 'value'


  Scenario Outline: non-admin user tries to get users with a certain role
    Given the administrator has assigned the role "<user-role>" to user "Alice" using the Graph API
    When the user "Alice" gets all users with role "<user-role-2>" using the Graph API
    Then the HTTP status code should be "403"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["error"],
        "properties": {
          "error": {
            "type": "object",
            "required": ["message"],
            "properties": {
              "message": {"const": "search term too short"}
            }
          }
        }
      }
      """
    Examples:
      | user-role   | user-role-2 |
      | Space Admin | Space Admin |
      | Space Admin | User        |
      | Space Admin | User Light  |
      | Space Admin | Admin       |
      | User        | Space Admin |
      | User        | User        |
      | User        | User Light  |
      | User        | Admin       |
      | User Light  | Space Admin |
      | User Light  | User        |
      | User Light  | User Light  |
      | User Light  | Admin       |

  @issue-6017
  Scenario Outline: admin user gets the drive information of a user with different user role
    Given the administrator has assigned the role "<user-role>" to user "Alice" using the Graph API
    And the administrator has assigned the role "<user-role-2>" to user "Brian" using the Graph API
    And user "Brian" has created folder "my_data"
    When user "Alice" gets the personal drive information of user "Brian" using Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["driveAlias","driveType","id","name","webUrl","owner","quota","root"],
        "properties": {
          "driveAlias": {"const": "personal/brian"},
          "driveType": {"const": "personal"},
          "id": {"pattern": "^%space_id_pattern%$"},
          "name": {"const": "Brian Murphy"},
          "webUrl": {"pattern": "^%base_url%/f/%space_id_pattern%$"},
          "owner": {
            "type": "object",
            "required": ["user"],
            "properties": {
              "user": {
                "type": "object",
                "required": ["displayName","id"],
                "properties": {
                  "displayName": {"const": ""},
                  "id": {"pattern": "^%user_id_pattern%$"}
                }
              }
            }
          },
          "qouta": {
            "type": "object",
            "required": ["state"],
            "properties": {
              "state": {"const": "normal"}
            }
          },
          "root": {
            "type": "object",
            "required": ["webDavUrl"],
            "properties": {
              "webDavUrl": {"pattern": "^%base_url%/dav/spaces/%space_id_pattern%$"}
            }
          }
        }
      }
      """
    Examples:
      | user-role   | user-role-2 |
      | Admin       | Admin       |
      | Admin       | Space Admin |
      | Admin       | User        |
      | Admin       | User Light  |
      | Space Admin | Admin       |
      | Space Admin | Space Admin |
      | Space Admin | User        |
      | Space Admin | User Light  |


  Scenario Outline: non-admin user tries to get drive information of other user with different user role
    Given the administrator has assigned the role "<user-role>" to user "Alice" using the Graph API
    And the administrator has assigned the role "<user-role-2>" to user "Brian" using the Graph API
    When user "Alice" gets the personal drive information of user "Brian" using Graph API
    Then the HTTP status code should be "404"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["error"],
        "properties": {
          "error": {
            "type": "object",
            "required": ["code","message"],
            "properties": {
              "code": {"const": "itemNotFound"},
              "message": {"const": "no drive returned from storage"}
            }
          }
        }
      }
      """
    Examples:
      | user-role  | user-role-2 |
      | User       | Admin       |
      | User       | Space Admin |
      | User       | User        |
      | User       | User Light  |
      | User Light | Admin       |
      | User Light | Space Admin |
      | User Light | User        |
      | User Light | User Light  |


  Scenario Outline: user with different user role gets his/her own drive information
    Given the administrator has assigned the role "<user-role>" to user "Alice" using the Graph API
    When user "Alice" gets own personal drive information using Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["driveAlias","driveType","id","name","webUrl","owner","quota","root"],
        "properties": {
          "driveAlias": {"const": "personal/alice"},
          "driveType": {"const": "personal"},
          "id": {"pattern": "^%space_id_pattern%$"},
          "name": {"const": "Alice Hansen"},
          "webUrl": {"pattern": "^%base_url%/f/%space_id_pattern%$"},
          "owner": {
            "type": "object",
            "required": ["user"],
            "properties": {
              "user": {
                "type": "object",
                "required": ["displayName","id"],
                "properties": {
                  "displayName": {"const": ""},
                  "id": {"pattern": "^%user_id_pattern%$"}
                }
              }
            }
          },
          "qouta": {
            "type": "object",
            "required": ["state"],
            "properties": {
              "state": {"const": "normal"}
            }
          },
          "root": {
            "type": "object",
            "required": ["webDavUrl"],
            "properties": {
              "webDavUrl": {"pattern": "^%base_url%/dav/spaces/%space_id_pattern%$"}
            }
          }
        }
      }
      """
    Examples:
      | user-role   |
      | Admin       |
      | Space Admin |
      | User        |
      | User Light  |


  Scenario: non-admin user searches other users by display name
    When user "Brian" searches for user "ali" using Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the search response should not contain user email
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
              "required": ["displayName","id","userType"],
              "properties": {
                "displayName": {"const": "Alice Hansen"},
                "id": {"pattern": "^%user_id_pattern%$"},
                "userType": {"const": "Member"}
              }
            }
          }
        }
      }
      """


  Scenario: non-admin user tries to search for a user by display name with less than 3 characters
    When user "Brian" tries to search for user "al" using Graph API
    Then the HTTP status code should be "403"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["error"],
        "properties": {
          "error": {
            "type": "object",
            "required": ["message"],
            "properties": {
              "message": {"const": "search term too short"}
            }
          }
        }
      }
      """

  @issue-7990
  Scenario Outline: user tries to search other users with invalid characters/token (search term without quotation)
    Given user "<user>" has been created with default attributes
    When user "Brian" tries to search for user "<user>" using Graph API
    Then the HTTP status code should be "400"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["error"],
        "properties": {
          "error": {
            "type": "object",
            "required": ["message"],
            "properties": {
              "message": {"const": "Token '<error-token>' is invalid"}
            }
          }
        }
      }
      """
    Examples:
      | user                  | error-token      |
      | Alice-From-Wonderland | -From-Wonderland |
      | Alice@From@Wonderland | @From@Wonderland |

  @issue-7990
  Scenario: non-admin user searches other users by e-mail (search term with quotation)
    When user "Brian" searches for user "%22alice@example.org%22" using Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the search response should not contain user email
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
              "required": ["displayName","id","userType"],
              "properties": {
                "displayName": {"const": "Alice Hansen"},
                "id": {"pattern": "^%user_id_pattern%$"},
                "userType": {"const": "Member"}
              }
            }
          }
        }
      }
      """


  Scenario: non-admin user searches for a disabled users
    Given the user "Admin" has disabled user "Alice"
    When user "Brian" searches for user "alice" using Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the search response should not contain user email
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
              "required": ["displayName","id","userType"],
              "properties": {
                "displayName": {"const": "Alice Hansen"},
                "id": {"pattern": "^%user_id_pattern%$"},
                "userType": {"const": "Member"}
              }
            }
          }
        }
      }
      """


  Scenario: non-admin user searches for multiple users having same displayname
    Given the user "Admin" has created a new user with the following attributes:
      | userName    | another-alice                |
      | displayName | Alice Hansen                 |
      | email       | another-alice@example.org    |
      | password    | containsCharacters(*:!;_+-&) |
    And the administrator has assigned the role "Admin" to user "Alice" using the Graph API
    When user "Brian" searches for user "alice" using Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the search response should not contain users email
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
                  "required": ["displayName","id","userType"],
                  "properties": {
                    "displayName": {"const": "Alice Hansen"},
                    "id": {"pattern": "^%user_id_pattern%$"},
                    "userType": {"const": "Member"}
                  }
                },
                {
                  "type": "object",
                  "required": ["displayName","id","userType"],
                  "properties": {
                    "displayName": {"const": "Alice Hansen"},
                    "id": {"pattern": "^%user_id_pattern%$"},
                    "userType": {"const": ["Admin"]}
                  }
                }
              ]
            }
          }
        }
      }
      """


  Scenario: user searches for a non-existent user/group
    When user "Brian" tries to search for user "nonexistent" using Graph API
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["value"],
        "properties": {
          "value": {
            "type": "array",
            "minItems": 0,
            "maxItems": 0
          }
        }
      }
      """

  @issue-7990
  Scenario Outline: user searches for other users having special characters in displayname (search term with quotation)
    Given the user "Admin" has created a new user with the following attributes:
      | userName    | specail-user                 |
      | displayName | <displayname>                |
      | email       | specialuser@example.org      |
    When user "Brian" searches for user '"<search-term>"' using Graph API
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
              "required": ["displayName", "id", "userType"],
              "properties": {
                "displayName": {"const": "<displayname>"},
                "id": {"pattern": "^%user_id_pattern%$"},
                "userType": {"const": "Member"}
              }
            }
          }
        }
      }
      """
    Examples:
      | displayname      | search-term |
      | -_.ocusr         | -_.         |
      | _ocusr@          | _oc         |
      | Alice-Wonderland | -Wonderland |
      | Alice@Wonderland | @Wonderland |

  @issue-7990 @issue-enterprise-6600
  Scenario Outline: non-admin user searches for other users having special characters in username (search term with quotation)
    Given the user "Admin" has created a new user with the following attributes:
      | userName    | <username>              |
      | displayName | special user            |
      | email       | specialuser@example.org |
    When user "Brian" searches for user '"<search-term>"' using Graph API
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
              "required": ["displayName", "id", "userType", "onPremisesSamAccountName"],
              "properties": {
                "attributes": {"const": []},
                "displayName": {"const": "special user"},
                "id": {"pattern": "^%user_id_pattern%$"},
                "userType": {"const": "Member"},
                "onPremisesSamAccountName": {"const": ""},
                "appRoleAssignments": {"const": null},
                "drives": {"const": null},
                "drive": {"const": null},
                "memberOf": {"const": null},
                "mail": {"const": null}
              }
            }
          }
        }
      }
      """
    Examples:
      | username         | search-term |
      | Alice-Wonderland | -Wonderland |
      | Alice@Wonderland | @Wonderland |

  @issue-enterprise-6600
  Scenario: admin user searches for group users having special characters in username (search term with quotation)
    Given the user "Admin" has created a new user with the following attributes:
      | userName    | special-user            |
      | displayName | alice wonderland        |
      | email       | specialuser@example.org |
    And group "tea-lover" has been created
    And user "special-user" has been added to group "tea-lover"
    When user "Admin" searches for user '"special-us"' of the group "tea-lover" using the Graph API
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
              "required": ["displayName", "id", "userType", "mail", "onPremisesSamAccountName", "surname"],
              "properties": {
                "displayName": {"const": "alice wonderland"},
                "id": {"pattern": "^%user_id_pattern%$"},
                "userType": {"const": "Member"},
                "mail": {"const": "specialuser@example.org"}
              }
            }
          }
        }
      }
      """

  @issue-enterprise-6600
  Scenario: non-admin user can't search for users of a group
    Given the user "Admin" has created a new user with the following attributes:
      | userName    | special-user            |
      | displayName | alice wonderland        |
      | email       | specialuser@example.org |
    And group "tea-lover" has been created
    And user "special-user" has been added to group "tea-lover"
    When user "Brian" tries to search for user '"special-us"' of the group "tea-lover" using the Graph API
    Then the HTTP status code should be "403"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": ["error"],
        "properties": {
          "error": {
            "type": "object",
            "required": ["message"],
            "properties": {
              "message": {"const": "filter has forbidden elements for regular users"}
            }
          }
        }
      }
      """

  @env-config
  Scenario: non-admin user searches for other users with attributes
    Given the config "OCIS_USER_SEARCH_DISPLAYED_ATTRIBUTES" has been set to "onPremisesSamAccountName,userType,mail" for "graph" service
    And the user "Admin" has created a new user with the following attributes:
      | userName    | Carol             |
      | displayName | carol             |
      | email       | carol@example.org |
    When user "Brian" searches for user '"carol"' using Graph API
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
              "required": ["displayName", "id", "userType", "onPremisesSamAccountName"],
              "properties": {
                "attributes": {
                  "type": "array",
                  "minItems": 3,
                  "maxItems": 3,
                  "uniqueItems": true,
                  "items": { "type": "string" },
                  "allOf": [
                    { "contains": { "const": "carol@example.org" } },
                    { "contains": { "const": "Carol" } },
                    { "contains": { "const": "Member" } }
                  ]
                },
                "displayName": {"const": "carol"},
                "id": {"pattern": "^%user_id_pattern%$"},
                "userType": {"const": "Member"},
                "onPremisesSamAccountName": {"const": ""},
                "appRoleAssignments": {"const": null},
                "drives": {"const": null},
                "drive": {"const": null},
                "memberOf": {"const": null},
                "email": {"const": "carol@example.org"}
              }
            }
          }
        }
      }
      """
