@issue-1284
Feature: UNLOCK locked items (sharing)
  As a user
  I want to unlock a shared resource that has been locked by me to be restricted
  So that other users cannot unlock the shared resource

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
    And user "Alice" has created folder "PARENT"
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "PARENT/parent.txt"


  Scenario Outline: unlocking a shared file that has been locked by the file owner is not feasible unless the owner lock tocken is used
    Given using <dav-path-version> DAV path
    And user "Alice" has locked file "PARENT/parent.txt" setting the following properties
      | lockscope | <lock-scope> |
    And user "Alice" has shared file "PARENT/parent.txt" with user "Brian"
    And user "Brian" has accepted share "/parent.txt" offered by user "Alice"
    When user "Brian" unlocks file "Shares/parent.txt" with the last created lock of file "PARENT/parent.txt" of user "Alice" using the WebDAV API
    Then the HTTP status code should be "403"
    And 1 locks should be reported for file "PARENT/parent.txt" of user "Alice" by the WebDAV API
    And 1 locks should be reported for file "Shares/parent.txt" of user "Brian" by the WebDAV API
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


  Scenario Outline: sharee cannot unlock a file within a shared folder when it is locked by the file owner unless the owner lock token is used
    Given using <dav-path-version> DAV path
    And user "Alice" has locked file "PARENT/parent.txt" setting the following properties
      | lockscope | <lock-scope> |
    And user "Alice" has shared folder "PARENT" with user "Brian"
    And user "Brian" has accepted share "/PARENT" offered by user "Alice"
    When user "Brian" unlocks file "Shares/PARENT/parent.txt" with the last created lock of file "PARENT/parent.txt" of user "Alice" using the WebDAV API
    Then the HTTP status code should be "403"
    And 1 locks should be reported for file "PARENT/parent.txt" of user "Alice" by the WebDAV API
    And 1 locks should be reported for file "Shares/PARENT/parent.txt" of user "Brian" by the WebDAV API
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


  Scenario Outline: sharee unlock a shared file
    Given using <dav-path-version> DAV path
    And user "Alice" has shared file "PARENT/parent.txt" with user "Brian"
    And user "Brian" has accepted share "/parent.txt" offered by user "Alice"
    And user "Brian" has locked file "Shares/parent.txt" setting the following properties
      | lockscope | <lock-scope> |
    When user "Brian" unlocks the last created lock of file "Shares/parent.txt" using the WebDAV API
    Then the HTTP status code should be "204"
    And 0 locks should be reported for file "PARENT/parent.txt" of user "Alice" by the WebDAV API
    And 0 locks should be reported for file "Shares/parent.txt" of user "Brian" by the WebDAV API
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


  Scenario Outline: as owner unlocking a shared file locked by the receiver is not possible. To unlock use the receivers locktoken
    Given using <dav-path-version> DAV path
    And user "Alice" has shared file "PARENT/parent.txt" with user "Brian"
    And user "Brian" has accepted share "/parent.txt" offered by user "Alice"
    And user "Brian" has locked file "Shares/parent.txt" setting the following properties
      | lockscope | <lock-scope> |
    When user "Alice" unlocks file "PARENT/parent.txt" with the last created lock of file "Shares/parent.txt" of user "Brian" using the WebDAV API
    Then the HTTP status code should be "403"
    And 1 locks should be reported for file "PARENT/parent.txt" of user "Alice" by the WebDAV API
    And 1 locks should be reported for file "Shares/parent.txt" of user "Brian" by the WebDAV API
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


  Scenario Outline: unlocking a file in a shared folder, which was locked by the sharee is not possible for the owner unless the receiver's locktoken is used
    Given using <dav-path-version> DAV path
    And user "Alice" has shared folder "PARENT" with user "Brian"
    And user "Brian" has accepted share "/PARENT" offered by user "Alice"
    And user "Brian" has locked file "Shares/PARENT/parent.txt" setting the following properties
      | lockscope | <lock-scope> |
    When user "Alice" unlocks file "PARENT/parent.txt" with the last created lock of file "Shares/PARENT/parent.txt" of user "Brian" using the WebDAV API
    Then the HTTP status code should be "403"
    And 1 locks should be reported for file "PARENT/parent.txt" of user "Alice" by the WebDAV API
    And 1 locks should be reported for file "Shares/PARENT/parent.txt" of user "Brian" by the WebDAV API
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
