@api @skipOnOcV10
Feature: Upload files into a space
  As a user
  I want to be able to work with project spaces and quota

  Note - this feature is run in CI with ACCOUNTS_HASH_DIFFICULTY set to the default for production
  See https://github.com/owncloud/ocis/issues/1542 and https://github.com/owncloud/ocis/pull/839

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files
    And user "Bob" has been created with default attributes and without skeleton files
    And the administrator has given "Alice" the role "Admin" using the settings api

  Scenario: A user can create a folder in a Space via the Graph API
    Given user "Alice" has created a space "Project Ceres" of type "project" with quota "2000"
    When user "Alice" creates a folder "mainFolder" in space "Project Ceres" using the WebDav Api
    Then the HTTP status code should be "201"
    And for user "Alice" the space "Project Ceres" should contain these entries:
      | mainFolder        |

  Scenario: A user cannot create a folder in a Space if they do not have permission
    Given user "Alice" has created a space "Project Merkur" of type "project" with quota "2000"
    And user "Bob" creates a folder "forAlice" in space "Project Merkur" owned by the user "Alice" using the WebDav Api
    Then the HTTP status code should be "404"
    And for user "Alice" the space "Project Merkur" should not contain these entries:
      | forAlice        |

  Scenario: A user can create a folder and upload a file to a Space
    Given user "Alice" has created a space "Project Moon" of type "project" with quota "2000"
    When user "Alice" creates a folder "NewFolder" in space "Project Moon" using the WebDav Api
    Then the HTTP status code should be "201"
    And user "Alice" uploads a file inside space "Project Moon" with content "Test" to "test.txt" using the WebDAV API
    Then the HTTP status code should be "201"
    And for user "Alice" the space "Project Moon" should contain these entries:
      | NewFolder        |
      | test.txt         |

  Scenario: A user cannot upload a file in a Space if they do not have permission
    Given user "Alice" has created a space "Project Pluto" of type "project" with quota "2000"
    When user "Bob" uploads a file inside space "Project Pluto" owned by the user "Alice" with content "Test" to "test.txt" using the WebDAV API
    Then the HTTP status code should be "404"
    And for user "Alice" the space "Project Pluto" should not contain these entries:
      | test.txt        |

  Scenario: A user can upload a file in a Space and see the remaining quota
    When user "Alice" creates a space "Project Saturn" of type "project" with quota "2000" using the GraphApi
    And the json responded should contain a space "Project Saturn" with these key and value pairs:
      | key              | value         |
      | driveType        | project       |
      | id               | %space_id%    |
      | name             | Project Saturn|
      | quota@@@total    | 2000          |
    And user "Alice" uploads a file inside space "Project Saturn" with content "Test" to "test.txt" using the WebDAV API
    Then the HTTP status code should be "201"
    And for user "Alice" the space "Project Saturn" should contain these entries:
      | test.txt         |
    When user "Alice" lists all available spaces via the GraphApi
    Then the json responded should contain a space "Project Saturn" with these key and value pairs:
      | key              | value         |
      | driveType        | project       |
      | id               | %space_id%    |
      | name             | Project Saturn|
      | quota@@@state    | normal        |
      | quota@@@total    | 2000          |
      | quota@@@remaining| 1996          |
      | quota@@@used     | 4             |

  Scenario: A user with role editor can create a folder in shared Space via the Graph API
    Given user "Alice" has created a space "Editor can create folder" of type "project" with quota "2000"
    And user "Alice" has shared a space "Editor can create folder" to user "Bob" with role "editor"
    When user "Bob" creates a folder "mainFolder" in space "Editor can create folder" using the WebDav Api
    Then the HTTP status code should be "201"
    And for user "Bob" the space "Editor can create folder" should contain these entries:
      | mainFolder        |
    And for user "Alice" the space "Editor can create folder" should contain these entries:
      | mainFolder        |


  Scenario: A user with role viewer cannot create a folder in shared Space via the Graph API
    Given user "Alice" has created a space "Viewer cannot create folder" of type "project" with quota "2000"
    And user "Alice" has shared a space "Viewer cannot create folder" to user "Bob" with role "viewer"
    When user "Bob" creates a folder "mainFolder" in space "Viewer cannot create folder" using the WebDav Api
    Then the HTTP status code should be "403"
    And for user "Bob" the space "Viewer cannot create folder" should not contain these entries:
      | mainFolder        |
    And for user "Alice" the space "Viewer cannot create folder" should not contain these entries:
      | mainFolder        |


  Scenario: A user with role editor can upload a file in shared Space via the Graph API
    Given user "Alice" has created a space "Editor can upload file" of type "project" with quota "20"
    And user "Alice" has shared a space "Editor can upload file" to user "Bob" with role "editor"
    When user "Bob" uploads a file inside space "Editor can upload file" with content "Test" to "test.txt" using the WebDAV API
    Then the HTTP status code should be "201"
    And for user "Bob" the space "Editor can upload file" should contain these entries:
      | test.txt         |
    And for user "Alice" the space "Editor can upload file" should contain these entries:
      | test.txt         |
    When user "Bob" lists all available spaces via the GraphApi
    Then the json responded should contain a space "Editor can upload file" with these key and value pairs:
      | key              | value                 |
      | name             | Editor can upload file |
      | quota@@@used     | 4                     |


  Scenario: A user with role viewer cannot upload a file in shared Space via the Graph API
    Given user "Alice" has created a space "Viewer cannot upload file" of type "project" with quota "20"
    And user "Alice" has shared a space "Viewer cannot upload file" to user "Bob" with role "viewer"
    When user "Bob" uploads a file inside space "Viewer cannot upload file" with content "Test" to "test.txt" using the WebDAV API
    Then the HTTP status code should be "403"
    And for user "Bob" the space "Viewer cannot upload file" should not contain these entries:
      | test.txt         |
    And for user "Alice" the space "Viewer cannot upload file" should not contain these entries:
      | test.txt         |
      