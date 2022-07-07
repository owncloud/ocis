@api @skipOnOcV10
Feature: Change data of space
  As a user with space admin rights
  I want to be able to change the data of a created space (increase the quota, change name, etc.)

  Note - this feature is run in CI with ACCOUNTS_HASH_DIFFICULTY set to the default for production
  See https://github.com/owncloud/ocis/issues/1542 and https://github.com/owncloud/ocis/pull/839

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
      | Bob      |
    And the administrator has given "Alice" the role "Space Admin" using the settings api
    And user "Alice" has created a space "Project Jupiter" of type "project" with quota "20"
    And user "Alice" has shared a space "Project Jupiter" to user "Brian" with role "editor"
    And user "Alice" has shared a space "Project Jupiter" to user "Bob" with role "viewer"


  Scenario Outline: Only space admin user can change the name of a Space via the Graph API
    When user "<user>" changes the name of the "Project Jupiter" space to "Project Death Star"
    Then the HTTP status code should be "<code>"
    And the user "<user>" should have a space called "<expectedName>" with these key and value pairs:
      | key       | value          |
      | driveType | project        |
      | name      | <expectedName> |
    Examples:
      | user  | code | expectedName       |
      | Alice | 200  | Project Death Star |
      | Brian | 403  | Project Jupiter    |
      | Bob   | 403  | Project Jupiter    |


  Scenario: Only space admin user can change the description(subtitle) of a Space via the Graph API
    When user "Alice" changes the description of the "Project Jupiter" space to "The Death Star is a fictional mobile space station"
    Then the HTTP status code should be "200"
    And the user "Alice" should have a space called "Project Jupiter" with these key and value pairs:
      | key         | value                                              |
      | driveType   | project                                            |
      | name        | Project Jupiter                                    |
      | description | The Death Star is a fictional mobile space station |


  Scenario Outline: Viewer and editor cannot change the description(subtitle) of a Space via the Graph API
    When user "<user>" changes the description of the "Project Jupiter" space to "The Death Star is a fictional mobile space station"
    Then the HTTP status code should be "<code>"
    Examples:
      | user  | code |
      | Brian | 403  |
      | Bob   | 403  |


  Scenario Outline: An user tries to increase the quota of a Space via the Graph API
    When user "<user>" changes the quota of the "Project Jupiter" space to "100"
    Then the HTTP status code should be "<code>"
    And the user "<user>" should have a space called "Project Jupiter" with these key and value pairs:
      | key           | value                |
      | name          | Project Jupiter      |
      | quota@@@total | <expectedQuataValue> |
    Examples:
      | user  | code | expectedQuataValue |
      | Alice | 200  | 100                |
      | Brian | 401  | 20                 |
      | Bob   | 401  | 20                 |


  Scenario Outline: An space admin user set no restriction quota of a Space via the Graph API
    When user "Alice" changes the quota of the "Project Jupiter" space to "<quotaValue>"
    Then the HTTP status code should be "200"
    When user "Alice" uploads a file inside space "Project Jupiter" with content "some content" to "file.txt" using the WebDAV API
    Then the HTTP status code should be "201"
    And the user "Alice" should have a space called "Project Jupiter" with these key and value pairs:
      | key          | value           |
      | name         | Project Jupiter |
      | quota@@@used | 12              |
    Examples:
      | quotaValue |
      | 0          |
      | -1         |


  Scenario: An user space admin set readme file as description of the space via the Graph API
    Given user "Alice" has created a folder ".space" in space "Project Jupiter"
    And user "Alice" has uploaded a file inside space "Project Jupiter" with content "space description" to ".space/readme.md"
    When user "Alice" sets the file ".space/readme.md" as a description in a special section of the "Project Jupiter" space
    Then the HTTP status code should be "200"
    And the user "Alice" should have a space called "Project Jupiter" owned by "Alice" with description file ".space/readme.md" with these key and value pairs:
      | key                                | value           |
      | name                               | Project Jupiter |
      | special@@@0@@@size                 | 17              |
      | special@@@0@@@name                 | readme.md       |
      | special@@@0@@@specialFolder@@@name | readme          |
      | special@@@0@@@file@@@mimeType      | text/markdown   |
      | special@@@0@@@id                   | %file_id%       |
      | special@@@0@@@eTag                 | %eTag%          |
    And for user "Alice" folder ".space/" of the space "Project Jupiter" should contain these entries:
      | readme.md |
    And for user "Alice" the content of the file ".space/readme.md" of the space "Project Jupiter" should be "space description"


  Scenario Outline: An user member of the space changes readme file
    Given user "Alice" has created a folder ".space" in space "Project Jupiter"
    And user "Alice" has uploaded a file inside space "Project Jupiter" with content "space description" to ".space/readme.md"
    And user "Alice" has set the file ".space/readme.md" as a description in a special section of the "Project Jupiter" space
    When user "<user>" uploads a file inside space "Project Jupiter" with content "new description" to ".space/readme.md" using the WebDAV API
    Then the HTTP status code should be "<code>"
    And the user "<user>" should have a space called "Project Jupiter" owned by "Alice" with description file ".space/readme.md" with these key and value pairs:
      | key                                | value           |
      | name                               | Project Jupiter |
      | special@@@0@@@size                 | <size>          |
      | special@@@0@@@name                 | readme.md       |
      | special@@@0@@@specialFolder@@@name | readme          |
      | special@@@0@@@file@@@mimeType      | text/markdown   |
      | special@@@0@@@id                   | %file_id%       |
      | special@@@0@@@eTag                 | %eTag%          |
    And for user "<user>" folder ".space/" of the space "Project Jupiter" should contain these entries:
      | readme.md |
    And for user "<user>" the content of the file ".space/readme.md" of the space "Project Jupiter" should be "<content>"
    Examples:
      | user  | code | size | content           |
      | Brian | 204  | 15   | new description   |
      | Bob   | 403  | 17   | space description |


  Scenario Outline: An user space admin and editor set image file as space image of the space via the Graph API
    Given user "Alice" has created a folder ".space" in space "Project Jupiter"
    And user "<user>" has uploaded a file inside space "Project Jupiter" with content "" to ".space/<fileName>"
    When user "<user>" sets the file ".space/<fileName>" as a space image in a special section of the "Project Jupiter" space
    Then the HTTP status code should be "200"
    And the user "Alice" should have a space called "Project Jupiter" owned by "Alice" with description file ".space/<fileName>" with these key and value pairs:
      | key                                | value            |
      | name                               | Project Jupiter  |
      | special@@@0@@@size                 | 0                |
      | special@@@0@@@name                 | <nameInResponse> |
      | special@@@0@@@specialFolder@@@name | image            |
      | special@@@0@@@file@@@mimeType      | <mimeType>       |
      | special@@@0@@@id                   | %file_id%        |
      | special@@@0@@@eTag                 | %eTag%           |
    And for user "<user>" folder ".space/" of the space "Project Jupiter" should contain these entries:
      | <fileName> |
    Examples:
      | user  | fileName        | nameInResponse  | mimeType   |
      | Alice | spaceImage.jpeg | spaceImage.jpeg | image/jpeg |
      | Brian | spaceImage.png  | spaceImage.png  | image/png  |
      | Alice | spaceImage.gif  | spaceImage.gif  | image/gif  |


  Scenario: An user viewer cannot set image file as space image of the space via the Graph API
    Given user "Alice" has created a folder ".space" in space "Project Jupiter"
    And user "Alice" has uploaded a file inside space "Project Jupiter" with content "" to ".space/someImageFile.jpg"
    When user "Bob" sets the file ".space/someImageFile.jpg" as a space image in a special section of the "Project Jupiter" space
    Then the HTTP status code should be "403"


  Scenario Outline: An user set new readme file as description of the space via the Graph API
    Given user "Alice" has created a folder ".space" in space "Project Jupiter"
    And user "Alice" has uploaded a file inside space "Project Jupiter" with content "space description" to ".space/readme.md"
    And user "Alice" has set the file ".space/readme.md" as a description in a special section of the "Project Jupiter" space
    When user "<user>" uploads a file inside space "Project Jupiter" owned by the user "Alice" with content "new content" to ".space/readme.md" using the WebDAV API
    Then the HTTP status code should be "<code>"
    And for user "<user>" the content of the file ".space/readme.md" of the space "Project Jupiter" should be "<expectedContent>"
    And the user "<user>" should have a space called "Project Jupiter" owned by "Alice" with description file ".space/readme.md" with these key and value pairs:
      | key                                | value           |
      | name                               | Project Jupiter |
      | special@@@0@@@size                 | <expectedSize>  |
      | special@@@0@@@name                 | readme.md       |
      | special@@@0@@@specialFolder@@@name | readme          |
      | special@@@0@@@file@@@mimeType      | text/markdown   |
      | special@@@0@@@id                   | %file_id%       |
      | special@@@0@@@eTag                 | %eTag%          |
    Examples:
      | user  | code | expectedSize | expectedContent   |
      | Alice | 204  | 11           | new content       |
      | Brian | 204  | 11           | new content       |
      | Bob   | 403  | 17           | space description |


  Scenario Outline: An user set new image file as space image of the space via the Graph API
    Given user "Alice" has created a folder ".space" in space "Project Jupiter"
    And user "Alice" has uploaded a file inside space "Project Jupiter" with content "" to ".space/spaceImage.jpeg"
    And user "Alice" has set the file ".space/spaceImage.jpeg" as a space image in a special section of the "Project Jupiter" space
    When user "<user>" has uploaded a file inside space "Project Jupiter" with content "" to ".space/newSpaceImage.png"
    And user "<user>" sets the file ".space/newSpaceImage.png" as a space image in a special section of the "Project Jupiter" space
    Then the HTTP status code should be "200"
     And the user "<user>" should have a space called "Project Jupiter" owned by "Alice" with space image ".space/newSpaceImage.png" with these key and value pairs:
      | key                                | value             |
      | name                               | Project Jupiter   |
      | special@@@0@@@size                 | 0                 |
      | special@@@0@@@name                 | newSpaceImage.png |
      | special@@@0@@@specialFolder@@@name | image             |
      | special@@@0@@@file@@@mimeType      | image/png         |
      | special@@@0@@@id                   | %file_id%         |
      | special@@@0@@@eTag                 | %eTag%            |
    Examples:
      | user  |
      | Alice |
      | Brian |