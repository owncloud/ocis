@api @issue-ocis-reva-172
Feature: UNLOCK locked items

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files


  Scenario Outline: unlock one of multiple locks set by the user itself
    Given using <dav-path> DAV path
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "textfile0.txt"
    And user "Alice" has locked file "textfile0.txt" setting the following properties
      | lockscope | shared |
    And user "Alice" has locked file "textfile0.txt" setting the following properties
      | lockscope | shared |
    When user "Alice" unlocks the last created lock of file "textfile0.txt" using the WebDAV API
    Then the HTTP status code should be "204"
    And 1 locks should be reported for file "textfile0.txt" of user "Alice" by the WebDAV API
    Examples:
      | dav-path |
      | old      |
      | new      |

    @personalSpace @skipOnOcV10
    Examples:
      | dav-path |
      | spaces   |

  @files_sharing-app-required
  Scenario Outline: as public unlocking a file in a share that was locked by the file owner is not possible. To unlock use the owners locktoken
    Given user "Alice" has created folder "PARENT"
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "PARENT/parent.txt"
    And user "Alice" has created a public link share of folder "PARENT" with change permission
    And user "Alice" has locked file "PARENT/parent.txt" setting the following properties
      | lockscope | <lock-scope> |
    When the public unlocks file "/parent.txt" with the last created lock of file "PARENT/parent.txt" of user "Alice" using the WebDAV API
    Then the HTTP status code should be "403"
    And 1 locks should be reported for file "PARENT/parent.txt" of user "Alice" by the WebDAV API
    Examples:
      | lock-scope |
      | shared     |
      | exclusive  |


  Scenario Outline: unlocking a file or folder does not unlock another file with the same name in another part of the file system
    Given using <dav-path> DAV path
    And user "Alice" has created folder "locked"
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "/locked/textfile0.txt"
    And user "Alice" has created folder "notlocked"
    And user "Alice" has created folder "notlocked/textfile0.txt"
    And user "Alice" has locked file "locked/textfile0.txt" setting the following properties
      | lockscope | <lock-scope> |
    And user "Alice" has locked file "notlocked/textfile0.txt" setting the following properties
      | lockscope | <lock-scope> |
    And user "Alice" has locked file "textfile0.txt" setting the following properties
      | lockscope | <lock-scope> |
    When user "Alice" unlocks the last created lock of file "notlocked/textfile0.txt" using the WebDAV API
    And user "Alice" unlocks the last created lock of file "textfile0.txt" using the WebDAV API
    Then user "Alice" should be able to upload file "filesForUpload/lorem.txt" to "/notlocked/textfile0.txt/real-file.txt"
    And user "Alice" should be able to upload file "filesForUpload/lorem.txt" to "/textfile0.txt"
    But user "Alice" should not be able to upload file "filesForUpload/lorem.txt" to "/locked/textfile0.txt"
    Examples:
      | dav-path | lock-scope |
      | old      | shared     |
      | old      | exclusive  |
      | new      | shared     |
      | new      | exclusive  |

    @personalSpace @skipOnOcV10
    Examples:
      | dav-path | lock-scope |
      | spaces   | shared     |
      | spaces   | exclusive  |
