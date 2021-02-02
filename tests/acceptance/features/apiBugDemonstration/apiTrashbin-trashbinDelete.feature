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
