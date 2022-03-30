@api @skipOnOcV10
Feature: Change data of space
  As a user with admin rights
  I want to be able to change the data of a created space (increase the quota, change name, etc.)

  Note - this feature is run in CI with ACCOUNTS_HASH_DIFFICULTY set to the default for production
  See https://github.com/owncloud/ocis/issues/1542 and https://github.com/owncloud/ocis/pull/839

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files
    And the administrator has given "Alice" the role "Admin" using the settings api

  Scenario: An admin user can change the name and description of a Space via the Graph API
    Given user "Alice" has created a space "Project Jupiter" of type "project" with quota "20"
    When user "Alice" changes the name of the "Project Jupiter" space to "Project Death Star"
    And user "Alice" changes the description of the "Project Death Star" space to "The Death Star is a fictional mobile space station"
    Then the HTTP status code should be "200"
    When user "Alice" lists all available spaces via the GraphApi
    Then the json responded should contain a space "Project Death Star" with these key and value pairs:
      | key              | value                                              |
      | driveType        | project                                            |
      | name             | Project Death Star                                 |
      | description      | The Death Star is a fictional mobile space station |
      | quota@@@total    | 20                                                 |
      | root@@@webDavUrl | %base_url%/dav/spaces/%space_id%                   |

  Scenario: An admin user can increase the quota of a Space via the Graph API
    Given user "Alice" has created a space "Project Earth" of type "project" with quota "20"
    When user "Alice" changes the quota of the "Project Earth" space to "100"
    Then the HTTP status code should be "200"
    When user "Alice" lists all available spaces via the GraphApi
    Then the json responded should contain a space "Project Earth" with these key and value pairs:
      | key           | value         |
      | name          | Project Earth |
      | quota@@@total | 100           |


  Scenario: An user set readme file as description of the space via the Graph API
    Given user "Alice" has created a space "add special section" with the default quota using the GraphApi
    And user "Alice" has created a folder ".space" in space "add special section"
    And user "Alice" has uploaded a file inside space "add special section" with content "space description" to ".space/readme.md"
    When user "Alice" sets the file ".space/readme.md" as a description in a special section of the "add special section" space
    Then the HTTP status code should be "200"
    When user "Alice" lists all available spaces via the GraphApi
    Then the json responded should contain a space "add special section" owned by "Alice" with description file ".space/readme.md" with these key and value pairs:
      | key                                | value               |
      | name                               | add special section |
      | special@@@0@@@size                 | 17                  |
      | special@@@0@@@name                 | readme.md           |
      | special@@@0@@@specialFolder@@@name | readme              |
      | special@@@0@@@file@@@mimeType      | text/markdown       |
      | special@@@0@@@id                   | %file_id%           |
      | special@@@0@@@eTag                 | %eTag%              |
    And for user "Alice" folder ".space/" of the space "add special section" should contain these entries:
      | readme.md |
    And for user "Alice" the content of the file ".space/readme.md" of the space "add special section" should be "space description"


  Scenario Outline: An user set image file as space image of the space via the Graph API
    Given user "Alice" has created a space "add special section" with the default quota using the GraphApi
    And user "Alice" has created a folder ".space" in space "add special section"
    And user "Alice" has uploaded a file inside space "add special section" with content "" to ".space/<fileName>"
    When user "Alice" sets the file ".space/<fileName>" as a space image in a special section of the "add special section" space
    Then the HTTP status code should be "200"
    When user "Alice" lists all available spaces via the GraphApi
    Then the json responded should contain a space "add special section" owned by "Alice" with description file ".space/<fileName>" with these key and value pairs:
      | key                                | value               |
      | name                               | add special section |
      | special@@@0@@@size                 | 0                   |
      | special@@@0@@@name                 | <nameInResponse>    |
      | special@@@0@@@specialFolder@@@name | image               |
      | special@@@0@@@file@@@mimeType      | <mimeType>          |
      | special@@@0@@@id                   | %file_id%           |
      | special@@@0@@@eTag                 | %eTag%              |
    And for user "Alice" folder ".space/" of the space "add special section" should contain these entries:
      | <fileName> |
    Examples:
      | fileName        | nameInResponse  | mimeType   |
      | spaceImage.jpeg | spaceImage.jpeg | image/jpeg |
      | spaceImage.png  | spaceImage.png  | image/png  |
      | spaceImage.gif  | spaceImage.gif  | image/gif  |


  Scenario: An user set new readme file as description of the space via the Graph API
    Given user "Alice" has created a space "add special section" with the default quota using the GraphApi
    And user "Alice" has created a folder ".space" in space "add special section"
    And user "Alice" has uploaded a file inside space "add special section" with content "space description" to ".space/readme.md"
    And user "Alice" has set the file ".space/readme.md" as a description in a special section of the "add special section" space
    When user "Alice" has uploaded a file inside space "add special section" with content "new space description" to ".space/newReadme.md"
    And user "Alice" sets the file ".space/newReadme.md" as a description in a special section of the "add special section" space
    Then the HTTP status code should be "200"
    And for user "Alice" the content of the file ".space/newReadme.md" of the space "add special section" should be "new space description"
    When user "Alice" lists all available spaces via the GraphApi
    Then the json responded should contain a space "add special section" owned by "Alice" with description file ".space/newReadme.md" with these key and value pairs:
      | key                                | value               |
      | name                               | add special section |
      | special@@@0@@@size                 | 21                  |
      | special@@@0@@@name                 | newReadme.md        |
      | special@@@0@@@specialFolder@@@name | readme              |
      | special@@@0@@@file@@@mimeType      | text/markdown       |
      | special@@@0@@@id                   | %file_id%           |
      | special@@@0@@@eTag                 | %eTag%              |


  Scenario: An user set new image file as space image of the space via the Graph API
    Given user "Alice" has created a space "add special section" with the default quota using the GraphApi
    And user "Alice" has created a folder ".space" in space "add special section"
    And user "Alice" has uploaded a file inside space "add special section" with content "" to ".space/spaceImage.jpeg"
    And user "Alice" has set the file ".space/spaceImage.jpeg" as a space image in a special section of the "add special section" space
    When user "Alice" has uploaded a file inside space "add special section" with content "" to ".space/newSpaceImage.png"
    And user "Alice" sets the file ".space/newSpaceImage.png" as a space image in a special section of the "add special section" space
    Then the HTTP status code should be "200"
    When user "Alice" lists all available spaces via the GraphApi
    Then the json responded should contain a space "add special section" owned by "Alice" with description file ".space/newSpaceImage.png" with these key and value pairs:
      | key                                | value               |
      | name                               | add special section |
      | special@@@0@@@size                 | 0                   |
      | special@@@0@@@name                 | newSpaceImage.png   |
      | special@@@0@@@specialFolder@@@name | image               |
      | special@@@0@@@file@@@mimeType      | image/png           |
      | special@@@0@@@id                   | %file_id%           |
      | special@@@0@@@eTag                 | %eTag%              |
      