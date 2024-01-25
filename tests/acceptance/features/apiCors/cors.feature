# NOTE: set env OCIS_CORS_ALLOW_ORIGINS=https://aphno.badal while running ocis server
@env-config
Feature: CORS headers
  As a user
  I want to send a cross-origin request
  So that I can check if the correct headers are set

  Background:
    Given user "Alice" has been created with default attributes and without skeleton files
  And the config "OCIS_CORS_ALLOW_ORIGINS" has been set to "https://aphno.badal"

  @issue-5195
  Scenario Outline: CORS headers should be returned when setting CORS domain sending origin header
    Given using OCS API version "<ocs_api_version>"
    When user "Alice" sends HTTP method "GET" to OCS API endpoint "<endpoint>" with headers
      | header | value               |
      | Origin | https://aphno.badal |
    Then the OCS status code should be "<ocs-code>"
    And the HTTP status code should be "<http-code>"
    And the following headers should be set
      | header                           | value               |
      | Access-Control-Expose-Headers    | Location            |
      | Access-Control-Allow-Origin      | https://aphno.badal |
      | Access-Control-Allow-Credentials | true                |
    Examples:
      | ocs_api_version | endpoint                          | ocs-code | http-code |
      | 1               | /config                           | 100      | 200       |
      | 2               | /config                           | 200      | 200       |
      | 1               | /apps/files_sharing/api/v1/shares | 100      | 200       |
      | 2               | /apps/files_sharing/api/v1/shares | 200      | 200       |


  Scenario Outline: CORS headers should not be returned when CORS domain does not match origin header
    Given using OCS API version "<ocs_api_version>"
    When user "Alice" sends HTTP method "GET" to OCS API endpoint "<endpoint>" with headers
      | header | value              |
      | Origin | https://mero.badal |
    Then the OCS status code should be "<ocs-code>"
    And the HTTP status code should be "<http-code>"
    And the following headers should not be set
      | header                        |
      | Access-Control-Allow-Headers  |
      | Access-Control-Expose-Headers |
      | Access-Control-Allow-Origin   |
      | Access-Control-Allow-Methods  |
    Examples:
      | ocs_api_version | endpoint                          | ocs-code | http-code |
      | 1               | /config                           | 100      | 200       |
      | 2               | /config                           | 200      | 200       |
      | 1               | /apps/files_sharing/api/v1/shares | 100      | 200       |
      | 2               | /apps/files_sharing/api/v1/shares | 200      | 200       |

  @issue-5194
  Scenario Outline: CORS headers should be returned when an preflight request is sent
    Given using OCS API version "<ocs_api_version>"
    When user "Alice" sends HTTP method "OPTIONS" to OCS API endpoint "<endpoint>" with headers
      | header                         | value                                                                                                                                                                                                                                                                                                                                                 |
      | Origin                         | https://aphno.badal                                                                                                                                                                                                                                                                                                                                   |
      | Access-Control-Request-Headers | Origin, Accept, Content-Type, Depth, Authorization, Ocs-Apirequest, If-None-Match, If-Match, Destination, Overwrite, X-Request-Id, X-Requested-With, Tus-Resumable, Tus-Checksum-Algorithm, Upload-Concat, Upload-Length, Upload-Metadata, Upload-Defer-Length, Upload-Expires, Upload-Checksum, Upload-Offset, X-Http-Method-Override, Cache-Control |
      | Access-Control-Request-Method  | <request_method>                                                                                                                                                                                                                                                                                                                                      |
    And the HTTP status code should be "204"
    And the following headers should be set
      | header                       | value                                                                                                                                                                                                                                                                                                                                                 |
      | Access-Control-Allow-Headers | Origin, Accept, Content-Type, Depth, Authorization, Ocs-Apirequest, If-None-Match, If-Match, Destination, Overwrite, X-Request-Id, X-Requested-With, Tus-Resumable, Tus-Checksum-Algorithm, Upload-Concat, Upload-Length, Upload-Metadata, Upload-Defer-Length, Upload-Expires, Upload-Checksum, Upload-Offset, X-Http-Method-Override, Cache-Control |
      | Access-Control-Allow-Origin  | https://aphno.badal                                                                                                                                                                                                                                                                                                                                   |
      | Access-Control-Allow-Methods | <request_method>                                                                                                                                                                                                                                                                                                                                      |
    Examples:
      | ocs_api_version |  | endpoint                          | request_method |
      | 1               |  | /apps/files_sharing/api/v1/shares | GET            |
      | 2               |  | /apps/files_sharing/api/v1/shares | PUT            |
      | 1               |  | /apps/files_sharing/api/v1/shares | DELETE         |
      | 2               |  | /apps/files_sharing/api/v1/shares | POST           |


  Scenario: CORS headers should be returned when setting CORS domain sending origin header in the Graph api
    When user "Alice" lists all available spaces with headers using the Graph API
      | header | value               |
      | Origin | https://aphno.badal |
    Then the HTTP status code should be "200"
    And the following headers should be set
      | header                      | value               |
      | Access-Control-Allow-Origin | https://aphno.badal |

  @issue-8231
  Scenario: CORS headers should be returned when setting CORS domain sending origin header in the Webdav api
    When user "Alice" sends PROPFIND request to space "Alice Hansen" with headers using the WebDAV API
      | header | value               |
      | Origin | https://aphno.badal |
    Then the HTTP status code should be "207"
    And the following headers should be set
      | header                      | value               |
      | Access-Control-Allow-Origin | https://aphno.badal |


  Scenario: CORS headers should be returned when setting CORS domain sending origin header in the settings api
    When user "Alice" lists values-list with headers using the Settings API
      | header | value               |
      | Origin | https://aphno.badal |
    Then the HTTP status code should be "201"
    And the following headers should be set
      | header                      | value               |
      | Access-Control-Allow-Origin | https://aphno.badal |
      