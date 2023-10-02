@env-config
Feature: PROPFIND with depth:infinity
  As a user
  I want to retrieve all properties of a resource
  So that I can get the information about a resource

  Background:
    Given the config "OCDAV_ALLOW_PROPFIND_DEPTH_INFINITY" has been set to "true"
    And user "Alice" has been created with default attributes and without skeleton files
    And user "Alice" has created the following folders
      | path                                        |
      | simple-folder                               |
      | simple-folder/simple-folder1                |
      | simple-folder/simple-empty-folder           |
      | simple-folder/simple-folder1/simple-folder2 |
    And user "Alice" has uploaded the following files with content "simple-test-content"
      | path                                                      |
      | textfile0.txt                                             |
      | welcome.txt                                               |
      | simple-folder/textfile0.txt                               |
      | simple-folder/welcome.txt                                 |
      | simple-folder/simple-folder1/textfile0.txt                |
      | simple-folder/simple-folder1/welcome.txt                  |
      | simple-folder/simple-folder1/simple-folder2/textfile0.txt |
      | simple-folder/simple-folder1/simple-folder2/welcome.txt   |


  Scenario Outline: get the list of resources with depth infinity
    Given using <dav-path-version> DAV path
    When user "Alice" lists the resources in "/" with depth "infinity" using the WebDAV API
    Then the HTTP status code should be "207"
    And the last DAV response for user "Alice" should contain these nodes
      | name                                                      |
      | textfile0.txt                                             |
      | welcome.txt                                               |
      | simple-folder/                                            |
      | simple-folder/textfile0.txt                               |
      | simple-folder/welcome.txt                                 |
      | simple-folder/simple-empty-folder/                        |
      | simple-folder/simple-folder1/                             |
      | simple-folder/simple-folder1/simple-folder2               |
      | simple-folder/simple-folder1/textfile0.txt                |
      | simple-folder/simple-folder1/welcome.txt                  |
      | simple-folder/simple-folder1/simple-folder2/textfile0.txt |
      | simple-folder/simple-folder1/simple-folder2/welcome.txt   |
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: get the list of resources of a folder with depth infinity
    Given using <dav-path-version> DAV path
    When user "Alice" lists the resources in "/simple-folder" with depth "infinity" using the WebDAV API
    Then the HTTP status code should be "207"
    And the last DAV response for user "Alice" should contain these nodes
      | name                                                      |
      | /simple-folder/textfile0.txt                              |
      | /simple-folder/welcome.txt                                |
      | /simple-folder/simple-folder1/                            |
      | simple-folder/simple-folder1/simple-folder2               |
      | simple-folder/simple-folder1/textfile0.txt                |
      | simple-folder/simple-folder1/welcome.txt                  |
      | simple-folder/simple-folder1/simple-folder2/textfile0.txt |
      | simple-folder/simple-folder1/simple-folder2/welcome.txt   |
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario: get the list of resources in a folder shared through public link with depth infinity
    Given using new DAV path
    And user "Alice" has created a public link share of folder "simple-folder"
    When the public lists the resources in the last created public link with depth "infinity" using the WebDAV API
    Then the HTTP status code should be "207"
    And the last public link DAV response should contain these nodes
      | name                                                         |
      | /textfile0.txt                                               |
      | /welcome.txt                                                 |
      | /simple-folder1/                                             |
      | /simple-folder1/welcome.txt                                  |
      | /simple-folder1/simple-folder2                               |
      | /simple-folder1/textfile0.txt                                |
      | /simple-folder1/simple-folder2/textfile0.txt                 |
      | /simple-folder1/simple-folder2/welcome.txt                   |


  Scenario Outline: get the list of files in the trashbin with depth infinity
    Given using <dav-path-version> DAV path
    And user "Alice" has deleted the following resources
      | path           |
      | textfile0.txt  |
      | welcome.txt    |
      | simple-folder/ |
    When user "Alice" lists the resources in the trashbin with depth "infinity" using the WebDAV API
    Then the HTTP status code should be "207"
    And the trashbin DAV response should contain these nodes
      | name                                                      |
      | textfile0.txt                                             |
      | welcome.txt                                               |
      | simple-folder/                                            |
      | simple-folder/textfile0.txt                               |
      | simple-folder/welcome.txt                                 |
      | simple-folder/simple-folder1/textfile0.txt                |
      | simple-folder/simple-folder1/welcome.txt                  |
      | simple-folder/simple-folder1/simple-folder2/textfile0.txt |
      | simple-folder/simple-folder1/simple-folder2/welcome.txt   |
    Examples:
      | dav-path-version |
      | new              |
      | spaces           |


  Scenario: get the list of resources in a folder shared through public link with depth infinity when depth infinity is not allowed
    Given the config "OCDAV_ALLOW_PROPFIND_DEPTH_INFINITY" has been set to "false"
    And using new DAV path
    And user "Alice" has created a public link share of folder "simple-folder"
    When the public lists the resources in the last created public link with depth "infinity" using the WebDAV API
    Then the HTTP status code should be "400"


  Scenario Outline: get the list of files in the trashbin with depth infinity when depth infinity is not allowed
    Given the config "OCDAV_ALLOW_PROPFIND_DEPTH_INFINITY" has been set to "false"
    And using <dav-path-version> DAV path
    And user "Alice" has deleted the following resources
      | path           |
      | textfile0.txt  |
      | welcome.txt    |
      | simple-folder/ |
    When user "Alice" lists the resources in the trashbin with depth "infinity" using the WebDAV API
    Then the HTTP status code should be "400"
    Examples:
      | dav-path-version |
      | new              |
      | spaces           |