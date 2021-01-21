@api @files_trashbin-app-required
Feature: files and folders can be deleted from the trashbin
  As a user
  I want to delete files and folders from the trashbin
  So that I can control my trashbin space and which files are kept in that space

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files
    And user "Alice" has uploaded file with content "to delete" to "/textfile0.txt"
    And user "Alice" has uploaded file with content "to delete" to "/textfile1.txt"
    And user "Alice" has created folder "PARENT"
    And user "Alice" has created folder "PARENT/CHILD"
    And user "Alice" has uploaded file with content "to delete" to "/PARENT/parent.txt"
    And user "Alice" has uploaded file with content "to delete" to "/PARENT/CHILD/child.txt"

  @smokeTest
  @issue-ocis-reva-118
  @issue-product-179
  @skipOnOcis-OCIS-Storage
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario: delete a single file from the trashbin
    Given user "Alice" has deleted file "/textfile0.txt"
    And user "Alice" has deleted file "/textfile1.txt"
    And user "Alice" has deleted file "/PARENT/parent.txt"
    And user "Alice" has deleted file "/PARENT/CHILD/child.txt"
    When user "Alice" deletes the file with original path "textfile1.txt" from the trashbin using the trashbin API
    Then the HTTP status code should be "405"
    And as "Alice" the file with original path "/textfile1.txt" should exist in the trashbin
    But as "Alice" the file with original path "/textfile0.txt" should exist in the trashbin
    And as "Alice" the file with original path "/PARENT/parent.txt" should exist in the trashbin
    And as "Alice" the file with original path "/PARENT/CHILD/child.txt" should exist in the trashbin
