@api @issue-1328
Feature: resharing a resource with an expiration date

  Background:
    Given these users have been created with default attributes and without skeleton files:
      | username |
      | Alice    |
      | Brian    |
    And user "Alice" has uploaded file "filesForUpload/textfile.txt" to "/textfile0.txt"


  Scenario Outline: User should be able to set expiration while resharing a file with user
    Given using OCS API version "<ocs_api_version>"
    And user "Carol" has been created with default attributes and without skeleton files
    And user "Alice" has shared file "/textfile0.txt" with user "Brian" with permissions "read,update,share"
    And user "Brian" has accepted share "/textfile0.txt" offered by user "Alice"
    When user "Brian" creates a share using the sharing API with settings
      | path        | /Shares/textfile0.txt |
      | shareType   | user                  |
      | permissions | change                |
      | shareWith   | Carol                 |
      | expireDate  | +3 days               |
    Then the HTTP status code should be "200"
    And the OCS status code should be "<ocs_status_code>"
    And user "Carol" should be able to accept pending share "/textfile0.txt" offered by user "Brian"
    And the information of the last share of user "Brian" should include
      | expiration | +3 days |
    And the response when user "Carol" gets the info of the last share should include
      | expiration | +3 days |
    Examples:
      | ocs_api_version | ocs_status_code |
      | 1               | 100             |
      | 2               | 200             |

  @issue-1289
  Scenario Outline: User should be able to set expiration while resharing a file with group
    Given using OCS API version "<ocs_api_version>"
    And user "Carol" has been created with default attributes and without skeleton files
    And group "grp1" has been created
    And user "Carol" has been added to group "grp1"
    And user "Alice" has shared file "/textfile0.txt" with user "Brian" with permissions "read,update,share"
    And user "Brian" has accepted share "/textfile0.txt" offered by user "Alice"
    When user "Brian" creates a share using the sharing API with settings
      | path        | /Shares/textfile0.txt |
      | shareType   | group                 |
      | permissions | change                |
      | shareWith   | grp1                  |
      | expireDate  | +3 days               |
    Then the HTTP status code should be "200"
    And the OCS status code should be "<ocs_status_code>"
    And user "Carol" should be able to accept pending share "/textfile0.txt" offered by user "Brian"
    And the information of the last share of user "Brian" should include
      | expiration | +3 days |
    And the response when user "Carol" gets the info of the last share should include
      | expiration | +3 days |
    Examples:
      | ocs_api_version | ocs_status_code |
      | 1               | 100             |
      | 2               | 200             |


  Scenario Outline: resharing with user using the sharing API with expire days set and combinations of default/enforce expire date enabled
    Given using OCS API version "<ocs_api_version>"
    And user "Carol" has been created with default attributes and without skeleton files
    And user "Alice" has shared file "/textfile0.txt" with user "Brian" with permissions "read,update,share"
    And user "Brian" has accepted share "/textfile0.txt" offered by user "Alice"
    When user "Brian" creates a share using the sharing API with settings
      | path        | /Shares/textfile0.txt |
      | shareType   | user                  |
      | permissions | change                |
      | shareWith   | Carol                 |
    Then the HTTP status code should be "200"
    And the OCS status code should be "<ocs_status_code>"
    And user "Carol" should be able to accept pending share "/textfile0.txt" offered by user "Brian"
    And the information of the last share of user "Brian" should include
      | expiration | <expected-expire-date> |
    And the response when user "Carol" gets the info of the last share should include
      | expiration | <expected-expire-date> |
    Examples:
      | ocs_api_version | expected-expire-date | ocs_status_code |
      | 1               | +30 days             | 100             |
      | 2               | +30 days             | 200             |
      | 1               |                      | 100             |
      | 2               |                      | 200             |

  @issue-1289
  Scenario Outline: resharing with group using the sharing API with expire days set and combinations of default/enforce expire date enabled
    Given using OCS API version "<ocs_api_version>"
    And user "Carol" has been created with default attributes and without skeleton files
    And group "grp1" has been created
    And user "Carol" has been added to group "grp1"
    And user "Alice" has shared file "/textfile0.txt" with user "Brian" with permissions "read,update,share"
    And user "Brian" has accepted share "/textfile0.txt" offered by user "Alice"
    When user "Brian" creates a share using the sharing API with settings
      | path        | /Shares/textfile0.txt |
      | shareType   | group                 |
      | permissions | change                |
      | shareWith   | grp1                  |
    Then the HTTP status code should be "200"
    And the OCS status code should be "<ocs_status_code>"
    And user "Carol" should be able to accept pending share "/textfile0.txt" offered by user "Brian"
    And the information of the last share of user "Brian" should include
      | expiration | <expected-expire-date> |
    And the response when user "Carol" gets the info of the last share should include
      | expiration | <expected-expire-date> |
    Examples:
      | ocs_api_version | expected-expire-date | ocs_status_code |
      | 1               | +30 days             | 100             |
      | 2               | +30 days             | 200             |
      | 1               |                      | 100             |
      | 2               |                      | 200             |


  Scenario Outline: resharing with user using the sharing API without expire days set and with combinations of default/enforce expire date enabled
    Given using OCS API version "<ocs_api_version>"
    And user "Carol" has been created with default attributes and without skeleton files
    And user "Alice" has shared file "/textfile0.txt" with user "Brian" with permissions "read,update,share"
    And user "Brian" has accepted share "/textfile0.txt" offered by user "Alice"
    When user "Brian" creates a share using the sharing API with settings
      | path        | /Shares/textfile0.txt |
      | shareType   | user                  |
      | permissions | change                |
      | shareWith   | Carol                 |
    Then the HTTP status code should be "200"
    And the OCS status code should be "<ocs_status_code>"
    And user "Carol" should be able to accept pending share "/textfile0.txt" offered by user "Brian"
    And the information of the last share of user "Brian" should include
      | expiration | <expected-expire-date> |
    And the response when user "Carol" gets the info of the last share should include
      | expiration | <expected-expire-date> |
    Examples:
      | ocs_api_version | expected-expire-date | ocs_status_code |
      | 1               | +7 days              | 100             |
      | 2               | +7 days              | 200             |
      | 1               |                      | 100             |
      | 2               |                      | 200             |

  @issue-1289
  Scenario Outline: resharing with group using the sharing API without expire days set and with combinations of default/enforce expire date enabled
    Given using OCS API version "<ocs_api_version>"
    And user "Carol" has been created with default attributes and without skeleton files
    And group "grp1" has been created
    And user "Carol" has been added to group "grp1"
    And user "Alice" has shared file "/textfile0.txt" with user "Brian" with permissions "read,update,share"
    And user "Brian" has accepted share "/textfile0.txt" offered by user "Alice"
    When user "Brian" creates a share using the sharing API with settings
      | path        | /Shares/textfile0.txt |
      | shareType   | group                 |
      | permissions | change                |
      | shareWith   | grp1                  |
    Then the HTTP status code should be "200"
    And the OCS status code should be "<ocs_status_code>"
    And user "Carol" should be able to accept pending share "/textfile0.txt" offered by user "Brian"
    And the information of the last share of user "Brian" should include
      | expiration | <expected-expire-date> |
    And the response when user "Carol" gets the info of the last share should include
      | expiration | <expected-expire-date> |
    Examples:
      | ocs_api_version | expected-expire-date | ocs_status_code |
      | 1               | +7 days              | 100             |
      | 2               | +7 days              | 200             |
      | 1               |                      | 100             |
      | 2               |                      | 200             |


  Scenario Outline: resharing with user using the sharing API with expire days set and with combinations of default/enforce expire date enabled and specify expire date in share
    Given using OCS API version "<ocs_api_version>"
    And user "Carol" has been created with default attributes and without skeleton files
    And user "Alice" has shared file "/textfile0.txt" with user "Brian" with permissions "read,update,share"
    And user "Brian" has accepted share "/textfile0.txt" offered by user "Alice"
    When user "Brian" creates a share using the sharing API with settings
      | path        | /Shares/textfile0.txt |
      | shareType   | user                  |
      | permissions | change                |
      | shareWith   | Carol                 |
      | expireDate  | +20 days              |
    Then the HTTP status code should be "200"
    And the OCS status code should be "<ocs_status_code>"
    And user "Carol" should be able to accept pending share "/textfile0.txt" offered by user "Brian"
    And the information of the last share of user "Brian" should include
      | expiration | +20 days |
    And the response when user "Carol" gets the info of the last share should include
      | expiration | +20 days |
    Examples:
      | ocs_api_version | ocs_status_code |
      | 1               | 100             |
      | 2               | 200             |
      | 1               | 100             |
      | 2               | 200             |

  @issue-1289
  Scenario Outline: resharing with group using the sharing API with expire days set and with combinations of default/enforce expire date enabled and specify expire date in share
    Given using OCS API version "<ocs_api_version>"
    And user "Carol" has been created with default attributes and without skeleton files
    And group "grp1" has been created
    And user "Carol" has been added to group "grp1"
    And user "Alice" has shared file "/textfile0.txt" with user "Brian" with permissions "read,update,share"
    And user "Brian" has accepted share "/textfile0.txt" offered by user "Alice"
    When user "Brian" creates a share using the sharing API with settings
      | path        | /Shares/textfile0.txt |
      | shareType   | group                 |
      | permissions | change                |
      | shareWith   | grp1                  |
      | expireDate  | +20 days              |
    Then the HTTP status code should be "200"
    And the OCS status code should be "<ocs_status_code>"
    And user "Carol" should be able to accept pending share "/textfile0.txt" offered by user "Brian"
    And the information of the last share of user "Brian" should include
      | expiration | +20 days |
    And the response when user "Carol" gets the info of the last share should include
      | expiration | +20 days |
    Examples:
      | ocs_api_version | ocs_status_code |
      | 1               | 100             |
      | 2               | 200             |
      | 1               | 100             |
      | 2               | 200             |


  Scenario Outline: Setting default expiry date and enforcement after the share is created
    Given using OCS API version "<ocs_api_version>"
    And user "Carol" has been created with default attributes and without skeleton files
    And user "Alice" has shared file "/textfile0.txt" with user "Brian" with permissions "read,update,share"
    And user "Brian" has accepted share "/textfile0.txt" offered by user "Alice"
    And user "Brian" has shared file "/Shares/textfile0.txt" with user "Carol"
    And user "Carol" has accepted share "/textfile0.txt" offered by user "Brian"
    When user "Brian" gets the info of the last share using the sharing API
    Then the HTTP status code should be "200"
    And the OCS status code should be "<ocs_status_code>"
    And the information of the last share of user "Brian" should include
      | expiration |  |
    And the response when user "Carol" gets the info of the last share should include
      | expiration |  |
    Examples:
      | ocs_api_version | ocs_status_code |
      | 1               | 100             |
      | 2               | 200             |
      | 1               | 100             |
      | 2               | 200             |

  @issue-1289
  Scenario Outline: resharing group share with user using the sharing API with default expire date set and with combinations of default/enforce expire date enabled
    Given using OCS API version "<ocs_api_version>"
    And group "grp1" has been created
    And user "Carol" has been created with default attributes and without skeleton files
    And user "Brian" has been added to group "grp1"
    And user "Alice" has shared file "/textfile0.txt" with group "grp1" with permissions "read,update,share"
    And user "Brian" has accepted share "/textfile0.txt" offered by user "Alice"
    When user "Brian" creates a share using the sharing API with settings
      | path        | /Shares/textfile0.txt |
      | shareType   | user                  |
      | permissions | change                |
      | shareWith   | Carol                 |
    Then the HTTP status code should be "200"
    And the OCS status code should be "<ocs_status_code>"
    And user "Carol" should be able to accept pending share "/textfile0.txt" offered by user "Brian"
    And the information of the last share of user "Brian" should include
      | expiration | <expected-expire-date> |
    And the response when user "Carol" gets the info of the last share should include
      | expiration | <expected-expire-date> |
    Examples:
      | ocs_api_version | expected-expire-date | ocs_status_code |
      | 1               | +30 days             | 100             |
      | 2               | +30 days             | 200             |
      | 1               |                      | 100             |
      | 2               |                      | 200             |

  @issue-1289
  Scenario Outline: resharing group share with user using the sharing API with default expire date set and specifying expiration on share and with combinations of default/enforce expire date enabled
    Given using OCS API version "<ocs_api_version>"
    And group "grp1" has been created
    And user "Carol" has been created with default attributes and without skeleton files
    And user "Brian" has been added to group "grp1"
    And user "Alice" has shared file "/textfile0.txt" with group "grp1" with permissions "read,update,share"
    And user "Brian" has accepted share "/textfile0.txt" offered by user "Alice"
    When user "Brian" creates a share using the sharing API with settings
      | path        | /Shares/textfile0.txt |
      | shareType   | user                  |
      | permissions | change                |
      | shareWith   | Carol                 |
      | expireDate  | +20 days              |
    Then the HTTP status code should be "200"
    And the OCS status code should be "<ocs_status_code>"
    And user "Carol" should be able to accept pending share "/textfile0.txt" offered by user "Brian"
    And the information of the last share of user "Brian" should include
      | expiration | +20 days |
    And the response when user "Carol" gets the info of the last share should include
      | expiration | +20 days |
    Examples:
      | ocs_api_version | ocs_status_code |
      | 1               | 100             |
      | 2               | 200             |
      | 1               | 100             |
      | 2               | 200             |


  Scenario Outline: resharing using the sharing API with default expire date set but not enforced
    Given using OCS API version "<ocs_api_version>"
    And user "Carol" has been created with default attributes and without skeleton files
    And user "Alice" has shared file "/textfile0.txt" with user "Brian" with permissions "read,update,share"
    And user "Brian" has accepted share "/textfile0.txt" offered by user "Alice"
    When user "Brian" creates a share using the sharing API with settings
      | path        | /Shares/textfile0.txt |
      | shareType   | user                  |
      | permissions | change                |
      | shareWith   | Carol                 |
    Then the HTTP status code should be "200"
    And the OCS status code should be "<ocs_status_code>"
    And user "Carol" should be able to accept pending share "/textfile0.txt" offered by user "Brian"
    And the information of the last share of user "Brian" should include
      | expiration |  |
    And the response when user "Carol" gets the info of the last share should include
      | expiration |  |
    Examples:
      | ocs_api_version | ocs_status_code |
      | 1               | 100             |
      | 2               | 200             |
      | 1               | 100             |
      | 2               | 200             |


  Scenario Outline: reshare extends the received expiry date up to the default by default
    Given using OCS API version "<ocs_api_version>"
    And user "Carol" has been created with default attributes and without skeleton files
    And user "Alice" has created a share with settings
      | path        | textfile0.txt |
      | shareType   | user          |
      | permissions | all           |
      | shareWith   | Brian         |
      | expireDate  | +20 days      |
    And user "Brian" has accepted share "/textfile0.txt" offered by user "Alice"
    When user "Brian" creates a share using the sharing API with settings
      | path        | /Shares/textfile0.txt |
      | shareType   | user                  |
      | permissions | change                |
      | shareWith   | Carol                 |
    Then the HTTP status code should be "200"
    And the OCS status code should be "<ocs_status_code>"
    And user "Carol" should be able to accept pending share "/textfile0.txt" offered by user "Brian"
    And the information of the last share of user "Alice" should include
      | expiration | +20 days |
    And the response when user "Carol" gets the info of the last share should include
      | expiration | +20 |
    Examples:
      | ocs_api_version | ocs_status_code |
      | 1               | 100             |
      | 2               | 200             |
      | 1               | 100             |
      | 2               | 200             |
      | 1               | 100             |
      | 2               | 200             |


  Scenario Outline: reshare cannot extend the received expiry date further into the future
    Given using OCS API version "<ocs_api_version>"
    And user "Carol" has been created with default attributes and without skeleton files
    And user "Alice" has created a share with settings
      | path        | textfile0.txt |
      | shareType   | user          |
      | permissions | all           |
      | shareWith   | Brian         |
      | expireDate  | +20 days      |
    And user "Brian" has accepted share "/textfile0.txt" offered by user "Alice"
    When user "Brian" creates a share using the sharing API with settings
      | path        | /Shares/textfile0.txt |
      | shareType   | user                  |
      | permissions | change                |
      | shareWith   | Carol                 |
      | expireDate  | +40 days              |
    #The action of changing the expiration date while resharing should be forbidden
    Then the HTTP status code should be "403"
    And the OCS status code should be "403"
    And the information of the last share of user "Alice" should include
      | expiration | +20 days |
    Examples:
      | ocs_api_version |
      | 1               |
      | 2               |
      | 1               |
      | 2               |


  Scenario Outline: reshare cannot extend the received expiry date past the default when the default is enforced
    Given using OCS API version "<ocs_api_version>"
    And user "Carol" has been created with default attributes and without skeleton files
    And user "Alice" has created a share with settings
      | path        | textfile0.txt |
      | shareType   | user          |
      | permissions | all           |
      | shareWith   | Brian         |
      | expireDate  | +20 days      |
    And user "Brian" has accepted share "/textfile0.txt" offered by user "Alice"
    When user "Brian" creates a share using the sharing API with settings
      | path        | /Shares/textfile0.txt |
      | shareType   | user                  |
      | permissions | change                |
      | shareWith   | Carol                 |
      | expireDate  | +40 days              |
    Then the HTTP status code should be "403"
    And the OCS status code should be "403"
    And the information of the last share of user "Alice" should include
      | expiration | +20 days |
    Examples:
      | ocs_api_version |
      | 1               |
      | 2               |
