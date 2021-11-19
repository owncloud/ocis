
# Sharing tests currently doesn't work
# Accessing oc10 shares from ocis still WIP in PR #2232
# https://github.com/owncloud/ocis/pull/2232

Feature: sharing files and folders
  As a user
  I want to share files/folders with other users
  So that I can give access to my files/folders to others


  Background:
    Given using "oc10" as owncloud selector
    And user "Alice" has been created with default attributes and without skeleton files
    And user "Brian" has been created with default attributes and without skeleton files


  Scenario: share a file with a user
    Given user "Alice" has created folder "PARENT"
    And user "Alice" has shared folder "PARENT" with user "Brian"
    Then the HTTP status code should be "200"
    When using "ocis" as owncloud selector
    # And user "Brian" accepts share "PARENT" offered by user "Alice" using the sharing API
    # Then the HTTP status code should be "200"