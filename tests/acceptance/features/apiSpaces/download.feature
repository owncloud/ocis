@api @skipOnOcV10
Feature: Download file in project space
    As a user, I want to be able to download files

    Note - this feature is run in CI with ACCOUNTS_HASH_DIFFICULTY set to the default for production
    See https://github.com/owncloud/ocis/issues/1542 and https://github.com/owncloud/ocis/pull/839

    Background:
        Given these users have been created with default attributes and without skeleton files:
            | username |
            | Alice    |
            | Brian    |
        And the administrator has given "Alice" the role "Admin" using the settings api
        And user "Alice" has created a space "download file" with the default quota using the GraphApi
        And user "Alice" has uploaded a file inside space "download file" with content "some content" to "file.txt"


    Scenario: An user downloads a file in the project space
        When user "Alice" downloads the file "file.txt" of the space "download file" using the WebDAV API
        Then the HTTP status code should be "200"
        And the following headers should be set
            | header         | value |
            | Content-Length | 12    |


    Scenario: An user downloads an old version of the file in the project space
        Given user "Alice" has uploaded a file inside space "download file" with content "new content" to "file.txt"
        And user "Alice" has uploaded a file inside space "download file" with content "newest content" to "file.txt"
        When user "Alice" downloads version of the file "file.txt" with the index "1" of the space "download file" using the WebDAV API
        Then the HTTP status code should be "200"
        And the following headers should be set
            | header         | value |
            | Content-Length | 11    |
        When user "Alice" downloads version of the file "file.txt" with the index "2" of the space "download file" using the WebDAV API
        Then the HTTP status code should be "200"
        And the following headers should be set
            | header         | value |
            | Content-Length | 12    |
