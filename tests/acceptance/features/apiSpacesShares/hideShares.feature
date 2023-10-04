Feature: Share a file or folder that is inside a space
  As a user with manager space role
  I want to be able to share the data inside the space
  So that other users can have access to it

  Note - this feature is run in CI with ACCOUNTS_HASH_DIFFICULTY set to the default for production
  See https://github.com/owncloud/ocis/issues/1542 and https://github.com/owncloud/ocis/pull/839

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
    And user "Alice" has created a folder "folder" in space "Alice Hansen"

  
  Scenario: user hides accepted share
    Given user "Alice" has created a share inside of space "Alice Hansen" with settings:
      | path       | folder                   |
      | shareWith  | Brian                    |
      | role       | viewer                   |
    And user "Brian" has accepted share "/folder" offered by user "Alice"
    When user "Brian" hiddes share "/Shares/folder" of the accepted state offered by user "Alice" using the sharing API
    Then the HTTP status code should be "200"
    

  Scenario: user hides pending share
    Given user "Alice" has created a share inside of space "Alice Hansen" with settings:
      | path       | folder                   |
      | shareWith  | Brian                    |
      | role       | viewer                   |
    When user "Brian" hiddes share "/folder" of the pending state offered by user "Alice" using the sharing API
    Then the HTTP status code should be "200"

  
  Scenario: user hides declined share
    Given user "Alice" has created a share inside of space "Alice Hansen" with settings:
      | path       | folder                   |
      | shareWith  | Brian                    |
      | role       | viewer                   |
    And user "Brian" has declined share "/folder" offered by user "Alice"
    When user "Brian" hiddes share "/folder" of the declined state offered by user "Alice" using the sharing API
    Then the HTTP status code should be "200"