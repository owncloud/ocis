@api @skipOnOcV10
Feature: Share spaces via link
    As the manager of a space
    I want to be able to share a space via public link

    Note - this feature is run in CI with ACCOUNTS_HASH_DIFFICULTY set to the default for production
    See https://github.com/owncloud/ocis/issues/1542 and https://github.com/owncloud/ocis/pull/839

    Background:
        Given these users have been created with default attributes and without skeleton files:
            | username |
            | Alice    |
            | Brian    |
        And the administrator has given "Alice" the role "Space Admin" using the settings api
        And user "Alice" has created a space "share space" with the default quota using the GraphApi
        And user "Alice" has uploaded a file inside space "share space" with content "some content" to "test.txt"


    Scenario Outline: A manager can share a space to public via link with different permissions
        When user "Alice" creates public link share of the space "share space" with settings:
            | shareType   | 3             |
            | permissions | <permissions> |
            | password    | <password>    |
            | name        | <linkName>    |
            | expireDate  | <expireDate>  |
        Then the HTTP status code should be "200"
        And the OCS status code should be "200"
        And the OCS status message should be "OK"
        And the fields of the last response to user "Alice" should include
            | item_type              | folder                |
            | mimetype               | httpd/unix-directory  |
            | file_target            | /                     |
            | path                   | /                     |
            | permissions            | <permissionsMatching> |
            | share_type             | public_link           |
            | displayname_file_owner | %displayname%         |
            | displayname_owner      | %displayname%         |
            | uid_file_owner         | %username%            |
            | uid_owner              | %username%            |
            | name                   | <linkName>            |
        And the public should be able to download file "/test.txt" from inside the last public link shared folder using the new public WebDAV API with password "123"
        And the downloaded content should be "some content"
        Examples:
            | permissions | permissionsMatching       | password | linkName | expireDate               |
            | 1           | read                      | 123      | link     | 2042-03-25T23:59:59+0100 |
            | 5           | read,create               | 123      |          | 2042-03-25T23:59:59+0100 |
            | 15          | read,update,create,delete |          | link     |                          |


    Scenario: An uploader should be abble to upload a file
        When user "Alice" creates public link share of the space "share space" with settings:
            | shareType   | 3                        |
            | permissions | 4                        |
            | password    | 123                      |
            | name        | forUpload                |
            | expireDate  | 2042-03-25T23:59:59+0100 |
        Then the HTTP status code should be "200"
        And the OCS status code should be "200"
        And the OCS status message should be "OK"
        And the fields of the last response to user "Alice" should include
            | item_type              | folder               |
            | mimetype               | httpd/unix-directory |
            | file_target            | /                    |
            | path                   | /                    |
            | permissions            | create               |
            | share_type             | public_link          |
            | displayname_file_owner | %displayname%        |
            | displayname_owner      | %displayname%        |
            | uid_file_owner         | %username%           |
            | uid_owner              | %username%           |
            | name                   | forUpload            |
        And the public should be able to upload file "lorem.txt" into the last public link shared folder using the new public WebDAV API with password "123"
        And for user "Alice" the space "share space" should contain these entries:
            | lorem.txt |


    Scenario Outline: An user without manager role cannot share a space to public via link
        Given user "Alice" has shared a space "share space" to user "Brian" with role "<role>"
        When user "Brian" creates public link share of the space "share space" with settings:
            | permissions | 1 |
        Then the HTTP status code should be "404"
        And the OCS status code should be "404"
        And the OCS status message should be "No share permission"
        Examples:
            | role   |
            | viewer |
            | editor |
