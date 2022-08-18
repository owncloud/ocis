@api @skipOnOcV10
Feature: move (rename) file
  As a user
  I want to be able to move and rename files
  So that I can manage my file system

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |

  Scenario Outline: Moving a file within same space project with role manager and editor
    Given the administrator has given "Brian" the role "Space Admin" using the settings api
    And user "Brian" has created a space "Project" with the default quota using the GraphApi
    And user "Brian" has created a folder "newfolder" in space "Project"
    And user "Brian" has uploaded a file inside space "Project" with content "some content" to "insideSpace.txt"
    And user "Brian" has shared a space "Project" to user "Alice" with role "<role>"
    When user "Alice" moves file "insideSpace.txt" to "newfolder/insideSpace.txt" in space "Project" using the WebDAV API
    Then the HTTP status code should be "201"
    And for user "Alice" the space "Project" should contain these entries:
      | newfolder/insideSpace.txt       |
    And for user "Alice" the space "Project" should not contain these entries:
      | insideSpace.txt       |
    Examples:
      | role    |
      | manager |
      | editor  |


  Scenario: Moving a file within same space project with role viewer
    Given the administrator has given "Brian" the role "Space Admin" using the settings api
    And user "Brian" has created a space "Project" with the default quota using the GraphApi
    And user "Brian" has created a folder "newfolder" in space "Project"
    And user "Brian" has uploaded a file inside space "Project" with content "some content" to "insideSpace.txt"
    And user "Brian" has shared a space "Project" to user "Alice" with role "viewer"
    When user "Alice" moves file "insideSpace.txt" to "newfolder/insideSpace.txt" in space "Project" using the WebDAV API
    Then the HTTP status code should be "403"
    And for user "Alice" the space "Project" should not contain these entries:
      | newfolder/insideSpace.txt       |
    And for user "Alice" the space "Project" should contain these entries:
      | insideSpace.txt       |


  Scenario Outline: User moves a file from a space project with different a role to a space project with different role
    Given the administrator has given "Brian" the role "Space Admin" using the settings api
    And user "Brian" has created a space "Project1" with the default quota using the GraphApi
    And user "Brian" has created a space "Project2" with the default quota using the GraphApi
    And user "Brian" has uploaded a file inside space "Project1" with content "Project1 content" to "project1.txt"
    And user "Brian" has shared a space "Project2" to user "Alice" with role "<to_role>"
    And user "Brian" has shared a space "Project1" to user "Alice" with role "<from_role>"
    When user "Alice" moves file "project1.txt" from space "Project1" to "project1.txt" inside space "Project2" using the WebDAV API
    Then the HTTP status code should be "<https_status_code>"
    And for user "Alice" the space "Project1" should contain these entries:
      | project1.txt       |
    And for user "Alice" the space "Project2" should not contain these entries:
      | project1.txt       |
    Examples:
      | from_role | to_role | https_status_code |
      | manager   | manager | 502               |
      | editor    | manager | 502               |
      | manager   | editor  | 502               |
      | editor    | editor  | 502               |
      | manager   | viewer  | 403               |
      | editor    | viewer  | 403               |
      | viewer    | manager | 403               |
      | viewer    | editor  | 403               |
      | viewer    | viewer  | 403               |


  Scenario Outline: User moves a file from a space project with different role to a space personal
    Given the administrator has given "Brian" the role "Space Admin" using the settings api
    And user "Brian" has created a space "Project" with the default quota using the GraphApi
    And user "Brian" has uploaded a file inside space "Project" with content "Project content" to "project.txt"
    And user "Brian" has shared a space "Project" to user "Alice" with role "<role>"
    When user "Alice" moves file "project.txt" from space "Project" to "project.txt" inside space "Personal" using the WebDAV API
    Then the HTTP status code should be "<https_status_code>"
    And for user "Alice" the space "Project" should contain these entries:
      | project.txt       |
    And for user "Alice" the space "Personal" should not contain these entries:
      | project.txt       |
    Examples:
      | role    | https_status_code |
      | manager | 502               |
      | editor  | 502               |
      | viewer  | 403               |


  Scenario Outline: User moves a file from space project with different role to space shares jail with different role (permission)
    Given the administrator has given "Brian" the role "Space Admin" using the settings api
    And user "Brian" has created a space "Project" with the default quota using the GraphApi
    And user "Brian" has created folder "/testshare"
    And user "Brian" has uploaded a file inside space "Project" with content "Project content" to "project.txt"
    And user "Brian" has shared a space "Project" to user "Alice" with role "<role>"
    And user "Brian" has shared folder "/testshare" with user "Alice" with permissions "<permissions>"
    And user "Alice" has accepted share "/testshare" offered by user "Brian"
    When user "Alice" moves file "project.txt" from space "Project" to "/testshare/project.txt" inside space "Shares Jail" using the WebDAV API
    Then the HTTP status code should be "502"
    And for user "Alice" the space "Project" should contain these entries:
      | project.txt       |
    And for user "Alice" the space "Shares Jail" should not contain these entries:
      | /testshare/project.txt       |
    Examples:
      | role    | permissions |
      | manager | 31          |
      | editor  | 31          |
      | manager | 17          |
      | editor  | 17          |
      | viewer  | 31          |
      | viewer  | 17          |


  Scenario Outline: User moves a file from space personal to space project with different role
    Given the administrator has given "Brian" the role "Space Admin" using the settings api
    And user "Brian" has created a space "Project" with the default quota using the GraphApi
    And user "Brian" has shared a space "Project" to user "Alice" with role "<role>"
    And user "Alice" has uploaded file with content "personal space content" to "/personal.txt"
    When user "Alice" moves file "personal.txt" from space "Personal" to "personal.txt" inside space "Project" using the WebDAV API
    Then the HTTP status code should be "<https_status_code>"
    And for user "Alice" the space "Personal" should contain these entries:
      | personal.txt       |
    And for user "Alice" the space "Project" should not contain these entries:
      | personal.txt       |
    Examples:
      | role    | https_status_code |
      | manager | 502               |
      | editor  | 502               |
      | viewer  | 403               |


  Scenario Outline: User moves a file from space personal to space shares jail with different role (permission)
    Given user "Brian" has created folder "/testshare"
    And user "Brian" has shared folder "/testshare" with user "Alice" with permissions "<permissions>"
    And user "Alice" has accepted share "/testshare" offered by user "Brian"
    And user "Alice" has uploaded file with content "personal content" to "personal.txt"
    When user "Alice" moves file "personal.txt" from space "Personal" to "/testshare/personal.txt" inside space "Shares Jail" using the WebDAV API
    Then the HTTP status code should be "502"
    And for user "Alice" the space "Personal" should contain these entries:
      | personal.txt       |
    And for user "Alice" the space "Shares Jail" should not contain these entries:
      | /testshare/personal.txt       |
    Examples:
      | permissions |
      | 31          |
      | 17          |


  Scenario Outline: User moves a file from space shares jail with different role (permissions) to space personal
    Given user "Brian" has created folder "/testshare"
    And user "Brian" has uploaded file with content "testshare content" to "/testshare/testshare.txt"
    And user "Brian" has shared folder "/testshare" with user "Alice" with permissions "<permissions>"
    And user "Alice" has accepted share "/testshare" offered by user "Brian"
    When user "Alice" moves file "/testshare/testshare.txt" from space "Shares Jail" to "testshare.txt" inside space "Personal" using the WebDAV API
    Then the HTTP status code should be "502"
    And for user "Alice" the space "Personal" should not contain these entries:
      | testshare.txt |
    And for user "Alice" folder "testshare" of the space "Shares Jail" should contain these entries:
      | testshare.txt |
    Examples:
      | permissions |
      | 31          |
      | 17          |


  Scenario Outline: User moves a file from space shares jail with different role (permissions) to space project with different role
    Given the administrator has given "Brian" the role "Space Admin" using the settings api
    And user "Brian" has created a space "Project" with the default quota using the GraphApi
    And user "Brian" has shared a space "Project" to user "Alice" with role "<role>"
    And user "Brian" has created folder "/testshare"
    And user "Brian" has uploaded file with content "testshare content" to "/testshare/testshare.txt"
    And user "Brian" has shared folder "/testshare" with user "Alice" with permissions "<permissions>"
    And user "Alice" has accepted share "/testshare" offered by user "Brian"
    When user "Alice" moves file "/testshare/testshare.txt" from space "Shares Jail" to "testshare.txt" inside space "Project" using the WebDAV API
    Then the HTTP status code should be "502"
    And for user "Alice" the space "Project" should not contain these entries:
      | /testshare.txt |
    And for user "Alice" folder "testshare" of the space "Shares Jail" should contain these entries:
      | testshare.txt |
    Examples:
      | role    | permissions |
      | manager | 31          |
      | editor  | 31          |
      | viewer  | 31          |
      | manager | 17          |
      | editor  | 17          |
      | viewer  | 17          |


  Scenario: User moves a file from space shares jail with role editor to space shares jail with role editor
    Given user "Brian" has created folder "/testshare1"
    And user "Brian" has created folder "/testshare2"
    And user "Brian" has uploaded file with content "testshare1 content" to "/testshare1/testshare1.txt"
    And user "Brian" has shared folder "/testshare1" with user "Alice" with permissions "31"
    And user "Brian" has shared folder "/testshare2" with user "Alice" with permissions "31"
    And user "Alice" has accepted share "/testshare1" offered by user "Brian"
    And user "Alice" has accepted share "/testshare2" offered by user "Brian"
    When user "Alice" moves file "/testshare1/testshare1.txt" from space "Shares Jail" to "/testshare2/testshare1.txt" inside space "Shares Jail" using the WebDAV API
    Then the HTTP status code should be "201"
    And for user "Alice" folder "testshare2" of the space "Shares Jail" should contain these entries:
      | testshare1.txt |
    And for user "Alice" folder "testshare1" of the space "Shares Jail" should not contain these entries:
      | testshare1.txt |
    And for user "Brian" the space "Personal" should contain these entries:
      | /testshare2/testshare1.txt |


  Scenario: User moves a file from space shares jail with role editor to space shares jail with role viewer
    Given user "Brian" has created folder "/testshare1"
    And user "Brian" has created folder "/testshare2"
    And user "Brian" has uploaded file with content "testshare1 content" to "/testshare1/testshare1.txt"
    And user "Brian" has shared folder "/testshare1" with user "Alice" with permissions "31"
    And user "Brian" has shared folder "/testshare2" with user "Alice" with permissions "17"
    And user "Alice" has accepted share "/testshare1" offered by user "Brian"
    And user "Alice" has accepted share "/testshare2" offered by user "Brian"
    When user "Alice" moves file "/testshare1/testshare1.txt" from space "Shares Jail" to "/testshare2/testshare1.txt" inside space "Shares Jail" using the WebDAV API
    Then the HTTP status code should be "403"
    And for user "Alice" the space "Shares Jail" should not contain these entries:
      | /testshare2/testshare1.txt |
    And for user "Brian" the space "Personal" should not contain these entries:
      | /testshare2/testshare1.txt |


  Scenario: User moves a file from space shares jail with role viewer to space shares jail with role editor
    Given user "Brian" has created folder "/testshare1"
    And user "Brian" has created folder "/testshare2"
    And user "Brian" has uploaded file with content "testshare1 content" to "/testshare1/testshare1.txt"
    And user "Brian" has shared folder "/testshare1" with user "Alice" with permissions "17"
    And user "Brian" has shared folder "/testshare2" with user "Alice" with permissions "31"
    And user "Alice" has accepted share "/testshare1" offered by user "Brian"
    And user "Alice" has accepted share "/testshare2" offered by user "Brian"
    When user "Alice" moves file "/testshare1/testshare1.txt" from space "Shares Jail" to "/testshare2/testshare1.txt" inside space "Shares Jail" using the WebDAV API
    Then the HTTP status code should be "403"
    And for user "Alice" the space "Shares Jail" should not contain these entries:
      | /testshare2/testshare1.txt |
    And for user "Brian" the space "Personal" should not contain these entries:
      | /testshare2/testshare1.txt |
