Feature: lock files
  As a user
  I want to lock files

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |


  Scenario Outline: lock a file
    Given using <dav-path-version> DAV path
    And user "Alice" has uploaded a file inside space "Alice Hansen" with content "some content" to "textfile.txt"
    When user "Alice" locks file "textfile.txt" using the WebDAV API setting the following properties
      | lockscope | exclusive |
    Then the HTTP status code should be "200"
    When user "Alice" sends PROPFIND request from the space "Alice Hansen" to the resource "textfile.txt" using the WebDAV API
    Then the HTTP status code should be "207"
    And the "PROPFIND" response to user "Alice" should contain a space "Alice Hansen" with these key and value pairs:
      | key                                                  | value        |
      | d:lockdiscovery/d:activelock/d:lockscope/d:exclusive |              |
      | d:lockdiscovery/d:activelock/d:depth                 | Infinity     |
      | d:lockdiscovery/d:activelock/d:timeout               | Infinity     |
      # | d:lockdiscovery/d:activelock/oc:ownername            | Alice Hansen |  no "oc:ownername" property in stable-4.0
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: lock a file with a timeout
    Given using <dav-path-version> DAV path
    And user "Alice" has uploaded a file inside space "Alice Hansen" with content "some content" to "textfile.txt"
    When user "Alice" locks file "textfile.txt" using the WebDAV API setting the following properties
      | lockscope | exclusive   |
      | timeout   | Second-5000 |
    Then the HTTP status code should be "200"
    When user "Alice" sends PROPFIND request from the space "Alice Hansen" to the resource "textfile.txt" using the WebDAV API
    Then the HTTP status code should be "207"
    And the "PROPFIND" response to user "Alice" should contain a space "Alice Hansen" with these key and value pairs:
      | key                                                  | value        |
      | d:lockdiscovery/d:activelock/d:lockscope/d:exclusive |              |
      | d:lockdiscovery/d:activelock/d:depth                 | Infinity     |
      | d:lockdiscovery/d:activelock/d:timeout               | Second-5000  |
      # | d:lockdiscovery/d:activelock/oc:ownername            | Alice Hansen |  no "oc:ownername" property in stable-4.0
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: lock a file using file-id
    Given user "Alice" has uploaded a file inside space "Alice Hansen" with content "some content" to "textfile.txt"
    And we save it into "FILEID"
    When user "Alice" locks file using file-id path "<dav-path>" using the WebDAV API setting the following properties
      | lockscope | exclusive   |
      | timeout   | Second-3600 |
    Then the HTTP status code should be "200"
    When user "Alice" sends PROPFIND request from the space "Alice Hansen" to the resource "textfile.txt" using the WebDAV API
    Then the HTTP status code should be "207"
    And the "PROPFIND" response to user "Alice" should contain a space "Alice Hansen" with these key and value pairs:
      | key                                                  | value        |
      | d:lockdiscovery/d:activelock/d:lockscope/d:exclusive |              |
      | d:lockdiscovery/d:activelock/d:depth                 | Infinity     |
      | d:lockdiscovery/d:activelock/d:timeout               | Second-3600  |
      # | d:lockdiscovery/d:activelock/oc:ownername            | Alice Hansen |  no "oc:ownername" property in stable-4.0
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


  Scenario Outline: user cannot lock file twice
    Given using <dav-path-version> DAV path
    And user "Alice" has uploaded a file inside space "Alice Hansen" with content "some content" to "textfile.txt"
    And user "Alice" has locked file "textfile.txt" setting the following properties
      | lockscope | exclusive |
    When user "Alice" tries to lock file "textfile.txt" using the WebDAV API setting the following properties
      | lockscope | exclusive |
    Then the HTTP status code should be "423"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: lock a file in the project space
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And using spaces DAV path
    And user "Alice" has created a space "Project" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "Project" with content "some content" to "textfile.txt"
    And user "Alice" has shared a space "Project" with settings:
      | shareWith | Brian  |
      | role      | <role> |
    When user "Brian" locks file "textfile.txt" inside the space "Project" using the WebDAV API setting the following properties
      | lockscope | exclusive   |
      | timeout   | Second-3600 |
    Then the HTTP status code should be "200"
    When user "Brian" sends PROPFIND request from the space "Project" to the resource "textfile.txt" using the WebDAV API
    Then the HTTP status code should be "207"
    And the "PROPFIND" response to user "Brian" should contain a space "Project" with these key and value pairs:
      | key                                                  | value        |
      | d:lockdiscovery/d:activelock/d:lockscope/d:exclusive |              |
      | d:lockdiscovery/d:activelock/d:depth                 | Infinity     |
      | d:lockdiscovery/d:activelock/d:timeout               | Second-3600  |
      # | d:lockdiscovery/d:activelock/oc:ownername            | Brian Murphy |  no "oc:ownername" property in stable-4.0
    Examples:
      | role    |
      | manager |
      | editor  |


  Scenario Outline: lock a file in the project space using file-id
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And using spaces DAV path
    And user "Alice" has created a space "Project" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "Project" with content "some content" to "textfile.txt"
    And we save it into "FILEID"
    And user "Alice" has shared a space "Project" with settings:
      | shareWith | Brian  |
      | role      | <role> |
    When user "Brian" locks file using file-id path "<dav-path>" using the WebDAV API setting the following properties
      | lockscope | exclusive   |
      | timeout   | Second-3600 |
    Then the HTTP status code should be "200"
    When user "Brian" sends PROPFIND request from the space "Project" to the resource "textfile.txt" using the WebDAV API
    Then the HTTP status code should be "207"
    And the "PROPFIND" response to user "Brian" should contain a space "Project" with these key and value pairs:
      | key                                                  | value        |
      | d:lockdiscovery/d:activelock/d:lockscope/d:exclusive |              |
      | d:lockdiscovery/d:activelock/d:depth                 | Infinity     |
      | d:lockdiscovery/d:activelock/d:timeout               | Second-3600  |
      # | d:lockdiscovery/d:activelock/oc:ownername            | Brian Murphy |  no "oc:ownername" property in stable-4.0
    Examples:
      | role    | dav-path                          |
      | manager | /remote.php/dav/spaces/<<FILEID>> |
      | editor  | /dav/spaces/<<FILEID>>            |


  Scenario: viewer cannot lock a file in the project space
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And using spaces DAV path
    And user "Alice" has created a space "Project" with the default quota using the Graph API
    And user "Alice" has uploaded a file inside space "Project" with content "some content" to "textfile.txt"
    And we save it into "FILEID"
    And user "Alice" has shared a space "Project" with settings:
      | shareWith | Brian  |
      | role      | viewer |
    When user "Brian" locks file "textfile.txt" inside the space "Project" using the WebDAV API setting the following properties
      | lockscope | exclusive |
    Then the HTTP status code should be "403"
    When user "Brian" locks file using file-id path "/dav/spaces/<<FILEID>>" using the WebDAV API setting the following properties
      | lockscope | exclusive |
    Then the HTTP status code should be "403"


  Scenario Outline: lock a file in the shares
    Given using <dav-path-version> DAV path
    And user "Alice" has uploaded a file inside space "Alice Hansen" with content "some content" to "textfile.txt"
    And user "Alice" has created a share inside of space "Alice Hansen" with settings:
      | path      | textfile.txt |
      | shareWith | Brian        |
      | role      | editor       |
    When user "Brian" locks file "/Shares/textfile.txt" using the WebDAV API setting the following properties
      | lockscope | exclusive |
    Then the HTTP status code should be "200"
    When user "Alice" sends PROPFIND request from the space "Alice Hansen" to the resource "textfile.txt" using the WebDAV API
    Then the HTTP status code should be "207"
    And the "PROPFIND" response to user "Alice" should contain a space "Alice Hansen" with these key and value pairs:
      | key                                                  | value        |
      | d:lockdiscovery/d:activelock/d:lockscope/d:exclusive |              |
      # | d:lockdiscovery/d:activelock/oc:ownername            | Brian Murphy |  no "oc:ownername" property in stable-4.0
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: lock a file in the shares using file-id
    Given user "Alice" has uploaded a file inside space "Alice Hansen" with content "some content" to "textfile.txt"
    And we save it into "FILEID"
    And user "Alice" has created a share inside of space "Alice Hansen" with settings:
      | path      | textfile.txt |
      | shareWith | Brian        |
      | role      | editor       |
    When user "Brian" locks file using file-id path "<dav-path>" using the WebDAV API setting the following properties
      | lockscope | exclusive   |
      | timeout   | Second-3600 |
    Then the HTTP status code should be "200"
    When user "Alice" sends PROPFIND request from the space "Alice Hansen" to the resource "textfile.txt" using the WebDAV API
    Then the HTTP status code should be "207"
    And the "PROPFIND" response to user "Alice" should contain a space "Alice Hansen" with these key and value pairs:
      | key                                                  | value        |
      | d:lockdiscovery/d:activelock/d:lockscope/d:exclusive |              |
      # | d:lockdiscovery/d:activelock/oc:ownername            | Brian Murphy |  no "oc:ownername" property in stable-4.0
    Examples:
      | dav-path                          |
      | /remote.php/dav/spaces/<<FILEID>> |
      | /dav/spaces/<<FILEID>>            |


   Scenario: viewer cannot lock a file in the shares using file-id
    Given user "Alice" has uploaded a file inside space "Alice Hansen" with content "some content" to "textfile.txt"
    And we save it into "FILEID"
    And user "Alice" has created a share inside of space "Alice Hansen" with settings:
      | path      | textfile.txt |
      | shareWith | Brian        |
      | role      | viewer       |
    When user "Brian" locks file using file-id path "<dav-path>" using the WebDAV API setting the following properties
      | lockscope | exclusive   |
    Then the HTTP status code should be "403"
