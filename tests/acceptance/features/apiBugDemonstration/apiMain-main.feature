@api
Feature: Other tests related to api

  @issue-ocis-reva-100
  # after fixing all issues delete this Scenario and use the one from oC10 core
  Scenario: robots.txt file should be accessible
    When a user requests "/robots.txt" with "GET" and no authentication
    Then the HTTP status code should be "401" or "404"
