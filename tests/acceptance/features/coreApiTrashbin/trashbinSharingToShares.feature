@skipOnReva
Feature: using trashbin together with sharing
  As a user
  I want the deletion of the resources that I shared to end up in my trashbin
  So that I can restore the resources that were accidentally deleted

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file with content "file to delete" to "/textfile0.txt"

  @smokeTest @issue-7555
  Scenario Outline: deleting a received folder doesn't move it to trashbin
    Given using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has created folder "/shared"
    And user "Alice" has moved file "/textfile0.txt" to "/shared/shared_file.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | shared   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "shared" synced
    And user "Brian" has moved folder "/Shares/shared" to "/Shares/renamed_shared"
    When user "Brian" deletes folder "/Shares/renamed_shared" using the WebDAV API
    Then the HTTP status code should be "204"
    And as "Brian" the folder with original path "/Shares/renamed_shared" should not exist in the trashbin
    Examples:
      | dav-path-version |
      | new              |
      | spaces           |

  @issue-7555
  Scenario Outline: deleting a file in a received folder moves it to trashbin of both users
    Given using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has created folder "/shared"
    And user "Alice" has moved file "/textfile0.txt" to "/shared/shared_file.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | shared   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "shared" synced
    And user "Brian" has moved file "/Shares/shared" to "/Shares/renamed_shared"
    When user "Brian" deletes file "/Shares/renamed_shared/shared_file.txt" using the WebDAV API
    Then the HTTP status code should be "204"
    And as "Brian" the file with original path "/Shares/renamed_shared/shared_file.txt" should exist in the trashbin
    And as "Alice" the file with original path "/shared/shared_file.txt" should exist in the trashbin
    Examples:
      | dav-path-version |
      | new              |
      | spaces           |


  Scenario Outline: sharee deleting a file in a group-shared folder moves it to the trashbin of sharee and sharer only
    Given using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Carol" has been created with default attributes and without skeleton files
    And group "grp1" has been created
    And user "Brian" has been added to group "grp1"
    And user "Carol" has been added to group "grp1"
    And user "Alice" has created folder "/shared"
    And user "Alice" has moved file "/textfile0.txt" to "/shared/shared_file.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | shared   |
      | space           | Personal |
      | sharee          | grp1     |
      | shareType       | group    |
      | permissionsRole | Editor   |
    And user "Brian" has a share "shared" synced
    And user "Carol" has a share "shared" synced
    When user "Brian" deletes file "/Shares/shared/shared_file.txt" using the WebDAV API
    Then the HTTP status code should be "204"
    And as "Brian" the file with original path "/Shares/shared/shared_file.txt" should exist in the trashbin
    And as "Alice" the file with original path "/shared/shared_file.txt" should exist in the trashbin
    And as "Carol" the file with original path "/Shares/shared/shared_file.txt" should not exist in the trashbin
    Examples:
      | dav-path-version |
      | new              |
      | spaces           |


  Scenario Outline: sharer deleting a file in a group-shared folder moves it to the trashbin of sharer only
    Given using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Carol" has been created with default attributes and without skeleton files
    And group "grp1" has been created
    And user "Brian" has been added to group "grp1"
    And user "Carol" has been added to group "grp1"
    And user "Alice" has created folder "/shared"
    And user "Alice" has moved file "/textfile0.txt" to "/shared/shared_file.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | shared   |
      | space           | Personal |
      | sharee          | grp1     |
      | shareType       | group    |
      | permissionsRole | Editor   |
    And user "Brian" has a share "shared" synced
    And user "Carol" has a share "shared" synced
    When user "Alice" deletes file "/shared/shared_file.txt" using the WebDAV API
    Then the HTTP status code should be "204"
    And as "Alice" the file with original path "/shared/shared_file.txt" should exist in the trashbin
    And as "Brian" the file with original path "/Shares/shared/shared_file.txt" should not exist in the trashbin
    And as "Carol" the file with original path "/Shares/shared/shared_file.txt" should not exist in the trashbin
    Examples:
      | dav-path-version |
      | new              |
      | spaces           |


  Scenario Outline: sharee deleting a folder in a group-shared folder moves it to the trashbin of sharee and sharer only
    Given using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Carol" has been created with default attributes and without skeleton files
    And group "grp1" has been created
    And user "Brian" has been added to group "grp1"
    And user "Carol" has been added to group "grp1"
    And user "Alice" has created folder "/shared"
    And user "Alice" has created folder "/shared/sub"
    And user "Alice" has moved file "/textfile0.txt" to "/shared/sub/shared_file.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | shared   |
      | space           | Personal |
      | sharee          | grp1     |
      | shareType       | group    |
      | permissionsRole | Editor   |
    And user "Brian" has a share "shared" synced
    And user "Carol" has a share "shared" synced
    When user "Brian" deletes file "/Shares/shared/sub/shared_file.txt" using the WebDAV API
    Then the HTTP status code should be "204"
    And as "Brian" the file with original path "/Shares/shared/sub/shared_file.txt" should exist in the trashbin
    And as "Alice" the file with original path "/shared/sub/shared_file.txt" should exist in the trashbin
    And as "Carol" the file with original path "/Shares/sub/shared/shared_file.txt" should not exist in the trashbin
    Examples:
      | dav-path-version |
      | new              |
      | spaces           |


  Scenario Outline: sharer deleting a folder in a group-shared folder moves it to the trashbin of sharer only
    Given using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Carol" has been created with default attributes and without skeleton files
    And group "grp1" has been created
    And user "Brian" has been added to group "grp1"
    And user "Carol" has been added to group "grp1"
    And user "Alice" has created folder "/shared"
    And user "Alice" has created folder "/shared/sub"
    And user "Alice" has moved file "/textfile0.txt" to "/shared/sub/shared_file.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | shared   |
      | space           | Personal |
      | sharee          | grp1     |
      | shareType       | group    |
      | permissionsRole | Editor   |
    And user "Brain" has a share "shared" synced
    And user "Carol" has a share "shared" synced
    When user "Alice" deletes file "/shared/sub/shared_file.txt" using the WebDAV API
    Then the HTTP status code should be "204"
    And as "Alice" the file with original path "/shared/sub/shared_file.txt" should exist in the trashbin
    And as "Brian" the file with original path "/Shares/shared/sub/shared_file.txt" should not exist in the trashbin
    And as "Carol" the file with original path "/Shares/shared/sub/shared_file.txt" should not exist in the trashbin
    Examples:
      | dav-path-version |
      | new              |
      | spaces           |

  @issue-7555
  Scenario Outline: deleting a file in a received folder when restored it comes back to the original path
    Given using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has created folder "/shared"
    And user "Alice" has moved file "/textfile0.txt" to "/shared/shared_file.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | shared   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brain" has a share "shared" synced
    And user "Brian" has moved folder "/Shares/shared" to "/Shares/renamed_shared"
    And user "Brian" has deleted file "/Shares/renamed_shared/shared_file.txt"
    When user "Brian" restores the file with original path "/Shares/renamed_shared/shared_file.txt" using the trashbin API
    Then the HTTP status code should be "201"
    And the following headers should match these regular expressions for user "Brian"
      | ETag | /^"[a-f0-9:\.]{1,32}"$/ |
    And as "Brian" the file with original path "/Shares/renamed_shared/shared_file.txt" should not exist in the trashbin
    And user "Brian" should see the following elements
      | /Shares/renamed_shared/                |
      | /Shares/renamed_shared/shared_file.txt |
    And the content of file "/Shares/renamed_shared/shared_file.txt" for user "Brian" should be "file to delete"
    Examples:
      | dav-path-version |
      | new              |
      | spaces           |


  Scenario Outline: restoring a file to a read-only folder
    Given using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Brian" has created folder "shareFolderParent"
    And user "Brian" has sent the following resource share invitation:
      | resource        | shareFolderParent |
      | space           | Personal          |
      | sharee          | Alice             |
      | shareType       | user              |
      | permissionsRole | Viewer            |
    And user "Alice" has a share "shareFolderParent" synced
    And as "Alice" folder "/Shares/shareFolderParent" should exist
    And user "Alice" has deleted file "/textfile0.txt"
    When user "Alice" restores the file with original path "/textfile0.txt" to "/Shares/shareFolderParent/textfile0.txt" using the trashbin API
    Then the HTTP status code should be "403"
    And as "Alice" the file with original path "/textfile0.txt" should exist in the trashbin
    And as "Alice" file "/Shares/shareFolderParent/textfile0.txt" should not exist
    And as "Brian" file "/shareFolderParent/textfile0.txt" should not exist
    Examples:
      | dav-path-version |
      | new              |
      | spaces           |


  Scenario Outline: restoring a file to a read-only sub-folder
    Given using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Brian" has created folder "shareFolderParent"
    And user "Brian" has created folder "shareFolderParent/shareFolderChild"
    And user "Brian" has sent the following resource share invitation:
      | resource        | shareFolderParent |
      | space           | Personal          |
      | sharee          | Alice             |
      | shareType       | user              |
      | permissionsRole | Viewer            |
    And user "Alice" has a share "shareFolderParent" synced
    And as "Alice" folder "/Shares/shareFolderParent/shareFolderChild" should exist
    And user "Alice" has deleted file "/textfile0.txt"
    When user "Alice" restores the file with original path "/textfile0.txt" to "/Shares/shareFolderParent/shareFolderChild/textfile0.txt" using the trashbin API
    Then the HTTP status code should be "403"
    And as "Alice" the file with original path "/textfile0.txt" should exist in the trashbin
    And as "Alice" file "/Shares/shareFolderParent/shareFolderChild/textfile0.txt" should not exist
    And as "Brian" file "/shareFolderParent/shareFolderChild/textfile0.txt" should not exist
    Examples:
      | dav-path-version |
      | new              |
      | spaces           |
