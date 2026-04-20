Feature: metadata type search
  As a user
  I want to search files by metadata
  So that I can find the files with specific metadata

  Background:
    Given these users have been created with default attributes:
      | username |
      | Alice    |
      | Brian    |
    And using spaces DAV path


  Scenario: search for files using metadata
    And user "Alice" has uploaded file "filesForUpload/testavatar.jpg" to "/testavatar.jpg"
    And user "Alice" has uploaded file "filesForUpload/testavatar.png" to "/testavatar.png"
    When user "Alice" searches for '<metadata><operation><value>' using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result should contain "1" entries
    And the search result of user "Alice" should contain these entries:
      | testavatar.jpg |
    Examples:
      | metadata                  | value           | operation |
      | photo.cameraMake          | NIKON           | =         |
      | photo.cameraModel         | "COOLPIX P6000" | =         |
      | photo.exposureDenominator |             178 | =         |
      | photo.exposureNumerator   |               1 | =         |
      | photo.fNumber             |             4.5 | =         |
      | photo.focalLength         |               6 | =         |
      | photo.orientation         |               1 | =         |
      | photo.takenDateTime       |      2008-10-22 | =         |
      | location.latitude         |       43.467157 | =         |
      | location.longitude        |       11.885395 | =         |
      | location.altitude         |             100 | =         |


  Scenario Outline: search for files inside sub folders using media type
    Given user "Alice" has created folder "uploadFolder"
    And user "Alice" has uploaded file "filesForUpload/testavatar.jpg" to "uploadFolder/testavatar.jpg"
    And user "Alice" has uploaded file "filesForUpload/testavatar.png" to "uploadFolder/testavatar.png"
    When user "Alice" searches for '<metadata><operation><value>' using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result should contain "1" entries
    And the search result of user "Alice" should contain these entries:
      | testavatar.jpg |
    Examples:
      | metadata                  | value           | operation |
      | photo.cameraMake          | NIKON           | =         |
      | photo.cameraModel         | "COOLPIX P6000" | =         |
      | photo.exposureDenominator |             178 | =         |
      | photo.exposureNumerator   |               1 | =         |
      | photo.fNumber             |             4.5 | =         |
      | photo.focalLength         |               6 | =         |
      | photo.orientation         |               1 | =         |
      | photo.takenDateTime       |      2008-10-22 | =         |
      | location.latitude         |       43.467157 | =         |
      | location.longitude        |       11.885395 | =         |
      | location.altitude         |             100 | =         |


  Scenario Outline: search for file inside project space using media type
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "project101" with the default quota using the Graph API
    And user "Alice" has uploaded a file "filesForUpload/testavatar.jpg" to "/testavatar.jpg" in space "project101"
    And user "Alice" has uploaded a file "filesForUpload/testavatar.png" to "/testavatar.png" in space "project101"
    When user "Alice" searches for '<metadata><operation><value>' using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result should contain "1" entries
    And the search result of user "Alice" should contain these entries:
      | testavatar.jpg |
    Examples:
      | metadata                  | value           | operation |
      | photo.cameraMake          | NIKON           | =         |
      | photo.cameraModel         | "COOLPIX P6000" | =         |
      | photo.exposureDenominator |             178 | =         |
      | photo.exposureNumerator   |               1 | =         |
      | photo.fNumber             |             4.5 | =         |
      | photo.focalLength         |               6 | =         |
      | photo.orientation         |               1 | =         |
      | photo.takenDateTime       |      2008-10-22 | =         |
      | location.latitude         |       43.467157 | =         |
      | location.longitude        |       11.885395 | =         |
      | location.altitude         |             100 | =         |


  Scenario Outline: sharee searches for shared files using media type
    Given user "Alice" has created folder "uploadFolder"
    And user "Alice" has uploaded file "filesForUpload/testavatar.jpg" to "uploadFolder/testavatar.jpg"
    And user "Alice" has uploaded file "filesForUpload/testavatar.png" to "uploadFolder/testavatar.png"
    And user "Alice" has sent the following resource share invitation:
      | resource        | uploadFolder |
      | space           | Personal     |
      | sharee          | Brian        |
      | shareType       | user         |
      | permissionsRole | Viewer       |
    And user "Brian" has a share "uploadFolder" synced
    When user "Alice" searches for '<metadata><operation><value>' using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result should contain "1" entries
    And the search result of user "Alice" should contain these entries:
      | testavatar.jpg |
    Examples:
      | metadata                  | value           | operation |
      | photo.cameraMake          | NIKON           | =         |
      | photo.cameraModel         | "COOLPIX P6000" | =         |
      | photo.exposureDenominator |             178 | =         |
      | photo.exposureNumerator   |               1 | =         |
      | photo.fNumber             |             4.5 | =         |
      | photo.focalLength         |               6 | =         |
      | photo.orientation         |               1 | =         |
      | photo.takenDateTime       |      2008-10-22 | =         |
      | location.latitude         |       43.467157 | =         |
      | location.longitude        |       11.885395 | =         |
      | location.altitude         |             100 | =         |


  Scenario Outline: space viewer searches for files using mediatype filter
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "project101" with the default quota using the Graph API
    And user "Alice" has uploaded a file "filesForUpload/testavatar.jpg" to "/testavatar.jpg" in space "project101"
    And user "Alice" has uploaded a file "filesForUpload/testavatar.png" to "/testavatar.png" in space "project101"
    And user "Alice" has sent the following space share invitation:
      | space           | project101   |
      | sharee          | Brian        |
      | shareType       | user         |
      | permissionsRole | Space Viewer |
    When user "Brian" searches for '<metadata><operation><value>' using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result should contain "1" entries
    And the search result of user "Alice" should contain these entries:
      | testavatar.jpg |
    Examples:
      | metadata                  | value           | operation |
      | photo.cameraMake          | NIKON           | =         |
      | photo.cameraModel         | "COOLPIX P6000" | =         |
      | photo.exposureDenominator |             178 | =         |
      | photo.exposureNumerator   |               1 | =         |
      | photo.fNumber             |             4.5 | =         |
      | photo.focalLength         |               6 | =         |
      | photo.orientation         |               1 | =         |
      | photo.takenDateTime       |      2008-10-22 | =         |
      | location.latitude         |       43.467157 | =         |
      | location.longitude        |       11.885395 | =         |
      | location.altitude         |             100 | =         |

