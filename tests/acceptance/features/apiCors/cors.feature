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
      | header                        | value                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                     |
      | Access-Control-Allow-Headers  | OC-Checksum,OC-Total-Length,OCS-APIREQUEST,X-OC-Mtime,OC-RequestAppPassword,Accept,Authorization,Brief,Content-Length,Content-Range,Content-Type,Date,Depth,Destination,Host,If,If-Match,If-Modified-Since,If-None-Match,If-Range,If-Unmodified-Since,Location,Lock-Token,Overwrite,Prefer,Range,Schedule-Reply,Timeout,User-Agent,X-Expected-Entity-Length,Accept-Language,Access-Control-Request-Method,Access-Control-Allow-Origin,Cache-Control,ETag,OC-Autorename,OC-CalDav-Import,OC-Chunked,OC-Etag,OC-FileId,OC-LazyOps,OC-Total-File-Length,Origin,X-Request-ID,X-Requested-With |
      | Access-Control-Expose-Headers | Content-Location,DAV,ETag,Link,Lock-Token,OC-ETag,OC-Checksum,OC-FileId,OC-JobStatus-Location,OC-RequestAppPassword,Vary,Webdav-Location,X-Sabre-Status                                                                                                                                                                                                                                                                                                                                                                                                                                   |
      | Access-Control-Allow-Origin   | https://aphno.badal                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       |
      | Access-Control-Allow-Methods  | GET,OPTIONS,POST,PUT,DELETE,MKCOL,PROPFIND,PATCH,PROPPATCH,REPORT                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                         |
    Examples:
      | ocs_api_version | endpoint                          | ocs-code | http-code |
      | 1               | /config                           | 100      | 200       |
      | 2               | /config                           | 200      | 200       |
      | 1               | /apps/files_sharing/api/v1/shares | 100      | 200       |
      | 2               | /apps/files_sharing/api/v1/shares | 200      | 200       |


  Scenario Outline: CORS headers should not be returned when CORS domain does not match origin header
    Given using OCS API version "<ocs_api_version>"
    When user "Alice" sends HTTP method "GET" to OCS API endpoint "<endpoint>" with headers
      | header | value               |
      | Origin | https://mero.badal  |
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
  Scenario Outline: CORS headers should be returned when an invalid password is used
    Given using OCS API version "<ocs_api_version>"
    When user "Alice" sends HTTP method "GET" to OCS API endpoint "<endpoint>" with headers using password "invalid"
      | header | value               |
      | Origin | https://aphno.badal |
    Then the OCS status code should be "997"
    And the HTTP status code should be "401"
    And the following headers should be set
      | header                        | value                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                     |
      | Access-Control-Allow-Headers  | OC-Checksum,OC-Total-Length,OCS-APIREQUEST,X-OC-Mtime,OC-RequestAppPassword,Accept,Authorization,Brief,Content-Length,Content-Range,Content-Type,Date,Depth,Destination,Host,If,If-Match,If-Modified-Since,If-None-Match,If-Range,If-Unmodified-Since,Location,Lock-Token,Overwrite,Prefer,Range,Schedule-Reply,Timeout,User-Agent,X-Expected-Entity-Length,Accept-Language,Access-Control-Request-Method,Access-Control-Allow-Origin,Cache-Control,ETag,OC-Autorename,OC-CalDav-Import,OC-Chunked,OC-Etag,OC-FileId,OC-LazyOps,OC-Total-File-Length,Origin,X-Request-ID,X-Requested-With |
      | Access-Control-Expose-Headers | Content-Location,DAV,ETag,Link,Lock-Token,OC-ETag,OC-Checksum,OC-FileId,OC-JobStatus-Location,OC-RequestAppPassword,Vary,Webdav-Location,X-Sabre-Status                                                                                                                                                                                                                                                                                                                                                                                                                                   |
      | Access-Control-Allow-Origin   | https://aphno.badal                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                       |
      | Access-Control-Allow-Methods  | GET,OPTIONS,POST,PUT,DELETE,MKCOL,PROPFIND,PATCH,PROPPATCH,REPORT                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                         |
    Examples:
      | ocs_api_version | endpoint                          |
      | 1               | /apps/files_sharing/api/v1/shares |
      | 2               | /apps/files_sharing/api/v1/shares |
