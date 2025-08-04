@env-config @cli-uploads-sessions
Feature: List upload sessions via CLI command
  As an administrator
  I want to list the upload sessions
  So that I can manage the upload sessions

  Background:
    Given user "Alice" has been created with default attributes


  Scenario: list all upload sessions
    Given user "Alice" has uploaded file with content "uploaded content" to "/file0.txt"
    And the config "POSTPROCESSING_DELAY" has been set to "10s" for "postprocessing" service
    And user "Alice" has uploaded file with content "uploaded content" to "/file1.txt"
    And user "Alice" has uploaded file with content "uploaded content" to "/file2.txt"
    When the administrator lists all the upload sessions
    Then the command should be successful
    And the CLI response should contain these entries:
      | file1.txt |
      | file2.txt |
    And the CLI response should not contain these entries:
      | file0.txt |


  Scenario: list all upload sessions that are currently in postprocessing
    Given the following configs have been set:
      | service        | config                           | value     |
      | postprocessing | POSTPROCESSING_STEPS             | virusscan |
      | antivirus      | ANTIVIRUS_INFECTED_FILE_HANDLING | abort     |
    And user "Alice" has uploaded file "filesForUpload/filesWithVirus/eicar.com" to "/virusFile.txt"
    And the config "POSTPROCESSING_DELAY" has been set to "10s" for "postprocessing" service
    And the administrator has waited for "2" seconds
    And user "Alice" has uploaded file with content "uploaded content" to "/file1.txt"
    And user "Alice" has uploaded file with content "uploaded content" to "/file2.txt"
    When the administrator lists all the upload sessions with flag "processing"
    Then the command should be successful
    And the CLI response should contain these entries:
      | file1.txt |
      | file2.txt |
    And the CLI response should not contain these entries:
      | virusFile.txt |


  Scenario: list all upload sessions that are infected by virus
    Given the following configs have been set:
      | service        | config                           | value     |
      | postprocessing | POSTPROCESSING_STEPS             | virusscan |
      | antivirus      | ANTIVIRUS_INFECTED_FILE_HANDLING | abort     |
    And user "Alice" has uploaded file "filesForUpload/filesWithVirus/eicar.com" to "/virusFile.txt"
    And user "Alice" has uploaded file with content "uploaded content" to "/file1.txt"
    When the administrator lists all the upload sessions with flag "has-virus"
    Then the command should be successful
    And the CLI response should contain these entries:
      | virusFile.txt |
    And the CLI response should not contain these entries:
      | file1.txt |


  Scenario: list and cleanup the expired upload sessions
    Given a file "large.zip" with the size of "2GB" has been created locally
    And the config "STORAGE_USERS_UPLOAD_EXPIRATION" has been set to "1" for "storageuser" service
    And user "Alice" has uploaded a file from "filesForUpload/textfile.txt" to "file.txt" via TUS inside of the space "Personal" using the WebDAV API
    And user "Alice" has tried to upload file "filesForUpload/large.zip" to "large.zip" inside space "Personal" via TUS
    When the administrator lists all the upload sessions with flag "expired"
    Then the command should be successful
    And the CLI response should contain these entries:
      | large.zip |
    And the CLI response should not contain these entries:
      | file.txt |
    When the administrator cleans upload sessions with the following flags:
      | expired |
    Then the command should be successful
    And the CLI response should contain these entries:
      | large.zip |
    And the CLI response should not contain these entries:
      | file.txt |


  Scenario: restart upload sessions that are in postprocessing
    Given the config "POSTPROCESSING_DELAY" has been set to "3s" for "postprocessing" service
    And user "Alice" has uploaded file with content "upload content" to "/file1.txt"
    And user "Alice" has uploaded file with content "upload content" to "/file2.txt"
    And the administrator has waited for "1" seconds
    And the administrator has stopped the server
    And the administrator has started the server
    When the administrator waits for "3" seconds
    Then for user "Alice" file "file1.txt" of space "Personal" should be in postprocessing
    And for user "Alice" file "file2.txt" of space "Personal" should be in postprocessing
    When the administrator restarts the upload sessions that are in postprocessing
    Then the command should be successful
    And the CLI response should contain these entries:
      | file2.txt |
      | file1.txt |
    When the administrator waits for "3" seconds
    Then the content of file "file1.txt" for user "Alice" should be "upload content"
    And the content of file "file2.txt" for user "Alice" should be "upload content"


  Scenario: restart upload sessions of a single file
    Given the config "POSTPROCESSING_DELAY" has been set to "3s" for "postprocessing" service
    And user "Alice" has uploaded file with content "uploaded content" to "file1.txt"
    And user "Alice" has uploaded file with content "uploaded content" to "file2.txt"
    And the administrator has waited for "1" seconds
    And the administrator has stopped the server
    And the administrator has started the server
    When the administrator waits for "3" seconds
    Then for user "Alice" file "file1.txt" of space "Personal" should be in postprocessing
    And for user "Alice" file "file2.txt" of space "Personal" should be in postprocessing
    When the administrator restarts the upload session of file "file1.txt" using the CLI
    Then the command should be successful
    And the CLI response should contain these entries:
      | file1.txt |
    And the CLI response should not contain these entries:
      | file2.txt |
    When the administrator waits for "3" seconds
    Then for user "Alice" file "file2.txt" of space "Personal" should be in postprocessing
    And the content of file "file1.txt" for user "Alice" should be "uploaded content"


  Scenario: clean all upload sessions that are not in post-processing
    Given the following configs have been set:
      | service        | config                           | value     |
      | postprocessing | POSTPROCESSING_STEPS             | virusscan |
      | antivirus      | ANTIVIRUS_INFECTED_FILE_HANDLING | abort     |
    And user "Alice" has uploaded file "filesForUpload/filesWithVirus/eicar.com" to "/virusFile.txt"
    And the config "POSTPROCESSING_DELAY" has been set to "10s" for "postprocessing" service
    And user "Alice" has uploaded file with content "upload content" to "/file1.txt"
    When the administrator cleans upload sessions with the following flags:
      | processing=false |
    Then the command should be successful
    And the CLI response should contain these entries:
      | virusFile.txt |
    And the CLI response should not contain these entries:
      | file1.txt |


  Scenario: clean upload sessions that are not in post-processing and is not virus infected
    Given the following configs have been set:
      | service        | config                           | value     |
      | postprocessing | POSTPROCESSING_STEPS             | virusscan |
      | antivirus      | ANTIVIRUS_INFECTED_FILE_HANDLING | abort     |
      | postprocessing | POSTPROCESSING_DELAY             | 10s       |
    And user "Alice" has uploaded file "filesForUpload/filesWithVirus/eicar.com" to "/virusFile.txt"
    And user "Alice" has uploaded file with content "upload content" to "/file1.txt"
    And user "Alice" has created a new TUS resource in the space "Personal" with the following headers:
      | Upload-Length   | 10                        |
      #    dGV4dEZpbGUudHh0 is the base64 encode of textFile.txt
      | Upload-Metadata | filename dGV4dEZpbGUudHh0 |
      | Tus-Resumable   | 1.0.0                     |
    And user "Alice" has uploaded file with checksum "SHA1 8cb2237d0679ca88db6464eac60da96345513964" to the last created TUS Location with offset "0" and content "12345" via TUS inside of the space "Personal" using the WebDAV API
    When the administrator cleans upload sessions with the following flags:
      | processing=false |
      | has-virus=false  |
    Then the command should be successful
    And the CLI response should contain these entries:
      | textFile.txt |
    And the CLI response should not contain these entries:
      | file1.txt     |
      | virusFile.txt |

  @issue-11290
  Scenario: resume all upload sessions
    Given the following configs have been set:
      | service        | config                           | value           |
      | postprocessing | POSTPROCESSING_STEPS             | virusscan,delay |
      | antivirus      | ANTIVIRUS_INFECTED_FILE_HANDLING | abort           |
      | postprocessing | POSTPROCESSING_DELAY             | 3s              |
    And user "Alice" has uploaded file with content "upload content" to "file.txt"
    And the administrator has waited for "1" seconds
    And the administrator has stopped the server
    And the administrator has started the server
    When the administrator waits for "3" seconds
    Then for user "Alice" file "file.txt" of space "Personal" should be in postprocessing
    When the administrator resumes all the upload sessions using the CLI
    Then the command should be successful
    And the CLI response should contain these entries:
      | file.txt |
    When the administrator waits for "3" seconds
    Then the content of file "file.txt" for user "Alice" should be "upload content"


  Scenario: resume upload session of a single file
    Given the config "POSTPROCESSING_DELAY" has been set to "3s"
    And user "Alice" has uploaded file with content "uploaded content" to "file1.txt"
    And user "Alice" has uploaded file with content "uploaded content" to "file2.txt"
    And the administrator has waited for "1" seconds
    And the administrator has stopped the server
    And the administrator has started the server
    When the administrator waits for "3" seconds
    Then for user "Alice" file "file1.txt" of space "Personal" should be in postprocessing
    And for user "Alice" file "file2.txt" of space "Personal" should be in postprocessing
    When the administrator resumes the upload session of file "file1.txt" using the CLI
    Then the command should be successful
    And the CLI response should contain these entries:
      | file1.txt |
    And the CLI response should not contain these entries:
      | file2.txt |
    When the administrator waits for "3" seconds
    Then for user "Alice" file "file2.txt" of space "Personal" should be in postprocessing
    And the content of file "file1.txt" for user "Alice" should be "uploaded content"


  Scenario: restart expired upload sessions
    Given a file "large.zip" with the size of "2GB" has been created locally
    And the config "STORAGE_USERS_UPLOAD_EXPIRATION" has been set to "1" for "storageuser" service
    And user "Alice" has uploaded a file from "filesForUpload/textfile.txt" to "file.txt" via TUS inside of the space "Personal" using the WebDAV API
    And user "Alice" has tried to upload file "filesForUpload/large.zip" to "large.zip" inside space "Personal" via TUS
    When the administrator restarts the expired upload sessions using the CLI
    Then the command should be successful
    And the CLI response should contain these entries:
      | large.zip |
    And the CLI response should not contain these entries:
      | file.txt |


  Scenario: resume a specific upload session using postprocessing command
    Given the config "POSTPROCESSING_DELAY" has been set to "3s"
    And user "Alice" has uploaded file with content "uploaded content" to "file1.txt"
    And user "Alice" has uploaded file with content "uploaded content" to "file2.txt"
    And the administrator has waited for "1" seconds
    And the administrator has stopped the server
    And the administrator has started the server
    When the administrator waits for "3" seconds
    Then for user "Alice" file "file1.txt" of space "Personal" should be in postprocessing
    And for user "Alice" file "file2.txt" of space "Personal" should be in postprocessing
    When the administrator resumes the upload session of file "file1.txt" using postprocessing command
    Then the command should be successful
    When the administrator waits for "3" seconds
    Then for user "Alice" file "file2.txt" of space "Personal" should be in postprocessing
    And the content of file "file1.txt" for user "Alice" should be "uploaded content"


  Scenario: restart a specific upload session using postprocessing command
    Given the config "POSTPROCESSING_DELAY" has been set to "3s"
    And user "Alice" has uploaded file with content "uploaded content" to "file1.txt"
    And user "Alice" has uploaded file with content "uploaded content" to "file2.txt"
    And the administrator has waited for "1" seconds
    And the administrator has stopped the server
    And the administrator has started the server
    When the administrator waits for "3" seconds
    Then for user "Alice" file "file1.txt" of space "Personal" should be in postprocessing
    And for user "Alice" file "file2.txt" of space "Personal" should be in postprocessing
    When the administrator restarts the upload session of file "file1.txt" using postprocessing command
    Then the command should be successful
    When the administrator waits for "3" seconds
    Then for user "Alice" file "file2.txt" of space "Personal" should be in postprocessing
    And the content of file "file1.txt" for user "Alice" should be "uploaded content"
