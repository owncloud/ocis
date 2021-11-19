Feature: download file
  As a user
  I want to be able to download files
  So that I can work wih local copies of files on my client system

  Background:
    Given using "oc10" as owncloud selector
    And user "Alice" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file with content "ownCloud test text file" to "textfile.txt"


  Scenario Outline: download a file
    Given using "ocis" as owncloud selector
    And using <dav_version> DAV path
    When user "Alice" downloads file "textfile.txt" using the WebDAV API
    Then the downloaded content should be "ownCloud test text file"
    Examples:
      | dav_version |
      | old         |
      | new         |


  Scenario Outline: download a file with range
    Given using "ocis" as owncloud selector
    And using <dav_version> DAV path
    When user "Alice" downloads file "textfile.txt" with range "bytes=0-7" using the WebDAV API
    Then the downloaded content should be "ownCloud"
    Examples:
      | dav_version |
      | old         |
      | new         |


  Scenario: Get the size of a file
    Given using "ocis" as owncloud selector
    When user "Alice" gets the size of file "textfile.txt" using the WebDAV API
    Then the HTTP status code should be "207"
    And the size of the file should be "23"


  Scenario Outline: Download a file with comma in the filename
    Given using <dav_version> DAV path
    And user "Alice" has uploaded file with content "file with comma in filename" to <filename>
    And using "ocis" as owncloud selector
    When user "Alice" downloads file <filename> using the WebDAV API
    Then the downloaded content should be "file with comma in filename"
    Examples:
      | dav_version | filename       |
      | old         | "sample,1.txt" |
      | old         | ",,,.txt"      |
      | old         | ",,,.,"        |
      | new         | "sample,1.txt" |
      | new         | ",,,.txt"      |
      | new         | ",,,.,"        |


  Scenario Outline: download a file with single part ranges
    Given using "ocis" as owncloud selector
    And using <dav_version> DAV path
    When user "Alice" downloads file "textfile.txt" with range "bytes=0-7" using the WebDAV API
    Then the HTTP status code should be "206"
    And the following headers should be set
      | header         | value         |
      | Content-Length | 8             |
      | Content-Range  | bytes 0-7/23  |
    And the downloaded content should be "ownCloud"
    Examples:
      | dav_version |
      | old         |
      | new         |


  Scenario Outline: download a file with last byte range out of bounds
    Given using "ocis" as owncloud selector
    And using <dav_version> DAV path
    When user "Alice" downloads file "textfile.txt" with range "bytes=0-24" using the WebDAV API
    Then the HTTP status code should be "206"
    And the downloaded content should be "ownCloud test text file"
    Examples:
      | dav_version |
      | old         |
      | new         |


  Scenario Outline: download a range at the end of a file
    Given using "ocis" as owncloud selector
    And using <dav_version> DAV path
    When user "Alice" downloads file "textfile.txt" with range "bytes=-4" using the WebDAV API
    Then the HTTP status code should be "206"
    And the downloaded content should be "file"
    Examples:
      | dav_version |
      | old         |
      | new         |


  Scenario Outline: download a file with range out of bounds
    Given using "ocis" as owncloud selector
    And using <dav_version> DAV path
    When user "Alice" downloads file "textfile.txt" with range "bytes=24-30" using the WebDAV API
    Then the HTTP status code should be "416"
    Examples:
      | dav_version |
      | old         |
      | new         |


  Scenario Outline: download a hidden file
    Given using <dav_version> DAV path
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has uploaded the following files with content "hidden file"
      | path                |
      | .hidden_file        |
      | FOLDER/.hidden_file |
    And using "ocis" as owncloud selector
    When user "Alice" downloads file ".hidden_file" using the WebDAV API
    Then the HTTP status code should be "200"
    And the downloaded content should be "hidden file"
    When user "Alice" downloads file "FOLDER/.hidden_file" using the WebDAV API
    Then the HTTP status code should be "200"
    And the downloaded content should be "hidden file"
    Examples:
      | dav_version |
      | old         |
      | new         |


  Scenario Outline: Downloading a file should serve security headers
    Given using "ocis" as owncloud selector
    And using <dav_version> DAV path
    When user "Alice" downloads file "textfile.txt" using the WebDAV API
    Then the following headers should be set
      | header                            | value                                                              | 
      | Content-Disposition               | attachment; filename*=UTF-8''textfile.txt; filename="textfile.txt" |
      | Content-Security-Policy           | default-src 'none';                                                |
      | X-Content-Type-Options            | nosniff                                                            |
      | X-Download-Options                | noopen                                                             |
      | X-Frame-Options                   | SAMEORIGIN                                                         |
      | X-Permitted-Cross-Domain-Policies | none                                                               |
      | X-Robots-Tag                      | none                                                               |
      | X-XSS-Protection                  | 1; mode=block                                                      |
    Examples:
      | dav_version |
      | old         |
      | new         |
