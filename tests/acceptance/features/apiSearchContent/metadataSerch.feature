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

  @issue-12230
  Scenario: search for files using metadata
    Given user "Alice" has uploaded file "filesForUpload/testavatar.jpg" to "/testavatar.jpg"
    When user "Alice" searches for '<query>' using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result should contain "1" entries
    And the search result of user "Alice" should contain these entries:
      | testavatar.jpg |
    Examples:
      | query                             |
      | photo.cameraMake=NIKON            |
      | photo.cameraModel="COOLPIX P6000" |
      | photo.exposureDenominator=178     |
      | photo.exposureNumerator=1         |
      | photo.fNumber=4.5                 |
      | photo.focalLength=6               |
      | photo.orientation=1               |
      | photo.takenDateTime=2008-10-22    |
      | location.latitude=43.467157       |
      | location.longitude=11.885395      |
      | location.altitude=100             |

  @issue-12230
  Scenario Outline: search for files inside sub folders using media type
    Given user "Alice" has created folder "uploadFolder"
    And user "Alice" has uploaded file "filesForUpload/testavatar.jpg" to "uploadFolder/testavatar.jpg"
    When user "Alice" searches for '<query>' using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result should contain "1" entries
    And the search result of user "Alice" should contain these entries:
      | uploadFolder/testavatar.jpg |
    Examples:
      | query                             |
      | photo.cameraMake=NIKON            |
      | photo.cameraModel="COOLPIX P6000" |
      | photo.exposureDenominator=178     |
      | photo.exposureNumerator=1         |
      | photo.fNumber=4.5                 |
      | photo.focalLength=6               |
      | photo.orientation=1               |
      | photo.takenDateTime=2008-10-22    |
      | location.latitude=43.467157       |
      | location.longitude=11.885395      |
      | location.altitude=100             |

  @issue-12230
  Scenario Outline: search for file inside project space using media type
    Given the administrator has assigned the role "Space Admin" to user "Alice" using the Graph API
    And user "Alice" has created a space "project101" with the default quota using the Graph API
    And user "Alice" has uploaded a file "filesForUpload/testavatar.jpg" to "/testavatar.jpg" in space "project101"
    And user "Alice" has uploaded a file "filesForUpload/testavatar.png" to "/testavatar.png" in space "project101"
    When user "Alice" searches for '<query>' using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result should contain "1" entries
    And the search result of user "Alice" should contain these entries:
      | testavatar.jpg |
    Examples:
      | query                             |
      | photo.cameraMake=NIKON            |
      | photo.cameraModel="COOLPIX P6000" |
      | photo.exposureDenominator=178     |
      | photo.exposureNumerator=1         |
      | photo.fNumber=4.5                 |
      | photo.focalLength=6               |
      | photo.orientation=1               |
      | photo.takenDateTime=2008-10-22    |
      | location.latitude=43.467157       |
      | location.longitude=11.885395      |
      | location.altitude=100             |

  @issue-12230
  Scenario Outline: sharee searches for shared files using media type
    Given user "Alice" has created folder "uploadFolder"
    And user "Alice" has uploaded file "filesForUpload/testavatar.jpg" to "uploadFolder/testavatar.jpg"
    And user "Alice" has sent the following resource share invitation:
      | resource        | uploadFolder |
      | space           | Personal     |
      | sharee          | Brian        |
      | shareType       | user         |
      | permissionsRole | Viewer       |
    And user "Brian" has a share "uploadFolder" synced
    When user "Brian" searches for '<query>' using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result should contain "1" entries
    And the search result of user "Brian" should contain these entries:
      | uploadFolder/testavatar.jpg |
    Examples:
      | query                             |
      | photo.cameraMake=NIKON            |
      | photo.cameraModel="COOLPIX P6000" |
      | photo.exposureDenominator=178     |
      | photo.exposureNumerator=1         |
      | photo.fNumber=4.5                 |
      | photo.focalLength=6               |
      | photo.orientation=1               |
      | photo.takenDateTime=2008-10-22    |
      | location.latitude=43.467157       |
      | location.longitude=11.885395      |
      | location.altitude=100             |

  @issue-12230
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
    When user "Brian" searches for '<query>' using the WebDAV API
    Then the HTTP status code should be "207"
    And the search result should contain "<expectedCount>" entries
    And the search result of user "Brian" should contain these entries:
      | <expectedFile1> |
      | <expectedFile2> |
    # Examples:
    #   | query                             |
    #   | photo.cameraMake=NIKON            |
    #   | photo.cameraModel="COOLPIX P6000" |
    #   | photo.exposureDenominator=178     |
    #   | photo.exposureNumerator=1         |
    #   | photo.fNumber=4.5                 |
    #   | photo.focalLength=6               |
    #   | photo.orientation=1               |
    #   | photo.takenDateTime=2008-10-22    |
    #   | location.latitude=43.467157       |
    #   | location.longitude=11.885395      |
    #   | location.altitude=100             |

    Examples:
      | query                                  | expectedCount | expectedFile1  | expectedFile2  |
      | photo.cameraMake=NIKON                 | 1             | testavatar.jpg |                |
      | photo.cameraModel="COOLPIX P6000"      | 1             | testavatar.jpg |                |
      | photo.exposureDenominator=178          | 1             | testavatar.jpg |                |
      | photo.exposureNumerator=1              | 1             | testavatar.jpg |                |
      | photo.fNumber=4.5                      | 1             | testavatar.jpg |                |
      | photo.focalLength=6                    | 1             | testavatar.jpg |                |
      | photo.orientation=1                    | 1             | testavatar.jpg |                |
      | photo.takenDateTime=2008-10-22         | 1             | testavatar.jpg |                |
      | location.latitude=43.467157            | 1             | testavatar.jpg |                |
      | location.longitude=11.885395           | 1             | testavatar.jpg |                |
      | location.altitude=100                  | 1             | testavatar.jpg |                |
      | NOT photo.cameraMake=CANON             | 2             | testavatar.jpg | testavatar.png |
      | NOT photo.exposureDenominator=180      | 2             | testavatar.jpg | testavatar.png |
      | NOT photo.orientation=2                | 2             | testavatar.jpg | testavatar.png |
      | NOT photo.takenDateTime=2007-01-01     | 2             | testavatar.jpg | testavatar.png |
      | photo.exposureDenominator<180          | 1             | testavatar.jpg |                |
      | photo.exposureDenominator<=178         | 1             | testavatar.jpg |                |
      | photo.fNumber>4                        | 1             | testavatar.jpg |                |
      | location.altitude<150                  | 1             | testavatar.jpg |                |
      | location.altitude>=100                 | 1             | testavatar.jpg |                |
      | photo.focalLength<10                   | 1             | testavatar.jpg |                |
      | photo.focalLength>4                    | 1             | testavatar.jpg |                |
      | photo.focalLength>=6                   | 1             | testavatar.jpg |                |