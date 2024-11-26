Feature: service health check


  Scenario: check service health
    When a user requests these endpoints with "GET" and no authentication
      | endpoint                              | service       |
      | %collaboration_hostname%:9304/healthz | collaboration |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: check service readiness
    When a user requests these endpoints with "GET" and no authentication
      | endpoint                             | service       |
      | %collaboration_hostname%:9304/readyz | collaboration |
    Then the HTTP status code of responses on all endpoints should be "200"
