@api @skipOnOcV10
Feature: download multiple resources bundled into an archive
  As a user
  I want to be able to download multiple items at once
  So that I don't have to execute repetitive tasks

  As a developer
  I want to be able to use the resource ID to download multiple items at once
  So that I don't have to know the full path of the resource

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files


  Scenario Outline: download a single file
    Given user "Alice" has uploaded file with content "some data" to "/textfile0.txt"
    When user "Alice" downloads the archive of "/textfile0.txt" using the resource id and setting these headers
      | header     | value        |
      | User-Agent | <user-agent> |
    Then the HTTP status code should be "200"
    And the downloaded <archive-type> archive should contain these files:
      | name          | content   |
      | textfile0.txt | some data |
    Examples:
      | user-agent | archive-type |
      | Linux      | tar          |
      | Windows NT | zip          |


  Scenario Outline: download a single folder
    Given user "Alice" has created folder "my_data"
    And user "Alice" has uploaded file with content "some data" to "/my_data/textfile0.txt"
    And user "Alice" has uploaded file with content "more data" to "/my_data/an_other_file.txt"
    When user "Alice" downloads the archive of "/my_data" using the resource id and setting these headers
      | header     | value        |
      | User-Agent | <user-agent> |
    Then the HTTP status code should be "200"
    And the downloaded <archive-type> archive should contain these files:
      | name                      | content   |
      | my_data/textfile0.txt     | some data |
      | my_data/an_other_file.txt | more data |
    Examples:
      | user-agent | archive-type |
      | Linux      | tar          |
      | Windows NT | zip          |


  Scenario: download multiple files and folders
    Given user "Alice" has uploaded file with content "some data" to "/textfile0.txt"
    And user "Alice" has uploaded file with content "other data" to "/textfile1.txt"
    And user "Alice" has created folder "my_data"
    And user "Alice" has uploaded file with content "some data" to "/my_data/textfile2.txt"
    And user "Alice" has created folder "more_data"
    And user "Alice" has uploaded file with content "more data" to "/more_data/an_other_file.txt"
    When user "Alice" downloads the archive of these items using the resource ids
      | textfile0.txt |
      | textfile1.txt |
      | my_data       |
      | more_data     |
    Then the HTTP status code should be "200"
    And the downloaded tar archive should contain these files:
      | name                        | content    |
      | textfile0.txt               | some data  |
      | textfile1.txt               | other data |
      | my_data/textfile2.txt       | some data  |
      | more_data/an_other_file.txt | more data  |


  Scenario: download a single file as different user
    Given user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file with content "some data" to "/textfile0.txt"
    When user "Brian" downloads the archive of "/textfile0.txt" of user "Alice" using the resource id
    Then the HTTP status code should be "400"


  Scenario: download multiple shared items as share receiver
    Given user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file with content "some data" to "/textfile0.txt"
    And user "Alice" has uploaded file with content "other data" to "/textfile1.txt"
    And user "Alice" has created folder "my_data"
    And user "Alice" has uploaded file with content "some data" to "/my_data/textfile2.txt"
    And user "Alice" has created folder "more_data"
    And user "Alice" has uploaded file with content "more data" to "/more_data/an_other_file.txt"
    And user "Alice" has shared file "textfile0.txt" with user "Brian"
    And user "Alice" has shared file "textfile1.txt" with user "Brian"
    And user "Alice" has shared folder "my_data" with user "Brian"
    And user "Alice" has shared folder "more_data" with user "Brian"
    And user "Brian" has accepted share "/textfile0.txt" offered by user "Alice"
    And user "Brian" has accepted share "/textfile1.txt" offered by user "Alice"
    And user "Brian" has accepted share "/my_data" offered by user "Alice"
    And user "Brian" has accepted share "/more_data" offered by user "Alice"
    When user "Brian" downloads the archive of these items using the resource ids
      | /Shares/textfile0.txt |
      | /Shares/textfile1.txt |
      | /Shares/my_data       |
      | /Shares/more_data     |
    Then the HTTP status code should be "200"
    And the downloaded tar archive should contain these files:
      | name                        | content    |
      | textfile0.txt               | some data  |
      | textfile1.txt               | other data |
      | my_data/textfile2.txt       | some data  |
      | more_data/an_other_file.txt | more data  |


  Scenario Outline: download the Shares folder as share receiver
    Given user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file with content "some data" to "/textfile0.txt"
    And user "Alice" has uploaded file with content "other data" to "/textfile1.txt"
    And user "Alice" has created folder "my_data"
    And user "Alice" has uploaded file with content "some data" to "/my_data/textfile2.txt"
    And user "Alice" has created folder "more_data"
    And user "Alice" has uploaded file with content "more data" to "/more_data/an_other_file.txt"
    And user "Alice" has shared file "textfile0.txt" with user "Brian"
    And user "Alice" has shared file "textfile1.txt" with user "Brian"
    And user "Alice" has shared folder "my_data" with user "Brian"
    And user "Alice" has shared folder "more_data" with user "Brian"
    And user "Brian" has accepted share "/textfile0.txt" offered by user "Alice"
    And user "Brian" has accepted share "/textfile1.txt" offered by user "Alice"
    And user "Brian" has accepted share "/my_data" offered by user "Alice"
    And user "Brian" has accepted share "/more_data" offered by user "Alice"
    When user "Brian" downloads the archive of "/Shares" using the resource id and setting these headers
      | header     | value        |
      | User-Agent | <user-agent> |
    Then the HTTP status code should be "200"
    And the downloaded <archive-type> archive should contain these files:
      | name                               | content    |
      | Shares/textfile0.txt               | some data  |
      | Shares/textfile1.txt               | other data |
      | Shares/my_data/textfile2.txt       | some data  |
      | Shares/more_data/an_other_file.txt | more data  |
    Examples:
      | user-agent | archive-type |
      | Linux      | tar          |
      | Windows NT | zip          |
