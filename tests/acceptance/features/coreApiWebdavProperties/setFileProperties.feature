Feature: set file properties
  As a user
  I want to be able to set meta-information about files
  So that I can record file meta-information (detailed requirement TBD)

  Background:
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes and without skeleton files

  @smokeTest @issue-1263
  Scenario Outline: setting custom DAV property and reading it
    Given using <dav-path-version> DAV path
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "/testcustomprop.txt"
    And user "Alice" has set property "very-custom-prop" with namespace "x1='http://whatever.org/ns'" of file "/testcustomprop.txt" to "veryCustomPropValue"
    When user "Alice" gets a custom property "very-custom-prop" with namespace "x1='http://whatever.org/ns'" of file "/testcustomprop.txt"
    Then the response should contain a custom "very-custom-prop" property with namespace "x1='http://whatever.org/ns'" and value "veryCustomPropValue"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |

  @issue-1297
  Scenario Outline: setting custom complex DAV property and reading it
    Given using <dav-path-version> DAV path
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "/testcustomprop.txt"
    And user "Alice" has set property "very-custom-prop" with namespace "x1='http://whatever.org/ns'" of file "/testcustomprop.txt" to "<foo xmlns='http://bar'/>"
    When user "Alice" gets a custom property "very-custom-prop" with namespace "x1='http://whatever.org/ns'" of file "/testcustomprop.txt"
    Then the response should contain a custom "very-custom-prop" property with namespace "x1='http://whatever.org/ns'" and complex value "<x2:foo xmlns:x2=\"http://bar\"/>"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |

  @issue-1263
  Scenario Outline: setting custom DAV property and reading it after the file is renamed
    Given using <dav-path-version> DAV path
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "/testcustompropwithmove.txt"
    And user "Alice" has set property "very-custom-prop" with namespace "x1='http://whatever.org/ns'" of file "/testcustompropwithmove.txt" to "valueForMovetest"
    And user "Alice" has moved file "/testcustompropwithmove.txt" to "/catchmeifyoucan.txt"
    When user "Alice" gets a custom property "very-custom-prop" with namespace "x1='http://whatever.org/ns'" of file "/catchmeifyoucan.txt"
    Then the response should contain a custom "very-custom-prop" property with namespace "x1='http://whatever.org/ns'" and value "valueForMovetest"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |

  @issue-1297
  Scenario Outline: setting custom DAV property on a shared file as an owner and reading as a recipient
    Given using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "/testcustompropshared.txt"
    And user "Alice" has sent the following resource share invitation:
      | resource        | testcustompropshared.txt |
      | space           | Personal                 |
      | sharee          | Brian                    |
      | shareType       | user                     |
      | permissionsRole | File Editor              |
    And user "Alice" has set property "very-custom-prop" with namespace "x1='http://whatever.org/ns'" of file "/testcustompropshared.txt" to "valueForSharetest"
    When user "Brian" gets a custom property "very-custom-prop" with namespace "x1='http://whatever.org/ns'" of file "/testcustompropshared.txt"
    Then the response should contain a custom "very-custom-prop" property with namespace "x1='http://whatever.org/ns'" and value "valueForSharetest"
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |

  @issue-1263
  Scenario Outline: setting custom DAV property using one endpoint and reading it with other endpoint
    Given using <dav-path-version> DAV path
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "/testnewold.txt"
    And user "Alice" has set property "very-custom-prop" with namespace "x1='http://whatever.org/ns'" of file "/testnewold.txt" to "lucky"
    And using <dav-path-version-2> DAV path
    When user "Alice" gets a custom property "very-custom-prop" with namespace "x1='http://whatever.org/ns'" of file "/testnewold.txt"
    Then the response should contain a custom "very-custom-prop" property with namespace "x1='http://whatever.org/ns'" and value "lucky"
    Examples:
      | dav-path-version | dav-path-version-2 |
      | old              | new                |
      | new              | old                |
      | spaces           | new                |
      | spaces           | old                |
      | new              | spaces             |
      | old              | spaces             |
