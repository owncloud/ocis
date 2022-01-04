# Sharing tests currently doesn't work
# Accessing oc10 shares from ocis still WIP in PR #2232
# https://github.com/owncloud/ocis/pull/2232
@api
Feature: sharing files and folders
  As a user
  I want to share files/folders with other users
  So that I can give access to my files/folders to others


  Background:
    Given using "oc10" as owncloud selector
    And the administrator has set the default folder for received shares to "Shares"
    And auto-accept shares has been disabled
    And using OCS API version "1"
    And using new DAV path
    And user "Alice" has been created with default attributes and without skeleton files
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file with content "ownCloud test text file" to "textfile.txt"


  Scenario: accept a pending share
    Given user "Alice" has shared folder "/textfile.txt" with user "Brian"
    And using "ocis" as owncloud selector
    When user "Brian" accepts share "/textfile.txt" offered by user "Alice" using the sharing API
    Then the OCS status code should be "100"
    And the HTTP status code should be "200"
    And the sharing API should report to user "Brian" that these shares are in the accepted state
      | path                 |
      | /Shares/textfile.txt |