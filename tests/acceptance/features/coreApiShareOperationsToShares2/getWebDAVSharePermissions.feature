@skipOnReva
Feature: sharing
  As a user
  I want to check the webdav share permissions
  So that I know the resources have proper permissions

  Background:
    Given using OCS API version "1"
    And these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |

  @smokeTest
  Scenario Outline: check webdav share-permissions for owned file
    Given using <dav-path-version> DAV path
    And user "Alice" has uploaded file with content "foo" to "/tmp.txt"
    When user "Alice" gets the following properties of file "/tmp.txt" using the WebDAV API
      | propertyName          |
      | ocs:share-permissions |
    Then the HTTP status code should be "207"
    And the single response should contain a property "ocs:share-permissions" with value "19"
    Examples:
      | dav-path-version |
      | old              |
      | new              |

    @skipOnRevaMaster
    Examples:
      | dav-path-version |
      | spaces           |

  @skipOnReva
  Scenario Outline: check webdav share-permissions for received file with edit
    Given using <dav-path-version> DAV path
    And user "Alice" has uploaded file with content "foo" to "/tmp.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | tmp.txt     |
      | space           | Personal    |
      | sharee          | Brian       |
      | shareType       | user        |
      | permissionsRole | File Editor |
    When user "Brian" gets the following properties of file "/Shares/tmp.txt" using the WebDAV API
      | propertyName          |
      | ocs:share-permissions |
    Then the HTTP status code should be "207"
    And the single response should contain a property "ocs:share-permissions" with value "3"
    Examples:
      | dav-path-version |
      | old              |
      | new              |


  Scenario Outline: check webdav share-permissions for received group shared file with edit
    Given using <dav-path-version> DAV path
    And group "grp1" has been created
    And user "Brian" has been added to group "grp1"
    And user "Alice" has uploaded file with content "foo" to "/tmp.txt"
    And user "Alice" has created a share with settings
      | path        | /tmp.txt    |
      | shareType   | group       |
      | permissions | update,read |
      | shareWith   | grp1        |
    When user "Brian" gets the following properties of file "/Shares/tmp.txt" using the WebDAV API
      | propertyName          |
      | ocs:share-permissions |
    Then the HTTP status code should be "207"
    And the single response should contain a property "ocs:share-permissions" with value "3"
    Examples:
      | dav-path-version |
      | old              |
      | new              |

  @skipOnReva @issue-2213
  Scenario Outline: check webdav share-permissions for received file with edit permissions but no reshare permissions
    Given using <dav-path-version> DAV path
    And user "Alice" has uploaded file with content "foo" to "/tmp.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | tmp.txt  |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Viewer   |
    And using SharingNG
    When user "Alice" updates the last share using the sharing API with
      | permissions | update,read |
    Then the HTTP status code should be "200"
    And as user "Brian" file "/Shares/tmp.txt" should contain a property "ocs:share-permissions" with value "3"
    Examples:
      | dav-path-version |
      | old              |
      | new              |

  @issue-2213
  Scenario Outline: check webdav share-permissions for received group shared file with edit permissions but no reshare permissions
    Given using <dav-path-version> DAV path
    And group "grp1" has been created
    And user "Brian" has been added to group "grp1"
    And user "Alice" has uploaded file with content "foo" to "/tmp.txt"
    And user "Alice" has created a share with settings
      | path        | /tmp.txt    |
      | shareType   | group       |
      | permissions | update,read |
      | shareWith   | grp1        |
    When user "Brian" gets the following properties of file "/Shares/tmp.txt" using the WebDAV API
      | propertyName          |
      | ocs:share-permissions |
    Then the HTTP status code should be "207"
    And the single response should contain a property "ocs:share-permissions" with value "3"
    Examples:
      | dav-path-version |
      | old              |
      | new              |

  @skipOnReva @issue-2213
  Scenario Outline: check webdav share-permissions for received file without edit permissions
    Given using <dav-path-version> DAV path
    And user "Alice" has uploaded file with content "foo" to "/tmp.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | tmp.txt  |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Viewer   |
    And using SharingNG
    When user "Alice" updates the last share using the sharing API with
      | permissions | read |
    Then the HTTP status code should be "200"
    And as user "Brian" file "/Shares/tmp.txt" should contain a property "ocs:share-permissions" with value "1"
    Examples:
      | dav-path-version |
      | old              |
      | new              |


  Scenario Outline: check webdav share-permissions for received group shared file with reshare permissions but no edit permissions
    Given using <dav-path-version> DAV path
    And group "grp1" has been created
    And user "Brian" has been added to group "grp1"
    And user "Alice" has uploaded file with content "foo" to "/tmp.txt"
    And user "Alice" has created a share with settings
      | path        | /tmp.txt |
      | shareType   | group    |
      | permissions | read     |
      | shareWith   | grp1     |
    When user "Brian" gets the following properties of file "/Shares/tmp.txt" using the WebDAV API
      | propertyName          |
      | ocs:share-permissions |
    Then the HTTP status code should be "207"
    And the single response should contain a property "ocs:share-permissions" with value "1"
    Examples:
      | dav-path-version |
      | old              |
      | new              |


  Scenario Outline: check webdav share-permissions for owned folder
    Given using <dav-path-version> DAV path
    And user "Alice" has created folder "/tmp"
    When user "Alice" gets the following properties of folder "/" using the WebDAV API
      | propertyName          |
      | ocs:share-permissions |
    Then the HTTP status code should be "207"
    And the single response should contain a property "ocs:share-permissions" with value "31"
    Examples:
      | dav-path-version |
      | old              |
      | new              |

    @skipOnRevaMaster
    Examples:
      | dav-path-version |
      | spaces           |

  @skipOnReva
  Scenario Outline: check webdav share-permissions for received folder with all permissions
    Given using <dav-path-version> DAV path
    And user "Alice" has created folder "/tmp"
    And user "Alice" has sent the following resource share invitation:
      | resource        | tmp      |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Editor   |
    When user "Brian" gets the following properties of folder "/Shares/tmp" using the WebDAV API
      | propertyName          |
      | ocs:share-permissions |
    Then the HTTP status code should be "207"
    And the single response should contain a property "ocs:share-permissions" with value "15"
    Examples:
      | dav-path-version |
      | old              |
      | new              |


  Scenario Outline: check webdav share-permissions for received group shared folder with all permissions
    Given using <dav-path-version> DAV path
    And group "grp1" has been created
    And user "Brian" has been added to group "grp1"
    And user "Alice" has created folder "/tmp"
    And user "Alice" has created a share with settings
      | path      | tmp   |
      | shareType | group |
      | shareWith | grp1  |
    When user "Brian" gets the following properties of folder "/Shares/tmp" using the WebDAV API
      | propertyName          |
      | ocs:share-permissions |
    Then the HTTP status code should be "207"
    And the single response should contain a property "ocs:share-permissions" with value "15"
    Examples:
      | dav-path-version |
      | old              |
      | new              |

  @skipOnReva @issue-2213
  Scenario Outline: check webdav share-permissions for received folder with all permissions but edit
    Given using <dav-path-version> DAV path
    And user "Alice" has created folder "/tmp"
    And user "Alice" has sent the following resource share invitation:
      | resource        | tmp      |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Viewer   |
    And using SharingNG
    When user "Alice" updates the last share using the sharing API with
      | permissions | delete,create,read |
    Then the HTTP status code should be "200"
    And as user "Brian" folder "/Shares/tmp" should contain a property "ocs:share-permissions" with value "13"
    Examples:
      | dav-path-version |
      | old              |
      | new              |


  Scenario Outline: check webdav share-permissions for received group shared folder with all permissions but edit
    Given using <dav-path-version> DAV path
    And group "grp1" has been created
    And user "Brian" has been added to group "grp1"
    And user "Alice" has created folder "/tmp"
    And user "Alice" has created a share with settings
      | path        | tmp                |
      | shareType   | group              |
      | shareWith   | grp1               |
      | permissions | delete,create,read |
    When user "Brian" gets the following properties of folder "/Shares/tmp" using the WebDAV API
      | propertyName          |
      | ocs:share-permissions |
    Then the HTTP status code should be "207"
    And the single response should contain a property "ocs:share-permissions" with value "13"
    Examples:
      | dav-path-version |
      | old              |
      | new              |

  @skipOnReva
  Scenario Outline: check webdav share-permissions for received folder with all permissions but create
    Given using <dav-path-version> DAV path
    And user "Alice" has created folder "/tmp"
    And user "Alice" has sent the following resource share invitation:
      | resource        | tmp      |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Viewer   |
    And using SharingNG
    When user "Alice" updates the last share using the sharing API with
      | permissions | delete,update,read |
    Then the HTTP status code should be "200"
    And as user "Brian" folder "/Shares/tmp" should contain a property "ocs:share-permissions" with value "11"
    Examples:
      | dav-path-version |
      | old              |
      | new              |


  Scenario Outline: check webdav share-permissions for received group shared folder with all permissions but create
    Given using <dav-path-version> DAV path
    And group "grp1" has been created
    And user "Brian" has been added to group "grp1"
    And user "Alice" has created folder "/tmp"
    And user "Alice" has created a share with settings
      | path        | tmp                |
      | shareType   | group              |
      | shareWith   | grp1               |
      | permissions | delete,update,read |
    When user "Brian" gets the following properties of folder "/Shares/tmp" using the WebDAV API
      | propertyName          |
      | ocs:share-permissions |
    Then the HTTP status code should be "207"
    And the single response should contain a property "ocs:share-permissions" with value "11"
    Examples:
      | dav-path-version |
      | old              |
      | new              |

  @skipOnReva
  Scenario Outline: check webdav share-permissions for received folder with all permissions but delete
    Given using <dav-path-version> DAV path
    And user "Alice" has created folder "/tmp"
    And user "Alice" has sent the following resource share invitation:
      | resource        | tmp      |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Viewer   |
    And using SharingNG
    When user "Alice" updates the last share using the sharing API with
      | permissions | create,update,read |
    Then the HTTP status code should be "200"
    And as user "Brian" folder "/Shares/tmp" should contain a property "ocs:share-permissions" with value "7"
    Examples:
      | dav-path-version |
      | old              |
      | new              |


  Scenario Outline: check webdav share-permissions for received group shared folder with all permissions but delete
    Given using <dav-path-version> DAV path
    And group "grp1" has been created
    And user "Brian" has been added to group "grp1"
    And user "Alice" has created folder "/tmp"
    And user "Alice" has created a share with settings
      | path        | tmp                |
      | shareType   | group              |
      | shareWith   | grp1               |
      | permissions | create,update,read |
    When user "Brian" gets the following properties of folder "/Shares/tmp" using the WebDAV API
      | propertyName          |
      | ocs:share-permissions |
    Then the HTTP status code should be "207"
    And the single response should contain a property "ocs:share-permissions" with value "7"
    Examples:
      | dav-path-version |
      | old              |
      | new              |

  @skipOnReva
  Scenario Outline: check webdav share-permissions for received folder with all permissions but share
    Given using <dav-path-version> DAV path
    And user "Alice" has created folder "/tmp"
    And user "Alice" has sent the following resource share invitation:
      | resource        | tmp      |
      | space           | Personal |
      | sharee          | Brian    |
      | shareType       | user     |
      | permissionsRole | Viewer   |
    And using SharingNG
    When user "Alice" updates the last share using the sharing API with
      | permissions | change |
    Then the HTTP status code should be "200"
    And as user "Brian" folder "/Shares/tmp" should contain a property "ocs:share-permissions" with value "15"
    Examples:
      | dav-path-version |
      | old              |
      | new              |


  Scenario Outline: check webdav share-permissions for received group shared folder with all permissions but share
    Given using <dav-path-version> DAV path
    And group "grp1" has been created
    And user "Brian" has been added to group "grp1"
    And user "Alice" has created folder "/tmp"
    And user "Alice" has created a share with settings
      | path        | tmp    |
      | shareType   | group  |
      | shareWith   | grp1   |
      | permissions | change |
    When user "Brian" gets the following properties of folder "/Shares/tmp" using the WebDAV API
      | propertyName          |
      | ocs:share-permissions |
    Then the HTTP status code should be "207"
    And the single response should contain a property "ocs:share-permissions" with value "15"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
