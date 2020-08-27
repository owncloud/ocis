@api @files_trashbin-app-required
Feature: files and folders exist in the trashbin after being deleted
  As a user
  I want deleted files and folders to be available in the trashbin
  So that I can recover data easily

  Background:
    Given user "Alice" has been created with default attributes and skeleton files

  @smokeTest
  @issue-product-178
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario: deleting a file moves it to trashbin
    Given using old DAV path
    When user "Alice" deletes file "/textfile0.txt" using the WebDAV API
    And using new DAV path
    When user "Alice" deletes file "/textfile1.txt" using the WebDAV API
    Then as "Alice" the file with original path "/textfile0.txt" should exist in the trashbin
    And as "Alice" the file with original path "/textfile1.txt" should exist in the trashbin
    And as "Alice" file "/textfile0.txt" should not exist
    And as "Alice" file "/textfile1.txt" should not exist

  @smokeTest
  @issue-product-178
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario: deleting a folder moves it to trashbin
    Given user "Alice" has created folder "/tmp1"
    And user "Alice" has created folder "/tmp2"
    And using old DAV path
    When user "Alice" deletes folder "/tmp1" using the WebDAV API
    And using new DAV path
    When user "Alice" deletes folder "/tmp2" using the WebDAV API
    Then as "Alice" the folder with original path "/tmp1" should exist in the trashbin
    And as "Alice" the folder with original path "/tmp2" should exist in the trashbin
    And as "Alice" folder "/tmp1" should not exist
    And as "Alice" folder "/tmp2" should not exist
