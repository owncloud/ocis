@issue-2914 @skipOnReva
Feature: file versions remember the author of each version
  As a user
  I want to know the author of each version of a file
  So that I can better understand the version history and choose a version to restore

  Background:
    Given using OCS API version "2"
    And using new DAV path
    And user "Alice" has been created with default attributes and without skeleton files
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Carol" has been created with default attributes and without skeleton files


  Scenario: enable file versioning and check the history of changes from multiple users
    Given user "David" has been created with default attributes and without skeleton files
    And user "Alice" has created folder "/test"
    And user "Alice" has sent the following resource share invitation:
      | resource        | test     |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Alice" has sent the following resource share invitation:
      | resource        | test     |
      | space           | Personal |
      | sharee          | Carol    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Alice" has sent the following resource share invitation:
      | resource        | test     |
      | space           | Personal |
      | sharee          | David    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Alice" has uploaded file with content "uploaded content alice" to "/test/textfile0.txt"
    And user "Brian" has uploaded file with content "uploaded content brian" to "/Shares/test/textfile0.txt"
    And user "Carol" has uploaded file with content "uploaded content carol" to "/Shares/test/textfile0.txt"
    And user "David" has uploaded file with content "uploaded content david" to "/Shares/test/textfile0.txt"
    When user "Alice" gets the number of versions of file "/test/textfile0.txt"
    Then the HTTP status code should be "207"
    And the number of versions should be "3"
    And the content of version index "1" of file "/test/textfile0.txt" for user "Alice" should be "uploaded content carol"
    And the content of version index "2" of file "/test/textfile0.txt" for user "Alice" should be "uploaded content brian"
    And the content of version index "3" of file "/test/textfile0.txt" for user "Alice" should be "uploaded content alice"
    And as users "Alice" the authors of the versions of file "/test/textfile0.txt" should be:
      | index | author |
      | 1     | Carol  |
      | 2     | Brian  |
      | 3     | Alice  |
    And as users "Brian,Carol,David" the authors of the versions of file "/Shares/test/textfile0.txt" should be:
      | index | author |
      | 1     | Carol  |
      | 2     | Brian  |
      | 3     | Alice  |


  Scenario: enable file versioning and check the history of changes from multiple users for shared folder in the group
    Given user "Alice" has created folder "/test"
    And group "grp1" has been created
    And user "Alice" has been added to group "grp1"
    And user "Brian" has been added to group "grp1"
    And user "Carol" has been added to group "grp1"
    And user "Alice" has sent the following resource share invitation:
      | resource        | test     |
      | space           | Personal |
      | sharee          | grp1     |
      | shareType       | group    |
      | permissionsRole | Editor   |
    And user "Alice" has uploaded file with content "uploaded content alice" to "/test/textfile0.txt"
    And user "Brian" has uploaded file with content "uploaded content brian" to "/Shares/test/textfile0.txt"
    And user "Carol" has uploaded file with content "uploaded content carol" to "/Shares/test/textfile0.txt"
    When user "Alice" gets the number of versions of file "/test/textfile0.txt"
    Then the HTTP status code should be "207"
    And the number of versions should be "2"
    And the content of version index "1" of file "/test/textfile0.txt" for user "Alice" should be "uploaded content brian"
    And the content of version index "2" of file "/test/textfile0.txt" for user "Alice" should be "uploaded content alice"
    And as users "Alice" the authors of the versions of file "/test/textfile0.txt" should be:
      | index | author |
      | 1     | Brian  |
      | 2     | Alice  |
    And as users "Brian,Carol" the authors of the versions of file "/Shares/test/textfile0.txt" should be:
      | index | author |
      | 1     | Brian  |
      | 2     | Alice  |


  Scenario: enable file versioning and check the history of changes from multiple users for shared file in the group
    Given group "grp1" has been created
    And user "Alice" has been added to group "grp1"
    And user "Brian" has been added to group "grp1"
    And user "Carol" has been added to group "grp1"
    And user "Alice" has uploaded file with content "uploaded content alice" to "/textfile0.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | textfile0.txt |
      | space           | Personal      |
      | sharee          | grp1          |
      | shareType       | group         |
      | permissionsRole | File Editor   |
    And user "Brian" has uploaded file with content "uploaded content brian" to "/Shares/textfile0.txt"
    And user "Carol" has uploaded file with content "uploaded content carol" to "/Shares/textfile0.txt"
    When user "Alice" gets the number of versions of file "textfile0.txt"
    Then the HTTP status code should be "207"
    And the number of versions should be "2"
    And the content of version index "1" of file "/textfile0.txt" for user "Alice" should be "uploaded content brian"
    And the content of version index "2" of file "/textfile0.txt" for user "Alice" should be "uploaded content alice"
    And as users "Alice" the authors of the versions of file "/textfile0.txt" should be:
      | index | author |
      | 1     | Brian  |
      | 2     | Alice  |
    And as users "Brian,Carol" the authors of the versions of file "/Shares/textfile0.txt" should be:
      | index | author |
      | 1     | Brian  |
      | 2     | Alice  |


  Scenario: enable file versioning and check the history of changes from multiple users while moving file in/out of a subfolder
    Given user "Alice" has created folder "/test"
    And group "grp1" has been created
    And user "Alice" has been added to group "grp1"
    And user "Brian" has been added to group "grp1"
    And user "Carol" has been added to group "grp1"
    And user "Alice" has sent the following resource share invitation:
      | resource        | test     |
      | space           | Personal |
      | sharee          | grp1     |
      | shareType       | group    |
      | permissionsRole | Editor   |
    And user "Alice" has uploaded file with content "uploaded content alice" to "/test/textfile0.txt"
    And user "Brian" has uploaded file with content "uploaded content brian" to "/Shares/test/textfile0.txt"
    And user "Carol" has uploaded file with content "uploaded content carol" to "/Shares/test/textfile0.txt"
    And user "Alice" has moved file "/test/textfile0.txt" to "/textfile0.txt"
    And user "Alice" has uploaded file with content "uploaded content alice after moving file outside subfolder" to "/textfile0.txt"
    And user "Alice" has moved file "/textfile0.txt" to "/test/textfile0.txt"
    When user "Alice" gets the number of versions of file "/test/textfile0.txt"
    Then the HTTP status code should be "207"
    And the number of versions should be "3"
    And the content of version index "1" of file "/test/textfile0.txt" for user "Alice" should be "uploaded content carol"
    And the content of version index "2" of file "/test/textfile0.txt" for user "Alice" should be "uploaded content brian"
    And the content of version index "3" of file "/test/textfile0.txt" for user "Alice" should be "uploaded content alice"
    And as users "Alice" the authors of the versions of file "/test/textfile0.txt" should be:
      | index | author |
      | 1     | Carol  |
      | 2     | Brian  |
      | 3     | Alice  |
    And as users "Brian,Carol" the authors of the versions of file "/Shares/test/textfile0.txt" should be:
      | index | author |
      | 1     | Carol  |
      | 2     | Brian  |
      | 3     | Alice  |


  Scenario: enable file versioning and check the history of changes from multiple users after renaming file by sharer
    Given group "grp1" has been created
    And user "Alice" has been added to group "grp1"
    And user "Brian" has been added to group "grp1"
    And user "Carol" has been added to group "grp1"
    And user "Alice" has uploaded file with content "uploaded content alice" to "/exist.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | exist.txt   |
      | space           | Personal    |
      | sharee          | grp1        |
      | shareType       | group       |
      | permissionsRole | File Editor |
    And user "Brian" has uploaded file with content "uploaded content brian" to "/Shares/exist.txt"
    And user "Carol" has uploaded file with content "uploaded content carol" to "/Shares/exist.txt"
    And user "Alice" has moved file "/exist.txt" to "/textfile0.txt"
    When user "Alice" gets the number of versions of file "textfile0.txt"
    Then the HTTP status code should be "207"
    And the number of versions should be "2"
    And the content of version index "1" of file "/textfile0.txt" for user "Alice" should be "uploaded content brian"
    And the content of version index "2" of file "/textfile0.txt" for user "Alice" should be "uploaded content alice"
    And as user "Alice" the authors of the versions of file "/textfile0.txt" should be:
      | index | author |
      | 1     | Brian  |
      | 2     | Alice  |
    And as users "Brian,Carol" the authors of the versions of file "/Shares/exist.txt" should be:
      | index | author |
      | 1     | Brian  |
      | 2     | Alice  |

  @issue-7555
  Scenario: enable file versioning and check the history of changes in sharer after renaming file by sharee
    Given group "grp1" has been created
    And user "Alice" has been added to group "grp1"
    And user "Brian" has been added to group "grp1"
    And user "Carol" has been added to group "grp1"
    And user "Alice" has uploaded file with content "uploaded content alice" to "/exist.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | exist.txt   |
      | space           | Personal    |
      | sharee          | grp1        |
      | shareType       | group       |
      | permissionsRole | File Editor |
    And user "Brian" has uploaded file with content "uploaded content brian" to "/Shares/exist.txt"
    And user "Carol" has uploaded file with content "uploaded content carol" to "/Shares/exist.txt"
    And user "Brian" has moved file "/Shares/exist.txt" to "/Shares/textfile0.txt"
    When user "Alice" gets the number of versions of file "exist.txt"
    Then the HTTP status code should be "207"
    And the number of versions should be "2"
    And the content of version index "1" of file "/exist.txt" for user "Alice" should be "uploaded content brian"
    And the content of version index "2" of file "/exist.txt" for user "Alice" should be "uploaded content alice"
    And as users "Alice" the authors of the versions of file "/exist.txt" should be:
      | index | author |
      | 1     | Brian  |
      | 2     | Alice  |
    And as users "Carol" the authors of the versions of file "/Shares/exist.txt" should be:
      | index | author |
      | 1     | Brian  |
      | 2     | Alice  |
    And as user "Brian" the authors of the versions of file "/Shares/textfile0.txt" should be:
      | index | author |
      | 1     | Brian  |
      | 2     | Alice  |


  Scenario: enable file versioning and check the history of changes from multiple users when reshared after unshared by sharer
    Given user "Alice" has created folder "/test"
    And group "grp1" has been created
    And user "Alice" has been added to group "grp1"
    And user "Brian" has been added to group "grp1"
    And user "Carol" has been added to group "grp1"
    And user "Alice" has sent the following resource share invitation:
      | resource        | test     |
      | space           | Personal |
      | sharee          | grp1     |
      | shareType       | group    |
      | permissionsRole | Editor   |
    And user "Alice" has uploaded file with content "uploaded content alice" to "/test/textfile0.txt"
    And user "Brian" has uploaded file with content "uploaded content brian" to "/Shares/test/textfile0.txt"
    And user "Carol" has uploaded file with content "uploaded content carol" to "/Shares/test/textfile0.txt"
    And user "Alice" has deleted the last share
    And user "Alice" has uploaded file with content "uploaded content alice after unshared folder" to "/test/textfile0.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | test     |
      | space           | Personal |
      | sharee          | grp1     |
      | shareType       | group    |
      | permissionsRole | Editor   |
    When user "Alice" gets the number of versions of file "/test/textfile0.txt"
    Then the HTTP status code should be "207"
    And the number of versions should be "3"
    And the content of version index "1" of file "/test/textfile0.txt" for user "Alice" should be "uploaded content carol"
    And the content of version index "2" of file "/test/textfile0.txt" for user "Alice" should be "uploaded content brian"
    And the content of version index "3" of file "/test/textfile0.txt" for user "Alice" should be "uploaded content alice"
    And as users "Alice" the authors of the versions of file "/test/textfile0.txt" should be:
      | index | author |
      | 1     | Carol  |
      | 2     | Brian  |
      | 3     | Alice  |
    And as users "Brian,Carol" the authors of the versions of file "/Shares/test/textfile0.txt" should be:
      | index | author |
      | 1     | Carol  |
      | 2     | Brian  |
      | 3     | Alice  |


  Scenario: enable file versioning and check the history of changes from multiple users who have a matching folder/file
    Given user "David" has been created with default attributes and without skeleton files
    And user "Brian" has created folder "/test"
    And user "Brian" has uploaded file with content "duplicate brian" to "/test/textfile0.txt"
    And user "Brian" has uploaded file with content "overwrite brian" to "/test/textfile0.txt"
    And user "Brian" has sent the following resource share invitation:
      | resource        | test     |
      | space           | Personal |
      | sharee          | Carol    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has sent the following resource share invitation:
      | resource        | test     |
      | space           | Personal |
      | sharee          | David    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Carol" has created folder "/test"
    And user "Carol" has uploaded file with content "duplicate carol" to "/test/textfile0.txt"
    And user "Carol" has sent the following resource share invitation:
      | resource        | test     |
      | space           | Personal |
      | sharee          | David    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Alice" has created folder "/test"
    And user "Alice" has sent the following resource share invitation:
      | resource        | test     |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Alice" has sent the following resource share invitation:
      | resource        | test     |
      | space           | Personal |
      | sharee          | Carol    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Alice" has sent the following resource share invitation:
      | resource        | test     |
      | space           | Personal |
      | sharee          | David    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Alice" has uploaded file with content "uploaded content alice" to "/test/textfile0.txt"
    And user "Brian" has uploaded file with content "uploaded content brian" to "/Shares/test (1)/textfile0.txt"
    And user "Carol" has uploaded file with content "uploaded content carol" to "/Shares/test (1)/textfile0.txt"
    And user "David" has uploaded file with content "uploaded content david" to "/Shares/test (1)/textfile0.txt"
    When user "Alice" gets the number of versions of file "/test/textfile0.txt"
    Then the HTTP status code should be "207"
    And the number of versions should be "3"
    And the content of version index "1" of file "/test/textfile0.txt" for user "Alice" should be "uploaded content carol"
    And the content of version index "2" of file "/test/textfile0.txt" for user "Alice" should be "uploaded content brian"
    And the content of version index "3" of file "/test/textfile0.txt" for user "Alice" should be "uploaded content alice"
    And as users "Alice" the authors of the versions of file "/test/textfile0.txt" should be:
      | index | author |
      | 1     | Carol  |
      | 2     | Brian  |
      | 3     | Alice  |
    And as users "Brian,Carol,David" the authors of the versions of file "/Shares/test (1)/textfile0.txt" should be:
      | index | author |
      | 1     | Carol  |
      | 2     | Brian  |
      | 3     | Alice  |
    When user "Brian" gets the number of versions of file "/test/textfile0.txt"
    Then the HTTP status code should be "207"
    And the number of versions should be "1"
    And the content of version index "1" of file "/test/textfile0.txt" for user "Brian" should be "duplicate brian"
    And as user "Brian" the authors of the versions of file "/test/textfile0.txt" should be:
      | index | author |
      | 1     | Brian  |
    When user "Carol" gets the number of versions of file "/test/textfile0.txt"
    Then the HTTP status code should be "207"
    And the number of versions should be "0"


  Scenario: enable file versioning and check the history of changes from multiple users who have a matching file
    Given user "David" has been created with default attributes and without skeleton files
    And user "Brian" has uploaded file with content "duplicate brian" to "/textfile0.txt"
    And user "Brian" has uploaded file with content "overwrite brian" to "/textfile0.txt"
    And user "Brian" has sent the following resource share invitation:
      | resource        | textfile0.txt |
      | space           | Personal      |
      | sharee          | Carol         |
      | shareType       | user          |
      | permissionsRole | File Editor   |
    And user "Brian" has sent the following resource share invitation:
      | resource        | textfile0.txt |
      | space           | Personal      |
      | sharee          | David         |
      | shareType       | user          |
      | permissionsRole | File Editor   |
    And user "Carol" has uploaded file with content "duplicate carol" to "/textfile0.txt"
    And user "Carol" has sent the following resource share invitation:
      | resource        | textfile0.txt |
      | space           | Personal      |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | File Editor   |
    And user "Alice" has uploaded file with content "uploaded content alice" to "/textfile0.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | textfile0.txt |
      | space           | Personal      |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | File Editor   |
    And user "Alice" has sent the following resource share invitation:
      | resource        | textfile0.txt |
      | space           | Personal      |
      | sharee          | Carol         |
      | shareType       | user          |
      | permissionsRole | File Editor   |
    And user "Alice" has sent the following resource share invitation:
      | resource        | textfile0.txt |
      | space           | Personal      |
      | sharee          | David         |
      | shareType       | user          |
      | permissionsRole | File Editor   |
    And user "Brian" has uploaded file with content "uploaded content brian" to "/Shares/textfile0 (1).txt"
    And user "Carol" has uploaded file with content "uploaded content carol" to "/Shares/textfile0 (1).txt"
    And user "David" has uploaded file with content "uploaded content david" to "/Shares/textfile0 (1).txt"
    When user "Alice" gets the number of versions of file "/textfile0.txt"
    Then the HTTP status code should be "207"
    And the number of versions should be "3"
    And the content of version index "1" of file "/textfile0.txt" for user "Alice" should be "uploaded content carol"
    And the content of version index "2" of file "/textfile0.txt" for user "Alice" should be "uploaded content brian"
    And the content of version index "3" of file "/textfile0.txt" for user "Alice" should be "uploaded content alice"
    And as users "Alice" the authors of the versions of file "/textfile0.txt" should be:
      | index | author |
      | 1     | Carol  |
      | 2     | Brian  |
      | 3     | Alice  |
    And as users "Brian,Carol,David" the authors of the versions of file "/Shares/textfile0 (1).txt" should be:
      | index | author |
      | 1     | Carol  |
      | 2     | Brian  |
      | 3     | Alice  |
    When user "Brian" gets the number of versions of file "/textfile0.txt"
    Then the HTTP status code should be "207"
    And the number of versions should be "1"
    And the content of version index "1" of file "/textfile0.txt" for user "Brian" should be "duplicate brian"
    And as user "Brian" the authors of the versions of file "/textfile0.txt" should be:
      | index | author |
      | 1     | Brian  |
    When user "Carol" gets the number of versions of file "/textfile0.txt"
    Then the HTTP status code should be "207"
    And the number of versions should be "0"


  Scenario: enable file versioning and check the version author after restoring a version of a file inside a folder
    Given user "Alice" has created folder "/test"
    And user "Alice" has sent the following resource share invitation:
      | resource        | test     |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Alice" has sent the following resource share invitation:
      | resource        | test     |
      | space           | Personal |
      | sharee          | Carol    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Alice" has uploaded file with content "uploaded content alice" to "/test/textfile0.txt"
    And user "Brian" has uploaded file with content "uploaded content brian" to "/Shares/test/textfile0.txt"
    And user "Carol" has uploaded file with content "uploaded content carol" to "/Shares/test/textfile0.txt"
    When user "Brian" restores version index "1" of file "/Shares/test/textfile0.txt" using the WebDAV API
    Then the HTTP status code should be "204"
    And as user "Alice" the authors of the versions of file "/test/textfile0.txt" should be:
      | index | author |
      | 1     | Carol  |
      | 2     | Alice  |
    And as user "Brian,Carol" the authors of the versions of file "/Shares/test/textfile0.txt" should be:
      | index | author |
      | 1     | Carol  |
      | 2     | Alice  |


  Scenario: enable file versioning and check the version author after restoring a version of a file
    Given user "Alice" has uploaded file with content "uploaded content alice" to "/textfile0.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | textfile0.txt |
      | space           | Personal      |
      | sharee          | Brian         |
      | shareType       | user          |
      | permissionsRole | File Editor   |
    And user "Alice" has sent the following resource share invitation:
      | resource        | textfile0.txt |
      | space           | Personal      |
      | sharee          | Carol         |
      | shareType       | user          |
      | permissionsRole | File Editor   |
    And user "Brian" has uploaded file with content "uploaded content brian" to "/Shares/textfile0.txt"
    And user "Carol" has uploaded file with content "uploaded content carol" to "/Shares/textfile0.txt"
    When user "Brian" restores version index "1" of file "/Shares/textfile0.txt" using the WebDAV API
    Then the HTTP status code should be "204"
    And as user "Alice" the authors of the versions of file "/textfile0.txt" should be:
      | index | author |
      | 1     | Carol  |
      | 2     | Alice  |
    And as user "Brian,Carol" the authors of the versions of file "/Shares/textfile0.txt" should be:
      | index | author |
      | 1     | Carol  |
      | 2     | Alice  |
