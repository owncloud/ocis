@api
Feature: get file properties
  As a user
  I want to be able to get meta-information about files
  So that I can know file meta-information (detailed requirement TBD)

  Background:
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes and without skeleton files

  @issue-ocis-reva-214 @skipOnOcis-OCIS-Storage
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario Outline: Do a PROPFIND of various file names
    Given using <dav_version> DAV path
    And user "Alice" has uploaded file with content "uploaded content" to "<file_name>"
    When user "Alice" gets the properties of file "<file_name>" using the WebDAV API
    Then the properties response should contain an etag
    And the value of the item "//d:response/d:href" in the response to user "Alice" should match "/remote\.php\/<expected_href>/"
    Examples:
      | dav_version | file_name     | expected_href                             |
      | old         | /C++ file.cpp | webdav\/C\+\+%20file\.cpp                 |
      | old         | /file #2.txt  | webdav\/file%20%232\.txt                  |
      | old         | /file &2.txt  | webdav\/file%20&2\.txt                    |
      | new         | /C++ file.cpp | dav\/files\/%username%\/C\+\+%20file\.cpp |
      | new         | /file #2.txt  | dav\/files\/%username%\/file%20%232\.txt  |
      | new         | /file &2.txt  | dav\/files\/%username%\/file%20&2\.txt    |

  @issue-ocis-reva-214 @issue-ocis-reva-265 @skipOnOcis-EOS-Storage @skipOnOcis-OCIS-Storage
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario Outline: Do a PROPFIND of various file names
    Given using <dav_version> DAV path
    And user "Alice" has uploaded file with content "uploaded content" to "<file_name>"
    When user "Alice" gets the properties of file "<file_name>" using the WebDAV API
    Then the properties response should contain an etag
    And the value of the item "//d:response/d:href" in the response to user "Alice" should match "/remote\.php\/<expected_href>/"
    Examples:
      | dav_version | file_name    | expected_href                            |
      | old         | /file ?2.txt | webdav\/file%20%3F2\.txt                 |
      | new         | /file ?2.txt | dav\/files\/%username%\/file%20%3F2\.txt |

  @issue-ocis-reva-214 @skipOnOcis-OCIS-Storage
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario Outline: Do a PROPFIND of various folder names
    Given using <dav_version> DAV path
    And user "Alice" has created folder "<folder_name>"
    And user "Alice" has uploaded file with content "uploaded content" to "<folder_name>/file1.txt"
    And user "Alice" has uploaded file with content "uploaded content" to "<folder_name>/file2.txt"
    When user "Alice" gets the properties of folder "<folder_name>" with depth 1 using the WebDAV API
    Then the value of the item "//d:response[1]/d:href" in the response to user "Alice" should match "/remote\.php\/<expected_href>\//"
    And the value of the item "//d:response[2]/d:href" in the response to user "Alice" should match "/remote\.php\/<expected_href>\/file1.txt/"
    And the value of the item "//d:response[3]/d:href" in the response to user "Alice" should match "/remote\.php\/<expected_href>\/file2.txt/"
    Examples:
      | dav_version | folder_name     | expected_href                                                                  |
      | old         | /upload         | webdav\/upload                                                                 |
      | old         | /strängé folder | webdav\/str%C3%A4ng%C3%A9%20folder                                             |
      | old         | /C++ folder     | webdav\/C\+\+%20folder                                                         |
      | old         | /नेपाली         | webdav\/%E0%A4%A8%E0%A5%87%E0%A4%AA%E0%A4%BE%E0%A4%B2%E0%A5%80                 |
      | old         | /folder #2.txt  | webdav\/folder%20%232\.txt                                                     |
      | old         | /folder &2.txt  | webdav\/folder%20&2\.txt                                                       |
      | new         | /upload         | dav\/files\/%username%\/upload                                                 |
      | new         | /strängé folder | dav\/files\/%username%\/str%C3%A4ng%C3%A9%20folder                             |
      | new         | /C++ folder     | dav\/files\/%username%\/C\+\+%20folder                                         |
      | new         | /नेपाली         | dav\/files\/%username%\/%E0%A4%A8%E0%A5%87%E0%A4%AA%E0%A4%BE%E0%A4%B2%E0%A5%80 |
      | new         | /folder #2.txt  | dav\/files\/%username%\/folder%20%232\.txt                                     |
      | new         | /folder &2.txt  | dav\/files\/%username%\/folder%20&2\.txt                                       |

  @issue-ocis-reva-214 @skipOnOcis-EOS-Storage @issue-ocis-reva-265 @skipOnOcis-OCIS-Storage
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario Outline: Do a PROPFIND of various folder names
    Given using <dav_version> DAV path
    And user "Alice" has created folder "<folder_name>"
    And user "Alice" has uploaded file with content "uploaded content" to "<folder_name>/file1.txt"
    And user "Alice" has uploaded file with content "uploaded content" to "<folder_name>/file2.txt"
    When user "Alice" gets the properties of folder "<folder_name>" with depth 1 using the WebDAV API
    Then the value of the item "//d:response[1]/d:href" in the response to user "Alice" should match "/remote\.php\/<expected_href>\//"
    And the value of the item "//d:response[2]/d:href" in the response to user "Alice" should match "/remote\.php\/<expected_href>\/file1.txt/"
    And the value of the item "//d:response[3]/d:href" in the response to user "Alice" should match "/remote\.php\/<expected_href>\/file2.txt/"
    Examples:
      | dav_version | folder_name    | expected_href                              |
      | old         | /folder ?2.txt | webdav\/folder%20%3F2\.txt                 |
      | new         | /folder ?2.txt | dav\/files\/%username%\/folder%20%3F2\.txt |

  @skipOnOcis-OC-Storage @issue-ocis-reva-265 @skipOnOcis-OCIS-Storage
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario Outline: Do a PROPFIND of various folder names
    Given using <dav_version> DAV path
    And user "Alice" has created folder "/folder ?2.txt"
    When user "Alice" uploads to these filenames with content "uploaded content" using the webDAV API then the results should be as listed
      | filename                 | http-code | exists |
      | /folder ?2.txt/file1.txt | 500       | no     |
    Examples:
      | dav_version |
      | old         |
      | new         |

  @issue-ocis-reva-163
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario Outline: Do a PROPFIND to a non-existing URL
    And user "Alice" requests "<url>" with "PROPFIND" using basic auth
    Then the body of the response should be empty
    Examples:
      | url                                  |
      | /remote.php/dav/files/does-not-exist |
      | /remote.php/dav/does-not-exist       |
