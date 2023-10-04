Feature: PROPFIND
  As a user
  I want to retrieve all properties of a resource
  So that I can get the information about a resource

  @issue-751
  Scenario Outline: send PROPFIND request to "/remote.php/dav/(files|spaces)"
    Given user "Alice" has been created with default attributes and without skeleton files
    When user "Alice" requests "<dav_path>" with "PROPFIND" using basic auth
    Then the HTTP status code should be "405"
    Examples:
      | dav_path              |
      | /remote.php/dav/files |

    @skipOnRevaMaster
    Examples:
      | dav_path               |
      | /remote.php/dav/spaces |


  Scenario Outline: send PROPFIND request to "/remote.php/dav/(files|spaces)" with depth header
    Given user "Alice" has been created with default attributes and without skeleton files
    When user "Alice" requests "<dav-path>" with "PROPFIND" using basic auth and with headers
      | header | value   |
      | depth  | <depth> |
    Then the HTTP status code should be "<http-code>"
    Examples:
      | dav-path                    | depth    | http-code | 
      | /remote.php/webdav          | 0        | 207       |
      | /remote.php/webdav          | 1        | 207       |
      | /remote.php/dav/files/alice | 0        | 207       |
      | /remote.php/dav/files/alice | 1        | 207       |

    @skipOnRevaMaster
    Examples:
      | dav-path                         | depth    | http-code | 
      | /remote.php/dav/spaces/%spaceid% | 0        | 207       |
      | /remote.php/dav/spaces/%spaceid% | 1        | 207       |
      | /remote.php/dav/spaces/%spaceid% | infinity | 400       |


  Scenario: send PROPFIND request to a public link
    Given user "Alice" has been created with default attributes and without skeleton files
    And user "Alice" has created folder "/PARENT"
    And user "Alice" has created a public link share with settings
      | path        | /PARENT |
      | permissions | read    |
    When the public sends "PROPFIND" request to the last public link share using the new public WebDAV API
    Then the HTTP status code should be "207"
    And the value of the item "//d:href" in the response should match "/%base_path%\/remote.php\/dav\/public-files\/%public_token%\/$/"
    And the value of the item "//oc:public-link-share-owner" in the response should be "Alice"


  Scenario: send PROPFIND request to a public link shared with password
    Given user "Alice" has been created with default attributes and without skeleton files
    And user "Alice" has created folder "/PARENT"
    And user "Alice" has created a public link share with settings
      | path        | /PARENT |
      | permissions | read    |
      | password    | 1111    |
    When the public sends "PROPFIND" request to the last public link share using the new public WebDAV API with password "1111"
    Then the HTTP status code should be "207"
    And the value of the item "//d:href" in the response should match "/%base_path%\/remote.php\/dav\/public-files\/%public_token%\/$/"
    And the value of the item "//oc:public-link-share-owner" in the response should be "Alice"


  Scenario: send PROPFIND request to a public link shared with password (request without password)
    Given user "Alice" has been created with default attributes and without skeleton files
    And user "Alice" has created folder "/PARENT"
    And user "Alice" has created a public link share with settings
      | path        | /PARENT |
      | permissions | read    |
      | password    | 1111    |
    When the public sends "PROPFIND" request to the last public link share using the new public WebDAV API
    Then the HTTP status code should be "401"
    And the value of the item "/d:error/s:exception" in the response should be "Sabre\DAV\Exception\NotAuthenticated"


  Scenario: send PROPFIND request to a public link shared with password (request with incorrect password)
    Given user "Alice" has been created with default attributes and without skeleton files
    And user "Alice" has created folder "/PARENT"
    And user "Alice" has created a public link share with settings
      | path        | /PARENT |
      | permissions | read    |
      | password    | 1111    |
    When the public sends "PROPFIND" request to the last public link share using the new public WebDAV API with password "1234"
    Then the HTTP status code should be "401"
    And the value of the item "/d:error/s:exception" in the response should be "Sabre\DAV\Exception\NotAuthenticated"
