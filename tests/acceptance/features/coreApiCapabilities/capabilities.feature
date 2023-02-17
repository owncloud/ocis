@api @files_sharing-app-required
Feature: capabilities

  Background:
    Given using OCS API version "1"

  @smokeTest @issue-ocis-1285
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

  @issue-ocis-1285
  Scenario: getting default_permissions capability with admin user
    When the administrator retrieves the capabilities using the capabilities API
    Then the OCS status code should be "100"
    And the HTTP status code should be "200"
    And the capabilities should contain
      | capability    | path_to_element     | value |
      | files_sharing | default_permissions | 31    |

  @issue-ocis-1285
  Scenario: .htaccess is reported as a blacklisted file by default
    When the administrator retrieves the capabilities using the capabilities API
    Then the OCS status code should be "100"
    And the HTTP status code should be "200"
    And the capabilities should contain
      | capability | path_to_element                | value     |
      | files      | blacklisted_files@@@element[0] | .htaccess |

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
