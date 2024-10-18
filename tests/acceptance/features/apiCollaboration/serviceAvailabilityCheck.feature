Feature: service health check


  Scenario: health check defauts service
    When a user requests these endpoints:
      | endpoint                                        | service       |
      | %base_url_without_scheme_and_port%:9304/healthz | collaboration |
    Then the HTTP status code of responses on all endpoints should be "200"


  Scenario: service ready check
    When a user requests these endpoints:
      | endpoint                                       | service       |
      | %base_url_without_scheme_and_port%:9304/readyz | collaboration |
    Then the HTTP status code of responses on all endpoints should be "200"