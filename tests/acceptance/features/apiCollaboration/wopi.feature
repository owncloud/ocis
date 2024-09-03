Feature: collaboration (wopi)
  As a user
  I want to access files with collaboration service apps
  So that I can collaborate with other users

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files


  Scenario Outline: open file with .odt extension
    Given user "Alice" has uploaded file "filesForUpload/simple.odt" to "simple.odt"
    And we save it into "FILEID"
    When user "Alice" sends HTTP method "POST" to URL "<app-endpoint>"
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "app_url",
          "method",
          "form_parameters"
        ],
        "properties": {
          "app_url": {
            "type": "string",
            "pattern": "^.*\\?WOPISrc=.*wopi%2Ffiles%2F[a-fA-F0-9]{64}$"
          },
          "method": {
            "const": "POST"
          },
          "form_parameters": {
            "type": "object",
            "required": [
              "access_token",
              "access_token_ttl"
            ],
            "properties": {
              "access_token": {
                "type": "string"
              },
              "access_token_ttl": {
                "type": "string"
              }
            }
          }
        }
      }
      """
    Examples:
      | app-endpoint                                     |
      | /app/open?file_id=<<FILEID>>&app_name=FakeOffice |
      | /app/open?file_id=<<FILEID>>                     |


  Scenario: open text file with app name in url query (MIME type not registered in app-registry)
    Given user "Alice" has uploaded file "filesForUpload/lorem.txt" to "lorem.txt"
    And we save it into "FILEID"
    When user "Alice" sends HTTP method "POST" to URL "/app/open?file_id=<<FILEID>>&app_name=FakeOffice"
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "app_url",
          "method",
          "form_parameters"
        ],
        "properties": {
          "app_url": {
            "type": "string",
            "pattern": "^.*\\?WOPISrc=.*wopi%2Ffiles%2F[a-fA-F0-9]{64}$"
          },
          "method": {
            "const": "POST"
          },
          "form_parameters": {
            "type": "object",
            "required": [
              "access_token",
              "access_token_ttl"
            ],
            "properties": {
              "access_token": {
                "type": "string"
              },
              "access_token_ttl": {
                "type": "string"
              }
            }
          }
        }
      }
      """

  @issue-9928
  Scenario: open text file without app name in url query (MIME type not registered in app-registry)
    Given user "Alice" has uploaded file "filesForUpload/lorem.txt" to "lorem.txt"
    And we save it into "FILEID"
    When user "Alice" sends HTTP method "POST" to URL "/app/open?file_id=<<FILEID>>"
    Then the HTTP status code should be "500"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "code",
          "message"
        ],
        "properties": {
          "code": {
            "const": "SERVER_ERROR"
          },
          "message": {
            "const": "Error contacting the requested application, please use a different one or try again later"
          }
        }
      }
      """


  Scenario Outline: sharee open file with .odt extension
    Given user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file "filesForUpload/simple.odt" to "simple.odt"
    And we save it into "FILEID"
    And user "Alice" has sent the following resource share invitation:
      | resource        | simple.odt |
      | space           | Personal   |
      | sharee          | Brian      |
      | shareType       | user       |
      | permissionsRole | Viewer     |
    When user "Brian" sends HTTP method "POST" to URL "<app-endpoint>"
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "app_url",
          "method",
          "form_parameters"
        ],
        "properties": {
          "app_url": {
            "type": "string",
            "pattern": "^.*\\?WOPISrc=.*wopi%2Ffiles%2F[a-fA-F0-9]{64}$"
          },
          "method": {
            "const": "POST"
          },
          "form_parameters": {
            "type": "object",
            "required": [
              "access_token",
              "access_token_ttl"
            ],
            "properties": {
              "access_token": {
                "type": "string"
              },
              "access_token_ttl": {
                "type": "string"
              }
            }
          }
        }
      }
      """
    Examples:
      | app-endpoint                                     |
      | /app/open?file_id=<<FILEID>>&app_name=FakeOffice |
      | /app/open?file_id=<<FILEID>>                     |


  Scenario Outline: public user opens file with .odt extension
    Given user "Alice" has uploaded file "filesForUpload/simple.odt" to "simple.odt"
    And we save it into "FILEID"
    And user "Alice" has created the following resource link share:
      | resource        | simple.odt |
      | space           | Personal   |
      | permissionsRole | view       |
      | password        | %public%   |
    When the public sends HTTP method "POST" to URL "<app-endpoint>" with password "%public%"
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "app_url",
          "method",
          "form_parameters"
        ],
        "properties": {
          "app_url": {
            "type": "string",
            "pattern": "^.*\\?WOPISrc=.*wopi%2Ffiles%2F[a-fA-F0-9]{64}$"
          },
          "method": {
            "const": "POST"
          },
          "form_parameters": {
            "type": "object",
            "required": [
              "access_token",
              "access_token_ttl"
            ],
            "properties": {
              "access_token": {
                "type": "string"
              },
              "access_token_ttl": {
                "type": "string"
              }
            }
          }
        }
      }
      """
    Examples:
      | app-endpoint                                     |
      | /app/open?file_id=<<FILEID>>&app_name=FakeOffice |
      | /app/open?file_id=<<FILEID>>                     |

  @issue-9928
  Scenario Outline: open unsupported file format
    Given user "Alice" has uploaded file "filesForUpload/simple.pdf" to "simple.pdf"
    And we save it into "FILEID"
    When user "Alice" sends HTTP method "POST" to URL "<app-endpoint>"
    Then the HTTP status code should be "500"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "code",
          "message"
        ],
        "properties": {
          "code": {
            "const": "SERVER_ERROR"
          },
          "message": {
            "const": "Error contacting the requested application, please use a different one or try again later"
          }
        }
      }
      """
    Examples:
      | app-endpoint                                     |
      | /app/open?file_id=<<FILEID>>&app_name=FakeOffice |
      | /app/open?file_id=<<FILEID>>                     |


  Scenario Outline: open file with non-existing file id
    Given user "Alice" has uploaded file "filesForUpload/simple.odt" to "simple.odt"
    And we save it into "FILEID"
    And user "Alice" has deleted file "/simple.odt"
    When user "Alice" sends HTTP method "POST" to URL "<app-endpoint>"
    Then the HTTP status code should be "404"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "code",
          "message"
        ],
        "properties": {
          "code": {
            "const": "RESOURCE_NOT_FOUND"
          },
          "message": {
            "const": "file does not exist"
          }
        }
      }
      """
    Examples:
      | app-endpoint                                     |
      | /app/open?file_id=<<FILEID>>&app_name=FakeOffice |
      | /app/open?file_id=<<FILEID>>                     |

  @issue-9495
  Scenario Outline: open file with .odt extension (open-with-web)
    Given user "Alice" has uploaded file "filesForUpload/simple.odt" to "simple.odt"
    And we save it into "FILEID"
    When user "Alice" sends HTTP method "POST" to URL "<app-endpoint>"
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "uri"
        ],
        "properties": {
          "uri": {
            "type": "string",
             "pattern": "%base_url%/external\\?<url-query>contextRouteName=files-spaces-personal&fileId=%uuidv4_pattern%%24%uuidv4_pattern%%21%uuidv4_pattern%$"
          }
        }
      }
      """
    Examples:
      | app-endpoint                                              | url-query       |
      | /app/open-with-web?file_id=<<FILEID>>&app_name=FakeOffice | app=FakeOffice& |
      | /app/open-with-web?file_id=<<FILEID>>                     |                 |


  Scenario: open text file using open-with-web with app name in url query (MIME type not registered in app-registry)
    Given user "Alice" has uploaded file "filesForUpload/lorem.txt" to "lorem.txt"
    And we save it into "FILEID"
    When user "Alice" sends HTTP method "POST" to URL "/app/open-with-web?file_id=<<FILEID>>&app_name=FakeOffice"
    Then the HTTP status code should be "200"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "uri"
        ],
        "properties": {
          "uri": {
            "type": "string",
             "pattern": "%base_url%/external\\?app=FakeOffice&contextRouteName=files-spaces-personal&fileId=%uuidv4_pattern%%24%uuidv4_pattern%%21%uuidv4_pattern%$"
          }
        }
      }
      """

  @issue-9928
  Scenario: open text file using open-with-web without app name in url query (MIME type not registered in app-registry)
    Given user "Alice" has uploaded file "filesForUpload/lorem.txt" to "lorem.txt"
    And we save it into "FILEID"
    When user "Alice" sends HTTP method "POST" to URL "/app/open-with-web?file_id=<<FILEID>>"
    Then the HTTP status code should be "500"
    And the JSON data of the response should match
      """
      {
        "type": "object",
        "required": [
          "code",
          "message"
        ],
        "properties": {
          "code": {
            "const": "SERVER_ERROR"
          },
          "message": {
            "const": "Error contacting the requested application, please use a different one or try again later"
          }
        }
      }
      """
