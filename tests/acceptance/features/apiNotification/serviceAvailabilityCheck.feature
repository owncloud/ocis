Feature: service health check


  Scenario: check service health
    When a user requests these endpoints with "GET" and no authentication
      | endpoint                         | service      |
      | %base_url_hostname%:9174/healthz | notification |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: check service readiness
    When a user requests these endpoints with "GET" and no authentication
      | endpoint                        | service      |
      | %base_url_hostname%:9174/readyz | notification |
    Then the HTTP status code of responses on all endpoints should be "200"
