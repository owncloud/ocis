@env-config @cli-uploads-sessions
Feature: List upload sessions via CLI command
  As a user
  I want to list the upload sessions
  So that I can manage the upload sessions

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files


  Scenario: list all upload sessions
    Given user "Alice" has uploaded file with content "uploaded content" to "/file0.txt"
    And the config "POSTPROCESSING_DELAY" has been set to "10s"
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
      | config                           | value     |
      | POSTPROCESSING_STEPS             | virusscan |
      | ANTIVIRUS_INFECTED_FILE_HANDLING | abort     |
    And user "Alice" has uploaded file "filesForUpload/filesWithVirus/eicar.com" to "/virusFile.txt"
    And the config "POSTPROCESSING_DELAY" has been set to "10s"
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
      | config                           | value     |
      | POSTPROCESSING_STEPS             | virusscan |
      | ANTIVIRUS_INFECTED_FILE_HANDLING | abort     |
    And user "Alice" has uploaded file "filesForUpload/filesWithVirus/eicar.com" to "/virusFile.txt"
    And user "Alice" has uploaded file with content "uploaded content" to "/file1.txt"
    When the administrator lists all the upload sessions with flag "has-virus"
    Then the command should be successful
    And the CLI response should contain these entries:
      | virusFile.txt |
    And the CLI response should not contain these entries:
      | file1.txt |


  Scenario: list all expired upload sessions
    Given the config "POSTPROCESSING_DELAY" has been set to "10s"
    And user "Alice" has uploaded file with content "uploaded content" to "/file1.txt"
    And the config "STORAGE_USERS_UPLOAD_EXPIRATION" has been set to "0"
    And user "Alice" has uploaded file with content "uploaded content" to "/file2.txt"
    And user "Alice" has uploaded file with content "uploaded content" to "/file3.txt"
    When the administrator lists all the upload sessions with flag "expired"
    Then the command should be successful
    And the CLI response should contain these entries:
      | file2.txt |
      | file3.txt |
    And the CLI response should not contain these entries:
      | file1.txt |


  Scenario: clean all expired upload sessions
    Given the config "POSTPROCESSING_DELAY" has been set to "10s"
    And user "Alice" has uploaded file with content "upload content" to "/file1.txt"
    And the config "STORAGE_USERS_UPLOAD_EXPIRATION" has been set to "0"
    And user "Alice" has uploaded file with content "upload content" to "/file2.txt"
    And user "Alice" has uploaded file with content "upload content" to "/file3.txt"
    When the administrator cleans all expired upload sessions
    Then the command should be successful
    And the CLI response should contain these entries:
      | file2.txt |
      | file3.txt |
    And the CLI response should not contain these entries:
      | file1.txt |


  Scenario: restart upload sessions that are in postprocessing
    Given user "Alice" has uploaded file with content "upload content" to "/file1.txt"
    And the config "POSTPROCESSING_DELAY" has been set to "10s"
    And user "Alice" has uploaded file with content "upload content" to "/file2.txt"
    When the administrator restarts the upload sessions that are in postprocessing
    Then the command should be successful
    And the CLI response should contain these entries:
      | file2.txt |
    And the CLI response should not contain these entries:
      | file1.txt |


  Scenario: restart upload sessions of a single file
    Given the config "POSTPROCESSING_DELAY" has been set to "10s"
    And user "Alice" has uploaded file with content "upload content" to "/file1.txt"
    And user "Alice" has uploaded file with content "upload content" to "/file2.txt"
    When the administrator restarts the upload sessions of file "file1.txt"
    Then the command should be successful
    And the CLI response should contain these entries:
      | file1.txt |
    And the CLI response should not contain these entries:
      | file2.txt |
