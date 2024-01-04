@issue-1313 @skipOnReva
Feature: get quota
  As a user
  I want to be able to find out my available storage quota
  So that I can manage the use of my allocated storage

  Background:
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes and without skeleton files


  Scenario Outline: retrieving folder quota when no quota is set
    Given using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |

  @smokeTest
  Scenario Outline: retrieving folder quota when quota is set
    Given using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "10000"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "10000"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota of shared folder with quota when no quota is set for recipient
    Given using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "0"
    And user "Admin" has changed the quota of the personal space of "Brian Murphy" space to "10000"
    And user "Brian" has created folder "/testquota"
    And user "Brian" has created a share with settings
      | path        | testquota |
      | shareType   | user      |
      | permissions | all       |
      | shareWith   | Alice     |
    When user "Alice" gets the following properties of folder "<folder-path>" inside space "Shares" using the WebDAV API
      | propertyName            |
      | d:quota-available-bytes |
    Then the HTTP status code should be "207"
    And the single response should contain a property "d:quota-available-bytes" with value "10000"
    Examples:
      | dav-path-version | folder-path       |
      | old              | /Shares/testquota |
      | new              | /Shares/testquota |
      | spaces           | /testquota        |


  Scenario Outline: retrieving folder quota when quota is set and a file was uploaded
    Given using <dav-path-version> DAV path
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "10000"
    And user "Alice" has uploaded file "/prueba.txt" of size 93 bytes
    When user "Alice" gets the following properties of folder "/" using the WebDAV API
      | propertyName            |
      | d:quota-available-bytes |
    Then the HTTP status code should be "207"
    And the single response should contain a property "d:quota-available-bytes" with value "10000"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when quota is set and a file was received
    Given using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Admin" has changed the quota of the personal space of "Brian Murphy" space to "10000"
    And user "Alice" has uploaded file "/Alice.txt" of size 93 bytes
    And user "Alice" has shared file "Alice.txt" with user "Brian"
    When user "Brian" gets the following properties of folder "/" using the WebDAV API
      | propertyName            |
      | d:quota-available-bytes |
    Then the HTTP status code should be "207"
    And the single response should contain a property "d:quota-available-bytes" with value "10000"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |
