@issue-1276 @issue-1277 @issue-1269

Feature: changing a public link share
  As a user
  I want to set the permissions of a public link share
  So that people who have the public link only have the designated authorization

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
    And user "Alice" has created folder "PARENT"
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "PARENT/parent.txt"


  Scenario Outline: public can or cannot delete file through publicly shared link depending on having delete permissions using the public WebDAV API
    Given user "Alice" has created a public link share with settings
      | path        | /PARENT       |
      | permissions | <permissions> |
    When the public deletes file "parent.txt" from the last public link share using the new public WebDAV API
    Then the HTTP status code should be "<http-status-code>"
    And as "Alice" file "PARENT/parent.txt" <should-or-not> exist

    Examples:
      | permissions               | http-status-code | should-or-not |
      | read                      | 403              | should        |
      | read,create               | 403              | should        |
      | create                    | 403              | should        |
      | read,update,create,delete | 204              | should not    |


  Scenario: public link share permissions work correctly for renaming and share permissions read,update,create,delete using the public WebDAV API
    Given user "Alice" has created a public link share with settings
      | path        | /PARENT                   |
      | permissions | read,update,create,delete |
    When the public renames file "parent.txt" to "newparent.txt" from the last public link share using the new public WebDAV API
    Then the HTTP status code should be "201"
    And as "Alice" file "/PARENT/parent.txt" should not exist
    And as "Alice" file "/PARENT/newparent.txt" should exist


  Scenario: public link share permissions work correctly for upload with share permissions read,update,create,delete with the public WebDAV API
    Given user "Alice" has created a public link share with settings
      | path        | /PARENT                   |
      | permissions | read,update,create,delete |
    When the public uploads file "lorem.txt" with content "test" using the new public WebDAV API
    Then the HTTP status code should be "201"
    And the content of file "PARENT/lorem.txt" for user "Alice" should be "test"


  Scenario: public cannot delete file through publicly shared link with password using an invalid password with public WebDAV API
    Given user "Alice" has created a public link share with settings
      | path        | /PARENT   |
      | permissions | change    |
      | password    | newpasswd |
    When the public deletes file "parent.txt" from the last public link share using the password "invalid" and new public WebDAV API
    Then the HTTP status code should be "401"
    And as "Alice" file "PARENT/parent.txt" should exist


  Scenario: public can delete file through publicly shared link with password using the valid password with the public WebDAV API
    Given user "Alice" has created a public link share with settings
      | path        | /PARENT   |
      | permissions | change    |
      | password    | newpasswd |
    When the public deletes file "parent.txt" from the last public link share using the password "newpasswd" and new public WebDAV API
    Then the HTTP status code should be "204"
    And as "Alice" file "PARENT/parent.txt" should not exist


  Scenario: public tries to rename a file in a password protected share using an invalid password with the public WebDAV API
    Given user "Alice" has created a public link share with settings
      | path        | /PARENT   |
      | permissions | change    |
      | password    | newpasswd |
    When the public renames file "parent.txt" to "newparent.txt" from the last public link share using the password "invalid" and new public WebDAV API
    Then the HTTP status code should be "401"
    And as "Alice" file "/PARENT/newparent.txt" should not exist
    And as "Alice" file "/PARENT/parent.txt" should exist


  Scenario: public tries to rename a file in a password protected share using the valid password with the public WebDAV API
    Given user "Alice" has created a public link share with settings
      | path        | /PARENT   |
      | permissions | change    |
      | password    | newpasswd |
    When the public renames file "parent.txt" to "newparent.txt" from the last public link share using the password "newpasswd" and new public WebDAV API
    Then the HTTP status code should be "201"
    And as "Alice" file "/PARENT/newparent.txt" should exist
    And as "Alice" file "/PARENT/parent.txt" should not exist


  Scenario: public tries to upload to a password protected public share using an invalid password with the public WebDAV API
    Given user "Alice" has created a public link share with settings
      | path        | /PARENT   |
      | permissions | change    |
      | password    | newpasswd |
    When the public uploads file "lorem.txt" with password "invalid" and content "test" using the new public WebDAV API
    Then the HTTP status code should be "401"
    And as "Alice" file "/PARENT/lorem.txt" should not exist


  Scenario: public tries to upload to a password protected public share using the valid password with the public WebDAV API
    Given user "Alice" has created a public link share with settings
      | path        | /PARENT   |
      | permissions | change    |
      | password    | newpasswd |
    When the public uploads file "lorem.txt" with password "newpasswd" and content "test" using the new public WebDAV API
    Then the HTTP status code should be "201"
    And as "Alice" file "/PARENT/lorem.txt" should exist


  Scenario: public cannot rename a file in upload-write-only public link share with the public WebDAV API
    Given user "Alice" has created a public link share with settings
      | path        | /PARENT         |
      | permissions | uploadwriteonly |
    When the public renames file "parent.txt" to "newparent.txt" from the last public link share using the new public WebDAV API
    Then the HTTP status code should be "403"
    And as "Alice" file "/PARENT/parent.txt" should exist
    And as "Alice" file "/PARENT/newparent.txt" should not exist


  Scenario: public cannot delete a file in upload-write-only public link share with the public WebDAV API
    Given user "Alice" has created a public link share with settings
      | path        | /PARENT         |
      | permissions | uploadwriteonly |
    When the public deletes file "parent.txt" from the last public link share using the new public WebDAV API
    Then the HTTP status code should be "403"
    And as "Alice" file "PARENT/parent.txt" should exist
