@api @skipOnOcV10
Feature: lock
  # Note: This Feature includes all the tests from core (apiWebdavLock suite) related to /Shares since in core no implementation is there for space Shares Jail

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
      | Carol    |
    And user "Alice" has created folder "PARENT"
    And user "Brian" has created folder "PARENT"
    And user "Carol" has created folder "PARENT"
    And using spaces DAV path


  Scenario Outline: lock should propagate correctly when uploaded to a reshare that was locked by the owner
    Given user "Alice" has shared folder "PARENT" with user "Brian"
    And user "Brian" has accepted share "/PARENT" offered by user "Alice"
    And user "Brian" has shared the following entity "PARENT" inside of space "Shares Jail" with user "Carol" with role "editor"
    And user "Carol" has accepted share "/PARENT" offered by user "Brian"
    And user "Alice" has locked folder "/PARENT" inside space "Personal" setting the following properties
      | lockscope | <lock-scope> |
    When user "Carol" uploads a file inside space "Shares Jail" with content "uploaded by carol" to "PARENT/textfile.txt" using the WebDAV API
    And user "Brian" uploads a file inside space "Shares Jail" with content "uploaded by brian" to "PARENT/textfile.txt" using the WebDAV API
    And user "Alice" uploads file "filesForUpload/textfile.txt" to "/PARENT/textfile.txt" using the WebDAV API
    Then the HTTP status code of responses on all endpoints should be "423"
    And as "Alice" file "/PARENT/textfile.txt" should not exist
    Examples:
      | lock-scope |
      | shared     |
      | exclusive  |


  Scenario Outline: lock should propagate correctly when uploaded overwriting to a reshare that was locked by the owner
    Given user "Alice" has uploaded file with content "ownCloud test text file parent" to "PARENT/parent.txt"
    And user "Alice" has shared folder "PARENT" with user "Brian"
    And user "Brian" has accepted share "/PARENT" offered by user "Alice"
    And user "Brian" has shared the following entity "PARENT" inside of space "Shares Jail" with user "Carol" with role "editor"
    And user "Carol" has accepted share "/PARENT" offered by user "Brian"
    And user "Alice" has locked folder "/PARENT" inside space "Personal" setting the following properties
      | lockscope | <lock-scope> |
    When user "Carol" uploads a file inside space "Shares Jail" with content "uploaded by carol" to "PARENT/textfile.txt" using the WebDAV API
    And user "Brian" uploads a file inside space "Shares Jail" with content "uploaded by brian" to "PARENT/textfile.txt" using the WebDAV API
    And user "Alice" uploads file "filesForUpload/textfile.txt" to "/PARENT/parent.txt" using the WebDAV API
    Then the HTTP status code of responses on all endpoints should be "423"
    And the content of file "/PARENT/parent.txt" for user "Alice" should be "ownCloud test text file parent"
    Examples:
      | lock-scope |
      | shared     |
      | exclusive  |


  Scenario Outline: lock should propagate correctly when the public uploads to a reshared share that was locked by the original owner
    Given user "Alice" has shared folder "PARENT" with user "Brian"
    And user "Brian" has accepted share "/PARENT" offered by user "Alice"
    And user "Brian" has shared the following entity "PARENT" inside of space "Shares Jail" with user "Carol" with role "editor"
    And user "Carol" has accepted share "/PARENT" offered by user "Brian"
    And user "Carol" has created a public link share inside of space "Shares Jail" with settings:
      | path        | PARENT      |
      | shareType   | 3           |
      | permissions | 15          |
      | name        | public link |
    And user "Alice" has locked folder "/PARENT" inside space "Personal" setting the following properties
      | lockscope | <lock-scope> |
    When the public uploads file "test.txt" with content "test" using the new public WebDAV API
    Then the HTTP status code should be "423"
    And as "Alice" file "/PARENT/test.txt" should not exist
    Examples:
      | lock-scope |
      | shared     |
      | exclusive  |


  Scenario Outline: lock should propagate correctly when uploaded to a reshare that was locked by the resharing user
    Given user "Alice" has shared folder "PARENT" with user "Brian"
    And user "Brian" has accepted share "/PARENT" offered by user "Alice"
    And user "Brian" has shared the following entity "PARENT" inside of space "Shares Jail" with user "Carol" with role "editor"
    And user "Carol" has accepted share "/PARENT" offered by user "Brian"
    And user "Brian" has locked folder "/PARENT" inside space "Shares Jail" setting the following properties
      | lockscope | <lock-scope> |
    When user "Carol" uploads a file inside space "Shares Jail" with content "uploaded by carol" to "PARENT/textfile.txt" using the WebDAV API
    And user "Brian" uploads a file inside space "Shares Jail" with content "uploaded by brian" to "PARENT/textfile.txt" using the WebDAV API
    And user "Alice" uploads file "filesForUpload/textfile.txt" to "/PARENT/textfile.txt" using the WebDAV API
    Then the HTTP status code of responses on all endpoints should be "423"
    And as "Alice" file "/PARENT/textfile.txt" should not exist
    Examples:
      | lock-scope |
      | shared     |
      | exclusive  |

  Scenario Outline: set timeouts of LOCKS on shares as owner set timeout on folder as receiver check it
    Given user "Alice" has created folder "PARENT/CHILD"
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "PARENT/parent.txt"
    And user "Brian" has created folder "PARENT/CHILD"
    And user "Brian" has uploaded file "filesForUpload/textfile.txt" to "PARENT/parent.txt"
    And user "Alice" has shared folder "PARENT" with user "Brian"
    And user "Brian" has accepted share "/PARENT" offered by user "Alice"
    When user "Alice" locks folder "/PARENT" inside space "Personal" setting the following properties
      | lockscope | shared    |
      | timeout   | <timeout> |
    Then the HTTP status code should be "200"
    And as user "Brian" the lock discovery property "//d:timeout" of the folder "/PARENT" inside space "Shares Jail" should match "<result>"
    And as user "Brian" the lock discovery property "//d:timeout" of the folder "/PARENT/CHILD" inside space "Shares Jail" should match "<result>"
    And as user "Brian" the lock discovery property "//d:timeout" of the folder "/PARENT/parent.txt" inside space "Shares Jail" should match "<result>"
    Examples:
      | timeout         | result          |
      | Second-999      | /Second-\d{3}$/ |
      | Second-99999999 | /Second-\d{5}$/ |
      | Infinite        | /Second-\d{5}$/ |
      | Second-1        | /Second-\d{5}$/ |
      | Second-0        | /Second-\d{4}$/ |

  Scenario Outline: set timeouts of LOCKS on shares as share receiver set timeout on folder as owner check it
    Given user "Alice" has created folder "PARENT/CHILD"
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "PARENT/parent.txt"
    And user "Brian" has created folder "PARENT/CHILD"
    And user "Brian" has uploaded file "filesForUpload/textfile.txt" to "PARENT/parent.txt"
    And user "Alice" has shared folder "PARENT" with user "Brian"
    And user "Brian" has accepted share "/PARENT" offered by user "Alice"
    When user "Brian" locks folder "/PARENT" inside space "Shares Jail" setting the following properties
      | lockscope | shared    |
      | timeout   | <timeout> |
    Then the HTTP status code should be "200"
    And as user "Alice" the lock discovery property "//d:timeout" of the folder "/PARENT" inside space "Shares Jail" should match "<result>"
    And as user "Alice" the lock discovery property "//d:timeout" of the folder "/PARENT/CHILD" inside space "Shares Jail" should match "<result>"
    And as user "Alice" the lock discovery property "//d:timeout" of the folder "/PARENT/parent.txt" inside space "Shares Jail" should match "<result>"
    Examples:
      | timeout         | result          |
      | Second-999      | /Second-\d{3}$/ |
      | Second-99999999 | /Second-\d{5}$/ |
      | Infinite        | /Second-\d{5}$/ |
      | Second-1        | /Second-\d{5}$/ |
      | Second-0        | /Second-\d{4}$/ |
