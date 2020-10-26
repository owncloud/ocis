@api @issue-ocis-reva-57
Feature: set file properties
  As a user
  I want to be able to set meta-information about files
  So that I can reccord file meta-information (detailed requirement TBD)

  Background:
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes and without skeleton files

  @skipOnOcis-OC-Storage @issue-ocis-reva-276 @skipOnOcis-OCIS-Storage
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario Outline: Setting custom DAV property
    Given using <dav_version> DAV path
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "/testcustomprop.txt"
    When user "Alice" sets property "very-custom-prop"  with namespace "x1='http://whatever.org/ns'" of file "/testcustomprop.txt" to "veryCustomPropValue" using the WebDAV API
    Then the HTTP status code should be "500"
    Examples:
      | dav_version |
      | old         |
      | new         |
