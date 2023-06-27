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
