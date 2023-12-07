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
      | password    | %public%      |
    When the public deletes file "parent.txt" from the last public link share using the password "%public%" and new public WebDAV API
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
      | password    | %public%                  |
    When the public renames file "parent.txt" to "newparent.txt" from the last public link share using the password "%public%" and the new public WebDAV API
    Then the HTTP status code should be "201"
    And as "Alice" file "/PARENT/parent.txt" should not exist
    And as "Alice" file "/PARENT/newparent.txt" should exist


  Scenario: public link share permissions work correctly for upload with share permissions read,update,create,delete with the public WebDAV API
    Given user "Alice" has created a public link share with settings
      | path        | /PARENT                   |
      | permissions | read,update,create,delete |
      | password    | %public%                  |
    When the public uploads file "lorem.txt" with password "%public%" and content "test" using the new public WebDAV API
    Then the HTTP status code should be "201"
    And the content of file "PARENT/lorem.txt" for user "Alice" should be "test"


  Scenario: public cannot delete file through publicly shared link with password using an invalid password with public WebDAV API
    Given user "Alice" has created a public link share with settings
      | path        | /PARENT  |
      | permissions | change   |
      | password    | %public% |
    When the public deletes file "parent.txt" from the last public link share using the password "invalid" and new public WebDAV API
    Then the HTTP status code should be "401"
    And as "Alice" file "PARENT/parent.txt" should exist


  Scenario: public can delete file through publicly shared link with password using the valid password with the public WebDAV API
    Given user "Alice" has created a public link share with settings
      | path        | /PARENT  |
      | permissions | change   |
      | password    | %public% |
    When the public deletes file "parent.txt" from the last public link share using the password "%public%" and new public WebDAV API
    Then the HTTP status code should be "204"
    And as "Alice" file "PARENT/parent.txt" should not exist


  Scenario: public tries to rename a file in a password protected share using an invalid password with the public WebDAV API
    Given user "Alice" has created a public link share with settings
      | path        | /PARENT  |
      | permissions | change   |
      | password    | %public% |
    When the public renames file "parent.txt" to "newparent.txt" from the last public link share using the password "invalid" and the new public WebDAV API
    Then the HTTP status code should be "401"
    And as "Alice" file "/PARENT/newparent.txt" should not exist
    And as "Alice" file "/PARENT/parent.txt" should exist


  Scenario: public tries to rename a file in a password protected share using the valid password with the public WebDAV API
    Given user "Alice" has created a public link share with settings
      | path        | /PARENT  |
      | permissions | change   |
      | password    | %public% |
    When the public renames file "parent.txt" to "newparent.txt" from the last public link share using the password "%public%" and the new public WebDAV API
    Then the HTTP status code should be "201"
    And as "Alice" file "/PARENT/newparent.txt" should exist
    And as "Alice" file "/PARENT/parent.txt" should not exist


  Scenario: public tries to upload to a password protected public share using an invalid password with the public WebDAV API
    Given user "Alice" has created a public link share with settings
      | path        | /PARENT  |
      | permissions | change   |
      | password    | %public% |
    When the public uploads file "lorem.txt" with password "invalid" and content "test" using the new public WebDAV API
    Then the HTTP status code should be "401"
    And as "Alice" file "/PARENT/lorem.txt" should not exist


  Scenario: public tries to upload to a password protected public share using the valid password with the public WebDAV API
    Given user "Alice" has created a public link share with settings
      | path        | /PARENT  |
      | permissions | change   |
      | password    | %public% |
    When the public uploads file "lorem.txt" with password "%public%" and content "test" using the new public WebDAV API
    Then the HTTP status code should be "201"
    And as "Alice" file "/PARENT/lorem.txt" should exist


  Scenario: public cannot rename a file in upload-write-only public link share with the public WebDAV API
    Given user "Alice" has created a public link share with settings
      | path        | /PARENT         |
      | permissions | uploadwriteonly |
      | password    | %public%        |
    When the public renames file "parent.txt" to "newparent.txt" from the last public link share using the password "%public%" and the new public WebDAV API
    Then the HTTP status code should be "403"
    And as "Alice" file "/PARENT/parent.txt" should exist
    And as "Alice" file "/PARENT/newparent.txt" should not exist


  Scenario: public cannot delete a file in upload-write-only public link share with the public WebDAV API
    Given user "Alice" has created a public link share with settings
      | path        | /PARENT         |
      | permissions | uploadwriteonly |
      | password    | %public%        |
    When the public deletes file "parent.txt" from the last public link share using the password "%public%" and new public WebDAV API
    Then the HTTP status code should be "403"
    And as "Alice" file "PARENT/parent.txt" should exist


  Scenario Outline: normal user tries to remove password of a public link share (change/create permission)
    Given using OCS API version "<ocs_api_version>"
    And user "Alice" has created a public link share with settings
      | path        | /PARENT       |
      | permissions | <permissions> |
      | password    | %public%      |
    When user "Alice" updates the last public link share using the sharing API with
      | path        | /PARENT       |
      | permissions | <permissions> |
      | password    |               |
    Then the HTTP status code should be "<http_status_code>"
    And the OCS status code should be "400"
    And the OCS status message should be "missing required password"
    Examples:
      | ocs_api_version | permissions | http_status_code |
      | 1               | change      | 200              |
      | 2               | change      | 400              |
      | 1               | create      | 200              |
      | 2               | create      | 400              |

  @issue-7821
  Scenario Outline: normal user tries to remove password of a public link (update without sending permissions)
    Given using OCS API version "<ocs_api_version>"
    And user "Alice" has created a public link share with settings
      | path        | /PARENT  |
      | permissions | change   |
      | password    | %public% |
    When user "Alice" updates the last public link share using the sharing API with
      | path     | /PARENT |
      | password |         |
    Then the HTTP status code should be "<http_status_code>"
    And the OCS status code should be "104"
    And the OCS status message should be "user is not allowed to delete the password from the public link"
    Examples:
      | ocs_api_version | http_status_code |
      | 1               | 200              |
      | 2               | 403              |


  Scenario Outline: normal user removes password of a public link (invite only public link)
    Given using OCS API version "<ocs_api_version>"
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "PARENT/parent.txt"
    And user "Alice" has created a public link share with settings
      | path        | /PARENT  |
      | permissions | invite   |
      | password    | %public% |
    When user "Alice" updates the last public link share using the sharing API with
      | path        | /PARENT |
      | password    |         |
      | permissions | invite  |
    Then the HTTP status code should be "200"
    And the OCS status code should be "<ocs_status_code>"
    And the OCS status message should be "OK"
    Examples:
      | ocs_api_version | ocs_status_code |
      | 1               | 100             |
      | 2               | 200             |


  Scenario Outline: administrator removes password of a read-only public link
    Given using OCS API version "<ocs_api_version>"
    And admin has created folder "/PARENT"
    And user "admin" has uploaded file "filesForUpload/textfile.txt" to "PARENT/parent.txt"
    And user "admin" has created a public link share with settings
      | path        | /PARENT  |
      | permissions | read     |
      | password    | %public% |
    When user "admin" updates the last public link share using the sharing API with
      | path        | /PARENT |
      | permissions | read    |
      | password    |         |
    Then the HTTP status code should be "200"
    And the OCS status code should be "<ocs_status_code>"
    And the public should be able to download file "/parent.txt" from inside the last public link shared folder using the new public WebDAV API with password ""
    Examples:
      | ocs_api_version | ocs_status_code |
      | 1               | 100             |
      | 2               | 200             |


  Scenario Outline: administrator tries to remove password of a public link share (change/create permission)
    Given using OCS API version "<ocs_api_version>"
    And admin has created folder "/PARENT"
    And user "admin" has created a public link share with settings
      | path        | /PARENT       |
      | permissions | <permissions> |
      | password    | %public%      |
    When user "admin" updates the last public link share using the sharing API with
      | path        | /PARENT       |
      | permissions | <permissions> |
      | password    |               |
    Then the HTTP status code should be "<http_status_code>"
    And the OCS status code should be "400"
    And the OCS status message should be "missing required password"
    Examples:
      | ocs_api_version | permissions | http_status_code |
      | 1               | change      | 200              |
      | 2               | change      | 400              |
      | 1               | create      | 200              |
      | 2               | create      | 400              |
