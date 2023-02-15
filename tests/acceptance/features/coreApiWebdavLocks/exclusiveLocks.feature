@api @issue-ocis-reva-172
Feature: there can be only one exclusive lock on a resource

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files


  Scenario Outline: a second lock cannot be set on a folder when its exclusively locked
    Given using <dav-path> DAV path
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "textfile0.txt"
    And user "Alice" has locked file "textfile0.txt" setting the following properties
      | lockscope | exclusive |
    When user "Alice" locks file "textfile0.txt" using the WebDAV API setting the following properties
      | lockscope | <lock-scope> |
    Then the HTTP status code should be "423"
    And 1 locks should be reported for file "textfile0.txt" of user "Alice" by the WebDAV API
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

  @files_sharing-app-required
  Scenario Outline: a share receiver cannot lock a resource exclusively locked by itself
    Given using <dav-path> DAV path
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "textfile0.txt"
    And user "Brian" has uploaded file "filesForUpload/textfile.txt" to "textfile0.txt"
    And user "Alice" has shared file "textfile0.txt" with user "Brian"
    And user "Brian" has locked file "textfile0 (2).txt" setting the following properties
      | lockscope | exclusive |
    When user "Brian" locks file "textfile0 (2).txt" using the WebDAV API setting the following properties
      | lockscope | <lock-scope> |
    Then the HTTP status code should be "423"
    And 1 locks should be reported for file "textfile0.txt" of user "Alice" by the WebDAV API
    And 1 locks should be reported for file "textfile0 (2).txt" of user "Brian" by the WebDAV API
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

  @files_sharing-app-required
  Scenario Outline: a share receiver cannot lock a resource exclusively locked by the owner
    Given using <dav-path> DAV path
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "textfile0.txt"
    And user "Brian" has uploaded file "filesForUpload/textfile.txt" to "textfile0.txt"
    And user "Alice" has shared file "textfile0.txt" with user "Brian"
    And user "Alice" has locked file "textfile0.txt" setting the following properties
      | lockscope | exclusive |
    When user "Brian" locks file "textfile0 (2).txt" using the WebDAV API setting the following properties
      | lockscope | <lock-scope> |
    Then the HTTP status code should be "423"
    And 1 locks should be reported for file "textfile0.txt" of user "Alice" by the WebDAV API
    And 1 locks should be reported for file "textfile0 (2).txt" of user "Brian" by the WebDAV API
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

  @files_sharing-app-required
  Scenario Outline: a share owner cannot lock a resource exclusively locked by a share receiver
    Given using <dav-path> DAV path
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "textfile0.txt"
    And user "Brian" has uploaded file "filesForUpload/textfile.txt" to "textfile0.txt"
    And user "Alice" has shared file "textfile0.txt" with user "Brian"
    And user "Brian" has locked file "textfile0 (2).txt" setting the following properties
      | lockscope | exclusive |
    When user "Alice" locks file "textfile0.txt" using the WebDAV API setting the following properties
      | lockscope | <lock-scope> |
    Then the HTTP status code should be "423"
    And 1 locks should be reported for file "textfile0.txt" of user "Alice" by the WebDAV API
    And 1 locks should be reported for file "textfile0 (2).txt" of user "Brian" by the WebDAV API
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
