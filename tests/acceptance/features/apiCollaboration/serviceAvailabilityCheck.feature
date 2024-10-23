Feature: service health check


  Scenario: check service health
    When a user requests these endpoints:
      | endpoint                                        | service       |
      | %base_url_without_scheme_and_port%:9304/healthz | collaboration |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: check service readiness
    When a user requests these endpoints:
      | endpoint                                       | service       |
      | %base_url_without_scheme_and_port%:9304/readyz | collaboration |
    Then the HTTP status code of responses on all endpoints should be "200"