Feature: service health check


  Scenario: check service health
    When a user requests these URLs with "GET" and no authentication
      | endpoint                                | service   |
      | http://%base_url_hostname%:9277/healthz | antivirus |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: check service readiness
    When a user requests these URLs with "GET" and no authentication
      | endpoint                               | service   |
      | http://%base_url_hostname%:9277/readyz | antivirus |
    Then the HTTP status code of responses on all endpoints should be "200"
