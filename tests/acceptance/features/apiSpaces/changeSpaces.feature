@api @skipOnOcV10
Feature: Change data of space
  As an user-owner of the space or an user who is a participant with the role of manager or editor
  I want to be able to change the data of a created space (increase the quota, change name, etc.)
  
  system role:
  | admin        |
  | Spacemanager |
  | user         |

  sharing role:
  | manager |
  | editor  |
  | viewer  |

  Cases:
  | user                                       | change spaceName | change description | change quota | set readme file | set space image |
  | owner with the SpaceManager role           |        v         |          v         |      v       |       v        |        v        |
  | participant-user with manager role         |        v         |          v         |      v       |       v        |        v        |
  | participant-user with editor role          |        x         |          x         |      x       |       v        |        v        |
  | participant-user with viewer role          |        x         |          x         |      x       |       x        |        x        |
  | participant-SpaceManager with manager role |        v         |          v         |      v       |       v        |        v        |
  | participant-SpaceManager with editor role  |        x         |          x         |      x       |       v        |        v        |
  | participant-SpaceManager with viewer role  |        x         |          x         |      x       |       x        |        x        |

  Note - this feature is run in CI with ACCOUNTS_HASH_DIFFICULTY set to the default for production
  See https://github.com/owncloud/ocis/issues/1542 and https://github.com/owncloud/ocis/pull/839

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
      | Bob      |
    And the administrator has given "Alice" the role "Spacemanager" using the settings api
    And the administrator has given "Brian" the role "Spacemanager" using the settings api

  
  Scenario: an owner of space changes the name of a Space via the Graph API
    Given user "Alice" has created a space "Project Jupiter" of type "project" with quota "20"
    When user "Alice" changes the name of the "Project Jupiter" space to "Project Death Star"
    Then the HTTP status code should be "200"
    When user "Alice" lists all available spaces via the GraphApi
    Then the json responded should contain a space "Project Death Star" with these key and value pairs:
      | key              | value                            |
      | driveType        | project                          |
      | name             | Project Death Star               |
      | quota@@@total    | 20                               |
      | root@@@webDavUrl | %base_url%/dav/spaces/%space_id% |


  Scenario Outline: <participant> with <systemRole> role changes the name of a Space via the Graph API
    Given user "Alice" has created a space "Project Jupiter" of type "project" with quota "20"
    And user "Alice" has shared a space "Project Jupiter" to user "<participant>" with role "<sharingRole>"
    When user "<participant>" changes the name of the "Project Jupiter" space to "new space name"
    Then the HTTP status code should be "<statusCode>"
    When user "<participant>" lists all available spaces via the GraphApi
    Then the json responded should contain a space "<resultingName>" with these key and value pairs:
      | key              | value                            |
      | name             | <resultingName>                  |
      | quota@@@total    | 20                               |
      | root@@@webDavUrl | %base_url%/dav/spaces/%space_id% |
    Examples:
      | participant  | systemRole   | sharingRole | resultingName   | statusCode |
      | Brian        | Spacemanager | manager     | new space name  |    200     |
      | Brian        | Spacemanager | editor      | Project Jupiter |    403     |
      | Brian        | Spacemanager | viewer      | Project Jupiter |    403     |
      | Bob          | user         | manager     | new space name  |    200     |
      | Bob          | user         | editor      | Project Jupiter |    403     |
      | Bob          | user         | viewer      | Project Jupiter |    403     |


  Scenario: an owner of the space changes the description of a Space via the Graph API
    Given user "Alice" has created a space "Project Jupiter" of type "project" with quota "20"
    When user "Alice" changes the description of the "Project Jupiter" space to "The Death Star is a fictional mobile space station"
    Then the HTTP status code should be "200"
    When user "Alice" lists all available spaces via the GraphApi
    Then the json responded should contain a space "Project Jupiter" with these key and value pairs:
      | key              | value                            |
      | driveType        | project                          |
      | name             | Project Jupiter                  |
      | description      | The Death Star is a fictional mobile space station |
      | quota@@@total    | 20                               |
      | root@@@webDavUrl | %base_url%/dav/spaces/%space_id% |
    

  Scenario Outline: <participant> with <systemRole> role changes the description of a Space via the Graph API
    Given user "Alice" has created a space "Project Jupiter" of type "project" with quota "20"
    And user "Alice" has shared a space "Project Jupiter" to user "<participant>" with role "<sharingRole>"
    And user "Alice" has changed the description of the "Project Jupiter" space to "subtitle"
    When user "<participant>" changes the description of the "Project Jupiter" space to "new subtitle"
    Then the HTTP status code should be "<statusCode>"
    When user "<participant>" lists all available spaces via the GraphApi
    Then the json responded should contain a space "Project Jupiter" with these key and value pairs:
      | key              | value                            |
      | driveType        | project                          |
      | name             | Project Jupiter                  |
      | description      | <descriptionName>                |
      | quota@@@total    | 20                               |
      | root@@@webDavUrl | %base_url%/dav/spaces/%space_id% |
    Examples:
      | participant  | systemRole   | sharingRole | descriptionName | statusCode |
      | Brian        | Spacemanager | manager     | new subtitle    |    200     |
      | Brian        | Spacemanager | editor      | subtitle        |    403     |
      | Brian        | Spacemanager | viewer      | subtitle        |    403     |
      | Bob          | user         | manager     | new subtitle    |    200     | 
      | Bob          | user         | editor      | subtitle        |    403     |
      | Bob          | user         | viewer      | subtitle        |    403     |


  Scenario: an owner of the space increases the quota of a Space via the Graph API
    Given user "Alice" has created a space "Project Earth" of type "project" with quota "20"
    When user "Alice" changes the quota of the "Project Earth" space to "100"
    Then the HTTP status code should be "200"
    When user "Alice" lists all available spaces via the GraphApi
    Then the json responded should contain a space "Project Earth" with these key and value pairs:
      | key              | value         |
      | name             | Project Earth |
      | quota@@@total    | 100           |
  

  Scenario Outline: <participant> with <systemRole> role increases the quota of a Space via the Graph API
    Given user "Alice" has created a space "Project Earth" of type "project" with quota "20"
    And user "Alice" has shared a space "Project Earth" to user "<participant>" with role "<sharingRole>"
    When user "<participant>" changes the quota of the "Project Earth" space to "100"
    Then the HTTP status code should be "<statusCode>"
    When user "<participant>" lists all available spaces via the GraphApi
    Then the json responded should contain a space "Project Earth" with these key and value pairs:
      | key              | value         |
      | name             | Project Earth |
      | quota@@@total    | <quotaTotal>  |
    Examples:
      | participant  | systemRole   | sharingRole | quotaTotal | statusCode |
      | Brian        | Spacemanager | manager     |     100    |    200     |
      | Brian        | Spacemanager | editor      |     20     |    403     |
      | Brian        | Spacemanager | viewer      |     20     |    403     |
      | Bob          | user         | manager     |     100    |    200     |
      | Bob          | user         | editor      |     20     |    401     |
      | Bob          | user         | viewer      |     20     |    401     |


  Scenario: an owner of the space set readme file as description of the space via the Graph API
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
      | special@@@0@@@file@@@mimeType       | text/markdown       |
      | special@@@0@@@id                   | %file_id%            |
      | special@@@0@@@eTag                 | %eTag%              |


  Scenario: an owner of the space adds new file and set as description of the space via the Graph API
    Given user "Alice" has created a space "add special section" with the default quota using the GraphApi
    And user "Alice" has created a folder ".space" in space "add special section"
    And user "Alice" has uploaded a file inside space "add special section" with content "space description" to ".space/readme.md"
    And user "Alice" has set the file ".space/readme.md" as a description in a special section of the "add special section" space
    When user "Alice" uploads a file inside space "add special section" with content "new space description" to ".space/newReadme.md" using the WebDAV API
    And user "Alice" sets the file ".space/newReadme.md" as a description in a special section of the "add special section" space
    Then the HTTP status code should be "200"
    When user "Alice" lists all available spaces via the GraphApi
    Then the json responded should contain a space "add special section" owned by "Alice" with description file ".space/newReadme.md" with these key and value pairs:
      | key                                | value               |
      | name                               | add special section |
      | special@@@0@@@size                 | 21                  |
      | special@@@0@@@name                 | newReadme.md        |
      | special@@@0@@@specialFolder@@@name | readme              |
      | special@@@0@@@file@@@mimeType       | text/markdown       |
      | special@@@0@@@id                   | %file_id%            |
      | special@@@0@@@eTag                 | %eTag%              |


  Scenario Outline: <participant> with <systemRole> set new readme file as description of the space via the Graph API
    Given user "Alice" has created a space "add special section" with the default quota using the GraphApi
    And user "Alice" has shared a space "add special section" to user "<participant>" with role "<sharingRole>"
    And user "Alice" has created a folder ".space" in space "add special section"
    And user "Alice" has uploaded a file inside space "add special section" with content "space description" to ".space/readme.md"
    And user "Alice" has set the file ".space/readme.md" as a description in a special section of the "add special section" space
    When user "Alice" has uploaded a file inside space "add special section" with content "new space description" to ".space/newReadme.md"
    And user "<participant>" sets the file ".space/newReadme.md" as a description in a special section of the "add special section" space
    Then the HTTP status code should be "<statusCode>"
    When user "<participant>" lists all available spaces via the GraphApi
    Then the json responded should contain a space "add special section" owned by "Alice" with description file ".space/<fileName>" with these key and value pairs:
      | key                                | value               |
      | name                               | add special section |
      | special@@@0@@@size                 | <fileSize>           |
      | special@@@0@@@name                 | <fileName>           |
      | special@@@0@@@specialFolder@@@name | readme              |
      | special@@@0@@@file@@@mimeType       | text/markdown       |
      | special@@@0@@@id                   | %file_id%            |
      | special@@@0@@@eTag                 | %eTag%              |
    Examples:
      | participant  | systemRole   | sharingRole | fileName      | fileSize | statusCode |
      | Brian        | Spacemanager | manager     | newReadme.md |    21   |     200    |
      | Brian        | Spacemanager | editor      | newReadme.md |    21   |     200    |
      | Brian        | Spacemanager | viewer      | readme.md    |    17   |     403    |
      | Bob          | user         | manager     | newReadme.md |    21   |     200    |
      | Bob          | user         | editor      | newReadme.md |    21   |     200    |
      | Bob          | user         | viewer      | readme.md    |    17   |     403    |


  Scenario Outline: owner set image file as space image of the space via the Graph API
    Given user "Alice" has created a space "add space image" with the default quota using the GraphApi
    And user "Alice" has created a folder ".space" in space "add space image"
    And user "Alice" has uploaded a file inside space "add space image" with content "" to ".space/<fileName>"
    When user "Alice" sets the file ".space/<fileName>" as a space image in a special section of the "add space image" space
    Then the HTTP status code should be "200"
    When user "Alice" lists all available spaces via the GraphApi
    Then the json responded should contain a space "add space image" owned by "Alice" with description file ".space/<fileName>" with these key and value pairs:
      | key                                | value               |
      | name                               | add space image     |
      | special@@@0@@@size                 | 0                   |
      | special@@@0@@@name                 | <fileName>           |
      | special@@@0@@@specialFolder@@@name | image               |
      | special@@@0@@@file@@@mimeType       | <mimeType>          |
      | special@@@0@@@id                   | %file_id%            |
      | special@@@0@@@eTag                 | %eTag%              |
    Examples:
      | fileName         | nameInResponse  | mimeType   |
      | spaceImage.jpeg | spaceImage.jpeg | image/jpeg |
      | spaceImage.png  | spaceImage.png  | image/png  |
      | spaceImage.gif  | spaceImage.gif  | image/gif  |


  Scenario Outline: <participant> with <systemRole> set image file as space image of the space via the Graph API
    Given user "Alice" has created a space "add space image" with the default quota using the GraphApi
    And user "Alice" has shared a space "add space image" to user "<participant>" with role "<sharingRole>"
    And user "Alice" has created a folder ".space" in space "add space image"
    And user "Alice" has uploaded a file inside space "add space image" with content "" to ".space/spaceImage.jpeg"
    And user "Alice" has set the file ".space/spaceImage.jpeg" as a space image in a special section of the "add space image" space
    And user "Alice" has uploaded a file inside space "add space image" with content "" to ".space/newImage.jpeg"
    When user "<participant>" sets the file ".space/newImage.jpeg" as a space image in a special section of the "add space image" space
    Then the HTTP status code should be "<statusCode>"
    When user "<participant>" lists all available spaces via the GraphApi
    Then the json responded should contain a space "add space image" owned by "Alice" with description file ".space/<fileName>" with these key and value pairs:
      | key                                | value               |
      | name                               | add space image     |
      | special@@@0@@@size                 | 0                   |
      | special@@@0@@@name                 | <fileName>           |
      | special@@@0@@@specialFolder@@@name | image               |
      | special@@@0@@@file@@@mimeType       | image/jpeg          |
      | special@@@0@@@id                   | %file_id%            |
      | special@@@0@@@eTag                 | %eTag%              |
    Examples:
      | participant  | systemRole   | sharingRole | statusCode | fileName         |
      | Brian        | Spacemanager | manager     |     200    | newImage.jpeg   |
      | Brian        | Spacemanager | editor      |     200    | newImage.jpeg   |
      | Brian        | Spacemanager | viewer      |     403    | spaceImage.jpeg | 
      | Bob          | user         | manager     |     200    | newImage.jpeg   |
      | Bob          | user         | editor      |     200    | newImage.jpeg   |
      | Bob          | user         | viewer      |     403    | spaceImage.jpeg |
