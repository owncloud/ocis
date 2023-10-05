Feature: hide or show shared resources
  As a user I want to hide and show again shared resources

  Note - this feature is run in CI with ACCOUNTS_HASH_DIFFICULTY set to the default for production
  See https://github.com/owncloud/ocis/issues/1542 and https://github.com/owncloud/ocis/pull/839

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
    And user "Alice" has created a folder "folder" in space "Alice Hansen"

  
  Scenario Outline: user hides accepted share
    Given using OCS API version "<ocs_api_version>"
    And user "Alice" has created a share inside of space "Alice Hansen" with settings:
      | path       | folder                   |
      | shareWith  | Brian                    |
      | role       | viewer                   |
    And user "Brian" has accepted share "/folder" offered by user "Alice"
    When user "Brian" hiddes share "/Shares/folder" offered by user "Alice" using the sharing API
    Then the HTTP status code should be "200"
    When user "Brian" sends HTTP method "GET" to OCS API endpoint "/apps/files_sharing/api/v1/shares?state=all&shared_with_me=true&show_hidden=true"
    Then the OCS status code should be "<ocs_status_code>"
    And the HTTP status code should be "200"
    And folder "/Shares/folder" should be included as path in the response
    When user "Brian" sends HTTP method "GET" to OCS API endpoint "/apps/files_sharing/api/v1/shares?state=all&shared_with_me=true&show_hidden=false"
    Then the OCS status code should be "<ocs_status_code>"
    And the HTTP status code should be "200"
    And folder "/Shares/folder" should not be included as path in the response
    Examples:
      | ocs_api_version | ocs_status_code |
      | 1               | 100             |
      | 2               | 200             |
    

  Scenario Outline: user hides pending share
    Given using OCS API version "<ocs_api_version>"
    And user "Alice" has created a share inside of space "Alice Hansen" with settings:
      | path       | folder                   |
      | shareWith  | Brian                    |
      | role       | viewer                   |
    When user "Brian" hiddes share "/folder" offered by user "Alice" using the sharing API
    Then the HTTP status code should be "200"
    When user "Brian" sends HTTP method "GET" to OCS API endpoint "/apps/files_sharing/api/v1/shares?state=all&shared_with_me=true&show_hidden=true"
    Then the OCS status code should be "<ocs_status_code>"
    And the HTTP status code should be "200"
    And folder "/Shares/folder" should be included as path in the response
    When user "Brian" sends HTTP method "GET" to OCS API endpoint "/apps/files_sharing/api/v1/shares?state=all&shared_with_me=true&show_hidden=false"
    Then the OCS status code should be "<ocs_status_code>"
    And the HTTP status code should be "200"
    And folder "/Shares/folder" should not be included as path in the response
    Examples:
      | ocs_api_version | ocs_status_code |
      | 1               | 100             |
      | 2               | 200             |

  
  Scenario Outline: user hides declined share
    Given using OCS API version "<ocs_api_version>"
    And user "Alice" has created a share inside of space "Alice Hansen" with settings:
      | path       | folder                   |
      | shareWith  | Brian                    |
      | role       | viewer                   |
    And user "Brian" has declined share "/folder" offered by user "Alice"
    When user "Brian" hiddes share "/folder" offered by user "Alice" using the sharing API
    Then the HTTP status code should be "200"
    When user "Brian" sends HTTP method "GET" to OCS API endpoint "/apps/files_sharing/api/v1/shares?state=all&shared_with_me=true&show_hidden=true"
    Then the OCS status code should be "<ocs_status_code>"
    And the HTTP status code should be "200"
    And folder "/Shares/folder" should be included as path in the response
    When user "Brian" sends HTTP method "GET" to OCS API endpoint "/apps/files_sharing/api/v1/shares?state=all&shared_with_me=true&show_hidden=false"
    Then the OCS status code should be "<ocs_status_code>"
    And the HTTP status code should be "200"
    And folder "/Shares/folder" should not be included as path in the response
    Examples:
      | ocs_api_version | ocs_status_code |
      | 1               | 100             |
      | 2               | 200             |


  Scenario Outline: user hides a shared folder
    Given using OCS API version "<ocs_api_version>"
    And user "Alice" has created a share inside of space "Alice Hansen" with settings:
      | path       | folder                   |
      | shareWith  | Brian                    |
      | role       | viewer                   |
    And user "Brian" has accepted share "/folder" offered by user "Alice"
    When user "Alice" hiddes shared folder "/folder" using the sharing API
    Then the HTTP status code should be "200"
    When user "Alice" sends HTTP method "GET" to OCS API endpoint "/apps/files_sharing/api/v1/shares?reshares=true&show_hidden=true"
    Then the OCS status code should be "<ocs_status_code>"
    And the HTTP status code should be "200"
    And folder "/folder" should be included as path in the response
    When user "Alice" sends HTTP method "GET" to OCS API endpoint "/apps/files_sharing/api/v1/shares?reshares=true&show_hidden=false"
    Then the OCS status code should be "<ocs_status_code>"
    And the HTTP status code should be "200"
    And folder "/folder" should not be included as path in the response
    Examples:
      | ocs_api_version | ocs_status_code |
      | 1               | 100             |
      | 2               | 200             |


  Scenario Outline: user displays the hidden accepted share
    Given using OCS API version "<ocs_api_version>"
    And user "Alice" has created a share inside of space "Alice Hansen" with settings:
      | path       | folder                   |
      | shareWith  | Brian                    |
      | role       | viewer                   |
    And user "Brian" has accepted share "/folder" offered by user "Alice"
    And user "Brian" has hidden share "/Shares/folder" offered by user "Alice"
    When user "Brian" displays share "/Shares/folder" offered by user "Alice" using the sharing API
    Then the HTTP status code should be "200"
    When user "Brian" sends HTTP method "GET" to OCS API endpoint "/apps/files_sharing/api/v1/shares?state=all&shared_with_me=true&show_hidden=false"
    Then the OCS status code should be "<ocs_status_code>"
    And the HTTP status code should be "200"
    And folder "/Shares/folder" should be included as path in the response
    Examples:
      | ocs_api_version | ocs_status_code |
      | 1               | 100             |
      | 2               | 200             |
    

  Scenario Outline: user displays the hidden pending share
    Given using OCS API version "<ocs_api_version>"
    And user "Alice" has created a share inside of space "Alice Hansen" with settings:
      | path       | folder                   |
      | shareWith  | Brian                    |
      | role       | viewer                   |
    And user "Brian" has hidden share "/folder" offered by user "Alice"
    When user "Brian" displays share "/Shares/folder" offered by user "Alice" using the sharing API
    Then the HTTP status code should be "200"
    When user "Brian" sends HTTP method "GET" to OCS API endpoint "/apps/files_sharing/api/v1/shares?state=all&shared_with_me=true&show_hidden=false"
    Then the OCS status code should be "<ocs_status_code>"
    And the HTTP status code should be "200"
    And folder "/Shares/folder" should be included as path in the response
    Examples:
      | ocs_api_version | ocs_status_code |
      | 1               | 100             |
      | 2               | 200             |

  
  Scenario Outline: user displays the hidden declined share
    Given using OCS API version "<ocs_api_version>"
    And user "Alice" has created a share inside of space "Alice Hansen" with settings:
      | path       | folder                   |
      | shareWith  | Brian                    |
      | role       | viewer                   |
    And user "Brian" has declined share "/folder" offered by user "Alice"
    And user "Brian" has hidden share "/folder" offered by user "Alice"
    When user "Brian" displays share "/Shares/folder" offered by user "Alice" using the sharing API
    Then the HTTP status code should be "200"
    When user "Brian" sends HTTP method "GET" to OCS API endpoint "/apps/files_sharing/api/v1/shares?state=all&shared_with_me=true&show_hidden=false"
    Then the OCS status code should be "<ocs_status_code>"
    And the HTTP status code should be "200"
    And folder "/Shares/folder" should be included as path in the response
    Examples:
      | ocs_api_version | ocs_status_code |
      | 1               | 100             |
      | 2               | 200             |
