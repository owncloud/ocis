@api
Feature: checksums
  As a user
  I want to upload files with checksum
  So that I can make sure that the files are uploaded with correct checksums

  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |



  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |



  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |

  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |


  Scenario Outline: uploading a file with checksum should work for file with special characters
    Given user "Alice" has been created with default attributes and without skeleton files
    And using <dav_version> DAV path
    When user "Alice" uploads file "filesForUpload/textfile.txt" to <renamed_file> with checksum "MD5:d70b40f177b14b470d1756a3c12b963a" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file <renamed_file> for user "Alice" should be:
      """
      This is a testfile.

      Cheers.
      """
    Examples:
      | dav_version | renamed_file      |
      | old         | " oc?test=ab&cd " |
      | old         | "# %ab ab?=ed"    |
      | new         | " oc?test=ab&cd " |
      | new         | "# %ab ab?=ed"    |

    @skipOnRevaMaster
    Examples:
      | dav_version | renamed_file      |
      | spaces      | " oc?test=ab&cd " |
      | spaces      | "# %ab ab?=ed"    |
