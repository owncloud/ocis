@issue-1284
Feature: independent locks - make sure all locks are independent and don't interact with other items that have the same name
  As a user
  I want to lock resources independently
  So that resources with same name in other parts of the file system will not be locked

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files


  Scenario Outline: locking a file does not lock other items with the same name in other parts of the file system
    Given using <dav-path-version> DAV path
    And user "Alice" has created folder "locked"
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "/locked/textfile0.txt"
    And user "Alice" has created folder "notlocked"
    And user "Alice" has created folder "notlocked/textfile0.txt"
    When user "Alice" locks file "locked/textfile0.txt" using the WebDAV API setting the following properties
      | lockscope | <lock-scope> |
    Then the HTTP status code should be "200"
    And user "Alice" should be able to upload file "filesForUpload/lorem.txt" to "/notlocked/textfile0.txt/real-file.txt"
    And user "Alice" should be able to upload file "filesForUpload/lorem.txt" to "/textfile0.txt"
    But user "Alice" should not be able to upload file "filesForUpload/lorem.txt" to "/locked/textfile0.txt"
    Examples:
      | dav-path-version | lock-scope |
      | old              | shared     |
      | old              | exclusive  |
      | new              | shared     |
      | new              | exclusive  |

    @skipOnRevaMaster
    Examples:
      | dav-path-version | lock-scope |
      | spaces           | shared     |
      | spaces           | exclusive  |


  Scenario Outline: locking a file/folder with git specific names does not lock other items with the same name in other parts of the file system
    Given using <dav-path-version> DAV path
    And user "Alice" has created folder "locked/"
    And user "Alice" has created folder "locked/.git"
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "/locked/.git/config"
    And user "Alice" has created folder "notlocked/"
    And user "Alice" has created folder "notlocked/"
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "notlocked/.git/config"
    When user "Alice" locks file "locked/<to-lock>" using the WebDAV API setting the following properties
      | lockscope | <lock-scope> |
    Then the HTTP status code should be "200"
    And user "Alice" should be able to upload file "filesForUpload/lorem.txt" to "/notlocked/.git/file.txt"
    And user "Alice" should be able to upload file "filesForUpload/lorem.txt" to "/notlocked/.git/config"
    But user "Alice" should not be able to upload file "filesForUpload/lorem.txt" to "/locked/.git/config"
    Examples:
      | dav-path-version | lock-scope | to-lock     |
      | old              | shared     | .git        |
      | old              | shared     | .git/config |
      | old              | exclusive  | .git        |
      | old              | exclusive  | .git/config |
      | new              | shared     | .git        |
      | new              | shared     | .git/config |
      | new              | exclusive  | .git        |
      | new              | exclusive  | .git/config |

    @skipOnRevaMaster
    Examples:
      | dav-path-version | lock-scope | to-lock     |
      | spaces           | shared     | .git        |
      | spaces           | shared     | .git/config |
      | spaces           | exclusive  | .git        |
      | spaces           | exclusive  | .git/config |
