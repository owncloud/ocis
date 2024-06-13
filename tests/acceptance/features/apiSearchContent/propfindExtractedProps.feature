@tikaServiceNeeded
Feature: propfind extracted props
  As a user
  I want to get extracted properties of resource
  So that I can make sure that the response contains audio, location, image and photo properties


  Scenario: check extracted properties of a file from project space
    Given user "Alice" has been created with default attributes and without skeleton files
    And using spaces DAV path
    And the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "new-space" with the default quota using the Graph API
    And user "Alice" has uploaded a file "filesForUpload/testaudio.mp3" to "testaudio.mp3" in space "new-space"
    And user "Alice" has uploaded a file "filesForUpload/testavatar.jpg" to "testavatar.jpg" in space "new-space"
    When user "Alice" gets the following extracted properties of resource "testaudio.mp3" inside space "new-space" using the WebDAV API
      | propertyName |
      | oc:audio     |
    Then the HTTP status code should be "207"
    And the "PROPFIND" response should contain a space "new-space" with these key and value pairs:
      | key                | value                          |
      | oc:audio/oc:album  | ALBUM1234567890123456789012345 |
      | oc:audio/oc:artist | ARTIST123456789012345678901234 |
      | oc:audio/oc:genre  | Pop                            |
      | oc:audio/oc:title  | TITLE1234567890123456789012345 |
      | oc:audio/oc:track  | 1                              |
      | oc:audio/oc:year   | 2001                           |
    When user "Alice" gets the following extracted properties of resource "testavatar.jpg" inside space "new-space" using the WebDAV API
      | propertyName |
      | oc:image     |
      | oc:location  |
      | oc:photo     |
    Then the HTTP status code should be "207"
    And the "PROPFIND" response should contain a space "new-space" with these key and value pairs:
      | key                              | value                |
      | oc:image/oc:width                | 640                  |
      | oc:image/oc:height               | 480                  |
      | oc:location/oc:latitude          | 43.467157            |
      | oc:location/oc:longitude         | 11.885395            |
      | oc:photo/oc:camera-make          | NIKON                |
      | oc:photo/oc:camera-model         | COOLPIX P6000        |
      | oc:photo/oc:f-number             | 4.5                  |
      | oc:photo/oc:focal-length         | 6                    |
