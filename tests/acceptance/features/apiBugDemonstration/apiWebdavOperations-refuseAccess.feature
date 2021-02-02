@api
Feature: refuse access
  As an administrator
  I want to refuse access to unauthenticated and disabled users
  So that I can secure the system

  Background:
    Given using OCS API version "1"

  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario Outline: Unauthenticated call
    Given using <dav_version> DAV path
    When an unauthenticated client connects to the dav endpoint using the WebDAV API
    Then the HTTP status code should be "401"
    And there should be no duplicate headers
    And the following headers should be set
      | header           | value                                   |
      | WWW-Authenticate | Basic realm="%base_url_without_scheme%" |
    Examples:
      | dav_version |
      | old         |
      | new         |
