Feature: move (rename) file
  As a user
  I want to be able to move and rename files
  So that I can manage my file system

  Background:
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes and without skeleton files

  @smokeTest
  Scenario Outline: moving a file
    Given using <dav-path-version> DAV path
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has uploaded file with content "ownCloud test text file 0" to "textfile0.txt"
    When user "Alice" moves file "/textfile0.txt" to "/FOLDER/textfile0.txt" using the WebDAV API
    Then the HTTP status code should be "201"
    And the following headers should match these regular expressions for user "Alice"
      | ETag | /^"[a-f0-9:\.]{1,32}"$/ |
    And the content of file "/FOLDER/textfile0.txt" for user "Alice" should be "ownCloud test text file 0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |

    @skipOnRevaMaster
    Examples:
      | dav-path-version |
      | spaces           |

  @smokeTest
  Scenario Outline: moving and overwriting a file
    Given using <dav-path-version> DAV path
    And user "Alice" has uploaded file with content "ownCloud test text file 0" to "textfile0.txt"
    And user "Alice" has uploaded file with content "ownCloud test text file 1" to "textfile1.txt"
    When user "Alice" moves file "/textfile0.txt" to "/textfile1.txt" using the WebDAV API
    Then the HTTP status code should be "204"
    And the following headers should match these regular expressions for user "Alice"
      | ETag | /^"[a-f0-9:\.]{1,32}"$/ |
    And the content of file "/textfile1.txt" for user "Alice" should be "ownCloud test text file 0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |

    @skipOnRevaMaster
    Examples:
      | dav-path-version |
      | spaces           |


  Scenario Outline: moving (renaming) a file to be only different case
    Given using <dav-path-version> DAV path
    And user "Alice" has uploaded file with content "ownCloud test text file 0" to "textfile0.txt"
    When user "Alice" moves file "/textfile0.txt" to "/TextFile0.txt" using the WebDAV API
    Then the HTTP status code should be "201"
    And as "Alice" file "/textfile0.txt" should not exist
    And the content of file "/TextFile0.txt" for user "Alice" should be "ownCloud test text file 0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |

    @skipOnRevaMaster
    Examples:
      | dav-path-version |
      | spaces           |

  @smokeTest
  Scenario Outline: moving (renaming) a file to a file with only different case to an existing file
    Given using <dav-path-version> DAV path
    And user "Alice" has uploaded file with content "ownCloud test text file 0" to "textfile0.txt"
    And user "Alice" has uploaded file with content "ownCloud test text file 1" to "textfile1.txt"
    When user "Alice" moves file "/textfile1.txt" to "/TextFile0.txt" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file "/textfile0.txt" for user "Alice" should be "ownCloud test text file 0"
    And the content of file "/TextFile0.txt" for user "Alice" should be "ownCloud test text file 1"
    Examples:
      | dav-path-version |
      | old              |
      | new              |

    @skipOnRevaMaster
    Examples:
      | dav-path-version |
      | spaces           |


  Scenario Outline: moving (renaming) a file to a file in a folder with only different case to an existing file
    Given using <dav-path-version> DAV path
    And user "Alice" has created folder "PARENT"
    And user "Alice" has uploaded file with content "ownCloud test text file parent" to "PARENT/parent.txt"
    And user "Alice" has uploaded file with content "ownCloud test text file 1" to "textfile1.txt"
    When user "Alice" moves file "/textfile1.txt" to "/PARENT/Parent.txt" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file "/PARENT/parent.txt" for user "Alice" should be "ownCloud test text file parent"
    And the content of file "/PARENT/Parent.txt" for user "Alice" should be "ownCloud test text file 1"
    Examples:
      | dav-path-version |
      | old              |
      | new              |

    @skipOnRevaMaster
    Examples:
      | dav-path-version |
      | spaces           |


  Scenario Outline: move file into a not-existing folder
    Given using <dav-path-version> DAV path
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "fileToMove.txt"
    When user "Alice" moves file "/fileToMove.txt" to "/not-existing/fileToMove.txt" using the WebDAV API
    Then the HTTP status code should be "409"
    Examples:
      | dav-path-version |
      | old              |
      | new              |

    @skipOnRevaMaster
    Examples:
      | dav-path-version |
      | spaces           |

  @issue-1295
  Scenario Outline: rename a file into an invalid filename
    Given using <dav-path-version> DAV path
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "fileToRename.txt"
    When user "Alice" moves file "/fileToRename.txt" to "/a\\a" using the WebDAV API
    Then the HTTP status code should be "400"
    Examples:
      | dav-path-version |
      | old              |
      | new              |

    @skipOnRevaMaster
    Examples:
      | dav-path-version |
      | spaces           |


  Scenario Outline: checking file id after a move
    Given using <dav-path-version> DAV path
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "textfile0.txt"
    And user "Alice" has stored id of file "/textfile0.txt"
    When user "Alice" moves file "/textfile0.txt" to "/FOLDER/textfile0.txt" using the WebDAV API
    Then the HTTP status code should be "201"
    And user "Alice" file "/FOLDER/textfile0.txt" should have the previously stored id
    And user "Alice" should not see the following elements
      | /textfile0.txt |
    Examples:
      | dav-path-version |
      | old              |
      | new              |

    @skipOnRevaMaster
    Examples:
      | dav-path-version |
      | spaces           |

  @issue-1295
  Scenario Outline: renaming a file to a path with extension .part should not be possible
    Given using <dav-path-version> DAV path
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "fileToRename.txt"
    When user "Alice" moves file "/fileToRename.txt" to "/welcome.part" using the WebDAV API
    Then the HTTP status code should be "400"
    And the DAV exception should be "OCA\DAV\Connector\Sabre\Exception\InvalidPath"
    And the DAV message should be "Can`t upload files with extension .part because these extensions are reserved for internal use."
    And the DAV reason should be "Can`t upload files with extension .part because these extensions are reserved for internal use."
    And user "Alice" should see the following elements
      | /fileToRename.txt |
    But user "Alice" should not see the following elements
      | /fileToRename.part |
    Examples:
      | dav-path-version |
      | old              |
      | new              |

    @skipOnRevaMaster
    Examples:
      | dav-path-version |
      | spaces           |

  @sqliteDB
  Scenario Outline: renaming to a file with special characters
    Given using <dav-path-version> DAV path
    And user "Alice" has uploaded file with content "ownCloud test text file 0" to "textfile0.txt"
    And user "Alice" has uploaded file with content "ownCloud test text file 1" to "textfile1.txt"
    And user "Alice" has uploaded file with content "ownCloud test text file 2" to "textfile2.txt"
    And user "Alice" has uploaded file with content "ownCloud test text file 3" to "textfile3.txt"
    When user "Alice" moves the following file using the WebDAV API
      | source         | destination   |
      | /textfile0.txt | *a@b#c$e%f&g* |
      | /textfile1.txt | 1 2 3##.##    |
      | /textfile2.txt | file[2]       |
      | /textfile3.txt | file [ 3 ]    |
    Then the HTTP status code of responses on all endpoints should be "201"
    And the content of file "*a@b#c$e%f&g*" for user "Alice" should be "ownCloud test text file 0"
    And the content of file "1 2 3##.##" for user "Alice" should be "ownCloud test text file 1"
    And the content of file "file[2]" for user "Alice" should be "ownCloud test text file 2"
    And the content of file "file [ 3 ]" for user "Alice" should be "ownCloud test text file 3"
    Examples:
      | dav-path-version |
      | old              |
      | new              |

    @skipOnRevaMaster
    Examples:
      | dav-path-version |
      | spaces           |

  @issue-1259
  #after fixing the issues merge this Scenario into the one above
  Scenario Outline: renaming to a file with question mark in its name
    Given using <dav-path-version> DAV path
    And user "Alice" has uploaded file with content "ownCloud test text file 0" to "textfile0.txt"
    When user "Alice" moves file "/textfile0.txt" to "/#oc ab?cd=ef#" using the WebDAV API
    Then the HTTP status code should be "201"
    And the content of file "/#oc ab?cd=ef#" for user "Alice" should be "ownCloud test text file 0"
    Examples:
      | dav-path-version |
      | old              |
      | new              |

    @skipOnRevaMaster
    Examples:
      | dav-path-version |
      | spaces           |


  Scenario Outline: renaming file with dots in the path
    Given using <dav-path-version> DAV path
    And user "Alice" has created folder "<folder_name>"
    And user "Alice" has uploaded file with content "uploaded content for file name ending with a dot" to "<folder_name>/<file_name>"
    When user "Alice" moves file "<folder_name>/<file_name>" to "<folder_name>/abc.txt" using the WebDAV API
    Then the HTTP status code should be "201"
    And as "Alice" file "<folder_name>/abc.txt" should exist
    Examples:
      | dav-path-version | folder_name   | file_name   |
      | old              | /upload.      | abc.        |
      | old              | /upload.      | abc .       |
      | old              | /upload.1     | abc         |
      | old              | /upload...1.. | abc...txt.. |
      | old              | /...          | abcd.txt    |
      | old              | /..upload     | ..abc       |
      | new              | /upload.      | abc.        |
      | new              | /upload.      | abc .       |
      | new              | /upload.1     | ..abc.txt   |
      | new              | /upload...1.. | abc...txt.. |
      | new              | /...          | ...         |
      | new              | /..upload     | ..abc       |

    @skipOnRevaMaster
    Examples:
      | dav-path-version | folder_name   | file_name   |
      | spaces           | /upload.      | abc.        |
      | spaces           | /upload.      | abc .       |
      | spaces           | /upload.1     | abc         |
      | spaces           | /upload...1.. | abc...txt.. |
      | spaces           | /...          | abcd.txt    |
      | spaces           | /...          | ...         |

  @smokeTest
  Scenario Outline: user tries to move a file that doesnt exist into a folder
    Given using <dav-path-version> DAV path
    And user "Alice" has created folder "FOLDER"
    When user "Alice" moves file "/doesNotExist.txt" to "/FOLDER/textfile0.txt" using the WebDAV API
    Then the HTTP status code should be "404"
    And as "Alice" file "/FOLDER/textfile0.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |

    @skipOnRevaMaster
    Examples:
      | dav-path-version |
      | spaces           |

  @smokeTest
  Scenario Outline: user tries to rename a file that doesn't exist
    Given using <dav-path-version> DAV path
    When user "Alice" moves file "/doesNotExist.txt" to "/exist.txt" using the WebDAV API
    Then the HTTP status code should be "404"
    And as "Alice" file "/exist.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |

    @skipOnRevaMaster
    Examples:
      | dav-path-version |
      | spaces           |


  Scenario Outline: moving a hidden file
    Given using <dav-path-version> DAV path
    And user "Alice" has created folder "FOLDER"
    And user "Alice" has uploaded the following files with content "hidden file"
      | path                    |
      | .hidden_file101         |
      | /FOLDER/.hidden_file102 |
    When user "Alice" moves the following files using the WebDAV API
      | from                    | to                      |
      | .hidden_file101         | /FOLDER/.hidden_file101 |
      | /FOLDER/.hidden_file102 | .hidden_file102         |
    Then the HTTP status code of responses on all endpoints should be "201"
    And as "Alice" the following files should exist
      | path                    |
      | .hidden_file102         |
      | /FOLDER/.hidden_file101 |
    And the content of the following files for user "Alice" should be "hidden file"
      | path                    |
      | .hidden_file102         |
      | /FOLDER/.hidden_file101 |
    Examples:
      | dav-path-version |
      | old              |
      | new              |

    @skipOnRevaMaster
    Examples:
      | dav-path-version |
      | spaces           |


  Scenario Outline: renaming to/from a hidden file
    Given using <dav-path-version> DAV path
    And user "Alice" has uploaded the following files with content "hidden file"
      | path               |
      | .hidden_file101    |
      | hidden_file101.txt |
    When user "Alice" moves the following files using the WebDAV API
      | from               | to                 |
      | .hidden_file101    | hidden_file102.txt |
      | hidden_file101.txt | .hidden_file102    |
    Then the HTTP status code of responses on all endpoints should be "201"
    And as "Alice" the following files should exist
      | path               |
      | .hidden_file102    |
      | hidden_file102.txt |
    And the content of the following files for user "Alice" should be "hidden file"
      | path               |
      | .hidden_file102    |
      | hidden_file102.txt |
    Examples:
      | dav-path-version |
      | old              |
      | new              |

    @skipOnRevaMaster
    Examples:
      | dav-path-version |
      | spaces           |


  Scenario Outline: moving a file (deep moves with various folder and file names)
    Given using <dav-path-version> DAV path
    And user "Alice" has created folder "<source_folder>"
    And user "Alice" has created folder "<target_folder>"
    And user "Alice" has uploaded file with content "ownCloud test text file 0" to "/<source_folder>/<source_file>"
    When user "Alice" moves file "/<source_folder>/<source_file>" to "/<target_folder>/<target_file>" using the WebDAV API
    Then the HTTP status code should be "201"
    And the following headers should match these regular expressions for user "Alice"
      | ETag | /^"[a-f0-9:\.]{1,32}"$/ |
    And the content of file "/<target_folder>/<target_file>" for user "Alice" should be "ownCloud test text file 0"
    Examples:
      | dav-path-version | source_folder | source_file | target_folder | target_file |
      | old              | text          | file.txt    | 0             | file.txt    |
      | old              | text          | file.txt    | 1             | file.txt    |
      | old              | 0             | file.txt    | text          | file.txt    |
      | old              | 1             | file.txt    | text          | file.txt    |
      | old              | texta         | 0           | textb         | file.txt    |
      | old              | texta         | 1           | textb         | file.txt    |
      | old              | texta         | file.txt    | textb         | 0           |
      | old              | texta         | file.txt    | textb         | 1           |
      | new              | text          | file.txt    | 0             | file.txt    |
      | new              | text          | file.txt    | 1             | file.txt    |
      | new              | 0             | file.txt    | text          | file.txt    |
      | new              | 1             | file.txt    | text          | file.txt    |
      | new              | texta         | 0           | textb         | file.txt    |
      | new              | texta         | 1           | textb         | file.txt    |
      | new              | texta         | file.txt    | textb         | 0           |
      | new              | texta         | file.txt    | textb         | 1           |

    @skipOnRevaMaster
    Examples:
      | dav-path-version | source_folder | source_file | target_folder | target_file |
      | spaces           | text          | file.txt    | 0             | file.txt    |
      | spaces           | text          | file.txt    | 1             | file.txt    |
      | spaces           | 0             | file.txt    | text          | file.txt    |
      | spaces           | 1             | file.txt    | text          | file.txt    |
      | spaces           | texta         | 0           | textb         | file.txt    |
      | spaces           | texta         | 1           | textb         | file.txt    |
      | spaces           | texta         | file.txt    | textb         | 0           |
      | spaces           | texta         | file.txt    | textb         | 1           |


  Scenario Outline: moving a file from a folder to the root
    Given using <dav-path-version> DAV path
    And user "Alice" has created folder "<source_folder>"
    And user "Alice" has uploaded file with content "ownCloud test text file 0" to "/<source_folder>/<source_file>"
    When user "Alice" moves file "/<source_folder>/<source_file>" to "/<target_file>" using the WebDAV API
    Then the HTTP status code should be "201"
    And the following headers should match these regular expressions for user "Alice"
      | ETag | /^"[a-f0-9:\.]{1,32}"$/ |
    And the content of file "/<target_file>" for user "Alice" should be "ownCloud test text file 0"
    Examples:
      | dav-path-version | source_folder | source_file | target_file |
      | old              | 0             | file.txt    | file.txt    |
      | old              | 1             | file.txt    | file.txt    |
      | old              | texta         | 0           | file.txt    |
      | old              | texta         | 1           | file.txt    |
      | old              | texta         | file.txt    | 0           |
      | old              | texta         | file.txt    | 1           |
      | new              | 0             | file.txt    | file.txt    |
      | new              | 1             | file.txt    | file.txt    |
      | new              | texta         | 0           | file.txt    |
      | new              | texta         | 1           | file.txt    |
      | new              | texta         | file.txt    | 0           |
      | new              | texta         | file.txt    | 1           |

    @skipOnRevaMaster
    Examples:
      | dav-path-version | source_folder | source_file | target_file |
      | spaces           | 0             | file.txt    | file.txt    |
      | spaces           | 1             | file.txt    | file.txt    |
      | spaces           | texta         | 0           | file.txt    |
      | spaces           | texta         | 1           | file.txt    |
      | spaces           | texta         | file.txt    | 0           |
      | spaces           | texta         | file.txt    | 1           |


  Scenario Outline: move a file of size zero byte
    Given using <dav-path-version> DAV path
    And user "Alice" has uploaded file "filesForUpload/zerobyte.txt" to "/zerobyte.txt"
    And user "Alice" has created folder "/testZeroByte"
    When user "Alice" moves file "/zerobyte.txt" to "/testZeroByte/zerobyte.txt" using the WebDAV API
    Then the HTTP status code should be "201"
    And as "Alice" file "/testZeroByte/zerobyte.txt" should exist
    And as "Alice" file "/zerobyte.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |

    @skipOnRevaMaster
    Examples:
      | dav-path-version |
      | spaces           |


  Scenario Outline: rename a file of size zero byte
    Given using <dav-path-version> DAV path
    And user "Alice" has uploaded file "filesForUpload/zerobyte.txt" to "/zerobyte.txt"
    When user "Alice" moves file "/zerobyte.txt" to "/rename_zerobyte.txt" using the WebDAV API
    Then the HTTP status code should be "201"
    And as "Alice" file "/rename_zerobyte.txt" should exist
    And as "Alice" file "/zerobyte.txt" should not exist
    Examples:
      | dav-path-version |
      | old              |
      | new              |

    @skipOnRevaMaster
    Examples:
      | dav-path-version |
      | spaces           |
