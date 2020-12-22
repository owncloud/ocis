@api @files_sharing-app-required @public_link_share-feature-required @skipOnOcis-EOS-Storage @issue-ocis-reva-315 @issue-ocis-reva-316

Feature: upload to a public link share

  Background:
    Given user "Alice" has been created with default attributes and skeleton files

  @issue-ocis-reva-290
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario: Uploading file to a public upload-only share that was deleted does not work
    Given the administrator has enabled DAV tech_preview
    And user "Alice" has created a public link share with settings
      | path        | FOLDER |
      | permissions | create |
    When user "Alice" deletes file "/FOLDER" using the WebDAV API
    And the public uploads file "does-not-matter.txt" with content "does not matter" using the new public WebDAV API
    Then the HTTP status code should be "403"
    # actually it should be 404
