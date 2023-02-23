@api @issue-ocis-reva-15
Feature: users cannot upload a file to a blacklisted name using old chunking
  As an administrator
  I want to be able to prevent users from uploading files to specified file names
  So that I can prevent unwanted file names existing in the cloud storage

  Background:
    Given using OCS API version "1"
    And using old DAV path
    And user "Alice" has been created with default attributes and without skeleton files


  Scenario: Upload a file to a banned filename using old chunking
    When user "Alice" uploads file "filesForUpload/textfile.txt" to "blacklisted-file.txt" in 3 chunks using the WebDAV API
    Then the HTTP status code should be "403"
    And as "Alice" file "blacklisted-file.txt" should not exist


  Scenario Outline: upload a file to a filename that matches blacklisted_files_regex using old chunking
    When user "Alice" uploads file "filesForUpload/textfile.txt" to "<filename>" in 3 chunks using the WebDAV API
    Then the HTTP status code should be "<http-status>"
    And as "Alice" file "<filename>" should not exist
    Examples:
      | filename                      | http-status |
      | filename.ext                  | 403         |
      | bannedfilename.txt            | 403         |
      | this-ContainsBannedString.txt | 403         |


  Scenario: upload a file to a filename that does not match blacklisted_files_regex using old chunking
    When user "Alice" uploads file "filesForUpload/textfile.txt" to "not-contains-banned-string.txt" in 3 chunks using the WebDAV API
    Then the HTTP status code should be "201"
    And as "Alice" file "not-contains-banned-string.txt" should exist
