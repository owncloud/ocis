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


  Scenario Outline: upload a file with the virus to a public share
    Given using <dav-path-version> DAV path
    And user "Alice" has created folder "/uploadFolder"
    And user "Alice" has created a public link share with settings
      | path        | /uploadFolder            |
      | name        | sharedlink               |
      | permissions | change                   |
      | expireDate  | 2040-01-01T23:59:59+0100 |
    When user "Alice" uploads file "filesForUpload/filesWithVirus/eicar.com" to "/virusFile.txt" using the WebDAV API
    Then the HTTP status code should be "201"
    And user "Alice" should get a notification with subject "Virus found" and message:
      | message                                                                        |
      | Virus found in virusFile.txt. Upload not possible. Virus: Win.Test.EICAR_HDB-1 |
    And as "Alice" file "/uploadFolder/virusFile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: upload a file with the virus to a password-protected public share
    Given using <dav-path-version> DAV path
    And user "Alice" has created folder "/uploadFolder"
    And user "Alice" has created a public link share with settings
      | path        | /uploadFolder            |
      | name        | sharedlink               |
      | permissions | change                   |
      | password    | newpasswd                |
      | expireDate  | 2040-01-01T23:59:59+0100 |
    When user "Alice" uploads file "filesForUpload/filesWithVirus/eicar.com" to "/virusFile.txt" using the WebDAV API
    Then the HTTP status code should be "201"
    And user "Alice" should get a notification with subject "Virus found" and message:
      | message                                                                        |
      | Virus found in virusFile.txt. Upload not possible. Virus: Win.Test.EICAR_HDB-1 |
    And as "Alice" file "/uploadFolder/virusFile.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |
      | spaces           |


  Scenario Outline: upload a file with virus to a user share
    Given using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has created folder "uploadFolder"
    And user "Alice" has shared folder "uploadFolder" with user "Brian" with permissions "all"
    And user "Brian" has accepted share "/uploadFolder" offered by user "Alice"
    When user "Brian" uploads file "filesForUpload/filesWithVirus/<filename>" to "/Shares/uploadFolder/<newfilename>" using the WebDAV API
    Then the HTTP status code should be "201"
    And user "Brian" should get a notification with subject "Virus found" and message:
      | message                                                                        |
      | Virus found in <newfilename>. Upload not possible. Virus: Win.Test.EICAR_HDB-1 |
    And as "Brian" file "/Shares/uploadFolder/<newfilename>" should not exist
    And as "Alice" file "/uploadFolder/<newfilename>" should not exist
    Examples:
      | dav-path-version | filename      | newfilename    |
      | old              | eicar.com     | virusFile1.txt |
      | old              | eicar_com.zip | virusFile2.zip |
      | new              | eicar.com     | virusFile1.txt |
      | new              | eicar_com.zip | virusFile2.zip |


  Scenario Outline: upload a file with virus to a user share using spaces dav endpoint
    Given using spaces DAV path
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has created folder "uploadFolder"
    And user "Alice" has shared folder "uploadFolder" with user "Brian" with permissions "all"
    And user "Brian" has accepted share "/uploadFolder" offered by user "Alice"
    When user "Brian" uploads a file "filesForUpload/filesWithVirus/<filename>" to "/uploadFolder/<newfilename>" in space "Shares" using the WebDAV API
    Then the HTTP status code should be "201"
    And user "Brian" should get a notification with subject "Virus found" and message:
      | message                                                                        |
      | Virus found in <newfilename>. Upload not possible. Virus: Win.Test.EICAR_HDB-1 |
    And as "Brian" file "/Shares/uploadFolder/<newfilename>" should not exist
    And as "Alice" file "/uploadFolder/<newfilename>" should not exist
    Examples:
      | filename      | newfilename    |
      | eicar.com     | virusFile1.txt |
      | eicar_com.zip | virusFile2.zip |


  Scenario Outline: upload a file with virus to a group share
    Given using <dav-path-version> DAV path
    And user "Brian" has been created with default attributes and without skeleton files
    And group "group1" has been created
    And user "Brian" has been added to group "group1"
    And user "Alice" has created folder "uploadFolder"
    And user "Alice" has shared folder "uploadFolder" with group "group1"
    And user "Brian" has accepted share "/uploadFolder" offered by user "Alice"
    When user "Brian" uploads file "filesForUpload/filesWithVirus/<filename>" to "/Shares/uploadFolder/<newfilename>" using the WebDAV API
    Then the HTTP status code should be "201"
    And user "Brian" should get a notification with subject "Virus found" and message:
      | message                                                                        |
      | Virus found in <newfilename>. Upload not possible. Virus: Win.Test.EICAR_HDB-1 |
    And as "Brian" file "/Shares/uploadFolder/<newfilename>" should not exist
    And as "Alice" file "/uploadFolder/<newfilename>" should not exist
    Examples:
      | dav-path-version | filename      | newfilename    |
      | old              | eicar.com     | virusFile1.txt |
      | old              | eicar_com.zip | virusFile2.zip |
      | new              | eicar.com     | virusFile1.txt |
      | new              | eicar_com.zip | virusFile2.zip |


  Scenario Outline: upload a file with virus to a group share using spaces dav endpoint
    Given using spaces DAV path
    And user "Brian" has been created with default attributes and without skeleton files
    And group "group1" has been created
    And user "Brian" has been added to group "group1"
    And user "Alice" has created folder "uploadFolder"
    And user "Alice" has shared folder "uploadFolder" with group "group1"
    And user "Brian" has accepted share "/uploadFolder" offered by user "Alice"
    When user "Brian" uploads a file "filesForUpload/filesWithVirus/<filename>" to "/uploadFolder/<newfilename>" in space "Shares" using the WebDAV API
    Then the HTTP status code should be "201"
    And user "Brian" should get a notification with subject "Virus found" and message:
      | message                                                                        |
      | Virus found in <newfilename>. Upload not possible. Virus: Win.Test.EICAR_HDB-1 |
    And as "Brian" file "/Shares/uploadFolder/<newfilename>" should not exist
    And as "Alice" file "/uploadFolder/<newfilename>" should not exist
    Examples:
      | filename      | newfilename    |
      | eicar.com     | virusFile1.txt |
      | eicar_com.zip | virusFile2.zip |

