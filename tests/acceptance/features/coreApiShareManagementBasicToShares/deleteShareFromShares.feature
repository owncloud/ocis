@skipOnReva @issue-1328 @issue-1289
Feature: sharing
  As a user
  I want to delete shares
  So that I don't have redundant shares

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
    And user "Alice" has uploaded file with content "ownCloud test text file 0" to "/textfile0.txt"

  @issue-7555
  Scenario Outline: delete all group shares
    Given using OCS API version "<ocs-api-version>"
    And group "grp1" has been created
    And user "Brian" has been added to group "grp1"
    And user "Alice" has shared file "textfile0.txt" with group "grp1"
    And user "Brian" has moved file "/Shares/textfile0.txt" to "/Shares/anotherName.txt"
    When user "Alice" deletes the last share using the sharing API
    Then the OCS status code should be "<ocs-status-code>"
    And the HTTP status code should be "200"
    And user "Brian" should not see the share id of the last share
    And as "Brian" file "/Shares/textfile0.txt" should not exist
    And as "Brian" file "/Shares/anotherName.txt" should not exist
    Examples:
      | ocs-api-version | ocs-status-code |
      | 1               | 100             |
      | 2               | 200             |

  @smokeTest
  Scenario Outline: delete a share
    Given using OCS API version "<ocs-api-version>"
    And user "Alice" has shared file "textfile0.txt" with user "Brian"
    When user "Alice" deletes the last share using the sharing API
    Then the OCS status code should be "<ocs-status-code>"
    And the HTTP status code should be "200"
    And the last share id should not be included in the response
    And as "Brian" file "/Shares/textfile0.txt" should not exist
    Examples:
      | ocs-api-version | ocs-status-code |
      | 1               | 100             |
      | 2               | 200             |


  Scenario: orphaned shares
    Given using OCS API version "1"
    And user "Alice" has created folder "/common"
    And user "Alice" has created folder "/common/sub"
    And user "Alice" has shared folder "/common/sub" with user "Brian"
    When user "Alice" deletes folder "/common" using the WebDAV API
    Then the HTTP status code should be "204"
    And as "Brian" folder "/Shares/sub" should not exist
    And as "Brian" folder "/sub" should not exist

  @smokeTest
  Scenario: deleting a file out of a share as recipient creates a backup for the owner
    Given using OCS API version "1"
    And user "Alice" has created folder "/shared"
    And user "Alice" has moved file "/textfile0.txt" to "/shared/shared_file.txt"
    And user "Alice" has shared folder "/shared" with user "Brian"
    When user "Brian" deletes file "/Shares/shared/shared_file.txt" using the WebDAV API
    Then the HTTP status code should be "204"
    And as "Brian" file "/Shares/shared/shared_file.txt" should not exist
    And as "Alice" file "/shared/shared_file.txt" should not exist
    And as "Alice" file "/shared_file.txt" should exist in the trashbin
    And as "Brian" file "/shared_file.txt" should exist in the trashbin


  Scenario: deleting a folder out of a share as recipient creates a backup for the owner
    Given using OCS API version "1"
    And user "Alice" has created folder "/shared"
    And user "Alice" has created folder "/shared/sub"
    And user "Alice" has moved file "/textfile0.txt" to "/shared/sub/shared_file.txt"
    And user "Alice" has shared folder "/shared" with user "Brian"
    When user "Brian" deletes folder "/Shares/shared/sub" using the WebDAV API
    Then the HTTP status code should be "204"
    And as "Brian" folder "/Shares/shared/sub" should not exist
    And as "Alice" folder "/shared/sub" should not exist
    And as "Alice" folder "/sub" should exist in the trashbin
    And as "Alice" file "/sub/shared_file.txt" should exist in the trashbin
    And as "Brian" folder "/sub" should exist in the trashbin
    And as "Brian" file "/sub/shared_file.txt" should exist in the trashbin

  @smokeTest
  Scenario: unshare from self
    And group "grp1" has been created
    And these users have been created with default attributes and without skeleton files:
      | username |
      | Carol    |
    And user "Brian" has been added to group "grp1"
    And user "Carol" has been added to group "grp1"
    And user "Carol" has created folder "PARENT"
    And user "Carol" has uploaded file "filesForUpload/textfile.txt" to "PARENT/parent.txt"
    And user "Carol" has shared file "/PARENT/parent.txt" with group "grp1"
    And user "Carol" has stored etag of element "/PARENT"
    And user "Brian" has stored etag of element "/"
    And user "Brian" has stored etag of element "/Shares"
    When user "Brian" declines share "/Shares/parent.txt" offered by user "Carol" using the sharing API
    Then the HTTP status code should be "200"
    And the etag of element "/" of user "Brian" should have changed
    And the etag of element "/Shares" of user "Brian" should have changed
    And the etag of element "/PARENT" of user "Carol" should not have changed


  Scenario: sharee of a read-only share folder tries to delete the shared folder
    Given using OCS API version "1"
    And user "Alice" has created folder "/shared"
    And user "Alice" has moved file "/textfile0.txt" to "/shared/shared_file.txt"
    And user "Alice" has shared folder "shared" with user "Brian" with permissions "read"
    When user "Brian" deletes file "/Shares/shared/shared_file.txt" using the WebDAV API
    Then the HTTP status code should be "403"
    And as "Alice" file "/shared/shared_file.txt" should exist
    And as "Brian" file "/Shares/shared/shared_file.txt" should exist


  Scenario: sharee of a upload-only shared folder tries to delete a file in the shared folder
    Given using OCS API version "1"
    And user "Alice" has created folder "/shared"
    And user "Alice" has moved file "/textfile0.txt" to "/shared/shared_file.txt"
    And user "Alice" has shared folder "shared" with user "Brian" with permissions "create"
    When user "Brian" deletes file "/Shares/shared/shared_file.txt" using the WebDAV API
    Then the HTTP status code should be "403"
    And as "Alice" file "/shared/shared_file.txt" should exist
    # Note: for Brian, the file does not "exist" because he only has "create" permission, not "read"
    And as "Brian" file "/Shares/shared/shared_file.txt" should not exist


  Scenario: sharee of an upload-only shared folder tries to delete their file in the folder
    Given using OCS API version "1"
    And user "Alice" has created folder "/shared"
    And user "Alice" has shared folder "shared" with user "Brian" with permissions "create"
    And user "Brian" has uploaded file "filesForUpload/textfile.txt" to "/Shares/shared/textfile.txt"
    When user "Brian" deletes file "/Shares/shared/textfile.txt" using the WebDAV API
    Then the HTTP status code should be "403"
    And as "Alice" file "/shared/textfile.txt" should exist
    # Note: for Brian, the file does not "exist" because he only has "create" permission, not "read"
    And as "Brian" file "/Shares/shared/textfile.txt" should not exist


  Scenario Outline: group share recipient tries to delete the share
    Given using OCS API version "<ocs-api-version>"
    And group "grp1" has been created
    And these users have been created with default attributes and without skeleton files:
      | username |
      | Carol    |
    And user "Brian" has been added to group "grp1"
    And user "Carol" has been added to group "grp1"
    And user "Alice" has created folder "/shared"
    And user "Alice" has moved file "/textfile0.txt" to "/shared/shared_file.txt"
    And user "Alice" has shared entry "<entry-to-share>" with group "grp1"
    When user "Brian" deletes the last share of user "Alice" using the sharing API
    Then the OCS status code should be "404"
    And the HTTP status code should be "<http-status-code>"
    And as "Alice" entry "<entry-to-share>" should exist
    And as "Brian" entry "<received-entry>" should exist
    And as "Carol" entry "<received-entry>" should exist
    Examples:
      | entry-to-share          | ocs-api-version | http-status-code | received-entry          |
      | /shared/shared_file.txt | 1               | 200              | /Shares/shared_file.txt |
      | /shared/shared_file.txt | 2               | 404              | /Shares/shared_file.txt |
      | /shared                 | 1               | 200              | /Shares/shared          |
      | /shared                 | 2               | 404              | /Shares/shared          |


  Scenario Outline: individual share recipient tries to delete the share
    Given using OCS API version "<ocs-api-version>"
    And user "Alice" has created folder "/shared"
    And user "Alice" has moved file "/textfile0.txt" to "/shared/shared_file.txt"
    And user "Alice" has shared entry "<entry-to-share>" with user "Brian"
    When user "Brian" deletes the last share of user "Alice" using the sharing API
    Then the OCS status code should be "404"
    And the HTTP status code should be "<http-status-code>"
    And as "Alice" entry "<entry-to-share>" should exist
    And as "Brian" entry "<received-entry>" should exist
    Examples:
      | entry-to-share          | ocs-api-version | http-status-code | received-entry          |
      | /shared/shared_file.txt | 1               | 200              | /Shares/shared_file.txt |
      | /shared/shared_file.txt | 2               | 404              | /Shares/shared_file.txt |
      | /shared                 | 1               | 200              | /Shares/shared          |
      | /shared                 | 2               | 404              | /Shares/shared          |

  @issue-720
  Scenario Outline: request PROPFIND after sharer deletes the collaborator
    Given using OCS API version "<ocs-api-version>"
    And user "Alice" has shared file "textfile0.txt" with user "Brian"
    When user "Alice" deletes the last share using the sharing API
    Then the OCS status code should be "<ocs-status-code>"
    And the HTTP status code should be "200"
    When user "Brian" requests "/remote.php/dav/files/%username%" with "PROPFIND" using basic auth
    Then the HTTP status code should be "207"
    Examples:
      | ocs-api-version | ocs-status-code |
      | 1               | 100             |
      | 2               | 200             |

  @issue-1229
  Scenario Outline: delete a share with wrong authentication
    Given using OCS API version "<ocs-api-version>"
    And user "Alice" has shared file "textfile0.txt" with user "Brian"
    When user "Brian" tries to delete the last share of user "Alice" using the sharing API
    Then the OCS status code should be "404"
    And the HTTP status code should be "<http-status-code>"
    Examples:
      | ocs-api-version | http-status-code |
      | 1               | 200              |
      | 2               | 404              |


  Scenario Outline: unshare a shared resources
    Given using OCS API version "<ocs-api-version>"
    And user "Alice" has shared file "textfile0.txt" with user "Brian"
    When user "Alice" unshares file "textfile0.txt" shared to "Brian"
    Then the OCS status code should be "<ocs-status-code>"
    And the HTTP status code should be "200"
    And as "Brian" file "/Shares/textfile0.txt" should not exist
    Examples:
      | ocs-api-version | ocs-status-code |
      | 1               | 100             |
      | 2               | 200             |