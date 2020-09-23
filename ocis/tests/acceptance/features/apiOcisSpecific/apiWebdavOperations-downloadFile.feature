@api
Feature: download file
  As a user
  I want to be able to download files
  So that I can work wih local copies of files on my client system

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file with content "ownCloud test text file 0" to "/textfile0.txt"
    And user "Alice" has uploaded file with content "Welcome this is just an example file for developers." to "/welcome.txt"

  @skipOnOcis-OC-Storage @issue-ocis-reva-98
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario Outline: Get the content-length response header of a pdf file
    Given using <dav_version> DAV path
    And user "Alice" has uploaded file "filesForUpload/simple.pdf" to "/simple.pdf"
    When user "Alice" downloads file "/simple.pdf" using the WebDAV API
    And the following headers should not be set
      | header                |
      | OC-JobStatus-Location |
    Examples:
      | dav_version |
      | old         |
      | new         |
