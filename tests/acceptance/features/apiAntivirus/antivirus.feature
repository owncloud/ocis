@api @antivirus @skipOnReva
Feature: antivirus
  As a system administrator and user
  I want to protect myself and others from known viruses
  So that I can prevent files with viruses from being uploaded


  Background:
    Given user "Alice" has been created with default attributes and without skeleton files


  Scenario Outline: upload a normal file without virus
    Given using <dav-path-version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to "/normalfile.txt" using the WebDAV API
    Then the HTTP status code should be "201"
    And as "Alice" file "/normalfile.txt" should exist
    And the content of file "/normalfile.txt" for user "Alice" should be:
    """
    This is a testfile.

    Cheers.
    """
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |

  Scenario Outline: upload a file with virus
    Given using <dav-path-version> DAV path
    When user "Alice" uploads file "filesForUpload/filesWithVirus/eicar.com" to "/aFileWithVirus.txt" using the WebDAV API
    # antivirus service can scan files during post-processing. on demand scanning is currently not available
    Then the HTTP status code should be "201"
    And user "Alice" should get a notification with subject "Virus found" and message:
      | message                                                                             |
      | Virus found in aFileWithVirus.txt. Upload not possible. Virus: Win.Test.EICAR_HDB-1 |
    And as "Alice" file "/aFileWithVirus.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: upload a file with virus and a file without virus
    Given using <dav-path-version> DAV path
    When user "Alice" uploads file "filesForUpload/filesWithVirus/eicar.com" to "/aFileWithVirus.txt" using the WebDAV API
    # antivirus service can scan files during post-processing. on demand scanning is currently not available
    Then the HTTP status code should be "201"
    And user "Alice" uploads file "filesForUpload/textfile.txt" to "/normalfile.txt" using the WebDAV API
    And the HTTP status code should be "201"
    And user "Alice" should get a notification with subject "Virus found" and message:
      | message                                                                             |
      | Virus found in aFileWithVirus.txt. Upload not possible. Virus: Win.Test.EICAR_HDB-1 |
    And as "Alice" file "/aFileWithVirus.txt" should not exist
    But as "Alice" file "/normalfile.txt" should exist
    And the content of file "/normalfile.txt" for user "Alice" should be:
    """
    This is a testfile.

    Cheers.
    """
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: upload a file with virus in chunks
    Given using <dav-path-version> DAV path
    When user "Alice" uploads the following chunks to "/myChunkedFile.txt" with old chunking and using the WebDAV API
      | number | content                 |
      | 1      | X5O!P%@AP[4\PZX54(P^)7C |
      | 2      | C)7}$EICAR-STANDARD-ANT |
      | 3      | IVIRUS-TEST-FILE!$H+H*  |
    # antivirus service can scan files during post-processing. on demand scanning is currently not available
    Then the HTTP status code should be "201"
    And user "Alice" should get a notification with subject "Virus found" and message:
      | message                                                                            |
      | Virus found in myChunkedFile.txt. Upload not possible. Virus: Win.Test.EICAR_HDB-1 |
    And as "Alice" file "/myChunkedFile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | spaces           |
