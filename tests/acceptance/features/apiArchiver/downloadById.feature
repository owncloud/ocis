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

  Scenario: download a single file
    Given user "Alice" has uploaded file with content "some data" to "/textfile0.txt"
    When user "Alice" downloads the archive of "/textfile0.txt" using the resource id
    Then the HTTP status code should be "200"
    And the downloaded archive should contain these files:
      | name          | content   |
      | textfile0.txt | some data |

  Scenario: download a single folder
    Given user "Alice" has created folder "my_data"
    And user "Alice" has uploaded file with content "some data" to "/my_data/textfile0.txt"
    And user "Alice" has uploaded file with content "more data" to "/my_data/an_other_file.txt"
    When user "Alice" downloads the archive of "/my_data" using the resource id
    Then the HTTP status code should be "200"
    And the downloaded archive should contain these files:
      | name                      | content   |
      | my_data/textfile0.txt     | some data |
      | my_data/an_other_file.txt | more data |
