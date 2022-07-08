@api @skipOnOcV10
Feature: Preview file in project space
  As a user, I want to be able to download different files for the preview

  Note - this feature is run in CI with ACCOUNTS_HASH_DIFFICULTY set to the default for production
  See https://github.com/owncloud/ocis/issues/1542 and https://github.com/owncloud/ocis/pull/839

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
    And the administrator has given "Alice" the role "Space Admin" using the settings api
    And user "Alice" has created a space "previews of the files" with the default quota using the GraphApi


  Scenario Outline: An user can preview created txt files in the project space
    Given user "Alice" has uploaded a file inside space "previews of the files" with content "test" to "<entity>"
    When user "Alice" downloads the preview of "<entity>" of the space "previews of the files" with width "<width>" and height "<height>" using the WebDAV API
    Then the HTTP status code should be "200"
    Examples:
      | entity                | width | height |
      | /file.txt             | 36    | 36     |
      | /name with spaces.txt | 1200  | 1200   |


  Scenario Outline: An user can preview image files in the project space
    Given user "Alice" has uploaded a file "<entity>" via TUS inside of the space "previews of the files" using the WebDAV API
    When user "Alice" downloads the preview of "<entity>" of the space "previews of the files" with width "<width>" and height "<height>" using the WebDAV API
    Then the HTTP status code should be "200"
    Examples:
      | entity         | width | height |
      | qrcode.png     | 36    | 36     |
      | qrcode.png     | 1200  | 1200   |
      | qrcode.png     | 1920  | 1920   |
      | testavatar.jpg | 36    | 36     |
      | testavatar.jpg | 1200  | 1200   |
      | testavatar.jpg | 1920  | 1920   |
      | example.gif    | 36    | 36     |
      | example.gif    | 1200  | 1200   |
      | example.gif    | 1280  | 1280   |
