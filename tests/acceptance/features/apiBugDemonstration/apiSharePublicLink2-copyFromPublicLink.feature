@api @files_sharing-app-required @public_link_share-feature-required @issue-ocis-reva-310
Feature: copying from public link share

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files
    And user "Alice" has created folder "/PARENT"
    And the administrator has enabled DAV tech_preview

  @issue-ocis-reva-373 @issue-core-37683 @skipOnOcis-OCIS-Storage
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario: Copy file within a public link folder to a file with name same as an existing folder
    Given user "Alice" has uploaded file with content "some data" to "/PARENT/testfile.txt"
    And user "Alice" has created folder "/PARENT/new-folder"
    And user "Alice" has uploaded file with content "some data 1" to "/PARENT/new-folder/testfile1.txt"
    And user "Alice" has created a public link share with settings
      | path        | /PARENT                   |
      | permissions | read,update,create,delete |
    When the public copies file "/testfile.txt" to "/new-folder" using the new public WebDAV API
    Then the HTTP status code should be "204"
    And as "Alice" file "/PARENT/testfile.txt" should exist
    And as "Alice" file "/PARENT/new-folder" should exist
    And the content of file "/PARENT/testfile.txt" for user "Alice" should be "some data"
