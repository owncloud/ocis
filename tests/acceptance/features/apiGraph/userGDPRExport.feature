@api
Feature: user GDPR (General Data Protection Regulation) report
  As a user
  I want to generate my GDPR report
  So that I can review what events are stored by the server

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files
    And using spaces DAV path


  Scenario: generate a GDPR report and check user data in the downloaded report
    When user "Alice" exports her GDPR report to "/.personal_data_export.json" using the Graph API
    And user "Alice" downloads the content of GDPR report ".personal_data_export.json"
    Then the HTTP status code of responses on each endpoint should be "201, 200" respectively
    And the downloaded JSON content should contain key 'user' and match
    """
    {
      "type": "object",
      "required": [
        "id",
        "username",
        "mail",
        "display_name",
        "uid_number",
        "gid_number"
      ],
      "properties": {
        "id": {
          "type": "object",
          "required": [
            "idp",
            "opaque_id",
            "type"
          ],
          "properties": {
            "idp": {
              "type": "string",
              "pattern": "^%base_url%$"
            },
            "opaque_id": {
              "type": "string",
              "pattern": "^%user_id_pattern%$"
            },
            "type": {
              "type": "number",
              "enum": [1]
            }
          }
        },
        "username": {
          "type": "string",
          "enum": ["Alice"]
        },
        "mail": {
          "type": "string",
          "enum": ["alice@example.org"]
        },
        "display_name": {
          "type": "string",
          "enum": ["Alice Hansen"]
        },
        "uid_number": {
          "type": "number",
          "enum": [99]
        },
        "gid_number": {
          "type": "number",
          "enum": [99]
        }
      }
    }
    """

  Scenario: generate a GDPR report and check events when a user is created
    When user "Alice" exports her GDPR report to "/.personal_data_export.json" using the Graph API
    And user "Alice" downloads the content of GDPR report ".personal_data_export.json"
    Then the HTTP status code of responses on each endpoint should be "201, 200" respectively
    And the downloaded JSON content should contain event type "events.UserCreated" in item 'events' and should match
    """
    {
      "type": "object",
      "required": [
        "event"
      ],
      "properties": {
        "event" : {
          "type": "object",
          "required": [
            "Executant",
            "UserID"
          ],
          "properties": {
            "Executant": {
              "type": "object",
              "required": [
                "idp",
                "opaque_id",
                "type"
              ],
              "properties": {
                "idp": {
                  "type": "string",
                  "pattern": "^%base_url%$"
                },
                "opaque_id": {
                  "type": "string",
                  "pattern": "^%user_id_pattern%$"
                },
                "type": {
                  "type": "number",
                  "enum": [1]
                }
              }
            },
            "UserID": {
              "type": "string",
              "pattern": "^%user_id_pattern%$"
            }
          }
        }
      }
    }
    """
    And the downloaded JSON content should contain event type "events.SpaceCreated" in item 'events' and should match
    """
    {
      "type": "object",
      "required": [
        "event"
      ],
      "properties": {
        "event" : {
          "type": "object",
          "required": [
            "Executant",
            "Name",
            "Type",
            "Quota"
          ],
          "properties": {
            "Executant": {
              "type": "object",
              "required": [
                "idp",
                "opaque_id",
                "type"
              ],
              "properties": {
                "idp": {
                  "type": "string",
                  "pattern": "^%base_url%$"
                },
                "opaque_id": {
                  "type": "string",
                  "pattern": "^%user_id_pattern%$"
                },
                "type": {
                  "type": "number",
                  "enum": [1]
                }
              }
            },
            "Name": {
              "type": "string",
              "enum": ["Alice Hansen"]
            },
            "Type": {
              "type": "string",
              "enum": ["personal"]
            },
            "Quota": {
              "type": ["number", "null"],
              "enum": [null]
            }
          }
        }
      }
    }
    """


  Scenario: generate a GDPR report and check events when user uploads a file
    Given user "Alice" has uploaded file with content "sample text" to "lorem.txt"
    When user "Alice" exports her GDPR report to "/.personal_data_export.json" using the Graph API
    And user "Alice" downloads the content of GDPR report ".personal_data_export.json"
    Then the HTTP status code of responses on each endpoint should be "201, 200" respectively
    And the downloaded JSON content should contain event type "events.BytesReceived" in item 'events' and should match
    """
    {
      "type": "object",
      "required": [
        "event"
      ],
      "properties": {
        "event" : {
          "type": "object",
          "required": [
            "ExecutingUser",
            "Filename",
            "ResourceID",
            "Filesize",
            "UploadID",
            "SpaceOwner"
          ],
          "properties": {
            "ExecutingUser": {
              "type": "object",
              "required": [
                "username"
              ],
              "properties": {
                "username": {
                  "type": "string",
                  "enum": ["Alice"]
                }
              }
            },
            "Filename": {
              "type": "string",
              "enum": ["lorem.txt"]
            },
            "Filesize": {
              "type": "number",
              "enum": [11]
            },
            "Quota": {
              "type": ["number", "null"],
              "enum": [null]
            }
          }
        }
      }
    }
    """
    And the downloaded JSON content should contain event type "events.FileUploaded" in item 'events' and should match
    """
    {
      "type": "object",
      "required": [
        "event"
      ],
      "properties": {
        "event" : {
          "type": "object",
          "required": [
            "Executant",
            "Owner",
            "Ref",
            "SpaceOwner"
          ],
          "properties": {
            "Ref": {
              "type": "object",
              "required": [
                "path"
              ],
              "properties": {
                "path" : {
                  "type": "string",
                  "enum": ["./lorem.txt"]
                }
              }
            }
          }
        }
      }
    }
    """
    And the downloaded JSON content should contain event type "events.PostprocessingFinished" in item 'events' and should match
    """
    {
      "type": "object",
      "required": [
        "event"
      ],
      "properties": {
        "event" : {
          "type": "object",
          "required": [
            "ExecutingUser",
            "Filename"
          ],
          "properties": {
            "ExecutingUser": {
              "type": "object",
              "required": [
                "username"
              ],
              "properties": {
                "username": {
                  "type": "string",
                  "enum": ["Alice"]
                }
              }
            },
            "Filename": {
              "type": "string",
              "enum": ["lorem.txt"]
            }
          }
        }
      }
    }
    """


  Scenario: generate a GDPR report and check events when a user is added to a group
    Given group "tea-lover" has been created
    And user "Alice" has been added to group "tea-lover"
    When user "Alice" exports her GDPR report to "/.personal_data_export.json" using the Graph API
    And user "Alice" downloads the content of GDPR report ".personal_data_export.json"
    Then the HTTP status code of responses on each endpoint should be "201, 200" respectively
    And the downloaded JSON content should contain event type "events.GroupMemberAdded" in item 'events' and should match
    """
    {
      "type": "object",
      "required": [
        "event"
      ],
      "properties": {
        "event" : {
          "type": "object",
          "required": [
            "Executant",
            "GroupID",
            "UserID"
          ],
          "properties": {
            "Executant": {
              "type": "object",
              "required": [
                "idp",
                "opaque_id",
                "type"
              ],
              "properties": {
                "idp": {
                  "type": "string",
                  "pattern": "^%base_url%$"
                },
                "opaque_id": {
                  "type": "string",
                  "pattern": "^%user_id_pattern%$"
                },
                "type": {
                  "type": "number",
                  "enum": [1]
                }
              }
            },
            "GroupID": {
              "type": "string",
              "pattern": "^%group_id_pattern%$"
            },
            "UserID": {
              "type": "string",
              "pattern": "^%user_id_pattern%$"
            }
          }
        }
      }
    }
    """
    And the downloaded JSON content should contain key 'user' and match
    """
    {
      "type": "object",
      "required": [
        "id",
        "username",
        "mail",
        "display_name",
        "groups",
        "uid_number",
        "gid_number"
      ],
      "properties": {
        "username": {
          "type": "string",
          "enum": ["Alice"]
        },
        "mail": {
          "type": "string",
          "enum": ["alice@example.org"]
        },
        "display_name": {
          "type": "string",
          "enum": ["Alice Hansen"]
        },
        "groups": {
          "type": "array",
          "items": [
            {
              "type": "string",
              "pattern": "^%group_id_pattern%$"
            }
          ]
        },
        "uid_number": {
          "type": "number",
          "enum": [99]
        },
        "gid_number": {
          "type": "number",
          "enum": [99]
        }
      }
    }
    """


  Scenario: generate a GDPR report after the admin updates the quota of personal space
    Given user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "10000"
    When user "Alice" exports her GDPR report to "/.personal_data_export.json" using the Graph API
    And user "Alice" downloads the content of GDPR report ".personal_data_export.json"
    Then the HTTP status code of responses on each endpoint should be "201, 200" respectively
    And the downloaded JSON content should contain event type "events.SpaceUpdated" in item 'events' and should match
    """
    {
      "type": "object",
      "required": [
        "event"
      ],
      "properties": {
        "event" : {
          "type": "object",
          "required": [
            "Executant",
            "Space"
          ],
          "properties": {
            "Executant": {
              "type": "object",
              "required": [
                "idp",
                "opaque_id",
                "type"
              ],
              "properties": {
                "idp": {
                  "type": "string",
                  "pattern": "^%base_url%$"
                },
                "opaque_id": {
                  "type": "string",
                  "pattern": "^%user_id_pattern%$"
                },
                "type": {
                  "type": "number",
                  "enum": [1]
                }
              }
            },
            "Space": {
              "type": "object",
              "required": [
                "name",
                "quota",
                "space_type"
              ],
              "properties": {
                "name": {
                  "type": "string",
                  "enum": ["Alice Hansen"]
                },
                "quota": {
                  "type": "object",
                  "required": [
                    "quota_max_bytes",
                    "quota_max_files"
                  ],
                  "properties": {
                    "quota_max_bytes": {
                      "type": "number",
                      "enum": [10000]
                    },
                    "quota_max_files": {
                      "type": "number",
                      "enum": [18446744073709552000]
                    }
                  }
                },
                "space_type": {
                  "type": "string",
                  "enum": ["personal"]
                }
              }
            }
          }
        }
      }
    }
    """


  Scenario Outline: user tries to generate GDPR report of other users
    Given user "Brian" has been created with default attributes and without skeleton files
    And the administrator has given "Alice" the role "<userRole>" using the settings api
    And the administrator has given "Brian" the role "<role>" using the settings api
    When user "Alice" tries to export GDPR report of user "Brian" to "/.personal_data_export.json" using Graph API
    Then the HTTP status code should be "400"
    Examples:
      | userRole    | role        |
      | Space Admin | Space Admin |
      | Space Admin | User        |
      | Space Admin | Guest       |
      | Space Admin | Admin       |
      | User        | Space Admin |
      | User        | User        |
      | User        | Guest       |
      | User        | Admin       |
      | Guest       | Space Admin |
      | Guest       | User        |
      | Guest       | Guest       |
      | Guest       | Admin       |
      | Admin       | Space Admin |
      | Admin       | User        |
      | Admin       | Guest       |
      | Admin       | Admin       |
