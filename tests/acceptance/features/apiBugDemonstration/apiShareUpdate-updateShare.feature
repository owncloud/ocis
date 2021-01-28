@api @files_sharing-app-required
Feature: sharing

  Background:
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes and without skeleton files

  @skipOnOcis-EOS-Storage @issue-ocis-reva-243 @skipOnOcis-OCIS-Storage
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario: Share ownership change after moving a shared file to another share
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Brian    |
      | Carol    |
    And user "Alice" has created folder "/Alice-folder"
    And user "Alice" has created folder "/Alice-folder/folder2"
    And user "Carol" has created folder "/Carol-folder"
    And user "Alice" has shared folder "/Alice-folder" with user "Brian" with permissions "all"
    And user "Carol" has shared folder "/Carol-folder" with user "Brian" with permissions "all"
    When user "Brian" moves folder "/Alice-folder/folder2" to "/Carol-folder/folder2" using the WebDAV API
    And user "Carol" gets the info of the last share using the sharing API
    # Note: in the following fields, file_parent has been removed because OCIS does not report that
    Then the fields of the last response to user "Carol" sharing with user "Brian" should include
      | id                | A_STRING             |
      | item_type         | folder               |
      | item_source       | A_STRING             |
      | share_type        | user                 |
      | file_source       | A_STRING             |
      | file_target       | /Carol-folder        |
      | permissions       | all                  |
      | stime             | A_NUMBER             |
      | storage           | A_STRING             |
      | mail_send         | 0                    |
      | uid_owner         | %username%           |
      | displayname_owner | %displayname%        |
      | mimetype          | httpd/unix-directory |
    # Really folder2 should be gone from Alice-folder and be found in Carol-folder
    # like in these 2 suggested steps:
    # And as "Alice" folder "/Alice-folder/folder2" should not exist
    # And as "Carol" folder "/Carol-folder/folder2" should exist
    #
    # But this happens on OCIS:
    And as "Alice" folder "/Alice-folder/folder2" should exist
    And as "Carol" folder "/Carol-folder/folder2" should not exist

  @skipOnOcis-OC-Storage @issue-ocis-reva-243 @skipOnOcis-OCIS-Storage
  # same as oC10 core Scenario but without displayname_owner because EOS does not report it
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario: Share ownership change after moving a shared file to another share
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Brian    |
      | Carol    |
    And user "Alice" has created folder "/Alice-folder"
    And user "Alice" has created folder "/Alice-folder/folder2"
    And user "Carol" has created folder "/Carol-folder"
    And user "Alice" has shared folder "/Alice-folder" with user "Brian" with permissions "all"
    And user "Carol" has shared folder "/Carol-folder" with user "Brian" with permissions "all"
    When user "Brian" moves folder "/Alice-folder/folder2" to "/Carol-folder/folder2" using the WebDAV API
    And user "Carol" gets the info of the last share using the sharing API
    Then the fields of the last response to user "Carol" sharing with user "Brian" should include
      | id          | A_STRING             |
      | item_type   | folder               |
      | item_source | A_STRING             |
      | share_type  | user                 |
      | file_source | A_STRING             |
      | file_target | /Carol-folder        |
      | permissions | all                  |
      | stime       | A_NUMBER             |
      | storage     | A_STRING             |
      | mail_send   | 0                    |
      | uid_owner   | %username%           |
      | mimetype    | httpd/unix-directory |
    And as "Alice" folder "/Alice-folder/folder2" should exist
    And as "Carol" folder "/Carol-folder/folder2" should not exist
