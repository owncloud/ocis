@api @files_sharing-app-required
Feature: sharing

  Background:
    Given using OCS API version "1"
    And user "Alice" has been created with default attributes and skeleton files

  @skipOnOcis-EOS-Storage @toFixOnOCIS @issue-ocis-reva-243
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

  @skipOnOcis-OC-Storage @toFixOnOCIS @issue-ocis-reva-243
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
      | mimetype          | httpd/unix-directory |
    And as "Alice" folder "/Alice-folder/folder2" should exist
    And as "Carol" folder "/Carol-folder/folder2" should not exist

  @toFixOnOCIS @toFixOnOcV10 @issue-ocis-reva-350 @issue-ocis-reva-352 @issue-37653
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario Outline: API responds with a full set of parameters when owner changes the permission of a share
    Given using OCS API version "<ocs_api_version>"
    And user "Brian" has been created with default attributes and without skeleton files
    And user "Alice" has created folder "/Alice-folder"
    And user "Alice" has shared folder "/Alice-folder" with user "Brian" with permissions "read"
    When user "Alice" updates the last share using the sharing API with
      | permissions | all |
    Then the OCS status code should be "<ocs_status_code>"
    And the OCS status message should be "OK"
    And the HTTP status code should be "200"
    Then the fields of the last response to user "Alice" sharing with user "Brian" should include
      | id                         | A_STRING             |
      | share_type                 | user                 |
      | uid_owner                  | %username%           |
      | displayname_owner          | %displayname%        |
      | permissions                | all                  |
      | stime                      | A_NUMBER             |
      | parent                     |                      |
      | expiration                 |                      |
      | token                      |                      |
      | uid_file_owner             | %username%           |
      | displayname_file_owner     | %displayname%        |
      | additional_info_owner      |                      |
      | additional_info_file_owner |                      |
      | state                      | 0                    |
      | item_type                  | folder               |
      | item_source                | A_STRING             |
      | path                       | /Alice-folder        |
      | mimetype                   | httpd/unix-directory |
      | storage_id                 | A_STRING             |
      | storage                    | 0                    |
      | file_source                | A_STRING             |
      | file_target                | /Alice-folder        |
      | share_with                 | %username%           |
      | share_with_displayname     | %displayname%        |
      | share_with_additional_info |                      |
      | mail_send                  | 0                    |
      | name                       |                      |
    And the fields of the last response should not include
      | attributes |  |
#      | token      |  |
#      | name       |  |
    Examples:
      | ocs_api_version | ocs_status_code |
      | 1               | 100             |
      | 2               | 200             |
