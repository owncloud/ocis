@env-config
Feature: remove file versions via CLI command


  Scenario: remove all file versions
    Given user "Alice" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file with content "This is version 1" to "textfile.txt"
    And user "Alice" has uploaded file with content "This is version 2" to "textfile.txt"
    And user "Alice" has uploaded file with content "This is version 3" to "textfile.txt"
    When the administrator removes all the file versions using the CLI
    Then the command should be successful
    And the command output should contain "✅ Deleted 2 revisions (6 files / 2 blobs)"
    When user "Alice" gets the number of versions of file "textfile.txt"
    Then the HTTP status code should be "207"
    And the number of versions should be "0"
