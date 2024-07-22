Feature: sharing files and folders
  As a user
  I want to share files/folders with other users
  So that I can give access to my files/folders to others


  Background:
    Given using "oc10" as owncloud selector
    And using OCS API version "1"
    And using new DAV path
    And user "Alice" has been created with default attributes and without skeleton files
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file with content "ownCloud test text file" to "textfile.txt"


  Scenario: accept a pending share
    Given user "Alice" has sent the following resource share invitation:
      | resource        | textfile.txt |
      | space           | Personal     |
      | sharee          | Brian        |
      | shareType       | user         |
      | permissionsRole | Editor       |
    And user "Brian" has a share "textfile.txt" synced
    And using "ocis" as owncloud selector
    And the sharing API should report to user "Brian" that these shares are in the accepted state
      | path                 |
      | /Shares/textfile.txt |
