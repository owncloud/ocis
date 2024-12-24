Feature: a


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: uploading to a user shared folder with read/write permission when the sharer has insufficient quota does not work
    Given user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has sent the following resource share invitation:
      | resource        | FOLDER   |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    And user "Brian" has a share "FOLDER" synced
    And user "Admin" has changed the quota of the personal space of "Alice Hansen" space to "1"
    When user "Brian" uploads file "filesForUpload/textfile.txt" to "/Shares/FOLDER/myfile.txt" using the WebDAV API
    Then the HTTP status code should be "507"
    And as "Alice" file "/FOLDER/myfile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: retrieving folder quota when no quota is set
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes
    And using <dav-path-version> DAV path
    When user "Admin" changes the quota of the "Alice Hansen" space to "0"
    Then the HTTP status code should be "200"
    And as user "Alice" folder "/" should contain a property "d:quota-available-bytes" with value "0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |
