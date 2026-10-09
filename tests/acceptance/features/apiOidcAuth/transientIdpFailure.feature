Feature: transient IdP failure
  As a user
  I want my session to survive a short identity provider outage
  So that I am not logged out when the IdP is briefly unavailable

  Background:
    Given using spaces DAV path
    And these users have been created with default attributes:
      | username |
      | Alice    |
    And user "Alice" has been set up in oCIS

  @env-config
  Scenario: valid session gets a retryable 503 while the IdP is unreachable and works again afterwards
    Given user "Alice" has created folder "/beforeOutage"
    # nothing listens on port 1 -> connection refused -> IdP treated as transiently unavailable
    And the config "PROXY_OIDC_ISSUER" has been set to "https://localhost:1" for "proxy" service
    When user "Alice" creates folder "/duringOutage" using the WebDAV API
    Then the HTTP status code should be "503"
    And the following headers should be set
      | header      | value |
      | Retry-After | 3     |
    When the administrator restores the original config
    And user "Alice" creates folder "/afterOutage" using the WebDAV API
    Then the HTTP status code should be "201"
