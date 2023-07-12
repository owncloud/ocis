@api
Feature: full text search
  As a user
  I want to do full text search
  So that I can find the files with the content I am looking for

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files


  Scenario Outline: search files by content
    Given using <dav-path-version> DAV path
    And user "Alice" has uploaded file with content "hello world from nepal" to "keywordAtStart.txt"
    And user "Alice" has uploaded file with content "saying hello to the world" to "keywordAtMiddle.txt"
    And user "Alice" has uploaded file with content "nepal want to say hello" to "keywordAtLast.txt"
    And user "Alice" has uploaded file with content "namaste from nepal" to "hello.txt"
    When user "Alice" searches for "Content:hello" using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result of user "Alice" should contain only these files:
      | keywordAtStart.txt  |
      | keywordAtMiddle.txt |
      | keywordAtLast.txt   |
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |
