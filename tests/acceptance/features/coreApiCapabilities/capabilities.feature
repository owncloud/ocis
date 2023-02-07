@api @files_sharing-app-required @issue-ocis-reva-41
Feature: capabilities

  Background:
    Given using OCS API version "1"

  @smokeTest @skipOnOcis
  Scenario: getting new default capabilities in versions after 10.5.0 with admin user
    When the administrator retrieves the capabilities using the capabilities API
    Then the OCS status code should be "100"
    And the HTTP status code should be "200"
    And the capabilities should contain
      | capability | path_to_element                 | value |
      | files      | favorites                       | 1     |
      | files      | file_locking_support            | 1     |
      | files      | file_locking_enable_file_action | EMPTY |

  @smokeTest @skipOnOcis
  Scenario: lock file action can be enabled
    Given parameter "enable_lock_file_action" of app "files" has been set to "yes"
    When the administrator retrieves the capabilities using the capabilities API
    Then the OCS status code should be "100"
    And the HTTP status code should be "200"
    And the capabilities should contain
      | capability | path_to_element                 | value |
      | files      | file_locking_support            | 1     |
      | files      | file_locking_enable_file_action | 1     |

  @smokeTest @skipOnOcis
  Scenario: getting default capabilities with admin user
    When the administrator retrieves the capabilities using the capabilities API
    Then the OCS status code should be "100"
    And the HTTP status code should be "200"
    And the capabilities should contain
      | capability    | path_to_element        | value |
      | files_sharing | user@@@profile_picture | 1     |

  @files_trashbin-app-required @skipOnReva
  Scenario: getting trashbin app capability with admin user
    When the administrator retrieves the capabilities using the capabilities API
    Then the OCS status code should be "100"
    And the HTTP status code should be "200"
    And the capabilities should contain
      | capability | path_to_element | value |
      | files      | undelete        | 1     |

  @files_versions-app-required @skipOnReva
  Scenario: getting versions app capability with admin user
    When the administrator retrieves the capabilities using the capabilities API
    Then the OCS status code should be "100"
    And the HTTP status code should be "200"
    And the capabilities should contain
      | capability | path_to_element | value |
      | files      | versioning      | 1     |

  @skipOnOcis
  Scenario: getting default_permissions capability with admin user
    When the administrator retrieves the capabilities using the capabilities API
    Then the OCS status code should be "100"
    And the HTTP status code should be "200"
    And the capabilities should contain
      | capability    | path_to_element     | value |
      | files_sharing | default_permissions | 31    |

  @skipOnOcis
  Scenario: .htaccess is reported as a blacklisted file by default
    When the administrator retrieves the capabilities using the capabilities API
    Then the OCS status code should be "100"
    And the HTTP status code should be "200"
    And the capabilities should contain
      | capability | path_to_element                | value     |
      | files      | blacklisted_files@@@element[0] | .htaccess |

  @skipOnOcis
  Scenario: multiple files can be reported as blacklisted
    Given the administrator has updated system config key "blacklisted_files" with value '["test.txt",".htaccess"]' and type "json"
    When the administrator retrieves the capabilities using the capabilities API
    Then the OCS status code should be "100"
    And the HTTP status code should be "200"
    And the capabilities should contain
      | capability | path_to_element                | value     |
      | files      | blacklisted_files@@@element[0] | test.txt  |
      | files      | blacklisted_files@@@element[1] | .htaccess |

  #feature added in #31824 released in 10.0.10
  @smokeTest @skipOnOcis
  Scenario: getting capabilities with admin user
    When the administrator retrieves the capabilities using the capabilities API
    Then the OCS status code should be "100"
    And the HTTP status code should be "200"
    And the capabilities should contain
      | capability    | path_to_element | value |
      | files_sharing | can_share       | 1     |

  #feature added in #32414 released in 10.0.10
  @skipOnOcis
  Scenario: getting async capabilities when async operations are enabled
    Given the administrator has enabled async operations
    When the administrator retrieves the capabilities using the capabilities API
    Then the OCS status code should be "100"
    And the HTTP status code should be "200"
    And the capabilities should contain
      | capability | path_to_element | value |
      | async      |                 | 1.0   |


  Scenario: getting async capabilities when async operations are disabled
    Given the administrator has disabled async operations
    When the administrator retrieves the capabilities using the capabilities API
    Then the capabilities should contain
      | capability | path_to_element | value |
      | async      |                 | EMPTY |

  @skipOnOcis
  Scenario: blacklisted_files_regex is reported in capabilities
    When the administrator retrieves the capabilities using the capabilities API
    Then the OCS status code should be "100"
    And the HTTP status code should be "200"
    And the capabilities should contain
      | capability | path_to_element         | value    |
      | files      | blacklisted_files_regex | \.(part\|filepart)$ |

  @smokeTest
  Scenario: getting default capabilities with admin user
    When the administrator retrieves the capabilities using the capabilities API
    Then the capabilities should contain
      | capability    | path_to_element                           | value             |
      | core          | status@@@edition                          | %edition%         |
      | core          | status@@@product                          | %productname%     |
      | core          | status@@@productname                      | %productname%     |
      | core          | status@@@version                          | %version%         |
      | core          | status@@@versionstring                    | %versionstring%   |
    And the version data in the response should contain
      | name    | value             |
      | string  | %versionstring%   |
      | edition | %edition%         |
      | product | %productname%     |
    And the major-minor-micro version data in the response should match the version string
