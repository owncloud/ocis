@api
Feature: full text search
  As a user
  I want to do full text search
  So that I can find the files with the content I am looking for

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files


  Scenario Outline: search files using a tag
    Given using <dav-path-version> DAV path
    And user "Alice" has uploaded file with content "hello world" to "file1.txt"
    And user "Alice" has uploaded file with content "Namaste nepal" to "file2.txt"
    And user "Alice" has uploaded file with content "hello nepal" to "file3.txt"
    And user "Alice" has created the following tags for file "file1.txt" of the space "Personal":
      | tag1 |
    And user "Alice" has created the following tags for file "file2.txt" of the space "Personal":
      | tag1 |
    When user "Alice" searches for "Tags:tag1" using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result of user "Alice" should contain only these files:
      | file1.txt |
      | file2.txt |
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: search folders using a tag
    Given using <dav-path-version> DAV path
    And user "Alice" has created folder "uploadFolder1"
    And user "Alice" has created folder "uploadFolder2"
    And user "Alice" has created folder "uploadFolder3"
    And user "Alice" has created the following tags for folder "uploadFolder1" of the space "Personal":
      | tag1 |
    And user "Alice" has created the following tags for folder "uploadFolder2" of the space "Personal":
      | tag1 |
    When user "Alice" searches for "Tags:tag1" using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result of user "Alice" should contain only these entries:
      | uploadFolder1 |
      | uploadFolder2 |
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: sharee searches shared files using a tag
    Given using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has created folder "uploadFolder"
    And user "Alice" has uploaded file with content "hello world" to "uploadFolder/file1.txt"
    And user "Alice" has uploaded file with content "Namaste nepal" to "uploadFolder/file2.txt"
    And user "Alice" has uploaded file with content "hello nepal" to "uploadFolder/file3.txt"
    And user "Alice" has created the following tags for file "uploadFolder/file1.txt" of the space "Personal":
      | tag1 |
    And user "Alice" has shared folder "/uploadFolder" with user "Brian"
    And user "Brian" has accepted share "/uploadFolder" offered by user "Alice"
    And user "Brian" has created the following tags for file "uploadFolder/file2.txt" of the space "Shares":
      | tag1 |
    When user "Brian" searches for "Tags:tag1" using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result of user "Brian" should contain only these files:
      | file1.txt |
      | file2.txt |
    Examples:
      | dav-path-version |
      | old              |
      | new              |


  Scenario Outline: search files using a deleted tag
    Given using <dav-path-version> DAV path
    And user "Alice" has uploaded file with content "hello world" to "file1.txt"
    And user "Alice" has created the following tags for file "file1.txt" of the space "Personal":
      | tag1 |
    And user "Alice" has removed the following tags for file "file1.txt" of space "Personal":
      | tag1 |
    When user "Alice" searches for "Tags:tag1" using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result should contain "0" entries
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |
